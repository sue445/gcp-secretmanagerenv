# gcp-secretmanagerenv
Detect variable from environment variable or [GCP Secret Manager](https://cloud.google.com/secret-manager)

You can access Secret Manager with a syntax similar to `os.Getenv`

[![Latest Version](https://img.shields.io/github/v/tag/sue445/gcp-secretmanagerenv)](https://github.com/sue445/gcp-secretmanagerenv/tags)
[![Build Status](https://github.com/sue445/gcp-secretmanagerenv/workflows/test/badge.svg?branch=master)](https://github.com/sue445/gcp-secretmanagerenv/actions?query=workflow%3Atest)
[![Coverage Status](https://coveralls.io/repos/github/sue445/gcp-secretmanagerenv/badge.svg)](https://coveralls.io/github/sue445/gcp-secretmanagerenv)
[![Maintainability](https://api.codeclimate.com/v1/badges/0251ae90c0736a00fdd8/maintainability)](https://codeclimate.com/github/sue445/gcp-secretmanagerenv/maintainability)
[![GoDoc](https://godoc.org/github.com/sue445/gcp-secretmanagerenv?status.svg)](https://godoc.org/github.com/sue445/gcp-secretmanagerenv)
[![Go Report Card](https://goreportcard.com/badge/github.com/sue445/gcp-secretmanagerenv)](https://goreportcard.com/report/github.com/sue445/gcp-secretmanagerenv)

## Requirements
Add IAM role `roles/secretmanager.secretAccessor` to service account if necessary.

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

### Specification
When `c.GetValueFromEnvOrSecretManager(key, required)` is called, processing is performed in the following order

1. Returns environment variable if `key` is found
2. Returns latest version value of Secret Manager if `projectID` isn't empty and `key` is found
3. Returns `""` if `required == false`
4. Returns `error` if `required == true`

## Development
### Setup
requires https://github.com/direnv/direnv

```bash
cp .envrc.example
vi .envrc
direnv allow
```
