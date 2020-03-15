package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
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

	for _, m := range ms {
		fmt.Printf("%+v\n", m)
	}

	// metric with no HELP or TYPE info should be separate metric
	if m := findMetric(t, ms, "msdos_file_access_time_seconds"); m.Type != Untyped {
		t.Fatal("incorrect metric")
	}
}
