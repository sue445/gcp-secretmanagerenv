# gcp-secretmanagerenv
Detect variable from environment variable or [GCP Secret Manager](https://cloud.google.com/secret-manager)

You can access Secret Manager with a syntax similar to `os.Getenv`

[![Build Status](https://github.com/sue445/gcp-secretmanagerenv/workflows/test/badge.svg?branch=master)](https://github.com/sue445/gcp-secretmanagerenv/actions?query=workflow%3Atest)
[![Maintainability](https://api.codeclimate.com/v1/badges/0251ae90c0736a00fdd8/maintainability)](https://codeclimate.com/github/sue445/gcp-secretmanagerenv/maintainability)
[![GoDoc](https://godoc.org/github.com/sue445/gcp-secretmanagerenv?status.svg)](https://godoc.org/github.com/sue445/gcp-secretmanagerenv)
[![Go Report Card](https://goreportcard.com/badge/github.com/sue445/gcp-secretmanagerenv)](https://goreportcard.com/report/github.com/sue445/gcp-secretmanagerenv)

## Usage
```go
package main

import (
    "context"
    "github.com/sue445/gcp-secretmanagerenv"
)

func main() {
    projectID := "gcp-project-id"
    c, err := secretmanagerenv.NewClient(context.Background(), projectID)
    if err != nil {
        panic(err)
    }

    // get from environment variable
    value, err := c.GetValueFromEnvOrSecretManager("SOME_KEY", true)
    // => return value from environment variable or Secret Manager

    // When key is not found in both environment variable and Secret Manager, returned empty string (not error)
    value, err := c.GetValueFromEnvOrSecretManager("INVALID_KEY", false)
    // => ""

    // When key is not found in both environment variable and Secret Manager, returned error
    value, err := c.GetValueFromEnvOrSecretManager("INVALID_KEY", true)
    // => error
}
```

## Development
### Setup
requires https://github.com/direnv/direnv

```bash
cp .envrc.example
vi .envrc
direnv allow
```
