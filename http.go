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
	}

	// Use default cookies constant
	var cookies []*http.Cookie
	for name, value := range DefaultCookies {
		cookies = append(cookies, &http.Cookie{Name: name, Value: value})
	}

	// Create transport with optional proxy support
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// Configure proxy if provided
	if proxyURL != "" {
		proxyURLParsed, err := url.Parse(proxyURL)
		if err != nil {
			fmt.Printf("Warning: Invalid proxy URL '%s': %v\n", proxyURL, err)
		} else {
			// Set proxy authentication if provided
			if proxyUser != "" && proxyPass != "" {
				proxyURLParsed.User = url.UserPassword(proxyUser, proxyPass)
			}
			transport.Proxy = http.ProxyURL(proxyURLParsed)
			fmt.Printf("Using proxy: %s\n", proxyURL)
		}
	}

	return &HTTPClient{
		client: &http.Client{
			Timeout:   60 * time.Second,
			Transport: transport,
		},
		headers: headers,
		cookies: cookies,
	}
}

// post performs a POST request with the configured headers and cookies
func (hc *HTTPClient) post(url string, data url.Values) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
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
func buildURLTemplate(filter *Filter) string {
	url := "https://myip.ms/ajax_table/sites/%d"

	// Order: url -> countryID -> rank/rankii -> ipID/ipIDii -> own -> hostID -> dns -> cntVisitors/cntVisitorsii

	if filter.URLFilter != "" {
		url += fmt.Sprintf("/url/%s", filter.URLFilter)
	}

	if filter.CountryCode != "" {
		url += fmt.Sprintf("/countryID/%s", filter.CountryCode)
	}

	if filter.RankFrom > 0 && filter.RankTo > 0 {
		url += fmt.Sprintf("/rank/%d/rankii/%d", filter.RankFrom, filter.RankTo)
	}

	if filter.IPFrom != "" && filter.IPTo != "" {
		url += fmt.Sprintf("/ipID/%s/ipIDii/%s", filter.IPFrom, filter.IPTo)
	}

	if filter.OwnerID != 0 {
		url += fmt.Sprintf("/own/%d", filter.OwnerID)
	}

	if filter.HostID != 0 {
		url += fmt.Sprintf("/hostID/%d", filter.HostID)
	}

	if filter.DNSID != 0 {
		url += fmt.Sprintf("/dns/%d", filter.DNSID)
	}

	if filter.VisitorsFrom > 0 && filter.VisitorsTo > 0 {
		url += fmt.Sprintf("/cntVisitors/%d/cntVisitorsii/%d", filter.VisitorsFrom, filter.VisitorsTo)
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

	if resp.StatusCode != 200 {
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
			fmt.Println(htmlContent)
		}
	}

	return domains, nil
}

// UpdateHTTPClient updates the HTTP client (useful when cookies are refreshed)
func (s *Scraper) UpdateHTTPClient(httpClient *HTTPClient) {
	s.httpClient = httpClient
}
