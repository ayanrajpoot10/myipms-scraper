package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

//go:embed templates/captcha.html
var captchaHTML string

//go:embed static/*
var staticFiles embed.FS

const (
	CaptchaBaseURL   = "https://myip.ms"
	CaptchaTargetURL = CaptchaBaseURL + "/ajax_table/sites/1"
	CaptchaFinalURL  = CaptchaBaseURL + "/browse/sites/1"
	WebServerPort    = ":8080"
)

// CaptchaData holds the current captcha information
type CaptchaData struct {
	Token    string
	ImageURL string
	mu       sync.RWMutex
}

// CaptchaServer manages the web interface for captcha solving
type CaptchaServer struct {
	server      *http.Server
	captchaData *CaptchaData
	result      chan string
	client      *HTTPClient
}

// openBrowser opens the default web browser to the specified URL
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

// NewCaptchaServer creates a new captcha server instance
func NewCaptchaServer(client *HTTPClient) *CaptchaServer {
	return &CaptchaServer{
		captchaData: &CaptchaData{},
		result:      make(chan string, 1),
		client:      client,
	}
}

// Start starts the web server for captcha solving
func (cs *CaptchaServer) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", cs.handleIndex)

	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return fmt.Errorf("failed to create static file system: %v", err)
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	mux.HandleFunc("/captcha-image", cs.handleCaptchaImage)

	mux.HandleFunc("/submit-captcha", cs.handleSubmitCaptcha)

	mux.HandleFunc("/refresh-captcha", cs.handleRefreshCaptcha)

	cs.server = &http.Server{
		Addr:    WebServerPort,
		Handler: mux,
	}

	fmt.Printf("Starting captcha solver web interface at http://localhost%s\n", WebServerPort)

	if err := openBrowser("http://localhost" + WebServerPort); err != nil {
		fmt.Printf("Could not open browser automatically. Please manually navigate to: http://localhost%s\n", WebServerPort)
	}

	return cs.server.ListenAndServe()
}

// Stop stops the web server
func (cs *CaptchaServer) Stop() error {
	if cs.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return cs.server.Shutdown(ctx)
}

