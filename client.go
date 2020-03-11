package client

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// ErrUnexpectedHTTPStatusCode is returned when the HTTP status is not 200
var ErrUnexpectedHTTPStatusCode = fmt.Errorf("unexpected HTTP status code")

// ErrParseFail is returned when parsing of metrics data fails
var ErrParseFail = fmt.Errorf("failed to parse line")

// Metrics are the parsed metrics from the endpoint.
type Metrics struct {
	Gauges   []Gauge
	Counters []Counter
}

// Counter is a cumulative metric that represents a single monotonically increasing counter whose value can only increase or be reset to zero on restart.
type Counter struct {
	Name        string
	Description string
	Value       float64
}

// Gauge is a metric that represents a single numerical value that can arbitrarily go up and down.
type Gauge struct {
	Name        string
	Description string
	Value       float64
}

// Histogram samples observations (usually things like request durations or response sizes) and counts them in configurable buckets.
// type Histogram struct {
// 	Description string
// }

// Summary samples observations (usually things like request durations and response sizes). While it also provides a total count of observations and a sum of all observed values, it calculates configurable quantiles over a sliding time window.
// type Summary struct {
// 	Description string
// }

// PromMetricsClient is a simple client that fetches and parses metrics from a prometheus /metrics endpoint.
type PromMetricsClient struct {
	URL string
}

// GetMetrics retrieves metrics from the address
func (c *PromMetricsClient) GetMetrics() (*Metrics, error) {
	res, err := http.Get(c.URL)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, ErrUnexpectedHTTPStatusCode
	}

	defer res.Body.Close()
	return Parse(res.Body)
}

// Parse reads raw metrics data from the reader, parses it and returns the result
func Parse(r io.Reader) (*Metrics, error) {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var metrics Metrics
	lines := strings.Split(strings.Trim(string(bytes), "\n\r"), "\n")
	llen := len(lines)

	i := 0
	for {
		if i >= llen {
			break
		}

		helpl, err := parseHelpLine(lines[i], i)
		if err != nil {
			return nil, err
		}

		i++
		if i >= llen {
			break
		}

		typel, err := parseTypeLine(lines[i], i)
		if err != nil {
			return nil, err
		}

		i++
		if i >= llen {
			break
		}

		switch typel.Type {
		case "counter":
			counterl, err := parseCounterLine(lines[i], i)
			if err != nil {
				return nil, err
			}

			metrics.Counters = append(metrics.Counters, Counter{
				Name:        helpl.Name,
				Description: helpl.Description,
				Value:       counterl.Value,
			})

			i++
		case "gauge":
			gaugel, err := parseGaugeLine(lines[i], i)
			if err != nil {
				return nil, err
			}

			metrics.Gauges = append(metrics.Gauges, Gauge{
				Name:        helpl.Name,
				Description: helpl.Description,
				Value:       gaugel.Value,
			})

			i++
		default:
			// Currently unsupported
			for {
				if isHashLine(lines[i]) {
					break
				}
				i++
				if i >= llen {
					break
				}
			}
		}
	}

	return &metrics, nil
}

func isHashLine(l string) bool {
	return l[0:1] == "#"
}

type helpLine struct {
	Name        string
	Description string
}

func parseHelpLine(l string, n int) (*helpLine, error) {
	// fmt.Println("parseHelpLine", l)
	l = strings.Trim(l, "\n\r")

	if strings.Index(l, "# HELP ") != 0 {
		return nil, fmt.Errorf("expecting '# HELP ' at line %d: %w", n, ErrParseFail)
	}

	l = l[7:]
	si := strings.Index(l, " ")
	if si == -1 {
		return nil, fmt.Errorf("missing separator at line %d: %w", n, ErrParseFail)
	}

	name := l[0:si]
	desc := l[si+1:]

	return &helpLine{Name: name, Description: desc}, nil
}

type typeLine struct {
	Name string
	Type string
}

func parseTypeLine(l string, n int) (*typeLine, error) {
	// fmt.Println("parseTypeLine", l)
	l = strings.Trim(l, "\n\r")

	if strings.Index(l, "# TYPE ") != 0 {
		return nil, fmt.Errorf("expecting '# TYPE ' at line %d: %w", n, ErrParseFail)
	}

	l = l[7:]
	si := strings.Index(l, " ")
	if si == -1 {
		return nil, fmt.Errorf("missing separator at line %d: %w", n, ErrParseFail)
	}

	name := l[0:si]
	typ := l[si+1:]

	return &typeLine{Name: name, Type: typ}, nil
}

type counterLine struct {
	Name  string
	Value float64
}

func parseCounterLine(l string, n int) (*counterLine, error) {
	// fmt.Println("parseCounterLine", l)
	l = strings.Trim(l, "\n\r")
	split := strings.Split(l, " ")

	if len(split) != 2 {
		return nil, fmt.Errorf("invalid counter at line %d: %w", n, ErrParseFail)
	}

	val, err := strconv.ParseFloat(split[1], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid float64 at line %d: %w", n, ErrParseFail)
	}

	return &counterLine{Name: split[0], Value: val}, nil
}

type gaugeLine struct {
	Name  string
	Value float64
}

func parseGaugeLine(l string, n int) (*gaugeLine, error) {
	// fmt.Println("parseGaugeLine", l)
	l = strings.Trim(l, "\n\r")
	split := strings.Split(l, " ")

	if len(split) != 2 {
		return nil, fmt.Errorf("invalid gauge at line %d: %w", n, ErrParseFail)
	}

	val, err := strconv.ParseFloat(split[1], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid float64 at line %d: %w", n, ErrParseFail)
	}

	return &gaugeLine{Name: split[0], Value: val}, nil
}
