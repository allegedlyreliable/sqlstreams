package main

// Scenario 03 -- consume with retry and dead-lettering.
//
// This scenario produces three uploads and handles their outcomes: a corrupt
// upload will never succeed (terminal), unavailable storage may recover
// (retry), and an embargoed video waits without counting as a failure
// (delay). Its own consumer group leaves scenario 02's cursor and config alone.

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	sqlstreams "github.com/allegedlyreliable/sqlstreams/client"
)

var (
	errSourceCorrupt      = errors.New("source video is corrupt")
	errStorageUnavailable = errors.New("source storage is unavailable")
)

type VideoUploadedV1 struct {
	VideoId         string `json:"video_id"`
	OwnerId         string `json:"owner_id"`
	UploadId        string `json:"upload_id"`
	DurationMinutes int    `json:"duration_minutes"`
	SourceStatus    string `json:"source_status"` // "ready" | "corrupt" | "unavailable"
	ReleaseAtUnix   int64  `json:"release_at_unix"`
}

// increment on breaking changes
func (VideoUploadedV1) SchemaVersion() int { return 1 }

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

	uploads := client.Stream[VideoUploadedV1]("videos.uploaded")
	if _, err := uploads.Register(ctx, nil); err != nil {
		return err
	}
	producer, err := uploads.Producer().Register(ctx, nil)
	if err != nil {
		return err
	}

	transcoder := uploads.Consumer("transcoder-with-retry")
	consumer, err := transcoder.Register(ctx, &sqlstreams.ConsumerConfig{
		Message: &sqlstreams.MessageOptions{
			Timeout: 10 * time.Second,
			Retry:   &sqlstreams.RetryPolicy{MaxRetries: 3, BaseDelay: 2 * time.Second},
		},
	})
	if err != nil {
		return err
	}

	for _, video := range []VideoUploadedV1{
		{VideoId: "video-corrupt", OwnerId: "creator-7", UploadId: "upl-corrupt", SourceStatus: "corrupt"},
		{VideoId: "video-unavailable", OwnerId: "creator-7", UploadId: "upl-unavailable", SourceStatus: "unavailable"},
		{VideoId: "video-embargoed", OwnerId: "creator-7", UploadId: "upl-embargoed", SourceStatus: "ready", ReleaseAtUnix: time.Now().Add(5 * time.Second).Unix()},
	} {
		if _, err := producer.Produce(ctx, &video, nil); err != nil {
			return err
		}
	}

	return consumer.Consume(ctx, transcodeVideo, nil)
}

func transcodeVideo(ctx context.Context, video *VideoUploadedV1) error {
	meta, _ := sqlstreams.MetaFromContext(ctx)
	fmt.Printf("transcoding %s (message %d, attempt %d, delays %d)\n",
		video.VideoId, meta.Id, meta.Attempts+1, meta.Delays)

	switch video.SourceStatus {
	case "corrupt":
		// dead on this attempt; the cause lands in last_error
		return sqlstreams.Terminal(errSourceCorrupt)
	case "unavailable":
		// the first attempt retries with backoff; the next simulates recovery
		if meta.Attempts == 0 {
			return errStorageUnavailable
		}
	}
	releaseAt := time.Unix(video.ReleaseAtUnix, 0)
	if releaseAt.After(time.Now()) {
		// runs again after the release window
		return sqlstreams.Delay(time.Until(releaseAt))
	}
	return nil
}
