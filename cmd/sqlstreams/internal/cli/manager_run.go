package cli

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/allegedlyreliable/sqlstreams/otel"
	"github.com/allegedlyreliable/sqlstreams/pkg/common/logging"
	"github.com/allegedlyreliable/sqlstreams/pkg/migrate"
	"github.com/spf13/cobra"
)

func newManagerRunCmd(g *globalFlags) *cobra.Command {
	var metricsAddress string
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run the system manager until stopped",
		Long: `Run deployment maintenance and built-in alert processing, including retention,
cursor advancement, scheduled message production, and metrics collection.
Application consumers run separately. Suspended workers remain suspended.

Multiple replicas can run: one holds the system manager claim, and the others
retry until they can take over after it is released or expires. Worker claims
enforce each worker's instance target.

With --metrics-address, serve measurements at a Prometheus /metrics endpoint.
Logs go to stderr. --output json is not supported. Stop with SIGINT or SIGTERM.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			// the daemon's output is its log stream; there is no result document
			if g.jsonOutput() {
				return failUsage("manager run streams logs and produces no result document -- --output json does not apply")
			}

			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
			defer stop()
			// once a stop begins, re-arm default delivery so a second signal
			// force-kills instead of being swallowed mid-drain
			go func() { <-ctx.Done(); stop() }()

			// unlike the one-shot commands, the daemon's log stream IS its
			// output -- full info level, still on stderr by convention
			connection, err := newConnection(ctx, g.databaseURL, g.schema, slog.LevelInfo)
			if err != nil {
				return err
			}
			defer connection.Close()
			client := connection.client
			runLogger := logging.NewPipelineLogger(connection.config.Logger, &logging.PipelineLoggerConfig{
				Args: []any{"schema", connection.config.Schema},
			})

			// a server failure cancels runCtx so the manager drains too
			runCtx, cancelRun := context.WithCancel(ctx)
			defer cancelRun()
			serverFailed := make(chan error, 1)
			if metricsAddress != "" {
				exporter, err := otel.NewExporter(ctx, connection.pool, &otel.ExporterConfig{
					Schema: connection.config.Schema,
					Logger: connection.config.Logger,
					Retry:  connection.config.Retry,
				})
				if err != nil {
					return failOp("%s", err.Error())
				}
				defer func() {
					shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancelShutdown()
					_ = exporter.Close(shutdownCtx)
				}()
				if _, err := client.System().Metrics().Latest(ctx); err != nil {
					if errors.Is(err, migrate.ErrNotRegistered) {
						return errSystemNotRegistered()
					}
					return failOp("%s", err.Error())
				}

				mux := http.NewServeMux()
				mux.Handle("/metrics", exporter.Handler())
				server := &http.Server{Addr: metricsAddress, Handler: mux}
				go func() {
					if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
						serverFailed <- err
						cancelRun()
					}
				}()
				defer func() {
					shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancelShutdown()
					_ = server.Shutdown(shutdownCtx)
				}()
				runLogger.InfoContext(ctx, "metrics endpoint serving", "address", metricsAddress)
			}

			// a signal cancels ctx; Run drains its worker claims and returns nil
			if err := client.Manager().Run(runCtx); err != nil {
				if errors.Is(err, migrate.ErrNotRegistered) {
					return errSystemNotRegistered()
				}
				return failOp("%s", err.Error())
			}
			select {
			case err := <-serverFailed:
				return failOp("metrics endpoint: %s", err.Error())
			default:
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&metricsAddress, "metrics-address", "",
		"serve a Prometheus /metrics endpoint on this address (e.g. :9464)")
	return cmd
}
