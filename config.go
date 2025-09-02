package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// Config holds all configuration options
type Config struct {
	Owner         string
	Country       string
	Host          string
	DNSRecord     string
	URLFilter     string
	RankRange     string
	IPRange       string
	VisitorsRange string
	Output        string
	MaxPages      int
	StartPage     int
	ProxyURL      string
	ProxyUser     string
	ProxyPass     string
	Workers       int
	Delay         int
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
	flag.StringVar(&config.Country, "country", "", "Enable country filter with specified country code")
	flag.StringVar(&config.Host, "host", "", "Enable host filter with specified host name")
	flag.StringVar(&config.DNSRecord, "dns", "", "Enable DNS filter with specified DNS record")
	flag.StringVar(&config.URLFilter, "url", "", "Enable URL filter to search for domains containing specific text in URL")
	flag.StringVar(&config.RankRange, "rank", "", "Filter by popularity ranking range (format: from-to, e.g., 10-20)")
	flag.StringVar(&config.IPRange, "ip", "", "Filter by IP address range (format: from-to or CIDR, e.g., 104.16.0.0-104.16.255.255 or 192.168.0.0/24)")
	flag.StringVar(&config.VisitorsRange, "visitors", "", "Filter by visitor count range (format: from-to, e.g., 1000-20000)")
	flag.StringVar(&config.Output, "output", "domains.txt", "Output filename")
	flag.IntVar(&config.MaxPages, "pages", 0, "Maximum number of pages to scrape (0 = unlimited, default: unlimited)")
	flag.IntVar(&config.StartPage, "start", 1, "Starting page number (default: 1)")
	flag.StringVar(&config.ProxyURL, "proxy", "", "Proxy URL with optional auth (e.g., http://proxy.example.com:8080@user:pass)")
	flag.IntVar(&config.Workers, "workers", 3, "Number of concurrent workers (default: 3, max recommended: 10)")
	flag.IntVar(&config.Delay, "delay", 500, "Delay between requests in milliseconds (default: 500ms)")
	flag.BoolVar(&config.List, "list", false, "Show all available options for all filters")

	flag.Parse()
	return config
}

// parseIntRange parses a string in format "from-to" and returns the two integers
func parseIntRange(rangeStr string) (int, int, error) {
	if rangeStr == "" {
		return 0, 0, nil
	}

	parts := strings.Split(rangeStr, "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid range format, expected 'from-to' (e.g., '10-20')")
	}

	from, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid 'from' value: %v", err)
	}

	to, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid 'to' value: %v", err)
	}

	if from <= 0 || to <= 0 {
		return 0, 0, fmt.Errorf("range values must be positive integers")
	}

	if from > to {
		return 0, 0, fmt.Errorf("'from' value (%d) cannot be greater than 'to' value (%d)", from, to)
	}

	return from, to, nil
}

// parseCIDR converts CIDR notation to first and last IP addresses
func parseCIDR(cidr string) (string, string, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", "", fmt.Errorf("invalid CIDR notation: %v", err)
	}

	// Get network and broadcast addresses
	firstIP := ipNet.IP
	lastIP := make(net.IP, len(firstIP))
	copy(lastIP, firstIP)

	// Calculate the last IP by setting all host bits to 1
	for i, b := range ipNet.Mask {
		lastIP[i] |= ^b
	}

	// Handle IPv4-mapped IPv6 addresses
	if firstIP.To4() != nil {
		firstIP = firstIP.To4()
		lastIP = lastIP.To4()
	}

	return firstIP.String(), lastIP.String(), nil
}

// parseIPRange parses a string in format "from-to" or CIDR and returns the two IP addresses
func parseIPRange(rangeStr string) (string, string, error) {
	if rangeStr == "" {
		return "", "", nil
	}

	// Check if input is CIDR notation (contains '/')
	if strings.Contains(rangeStr, "/") {
		return parseCIDR(rangeStr)
	}

	// Parse as traditional range format (from-to)
	parts := strings.Split(rangeStr, "-")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid IP range format, expected 'from-to' (e.g., '104.16.0.0-104.16.255.255') or CIDR notation (e.g., '192.168.0.0/24')")
	}

	ipFrom := strings.TrimSpace(parts[0])
	ipTo := strings.TrimSpace(parts[1])

	// Validate IP address formats
	if net.ParseIP(ipFrom) == nil {
		return "", "", fmt.Errorf("invalid IP address format for 'from': %s", ipFrom)
	}

	if net.ParseIP(ipTo) == nil {
		return "", "", fmt.Errorf("invalid IP address format for 'to': %s", ipTo)
	}

	return ipFrom, ipTo, nil
}

