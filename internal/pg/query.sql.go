// Code generated by pggen. DO NOT EDIT.

package pg

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

// Querier is a typesafe Go interface backed by SQL queries.
//
// Methods ending with Batch enqueue a query to run later in a pgx.Batch. After
// calling SendBatch on pgx.Conn, pgxpool.Pool, or pgx.Tx, use the Scan methods
// to parse the results.
type Querier interface {
	FindEnumTypes(ctx context.Context, oids []uint32) ([]FindEnumTypesRow, error)
	// FindEnumTypesBatch enqueues a FindEnumTypes query into batch to be executed
	// later by the batch.
	FindEnumTypesBatch(batch *pgx.Batch, oids []uint32)
	// FindEnumTypesScan scans the result of an executed FindEnumTypesBatch query.
	FindEnumTypesScan(results pgx.BatchResults) ([]FindEnumTypesRow, error)

	FindArrayTypes(ctx context.Context, oids []uint32) ([]FindArrayTypesRow, error)
	// FindArrayTypesBatch enqueues a FindArrayTypes query into batch to be executed
	// later by the batch.
	FindArrayTypesBatch(batch *pgx.Batch, oids []uint32)
	// FindArrayTypesScan scans the result of an executed FindArrayTypesBatch query.
	FindArrayTypesScan(results pgx.BatchResults) ([]FindArrayTypesRow, error)

	// A composite type represents a row or record, defined implicitly for each
	// table, or explicitly with CREATE TYPE.
	// https://www.postgresql.org/docs/13/rowtypes.html
	FindCompositeTypes(ctx context.Context, oids []uint32) ([]FindCompositeTypesRow, error)
	// FindCompositeTypesBatch enqueues a FindCompositeTypes query into batch to be executed
	// later by the batch.
	FindCompositeTypesBatch(batch *pgx.Batch, oids []uint32)
	// FindCompositeTypesScan scans the result of an executed FindCompositeTypesBatch query.
	FindCompositeTypesScan(results pgx.BatchResults) ([]FindCompositeTypesRow, error)

	// Recursively expands all given OIDs to all descendants through composite
	// types.
	FindDescendantOIDs(ctx context.Context, oids []uint32) ([]pgtype.OID, error)
	// FindDescendantOIDsBatch enqueues a FindDescendantOIDs query into batch to be executed
	// later by the batch.
	FindDescendantOIDsBatch(batch *pgx.Batch, oids []uint32)
	// FindDescendantOIDsScan scans the result of an executed FindDescendantOIDsBatch query.
	FindDescendantOIDsScan(results pgx.BatchResults) ([]pgtype.OID, error)

	FindOIDByName(ctx context.Context, name string) (pgtype.OID, error)
	// FindOIDByNameBatch enqueues a FindOIDByName query into batch to be executed
	// later by the batch.
	FindOIDByNameBatch(batch *pgx.Batch, name string)
	// FindOIDByNameScan scans the result of an executed FindOIDByNameBatch query.
	FindOIDByNameScan(results pgx.BatchResults) (pgtype.OID, error)

	FindOIDName(ctx context.Context, oid pgtype.OID) (pgtype.Name, error)
	// FindOIDNameBatch enqueues a FindOIDName query into batch to be executed
	// later by the batch.
	FindOIDNameBatch(batch *pgx.Batch, oid pgtype.OID)
	// FindOIDNameScan scans the result of an executed FindOIDNameBatch query.
	FindOIDNameScan(results pgx.BatchResults) (pgtype.Name, error)

	FindOIDNames(ctx context.Context, oid []uint32) ([]FindOIDNamesRow, error)
	// FindOIDNamesBatch enqueues a FindOIDNames query into batch to be executed
	// later by the batch.
	FindOIDNamesBatch(batch *pgx.Batch, oid []uint32)
	// FindOIDNamesScan scans the result of an executed FindOIDNamesBatch query.
	FindOIDNamesScan(results pgx.BatchResults) ([]FindOIDNamesRow, error)
}

type DBQuerier struct {
	conn  genericConn   // underlying Postgres transport to use
	types *typeResolver // resolve types by name
}

