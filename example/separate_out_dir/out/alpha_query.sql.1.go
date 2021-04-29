// Code generated by pggen. DO NOT EDIT.

package out

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
)

const alphaSQL = `SELECT 'alpha' as output;`

// Alpha implements Querier.Alpha.
func (q *DBQuerier) Alpha(ctx context.Context) (string, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "Alpha")
	row := q.conn.QueryRow(ctx, alphaSQL)
	var item string
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query Alpha: %w", err)
	}
	return item, nil
}

// AlphaBatch implements Querier.AlphaBatch.
func (q *DBQuerier) AlphaBatch(batch *pgx.Batch) {
	batch.Queue(alphaSQL)
}

// AlphaScan implements Querier.AlphaScan.
func (q *DBQuerier) AlphaScan(results pgx.BatchResults) (string, error) {
	row := results.QueryRow()
	var item string
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan AlphaBatch row: %w", err)
	}
	return item, nil
}
