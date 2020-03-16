# prom-metrics-client

[![Build Status](https://img.shields.io/travis/com/libp2p/hydra-booster/master?style=flat-square)](https://travis-ci.org/alanshaw/prom-metrics-client)
[![Coverage](https://img.shields.io/codecov/c/github/alanshaw/prom-metrics-client?style=flat-square)](https://codecov.io/gh/alanshaw/prom-metrics-client)
[![Standard README](https://img.shields.io/badge/readme%20style-standard-brightgreen.svg?style=flat-square)](https://github.com/RichardLitt/standard-readme)

> A simple client that fetches and parses metrics from a prometheus `/metrics` endpoint.

https://prometheus.io/docs/instrumenting/exposition_formats/

## Install

```console
go get github.com/alanshaw/prom-metrics-client
```

## Usage

⚠️ Currently has no special handling for histogram or summary types.

```go
package main

import (
	pmc "github.com/alanshaw/prom-metrics-client"
)

func main() {
	c := pmc.PromMetricsClient{
		URL: "http://localhost:8888/metrics",
	}

	ms, _ := c.GetMetrics()

	for _, m := range ms {
		fmt.Printf("%+v\n", m)
	}
}

/*
Example output:
&{Name:http_requests_total Description:The total number of HTTP requests. Type:counter Samples:[0xc000132090 0xc0001320c0]}
&{Name:msdos_file_access_time_seconds Description: Type: Samples:[0xc0001320f0]}
&{Name:metric_without_timestamp_and_labels Description: Type: Samples:[0xc000132120]}
&{Name:something_weird Description: Type: Samples:[0xc000132150]}
&{Name:http_request_duration_seconds Description:A histogram of the request duration. Type:histogram Samples:[0xc000132180 0xc0001321b0 0xc0001321e0 0xc000132210 0xc000132240 0xc000132270 0xc0001322a0 0xc0001322d0]}
&{Name:rpc_duration_seconds Description:A summary of the RPC duration in seconds. Type:summary Samples:[0xc000132300 0xc000132330 0xc000132360 0xc000132390 0xc0001323c0 0xc0001323f0 0xc000132420]}
*/
```

## Contribute

Feel free to dive in! [Open an issue](https://github.com/alanshaw/prom-metrics-client/issues/new) or submit PRs.

## License

[MIT](LICENSE) © Alan Shaw
