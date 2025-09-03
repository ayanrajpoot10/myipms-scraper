package main

import (
	"fmt"
	"sort"
	"strings"
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
			"-list             Show all available options (or specific options with f flags)",
		},
		"NOTES": {
			"• All f options can be combined",
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

// showSpecificOptions displays f options based on provided flags
// If no flags are provided, shows all available options
func showSpecificOptions(owner, country, host, dnsRecord string) {
	hasAnyFlag := owner != "" || country != "" || host != "" || dnsRecord != ""

	if country != "" || !hasAnyFlag {
		displayCategory("COUNTRIES", countries, hasAnyFlag)
	}

	if owner != "" || !hasAnyFlag {
		displayCategory("OWNERS/HOSTING PROVIDERS", owners, hasAnyFlag)
	}

	if host != "" || !hasAnyFlag {
		displayCategory("HOSTS", hosts, hasAnyFlag)
	}

	if dnsRecord != "" || !hasAnyFlag {
		displayCategory("DNS RECORDS", dns, hasAnyFlag)
	}

	if !hasAnyFlag {
		fmt.Printf("Total: %d countries, %d owners, %d hosts, %d DNS records\n",
			len(countries), len(owners), len(hosts), len(dns))
	}
}

// displayScrapingFilter shows the current scraping configuration
func displayScrapingFilter(f *Filter, c *Config) {
	fmt.Printf("Filter: ")
	if f.DNSName != "" {
		fmt.Printf("DNS (%s - ID: %d)", f.DNSName, f.DNSID)
	} else if f.HostName != "" {
		fmt.Printf("Host (%s - ID: %d)", f.HostName, f.HostID)
	} else if f.CountryCode == "" && f.OwnerName == "" && f.URLFilter == "" && f.RankFrom == 0 && f.IPFrom == "" && f.VisitorsFrom == 0 {
		fmt.Print("Top Domains (default)")
	} else {
		var filters []string
		if f.URLFilter != "" {
			filters = append(filters, fmt.Sprintf("URL (%s)", f.URLFilter))
		}
		if f.CountryCode != "" {
			filters = append(filters, fmt.Sprintf("Country (%s - %s)", f.CountryName, f.CountryCode))
		}
		if f.RankFrom > 0 && f.RankTo > 0 {
			filters = append(filters, fmt.Sprintf("Rank (%d-%d)", f.RankFrom, f.RankTo))
		}
		if f.IPFrom != "" && f.IPTo != "" {
			filters = append(filters, fmt.Sprintf("IP Range (%s-%s)", f.IPFrom, f.IPTo))
		}
		if f.VisitorsFrom > 0 && f.VisitorsTo > 0 {
			filters = append(filters, fmt.Sprintf("Visitors (%d-%d)", f.VisitorsFrom, f.VisitorsTo))
		}
		if f.OwnerName != "" {
			filters = append(filters, fmt.Sprintf("Owner (%s - ID: %d)", f.OwnerName, f.OwnerID))
		}
		fmt.Print(strings.Join(filters, " + "))
	}

	proxyInfo := ""
	if c.ProxyURL != "" {
		proxyInfo = fmt.Sprintf("\nProxy: %s", c.ProxyURL)
		if c.ProxyUser != "" {
			proxyInfo += fmt.Sprintf(" (authenticated as %s)", c.ProxyUser)
		}
	}

	fmt.Printf("\nOutput: %s\nPages: %s (starting from page %d)%s\n",
		c.Output, getPagesDisplay(c.MaxPages), c.StartPage, proxyInfo)
}

// getPagesDisplay returns the appropriate display string for MaxPages
func getPagesDisplay(maxPages int) string {
	if maxPages == 0 {
		return "unlimited"
	}
	return fmt.Sprintf("%d", maxPages)
}
