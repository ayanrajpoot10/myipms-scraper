package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// showHelp displays the help message.
func showHelp() {
	sections := map[string][]string{
		"DESCRIPTION": {
			"Scrapes domain lists from myip.ms with various filtering options.",
		},
		"USAGE": {
			"scraper [OPTIONS]",
		},
		"FILTER OPTIONS": {
			"-country <name>    Filter by country name (e.g., \"United States\", \"United Kingdom\", \"Japan\")",
			"-owner <name>      Filter by hosting provider (e.g., \"Cloudflare, Inc\")",
			"-host <name>       Filter by specific host",
			"-dns <record>      Filter by DNS record",
			"-url <text>        Filter by URL containing specific text (e.g., wiki, blog)",
			"-rank <range>      Filter by popularity ranking range (e.g., 10-20)",
			"-ip <range>        Filter by IP address range (e.g., 104.16.0.0-104.16.255.255 or 192.168.0.0/24)",
			"-visitors <range>  Filter by visitor count range (e.g., 1000-20000)",
		},
		"OUTPUT OPTIONS": {
			"-output <file>     Output filename (default: domains.txt)",
			"-pages <num>       Max pages to scrape (0 = unlimited, default: unlimited)",
			"-start <num>       Starting page number (default: 1)",
		},
		"CONCURRENCY OPTIONS": {
			"-workers <num>     Number of concurrent workers (default: 3, max: 10)",
			"-delay <ms>        Delay between requests in milliseconds (default: 500ms)",
			"                   Use -workers=1 for sequential processing",
		},
		"PROXY OPTIONS": {
			"-proxy <url>       Proxy URL with optional authentication",
			"                   Format: protocol://host:port[@user:pass]",
			"                   Examples:",
			"                     http://proxy.com:8080",
			"                     http://proxy.com:8080@user:pass",
			"                     socks5://127.0.0.1:9050",
		},
		"OTHER": {
			"-help             Show this help message",
			"-list             Show all available options (or specific options with filter flags)",
		},
		"NOTES": {
			"• All filter options can be combined",
			"• Range filters use 'from-to' format (e.g., 10-20, 1000-5000)",
			"• Range values must be positive integers (from ≤ to)",
			"• IP ranges support both 'from-to' and CIDR (e.g., 192.168.0.0/24)",
			"• Quotes required for CIDR in shells: -ip=\"192.168.0.0/24\"",
			"• Start page allows resuming scraping from a specific page",
			"• Proxy format: protocol://host:port[@user:pass]",
			"• Supported protocols: HTTP, HTTPS, SOCKS5",
			"• Higher worker count = faster scraping but may trigger rate limits",
			"• Use delay to respect server limits",
			"• If scraping fails, change IP, or use proxy",
		},
	}

	// Print sections in order
	order := []string{
		"DESCRIPTION", "USAGE", "FILTER OPTIONS", "OUTPUT OPTIONS",
		"CONCURRENCY OPTIONS", "PROXY OPTIONS", "OTHER", "NOTES",
	}

	for _, section := range order {
		fmt.Println(section + ":")
		for _, line := range sections[section] {
			fmt.Println("  " + line)
		}
		fmt.Println()
	}
}

// Helper function to display a category of options with generic values
func displayCategory[T any](title string, items map[string]T, showTotal bool) {
	fmt.Printf("%s:\n", title)
	fmt.Printf("%s\n", strings.Repeat("-", len(title)+1))

	names := make([]string, 0, len(items))
	for name := range items {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		fmt.Printf("  %s\n", name)
	}

	if showTotal {
		fmt.Printf("\nTotal: %d %s\n\n", len(items), strings.ToLower(title))
	} else {
		fmt.Println()
	}
}

