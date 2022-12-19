package pgtest

// import (
// 	"context"
// 	"fmt"

// 	"github.com/jackc/pgx/v5/internal/stmtcache"
// 	"github.com/jackc/pgx/v5/pgconn"
// )

// // GuardedStmtCache errors if any name in names is used to get a cached statement.
// // Allows verifying that PrepareAllQueries works by creating prepared statements
// // ahead of time. pgx accesses a map of prepared statements directly rather than
// // calling Get.
// type GuardedStmtCache struct {
// 	*stmtcache.LRU
// 	names map[string]struct{}
// }

// func NewGuardedStmtCache(conn *pgconn.PgConn, names ...string) *GuardedStmtCache {
// 	size := 1024 // so we never expire
// 	nameMap := make(map[string]struct{})
// 	for _, n := range names {
// 		nameMap[n] = struct{}{}
// 	}
// 	return &GuardedStmtCache{
// 		LRU:   stmtcache.NewLRU(conn, stmtcache.ModePrepare, size),
// 		names: nameMap,
// 	}
// }

// func (sc *GuardedStmtCache) Get(ctx context.Context, sql string) (*pgconn.StatementDescription, error) {
// 	if _, ok := sc.names[sql]; ok {
// 		return nil, fmt.Errorf("guard statement cache attempted to get %s;"+
// 			" should already exist as prepared statement", sql)
// 	}
// 	return sc.LRU.Get(ctx, sql)
// }
