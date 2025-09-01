package main

import (
	"fmt"
	"sort"
	"strings"
)

// showHelp displays the help message
func showHelp() {
	fmt.Println("Scrapes domain lists from myip.ms with various filtering options.")
	fmt.Println()

	fmt.Println("USAGE:")
	fmt.Println("  scraper [OPTIONS]")
	fmt.Println()

	fmt.Println("FILTER OPTIONS:")
	fmt.Println("  -country <code>    Filter by country (e.g., USA, UK, JP)")
	fmt.Println("  -owner <name>      Filter by hosting provider (e.g., \"Cloudflare, Inc\")")
	fmt.Println("  -host <name>       Filter by specific host")
	fmt.Println("  -dns <record>      Filter by DNS record")
	fmt.Println("  -url <text>        Filter by URL containing specific text (e.g., wiki, blog)")
	fmt.Println("  -rank <range>      Filter by popularity ranking range (e.g., 10-20)")
	fmt.Println("  -ip <range>        Filter by IP address range (e.g., 104.16.0.0-104.16.255.255 or 192.168.0.0/24)")
	fmt.Println("  -visitors <range>  Filter by visitor count range (e.g., 1000-20000)")
	fmt.Println()

	fmt.Println("OUTPUT OPTIONS:")
	fmt.Println("  -output <file>     Output filename (default: domains.txt)")
	fmt.Println("  -pages <num>       Max pages to scrape (0 = unlimited, default: unlimited)")
	fmt.Println("  -start <num>       Starting page number (default: 1)")
	fmt.Println()

	fmt.Println("CONCURRENCY OPTIONS:")
	fmt.Println("  -workers <num>     Number of concurrent workers (default: 3, max: 10)")
	fmt.Println("  -delay <ms>        Delay between requests in milliseconds (default: 500ms)")
	fmt.Println("                     Use -workers=1 for sequential processing")
	fmt.Println()

	fmt.Println("PROXY OPTIONS:")
	fmt.Println("  -proxy <url>       Proxy URL with optional authentication")
	fmt.Println("                     Format: protocol://host:port[@user:pass]")
	fmt.Println("                     Examples: http://proxy.com:8080")
	fmt.Println("                              http://proxy.com:8080@user:pass")
	fmt.Println("                              socks5://127.0.0.1:9050")
	fmt.Println()

	fmt.Println("OTHER:")
	fmt.Println("  -help             Show this help message")
	fmt.Println("  -list             Show all available options (or specific options when used with filter flags)")
	fmt.Println()

	fmt.Println("EXAMPLES:")
	fmt.Println("  Basic usage:")
	fmt.Println("    scraper                                    # Get top domains globally")
	fmt.Println("    scraper -country=USA                       # Get domains in USA")
	fmt.Println("    scraper -owner=\"Cloudflare, Inc\"           # Get Cloudflare domains")
	fmt.Println("    scraper -url=wiki                          # Get domains with 'wiki' in URL")
	fmt.Println("    scraper -rank=10-20                        # Get domains ranked 10-20")
	fmt.Println("    scraper -ip=104.16.0.0-104.16.255.255      # Get domains in IP range")
	fmt.Println("    scraper -ip=\"192.168.0.0/24\"               # Get domains in CIDR range (use quotes)")
	fmt.Println("    scraper -ip=\"10.0.0.0/8\"                   # Get domains in large CIDR range")
	fmt.Println("    scraper -visitors=1000-20000               # Get domains with 1K-20K visitors")
	fmt.Println()

	fmt.Println("  Combined filters:")
	fmt.Println("    scraper -country=USA -owner=\"Amazon.com, Inc\"")
	fmt.Println("    scraper -country=JP -dns=\"01.dnsv.jp\"")
	fmt.Println("    scraper -owner=\"Google Inc\" -host=\"example-host-1\"")
	fmt.Println("    scraper -url=wiki -country=USA             # Get USA domains with 'wiki' in URL")
	fmt.Println("    scraper -url=blog -owner=\"Cloudflare, Inc\" # Get Cloudflare domains with 'blog' in URL")
	fmt.Println("    scraper -country=USA -rank=10-20           # Get USA domains ranked 10-20")
	fmt.Println("    scraper -url=shop -rank=1-100              # Get top 100 domains with 'shop' in URL")
	fmt.Println("    scraper -rank=1-100 -ip=104.16.0.0-104.16.255.255  # Cloudflare IP range, top 100")
	fmt.Println("    scraper -country=USA -ip=\"8.8.8.0/24\"              # USA domains in Google DNS CIDR")
	fmt.Println("    scraper -visitors=10000-50000 -country=USA           # USA domains with 10K-50K visitors")
	fmt.Println("    scraper -rank=1-50 -visitors=100000-1000000          # Top 50 domains with 100K-1M visitors")
	fmt.Println()

	fmt.Println("  Custom output:")
	fmt.Println("    scraper -country=USA -output=usa_domains.txt -pages=50")
	fmt.Println("    scraper -country=USA -start=5 -pages=10         # Start from page 5, scrape 10 pages")
	fmt.Println("    scraper -pages=0                                # Scrape unlimited pages (same as default)")
	fmt.Println()

	fmt.Println("  Concurrent scraping:")
	fmt.Println("    scraper -workers=5                              # Use 5 concurrent workers")
	fmt.Println("    scraper -workers=3 -delay=1000                  # 3 workers, 1 second delay")
	fmt.Println("    scraper -workers=1                              # Sequential processing")
	fmt.Println("    scraper -country=USA -workers=8 -delay=200      # Fast scraping with rate limiting")
	fmt.Println()

	fmt.Println("  With proxy:")
	fmt.Println("    scraper -proxy=http://proxy.example.com:8080")
	fmt.Println("    scraper -proxy=http://proxy.example.com:8080@username:password")
	fmt.Println("    scraper -proxy=socks5://127.0.0.1:9050           # Using Tor proxy")
	fmt.Println("    scraper -country=USA -proxy=http://proxy.com:8080@user:pass")
	fmt.Println()

	fmt.Println("  List available options:")
	fmt.Println("    scraper --list                             # Show all available options")
	fmt.Println("    scraper -owner=\"\" --list                   # Show only owner options")
	fmt.Println("    scraper -country=\"\" --list                 # Show only country options")
	fmt.Println("    scraper -dns=\"\" --list                     # Show only DNS options")
	fmt.Println()

	fmt.Println("NOTES:")
	fmt.Println("  • All filter options can be combined")
	fmt.Println("  • Range filters use the format 'from-to' (e.g., 10-20, 1000-5000)")
	fmt.Println("  • Range values must be positive integers (from ≤ to)")
	fmt.Println("  • IP addresses must be valid IPv4 or IPv6 format")
	fmt.Println("  • IP ranges support both 'from-to' format and CIDR notation (e.g., 192.168.0.0/24)")
	fmt.Println("  • Use quotes around CIDR notation in shell commands (e.g., -ip=\"192.168.0.0/24\")")
	fmt.Println("  • Start page allows resuming scraping from a specific page")
	fmt.Println("  • Proxy format: protocol://host:port[@user:pass] (auth optional)")
	fmt.Println("  • Supported protocols: HTTP, HTTPS, and SOCKS5")
	fmt.Println("  • Concurrent scraping uses multiple workers for faster results")
	fmt.Println("  • Higher worker count = faster scraping, but may trigger rate limits")
	fmt.Println("  • Use appropriate delay between requests to respect server limits")
	fmt.Println("  • First run will prompt for cookies (required for authentication)")
	fmt.Println("  • Cookies are saved to cookies.json for future runs")
	fmt.Println("  • If scraping fails, you may need to update cookies, change IP, or use a proxy")
}

