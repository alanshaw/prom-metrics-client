# prom-metrics-client

[![Build Status](https://img.shields.io/travis/com/libp2p/hydra-booster/master?style=flat-square)](https://travis-ci.org/alanshaw/prom-metrics-client)
[![Coverage](https://img.shields.io/codecov/c/github/alanshaw/prom-metrics-client?style=flat-square)](https://codecov.io/gh/alanshaw/prom-metrics-client)
[![Standard README](https://img.shields.io/badge/readme%20style-standard-brightgreen.svg?style=flat-square)](https://github.com/RichardLitt/standard-readme)

> A simple client that fetches and parses metrics from a prometheus `/metrics` endpoint.

## Install

```console
go get github.com/alanshaw/prom-metrics-client
```

## Usage

⚠️ Currently only supports counters and gauges and ignores keys.

```go
package main

import (
	pmc "github.com/alanshaw/prom-metrics-client"
)

func main() {
    c := pmc.PromMetricsClient{
        URL: "http://localhost:8888/metrics",
    }

    m, _ := c.GetMetrics()

    fmt.Println("Gauges:")
    for _, gauge := range m.Gauges {
		fmt.Printf("%+v\n", gauge)
	}

    fmt.Println("Counters:")
    for _, counter := range m.Counters {
        fmt.Printf("%+v\n", counter)
    }
}
```

## Contribute

Feel free to dive in! [Open an issue](https://github.com/alanshaw/prom-metrics-client/issues/new) or submit PRs.

## License

[MIT](LICENSE) © Alan Shaw
