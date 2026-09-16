package sqlstreams

import (
	"github.com/allegedlyreliable/sqlstreams/pkg/system"
)

// please GOPLS make aliases and go doc comments work better

// SystemConfig declares the built-in alert settings and metrics collector's
// poll rate. Every field is optional.
type SystemConfig struct {
	// PartitionCountAlert - the partition_count alert declaration.
	// Default: its own defaults.
	PartitionCountAlert *PartitionCountAlertConfig

	// CompactionReadCostAlert - the compaction_read_cost alert declaration.
	// Default: its own defaults.
	CompactionReadCostAlert *CompactionReadCostAlertConfig

	// WorkerLivenessAlert - the worker_liveness alert declaration.
	// Default: its own defaults.
	WorkerLivenessAlert *WorkerLivenessAlertConfig

	// MetricCollectorProgressAlert declares the collector-progress check. Default: its own defaults.
	MetricCollectorProgressAlert *MetricCollectorProgressAlertConfig

	// MetricCollector - the metrics_collector worker declaration.
	// Default: its own defaults.
	MetricCollector *MetricCollectorWorkerConfig
}

func (c *SystemConfig) WithDefaults() *SystemConfig {
	return (*SystemConfig)((*system.SystemConfig)(c).WithDefaults())
}

func (c *SystemConfig) Validate() error {
	return (*system.SystemConfig)(c).Validate()
}
