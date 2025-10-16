package metric

import (
	model "github.com/lin-snow/ech0/internal/model/metric"
)

type MetricCollector interface {
	Collect() (model.Metrics, error)
	Reset() error
}
