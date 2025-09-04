package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var requestData = url.Values{
	"getpage": []string{"yes"},
	"lang":    []string{"en"},
}

var Cookies = map[string]string{
	"PHPSESSID":           "ae6doi5fo94hv5k2ouqmopd47k",
	"s2_csrf_cookie_name": "cf0b4574d2c27713afd4b26879597e5d",
	"s2_theme_ui":         "red",
	"s2_uGoo":             "w6a162dd67b1968e6349944bcff010fdd63ee724",
	"s2_uLang":            "en",
	"sh":                  "72",
	"sw":                  "95.4",
}

// HTTPClient represents an HTTP client with headers and cookies
type HTTPClient struct {
	client  *http.Client
	headers map[string]string
	cookies []*http.Cookie
}

// Scraper holds the HTTP client and URL template for efficient scraping
type Scraper struct {
	httpClient  *HTTPClient
	urlTemplate string
}

// NewScraper creates a new scraper with the HTTP client and builds the URL template
func NewScraper(httpClient *HTTPClient, filter *Filter) *Scraper {
	return &Scraper{
		httpClient:  httpClient,
		urlTemplate: buildURLTemplate(filter),
	}
}

// newHTTPClient creates a new HTTP client with default cookies and optional proxy
func newHTTPClient(proxyURL, proxyUser, proxyPass string) *HTTPClient {
	headers := map[string]string{
		"Content-Type":     "application/x-www-form-urlencoded; charset=UTF-8",
		"X-Requested-With": "XMLHttpRequest",
		"Origin":           "https://myip.ms",
		"Referer":          "https://myip.ms/browse/sites/1",
		"User-Agent":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36",
		"Accept":           "*/*",
	}

	var cookies []*http.Cookie
	for name, value := range Cookies {
		cookies = append(cookies, &http.Cookie{Name: name, Value: value})
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	if proxyURL != "" {
		proxyURLParsed, err := url.Parse(proxyURL)
		if err != nil {
			fmt.Printf("Warning: Invalid proxy URL '%s': %v\n", proxyURL, err)
		} else {
			if proxyUser != "" && proxyPass != "" {
				proxyURLParsed.User = url.UserPassword(proxyUser, proxyPass)
			}
			transport.Proxy = http.ProxyURL(proxyURLParsed)
			fmt.Printf("Using proxy: %s\n", proxyURL)
		}
	}

	return &HTTPClient{
		client: &http.Client{
			Timeout:   30 * time.Second,
			Transport: transport,
		},
		headers: headers,
		cookies: cookies,
	}
}

// post performs a POST request with the configured headers and cookies
func (hc *HTTPClient) post(url string, data url.Values) (*http.Response, error) {
	return hc.makeRequest("POST", url, data)
}

// get performs a GET request with the configured headers and cookies
func (hc *HTTPClient) get(url string) (*http.Response, error) {
	return hc.makeRequest("GET", url, nil)
}

// makeRequest performs HTTP requests with the configured headers and cookies
func (hc *HTTPClient) makeRequest(method, url string, data url.Values) (*http.Response, error) {
	var req *http.Request
	var err error

	if method == "POST" && data != nil {
		req, err = http.NewRequest("POST", url, strings.NewReader(data.Encode()))
		if err != nil {
			return nil, err
		}
	} else {
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
	}

	for key, value := range hc.headers {
		req.Header.Set(key, value)
	}

	for _, cookie := range hc.cookies {
		req.AddCookie(cookie)
	}

	return hc.client.Do(req)
}

// buildURLTemplate constructs a URL template with all filters, leaving page as placeholder
func buildURLTemplate(f *Filter) string {
	url := "https://myip.ms/ajax_table/sites/%d"

	if f.URLFilter != "" {
		url += fmt.Sprintf("/url/%s", f.URLFilter)
	}

	if f.CountryCode != "" {
		url += fmt.Sprintf("/countryID/%s", f.CountryCode)
	}

	if f.RankFrom > 0 && f.RankTo > 0 {
		url += fmt.Sprintf("/rank/%d/rankii/%d", f.RankFrom, f.RankTo)
	}

	if f.IPFrom != "" && f.IPTo != "" {
		url += fmt.Sprintf("/ipID/%s/ipIDii/%s", f.IPFrom, f.IPTo)
	}

	if f.OwnerID != 0 {
		url += fmt.Sprintf("/own/%d", f.OwnerID)
	}

	if f.HostID != 0 {
		url += fmt.Sprintf("/hostID/%d", f.HostID)
	}

	if f.DNSID != 0 {
		url += fmt.Sprintf("/dns/%d", f.DNSID)
	}

	if f.VisitorsFrom > 0 && f.VisitorsTo > 0 {
		url += fmt.Sprintf("/cntVisitors/%d/cntVisitorsii/%d", f.VisitorsFrom, f.VisitorsTo)
	}

	return url
}

// fetchPage fetches a single page and returns the domains found
func (s *Scraper) fetchPage(page int) ([]string, error) {
	reqURL := fmt.Sprintf(s.urlTemplate, page)

	resp, err := s.httpClient.post(reqURL, requestData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	htmlContent := string(body)
	domains := extractDomains(htmlContent)

	if len(domains) == 0 {
		if isCookieExpired(htmlContent) {
			return nil, fmt.Errorf("cookies expired - human verification required")
		} else if isIPLimitExceeded(htmlContent) {
			return nil, fmt.Errorf("IP limit exceeded - you have exceeded page visit limit. Please use proxy/VPN or turn airplane mode on/off if using mobile internet to change IP address")
		} else {
			fmt.Println("Response content:")
			fmt.Println(htmlContent)
		}
	}

	return domains, nil
}
