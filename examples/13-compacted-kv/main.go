package main

// Scenario 13 -- a compacted stream used as a key/value store.
//
// One current processing document per video id. Read the current value,
// write a new one, and increment its attempt count safely under concurrent
// writers (read-modify-write).

import (
	"context"
	"fmt"
	"os"

	sqlstreams "github.com/allegedlyreliable/sqlstreams/client"
)

type VideoProcessingStateV1 struct {
	VideoId  string `json:"video_id"`
	Stage    string `json:"stage"`
	Attempts int    `json:"attempts"`
}

// increment on breaking changes
func (VideoProcessingStateV1) SchemaVersion() int { return 1 }

func main() {
	if err := run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := sqlstreams.LifecycleContext(nil)
	defer stop()

	pool, err := sqlstreams.NewPostgresPool(ctx, "example_user", "example_password", "localhost", "example_db", nil)
	if err != nil {
		return err
	}
	defer pool.Close()

	client, err := sqlstreams.NewClient(ctx, pool, nil)
	if err != nil {
		return err
	}

	states := client.Stream[VideoProcessingStateV1]("videos.processing-state")
	_, err = states.Register(ctx, nil)
	if err != nil {
		return err
	}

	producer, err := states.Producer().Register(ctx, nil)
	if err != nil {
		return err
	}

	video := states.Key("video-42")

	// Initialize under the same lock as updates, so reruns cannot reset the count.
	if err := client.InTransaction(ctx, func(ctx context.Context, tx sqlstreams.Tx) error {
		head, err := video.LockCompactionHead(ctx, tx)
		if err != nil {
			return err
		}
		if head != nil {
			return nil
		}
		_, err = producer.ProduceInTx(ctx, tx, &VideoProcessingStateV1{VideoId: "video-42", Stage: "transcoding", Attempts: 1},
			&sqlstreams.ProduceOptions{MessageKey: "video-42", Compaction: &sqlstreams.CompactionOptions{Enable: true}})
		return err
	}); err != nil {
		return err
	}

	// Get (outside a transaction) -- the stream handle's read
	current, err := video.CompactionHead(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("current: id=%d stage=%s attempts=%d\n", current.Id, current.Message.Stage, current.Message.Attempts)

	// Update: hold the head's row lock throughout the read-modify-write.
	if err := client.InTransaction(ctx, func(ctx context.Context, tx sqlstreams.Tx) error {
		head, err := video.LockCompactionHead(ctx, tx)
		if err != nil {
			return err
		}
		next := VideoProcessingStateV1{VideoId: "video-42", Stage: "transcoding"}
		if head != nil {
			next = *head.Message
		}
		next.Attempts++
		_, err = producer.ProduceInTx(ctx, tx, &next, &sqlstreams.ProduceOptions{MessageKey: "video-42", Compaction: &sqlstreams.CompactionOptions{Enable: true}})
		return err
	}); err != nil {
		return err
	}

	// History
	versions, err := video.Messages(ctx, 10)
	if err != nil {
		return err
	}
	for _, version := range versions {
		fmt.Printf("version id=%d attempts=%d\n", version.Id, version.Message.Attempts)
	}
	return nil
}
