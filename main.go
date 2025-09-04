package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	config := parseFlags()

	if config.List {
		showSpecificOptions(config.Owner, config.Country, config.Host, config.DNSRecord)
		return
	}

	filter, err := validateAndResolveFilters(config)
	if err != nil {
		handleValidationError(err)
		return
	}

	displayScrapingFilter(filter, config)

	httpClient := newHTTPClient(config.ProxyURL, config.ProxyUser, config.ProxyPass)
	scraper := NewScraper(httpClient, filter)

	file, err := os.Create(config.Output)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	totalDomains := 0
	page := config.StartPage
	maxPage := config.StartPage + config.MaxPages

	isUnlimited := config.MaxPages == 0

	for isUnlimited || page < maxPage {
		domains, err := scraper.fetchPage(page)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			if strings.Contains(err.Error(), "cookies expired") {
				captchaErr := solveCaptcha(httpClient)
				if captchaErr != nil {
					fmt.Printf("Failed to solve captcha: %v\n", captchaErr)
					fmt.Println("Please restart the program and try again.")
					os.Exit(1)
				}
				continue
			} else if strings.Contains(err.Error(), "IP limit exceeded") {
				fmt.Println("IP address has been rate limited. Please:")
				fmt.Println("1. Use a proxy or VPN to change your IP address")
				fmt.Println("2. If using mobile internet, turn airplane mode on/off to get a new IP")
				os.Exit(1)
			}
			continue
		}

		if len(domains) == 0 {
			fmt.Println("No domains found")
			break
		}

		for _, domain := range domains {
			file.WriteString(domain + "\n")
		}

		totalDomains += len(domains)
		fmt.Printf("Page %d: Found %d domains (Total: %d)\n", page, len(domains), totalDomains)
		page++
	}

	fmt.Printf("Scraping complete! Total domains: %d\n", totalDomains)
}