var _ Querier = &DBQuerier{}

// genericConn is a connection to a Postgres database. This is usually backed by
// *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
type genericConn interface {
	// Query executes sql with args. If there is an error the returned Rows will
	// be returned in an error state. So it is allowed to ignore the error
	// returned from Query and handle it in Rows.
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)

	// QueryRow is a convenience wrapper over Query. Any error that occurs while
	// querying is deferred until calling Scan on the returned Row. That Row will
	// error with pgx.ErrNoRows if no rows are returned.
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row

	// Exec executes sql. sql can be either a prepared statement name or an SQL
	// string. arguments should be referenced positionally from the sql string
	// as $1, $2, etc.
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

// NewQuerier creates a DBQuerier that implements Querier. conn is typically
// *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
func NewQuerier(conn genericConn) *DBQuerier {
	return NewQuerierConfig(conn, QuerierConfig{})
}

type QuerierConfig struct {
	// DataTypes contains pgtype.Value to use for encoding and decoding instead
	// of pggen-generated pgtype.ValueTranscoder.
	//
	// If OIDs are available for an input parameter type and all of its
	// transitive dependencies, pggen will use the binary encoding format for
	// the input parameter.
	DataTypes []pgtype.DataType
}

// NewQuerierConfig creates a DBQuerier that implements Querier with the given
// config. conn is typically *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
func NewQuerierConfig(conn genericConn, cfg QuerierConfig) *DBQuerier {
	return &DBQuerier{conn: conn, types: newTypeResolver(cfg.DataTypes)}
}

// WithTx creates a new DBQuerier that uses the transaction to run all queries.
func (q *DBQuerier) WithTx(tx pgx.Tx) (*DBQuerier, error) {
	return &DBQuerier{conn: tx}, nil
}

// preparer is any Postgres connection transport that provides a way to prepare
// a statement, most commonly *pgx.Conn.
type preparer interface {
	Prepare(ctx context.Context, name, sql string) (sd *pgconn.StatementDescription, err error)
}

// PrepareAllQueries executes a PREPARE statement for all pggen generated SQL
// queries in querier files. Typical usage is as the AfterConnect callback
// for pgxpool.Config
//
// pgx will use the prepared statement if available. Calling PrepareAllQueries
// is an optional optimization to avoid a network round-trip the first time pgx
// runs a query if pgx statement caching is enabled.
func PrepareAllQueries(ctx context.Context, p preparer) error {
	if _, err := p.Prepare(ctx, findEnumTypesSQL, findEnumTypesSQL); err != nil {
		return fmt.Errorf("prepare query 'FindEnumTypes': %w", err)
	}
	if _, err := p.Prepare(ctx, findArrayTypesSQL, findArrayTypesSQL); err != nil {
		return fmt.Errorf("prepare query 'FindArrayTypes': %w", err)
	}
	if _, err := p.Prepare(ctx, findCompositeTypesSQL, findCompositeTypesSQL); err != nil {
		return fmt.Errorf("prepare query 'FindCompositeTypes': %w", err)
	}
	if _, err := p.Prepare(ctx, findDescendantOIDsSQL, findDescendantOIDsSQL); err != nil {
		return fmt.Errorf("prepare query 'FindDescendantOIDs': %w", err)
	}
	if _, err := p.Prepare(ctx, findOIDByNameSQL, findOIDByNameSQL); err != nil {
		return fmt.Errorf("prepare query 'FindOIDByName': %w", err)
	}
	if _, err := p.Prepare(ctx, findOIDNameSQL, findOIDNameSQL); err != nil {
		return fmt.Errorf("prepare query 'FindOIDName': %w", err)
	}
	if _, err := p.Prepare(ctx, findOIDNamesSQL, findOIDNamesSQL); err != nil {
		return fmt.Errorf("prepare query 'FindOIDNames': %w", err)
	}
	return nil
}

// typeResolver looks up the pgtype.ValueTranscoder by Postgres type name.
type typeResolver struct {
	connInfo *pgtype.ConnInfo // types by Postgres type name
}

