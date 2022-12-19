// Code generated by pggen. DO NOT EDIT.

package enums

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
)

// Querier is a typesafe Go interface backed by SQL queries.
//
// Methods ending with Batch enqueue a query to run later in a pgx.Batch. After
// calling SendBatch on pgx.Conn, pgxpool.Pool, or pgx.Tx, use the Scan methods
// to parse the results.
type Querier interface {
	FindAllDevices(ctx context.Context) ([]FindAllDevicesRow, error)
	// FindAllDevicesBatch enqueues a FindAllDevices query into batch to be executed
	// later by the batch.
	FindAllDevicesBatch(batch genericBatch)
	// FindAllDevicesScan scans the result of an executed FindAllDevicesBatch query.
	FindAllDevicesScan(results pgx.BatchResults) ([]FindAllDevicesRow, error)

	InsertDevice(ctx context.Context, mac pgtype.Macaddr, typePg DeviceType) (pgconn.CommandTag, error)
	// InsertDeviceBatch enqueues a InsertDevice query into batch to be executed
	// later by the batch.
	InsertDeviceBatch(batch genericBatch, mac pgtype.Macaddr, typePg DeviceType)
	// InsertDeviceScan scans the result of an executed InsertDeviceBatch query.
	InsertDeviceScan(results pgx.BatchResults) (pgconn.CommandTag, error)

	// Select an array of all device_type enum values.
	FindOneDeviceArray(ctx context.Context) ([]DeviceType, error)
	// FindOneDeviceArrayBatch enqueues a FindOneDeviceArray query into batch to be executed
	// later by the batch.
	FindOneDeviceArrayBatch(batch genericBatch)
	// FindOneDeviceArrayScan scans the result of an executed FindOneDeviceArrayBatch query.
	FindOneDeviceArrayScan(results pgx.BatchResults) ([]DeviceType, error)

	// Select many rows of device_type enum values.
	FindManyDeviceArray(ctx context.Context) ([][]DeviceType, error)
	// FindManyDeviceArrayBatch enqueues a FindManyDeviceArray query into batch to be executed
	// later by the batch.
	FindManyDeviceArrayBatch(batch genericBatch)
	// FindManyDeviceArrayScan scans the result of an executed FindManyDeviceArrayBatch query.
	FindManyDeviceArrayScan(results pgx.BatchResults) ([][]DeviceType, error)

	// Select many rows of device_type enum values with multiple output columns.
	FindManyDeviceArrayWithNum(ctx context.Context) ([]FindManyDeviceArrayWithNumRow, error)
	// FindManyDeviceArrayWithNumBatch enqueues a FindManyDeviceArrayWithNum query into batch to be executed
	// later by the batch.
	FindManyDeviceArrayWithNumBatch(batch genericBatch)
	// FindManyDeviceArrayWithNumScan scans the result of an executed FindManyDeviceArrayWithNumBatch query.
	FindManyDeviceArrayWithNumScan(results pgx.BatchResults) ([]FindManyDeviceArrayWithNumRow, error)

	// Regression test for https://github.com/eddiefisher/pggen/issues/23.
	EnumInsideComposite(ctx context.Context) (Device, error)
	// EnumInsideCompositeBatch enqueues a EnumInsideComposite query into batch to be executed
	// later by the batch.
	EnumInsideCompositeBatch(batch genericBatch)
	// EnumInsideCompositeScan scans the result of an executed EnumInsideCompositeBatch query.
	EnumInsideCompositeScan(results pgx.BatchResults) (Device, error)
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
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

// genericBatch batches queries to send in a single network request to a
// Postgres server. This is usually backed by *pgx.Batch.
type genericBatch interface {
	// Queue queues a query to batch b. query can be an SQL query or the name of a
	// prepared statement. See Queue on *pgx.Batch.
	Queue(query string, arguments ...any) *pgx.QueuedQuery
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
	if _, err := p.Prepare(ctx, findAllDevicesSQL, findAllDevicesSQL); err != nil {
		return fmt.Errorf("prepare query 'FindAllDevices': %w", err)
	}
	if _, err := p.Prepare(ctx, insertDeviceSQL, insertDeviceSQL); err != nil {
		return fmt.Errorf("prepare query 'InsertDevice': %w", err)
	}
	if _, err := p.Prepare(ctx, findOneDeviceArraySQL, findOneDeviceArraySQL); err != nil {
		return fmt.Errorf("prepare query 'FindOneDeviceArray': %w", err)
	}
	if _, err := p.Prepare(ctx, findManyDeviceArraySQL, findManyDeviceArraySQL); err != nil {
		return fmt.Errorf("prepare query 'FindManyDeviceArray': %w", err)
	}
	if _, err := p.Prepare(ctx, findManyDeviceArrayWithNumSQL, findManyDeviceArrayWithNumSQL); err != nil {
		return fmt.Errorf("prepare query 'FindManyDeviceArrayWithNum': %w", err)
	}
	if _, err := p.Prepare(ctx, enumInsideCompositeSQL, enumInsideCompositeSQL); err != nil {
		return fmt.Errorf("prepare query 'EnumInsideComposite': %w", err)
	}
	return nil
}

// Device represents the Postgres composite type "device".
type Device struct {
	Mac  pgtype.Macaddr `json:"mac"`
	Type DeviceType     `json:"type"`
}

// newDeviceTypeEnum creates a new pgtype.ValueTranscoder for the
// Postgres enum type 'device_type'.
func newDeviceTypeEnum() pgtype.ValueTranscoder {
	return pgtype.NewEnumType(
		"device_type",
		[]string{
			string(DeviceTypeUndefined),
			string(DeviceTypePhone),
			string(DeviceTypeLaptop),
			string(DeviceTypeIpad),
			string(DeviceTypeDesktop),
			string(DeviceTypeIot),
		},
	)
}

// DeviceType represents the Postgres enum "device_type".
type DeviceType string

const (
	DeviceTypeUndefined DeviceType = "undefined"
	DeviceTypePhone     DeviceType = "phone"
	DeviceTypeLaptop    DeviceType = "laptop"
	DeviceTypeIpad      DeviceType = "ipad"
	DeviceTypeDesktop   DeviceType = "desktop"
	DeviceTypeIot       DeviceType = "iot"
)

func (d DeviceType) String() string { return string(d) }

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

type compositeField struct {
	name       string                 // name of the field
	typeName   string                 // Postgres type name
	defaultVal pgtype.ValueTranscoder // default value to use
}

func (tr *typeResolver) newCompositeValue(name string, fields ...compositeField) pgtype.ValueTranscoder {
	if _, val, ok := tr.findValue(name); ok {
		return val
	}
	fs := make([]pgtype.CompositeTypeField, len(fields))
	vals := make([]pgtype.ValueTranscoder, len(fields))
	isBinaryOk := true
	for i, field := range fields {
		oid, val, ok := tr.findValue(field.typeName)
		if !ok {
			oid = unknownOID
			val = field.defaultVal
		}
		isBinaryOk = isBinaryOk && oid != unknownOID
		fs[i] = pgtype.CompositeTypeField{Name: field.name, OID: oid}
		vals[i] = val
	}
	// Okay to ignore error because it's only thrown when the number of field
	// names does not equal the number of ValueTranscoders.
	typ, _ := pgtype.NewCompositeTypeValues(name, fs, vals)
	if !isBinaryOk {
		return textPreferrer{ValueTranscoder: typ, typeName: name}
	}
	return typ
}

func (tr *typeResolver) newArrayValue(name, elemName string, defaultVal func() pgtype.ValueTranscoder) pgtype.ValueTranscoder {
	if _, val, ok := tr.findValue(name); ok {
		return val
	}
	elemOID, elemVal, ok := tr.findValue(elemName)
	elemValFunc := func() pgtype.ValueTranscoder {
		return pgtype.NewValue(elemVal).(pgtype.ValueTranscoder)
	}
	if !ok {
		elemOID = unknownOID
		elemValFunc = defaultVal
	}
	typ := pgtype.NewArrayType(name, elemOID, elemValFunc)
	if elemOID == unknownOID {
		return textPreferrer{ValueTranscoder: typ, typeName: name}
	}
	return typ
}

// newDevice creates a new pgtype.ValueTranscoder for the Postgres
// composite type 'device'.
func (tr *typeResolver) newDevice() pgtype.ValueTranscoder {
	return tr.newCompositeValue(
		"device",
		compositeField{name: "mac", typeName: "macaddr", defaultVal: &pgtype.Macaddr{}},
		compositeField{name: "type", typeName: "device_type", defaultVal: newDeviceTypeEnum()},
	)
}

// newDeviceTypeArray creates a new pgtype.ValueTranscoder for the Postgres
// '_device_type' array type.
func (tr *typeResolver) newDeviceTypeArray() pgtype.ValueTranscoder {
	return tr.newArrayValue("_device_type", "device_type", newDeviceTypeEnum)
}

const findAllDevicesSQL = `SELECT mac, type
FROM device;`

type FindAllDevicesRow struct {
	Mac  pgtype.Macaddr `json:"mac"`
	Type DeviceType     `json:"type"`
}

// FindAllDevices implements Querier.FindAllDevices.
func (q *DBQuerier) FindAllDevices(ctx context.Context) ([]FindAllDevicesRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindAllDevices")
	rows, err := q.conn.Query(ctx, findAllDevicesSQL)
	if err != nil {
		return nil, fmt.Errorf("query FindAllDevices: %w", err)
	}
	defer rows.Close()
	items := []FindAllDevicesRow{}
	for rows.Next() {
		var item FindAllDevicesRow
		if err := rows.Scan(&item.Mac, &item.Type); err != nil {
			return nil, fmt.Errorf("scan FindAllDevices row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindAllDevices rows: %w", err)
	}
	return items, err
}

// FindAllDevicesBatch implements Querier.FindAllDevicesBatch.
func (q *DBQuerier) FindAllDevicesBatch(batch genericBatch) {
	batch.Queue(findAllDevicesSQL)
}

// FindAllDevicesScan implements Querier.FindAllDevicesScan.
func (q *DBQuerier) FindAllDevicesScan(results pgx.BatchResults) ([]FindAllDevicesRow, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query FindAllDevicesBatch: %w", err)
	}
	defer rows.Close()
	items := []FindAllDevicesRow{}
	for rows.Next() {
		var item FindAllDevicesRow
		if err := rows.Scan(&item.Mac, &item.Type); err != nil {
			return nil, fmt.Errorf("scan FindAllDevicesBatch row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindAllDevicesBatch rows: %w", err)
	}
	return items, err
}

const insertDeviceSQL = `INSERT INTO device (mac, type)
VALUES ($1, $2);`

// InsertDevice implements Querier.InsertDevice.
func (q *DBQuerier) InsertDevice(ctx context.Context, mac pgtype.Macaddr, typePg DeviceType) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertDevice")
	cmdTag, err := q.conn.Exec(ctx, insertDeviceSQL, mac, typePg)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query InsertDevice: %w", err)
	}
	return cmdTag, err
}

