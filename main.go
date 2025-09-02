package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// PageResult represents the result of scraping a single page
type PageResult struct {
	Page    int
	Domains []string
	Error   error
}

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

	// Use the unified scraping logic that handles both concurrent and sequential modes
	err = scrapePages(scraper, config)
	if err != nil {
		fmt.Printf("Error during scraping: %v\n", err)
		os.Exit(1)
	}
}

// scrapePages performs scraping with configurable concurrency
func scrapePages(scraper *Scraper, config *Config) error {
	file, err := os.Create(config.Output)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer file.Close()

	workers := config.Workers
	delay := config.Delay

	// Display mode information
	if workers == 1 {
		fmt.Println("Using sequential scraping (single worker)")
	} else {
		fmt.Printf("Using concurrent scraping with %d workers\n", workers)
	}

	// Create channels for jobs and results
	jobs := make(chan int, workers*2)
	results := make(chan PageResult, workers*2)

	// Shared state
	var totalDomains int
	var mutex sync.Mutex
	retryCount := make(map[int]int)
	maxRetries := 3

	// Start workers
	var wg sync.WaitGroup
	for i := 1; i <= workers; i++ {
		wg.Add(1)
		go worker(i, scraper, delay, jobs, results, &wg)
	}

	// Close results channel when all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Calculate page range
	startPage := config.StartPage
	maxPage := startPage + config.MaxPages
	isUnlimited := config.MaxPages == 0

	// Send initial batch of jobs
	activePage := startPage
	activeJobs := 0
	noMorePages := false

	// Send first batch of pages to workers
	for i := 0; i < workers && (isUnlimited || activePage < maxPage); i++ {
		jobs <- activePage
		activeJobs++
		activePage++
	}

	if workers > 1 {
		fmt.Printf("Started %d concurrent workers with %dms delay between requests\n", workers, config.Delay/time.Millisecond)
	}

	// Process results
	for result := range results {
		activeJobs--

		shouldContinue := handlePageResult(result, file, jobs, &totalDomains, &mutex,
			retryCount, maxRetries, scraper, config)
		if !shouldContinue {
			noMorePages = true
		}

		// Send next page if we should continue and haven't reached the limit
		if !noMorePages && (isUnlimited || activePage < maxPage) {
			jobs <- activePage
			activeJobs++
			activePage++
		}

		// If no more jobs and all active jobs completed, we're done
		if noMorePages && activeJobs == 0 {
			break
		}

		// If we've reached the page limit and no active jobs, we're done
		if !isUnlimited && activePage >= maxPage && activeJobs == 0 {
			break
		}
	}

	// Close jobs channel to signal workers to stop
	close(jobs)

	fmt.Printf("Scraping complete! Total domains: %d\n", totalDomains)
	return nil
}

// worker processes pages from the jobs channel and sends results to the results channel
func worker(id int, scraper *Scraper, delay time.Duration, jobs <-chan int, results chan<- PageResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for page := range jobs {
		// Add delay to respect server rate limits
		if delay > 0 {
			time.Sleep(delay)
		}

		if id > 1 {
			fmt.Printf("Worker %d fetching page %d... ", id, page)
		} else {
			fmt.Printf("Fetching page %d... ", page)
		}

		domains, err := scraper.fetchPage(page)

		result := PageResult{
			Page:    page,
			Domains: domains,
			Error:   err,
		}

		results <- result

		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Found %d domains\n", len(domains))
		}
	}
}

// handlePageResult processes a single page result and handles retries
func handlePageResult(result PageResult, file *os.File, jobs chan<- int, totalDomains *int,
	mutex *sync.Mutex, retryCount map[int]int, maxRetries int, scraper *Scraper, config *Config) bool {

	mutex.Lock()
	defer mutex.Unlock()

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "cookies expired") {
			fmt.Printf("Page %d: Cookies have expired. Attempting to solve captcha...\n", result.Page)

			captchaErr := solveCaptcha()
			if captchaErr != nil {
				fmt.Printf("Failed to solve captcha: %v\n", captchaErr)
				fmt.Println("Please restart the program and try again.")
				return false
			}

			// Update the scraper with fresh cookies
			newHTTPClient := newHTTPClient(config.ProxyURL, config.ProxyUser, config.ProxyPass)
			scraper.UpdateHTTPClient(newHTTPClient)
			fmt.Println("Retrying with fresh cookies...")

			// Retry the page
			if retryCount[result.Page] < maxRetries {
				retryCount[result.Page]++
				jobs <- result.Page
			}
			return true

		} else if strings.Contains(result.Error.Error(), "IP limit exceeded") {
			fmt.Println("IP address has been rate limited. Please:")
			fmt.Println("1. Use a proxy or VPN to change your IP address")
			fmt.Println("2. If using mobile internet, turn airplane mode on/off to get a new IP")
			fmt.Println("3. Wait some time before trying again")
			fmt.Println("4. Try reducing the scraping speed or page count")
			if config.Workers > 1 {
				fmt.Println("5. Reduce the number of workers with --workers flag")
			}
			return false

		} else {
			// Other errors - retry if we haven't exceeded max retries
			if retryCount[result.Page] < maxRetries {
				retryCount[result.Page]++
				fmt.Printf("Retrying page %d (attempt %d/%d)\n", result.Page, retryCount[result.Page], maxRetries)
				time.Sleep(time.Second * time.Duration(retryCount[result.Page])) // Exponential backoff
				jobs <- result.Page
			}
			return true
		}
	}

	// Success case
	if len(result.Domains) == 0 {
		fmt.Printf("Page %d: No domains found, assuming end of results\n", result.Page)
		return false
	}

	// Write domains to file
	for _, domain := range result.Domains {
		file.WriteString(domain + "\n")
	}

	*totalDomains += len(result.Domains)
	fmt.Printf("Page %d: Found %d domains (Total: %d)\n", result.Page, len(result.Domains), *totalDomains)

	return true
}
