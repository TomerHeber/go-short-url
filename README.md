[![Go Report Card](https://goreportcard.com//badge/github.com/TomerHeber/go-short-url)](https://goreportcard.com/report/github.com/TomerHeber/go-short-url)

# go-short-url
A URL shortener library implemented in Go.

## Installation

```
go get github.com/TomerHeber/go-short-url
```

## Documentation

Documentation is available at: [https://pkg.go.dev/github.com/TomerHeber/go-short-url](https://pkg.go.dev/github.com/TomerHeber/go-short-url)

## Example

Below is a basic example (Note: errors are not checked):

```
// Create a new shortner.
s, _ := short.NewShortener()

// Create a shortened url for "https://www.google.com?a=b&c=d". 
surl, _ := s.CreateShortenedUrl(context.TODO(), "https://www.google.com?a=b&c=d")

// Return the original url from the shortened url.
ourl, _ := s.GetUrlFromShortenedUrl(context.TODO(), surl.GetUrl())
```

## Running the webserver example

A more complete example is available under the example directory.
It's a tiny webserver that generates short urls.
  
```
docker-compose up -d
go run ./example/main.go
```

## MongoDB

To store the mappings between short urls and their original urls, A MongoDB database is required.

## Development

Install `golangci-lint`:  
Check [https://golangci-lint.run/usage/install/#local-installation](https://golangci-lint.run/usage/install/#local-installation) for installation instructions.

Install `pre-commit`:  
Check [https://pre-commit.com/](https://pre-commit.com/) for installation instructions.

Enable the git pre commit hooks:  
`pre-commit install `