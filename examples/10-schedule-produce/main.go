package main

// Scenario 10 -- a schedule that produces on a cron expression.
//
// A usage report every minute: register the schedule with its message and
// target stream, then consume that stream like any other. For a nightly
// report at 02:00 UTC, use "0 2 * * *" instead of "* * * * *".

import (
	"context"
	"fmt"
	"os"
	"time"

	sqlstreams "github.com/allegedlyreliable/sqlstreams/client"
	"golang.org/x/sync/errgroup"
)

type UsageReportRequestedV1 struct {
	Scope string `json:"scope"`
}

// increment on breaking changes
func (UsageReportRequestedV1) SchemaVersion() int { return 1 }

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
	reports := client.Stream[UsageReportRequestedV1]("usage.reports.requested")
	if _, err := reports.Register(ctx, &sqlstreams.StreamConfig{DeliveryLogMode: sqlstreams.DeliveryLogModeAll}); err != nil {
		return err
	}

	reportSchedule := client.Scheduler("usage.reports.every-minute")
	scheduler, err := reportSchedule.Register(ctx, "usage.reports.requested", "* * * * *", &UsageReportRequestedV1{Scope: "all-creators"}, nil)
	if err != nil {
		return err
	}

	builder := reports.Consumer("usage-report-builder")
	consumer, err := builder.Register(ctx, nil)
	if err != nil {
		return err
	}

	// Schedule is the process that produces the message; a program that
	// registers the schedule and exits leaves one that never runs
	fmt.Println("waiting for the minute schedule (first delivery can take up to two minutes)")
	routines, routinesCtx := errgroup.WithContext(ctx)
	routines.Go(func() error { return scheduler.Schedule(routinesCtx) })
	routines.Go(func() error { return consumer.Consume(routinesCtx, buildUsageReport, nil) })
	return routines.Wait()
}

func buildUsageReport(ctx context.Context, request *UsageReportRequestedV1) error {
	// the payload is the same every run; the scheduled time is on the delivery's meta
	meta, _ := sqlstreams.MetaFromContext(ctx)
	fmt.Printf("building %s usage report for %s\n", request.Scope, meta.ScheduledAt.Format(time.RFC3339))
	return nil
}