// parseProxyURL parses proxy URL in format "http://host:port@user:pass"
// where @user:pass is optional
func parseProxyURL(proxyURL string) (string, string, string, error) {
	if proxyURL == "" {
		return "", "", "", nil
	}

	// Split by @ to separate URL from credentials
	parts := strings.Split(proxyURL, "@")

	var baseURL, user, pass string

	if len(parts) == 1 {
		// No credentials provided
		baseURL = parts[0]
	} else if len(parts) == 2 {
		// Credentials provided
		baseURL = parts[0]
		credentials := parts[1]

		// Split credentials by :
		credParts := strings.Split(credentials, ":")
		if len(credParts) == 2 {
			user = credParts[0]
			pass = credParts[1]
		} else {
			return "", "", "", fmt.Errorf("invalid proxy credentials format, expected 'user:pass' (current: %s)", credentials)
		}
	} else {
		return "", "", "", fmt.Errorf("invalid proxy URL format, expected 'protocol://host:port[@user:pass]' (current: %s)", proxyURL)
	}

	// Validate the base URL format
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") && !strings.HasPrefix(baseURL, "socks5://") {
		return "", "", "", fmt.Errorf("proxy URL must start with http://, https://, or socks5:// (current: %s)", baseURL)
	}

	return baseURL, user, pass, nil
}