func newTypeResolver(types []pgtype.DataType) *typeResolver {
	ci := pgtype.NewConnInfo()
	for _, typ := range types {
		if txt, ok := typ.Value.(textPreferrer); ok && typ.OID != unknownOID {
			typ.Value = txt.ValueTranscoder
		}
		ci.RegisterDataType(typ)
	}
	return &typeResolver{connInfo: ci}
}

// findValue find the OID, and pgtype.ValueTranscoder for a Postgres type name.
func (tr *typeResolver) findValue(name string) (uint32, pgtype.ValueTranscoder, bool) {
	typ, ok := tr.connInfo.DataTypeForName(name)
	if !ok {
		return 0, nil, false
	}
	v := pgtype.NewValue(typ.Value)
	return typ.OID, v.(pgtype.ValueTranscoder), true
}

// setValue sets the value of a ValueTranscoder to a value that should always
// work and panics if it fails.
func (tr *typeResolver) setValue(vt pgtype.ValueTranscoder, val interface{}) pgtype.ValueTranscoder {
	if err := vt.Set(val); err != nil {
		panic(fmt.Sprintf("set ValueTranscoder %T to %+v: %s", vt, val, err))
	}
	return vt
}

const findEnumTypesSQL = `WITH enums AS (
  SELECT
    enumtypid::int8                                   AS enum_type,
    -- pg_enum row identifier.
    -- The OIDs for pg_enum rows follow a special rule: even-numbered OIDs
    -- are guaranteed to be ordered in the same way as the sort ordering of
    -- their enum type. That is, if two even OIDs belong to the same enum
    -- type, the smaller OID must have the smaller enumsortorder value.
    -- Odd-numbered OID values need bear no relationship to the sort order.
    -- This rule allows the enum comparison routines to avoid catalog
    -- lookups in many common cases. The routines that create and alter enum
    -- types attempt to assign even OIDs to enum values whenever possible.
    array_agg(oid::int8 ORDER BY enumsortorder)       AS enum_oids,
    -- The sort position of this enum value within its enum type. Starts as
    -- 1..n but can be fractional or negative.
    array_agg(enumsortorder ORDER BY enumsortorder)   AS enum_orders,
    -- The textual label for this enum value
    array_agg(enumlabel::text ORDER BY enumsortorder) AS enum_labels
  FROM pg_enum
  GROUP BY pg_enum.enumtypid)
SELECT
  typ.oid           AS oid,
  -- typename: Data type name.
  typ.typname::text AS type_name,
  enum.enum_oids    AS child_oids,
  enum.enum_orders  AS orders,
  enum.enum_labels  AS labels,
  -- typtype: b for a base type, c for a composite type (e.g., a table's
  -- row type), d for a domain, e for an enum type, p for a pseudo-type,
  -- or r for a range type.
  typ.typtype       AS type_kind,
  -- typdefault is null if the type has no associated default value. If
  -- typdefaultbin is not null, typdefault must contain a human-readable
  -- version of the default expression represented by typdefaultbin. If
  -- typdefaultbin is null and typdefault is not, then typdefault is the
  -- external representation of the type's default value, which can be fed
  -- to the type's input converter to produce a constant.
  COALESCE(typ.typdefault, '')    AS default_expr
FROM pg_type typ
  JOIN enums enum ON typ.oid = enum.enum_type
WHERE typ.typisdefined
  AND typ.typtype = 'e'
  AND typ.oid = ANY ($1::oid[]);`

type FindEnumTypesRow struct {
	OID         pgtype.OID   `json:"oid"`
	TypeName    string       `json:"type_name"`
	ChildOIDs   []int        `json:"child_oids"`
	Orders      []float32    `json:"orders"`
	Labels      []string     `json:"labels"`
	TypeKind    pgtype.QChar `json:"type_kind"`
	DefaultExpr string       `json:"default_expr"`
}

