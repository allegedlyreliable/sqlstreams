package datastore

import (
	"context"
	"fmt"

	"github.com/allegedlyreliable/sqlstreams/pkg/common"
	"github.com/allegedlyreliable/sqlstreams/pkg/datastore"
	"github.com/allegedlyreliable/sqlstreams/pkg/migrate"
)

// SystemOwner resolves the singleton system row to its owner, read on the
// pool. Returns ErrNotRegistered if the row, or the table itself, isn't there.
func (d *MigrateDatastore) SystemOwner(ctx context.Context) (*common.Owner, error) {
	var owner *common.Owner
	err := d.DatastoreRetry.WrapIdempotent(ctx, func() error {
		var err error
		owner, err = d.systemOwner(ctx)
		return err
	})
	return owner, err
}

func (d *MigrateDatastore) systemOwner(ctx context.Context) (*common.Owner, error) {
	return SystemOwner(ctx, d.Datastore.Pool, d.Datastore.Schema)
}

// SystemOwner resolves the singleton system row to its owner. Returns
// ErrNotRegistered if the row, or the table itself, isn't there.
func SystemOwner(ctx context.Context, q datastore.Querier, schema string) (*common.Owner, error) {
	installed, err := tableExists(ctx, q, schema, "system_config")
	if err != nil {
		return nil, err
	}
	if !installed {
		return nil, migrate.ErrNotRegistered
	}

	var id int64
	sql := fmt.Sprintf(`
		-- sqlstreams: migrate.SystemOwner
		SELECT id FROM %[1]s.system_config;
	`, schema)
	if err := q.QueryRow(ctx, sql).Scan(&id); err != nil {
		return nil, registrationError(err)
	}
	return common.NewSystemOwner(id)
}

// ***************
// *** HELPERS ***
// ***************

// tableExists asks the catalog before a read touches the table: a statement
// against a missing table is an ERROR in the server log, a NULL regclass is not.
func tableExists(ctx context.Context, q datastore.Querier, schema string, name string) (bool, error) {
	sql := `
		-- sqlstreams: migrate.tableExists
		SELECT to_regclass($1) IS NOT NULL;
	`
	var exists bool
	err := q.QueryRow(ctx, sql, schema+"."+name).Scan(&exists)
	return exists, err
}
