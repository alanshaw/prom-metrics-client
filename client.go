package client

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// ErrUnexpectedHTTPStatusCode is returned when the HTTP status is not 200
var ErrUnexpectedHTTPStatusCode = fmt.Errorf("unexpected HTTP status code")

// ErrParseFail is returned when parsing of metrics data fails
var ErrParseFail = fmt.Errorf("failed to parse line")

// MetricType describes the type of the metric that was collected. See https://prometheus.io/docs/concepts/metric_types/
type MetricType string

const (
	// Untyped is the metric type of metrics that have no specified type
	Untyped MetricType = ""
	// CounterType is a cumulative metric that represents a single monotonically increasing counter whose value can only increase or be reset to zero on restart.
	CounterType = "counter"
	// GaugeType is a metric that represents a single numerical value that can arbitrarily go up and down.
	GaugeType = "gauge"
	// HistogramType samples observations (usually things like request durations or response sizes) and counts them in configurable buckets.
	HistogramType = "histogram"
	// SummaryType samples observations (usually things like request durations and response sizes). While it also provides a total count of observations and a sum of all observed values, it calculates configurable quantiles over a sliding time window.
	SummaryType = "summary"
)

// Metric is a parsed metric from the endpoint.
type Metric struct {
	Name        string
	Description string
	Type        MetricType
	Samples     []*Sample
}

// HistogramMetric
// SummaryMetric

// Sample is a sample taken for a particular metric.
type Sample struct {
	Name      string
	Labels    map[string]string
	Value     float64
	Timestamp int64
}

// PromMetricsClient is a simple client that fetches and parses metrics from a prometheus /metrics endpoint.
type PromMetricsClient struct {
	URL string
}

// GetMetrics retrieves metrics from the stored URL
func (c *PromMetricsClient) GetMetrics() ([]*Metric, error) {
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

var (
	startHelpLineRE = regexp.MustCompile("^#\\s+HELP\\s+")
	startTypeLineRE = regexp.MustCompile("^#\\s+TYPE\\s+")
	wsRE            = regexp.MustCompile("\\s+")
	labelRE         = regexp.MustCompile("([a-zA-Z_][a-zA-Z0-9_]*)\\s*=\\s*\"((?:\\\\\"|[^\"])*)\"")
)

// Parse reads raw metrics data from the reader, parses it and returns the result
func Parse(r io.Reader) ([]*Metric, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var m *Metric
	var ms []*Metric

	lns := strings.Split(string(b), "\n")

	for i, ln := range lns {
		if ln == "" || isCommentLine(ln) {
			continue
		}

		nm, err := parseLine(m, ln, i)
		if err != nil {
			return nil, err
		}

		if m == nil {
			m = nm
		}

		if nm != m {
			// TODO: upgrade metric to HistogramMetric or SummaryMetric
			ms = append(ms, m)
			m = nm
		}
	}

	if m != nil {
		// TODO: upgrade metric to HistogramMetric or SummaryMetric
		ms = append(ms, m)
	}

	return ms, nil
}

func parseLine(m *Metric, l string, n int) (*Metric, error) {
	if isHelpLine(l) {
		return parseHelpLine(m, l, n)
	}
	if isTypeLine(l) {
		return parseTypeLine(m, l, n)
	}
	return parseSampleLine(m, l, n)
}

func isCommentLine(l string) bool {
	return l[0:1] == "#" && !isHelpLine(l) && !isTypeLine(l)
}

func isHelpLine(l string) bool {
	return startHelpLineRE.MatchString(l)
}

func isTypeLine(l string) bool {
	return startTypeLineRE.MatchString(l)
}

func parseHelpLine(m *Metric, l string, n int) (*Metric, error) {
	l = startHelpLineRE.ReplaceAllString(l, "")
	sp := wsRE.Split(l, 2)

	if m == nil || sp[0] != m.Name {
		m = &Metric{}
	}

	m.Name = sp[0]

	if len(sp) > 1 {
		m.Description = sp[1]
	}

	return m, nil
}

func parseTypeLine(m *Metric, l string, n int) (*Metric, error) {
	l = startTypeLineRE.ReplaceAllString(l, "")
	sp := wsRE.Split(l, 2)

	if len(sp) < 2 {
		return nil, fmt.Errorf("invalid TYPE at line %d: %w", n, ErrParseFail)
	}

	if m == nil || sp[0] != m.Name {
		m = &Metric{}
	}

	m.Name = sp[0]
	m.Type = MetricType(sp[1])

	return m, nil
}

func parseSampleLine(m *Metric, l string, n int) (*Metric, error) {
	var labels string

	fic := strings.Index(l, "{")
	if fic > -1 {
		lic := strings.LastIndex(l, "}")
		labels = l[fic : lic+1]
		l = l[0:fic] + l[lic+1:]
	}

	sp := wsRE.Split(l, -1)

	if len(sp) < 2 {
		return nil, fmt.Errorf("invalid sample at line %d: %w", n, ErrParseFail)
	}

	val, err := strconv.ParseFloat(sp[1], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid float64 value at line %d: %w", n, ErrParseFail)
	}

	s := Sample{
		Name:   sp[0],
		Value:  val,
		Labels: parseLabels(labels),
	}

	if len(sp) >= 3 {
		val, err := strconv.ParseInt(sp[2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid int64 timestamp at line %d: %w", n, ErrParseFail)
		}
		s.Timestamp = val
	}

	if m == nil || strings.Index(s.Name, m.Name) != 0 {
		m = &Metric{Name: s.Name}
	}

	m.Samples = append(m.Samples, &s)

	return m, nil
}

func parseLabels(s string) map[string]string {
	ma := labelRE.FindAllStringSubmatch(s, -1)
	lbs := make(map[string]string)
	for _, sm := range ma {
		lbs[sm[1]] = sm[2]
	}
	return lbs
}
