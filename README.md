README :: GoPee - Excel Formula Execution WebUI with GoLang
====
[![Build Status](https://travis-ci.org/khanhhua/gopee.svg?branch=master)](https://travis-ci.org/khanhhua/gopee)

## Installation

1. Install Go 1.10.x or greater, git, setup `$GOPATH`, and `PATH=$PATH:$GOPATH/bin`

2. Run the server
    ```
    cd $GOPATH/src/github.com/khanhhua/gopee
    go run main.go
    ```


## Environment Variables for Configuration

* **PORT:** The port. Default: `"8888"`

* **HTTP_CERT_FILE:** Path to cert file. Default: `""`

* **HTTP_KEY_FILE:** Path to key file. Default: `""`

* **HTTP_DRAIN_INTERVAL:** How long application will wait to drain old requests before restarting. Default: `"1s"`

* **COOKIE_SECRET:** Cookie secret for session. Default: Auto generated.
