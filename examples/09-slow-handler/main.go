package main

// Scenario 09 -- a longer timeout for a slow message.
//
// Simulate each minute of video with 50 milliseconds of work. A short video
// fits the consumer's one-second default; a 95-minute video takes 4.75 seconds
// and asks for ten seconds. Both finish within their timeout and lease.
//
// How timeout resolves, highest first:
//	consumer clamp (MessageMin/MessageMax) > produced message > producer defaults > consumer defaults > system defaults

import (
	"context"
	"fmt"
	"os"
	"time"

	sqlstreams "github.com/allegedlyreliable/sqlstreams/client"
)

type VideoTranscodeRequestedV1 struct {
	VideoId         string `json:"video_id"`
	DurationMinutes int    `json:"duration_minutes"`
}

// increment on breaking changes
func (VideoTranscodeRequestedV1) SchemaVersion() int { return 1 }

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

	requests := client.Stream[VideoTranscodeRequestedV1]("videos.transcode-requested")
	if _, err := requests.Register(ctx, nil); err != nil {
		return err
	}
	producer, err := requests.Producer().Register(ctx, nil)
	if err != nil {
		return err
	}

	// ConsumerConfig.Message    -> the timeout a message gets when it asks for nothing (most videos)
	// ConsumerConfig.MessageMax -> the most a message may ask for; a request above it is lowered to it
	transcoder := requests.Consumer("feature-transcoder")
	consumer, err := transcoder.Register(ctx, &sqlstreams.ConsumerConfig{
		Message:    &sqlstreams.MessageOptions{Timeout: time.Second},
		MessageMax: &sqlstreams.MessageOptions{Timeout: 10 * time.Second},
	})
	if err != nil {
		return err
	}

	if _, err := producer.Produce(ctx, &VideoTranscodeRequestedV1{VideoId: "video-42", DurationMinutes: 12}, nil); err != nil {
		return err
	}

	// The feature-length video asks for longer than the consumer's default.
	if _, err := producer.Produce(ctx, &VideoTranscodeRequestedV1{VideoId: "video-99", DurationMinutes: 95},
		&sqlstreams.ProduceOptions{Message: &sqlstreams.MessageOptions{Timeout: 10 * time.Second}}); err != nil {
		return err
	}

	return consumer.Consume(ctx, transcodeVideo, nil)
}

func transcodeVideo(ctx context.Context, video *VideoTranscodeRequestedV1) error {
	meta, _ := sqlstreams.MetaFromContext(ctx)
	work := time.Duration(video.DurationMinutes) * 50 * time.Millisecond
	fmt.Printf("transcoding %s: simulated work %s, timeout %s\n", video.VideoId, work, meta.Options.Timeout)

	// Observe the deadline so a timed-out handler stops its work before retry.
	timer := time.NewTimer(work)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		fmt.Printf("transcoded %s\n", video.VideoId)
	}
	return nil
}