// FindEnumTypes implements Querier.FindEnumTypes.
func (q *DBQuerier) FindEnumTypes(ctx context.Context, oids []uint32) ([]FindEnumTypesRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindEnumTypes")
	rows, err := q.conn.Query(ctx, findEnumTypesSQL, oids)
	if err != nil {
		return nil, fmt.Errorf("query FindEnumTypes: %w", err)
	}
	defer rows.Close()
	items := []FindEnumTypesRow{}
	for rows.Next() {
		var item FindEnumTypesRow
		if err := rows.Scan(&item.OID, &item.TypeName, &item.ChildOIDs, &item.Orders, &item.Labels, &item.TypeKind, &item.DefaultExpr); err != nil {
			return nil, fmt.Errorf("scan FindEnumTypes row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindEnumTypes rows: %w", err)
	}
	return items, err
}

// FindEnumTypesBatch implements Querier.FindEnumTypesBatch.
func (q *DBQuerier) FindEnumTypesBatch(batch *pgx.Batch, oids []uint32) {
	batch.Queue(findEnumTypesSQL, oids)
}

// FindEnumTypesScan implements Querier.FindEnumTypesScan.
func (q *DBQuerier) FindEnumTypesScan(results pgx.BatchResults) ([]FindEnumTypesRow, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query FindEnumTypesBatch: %w", err)
	}
	defer rows.Close()
	items := []FindEnumTypesRow{}
	for rows.Next() {
		var item FindEnumTypesRow
		if err := rows.Scan(&item.OID, &item.TypeName, &item.ChildOIDs, &item.Orders, &item.Labels, &item.TypeKind, &item.DefaultExpr); err != nil {
			return nil, fmt.Errorf("scan FindEnumTypesBatch row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindEnumTypesBatch rows: %w", err)
	}
	return items, err
}

const findArrayTypesSQL = `SELECT
  arr_typ.oid           AS oid,
  -- typename: Data type name.
  arr_typ.typname::text AS type_name,
  elem_typ.oid          AS elem_oid,
  -- typtype: b for a base type, c for a composite type (e.g., a table's
  -- row type), d for a domain, e for an enum type, p for a pseudo-type,
  -- or r for a range type.
  arr_typ.typtype       AS type_kind
FROM pg_type arr_typ
  JOIN pg_type elem_typ ON arr_typ.typelem = elem_typ.oid
WHERE arr_typ.typisdefined
  AND arr_typ.typtype = 'b' -- Array types are base types
  -- If typelem is not 0 then it identifies another row in pg_type. The current
  -- type can then be subscripted like an array yielding values of type typelem.
  -- A “true” array type is variable length (typlen = -1), but some
  -- fixed-length (typlen > 0) types also have nonzero typelem, for example
  -- name and point. If a fixed-length type has a typelem then its internal
  -- representation must be some number of values of the typelem data type with
  -- no other data. Variable-length array types have a header defined by the
  -- array subroutines.
  AND arr_typ.typelem > 0
  -- For a fixed-size type, typlen is the number of bytes in the internal
  -- representation of the type. But for a variable-length type, typlen is
  -- negative. -1 indicates a "varlena" type (one that has a length word), -2
  -- indicates a null-terminated C string.
  AND arr_typ.typlen = -1
  AND arr_typ.oid = ANY ($1::oid[]);`

type FindArrayTypesRow struct {
	OID      pgtype.OID   `json:"oid"`
	TypeName string       `json:"type_name"`
	ElemOID  pgtype.OID   `json:"elem_oid"`
	TypeKind pgtype.QChar `json:"type_kind"`
}

// FindArrayTypes implements Querier.FindArrayTypes.
func (q *DBQuerier) FindArrayTypes(ctx context.Context, oids []uint32) ([]FindArrayTypesRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindArrayTypes")
	rows, err := q.conn.Query(ctx, findArrayTypesSQL, oids)
	if err != nil {
		return nil, fmt.Errorf("query FindArrayTypes: %w", err)
	}
	defer rows.Close()
	items := []FindArrayTypesRow{}
	for rows.Next() {
		var item FindArrayTypesRow
		if err := rows.Scan(&item.OID, &item.TypeName, &item.ElemOID, &item.TypeKind); err != nil {
			return nil, fmt.Errorf("scan FindArrayTypes row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindArrayTypes rows: %w", err)
	}
	return items, err
}

// FindArrayTypesBatch implements Querier.FindArrayTypesBatch.
func (q *DBQuerier) FindArrayTypesBatch(batch *pgx.Batch, oids []uint32) {
	batch.Queue(findArrayTypesSQL, oids)
}

// FindArrayTypesScan implements Querier.FindArrayTypesScan.
func (q *DBQuerier) FindArrayTypesScan(results pgx.BatchResults) ([]FindArrayTypesRow, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query FindArrayTypesBatch: %w", err)
	}
	defer rows.Close()
	items := []FindArrayTypesRow{}
	for rows.Next() {
		var item FindArrayTypesRow
		if err := rows.Scan(&item.OID, &item.TypeName, &item.ElemOID, &item.TypeKind); err != nil {
			return nil, fmt.Errorf("scan FindArrayTypesBatch row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindArrayTypesBatch rows: %w", err)
	}
	return items, err
}

const findCompositeTypesSQL = `WITH table_cols AS (
  SELECT
    cls.relname                                         AS table_name,
    cls.oid                                             AS table_oid,
    array_agg(attr.attname::text ORDER BY attr.attnum)  AS col_names,
    array_agg(attr.atttypid::int8 ORDER BY attr.attnum) AS col_oids,
    array_agg(attr.attnum::int8 ORDER BY attr.attnum)   AS col_orders,
    array_agg(attr.attnotnull ORDER BY attr.attnum)     AS col_not_nulls,
    array_agg(typ.typname::text ORDER BY attr.attnum)   AS col_type_names
  FROM pg_attribute attr
    JOIN pg_class cls ON attr.attrelid = cls.oid
    JOIN pg_type typ ON typ.oid = attr.atttypid
  WHERE attr.attnum > 0 -- Postgres represents system columns with attnum <= 0
    AND NOT attr.attisdropped
  GROUP BY cls.relname, cls.oid
)
SELECT
  typ.typname::text AS table_type_name,
  typ.oid           AS table_type_oid,
  table_name,
  col_names,
  col_oids,
  col_orders,
  col_not_nulls,
  col_type_names
FROM pg_type typ
  JOIN table_cols cols ON typ.typrelid = cols.table_oid
WHERE typ.oid = ANY ($1::oid[])
  AND typ.typtype = 'c';`

type FindCompositeTypesRow struct {
	TableTypeName string           `json:"table_type_name"`
	TableTypeOID  pgtype.OID       `json:"table_type_oid"`
	TableName     pgtype.Name      `json:"table_name"`
	ColNames      []string         `json:"col_names"`
	ColOIDs       []int            `json:"col_oids"`
	ColOrders     []int            `json:"col_orders"`
	ColNotNulls   pgtype.BoolArray `json:"col_not_nulls"`
	ColTypeNames  []string         `json:"col_type_names"`
}

// FindCompositeTypes implements Querier.FindCompositeTypes.
func (q *DBQuerier) FindCompositeTypes(ctx context.Context, oids []uint32) ([]FindCompositeTypesRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindCompositeTypes")
	rows, err := q.conn.Query(ctx, findCompositeTypesSQL, oids)
	if err != nil {
		return nil, fmt.Errorf("query FindCompositeTypes: %w", err)
	}
	defer rows.Close()
	items := []FindCompositeTypesRow{}
	for rows.Next() {
		var item FindCompositeTypesRow
		if err := rows.Scan(&item.TableTypeName, &item.TableTypeOID, &item.TableName, &item.ColNames, &item.ColOIDs, &item.ColOrders, &item.ColNotNulls, &item.ColTypeNames); err != nil {
			return nil, fmt.Errorf("scan FindCompositeTypes row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindCompositeTypes rows: %w", err)
	}
	return items, err
}

// FindCompositeTypesBatch implements Querier.FindCompositeTypesBatch.
func (q *DBQuerier) FindCompositeTypesBatch(batch *pgx.Batch, oids []uint32) {
	batch.Queue(findCompositeTypesSQL, oids)
}

// FindCompositeTypesScan implements Querier.FindCompositeTypesScan.
func (q *DBQuerier) FindCompositeTypesScan(results pgx.BatchResults) ([]FindCompositeTypesRow, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query FindCompositeTypesBatch: %w", err)
	}
	defer rows.Close()
	items := []FindCompositeTypesRow{}
	for rows.Next() {
		var item FindCompositeTypesRow
		if err := rows.Scan(&item.TableTypeName, &item.TableTypeOID, &item.TableName, &item.ColNames, &item.ColOIDs, &item.ColOrders, &item.ColNotNulls, &item.ColTypeNames); err != nil {
			return nil, fmt.Errorf("scan FindCompositeTypesBatch row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindCompositeTypesBatch rows: %w", err)
	}
	return items, err
}

const findDescendantOIDsSQL = `WITH RECURSIVE oid_descs(oid) AS (
  -- Base case.
  SELECT oid
  FROM unnest($1::oid[]) AS t(oid)
  UNION
  -- Recursive case.
  SELECT oid
  FROM (
    WITH all_oids AS (SELECT oid FROM oid_descs)
    -- All composite children.
    SELECT attr.atttypid AS oid
    FROM pg_type typ
      JOIN pg_class cls ON typ.oid = cls.reltype
      JOIN pg_attribute attr ON attr.attrelid = cls.oid
      JOIN all_oids od ON typ.oid = od.oid
    WHERE attr.attnum > 0 -- Postgres represents system columns with attnum <= 0
      AND NOT attr.attisdropped
    UNION
    -- All array elements.
    SELECT elem_typ.oid
    FROM pg_type arr_typ
      JOIN pg_type elem_typ ON arr_typ.typelem = elem_typ.oid
      JOIN all_oids od ON arr_typ.oid = od.oid
  ) t
)
SELECT oid
FROM oid_descs;`

// FindDescendantOIDs implements Querier.FindDescendantOIDs.
func (q *DBQuerier) FindDescendantOIDs(ctx context.Context, oids []uint32) ([]pgtype.OID, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindDescendantOIDs")
	rows, err := q.conn.Query(ctx, findDescendantOIDsSQL, oids)
	if err != nil {
		return nil, fmt.Errorf("query FindDescendantOIDs: %w", err)
	}
	defer rows.Close()
	items := []pgtype.OID{}
	for rows.Next() {
		var item pgtype.OID
		if err := rows.Scan(&item); err != nil {
			return nil, fmt.Errorf("scan FindDescendantOIDs row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindDescendantOIDs rows: %w", err)
	}
	return items, err
}

// FindDescendantOIDsBatch implements Querier.FindDescendantOIDsBatch.
func (q *DBQuerier) FindDescendantOIDsBatch(batch *pgx.Batch, oids []uint32) {
	batch.Queue(findDescendantOIDsSQL, oids)
}

// FindDescendantOIDsScan implements Querier.FindDescendantOIDsScan.
func (q *DBQuerier) FindDescendantOIDsScan(results pgx.BatchResults) ([]pgtype.OID, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query FindDescendantOIDsBatch: %w", err)
	}
	defer rows.Close()
	items := []pgtype.OID{}
	for rows.Next() {
		var item pgtype.OID
		if err := rows.Scan(&item); err != nil {
			return nil, fmt.Errorf("scan FindDescendantOIDsBatch row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindDescendantOIDsBatch rows: %w", err)
	}
	return items, err
}

const findOIDByNameSQL = `SELECT oid
FROM pg_type
WHERE typname::text = $1
ORDER BY oid DESC
LIMIT 1;`

// FindOIDByName implements Querier.FindOIDByName.
func (q *DBQuerier) FindOIDByName(ctx context.Context, name string) (pgtype.OID, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindOIDByName")
	row := q.conn.QueryRow(ctx, findOIDByNameSQL, name)
	var item pgtype.OID
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query FindOIDByName: %w", err)
	}
	return item, nil
}

// FindOIDByNameBatch implements Querier.FindOIDByNameBatch.
func (q *DBQuerier) FindOIDByNameBatch(batch *pgx.Batch, name string) {
	batch.Queue(findOIDByNameSQL, name)
}

// FindOIDByNameScan implements Querier.FindOIDByNameScan.
func (q *DBQuerier) FindOIDByNameScan(results pgx.BatchResults) (pgtype.OID, error) {
	row := results.QueryRow()
	var item pgtype.OID
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan FindOIDByNameBatch row: %w", err)
	}
	return item, nil
}

const findOIDNameSQL = `SELECT typname AS name
FROM pg_type
WHERE oid = $1;`

// FindOIDName implements Querier.FindOIDName.
func (q *DBQuerier) FindOIDName(ctx context.Context, oid pgtype.OID) (pgtype.Name, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindOIDName")
	row := q.conn.QueryRow(ctx, findOIDNameSQL, oid)
	var item pgtype.Name
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query FindOIDName: %w", err)
	}
	return item, nil
}

// FindOIDNameBatch implements Querier.FindOIDNameBatch.
func (q *DBQuerier) FindOIDNameBatch(batch *pgx.Batch, oid pgtype.OID) {
	batch.Queue(findOIDNameSQL, oid)
}

// FindOIDNameScan implements Querier.FindOIDNameScan.
func (q *DBQuerier) FindOIDNameScan(results pgx.BatchResults) (pgtype.Name, error) {
	row := results.QueryRow()
	var item pgtype.Name
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan FindOIDNameBatch row: %w", err)
	}
	return item, nil
}

