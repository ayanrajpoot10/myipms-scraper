# MyIP.ms Domain Scraper

A powerful and comprehensive Go-based web scraper for extracting domain lists from myip.ms with advanced filtering capabilities and proxy support.

## 🚀 Overview

MyIP.ms Domain Scraper is a sophisticated command-line tool designed to efficiently extract domain information from myip.ms, one of the most comprehensive domain databases available. The tool supports multiple filtering options, proxy configurations, and includes advanced features like automatic cookie management and captcha handling.

## ✨ Key Features

- **Multi-Filter Support**: Filter domains by country, hosting provider, DNS records, IP ranges, visitor counts, and more
- **Advanced IP Range Filtering**: Supports both traditional IP ranges (from-to) and CIDR notation
- **Concurrent Scraping**: Multi-threaded scraping with configurable worker pools for faster results
- **Rate Limiting Control**: Adjustable delays between requests to respect server limits
- **Proxy Integration**: Full proxy support including HTTP, HTTPS, and SOCKS5 with authentication
- **Automatic Cookie Management**: Built-in cookie handling and captcha solving capabilities
- **Resume Functionality**: Start scraping from any page number to resume interrupted sessions
- **Intelligent Error Handling**: Automatic detection and handling of rate limits and authentication issues
- **Comprehensive Output Options**: Customizable output files and page limits
- **Smart Suggestions**: Helpful suggestions for misspelled filter options
- **Progress Tracking**: Real-time progress monitoring with detailed statistics

## 🛠️ Installation

### Prerequisites
- Go 1.21 or higher
- Internet connection

### Building from Source

```bash
# Install the scraper
go install github.com/ayanrajpoot10/myipms-scraper
```

