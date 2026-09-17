package controller

import (
	"time"

	"github.com/allegedlyreliable/sqlstreams/pkg/alert"
	"github.com/allegedlyreliable/sqlstreams/pkg/common"
	"github.com/allegedlyreliable/sqlstreams/pkg/metric"
)

func newWorkerLivenessAlert(owner *common.Owner, unclaimed []*metric.UnclaimedWorkerMetadata, at time.Time) (*alert.Alert, error) {
	workersByOwner := groupUnclaimedWorkersByOwner(unclaimed)
	summary, err := newWorkerLivenessSummary(workersByOwner)
	if err != nil {
		return nil, err
	}

	message := summary.message()
	options := &alert.AlertOptions{
		Hint: summary.hint(),
		Data: unclaimedWorkerData(unclaimed),
	}
	return alert.NewAlert(alert.AlertWorkerLiveness.Name, owner, alert.AlertStatusActive, alert.AlertSeverity(alert.AlertWorkerLiveness.Severity), message, at, options)
}

func groupUnclaimedWorkersByOwner(unclaimed []*metric.UnclaimedWorkerMetadata) [][]*metric.UnclaimedWorkerMetadata {
	groups := make([][]*metric.UnclaimedWorkerMetadata, 0)
	positions := make(map[common.Owner]int)
	for _, worker := range unclaimed {
		position, seen := positions[*worker.Owner]
		if !seen {
			position = len(groups)
			positions[*worker.Owner] = position
			groups = append(groups, nil)
		}
		groups[position] = append(groups[position], worker)
	}
	return groups
}

func unclaimedWorkerData(unclaimed []*metric.UnclaimedWorkerMetadata) map[string]any {
	rows := make([]map[string]any, 0, len(unclaimed))
	for _, worker := range unclaimed {
		rows = append(rows, map[string]any{
			"worker":           worker.Name,
			"owner":            worker.Owner.Name,
			"owner_kind":       string(worker.Owner.Kind()),
			"target_instances": worker.TargetInstances,
		})
	}
	return map[string]any{
		"unclaimed_count": len(unclaimed),
		"workers":         rows,
	}
}
