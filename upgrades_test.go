package client

import (
	"bytes"
	"io/ioutil"
	"math"
	"testing"
)

func TestUpgradeHistogram(t *testing.T) {
	content, err := ioutil.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}

	ms, err := Parse(bytes.NewBuffer(content))
	if err != nil {
		t.Fatal(err)
	}

	m := findMetric(t, ms, "http_request_duration_seconds")
	h, err := UpgradeHistogram(m)
	if err != nil {
		t.Fatal(err)
	}

	if len(h.Buckets) != 6 {
		t.Fatal("incorrect buckets length")
	}

	if h.Buckets[0].LE != 0.05 {
		t.Fatal("incorrect bucket le")
	}

	if h.Buckets[0].Value != 24054 {
		t.Fatal("incorrect bucket value")
	}

	if h.Buckets[1].LE != 0.1 {
		t.Fatal("incorrect bucket le")
	}

	if h.Buckets[1].Value != 33444 {
		t.Fatal("incorrect bucket value")
	}

	if h.Buckets[2].LE != 0.2 {
		t.Fatal("incorrect bucket le")
	}

	if h.Buckets[2].Value != 100392 {
		t.Fatal("incorrect bucket value")
	}

	if h.Buckets[3].LE != 0.5 {
		t.Fatal("incorrect bucket le")
	}

	if h.Buckets[3].Value != 129389 {
		t.Fatal("incorrect bucket value")
	}

	if h.Buckets[4].LE != 1 {
		t.Fatal("incorrect bucket le")
	}

	if h.Buckets[4].Value != 133988 {
		t.Fatal("incorrect bucket value")
	}

	if h.Buckets[5].LE != math.Inf(1) {
		t.Fatal("incorrect bucket le")
	}

	if h.Buckets[5].Value != 144320 {
		t.Fatal("incorrect bucket value")
	}

	if h.Sum != 53423 {
		t.Fatal("incorrect sum")
	}

	if h.Count != 144320 {
		t.Fatal("incorrect count")
	}
}

func TestUpgradeSummary(t *testing.T) {
	content, err := ioutil.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}

	ms, err := Parse(bytes.NewBuffer(content))
	if err != nil {
		t.Fatal(err)
	}

	m := findMetric(t, ms, "rpc_duration_seconds")
	s, err := UpgradeSummary(m)
	if err != nil {
		t.Fatal(err)
	}

	if len(s.Quantiles) != 5 {
		t.Fatal("incorrect quantiles length")
	}

	if s.Quantiles[0].Quantile != 0.01 {
		t.Fatal("incorrect quantile")
	}

	if s.Quantiles[0].Value != 3102 {
		t.Fatal("incorrect quantile value")
	}

	if s.Quantiles[1].Quantile != 0.05 {
		t.Fatal("incorrect quantile")
	}

	if s.Quantiles[1].Value != 3272 {
		t.Fatal("incorrect quantile value")
	}

	if s.Quantiles[2].Quantile != 0.5 {
		t.Fatal("incorrect quantile")
	}

	if s.Quantiles[2].Value != 4773 {
		t.Fatal("incorrect quantile value")
	}

	if s.Quantiles[3].Quantile != 0.9 {
		t.Fatal("incorrect quantile")
	}

	if s.Quantiles[3].Value != 9001 {
		t.Fatal("incorrect quantile value")
	}

	if s.Quantiles[4].Quantile != 0.99 {
		t.Fatal("incorrect quantile")
	}

	if s.Quantiles[4].Value != 76656 {
		t.Fatal("incorrect quantile value")
	}

	if s.Sum != 1.7560473e+07 {
		t.Fatal("incorrect sum")
	}

	if s.Count != 2693 {
		t.Fatal("incorrect count")
	}
}
