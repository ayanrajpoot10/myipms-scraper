package main

import (
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// IntRange represents a range of integers
type IntRange struct {
	From, To int
}

// String returns the string representation of IntRange
func (ir IntRange) String() string {
	if ir.From == 0 && ir.To == 0 {
		return ""
	}
	return fmt.Sprintf("%d-%d", ir.From, ir.To)
}

// Set parses a string in format "from-to" and sets the IntRange
func (ir *IntRange) Set(value string) error {
	if value == "" {
		ir.From = 0
		ir.To = 0
		return nil
	}

	parts := strings.Split(value, "-")
	if len(parts) != 2 {
		return fmt.Errorf("invalid range format, expected 'from-to' (e.g., '10-20')")
	}

	from, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return fmt.Errorf("invalid 'from' value: %v", err)
	}

	to, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return fmt.Errorf("invalid 'to' value: %v", err)
	}

	if from <= 0 || to <= 0 {
		return fmt.Errorf("range values must be positive integers")
	}

	if from > to {
		return fmt.Errorf("'from' value (%d) cannot be greater than 'to' value (%d)", from, to)
	}

	ir.From = from
	ir.To = to
	return nil
}

// IPRange represents a range of IP addresses
type IPRange struct {
	From, To net.IP
}

// String returns the string representation of IPRange
func (ipr IPRange) String() string {
	if ipr.From == nil && ipr.To == nil {
		return ""
	}
	return fmt.Sprintf("%s-%s", ipr.From.String(), ipr.To.String())
}

// Set parses a string in format "from-to" or CIDR and sets the IPRange
func (ipr *IPRange) Set(value string) error {
	if value == "" {
		ipr.From = nil
		ipr.To = nil
		return nil
	}

	if strings.Contains(value, "/") {
		fromStr, toStr, err := parseCIDR(value)
		if err != nil {
			return err
		}
		ipr.From = net.ParseIP(fromStr)
		ipr.To = net.ParseIP(toStr)
		return nil
	}

	parts := strings.Split(value, "-")
	if len(parts) != 2 {
		return fmt.Errorf("invalid IP range format, expected 'from-to' (e.g., '104.16.0.0-104.16.255.255') or CIDR notation (e.g., '192.168.0.0/24')")
	}

	ipFrom := strings.TrimSpace(parts[0])
	ipTo := strings.TrimSpace(parts[1])

	from := net.ParseIP(ipFrom)
	if from == nil {
		return fmt.Errorf("invalid IP address format for 'from': %s", ipFrom)
	}

	to := net.ParseIP(ipTo)
	if to == nil {
		return fmt.Errorf("invalid IP address format for 'to': %s", ipTo)
	}

	ipr.From = from
	ipr.To = to
	return nil
}

// OptionError represents an error for unknown options
type OptionError struct {
	Kind  string // "dns", "host", "owner", "country"
	Input string
}

func (e OptionError) Error() string {
	return fmt.Sprintf("unknown %s '%s'", e.Kind, e.Input)
}

// Config holds all configuration options
type Config struct {
	Owner         string
	Country       string
	Host          string
	DNSRecord     string
	URLFilter     string
	RankRange     IntRange
	IPRange       IPRange
	VisitorsRange IntRange
	Output        string
	MaxPages      int
	StartPage     int
	ProxyURL      string
	ProxyUser     string
	ProxyPass     string
	Workers       int
	Delay         time.Duration
	List          bool
}

// Filter holds resolved filter information
type Filter struct {
	OwnerName    string
	OwnerID      int
	CountryCode  string
	CountryName  string
	HostName     string
	HostID       int
	DNSName      string
	DNSID        int
	URLFilter    string
	RankFrom     int
	RankTo       int
	IPFrom       string
	IPTo         string
	VisitorsFrom int
	VisitorsTo   int
}

// parseFlags parses command-line flags and returns a Config
func parseFlags() *Config {
	config := &Config{}

	flag.Usage = func() {
		showHelp()
	}

	flag.StringVar(&config.Owner, "owner", "", "Enable owner filter with specified owner name")
	flag.StringVar(&config.Country, "country", "", "Enable country filter with specified country name")
	flag.StringVar(&config.Host, "host", "", "Enable host filter with specified host name")
	flag.StringVar(&config.DNSRecord, "dns", "", "Enable DNS filter with specified DNS record")
	flag.StringVar(&config.URLFilter, "url", "", "Enable URL filter to search for domains containing specific text in URL")
	flag.Var(&config.RankRange, "rank", "Filter by popularity ranking range (format: from-to, e.g., 10-20)")
	flag.Var(&config.IPRange, "ip", "Filter by IP address range (format: from-to or CIDR, e.g., 104.16.0.0-104.16.255.255 or 192.168.0.0/24)")
	flag.Var(&config.VisitorsRange, "visitors", "Filter by visitor count range (format: from-to, e.g., 1000-20000)")
	flag.StringVar(&config.Output, "output", "domains.txt", "Output filename")
	flag.IntVar(&config.MaxPages, "pages", 0, "Maximum number of pages to scrape (0 = unlimited, default: unlimited)")
	flag.IntVar(&config.StartPage, "start", 1, "Starting page number (default: 1)")
	flag.StringVar(&config.ProxyURL, "proxy", "", "Proxy URL with optional auth (e.g., http://proxy.example.com:8080@user:pass)")
	flag.IntVar(&config.Workers, "workers", 3, "Number of concurrent workers (default: 3, max recommended: 10)")
	flag.DurationVar(&config.Delay, "delay", 500*time.Millisecond, "Delay between requests (default: 500ms)")
	flag.BoolVar(&config.List, "list", false, "Show all available options for all filters")

	flag.Parse()
	return config
}

