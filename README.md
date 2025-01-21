# Tiny Url Project

## Overview

Tiny URL Service is a web-based application that allows users to shorten long URLs into custom, shorter URLs and retrieve the original URLs associated with those shortened links.

## Features

- **Shorten URL**: Allows users to shorten long URLs into unique, custom short URLs.
- **Fetch Original URL**: Retrieves the original URL associated with a given short URL.
- **Static File Hosting**: Serves a static index.html page.

## Technologies Used

- Go
- net/http
- context
- JSON encoding/decoding
- SQLite/PostgreSQL database
- custom packages for URL generation and database handling

## Installation

1. Clone the repository:
```bash
git clone https://github.com/your-repo/tiny-url.git  
```

2. Navigate to the project directory:
```bash
cd tiny-url
```

3. Install dependencies:
```bash
go mod tidy  
```

4. Configure the application by editing the `config.json` file as needed.

5. Run the application:
```bash
go run main.go  
```

6. Access the service at `http://localhost:8080`.

## Endpoints

- **GET /static/**: Serves static files (e.g., `index.html`).
- **POST /shorten**: Shortens a URL. Requires a JSON payload with the field `"url"`.
- **GET /fetch?short_url={short_url}**: Retrieves the original URL for a given short URL.

## Configuration

The application uses a configuration file located at `./config.json`. The configuration includes database settings and other configurations required for the Tiny URL service.

## License

[MIT](LICENSE)
