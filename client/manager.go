package sqlstreams

import "context"

// ManagerHandle is a handle on the client's system manager.
type ManagerHandle struct {
	client *Client
}

// Manager returns the client system manager's handle. No I/O and no failure.
func (c *Client) Manager() *ManagerHandle {
	return &ManagerHandle{client: c}
}

// Run maintains registered streams, produces scheduled messages, collects
// metrics, and runs built-in alert consumers until ctx is cancelled.
// It does not run application consumer handlers. Cancellation returns nil.
// Multiple processes may call Run. A database lease admits one reconciliation
// loop at a time. An unregistered system returns ErrNotRegistered.
func (m *ManagerHandle) Run(ctx context.Context) error {
	return m.client.manager.Run(ctx)
}
