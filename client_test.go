package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
)

func findGauge(gauges []Gauge, name string) *Gauge {
	for _, g := range gauges {
		if g.Name == name {
			return &g
		}
	}
	return nil
}

func findCounter(counters []Counter, name string) *Counter {
	for _, c := range counters {
		if c.Name == name {
			return &c
		}
	}
	return nil
}

func TestParse(t *testing.T) {
	content, err := ioutil.ReadFile("testdata/sample0.txt")
	if err != nil {
		t.Fatal(err)
	}

	m, err := Parse(bytes.NewBuffer(content))
	if err != nil {
		t.Fatal(err)
	}

	if g := findGauge(m.Gauges, "go_goroutines"); g.Value != 166 {
		t.Fatal("incorrect gauge value")
	}

	if g := findGauge(m.Gauges, "go_memstats_heap_objects"); g.Value != 33814 {
		t.Fatal("incorrect gauge value")
	}

	if c := findCounter(m.Counters, "go_memstats_mallocs_total"); c.Value != 44939 {
		t.Fatal("incorrect counter value")
	}
}

func TestGet(t *testing.T) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, "testdata/sample0.txt")
	})

	go http.Serve(listener, mux)
	defer listener.Close()

	c := PromMetricsClient{
		URL: fmt.Sprintf("http://%v/metrics", listener.Addr().String()),
	}

	m, err := c.GetMetrics()
	if err != nil {
		t.Fatal(err)
	}

	if g := findGauge(m.Gauges, "go_goroutines"); g.Value != 166 {
		t.Fatal("incorrect gauge value")
	}

	if g := findGauge(m.Gauges, "go_memstats_heap_objects"); g.Value != 33814 {
		t.Fatal("incorrect gauge value")
	}

	if c := findCounter(m.Counters, "go_memstats_mallocs_total"); c.Value != 44939 {
		t.Fatal("incorrect counter value")
	}
}
