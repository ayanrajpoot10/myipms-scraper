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
func (r IntRange) String() string {
	if r.From == 0 && r.To == 0 {
		return ""
	}
	return fmt.Sprintf("%d-%d", r.From, r.To)
}

// Set parses a string in format "from-to" and sets the IntRange
func (r *IntRange) Set(s string) error {
	if s == "" {
		*r = IntRange{}
		return nil
	}

	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return fmt.Errorf("invalid range format, expected 'from-to' (e.g., '10-20')")
	}

	from, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
	to, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err1 != nil || err2 != nil || from <= 0 || to <= 0 || from > to {
		return fmt.Errorf("invalid range: %q", s)
	}

	r.From, r.To = from, to
	return nil
}

// IPRange represents a range of IP addresses
type IPRange struct {
	From, To net.IP
}

// String returns the string representation of IPRange
func (r IPRange) String() string {
	if r.From == nil && r.To == nil {
		return ""
	}
	return fmt.Sprintf("%s-%s", r.From, r.To)
}

// Set parses a string in format "from-to" or CIDR and sets the IPRange
func (r *IPRange) Set(s string) error {
	if s == "" {
		*r = IPRange{}
		return nil
	}

	if strings.Contains(s, "/") {
		from, to, err := parseCIDR(s)
		if err != nil {
			return err
		}
		r.From, r.To = net.ParseIP(from), net.ParseIP(to)
		return nil
	}

	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return fmt.Errorf("invalid IP range or CIDR: %q", s)
	}

	from, to := net.ParseIP(strings.TrimSpace(parts[0])), net.ParseIP(strings.TrimSpace(parts[1]))
	if from == nil || to == nil {
		return fmt.Errorf("invalid IP range: %q", s)
	}

	r.From, r.To = from, to
	return nil
}

// OptionError represents an error for unknown options
type OptionError struct {
	Kind  string
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

	flag.StringVar(&config.Owner, "owner", "", "Owner filter")
	flag.StringVar(&config.Country, "country", "", "Country filter")
	flag.StringVar(&config.Host, "host", "", "Host filter")
	flag.StringVar(&config.DNSRecord, "dns", "", "DNS filter")
	flag.StringVar(&config.URLFilter, "url", "", "URL filter substring")
	flag.Var(&config.RankRange, "rank", "Popularity ranking range (from-to)")
	flag.Var(&config.IPRange, "ip", "IP range (from-to or CIDR)")
	flag.Var(&config.VisitorsRange, "visitors", "Visitors range (from-to)")
	flag.StringVar(&config.Output, "output", "domains.txt", "Output file")
	flag.IntVar(&config.MaxPages, "pages", 0, "Max pages (0=unlimited)")
	flag.IntVar(&config.StartPage, "start", 1, "Starting page")
	flag.StringVar(&config.ProxyURL, "proxy", "", "Proxy URL (http[s]/socks5://user:pass@host:port)")
	flag.IntVar(&config.Workers, "workers", 3, "Concurrent workers (1-10)")
	flag.DurationVar(&config.Delay, "delay", 500*time.Millisecond, "Delay between requests")
	flag.BoolVar(&config.List, "list", false, "List all available options")

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
			return nil, OptionError{Kind: "DNS", Input: config.DNSRecord}
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
		case "DNS":
			suggestOptions(Err.Input, dns, "DNS")
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
