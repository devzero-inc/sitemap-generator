package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type UrlSet struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	URLs    []Url    `xml:"url"`
}

type Url struct {
	Loc string `xml:"loc"`
}

// Custom flag to collect multiple URLs
type urlList []string

func (u *urlList) String() string {
	return strings.Join(*u, ", ")
}

func (u *urlList) Set(value string) error {
	*u = append(*u, value)
	return nil
}

var verbose bool

func main() {
	domain := flag.String("domain", "", "The domain to generate sitemap for")
	var additionalUrls urlList
	flag.Var(&additionalUrls, "add-url", "Additional URL to include in the sitemap (can be repeated)")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose mode")

	flag.Parse()

	if *domain == "" {
		fmt.Println("Please provide a domain using the -domain flag")
		return
	}

	baseURL, err := url.Parse(*domain)
	if err != nil {
		fmt.Println("Invalid domain:", err)
		return
	}

	visited := make(map[string]bool)
	var urls []string

	// Crawl the base URL
	crawl(baseURL, &urls, visited)

	// Add user-specified URLs
	for _, additionalUrl := range additionalUrls {
		fullUrl, err := baseURL.Parse(additionalUrl)
		if err != nil {
			fmt.Println("Invalid URL added via flag:", additionalUrl)
		} else {
			crawl(fullUrl, &urls, visited)
		}
	}

	writeSitemap(urls)
}

func crawl(u *url.URL, urls *[]string, visited map[string]bool) {
	normalizedURL := normalizeURL(u)
	if normalizedURL == "" || visited[normalizedURL] {
		return
	}

	if verbose {
		fmt.Println("Crawling:", normalizedURL)
	}

	visited[normalizedURL] = true
	*urls = append(*urls, normalizedURL)

	resp, err := http.Get(normalizedURL)
	if err != nil {
		if verbose {
			fmt.Println("Error fetching:", normalizedURL, "Error:", err)
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if verbose {
			fmt.Println("Non-OK HTTP status:", resp.StatusCode, "for", normalizedURL)
		}
		return
	}

	if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		if verbose {
			fmt.Println("Skipping non-HTML content at:", normalizedURL)
		}
		return
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		if verbose {
			fmt.Println("Error parsing HTML at:", normalizedURL, "Error:", err)
		}
		return
	}

	visitLinks(u, doc, urls, visited)
}

func visitLinks(baseURL *url.URL, n *html.Node, urls *[]string, visited map[string]bool) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				link, err := baseURL.Parse(attr.Val)
				if err == nil && link.Host == baseURL.Host {
					normalizedLink := normalizeURL(link)
					if normalizedLink != "" && !visited[normalizedLink] {
						if verbose {
							fmt.Println("Found link:", normalizedLink)
						}
						crawl(link, urls, visited)
					}
				}
				break
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		visitLinks(baseURL, c, urls, visited)
	}
}

func writeSitemap(urls []string) {
	urlSet := UrlSet{Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9"}
	for _, u := range urls {
		urlSet.URLs = append(urlSet.URLs, Url{Loc: u})
	}

	file, err := os.Create("sitemap.xml")
	if err != nil {
		fmt.Println("Error creating sitemap.xml:", err)
		return
	}
	defer file.Close()

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	if err := encoder.Encode(urlSet); err != nil {
		fmt.Println("Error encoding sitemap:", err)
		return
	}

	if verbose {
		fmt.Println("Sitemap generated with", len(urls), "URLs.")
	}
}

// normalizeURL ensures that URLs are compared in a consistent manner and strips query parameters.
func normalizeURL(u *url.URL) string {
	// Ignore URLs with query parameters
	if u.RawQuery != "" {
		return ""
	}
	normalized := u.Scheme + "://" + u.Host + u.Path
	// Remove trailing slashes
	normalized = strings.TrimSuffix(normalized, "/")
	return normalized
}