// showSpecificOptions displays filter options based on provided flags
// If no flags are provided, shows all available options
func showSpecificOptions(owner, country, host, dnsRecord string) {
	hasAnyFlag := owner != "" || country != "" || host != "" || dnsRecord != ""

	// Show countries if requested specifically or if showing all options
	if country != "" || !hasAnyFlag {
		displayCategory("COUNTRIES", countries, hasAnyFlag)
	}

	// Show owners if requested specifically or if showing all options
	if owner != "" || !hasAnyFlag {
		displayCategory("OWNERS/HOSTING PROVIDERS", owners, hasAnyFlag)
	}

	// Show hosts if requested specifically or if showing all options
	if host != "" || !hasAnyFlag {
		displayCategory("HOSTS", hosts, hasAnyFlag)
	}

	// Show DNS records if requested specifically or if showing all options
	if dnsRecord != "" || !hasAnyFlag {
		displayCategory("DNS RECORDS", dns, hasAnyFlag)
	}

	// Show grand total only when displaying all options
	if !hasAnyFlag {
		fmt.Printf("Total: %d countries, %d owners, %d hosts, %d DNS records\n",
			len(countries), len(owners), len(hosts), len(dns))
	}
}

// displayScrapingFilter shows the current scraping configuration
func displayScrapingFilter(filter *Filter, config *Config) {
	fmt.Printf("Filter: ")
	if filter.DNSName != "" {
		fmt.Printf("DNS (%s - ID: %d)", filter.DNSName, filter.DNSID)
	} else if filter.HostName != "" {
		fmt.Printf("Host (%s - ID: %d)", filter.HostName, filter.HostID)
	} else if filter.CountryCode == "" && filter.OwnerName == "" && filter.URLFilter == "" && filter.RankFrom == 0 && filter.IPFrom == "" && filter.VisitorsFrom == 0 {
		fmt.Print("Top Domains (default)")
	} else {
		var filters []string
		if filter.URLFilter != "" {
			filters = append(filters, fmt.Sprintf("URL (%s)", filter.URLFilter))
		}
		if filter.CountryCode != "" {
			filters = append(filters, fmt.Sprintf("Country (%s - %s)", filter.CountryName, filter.CountryCode))
		}
		if filter.RankFrom > 0 && filter.RankTo > 0 {
			filters = append(filters, fmt.Sprintf("Rank (%d-%d)", filter.RankFrom, filter.RankTo))
		}
		if filter.IPFrom != "" && filter.IPTo != "" {
			filters = append(filters, fmt.Sprintf("IP Range (%s-%s)", filter.IPFrom, filter.IPTo))
		}
		if filter.VisitorsFrom > 0 && filter.VisitorsTo > 0 {
			filters = append(filters, fmt.Sprintf("Visitors (%d-%d)", filter.VisitorsFrom, filter.VisitorsTo))
		}
		if filter.OwnerName != "" {
			filters = append(filters, fmt.Sprintf("Owner (%s - ID: %d)", filter.OwnerName, filter.OwnerID))
		}
		fmt.Print(strings.Join(filters, " + "))
	}

	// Show proxy configuration if enabled
	proxyInfo := ""
	if config.ProxyURL != "" {
		proxyInfo = fmt.Sprintf("\nProxy: %s", config.ProxyURL)
		if config.ProxyUser != "" {
			proxyInfo += fmt.Sprintf(" (authenticated as %s)", config.ProxyUser)
		}
	}

	// Show concurrency configuration
	concurrencyInfo := fmt.Sprintf("\nConcurrency: %d workers, %dms delay", config.Workers, config.Delay/time.Millisecond)
	if config.Workers == 1 {
		concurrencyInfo = "\nMode: Sequential (single worker)"
	}

	fmt.Printf("\nOutput: %s\nPages: %s (starting from page %d)%s%s\n",
		config.Output, getPagesDisplay(config.MaxPages), config.StartPage, proxyInfo, concurrencyInfo)
}

// getPagesDisplay returns the appropriate display string for MaxPages
func getPagesDisplay(maxPages int) string {
	if maxPages == 0 {
		return "unlimited"
	}
	return fmt.Sprintf("%d", maxPages)
}
