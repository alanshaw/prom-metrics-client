package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"testing"
)

func findMetric(t *testing.T, ms []*Metric, name string) *Metric {
	for _, m := range ms {
		if m.Name == name {
			return m
		}
	}
	t.Fatal("missing metric", name)
	return nil
}

func TestIsHelpLine(t *testing.T) {
	if !isHelpLine("# HELP ") {
		t.Fatal("expected to be help line")
	}
}

func TestParseMemstatsTxt(t *testing.T) {
	content, err := ioutil.ReadFile("testdata/memstats.txt")
	if err != nil {
		t.Fatal(err)
	}

	ms, err := Parse(bytes.NewBuffer(content))
	if err != nil {
		t.Fatal(err)
	}

	if m := findMetric(t, ms, "go_goroutines"); m.Samples[0].Value != 166 {
		t.Fatal("incorrect gauge value")
	}

	if m := findMetric(t, ms, "go_memstats_heap_objects"); m.Samples[0].Value != 33814 {
		t.Fatal("incorrect gauge value")
	}

	if m := findMetric(t, ms, "go_memstats_mallocs_total"); m.Samples[0].Value != 44939 {
		t.Fatal("incorrect counter value")
	}
}

func TestGetMemstatsTxt(t *testing.T) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, "testdata/memstats.txt")
	})

	go http.Serve(listener, mux)
	defer listener.Close()

	c := PromMetricsClient{
		URL: fmt.Sprintf("http://%v/metrics", listener.Addr().String()),
	}

	ms, err := c.GetMetrics()
	if err != nil {
		t.Fatal(err)
	}

	if m := findMetric(t, ms, "go_goroutines"); m.Samples[0].Value != 166 {
		t.Fatal("incorrect gauge value")
	}

	if m := findMetric(t, ms, "go_memstats_heap_objects"); m.Samples[0].Value != 33814 {
		t.Fatal("incorrect gauge value")
	}

	if m := findMetric(t, ms, "go_memstats_mallocs_total"); m.Samples[0].Value != 44939 {
		t.Fatal("incorrect counter value")
	}
}

func TestParseMultisampleTxt(t *testing.T) {
	content, err := ioutil.ReadFile("testdata/multisample.txt")
	if err != nil {
		t.Fatal(err)
	}

	ms, err := Parse(bytes.NewBuffer(content))
	if err != nil {
		t.Fatal(err)
	}

	m := findMetric(t, ms, "hydrabooster_connected_peers")
	if len(m.Samples) != 2 {
		t.Fatal("incorrect counter values len")
	}

	if m.Samples[0].Value != 12 {
		t.Fatal("incorrect counter value [0]")
	}

	if m.Samples[1].Value != 6 {
		t.Fatal("incorrect counter value [1]")
	}
}

