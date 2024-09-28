# Sitemap Crawler

A robust concurrent sitemap crawler developed in Go using `goquery`. This tool efficiently scrapes URLs from sitemaps while utilizing a variety of random User-Agent strings to minimize detection and scraping limitations. The primary objective of this project is to gather SEO data from competitive websites in order to enhance our own site's performance through strategic implementation of similar techniques.

## Features

- **Concurrent URL Scraping**: Leveraging goroutines for efficient processing.
- **Randomized User-Agent Rotation**: Automatically changes User-Agent for each request to avoid detection.
- **User-Friendly Command-Line Interface**: Simple commands for ease of use.
- **Structured Output**: Results can be exported in JSON format.
- **Configurable Depth**: Customize crawling depth based on your needs.

## Requirements

- Go 1.23 or higher
- `goquery` package

## Installation

1. Clone the repository:

```bash
   mkdir sitemap-crawler
   cd sitemap-crawler
   git clone https://github.com/rohanyh101/go-sitemap_index.xml-crawler .
```

2. Install the required Go packages:

```bash
  go mod tidy
```

## Usage
Run the crawler using the following command:
1. Modify the function `ScrapeSiteMap("https://<url_here>/sitemap_index.xm", p, 10)` in main.go
2. Execute the crawler:

```bash
  go run main.go
```

## Configuration
You can adjust the crawler's settings by editing the constants in the config.go file:

 - MAX_DEPTH: Set the maximum depth for crawling by modifying the variable `n` in the `scrapeURLs` function.
 - USER_AGENTS: Add or remove User-Agent strings for randomization as needed.

## Output flow
1. crawl all `.xml` links
2. gather all URLs associated with each .xml by traversing recursively
3. then get the results of all URLs in `SeoData struct` format

## Contributing
Contributions are welcome! Please fork the repository and create a pull request for any improvements or bug fixes.

<!--
## License
This project is licensed under the MIT License. See the LICENSE file for details.
-->

## Acknowledgments
- `github.com/PuerkitoBio/goquery` for HTML parsing and scraping
- `Goroutines` for concurrent processing
