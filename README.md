# Sitemap Generator

## Overview
This project is a Go application designed to generate sitemaps for websites. It crawls the specified domain, normalizes URLs, and generates a sitemap in XML format. The application also supports adding additional URLs and provides a verbose mode for detailed output.

## Features
Sitemap Generation: Automatically generates sitemaps for websites.
URL Normalization: Removes trailing slashes and ignores URLs with query parameters to maintain consistency.
Additional URLs: Allows inclusion of additional URLs in the sitemap.
Verbose Mode: Provides detailed output for debugging and monitoring.

## Prerequisites
Go programming language (version 1.22 or later)

## Installation

1. Clone the repository
```
git clone https://github.com/yourusername/sitemap-generator.git
cd sitemap-generator
```

2. Build the project

```
go build -o sitemap-generator
```

### Usage
1. Run the binary:

```
./sitemap-generator -domain https://example.com
```

2. Command-line options:

 - `-domain`: The base domain of the website for which you want to generate the sitemap.
 - `-add-url`: Additional URL to include in the sitemap (can be repeated).
 - `-verbose`: Enable verbose mode for detailed output.

### Example
To generate a sitemap for https://example.com and include additional URLs, run:

```
./sitemap-generator -domain https://example.com -add-url /page1 -add-url /page2 -verbose
```

## License
This project is licensed under the MIT License - see the LICENSE file for details.
