# MyIP.ms Scraper - Command Reference Guide

This comprehensive guide covers all command-line options, usage patterns, and advanced configurations for the MyIP.ms Domain Scraper.

## Table of Contents

- [Basic Usage](#basic-usage)
- [Command Syntax](#command-syntax)
- [Filter Options](#filter-options)
- [Output Configuration](#output-configuration)
- [Proxy Configuration](#proxy-configuration)
- [Advanced Usage Patterns](#advanced-usage-patterns)
- [Error Handling](#error-handling)
- [Best Practices](#best-practices)

## Basic Usage

### Default Behavior
```bash
# Scrape top domains globally with default settings
./myipms-scraper
```

**Default Values:**
- Output file: `domains.txt`
- Pages: Unlimited (0)
- Start page: 1
- No proxy
- No filters applied

### Getting Help
```bash
# Display help message
./myipms-scraper -help
./myipms-scraper --help

# List all available filter options
./myipms-scraper --list
```

## Command Syntax

### General Format
```bash
./myipms-scraper [FILTER_OPTIONS] [OUTPUT_OPTIONS] [PROXY_OPTIONS] [OTHER_OPTIONS]
```

### Flag Conventions
- Flags use single dash: `-country`, `-owner`, `-pages`
- Boolean flags: `-help`, `-list`
- String values can be quoted: `-owner="Cloudflare, Inc"`
- No spaces around equals sign: `-country=USA` (not `-country = USA`)

## Filter Options

### Geographic Filters

#### Country Filter
Filter domains by country name (only full country names accepted).

```bash
# By country name (case-sensitive)
./myipms-scraper -country="United States"
./myipms-scraper -country="United Kingdom"
./myipms-scraper -country="Japan"

# Get list of available countries
./myipms-scraper -country="" --list
```

**Popular Countries (use full names):**
- `United States`
- `United Kingdom`
- `Germany`
- `France`
- `Japan`
- `Canada`
- `Australia`
- `Brazil`
- `India`
- `Russia`

### Hosting & Infrastructure Filters

#### Owner/Hosting Provider Filter
Filter by hosting company or organization.

```bash
# Major cloud providers
./myipms-scraper -owner="Amazon.com, Inc"
./myipms-scraper -owner="Cloudflare, Inc"
./myipms-scraper -owner="Google Inc"
./myipms-scraper -owner="Microsoft Corp"

# Web hosting companies
./myipms-scraper -owner="GoDaddy.com, LLC"
./myipms-scraper -owner="Namecheap, Inc"
./myipms-scraper -owner="DigitalOcean, Inc"
./myipms-scraper -owner="Linode, LLC"

# Get list of available owners
./myipms-scraper -owner="" --list
```

#### Host Filter
Filter by specific hosting service.

```bash
# Examples
./myipms-scraper -host="amazonaws.com"
./myipms-scraper -host="cloudflare.com"
./myipms-scraper -host="googleusercontent.com"

# Get list of available hosts
./myipms-scraper -host="" --list
```

#### DNS Record Filter
Filter by DNS server.

```bash
# Cloudflare DNS
./myipms-scraper -dns="andy.ns.cloudflare.com"
./myipms-scraper -dns="lisa.ns.cloudflare.com"

# Google DNS
./myipms-scraper -dns="ns-cloud-a1.googledomains.com"

# Get list of available DNS records
./myipms-scraper -dns="" --list
```

### Content & Performance Filters

#### URL Content Filter
Search for domains containing specific text in their URLs.

```bash
# Content-based filtering
./myipms-scraper -url=wiki
./myipms-scraper -url=blog
./myipms-scraper -url=shop
./myipms-scraper -url=api
./myipms-scraper -url=cdn

# Multiple keywords (separate commands needed)
./myipms-scraper -url=admin
./myipms-scraper -url=dashboard
```

#### Ranking Filter
Filter by website popularity ranking.

```bash
# Top 100 websites
./myipms-scraper -rank=1-100

# Websites ranked 1000-5000
./myipms-scraper -rank=1000-5000

# Mid-tier websites
./myipms-scraper -rank=10000-50000

# Format: from-to (both numbers must be positive, from ≤ to)
```

#### Visitor Count Filter
Filter by estimated visitor numbers.

```bash
# High-traffic sites (100K-1M visitors)
./myipms-scraper -visitors=100000-1000000

# Medium-traffic sites (10K-100K visitors)
./myipms-scraper -visitors=10000-100000

# Low-traffic sites (1K-10K visitors)
./myipms-scraper -visitors=1000-10000

# Format: from-to (both numbers must be positive)
```

### IP Address Filters

#### IP Range Filter
Filter by IP address ranges using either traditional ranges or CIDR notation.

```bash
# Traditional IP range format
./myipms-scraper -ip=104.16.0.0-104.16.255.255
./myipms-scraper -ip=8.8.8.0-8.8.8.255

# CIDR notation (use quotes in shell)
./myipms-scraper -ip="192.168.0.0/24"
./myipms-scraper -ip="10.0.0.0/8"
./myipms-scraper -ip="172.16.0.0/12"

# CloudFlare IP ranges (examples)
./myipms-scraper -ip="103.21.244.0/22"
./myipms-scraper -ip="103.22.200.0/22"

# AWS IP ranges (examples)  
./myipms-scraper -ip="52.0.0.0/11"
./myipms-scraper -ip="54.0.0.0/8"
```

**Important IP Range Notes:**
- CIDR notation must be quoted in most shells
- Supports both IPv4 and IPv6 addresses
- Range format: `from_ip-to_ip` or `network/prefix_length`
- Invalid formats will show error with examples

## Output Configuration

### Output File Options
```bash
# Custom output filename
./myipms-scraper -output=custom_domains.txt
./myipms-scraper -output=/path/to/domains.txt

# Organize by date
./myipms-scraper -output=domains_$(date +%Y%m%d).txt

# Organize by filter type
./myipms-scraper -country=USA -output=usa_domains.txt
./myipms-scraper -owner="Cloudflare, Inc" -output=cloudflare_domains.txt
```

### Page Control Options
```bash
# Limit number of pages
./myipms-scraper -pages=50              # Scrape exactly 50 pages
./myipms-scraper -pages=0               # Unlimited pages (default)

# Start from specific page (useful for resuming)
./myipms-scraper -start=10              # Start from page 10
./myipms-scraper -start=5 -pages=20     # Pages 5-24 (20 pages total)

# Resume interrupted scraping
./myipms-scraper -country=USA -start=25 -pages=25 -output=usa_domains_continued.txt
```

## Proxy Configuration

### Basic Proxy Setup
```bash
# HTTP proxy
./myipms-scraper -proxy=http://proxy.example.com:8080

# HTTPS proxy  
./myipms-scraper -proxy=https://secure-proxy.example.com:8080

# SOCKS5 proxy
./myipms-scraper -proxy=socks5://socks-proxy.example.com:1080
```

### Authenticated Proxy
```bash
# HTTP proxy with authentication
./myipms-scraper -proxy=http://proxy.example.com:8080@username:password

# HTTPS proxy with authentication
./myipms-scraper -proxy=https://proxy.example.com:8080@user:pass

# SOCKS5 proxy with authentication  
./myipms-scraper -proxy=socks5://proxy.example.com:1080@myuser:mypass
```

### Proxy with Filtering
```bash
# Combine proxy with country filter
./myipms-scraper -proxy=http://proxy.com:8080@user:pass -country=USA

# Use Tor for anonymity
./myipms-scraper -proxy=socks5://127.0.0.1:9050 -country=RUS -pages=10
```

**Proxy Format Rules:**
- Format: `protocol://host:port[@username:password]`
- Authentication is optional (use `@username:password`)
- Supported protocols: `http://`, `https://`, `socks5://`
- Invalid formats will show detailed error messages

## Advanced Usage Patterns

### Complex Multi-Filter Queries

#### Geographic + Hosting Combinations
```bash
# US domains hosted by Amazon
./myipms-scraper -country=USA -owner="Amazon.com, Inc" -output=aws_usa.txt

# German domains using Cloudflare
./myipms-scraper -country=DEU -owner="Cloudflare, Inc" -output=cloudflare_germany.txt

# Japanese domains with specific DNS
./myipms-scraper -country=JPN -dns="ns1.google.com" -output=google_dns_japan.txt
```

#### Performance-Based Filtering
```bash
# High-traffic US e-commerce sites
./myipms-scraper -country=USA -url=shop -visitors=50000-500000 -output=us_ecommerce.txt

# Top 1000 global domains hosted on AWS
./myipms-scraper -rank=1-1000 -owner="Amazon.com, Inc" -output=top1000_aws.txt

# Medium-traffic blog sites
./myipms-scraper -url=blog -visitors=10000-100000 -rank=5000-50000 -output=blog_sites.txt
```

#### Infrastructure Analysis
```bash
# All domains in Cloudflare IP ranges
./myipms-scraper -ip="104.16.0.0/12" -output=cloudflare_ips.txt

# Domains using specific DNS with visitor filtering  
./myipms-scraper -dns="ns1.cloudflare.com" -visitors=100000-1000000 -output=cf_dns_popular.txt

# Geographic distribution analysis
./myipms-scraper -country=USA -ip="8.8.8.0/24" -output=usa_google_dns.txt
```

### Batch Processing Examples

#### Daily Monitoring Script
```bash
#!/bin/bash
DATE=$(date +%Y%m%d)
PROXY="http://proxy.example.com:8080@user:pass"

# Monitor top US domains
./myipms-scraper -country=USA -rank=1-1000 -proxy="$PROXY" -output="usa_top1000_$DATE.txt"

# Monitor Cloudflare infrastructure
./myipms-scraper -owner="Cloudflare, Inc" -proxy="$PROXY" -output="cloudflare_$DATE.txt"

# Monitor specific IP ranges
./myipms-scraper -ip="104.16.0.0/12" -proxy="$PROXY" -output="cf_ips_$DATE.txt"
```

#### Competitive Analysis
```bash
# E-commerce landscape analysis
./myipms-scraper -url=shop -country=USA -visitors=10000-1000000 -output=us_ecommerce_analysis.txt

# CDN usage analysis  
./myipms-scraper -owner="Cloudflare, Inc" -rank=1-10000 -output=cf_top_domains.txt
./myipms-scraper -owner="Amazon.com, Inc" -rank=1-10000 -output=aws_top_domains.txt

# Geographic distribution
./myipms-scraper -country=USA -pages=100 -output=usa_sample.txt
./myipms-scraper -country=GBR -pages=100 -output=uk_sample.txt
```

### Resume & Continue Patterns
```bash
# Original scraping (interrupted at page 47)
./myipms-scraper -country=USA -owner="Amazon.com, Inc" -output=aws_usa.txt

# Resume from where it stopped
./myipms-scraper -country=USA -owner="Amazon.com, Inc" -start=48 -output=aws_usa_continued.txt

# Combine results later
cat aws_usa.txt aws_usa_continued.txt > aws_usa_complete.txt
```

## Error Handling

### Common Error Scenarios

#### Cookie Expiration
```bash
# When you see: "cookies expired - human verification required"
# The tool will automatically prompt for captcha solving

# Process will be:
# 1. Captcha image downloaded as captcha_image.png
# 2. User enters captcha text when prompted  
# 3. New cookies are obtained automatically
# 4. Scraping continues
```

#### IP Rate Limiting
```bash
# When you see: "IP limit exceeded"
# Solutions (in order of preference):

# 1. Use a proxy
./myipms-scraper -proxy=http://proxy.com:8080 -country=USA

# 2. Use Tor (if available)
./myipms-scraper -proxy=socks5://127.0.0.1:9050 -country=USA  

# 3. Change IP (mobile users)
# Turn airplane mode on/off to get new IP, then retry

# 4. Reduce scraping speed/volume
./myipms-scraper -country=USA -pages=10  # Smaller batch
```

#### Invalid Filter Values
```bash
# When you see: "unknown country 'INVALID'"
# Tool provides suggestions:

# Example error output:
# Error: unknown country 'USA'
# Did you mean one of these?
#   1. United States
#   2. USA  
#   3. US Virgin Islands

# Use suggested values:
./myipms-scraper -country="United States"
```

### Validation Errors

#### Range Format Errors
```bash
# Invalid range formats and fixes:

# ❌ Wrong: -rank=100-10 (from > to)
# ✅ Correct: -rank=10-100

# ❌ Wrong: -visitors=-1000-5000 (negative values)  
# ✅ Correct: -visitors=1000-5000

# ❌ Wrong: -ip=invalid-range
# ✅ Correct: -ip=192.168.1.0-192.168.1.255 or -ip="192.168.1.0/24"
```

#### Proxy Configuration Errors
```bash
# Common proxy errors and solutions:

# ❌ Wrong: -proxy=proxy.com:8080 (missing protocol)
# ✅ Correct: -proxy=http://proxy.com:8080

# ❌ Wrong: -proxy=http://proxy.com:8080:user:pass (wrong auth format)
# ✅ Correct: -proxy=http://proxy.com:8080@user:pass

# ❌ Wrong: -proxy=http://proxy.com (missing port) 
# ✅ Correct: -proxy=http://proxy.com:8080
```

## Best Practices

### Performance Optimization

#### Efficient Filtering
```bash
# ✅ Good: Use specific filters to reduce data volume
./myipms-scraper -country=USA -rank=1-1000 -output=focused_results.txt

# ❌ Avoid: Scraping everything then filtering locally  
./myipms-scraper -pages=1000 | grep "amazonaws"  # Inefficient
```

#### Batch Size Management  
```bash
# ✅ Good: Reasonable batch sizes
./myipms-scraper -country=USA -pages=50 -output=batch1.txt

# ❌ Avoid: Extremely large batches (may trigger rate limits)
./myipms-scraper -pages=10000  # Too aggressive
```

### Reliability Patterns

#### Use Proxy for Large Jobs
```bash
# ✅ Recommended for large-scale scraping
./myipms-scraper -proxy=http://proxy.com:8080@user:pass -country=USA -pages=200

# Consider rotating proxies for very large jobs
```

#### Resume Strategy
```bash
# ✅ Always plan for interruption recovery
./myipms-scraper -country=USA -start=1 -pages=50 -output=batch1.txt
./myipms-scraper -country=USA -start=51 -pages=50 -output=batch2.txt  
./myipms-scraper -country=USA -start=101 -pages=50 -output=batch3.txt
```

### Data Organization

#### Structured Output Naming
```bash
# ✅ Good: Descriptive, organized filenames
./myipms-scraper -country=USA -output="$(date +%Y%m%d)_usa_domains.txt"
./myipms-scraper -owner="Cloudflare, Inc" -output="$(date +%Y%m%d)_cloudflare_domains.txt"

# ✅ Good: Directory organization
mkdir -p data/$(date +%Y%m%d)
./myipms-scraper -country=USA -output="data/$(date +%Y%m%d)/usa_domains.txt"
```

### Monitoring & Logging

#### Progress Tracking
```bash
# ✅ Use page limits for predictable progress
./myipms-scraper -country=USA -pages=100 -output=usa_domains.txt 2>&1 | tee scraping.log

# ✅ Monitor output file growth
tail -f usa_domains.txt &  # In separate terminal
./myipms-scraper -country=USA -output=usa_domains.txt
```

### Legal & Ethical Guidelines

#### Respect Rate Limits
```bash  
# ✅ Use proxies for legitimate large-scale needs
# ✅ Avoid overwhelming the service

# ❌ Don't try to bypass rate limiting artificially
# ❌ Don't run multiple instances simultaneously without proxies
```

#### Data Usage Ethics
```bash
# ✅ Use scraped data for legitimate research/analysis
# ✅ Respect intellectual property rights
# ✅ Follow applicable data protection laws

# Document your usage:
echo "Scraped $(date): Country=USA, Purpose=Market Research" >> scraping_log.txt
./myipms-scraper -country=USA -output=market_research_usa.txt
```

---

## Quick Reference Summary

### Essential Commands
```bash
# Basic usage
./myipms-scraper                                    # Default scraping
./myipms-scraper -help                              # Show help
./myipms-scraper --list                             # List all options

# Common filters  
./myipms-scraper -country=USA                       # By country
./myipms-scraper -owner="Cloudflare, Inc"           # By hosting provider
./myipms-scraper -rank=1-100                        # Top 100 sites
./myipms-scraper -ip="192.168.0.0/24"               # IP range (CIDR)

# Output control
./myipms-scraper -output=custom.txt -pages=50       # Custom output, limited pages
./myipms-scraper -start=10 -pages=20                # Resume from page 10

# Proxy usage
./myipms-scraper -proxy=http://proxy.com:8080@user:pass  # Authenticated proxy

# Complex combinations
./myipms-scraper -country=USA -owner="Amazon.com, Inc" -rank=1-1000 -output=aws_usa_top1000.txt
```

### Flag Reference
| Flag | Type | Description | Example |
|------|------|-------------|---------|
| `-country` | string | Filter by country name | `-country="United States"` |
| `-owner` | string | Filter by hosting provider | `-owner="Cloudflare, Inc"` |
| `-host` | string | Filter by host | `-host=amazonaws.com` |
| `-dns` | string | Filter by DNS record | `-dns=ns1.cloudflare.com` |
| `-url` | string | Filter by URL content | `-url=wiki` |
| `-rank` | string | Filter by ranking range | `-rank=1-100` |
| `-visitors` | string | Filter by visitor range | `-visitors=1000-10000` |
| `-ip` | string | Filter by IP range/CIDR | `-ip="192.168.0.0/24"` |
| `-output` | string | Output filename | `-output=domains.txt` |
| `-pages` | int | Max pages (0=unlimited) | `-pages=50` |
| `-start` | int | Starting page number | `-start=10` |
| `-proxy` | string | Proxy URL with optional auth | `-proxy=http://proxy.com:8080@user:pass` |
| `-help` | bool | Show help message | `-help` |
| `-list` | bool | List available options | `--list` |

This completes the comprehensive command reference guide. Use this as your go-to resource for all command-line operations with the MyIP.ms Domain Scraper.
