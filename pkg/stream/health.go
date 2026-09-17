package stream

import (
	"github.com/allegedlyreliable/sqlstreams/pkg/metric"
)

// StreamVersionHealth reports whether a payload version can be retired.
// Safe requires no compaction heads, unread messages, or unresolved exceptions.
type StreamVersionHealth struct {
	Stream          *Stream                                `json:"stream"`
	Version         int                                    `json:"version"`
	Messages        int64                                  `json:"messages"`         // rows in the log at this version
	CompactionHeads int64                                  `json:"compaction_heads"` // keys whose current head is at this version
	Groups          []metric.ConsumerGroupSchemaVersionLag `json:"groups"`           // each group's unread and unresolved rows at it
	Safe            bool                                   `json:"safe"`             // CompactionHeads is 0 and every group's counts are 0
	Reason          string                                 `json:"reason"`           // the verdict in words: which heads or groups still hold it
}
