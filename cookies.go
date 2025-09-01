package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	CaptchaBaseURL   = "https://myip.ms"
	CaptchaTargetURL = CaptchaBaseURL + "/ajax_table/sites/1"
	CaptchaFinalURL  = CaptchaBaseURL + "/browse/sites/1"
)

// Default cookies as constants
var DefaultCookies = map[string]string{
	"PHPSESSID":           "w64p0gu6cdatego5pn9bc0pavd",
	"s2_csrf_cookie_name": "14a418c9c3be22cb5c9d1ecb099be45b",
	"s2_theme_ui":         "red",
	"s2_uGoo":             "bf531b53aa1246aacc4c9ca44ae7fbaae849b44f",
	"s2_uLang":            "en",
	"sh":                  "72",
	"sw":                  "95.4",
}

// CookieData represents the structure for storing cookies
type CookieData struct {
	Cookies map[string]string `json:"cookies"`
}

// CaptchaClient handles the captcha solving process
type CaptchaClient struct {
	client  *http.Client
	cookies map[string]string
}

// NewCaptchaClient creates a new captcha client with default cookies
func NewCaptchaClient() *CaptchaClient {
	client := &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return &CaptchaClient{
		client:  client,
		cookies: DefaultCookies,
	}
}

// makeRequest performs HTTP requests with cookies
func (c *CaptchaClient) makeRequest(method, url string, data url.Values) (*http.Response, error) {
	var req *http.Request
	var err error

	if method == "POST" && data != nil {
		req, err = http.NewRequest("POST", url, strings.NewReader(data.Encode()))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "*/*")

	for name, value := range c.cookies {
		req.AddCookie(&http.Cookie{Name: name, Value: value})
	}

	return c.client.Do(req)
}

// downloadCaptchaImage downloads the captcha image to a file
func (c *CaptchaClient) downloadCaptchaImage(captchaURL, filename string) error {
	resp, err := c.makeRequest("GET", captchaURL, nil)
	if err != nil {
		return fmt.Errorf("error downloading captcha image: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d when downloading captcha image", resp.StatusCode)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	return nil
}

// getUserInput prompts the user for input
func getUserInput(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

// solveCaptcha handles the complete captcha solving process
func solveCaptcha() error {
	fmt.Println("Starting automated captcha solving process...")

	client := NewCaptchaClient()

	postData := url.Values{
		"x":                    {"150"},
		"y":                    {"58"},
		"g_recaptcha_loaded":   {"no"},
		"captcha_token":        {"3030a682c435137811b1a547881bc70e93f746089e2dc466229537ad5afcf449"},
		"g_recaptcha_response": {""},
	}

	resp, err := client.makeRequest("POST", CaptchaTargetURL, postData)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	html := string(body)

	captchaToken := extractCaptchaToken(html)
	captchaURL := extractCaptchaURL(html)

	fmt.Printf("Captcha Token: %s\n", captchaToken)
	fmt.Printf("Captcha Image URL: %s\n", captchaURL)

	if captchaURL == "" || captchaToken == "" {
		return fmt.Errorf("no captcha image URL or token found")
	}

	filename := "captcha_image.png"
	err = client.downloadCaptchaImage(captchaURL, filename)
	if err != nil {
		return fmt.Errorf("error downloading captcha: %v", err)
	}

	fileInfo, err := os.Stat(filename)
	if err != nil {
		fmt.Printf("Error getting file info: %v\n", err)
	} else {
		fmt.Printf("Captcha image downloaded as: %s\n", filename)
		fmt.Printf("Image size: %d bytes\n", fileInfo.Size())
	}

	fmt.Println("\nPlease check the captcha image and enter the captcha text:")
	captchaResponse := getUserInput("Enter captcha: ")

	if captchaResponse == "" {
		return fmt.Errorf("no captcha response provided")
	}

	finalData := url.Values{
		"x":                  {"0"},
		"y":                  {"0"},
		"g_recaptcha_loaded": {"no"},
		"captcha_token":      {captchaToken},
		"p_captcha_response": {captchaResponse},
	}

	finalResp, err := client.makeRequest("POST", CaptchaFinalURL, finalData)
	if err != nil {
		return fmt.Errorf("error submitting captcha: %v", err)
	}
	defer finalResp.Body.Close()

	fmt.Printf("\nFinal Response Status: %d\n", finalResp.StatusCode)
	fmt.Printf("Final Response URL: %s\n", finalResp.Request.URL.String())

	fmt.Println("Captcha solving process completed successfully!")
	return nil
}