// InsertDeviceBatch implements Querier.InsertDeviceBatch.
func (q *DBQuerier) InsertDeviceBatch(batch genericBatch, mac pgtype.Macaddr, typePg DeviceType) {
	batch.Queue(insertDeviceSQL, mac, typePg)
}

// InsertDeviceScan implements Querier.InsertDeviceScan.
func (q *DBQuerier) InsertDeviceScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec InsertDeviceBatch: %w", err)
	}
	return cmdTag, err
}

const findOneDeviceArraySQL = `SELECT enum_range(NULL::device_type) AS device_types;`

// FindOneDeviceArray implements Querier.FindOneDeviceArray.
func (q *DBQuerier) FindOneDeviceArray(ctx context.Context) ([]DeviceType, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindOneDeviceArray")
	row := q.conn.QueryRow(ctx, findOneDeviceArraySQL)
	item := []DeviceType{}
	deviceTypesArray := q.types.newDeviceTypeArray()
	if err := row.Scan(deviceTypesArray); err != nil {
		return item, fmt.Errorf("query FindOneDeviceArray: %w", err)
	}
	if err := deviceTypesArray.AssignTo(&item); err != nil {
		return item, fmt.Errorf("assign FindOneDeviceArray row: %w", err)
	}
	return item, nil
}