func TestParseExampleTxt(t *testing.T) {
	content, err := ioutil.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}

	ms, err := Parse(bytes.NewBuffer(content))
	if err != nil {
		t.Fatal(err)
	}

	var m *Metric
	m = findMetric(t, ms, "http_requests_total")

	if m.Name != "http_requests_total" {
		t.Fatal("incorrect name")
	}

	if m.Description != "The total number of HTTP requests." {
		t.Fatal("incorrect description")
	}

	if m.Type != CounterType {
		t.Fatal("incorrect type")
	}

	if len(m.Samples) != 2 {
		t.Fatal("incorrect samples length")
	}

	if m.Samples[0].Name != "http_requests_total" {
		t.Fatal("incorrect sample name")
	}

	if len(m.Samples[0].Labels) != 2 {
		t.Fatal("incorrect sample labels length")
	}

	if m.Samples[0].Labels["method"] != "post" {
		t.Fatal("incorrect sample label value")
	}

	if m.Samples[0].Labels["code"] != "200" {
		t.Fatal("incorrect sample label value")
	}

	if m.Samples[0].Value != 1027 {
		t.Fatal("incorrect sample value")
	}

	if m.Samples[0].Timestamp != 1395066363000 {
		t.Fatal("incorrect sample timestamp")
	}

	if m.Samples[1].Name != "http_requests_total" {
		t.Fatal("incorrect sample name")
	}

	if len(m.Samples[1].Labels) != 2 {
		t.Fatal("incorrect sample labels length")
	}

	if m.Samples[1].Labels["method"] != "post" {
		t.Fatal("incorrect sample label value")
	}

	if m.Samples[1].Labels["code"] != "400" {
		t.Fatal("incorrect sample label value")
	}

	if m.Samples[1].Value != 3 {
		t.Fatal("incorrect sample value")
	}

	if m.Samples[1].Timestamp != 1395066363000 {
		t.Fatal("incorrect sample timestamp")
	}

	m = findMetric(t, ms, "msdos_file_access_time_seconds")

	if m.Name != "msdos_file_access_time_seconds" {
		t.Fatal("incorrect metric name")
	}

	if m.Description != "" {
		t.Fatal("incorrect metric description")
	}

	if m.Type != Untyped {
		t.Fatal("incorrect metric type")
	}

	if len(m.Samples) != 1 {
		t.Fatal("incorrect samples length")
	}

	if m.Samples[0].Name != "msdos_file_access_time_seconds" {
		t.Fatal("incorrect sample name")
	}

	if len(m.Samples[0].Labels) != 2 {
		t.Fatal("incorrect sample labels length")
	}

	if m.Samples[0].Labels["path"] != "C:\\\\DIR\\\\FILE.TXT" {
		t.Fatal("incorrect sample label value")
	}

	if m.Samples[0].Labels["error"] != "Cannot find file:\\n\\\"FILE.TXT\\\"" {
		t.Fatal("incorrect sample label value")
	}

	if m.Samples[0].Value != 1.458255915e9 {
		t.Fatal("incorrect sample value")
	}

	if m.Samples[0].Timestamp != 0 {
		t.Fatal("incorrect sample timestamp")
	}

	m = findMetric(t, ms, "metric_without_timestamp_and_labels")

	if m.Name != "metric_without_timestamp_and_labels" {
		t.Fatal("incorrect metric name")
	}

	if m.Description != "" {
		t.Fatal("incorrect metric description")
	}

	if m.Type != Untyped {
		t.Fatal("incorrect metric type")
	}

	if len(m.Samples) != 1 {
		t.Fatal("incorrect samples length")
	}

	if m.Samples[0].Name != "metric_without_timestamp_and_labels" {
		t.Fatal("incorrect sample name")
	}

	if len(m.Samples[0].Labels) != 0 {
		t.Fatal("incorrect sample labels length")
	}

	if m.Samples[0].Value != 12.47 {
		t.Fatal("incorrect sample value")
	}

	if m.Samples[0].Timestamp != 0 {
		t.Fatal("incorrect sample timestamp")
	}

	m = findMetric(t, ms, "something_weird")

	if m.Name != "something_weird" {
		t.Fatal("incorrect metric name")
	}

	if m.Description != "" {
		t.Fatal("incorrect metric description")
	}

	if m.Type != Untyped {
		t.Fatal("incorrect metric type")
	}

	if len(m.Samples) != 1 {
		t.Fatal("incorrect samples length")
	}

	if m.Samples[0].Name != "something_weird" {
		t.Fatal("incorrect sample name")
	}

	if len(m.Samples[0].Labels) != 1 {
		t.Fatal("incorrect sample labels length")
	}

	if m.Samples[0].Labels["problem"] != "division by zero" {
		t.Fatal("incorrect sample label value")
	}

	if m.Samples[0].Value != math.Inf(1) {
		t.Fatal("incorrect sample value")
	}

	if m.Samples[0].Timestamp != -3982045 {
		t.Fatal("incorrect sample timestamp")
	}

	m = findMetric(t, ms, "http_request_duration_seconds")

	if m.Name != "http_request_duration_seconds" {
		t.Fatal("incorrect name")
	}

	if m.Description != "A histogram of the request duration." {
		t.Fatal("incorrect description")
	}

	if m.Type != HistogramType {
		t.Fatal("incorrect type")
	}

	if len(m.Samples) != 8 {
		t.Fatal("incorrect samples length")
	}

	if m.Samples[0].Name != "http_request_duration_seconds_bucket" {
		t.Fatal("incorrect sample name")
	}

	if len(m.Samples[0].Labels) != 1 {
		t.Fatal("incorrect sample labels length")
	}

	if m.Samples[0].Labels["le"] != "0.05" {
		t.Fatal("incorrect sample label value")
	}

	if m.Samples[0].Value != 24054 {
		t.Fatal("incorrect sample value")
	}

	if m.Samples[0].Timestamp != 0 {
		t.Fatal("incorrect sample timestamp")
	}

	if m.Samples[1].Name != "http_request_duration_seconds_bucket" {
		t.Fatal("incorrect sample name")
	}

	if len(m.Samples[1].Labels) != 1 {
		t.Fatal("incorrect sample labels length")
	}

	if m.Samples[1].Labels["le"] != "0.1" {
		t.Fatal("incorrect sample label value")
	}

	if m.Samples[1].Value != 33444 {
		t.Fatal("incorrect sample value")
	}

	if m.Samples[1].Timestamp != 0 {
		t.Fatal("incorrect sample timestamp")
	}

	if m.Samples[2].Name != "http_request_duration_seconds_bucket" {
		t.Fatal("incorrect sample name")
	}

	if len(m.Samples[2].Labels) != 1 {
		t.Fatal("incorrect sample labels length")
	}

	if m.Samples[2].Labels["le"] != "0.2" {
		t.Fatal("incorrect sample label value")
	}

	if m.Samples[2].Value != 100392 {
		t.Fatal("incorrect sample value")
	}

	if m.Samples[2].Timestamp != 0 {
		t.Fatal("incorrect sample timestamp")
	}

	if m.Samples[3].Name != "http_request_duration_seconds_bucket" {
		t.Fatal("incorrect sample name")
	}

	if len(m.Samples[3].Labels) != 1 {
		t.Fatal("incorrect sample labels length")
	}

	if m.Samples[3].Labels["le"] != "0.5" {
		t.Fatal("incorrect sample label value")
	}

	if m.Samples[3].Value != 129389 {
		t.Fatal("incorrect sample value")
	}

	if m.Samples[3].Timestamp != 0 {
		t.Fatal("incorrect sample timestamp")
	}

	if m.Samples[4].Name != "http_request_duration_seconds_bucket" {
		t.Fatal("incorrect sample name")
	}

	if len(m.Samples[4].Labels) != 1 {
		t.Fatal("incorrect sample labels length")
	}

	if m.Samples[4].Labels["le"] != "1" {
		t.Fatal("incorrect sample label value")
	}

	if m.Samples[4].Value != 133988 {
		t.Fatal("incorrect sample value")
	}

	if m.Samples[4].Timestamp != 0 {
		t.Fatal("incorrect sample timestamp")
	}

	if m.Samples[5].Name != "http_request_duration_seconds_bucket" {
		t.Fatal("incorrect sample name")
	}

	if len(m.Samples[5].Labels) != 1 {
		t.Fatal("incorrect sample labels length")
	}

	if m.Samples[5].Labels["le"] != "+Inf" {
		t.Fatal("incorrect sample label value")
	}

	if m.Samples[5].Value != 144320 {
		t.Fatal("incorrect sample value")
	}

	if m.Samples[5].Timestamp != 0 {
		t.Fatal("incorrect sample timestamp")
	}

	if m.Samples[6].Name != "http_request_duration_seconds_sum" {
		t.Fatal("incorrect sample name")
	}

	if len(m.Samples[6].Labels) != 0 {
		t.Fatal("incorrect sample labels length")
	}

	if m.Samples[6].Value != 53423 {
		t.Fatal("incorrect sample value")
	}

	if m.Samples[6].Timestamp != 0 {
		t.Fatal("incorrect sample timestamp")
	}

	if m.Samples[7].Name != "http_request_duration_seconds_count" {
		t.Fatal("incorrect sample name")
	}

	if len(m.Samples[7].Labels) != 0 {
		t.Fatal("incorrect sample labels length")
	}

	if m.Samples[7].Value != 144320 {
		t.Fatal("incorrect sample value")
	}

	if m.Samples[7].Timestamp != 0 {
		t.Fatal("incorrect sample timestamp")
	}

	m = findMetric(t, ms, "rpc_duration_seconds")

	if m.Name != "rpc_duration_seconds" {
		t.Fatal("incorrect name")
	}

	if m.Description != "A summary of the RPC duration in seconds." {
		t.Fatal("incorrect description")
	}

	if m.Type != SummaryType {
		t.Fatal("incorrect type")
	}

	if len(m.Samples) != 7 {
		t.Fatal("incorrect samples length")
	}

	if m.Samples[0].Name != "rpc_duration_seconds" {
		t.Fatal("incorrect sample name")
	}

	if len(m.Samples[0].Labels) != 1 {
		t.Fatal("incorrect sample labels length")
	}

	if m.Samples[0].Labels["quantile"] != "0.01" {
		t.Fatal("incorrect sample label value")
	}

	if m.Samples[0].Value != 3102 {
		t.Fatal("incorrect sample value")
	}

	if m.Samples[0].Timestamp != 0 {
		t.Fatal("incorrect sample timestamp")
	}

	if m.Samples[1].Name != "rpc_duration_seconds" {
		t.Fatal("incorrect sample name")
	}

	if len(m.Samples[1].Labels) != 1 {
		t.Fatal("incorrect sample labels length")
	}

	if m.Samples[1].Labels["quantile"] != "0.05" {
		t.Fatal("incorrect sample label value")
	}

	if m.Samples[1].Value != 3272 {
		t.Fatal("incorrect sample value")
	}

	if m.Samples[1].Timestamp != 0 {
		t.Fatal("incorrect sample timestamp")
	}

	if m.Samples[2].Name != "rpc_duration_seconds" {
		t.Fatal("incorrect sample name")
	}

	if len(m.Samples[2].Labels) != 1 {
		t.Fatal("incorrect sample labels length")
	}

	if m.Samples[2].Labels["quantile"] != "0.5" {
		t.Fatal("incorrect sample label value")
	}

	if m.Samples[2].Value != 4773 {
		t.Fatal("incorrect sample value")
	}

	if m.Samples[2].Timestamp != 0 {
		t.Fatal("incorrect sample timestamp")
	}

	if m.Samples[3].Name != "rpc_duration_seconds" {
		t.Fatal("incorrect sample name")
	}

	if len(m.Samples[3].Labels) != 1 {
		t.Fatal("incorrect sample labels length")
	}

	if m.Samples[3].Labels["quantile"] != "0.9" {
		t.Fatal("incorrect sample label value")
	}

	if m.Samples[3].Value != 9001 {
		t.Fatal("incorrect sample value")
	}

	if m.Samples[3].Timestamp != 0 {
		t.Fatal("incorrect sample timestamp")
	}

	if m.Samples[4].Name != "rpc_duration_seconds" {
		t.Fatal("incorrect sample name")
	}

	if len(m.Samples[4].Labels) != 1 {
		t.Fatal("incorrect sample labels length")
	}

	if m.Samples[4].Labels["quantile"] != "0.99" {
		t.Fatal("incorrect sample label value")
	}

	if m.Samples[4].Value != 76656 {
		t.Fatal("incorrect sample value")
	}

	if m.Samples[4].Timestamp != 0 {
		t.Fatal("incorrect sample timestamp")
	}

	if m.Samples[5].Name != "rpc_duration_seconds_sum" {
		t.Fatal("incorrect sample name")
	}

	if len(m.Samples[5].Labels) != 0 {
		t.Fatal("incorrect sample labels length")
	}

	if m.Samples[5].Value != 1.7560473e+07 {
		t.Fatal("incorrect sample value")
	}

	if m.Samples[5].Timestamp != 0 {
		t.Fatal("incorrect sample timestamp")
	}

	if m.Samples[6].Name != "rpc_duration_seconds_count" {
		t.Fatal("incorrect sample name")
	}

	if len(m.Samples[6].Labels) != 0 {
		t.Fatal("incorrect sample labels length")
	}

	if m.Samples[6].Value != 2693 {
		t.Fatal("incorrect sample value")
	}

	if m.Samples[6].Timestamp != 0 {
		t.Fatal("incorrect sample timestamp")
	}
}
