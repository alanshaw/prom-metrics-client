package client

import (
	"fmt"
	"strconv"
)

// HistogramMetric is a metric of type HistogramType
type HistogramMetric struct {
	Metric
	Buckets []*HistogramBucket
	Sum     float64
	Count   float64
}

// HistogramBucket define the upper bound and value
type HistogramBucket struct {
	LE    float64
	Value float64
}

// SummaryMetric is a metric of type SummaryType
type SummaryMetric struct {
	Metric
	Quantiles []*SummaryQuantile
	Sum       float64
	Count     float64
}

// SummaryQuantile defines the cut point and value
type SummaryQuantile struct {
	Quantile float64
	Value    float64
}

func findSamplesByName(m *Metric, n string) []*Sample {
	var ss []*Sample
	for _, s := range m.Samples {
		if s.Name == n {
			ss = append(ss, s)
		}
	}
	return ss
}

// UpgradeHistogram upgrades a Metric of type HistorgramType to a HistogramMetric
func UpgradeHistogram(m *Metric) (*HistogramMetric, error) {
	if m.Type != HistogramType {
		return nil, fmt.Errorf("metric is not a histogram")
	}

	sbs := findSamplesByName(m, m.Name+"_bucket")
	var hbs []*HistogramBucket
	for _, sb := range sbs {
		le, err := strconv.ParseFloat(sb.Labels["le"], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid le float64 value: %w", err)
		}
		hbs = append(hbs, &HistogramBucket{
			LE:    le,
			Value: sb.Value,
		})
	}

	ss := findSamplesByName(m, m.Name+"_sum")
	if len(ss) != 1 {
		return nil, fmt.Errorf("missing or multiple sum sample(s)")
	}

	cs := findSamplesByName(m, m.Name+"_count")
	if len(cs) != 1 {
		return nil, fmt.Errorf("missing or multiple count sample(s)")
	}

	return &HistogramMetric{
		Metric: Metric{
			Name:        m.Name,
			Description: m.Description,
			Type:        m.Type,
			Samples:     m.Samples,
		},
		Buckets: hbs,
		Sum:     ss[0].Value,
		Count:   cs[0].Value,
	}, nil
}

// UpgradeSummary upgrades a Metric of type SummaryType to a SummaryMetric
func UpgradeSummary(m *Metric) (*SummaryMetric, error) {
	if m.Type != SummaryType {
		return nil, fmt.Errorf("metric is not a summary")
	}

	sqs := findSamplesByName(m, m.Name)
	var qs []*SummaryQuantile
	for _, sb := range sqs {
		q, err := strconv.ParseFloat(sb.Labels["quantile"], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid quantile float64 value: %w", err)
		}
		qs = append(qs, &SummaryQuantile{
			Quantile: q,
			Value:    sb.Value,
		})
	}

	ss := findSamplesByName(m, m.Name+"_sum")
	if len(ss) != 1 {
		return nil, fmt.Errorf("missing or multiple sum sample(s)")
	}

	cs := findSamplesByName(m, m.Name+"_count")
	if len(cs) != 1 {
		return nil, fmt.Errorf("missing or multiple count sample(s)")
	}

	return &SummaryMetric{
		Metric: Metric{
			Name:        m.Name,
			Description: m.Description,
			Type:        m.Type,
			Samples:     m.Samples,
		},
		Quantiles: qs,
		Sum:       ss[0].Value,
		Count:     cs[0].Value,
	}, nil
}
