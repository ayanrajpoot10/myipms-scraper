package main

import (
	"regexp"
	"strings"
)

// extractDomains extracts domain names from HTML content
func extractDomains(html string) []string {
	re := regexp.MustCompile(`<td class=['"]row_name['"][^>]*><a[^>]*>([^<]+)</a>`)
	matches := re.FindAllStringSubmatch(html, -1)

	var domains []string
	for _, match := range matches {
		if len(match) > 1 {
			domain := strings.TrimSpace(match[1])
			if domain != "" {
				domains = append(domains, domain)
			}
		}
	}

	return domains
}

// extractCaptchaToken extracts the captcha token from HTML
func extractCaptchaToken(html string) string {
	re := regexp.MustCompile(`<input[^>]*name=['"]captcha_token['"][^>]*value=['"]([^'"]*)['"']`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// extractCaptchaURL extracts the captcha image URL from HTML
func extractCaptchaURL(html string) string {
	re := regexp.MustCompile(`<img[^>]*src=['"]([^'"]*captcha\.php[^'"]*)['"']`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		captchaPath := matches[1]
		if strings.HasPrefix(captchaPath, "/") {
			return CaptchaBaseURL + captchaPath
		}
		return captchaPath
	}
	return ""
}

// isCookieExpired checks if the response indicates expired cookies
func isCookieExpired(html string) bool {
	return strings.Contains(html, "Human Verification")
}

// isIPLimitExceeded checks if the response indicates IP limit exceeded
func isIPLimitExceeded(html string) bool {
	return strings.Contains(html, "You have exceeded page visit limit")
}