// FindOneDeviceArrayBatch implements Querier.FindOneDeviceArrayBatch.
func (q *DBQuerier) FindOneDeviceArrayBatch(batch genericBatch) {
	batch.Queue(findOneDeviceArraySQL)
}

// FindOneDeviceArrayScan implements Querier.FindOneDeviceArrayScan.
func (q *DBQuerier) FindOneDeviceArrayScan(results pgx.BatchResults) ([]DeviceType, error) {
	row := results.QueryRow()
	item := []DeviceType{}
	deviceTypesArray := q.types.newDeviceTypeArray()
	if err := row.Scan(deviceTypesArray); err != nil {
		return item, fmt.Errorf("scan FindOneDeviceArrayBatch row: %w", err)
	}
	if err := deviceTypesArray.AssignTo(&item); err != nil {
		return item, fmt.Errorf("assign FindOneDeviceArray row: %w", err)
	}
	return item, nil
}

const findManyDeviceArraySQL = `SELECT enum_range('ipad'::device_type, 'iot'::device_type) AS device_types
UNION ALL
SELECT enum_range(NULL::device_type) AS device_types;`

// FindManyDeviceArray implements Querier.FindManyDeviceArray.
func (q *DBQuerier) FindManyDeviceArray(ctx context.Context) ([][]DeviceType, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindManyDeviceArray")
	rows, err := q.conn.Query(ctx, findManyDeviceArraySQL)
	if err != nil {
		return nil, fmt.Errorf("query FindManyDeviceArray: %w", err)
	}
	defer rows.Close()
	items := [][]DeviceType{}
	deviceTypesArray := q.types.newDeviceTypeArray()
	for rows.Next() {
		var item []DeviceType
		if err := rows.Scan(deviceTypesArray); err != nil {
			return nil, fmt.Errorf("scan FindManyDeviceArray row: %w", err)
		}
		if err := deviceTypesArray.AssignTo(&item); err != nil {
			return nil, fmt.Errorf("assign FindManyDeviceArray row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindManyDeviceArray rows: %w", err)
	}
	return items, err
}

// FindManyDeviceArrayBatch implements Querier.FindManyDeviceArrayBatch.
func (q *DBQuerier) FindManyDeviceArrayBatch(batch genericBatch) {
	batch.Queue(findManyDeviceArraySQL)
}

// FindManyDeviceArrayScan implements Querier.FindManyDeviceArrayScan.
func (q *DBQuerier) FindManyDeviceArrayScan(results pgx.BatchResults) ([][]DeviceType, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query FindManyDeviceArrayBatch: %w", err)
	}
	defer rows.Close()
	items := [][]DeviceType{}
	deviceTypesArray := q.types.newDeviceTypeArray()
	for rows.Next() {
		var item []DeviceType
		if err := rows.Scan(deviceTypesArray); err != nil {
			return nil, fmt.Errorf("scan FindManyDeviceArrayBatch row: %w", err)
		}
		if err := deviceTypesArray.AssignTo(&item); err != nil {
			return nil, fmt.Errorf("assign FindManyDeviceArray row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindManyDeviceArrayBatch rows: %w", err)
	}
	return items, err
}

const findManyDeviceArrayWithNumSQL = `SELECT 1 AS num, enum_range('ipad'::device_type, 'iot'::device_type) AS device_types
UNION ALL
SELECT 2 as num, enum_range(NULL::device_type) AS device_types;`

type FindManyDeviceArrayWithNumRow struct {
	Num         *int32       `json:"num"`
	DeviceTypes []DeviceType `json:"device_types"`
}

// FindManyDeviceArrayWithNum implements Querier.FindManyDeviceArrayWithNum.
func (q *DBQuerier) FindManyDeviceArrayWithNum(ctx context.Context) ([]FindManyDeviceArrayWithNumRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindManyDeviceArrayWithNum")
	rows, err := q.conn.Query(ctx, findManyDeviceArrayWithNumSQL)
	if err != nil {
		return nil, fmt.Errorf("query FindManyDeviceArrayWithNum: %w", err)
	}
	defer rows.Close()
	items := []FindManyDeviceArrayWithNumRow{}
	deviceTypesArray := q.types.newDeviceTypeArray()
	for rows.Next() {
		var item FindManyDeviceArrayWithNumRow
		if err := rows.Scan(&item.Num, deviceTypesArray); err != nil {
			return nil, fmt.Errorf("scan FindManyDeviceArrayWithNum row: %w", err)
		}
		if err := deviceTypesArray.AssignTo(&item.DeviceTypes); err != nil {
			return nil, fmt.Errorf("assign FindManyDeviceArrayWithNum row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindManyDeviceArrayWithNum rows: %w", err)
	}
	return items, err
}

// FindManyDeviceArrayWithNumBatch implements Querier.FindManyDeviceArrayWithNumBatch.
func (q *DBQuerier) FindManyDeviceArrayWithNumBatch(batch genericBatch) {
	batch.Queue(findManyDeviceArrayWithNumSQL)
}

// FindManyDeviceArrayWithNumScan implements Querier.FindManyDeviceArrayWithNumScan.
func (q *DBQuerier) FindManyDeviceArrayWithNumScan(results pgx.BatchResults) ([]FindManyDeviceArrayWithNumRow, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query FindManyDeviceArrayWithNumBatch: %w", err)
	}
	defer rows.Close()
	items := []FindManyDeviceArrayWithNumRow{}
	deviceTypesArray := q.types.newDeviceTypeArray()
	for rows.Next() {
		var item FindManyDeviceArrayWithNumRow
		if err := rows.Scan(&item.Num, deviceTypesArray); err != nil {
			return nil, fmt.Errorf("scan FindManyDeviceArrayWithNumBatch row: %w", err)
		}
		if err := deviceTypesArray.AssignTo(&item.DeviceTypes); err != nil {
			return nil, fmt.Errorf("assign FindManyDeviceArrayWithNum row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindManyDeviceArrayWithNumBatch rows: %w", err)
	}
	return items, err
}

const enumInsideCompositeSQL = `SELECT ROW('08:00:2b:01:02:03'::macaddr, 'phone'::device_type) ::device;`

// EnumInsideComposite implements Querier.EnumInsideComposite.
func (q *DBQuerier) EnumInsideComposite(ctx context.Context) (Device, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "EnumInsideComposite")
	row := q.conn.QueryRow(ctx, enumInsideCompositeSQL)
	var item Device
	rowRow := q.types.newDevice()
	if err := row.Scan(rowRow); err != nil {
		return item, fmt.Errorf("query EnumInsideComposite: %w", err)
	}
	if err := rowRow.AssignTo(&item); err != nil {
		return item, fmt.Errorf("assign EnumInsideComposite row: %w", err)
	}
	return item, nil
}

// EnumInsideCompositeBatch implements Querier.EnumInsideCompositeBatch.
func (q *DBQuerier) EnumInsideCompositeBatch(batch genericBatch) {
	batch.Queue(enumInsideCompositeSQL)
}

// EnumInsideCompositeScan implements Querier.EnumInsideCompositeScan.
func (q *DBQuerier) EnumInsideCompositeScan(results pgx.BatchResults) (Device, error) {
	row := results.QueryRow()
	var item Device
	rowRow := q.types.newDevice()
	if err := row.Scan(rowRow); err != nil {
		return item, fmt.Errorf("scan EnumInsideCompositeBatch row: %w", err)
	}
	if err := rowRow.AssignTo(&item); err != nil {
		return item, fmt.Errorf("assign EnumInsideComposite row: %w", err)
	}
	return item, nil
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
	return textPreferrer{ValueTranscoder: pgtype.NewValue(t.ValueTranscoder).(pgtype.ValueTranscoder), typeName: t.typeName}
}

func (t textPreferrer) TypeName() string {
	return t.typeName
}

// unknownOID means we don't know the OID for a type. This is okay for decoding
// because pgx call DecodeText or DecodeBinary without requiring the OID. For
// encoding parameters, pggen uses textPreferrer if the OID is unknown.
const unknownOID = 0