### Download Binary
Pre-built binaries are available for multiple platforms in the [Releases](https://github.com/ayanrajpoot10/myipms-scraper/releases) section.

## 🎯 Quick Start

```bash
# Basic usage - scrape top domains globally
./myipms-scraper

# Filter by country
./myipms-scraper -country="United States"

# Filter by hosting provider
./myipms-scraper -owner="Cloudflare, Inc"

# Combined filters with custom output
./myipms-scraper -country="United States" -owner="Amazon.com, Inc" -output=aws-usa-domains.txt -pages=10

# Fast concurrent scraping
./myipms-scraper -country="United States" -workers=5 -delay=300

# Conservative scraping with single worker (classic mode)
./myipms-scraper -country="United States" -workers=1 -delay=1000
```

## 📊 Filter Options

### Geographic Filters
- **Country**: Filter by country name (e.g., `United States`, `United Kingdom`, `Japan`)
- **IP Range**: Filter by IP address ranges or CIDR blocks

### Hosting & Infrastructure
- **Owner**: Filter by hosting provider or company name
- **Host**: Filter by specific hosting service
- **DNS**: Filter by DNS record type

### Performance Metrics
- **Rank**: Filter by popularity ranking range
- **Visitors**: Filter by visitor count range
- **URL**: Search for domains containing specific text in their URLs

## 🌐 Proxy Support

The scraper includes comprehensive proxy support for enhanced reliability and bypassing rate limits:

### Supported Protocols
- HTTP/HTTPS proxies
- SOCKS5 proxies
- Authenticated proxies

### Configuration Examples
```bash
# Basic proxy
./myipms-scraper -proxy=http://proxy.example.com:8080

# Authenticated proxy
./myipms-scraper -proxy=http://proxy.example.com:8080@username:password

# SOCKS5 proxy (Tor)
./myipms-scraper -proxy=socks5://127.0.0.1:9050
```

## ⚡ Concurrent Scraping

The scraper now supports concurrent/parallel scraping for significantly faster results:

### Key Benefits
- **Faster Scraping**: Multiple pages processed simultaneously
- **Configurable Workers**: 1-10 concurrent workers (default: 3)
- **Smart Rate Limiting**: Adjustable delays to avoid server overload
- **Error Resilience**: Automatic retries and graceful degradation

### Concurrency Configuration
```bash
# Use 5 concurrent workers (recommended for fast scraping)
./myipms-scraper -workers=5

# Conservative approach with more delay between requests
./myipms-scraper -workers=3 -delay=1000

# Maximum speed (use with caution, may trigger rate limits)
./myipms-scraper -workers=8 -delay=200

# Classic sequential mode (single worker)
./myipms-scraper -workers=1
```

### Performance Guidelines
- **Default (3 workers, 500ms delay)**: Good balance of speed and stability
- **Fast (5-8 workers, 200-300ms delay)**: For bulk scraping with proxy/VPN
- **Conservative (1-2 workers, 1000ms delay)**: For avoiding rate limits
- **Maximum**: Up to 10 workers, but may trigger anti-bot measures

### Considerations
- More workers = faster scraping but higher chance of rate limiting
- Use appropriate delays to respect server resources
- Consider using proxies with high worker counts
- Monitor for IP rate limiting and reduce workers if needed

## 📈 Performance & Reliability

### Built-in Rate Limiting
- Configurable delays between requests (default: 500ms)
- Intelligent backoff for rate limit detection
- IP rotation recommendations for large-scale scraping
- Concurrent processing with worker pool management

### Error Recovery
- Automatic cookie refresh when expired
- Captcha solving workflow for authentication challenges
- Resume capability from any page number
- Comprehensive error reporting and suggestions
- Per-page retry logic with exponential backoff

### Output Management
- Real-time progress tracking with worker statistics
- Configurable output files
- Automatic deduplication
- Statistics reporting for concurrent operations

## 🔧 Configuration

### Cookie Management
First-time users will be prompted to solve a captcha for authentication. Cookies are automatically saved and reused for subsequent runs.

## 📋 Available Filters

### Countries (Sample)
- United States, Canada, United Kingdom, Germany, France, Japan, Australia, Brazil, India, Russia
- Only accepts full country names (not country codes)

### Major Hosting Providers (Sample)
- Amazon.com, Inc, Cloudflare Inc, Google Inc, Microsoft Corp
- DigitalOcean, Linode, Vultr Holdings LLC, OVH SAS
- GoDaddy.com LLC, Namecheap Inc, Hostinger International

### DNS Records (Sample)
- ns1.cloudflare.com, ns2.cloudflare.com
- ns1.amazonaws.com, ns2.amazonaws.com
- dns1.registrar-servers.com, dns2.registrar-servers.com

*Use `./myipms-scraper --list` to see all available options*

## 🚨 Important Notes

### Rate Limiting & Best Practices
- The tool includes configurable rate limiting (default: 500ms delays)
- Concurrent scraping allows faster processing while respecting limits
- For large-scale scraping, consider using proxies or VPNs
- Respect the website's terms of service
- Monitor for IP rate limiting and adjust worker count/delays accordingly
- Start with default settings (3 workers, 500ms delay) and adjust based on results

### Authentication Requirements
- First-time usage requires solving a captcha
- Cookies are saved locally for future sessions
- Periodic re-authentication may be required

### Legal Considerations
- Ensure compliance with local laws and regulations
- Respect website terms of service
- Use scraped data responsibly and ethically

## 🔍 Troubleshooting

### Common Issues

1. **Cookies Expired Error**
   - Solution: Run the captcha solving process again
   - The tool will automatically prompt for captcha resolution

2. **IP Limit Exceeded**
   - Solution: Use a proxy server or change your IP address
   - Consider using mobile internet with airplane mode toggle

3. **No Domains Found**
   - Check filter criteria for typos
   - Use `--list` to verify available options
   - Try broader filter criteria

4. **Connection Issues**
   - Verify internet connectivity
   - Check proxy configuration if using one
   - Consider firewall or DNS issues

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Thanks to myip.ms for providing comprehensive domain data
- Built with Go's powerful standard library

## 📞 Support

- Create an [Issue](https://github.com/ayanrajpoot10/myipms-scraper/issues) for bug reports or feature requests
- Check the [Command Reference](COMMANDS.md) for detailed usage instructions
- Review existing issues before creating new ones

---

**⚠️ Disclaimer**: This tool is for educational and research purposes. Users are responsible for ensuring compliance with applicable laws and the target website's terms of service.
