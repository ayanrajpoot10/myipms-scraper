package main

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
)

const (
	CaptchaBaseURL   = "https://myip.ms"
	CaptchaTargetURL = CaptchaBaseURL + "/ajax_table/sites/1"
	CaptchaFinalURL  = CaptchaBaseURL + "/browse/sites/1"
)

// getUserInput prompts the user for input
func getUserInput(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

// solveCaptcha handles the complete captcha solving process with retry logic
func solveCaptcha(client *HTTPClient) error {
	fmt.Println("Starting captcha solving process...")

	postData := url.Values{
		"x":                    {"150"},
		"y":                    {"58"},
		"g_recaptcha_loaded":   {"no"},
		"captcha_token":        {""},
		"g_recaptcha_response": {""},
	}

	resp, err := client.post(CaptchaTargetURL, postData)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	html := string(body)

	for {
		captchaToken := extractCaptchaToken(html)
		captchaURL := extractCaptchaURL(html)

		if captchaURL == "" || captchaToken == "" {
			return fmt.Errorf("no captcha image URL or token found")
		}

		filename := "captcha_image.png"
		if err = client.downloadImage(captchaURL, filename); err != nil {
			return fmt.Errorf("error downloading captcha: %v", err)
		}

		if fileInfo, err := os.Stat(filename); err != nil {
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

		resp, err := client.post(CaptchaFinalURL, finalData)
		if err != nil {
			return fmt.Errorf("error submitting captcha: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading final response: %v", err)
		}

		if !strings.Contains(string(body), "captcha_token") {
			fmt.Println("Captcha solving process completed successfully!")
			return nil
		}

		fmt.Println("Captcha verification failed.")
		html = string(body)
		fmt.Println("Retrying...")
	}
}
