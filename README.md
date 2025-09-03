# MyIP.ms Domain Scraper

A powerful and comprehensive Go-based web scraper for extracting domain lists from myip.ms with advanced filtering capabilities and proxy support.

## 🚀 Overview

MyIP.ms Domain Scraper is a sophisticated command-line tool designed to efficiently extract domain information from myip.ms, one of the most comprehensive domain databases available. The tool supports multiple filtering options, proxy configurations, and includes advanced features like automatic cookie management and captcha handling.

## ✨ Key Features

- **Multi-Filter Support**: Filter domains by country, hosting provider, DNS records, IP ranges, visitor counts, and more
- **Advanced IP Range Filtering**: Supports both traditional IP ranges (from-to) and CIDR notation
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

## 📈 Performance & Reliability

### Built-in Rate Limiting
- Configurable delays between requests (default: 500ms)
- Intelligent backoff for rate limit detection
- IP rotation recommendations for large-scale scraping

### Error Recovery
- Automatic cookie refresh when expired
- Captcha solving workflow for authentication challenges
- Resume capability from any page number
- Comprehensive error reporting and suggestions
- Per-page retry logic with exponential backoff

### Output Management
- Configurable output files
- Automatic deduplication

## 🔧 Configuration

### Cookie Management
The scraper now features an **improved web-based captcha solver** for better user experience:

- **Web Interface**: Captcha images are served through a local web server at `http://localhost:8080`
- **Auto-Launch Browser**: The default browser automatically opens to the captcha solving page
- **Real-time Feedback**: Immediate feedback on captcha submission success or failure
- **No File Management**: No need to manually check captcha image files
- **Retry Logic**: Automatic captcha refresh on failed attempts
- **Mobile Friendly**: Works on any device that can access localhost

When cookies expire, the tool will:
1. Automatically start a local web server
2. Open your browser to the captcha solving interface
3. Display the captcha image directly in the browser
4. Allow you to submit the solution through a user-friendly form
5. Automatically continue scraping once solved

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
- For large-scale scraping, consider using proxies or VPNs
- Respect the website's terms of service

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