// handleIndex serves the main captcha solving interface
func (cs *CaptchaServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	tmpl, err := template.New("captcha").Parse(captchaHTML)
	if err != nil {
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleCaptchaImage serves the current captcha image
func (cs *CaptchaServer) handleCaptchaImage(w http.ResponseWriter, r *http.Request) {
	cs.captchaData.mu.RLock()
	imageURL := cs.captchaData.ImageURL
	cs.captchaData.mu.RUnlock()

	if imageURL == "" {
		http.Error(w, "No captcha image available", http.StatusNotFound)
		return
	}

	resp, err := cs.client.get(imageURL)
	if err != nil {
		http.Error(w, "Failed to fetch captcha image", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch captcha image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	if _, err := io.Copy(w, resp.Body); err != nil {
		fmt.Printf("Error serving captcha image: %v\n", err)
	}
}

// handleSubmitCaptcha handles captcha submission from the web interface
func (cs *CaptchaServer) handleSubmitCaptcha(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	captchaText := strings.TrimSpace(r.FormValue("captcha"))
	if captchaText == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "message": "Captcha text is required"}`))
		return
	}

	captchaToken := cs.GetCaptchaToken()
	if captchaToken == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success": false, "message": "No captcha token available"}`))
		return
	}

	finalData := url.Values{
		"x":                  {"0"},
		"y":                  {"0"},
		"g_recaptcha_loaded": {"no"},
		"captcha_token":      {captchaToken},
		"p_captcha_response": {captchaText},
	}

	resp, err := cs.client.post(CaptchaFinalURL, finalData)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success": false, "message": "Error validating captcha: ` + err.Error() + `"}`))
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success": false, "message": "Error reading validation response"}`))
		return
	}

	if !strings.Contains(string(body), "captcha_token") {
		select {
		case cs.result <- captchaText:
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"success": true, "message": "Captcha validated successfully"}`))
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"success": false, "message": "Server is not ready to accept captcha"}`))
		}
	} else {
		newCaptchaToken := extractCaptchaToken(string(body))
		newCaptchaURL := extractCaptchaURL(string(body))

		if newCaptchaToken != "" && newCaptchaURL != "" {
			cs.SetCaptchaData(newCaptchaToken, newCaptchaURL)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "message": "Incorrect captcha. Please try again with the new image."}`))
	}
}

// WaitForCaptcha waits for user input from the web interface
func (cs *CaptchaServer) WaitForCaptcha() (string, error) {
	select {
	case result := <-cs.result:
		return result, nil
	case <-time.After(5 * time.Minute):
		return "", fmt.Errorf("captcha input timeout")
	}
}

// SetCaptchaData updates the current captcha data
func (cs *CaptchaServer) SetCaptchaData(token, imageURL string) {
	cs.captchaData.mu.Lock()
	defer cs.captchaData.mu.Unlock()

	cs.captchaData.Token = token
	cs.captchaData.ImageURL = imageURL
}

// GetCaptchaToken returns the current captcha token
func (cs *CaptchaServer) GetCaptchaToken() string {
	cs.captchaData.mu.RLock()
	defer cs.captchaData.mu.RUnlock()

	return cs.captchaData.Token
}

// handleRefreshCaptcha handles requests to refresh captcha data
func (cs *CaptchaServer) handleRefreshCaptcha(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	postData := url.Values{
		"x":                    {"150"},
		"y":                    {"58"},
		"g_recaptcha_loaded":   {"no"},
		"captcha_token":        {""},
		"g_recaptcha_response": {""},
	}

	resp, err := cs.client.post(CaptchaTargetURL, postData)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success": false, "message": "Error fetching new captcha"}`))
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success": false, "message": "Error reading captcha response"}`))
		return
	}

	html := string(body)
	captchaToken := extractCaptchaToken(html)
	captchaURL := extractCaptchaURL(html)

	if captchaURL == "" || captchaToken == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success": false, "message": "No captcha data found in response"}`))
		return
	}

	cs.SetCaptchaData(captchaToken, captchaURL)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true, "message": "Captcha refreshed successfully"}`))
}

// solveCaptcha handles the complete captcha solving process with web interface
func solveCaptcha(client *HTTPClient) error {
	fmt.Println("Starting captcha solving process...")

	captchaServer := NewCaptchaServer(client)

	go func() {
		if err := captchaServer.Start(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting captcha server: %v\n", err)
		}
	}()

	time.Sleep(1 * time.Second)

	defer func() {
		if err := captchaServer.Stop(); err != nil {
			fmt.Printf("Error stopping captcha server: %v\n", err)
		}
	}()

	postData := url.Values{
		"x":                    {"150"},
		"y":                    {"58"},
		"g_recaptcha_loaded":   {"no"},
		"captcha_token":        {""},
		"g_recaptcha_response": {""},
	}

	resp, err := client.post(CaptchaTargetURL, postData)
	if err != nil {
		return fmt.Errorf("error making initial request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	html := string(body)
	captchaToken := extractCaptchaToken(html)
	captchaURL := extractCaptchaURL(html)

	if captchaURL == "" || captchaToken == "" {
		return fmt.Errorf("no captcha image URL or token found")
	}

	captchaServer.SetCaptchaData(captchaToken, captchaURL)

	fmt.Println("Captcha image is ready! Please solve it in your web browser.")
	fmt.Printf("If the browser didn't open automatically, navigate to: http://localhost%s\n", WebServerPort)

	captchaResponse, err := captchaServer.WaitForCaptcha()
	if err != nil {
		return fmt.Errorf("error getting captcha response: %v", err)
	}

	if captchaResponse == "" {
		return fmt.Errorf("no captcha response provided")
	}

	fmt.Println("Captcha solving process completed successfully!")
	return nil
}