// validateAndResolveFilters validates the input filters and resolves them to IDs
func validateAndResolveFilters(config *Config) (*Filter, error) {
	filter := &Filter{}

	// Validate start page
	if config.StartPage < 1 {
		return nil, fmt.Errorf("start page must be a positive integer (current: %d)", config.StartPage)
	}

	// Validate workers count
	if config.Workers < 1 {
		return nil, fmt.Errorf("workers count must be at least 1 (current: %d)", config.Workers)
	}
	if config.Workers > 10 {
		return nil, fmt.Errorf("workers count should not exceed 10 to avoid server overload (current: %d)", config.Workers)
	}

	// Validate delay
	if config.Delay < 0 {
		return nil, fmt.Errorf("delay must be non-negative (current: %d)", config.Delay)
	}

	// Parse and validate proxy URL if provided
	if config.ProxyURL != "" {
		baseURL, user, pass, err := parseProxyURL(config.ProxyURL)
		if err != nil {
			return nil, err
		}

		// Update config with parsed values
		config.ProxyURL = baseURL
		if user != "" {
			config.ProxyUser = user
		}
		if pass != "" {
			config.ProxyPass = pass
		}
	}

	// Validate DNS record
	if config.DNSRecord != "" {
		var exists bool
		filter.DNSID, exists = dns[config.DNSRecord]
		if !exists {
			return nil, fmt.Errorf("unknown DNS record '%s'", config.DNSRecord)
		}
		filter.DNSName = config.DNSRecord
	}

	// Validate host
	if config.Host != "" {
		var exists bool
		filter.HostID, exists = hosts[config.Host]
		if !exists {
			return nil, fmt.Errorf("unknown host '%s'", config.Host)
		}
		filter.HostName = config.Host
	}

	// Validate owner
	if config.Owner != "" {
		var exists bool
		filter.OwnerID, exists = owners[config.Owner]
		if !exists {
			return nil, fmt.Errorf("unknown owner '%s'", config.Owner)
		}
		filter.OwnerName = config.Owner
	}

	// Validate country
	if config.Country != "" {
		// First try to find by country code
		if countryName, exists := countries[config.Country]; exists {
			filter.CountryName = countryName
			filter.CountryCode = config.Country
		} else {
			// If not found by code, try to find by country name
			found := false
			for code, name := range countries {
				if name == config.Country {
					filter.CountryName = name
					filter.CountryCode = code
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("unknown country code '%s'", config.Country)
			}
		}
	}

	// Set URL filter (no validation needed, any string is acceptable)
	if config.URLFilter != "" {
		filter.URLFilter = config.URLFilter
	}

	// Validate ranking filters
	if config.RankRange != "" {
		rankFrom, rankTo, err := parseIntRange(config.RankRange)
		if err != nil {
			return nil, fmt.Errorf("invalid rank range: %v", err)
		}
		filter.RankFrom = rankFrom
		filter.RankTo = rankTo
	}

	// Validate IP address filters
	if config.IPRange != "" {
		ipFrom, ipTo, err := parseIPRange(config.IPRange)
		if err != nil {
			return nil, fmt.Errorf("invalid IP range: %v", err)
		}
		filter.IPFrom = ipFrom
		filter.IPTo = ipTo
	}

	// Validate visitor count filters
	if config.VisitorsRange != "" {
		visitorsFrom, visitorsTo, err := parseIntRange(config.VisitorsRange)
		if err != nil {
			return nil, fmt.Errorf("invalid visitors range: %v", err)
		}
		filter.VisitorsFrom = visitorsFrom
		filter.VisitorsTo = visitorsTo
	}

	return filter, nil
}

// handleValidationError handles validation errors with suggestions
func handleValidationError(err error, config *Config) {
	fmt.Printf("Error: %v\n", err)

	errStr := err.Error()

	// Check error type by string content
	if config.DNSRecord != "" && errStr == fmt.Sprintf("unknown DNS record '%s'", config.DNSRecord) {
		suggestDNSOptions(config.DNSRecord)
	} else if config.Host != "" && errStr == fmt.Sprintf("unknown host '%s'", config.Host) {
		suggestHostOptions(config.Host)
	} else if config.Owner != "" && errStr == fmt.Sprintf("unknown owner '%s'", config.Owner) {
		suggestOwnerOptions(config.Owner)
	} else if config.Country != "" && errStr == fmt.Sprintf("unknown country code '%s'", config.Country) {
		suggestCountryOptions(config.Country)
	}

	os.Exit(1)
}

// suggestDNSOptions suggests similar DNS options
func suggestDNSOptions(input string) {
	var dnsNames []string
	for name := range dns {
		dnsNames = append(dnsNames, name)
	}

	bestMatches := findBestMatches(input, dnsNames, 3)
	if len(bestMatches) > 0 {
		fmt.Println("\nDid you mean one of these?")
		for i, match := range bestMatches {
			fmt.Printf("  %d. %s (ID: %d)\n", i+1, match, dns[match])
		}
	} else {
		fmt.Println("\nNo similar DNS records found.")
	}

	fmt.Println("\nUse --list to see all available options.")
}

// suggestHostOptions suggests similar host options
func suggestHostOptions(input string) {
	var hostNames []string
	for name := range hosts {
		hostNames = append(hostNames, name)
	}

	bestMatches := findBestMatches(input, hostNames, 3)
	if len(bestMatches) > 0 {
		fmt.Println("\nDid you mean one of these?")
		for i, match := range bestMatches {
			fmt.Printf("  %d. %s (ID: %d)\n", i+1, match, hosts[match])
		}
	} else {
		fmt.Println("\nNo similar hosts found.")
	}

	fmt.Println("\nUse --list to see all available options.")
}

// suggestOwnerOptions suggests similar owner options
func suggestOwnerOptions(input string) {
	var ownerNames []string
	for name := range owners {
		ownerNames = append(ownerNames, name)
	}

	bestMatches := findBestMatches(input, ownerNames, 3)
	if len(bestMatches) > 0 {
		fmt.Println("\nDid you mean one of these?")
		for i, match := range bestMatches {
			fmt.Printf("  %d. %s\n", i+1, match)
		}
	} else {
		fmt.Println("\nNo similar owners found.")
	}

	fmt.Println("\nUse --list to see all available options.")
}

// suggestCountryOptions suggests similar country options
func suggestCountryOptions(input string) {
	var countryOptions []string
	for code, name := range countries {
		countryOptions = append(countryOptions, code)
		countryOptions = append(countryOptions, name)
	}

	bestMatches := findBestMatches(input, countryOptions, 5)
	if len(bestMatches) > 0 {
		fmt.Println("\nDid you mean one of these?")
		matchCount := 0
		for _, match := range bestMatches {
			if countryName, isCode := countries[match]; isCode {
				fmt.Printf("  %d. %s\n", matchCount+1, countryName)
				matchCount++
			} else {
				for _, name := range countries {
					if name == match {
						fmt.Printf("  %d. %s\n", matchCount+1, name)
						matchCount++
						break
					}
				}
			}
			if matchCount >= 3 {
				break
			}
		}
	} else {
		fmt.Println("\nNo similar countries found.")
	}

	fmt.Println("\nUse --list to see all available options.")
}