// parseCIDR converts CIDR notation to first and last IP addresses
func parseCIDR(cidr string) (string, string, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", "", fmt.Errorf("invalid CIDR notation: %v", err)
	}

	firstIP := ipNet.IP
	lastIP := make(net.IP, len(firstIP))
	copy(lastIP, firstIP)

	for i, b := range ipNet.Mask {
		lastIP[i] |= ^b
	}

	if firstIP.To4() != nil {
		firstIP = firstIP.To4()
		lastIP = lastIP.To4()
	}

	return firstIP.String(), lastIP.String(), nil
}

// parseProxyURL parses proxy URL using net/url and returns base URL, user, and password
func parseProxyURL(proxyURL string) (string, string, string, error) {
	u, err := url.Parse(proxyURL)
	if err != nil {
		return "", "", "", fmt.Errorf("invalid proxy URL format: %v", err)
	}

	if u.Scheme != "http" && u.Scheme != "https" && u.Scheme != "socks5" {
		return "", "", "", fmt.Errorf("proxy URL must use http, https, or socks5 scheme (current: %s)", u.Scheme)
	}

	var user, pass string
	if u.User != nil {
		user = u.User.Username()
		pass, _ = u.User.Password()
	}

	baseURL := fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	if u.Path != "" && u.Path != "/" {
		baseURL += u.Path
	}

	return baseURL, user, pass, nil
}

// validateAndResolveFilters validates the input filters and resolves them to IDs
func validateAndResolveFilters(config *Config) (*Filter, error) {
	filter := &Filter{}

	if config.StartPage < 1 {
		return nil, fmt.Errorf("start page must be a positive integer (current: %d)", config.StartPage)
	}

	if config.Workers < 1 {
		return nil, fmt.Errorf("workers count must be at least 1 (current: %d)", config.Workers)
	}
	if config.Workers > 10 {
		return nil, fmt.Errorf("workers count should not exceed 10 to avoid server overload (current: %d)", config.Workers)
	}

	if config.Delay < 0 {
		return nil, fmt.Errorf("delay must be non-negative (current: %d)", config.Delay)
	}

	if config.ProxyURL != "" {
		baseURL, user, pass, err := parseProxyURL(config.ProxyURL)
		if err != nil {
			return nil, err
		}

		config.ProxyURL = baseURL
		if user != "" {
			config.ProxyUser = user
		}
		if pass != "" {
			config.ProxyPass = pass
		}
	}

	if config.DNSRecord != "" {
		var exists bool
		filter.DNSID, exists = dns[config.DNSRecord]
		if !exists {
			return nil, OptionError{Kind: "DNS record", Input: config.DNSRecord}
		}
		filter.DNSName = config.DNSRecord
	}

	if config.Host != "" {
		var exists bool
		filter.HostID, exists = hosts[config.Host]
		if !exists {
			return nil, OptionError{Kind: "host", Input: config.Host}
		}
		filter.HostName = config.Host
	}

	if config.Owner != "" {
		var exists bool
		filter.OwnerID, exists = owners[config.Owner]
		if !exists {
			return nil, OptionError{Kind: "owner", Input: config.Owner}
		}
		filter.OwnerName = config.Owner
	}

	if config.Country != "" {
		var exists bool
		filter.CountryCode, exists = countries[config.Country]
		if !exists {
			return nil, OptionError{Kind: "country", Input: config.Country}
		}
		filter.CountryName = config.Country
	}

	if config.URLFilter != "" {
		filter.URLFilter = config.URLFilter
	}

	if config.RankRange.From != 0 && config.RankRange.To != 0 {
		filter.RankFrom = config.RankRange.From
		filter.RankTo = config.RankRange.To
	}

	if config.IPRange.From != nil && config.IPRange.To != nil {
		filter.IPFrom = config.IPRange.From.String()
		filter.IPTo = config.IPRange.To.String()
	}

	if config.VisitorsRange.From != 0 && config.VisitorsRange.To != 0 {
		filter.VisitorsFrom = config.VisitorsRange.From
		filter.VisitorsTo = config.VisitorsRange.To
	}

	return filter, nil
}

// handleValidationError handles validation errors and provides suggestions
func handleValidationError(err error) {
	fmt.Printf("Error: %v\n", err)

	if Err, ok := err.(OptionError); ok {
		switch Err.Kind {
		case "DNS record":
			suggestOptions(Err.Input, dns, "DNS records")
		case "host":
			suggestOptions(Err.Input, hosts, "hosts")
		case "owner":
			suggestOptions(Err.Input, owners, "owners")
		case "country":
			suggestOptions(Err.Input, countries, "countries")
		}
	}
	os.Exit(1)
}

// suggestOptions provides generic suggestions for unknown options
func suggestOptions[T any](input string, options map[string]T, label string) {
	optionNames := make([]string, 0, len(options))
	for name := range options {
		optionNames = append(optionNames, name)
	}

	bestMatches := findBestMatches(input, optionNames, 3)
	if len(bestMatches) > 0 {
		fmt.Println("\nDid you mean one of these?")
		for i, match := range bestMatches {
			fmt.Printf("  %d. %s\n", i+1, match)
		}
	} else {
		fmt.Printf("\nNo similar %s found.\n", label)
	}

	fmt.Println("\nUse --list to see all available options.")
}
