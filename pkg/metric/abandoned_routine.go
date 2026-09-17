package metric

import "time"

// AbandonedRoutineSnapshot pairs retained abandoned and cleared events on
// __system.metrics for one stream and consumer group. Counts describe the
// retained window, not lifetime totals.
type AbandonedRoutineSnapshot struct {
	Outstanding         int64         `json:"outstanding"`            // abandoned events with no matching cleared
	Total               int64         `json:"total"`                  // distinct abandoned keys currently in the window
	SelfClearLatencyAvg time.Duration `json:"self_clear_latency_avg"` // mean(cleared.At - abandoned.At) over matched pairs; 0 if no pair has cleared yet
}
