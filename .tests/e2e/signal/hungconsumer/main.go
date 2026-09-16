package main

// hungconsumer prints "consuming" and consumes under the lifecycle context;
// its handler prints "handler blocked" and ignores its context until SIGUSR1.
// This lets the parent release work during graceful shutdown or leave it
// blocked to test forced exit. The message timeout is five minutes.

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/allegedlyreliable/sqlstreams/.tests/e2e/common"
	"github.com/allegedlyreliable/sqlstreams/client"
)

var streamName = flag.String("stream", "", "the stream to consume")

func main() {
	flag.Parse()
	if err := run(); err != nil {
		fmt.Printf("\n❌ E2E TEST FAILED: %s\n", err.Error())
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := sqlstreams.LifecycleContext(nil)
	defer stop()
	pool, err := common.NewPool(ctx, nil)
	if err != nil {
		return err
	}
	defer pool.Close()
	var warnings atomic.Int64
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelWarn,
		ReplaceAttr: func(groups []string, attribute slog.Attr) slog.Attr {
			if attribute.Key == slog.LevelKey && attribute.Value.Any().(slog.Level) >= slog.LevelWarn {
				warnings.Add(1)
			}
			return attribute
		},
	}))
	client, err := sqlstreams.NewClient(ctx, pool, &sqlstreams.ClientConfig{Logger: logger})
	if err != nil {
		return err
	}
	cfg := &sqlstreams.ConsumerConfig{Start: sqlstreams.Head(), Message: &sqlstreams.MessageOptions{Timeout: 5 * time.Minute}}
	consumer, err := client.Stream[common.Work](*streamName).Consumer("signal.e2e.hung").Register(ctx, cfg)
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		fmt.Println("draining")
	}()
	fmt.Println("consuming")
	err = consumer.Consume(ctx, handle, &sqlstreams.ConsumeOptions{ClaimPollRate: 200 * time.Millisecond})
	if warnings.Load() != 0 {
		return fmt.Errorf("graceful drain emitted %d warning/error records, want 0", warnings.Load())
	}
	if errors.Is(err, context.Canceled) {
		return nil
	}
	return err
}

func handle(ctx context.Context, work *common.Work) error {
	release := make(chan os.Signal, 1)
	signal.Notify(release, syscall.SIGUSR1)
	defer signal.Stop(release)
	fmt.Println("handler blocked")
	<-release
	return nil
}
