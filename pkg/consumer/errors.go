package consumer

// appended to errors.ErrLifecycleContextNotCancellable at the Consume call site
const lifecycleContextHelp = `
Consume's context is the instance's lifetime -- cancelling it starts graceful
shutdown, and context.Background/TODO can never be cancelled.

Pass your application's shutdown context:

    import sqlstreams "github.com/allegedlyreliable/sqlstreams/client"

    ctx, stop := sqlstreams.LifecycleContext(nil)
    defer stop()

Or run a session that only stops with the process:

    instance.Consume(ctx, handler, &sqlstreams.ConsumeOptions{DisableGracefulShutdown: true})`
