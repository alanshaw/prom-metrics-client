# prom-metrics-client

[![Build Status](https://img.shields.io/travis/com/libp2p/hydra-booster/master?style=flat-square)](https://travis-ci.org/alanshaw/prom-metrics-client)
[![Coverage](https://img.shields.io/codecov/c/github/alanshaw/prom-metrics-client?style=flat-square)](https://codecov.io/gh/alanshaw/prom-metrics-client)
[![Standard README](https://img.shields.io/badge/readme%20style-standard-brightgreen.svg?style=flat-square)](https://github.com/RichardLitt/standard-readme)
[![GoDoc](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/alanshaw/prom-metrics-client)
[![golang version](https://img.shields.io/badge/golang-%3E%3D1.14.0-orange.svg?style=flat-square)](https://golang.org/)

> A simple client that fetches and parses metrics from a prometheus `/metrics` endpoint.

https://prometheus.io/docs/instrumenting/exposition_formats/

## Install

```console
go get github.com/alanshaw/prom-metrics-client
```

## Usage

### Example

```go
package main

import (
	pmc "github.com/alanshaw/prom-metrics-client"
)

func main() {
	c := pmc.PromMetricsClient{
		URL: "http://localhost:8888/metrics",
	}

	ms, _ := c.GetMetrics() // returns []*Metric

	for _, m := range ms {
		// histogram and summary metrics can be upgraded to "richer" types
		if m.Type == HistogramType {
			hm, _ := c.UpgradeHistogram(m)
			fmt.Printf("%+v\n", hm)
		} else if m.Type == SummaryType {
			sm, _ := c.UpgradeSummary(m)
			fmt.Printf("%+v\n", sm)
		} else {
			fmt.Printf("%+v\n", m)
		}
	}
}

/*
Example output:

&{Name:http_requests_total Description:The total number of HTTP requests. Type:counter Samples:[0xc00033e8a0 0xc00033e960]}
&{Name:msdos_file_access_time_seconds Description: Type: Samples:[0xc00033ea20]}
&{Name:metric_without_timestamp_and_labels Description: Type: Samples:[0xc00033eae0]}
&{Name:something_weird Description: Type: Samples:[0xc00033eb40]}
&{Metric:{Name:http_request_duration_seconds Description:A histogram of the request duration. Type:histogram Samples:[0xc00033ebd0 0xc00033ec60 0xc00033ecf0 0xc00033ed80 0xc00033ee10 0xc00033eea0 0xc00033ef30 0xc00033ef90]} Buckets:[0xc000015520 0xc000015530 0xc000015540 0xc000015550 0xc000015560 0xc000015570] Sum:53423 Count:144320}
&{Metric:{Name:rpc_duration_seconds Description:A summary of the RPC duration in seconds. Type:summary Samples:[0xc00033eff0 0xc00033f080 0xc00033f110 0xc00033f1a0 0xc00033f230 0xc00033f2c0 0xc00033f320]} Quantiles:[0xc0000155d0 0xc0000155e0 0xc0000155f0 0xc000015600 0xc000015610] Sum:1.7560473e+07 Count:2693}
*/
```

## API

[GoDoc Reference](https://godoc.org/github.com/alanshaw/prom-metrics-client)

## Contribute

Feel free to dive in! [Open an issue](https://github.com/alanshaw/prom-metrics-client/issues/new) or submit PRs.

## License

[MIT](LICENSE) © Alan Shaw
