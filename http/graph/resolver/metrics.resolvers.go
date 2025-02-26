package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/datarhei/core/http/graph/models"
	"github.com/datarhei/core/http/graph/scalars"
	"github.com/datarhei/core/monitor/metric"
)

func (r *queryResolver) Metrics(ctx context.Context, query models.MetricsInput) (*models.Metrics, error) {
	patterns := []metric.Pattern{}

	for _, m := range query.Metrics {
		labels := []string{}
		for k, v := range m.Labels {
			switch v.(type) {
			case string:
			default:
				continue
			}

			labels = append(labels, k, v.(string))
		}

		pattern := metric.NewPattern(m.Name, labels...)

		patterns = append(patterns, pattern)
	}

	response := &models.Metrics{
		Metrics: []*models.Metric{},
	}

	series := make(map[string]*models.Metric)

	var timeframe int64 = 0
	var interval int64 = 0

	if query.TimerangeSeconds != nil {
		timeframe = int64(*query.TimerangeSeconds)
	}

	if query.IntervalSeconds != nil {
		interval = int64(*query.IntervalSeconds)
	}

	if timeframe == 0 {
		// current data
		now := time.Now()
		data := r.Monitor.Collect(patterns)

		for _, v := range data.All() {
			hash := v.Hash()

			if _, ok := series[hash]; !ok {
				series[hash] = &models.Metric{
					Name:   v.Name(),
					Labels: map[string]interface{}{},
					Values: []*scalars.MetricsResponseValue{},
				}
			}

			k := series[hash]

			for lk, lv := range v.Labels() {
				k.Labels[lk] = lv
			}

			k.Values = append(k.Values, &scalars.MetricsResponseValue{
				TS:    now,
				Value: v.Val(),
			})

			series[hash] = k
		}
	} else {
		// historic data
		data := r.Monitor.History(time.Second*time.Duration(timeframe), time.Second*time.Duration(interval), patterns)

		for _, d := range data {
			if d.Metrics == nil {
				continue
			}

			for _, v := range d.Metrics.All() {
				hash := v.Hash()

				if _, ok := series[hash]; !ok {
					series[hash] = &models.Metric{
						Name:   v.Name(),
						Labels: map[string]interface{}{},
						Values: []*scalars.MetricsResponseValue{},
					}
				}

				k := series[hash]

				for lk, lv := range v.Labels() {
					k.Labels[lk] = lv
				}

				k.Values = append(k.Values, &scalars.MetricsResponseValue{
					TS:    d.TS,
					Value: v.Val(),
				})

				series[hash] = k
			}
		}
	}

	for _, metric := range series {
		response.Metrics = append(response.Metrics, metric)
	}

	resolutionTimerange, resolutionInterval := r.Monitor.Resolution()

	resolutionTimerangeInt := int(resolutionTimerange.Seconds())
	resolutionIntervalInt := int(resolutionInterval.Seconds())

	response.TimerangeSeconds = &resolutionTimerangeInt
	response.IntervalSeconds = &resolutionIntervalInt

	return response, nil
}