// Helper function to display a category of options with int values
func displayCategoryInt(title string, items map[string]int, showTotal bool) {
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

// Helper function to display a category of options with string values
func displayCategoryString(title string, items map[string]string, showTotal bool) {
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

	// Show title based on whether specific filters are requested
	if hasAnyFlag {
		fmt.Println("Available Options for Specified Filters")
		fmt.Println("=======================================")
	} else {
		fmt.Println("All Available Filter Options")
		fmt.Println("============================")
	}
	fmt.Println()

	// Show countries if requested specifically or if showing all options
	if country != "" || !hasAnyFlag {
		displayCategoryString("COUNTRIES", countries, hasAnyFlag)
	}

	// Show owners if requested specifically or if showing all options
	if owner != "" || !hasAnyFlag {
		displayCategoryInt("OWNERS/HOSTING PROVIDERS", owners, hasAnyFlag)
	}

	// Show hosts if requested specifically or if showing all options
	if host != "" || !hasAnyFlag {
		displayCategoryInt("HOSTS", hosts, hasAnyFlag)
	}

	// Show DNS records if requested specifically or if showing all options
	if dnsRecord != "" || !hasAnyFlag {
		displayCategoryInt("DNS RECORDS", dns, hasAnyFlag)
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
			filters = append(filters, fmt.Sprintf("Country (%s - %s)", filter.CountryCode, filter.CountryName))
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
	concurrencyInfo := fmt.Sprintf("\nConcurrency: %d workers, %dms delay", config.Workers, config.Delay)
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