const findOIDNamesSQL = `SELECT oid, typname AS name, typtype AS kind
FROM pg_type
WHERE oid = ANY ($1::oid[]);`

type FindOIDNamesRow struct {
	OID  pgtype.OID   `json:"oid"`
	Name pgtype.Name  `json:"name"`
	Kind pgtype.QChar `json:"kind"`
}

// FindOIDNames implements Querier.FindOIDNames.
func (q *DBQuerier) FindOIDNames(ctx context.Context, oid []uint32) ([]FindOIDNamesRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindOIDNames")
	rows, err := q.conn.Query(ctx, findOIDNamesSQL, oid)
	if err != nil {
		return nil, fmt.Errorf("query FindOIDNames: %w", err)
	}
	defer rows.Close()
	items := []FindOIDNamesRow{}
	for rows.Next() {
		var item FindOIDNamesRow
		if err := rows.Scan(&item.OID, &item.Name, &item.Kind); err != nil {
			return nil, fmt.Errorf("scan FindOIDNames row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindOIDNames rows: %w", err)
	}
	return items, err
}

// FindOIDNamesBatch implements Querier.FindOIDNamesBatch.
func (q *DBQuerier) FindOIDNamesBatch(batch *pgx.Batch, oid []uint32) {
	batch.Queue(findOIDNamesSQL, oid)
}

// FindOIDNamesScan implements Querier.FindOIDNamesScan.
func (q *DBQuerier) FindOIDNamesScan(results pgx.BatchResults) ([]FindOIDNamesRow, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query FindOIDNamesBatch: %w", err)
	}
	defer rows.Close()
	items := []FindOIDNamesRow{}
	for rows.Next() {
		var item FindOIDNamesRow
		if err := rows.Scan(&item.OID, &item.Name, &item.Kind); err != nil {
			return nil, fmt.Errorf("scan FindOIDNamesBatch row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindOIDNamesBatch rows: %w", err)
	}
	return items, err
}

// textPreferrer wraps a pgtype.ValueTranscoder and sets the preferred encoding
// format to text instead binary (the default). pggen uses the text format
// when the OID is unknownOID because the binary format requires the OID.
// Typically occurs if the results from QueryAllDataTypes aren't passed to
// NewQuerierConfig.
type textPreferrer struct {
	pgtype.ValueTranscoder
	typeName string
}

// PreferredParamFormat implements pgtype.ParamFormatPreferrer.
func (t textPreferrer) PreferredParamFormat() int16 { return pgtype.TextFormatCode }

func (t textPreferrer) NewTypeValue() pgtype.Value {
	return textPreferrer{pgtype.NewValue(t.ValueTranscoder).(pgtype.ValueTranscoder), t.typeName}
}

func (t textPreferrer) TypeName() string {
	return t.typeName
}

// unknownOID means we don't know the OID for a type. This is okay for decoding
// because pgx call DecodeText or DecodeBinary without requiring the OID. For
// encoding parameters, pggen uses textPreferrer if the OID is unknown.
const unknownOID = 0
