package indexer

import (
	"bufio"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB is a global variable to hold the GORM database connection.

func LabelsTableName(blockchain string) string {
	return fmt.Sprintf(blockchain + "_labels")
}

func TransactionsTableName(blockchain string) string {
	return fmt.Sprintf(blockchain + "_transactions")
}

func LogsTableName(blockchain string) string {
	return fmt.Sprintf(blockchain + "_logs")
}

func BlocksTableName(blockchain string) string {
	return fmt.Sprintf(blockchain + "_blocks")
}

func hexStringToInt(hexString string) (int64, error) {
	// Remove the "0x" prefix from the hexadecimal string
	hexString = strings.TrimPrefix(hexString, "0x")

	// Parse the hexadecimal string to an integer
	intValue, err := strconv.ParseInt(hexString, 16, 64)
	if err != nil {
		return 0, err
	}

	return intValue, nil
}

// https://klotzandrew.com/blog/postgres-passing-65535-parameter-limit/ insted of batching
type UnnestInsertValueStruct struct {
	Type   string        `json:"type"`   // e.g. "BIGINT" or "TEXT" or any other PostgreSQL data type
	Values []interface{} `json:"values"` // e.g. [1, 2, 3, 4, 5]
}

func IsBlockchainWithL1Chain(blockchain string) bool {
	switch blockchain {
	case "ethereum":
		return false
	case "polygon":
		return false
	case "arbitrum_one":
		return true
	case "arbitrum_sepolia":
		return true
	case "game7_orbit_arbitrum_sepolia":
		return true
	case "game7_testnet":
		return true
	case "xai":
		return true
	case "xai_sepolia":
		return true
	case "mantle":
		return false
	case "mantle_sepolia":
		return false
	default:
		return false
	}
}

type PostgreSQLpgx struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLpgx() (*PostgreSQLpgx, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)

	defer cancel()

	pool, err := pgxpool.New(ctx, MOONSTREAM_DB_V3_INDEXES_URI)
	if err != nil {
		return nil, err
	}

	return &PostgreSQLpgx{
		pool: pool,
	}, nil
}

func NewPostgreSQLpgxWithCustomURI(uri string) (*PostgreSQLpgx, error) {

	//  create a connection to the database

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	pool, err := pgxpool.New(ctx, uri)
	if err != nil {
		return nil, err
	}

	return &PostgreSQLpgx{
		pool: pool,
	}, nil
}

func (p *PostgreSQLpgx) Close() {
	p.pool.Close()
}

func (p *PostgreSQLpgx) GetPool() *pgxpool.Pool {
	return p.pool
}

// read from database

func (p *PostgreSQLpgx) ReadBlockIndex(ctx context.Context, startBlock uint64, endBlock uint64) ([]BlockIndex, error) {

	pool := p.GetPool()

	conn, err := pool.Acquire(ctx)

	if err != nil {
		return nil, err
	}

	defer conn.Release()

	rows, err := conn.Query(ctx, "SELECT * FROM products WHERE block_number >= $1 AND block_number <= $2", startBlock, endBlock)

	if err != nil {
		return nil, err
	}

	blocksIndex, err := pgx.CollectRows(rows, pgx.RowToStructByName[BlockIndex])

	if err != nil {
		return nil, err
	}

	return blocksIndex, nil

}

func (p *PostgreSQLpgx) ReadTransactionIndex(startBlock uint64, endBlock uint64) ([]TransactionIndex, error) {
	pool := p.GetPool()

	conn, err := pool.Acquire(context.Background())

	if err != nil {
		return nil, err
	}

	defer conn.Release()

	rows, err := conn.Query(context.Background(), "SELECT * FROM transactions WHERE block_number >= $1 AND block_number <= $2", startBlock, endBlock)

	if err != nil {
		return nil, err
	}

	transactionsIndex, err := pgx.CollectRows(rows, pgx.RowToStructByName[TransactionIndex])

	if err != nil {
		return nil, err
	}

	return transactionsIndex, nil

}

func (p *PostgreSQLpgx) ReadLogIndex(startBlock uint64, endBlock uint64, addresses []string) ([]LogIndex, error) {
	pool := p.GetPool()

	conn, err := pool.Acquire(context.Background())

	if err != nil {

		return nil, err

	}

	defer conn.Release()

	rows, err := conn.Query(context.Background(), "SELECT * FROM logs WHERE block_number >= $1 AND block_number <= $2 AND address = ANY($3)", startBlock, endBlock, addresses)

	if err != nil {

		return nil, err

	}

	logsIndex, err := pgx.CollectRows(rows, pgx.RowToStructByName[LogIndex])

	if err != nil {

		return nil, err

	}

	return logsIndex, nil
}

func (p *PostgreSQLpgx) ReadIndexOnRange(tableName string, startBlock uint64, endBlock uint64) ([]interface{}, error) {
	pool := p.GetPool()

	conn, err := pool.Acquire(context.Background())

	if err != nil {
		return nil, err
	}

	defer conn.Release()

	var indices []interface{}

	rows, err := conn.Query(context.Background(), "SELECT bt.block_number, bt.block_hash, bt.block_timestamp, tt.hash, tt.index, tt.path as transaction_, tt.input as transaction_input, lt.selector, lt.topic1, lt.topic2, lt.transaction_hash, lt.log_index, lt.path as event_path FROM block_index bt LEFT JOIN transaction_index tt ON bt.block_number = tt.block_number LEFT JOIN log_index lt ON tt.hash = lt.transaction_hash WHERE bt.block_number >= $1 AND bt.block_number <= $2", startBlock, endBlock)

	if err != nil {
		return nil, err
	}

	for rows.Next() {

		var index interface{}

		err = rows.Scan(&index)

		if err != nil {
			return nil, err
		}

		indices = append(indices, index)
	}

	return indices, nil
}

func (p *PostgreSQLpgx) ReadLastLabel(blockchain string) (uint64, error) {
	pool := p.GetPool()

	conn, err := pool.Acquire(context.Background())

	if err != nil {
		return 0, err
	}

	defer conn.Release()

	var label uint64

	query := fmt.Sprintf("SELECT block_number FROM %s ORDER BY block_number DESC LIMIT 1", blockchain+"_labels")

	err = conn.QueryRow(context.Background(), query).Scan(&label)

	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}
		return 0, err

	}

	return label, nil
}

func decodeAddress(address string) ([]byte, error) {
	if len(address) < 2 {
		return []byte{0x00}, nil
	}
	return hex.DecodeString(address[2:])
}

// updateValues updates the values in the map for a given key
func updateValues(valuesMap map[string]UnnestInsertValueStruct, key string, value interface{}) {
	tmp := valuesMap[key]
	tmp.Values = append(tmp.Values, value)
	valuesMap[key] = tmp
}

func (p *PostgreSQLpgx) WriteIndexes(blockchain string, blocksIndexPack []BlockIndex, transactionsIndexPack []TransactionIndex, logsIndexPack []LogIndex) error {

	ctx := context.Background()
	pool := p.GetPool()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		fmt.Println("Connection error", err)
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback(ctx)
			panic(err)
		} else if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	// Write blocks index
	if len(blocksIndexPack) > 0 {
		err = p.writeBlockIndexToDB(tx, blockchain, blocksIndexPack)
		if err != nil {
			return err
		}
	}

	// Write transactions index
	if len(transactionsIndexPack) > 0 {
		err = p.writeTransactionIndexToDB(tx, blockchain, transactionsIndexPack)
		if err != nil {
			return err
		}
	}

	// Write logs index
	if len(logsIndexPack) > 0 {
		err = p.writeLogIndexToDB(tx, blockchain, logsIndexPack)
		if err != nil {
			return err
		}
	}

	return nil
}

// Batch insert
func (p *PostgreSQLpgx) executeBatchInsert(tx pgx.Tx, ctx context.Context, tableName string, columns []string, values map[string]UnnestInsertValueStruct, conflictClause string) error {

	types := make([]string, 0)

	for index, column := range columns {
		// constract  unnest($1::int[], $2::int[] ...)
		types = append(types, fmt.Sprintf("$%d::%s[]", index+1, values[column].Type))
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) SELECT * FROM unnest(%s) %s", tableName, strings.Join(columns, ","), strings.Join(types, ","), conflictClause)

	// create a slices of values
	var valuesSlice []interface{}

	for _, column := range columns {

		valuesSlice = append(valuesSlice, values[column].Values)
	}

	// track execution time

	if _, err := tx.Exec(ctx, query, valuesSlice...); err != nil {
		fmt.Println("Error executing bulk insert", err)
		return fmt.Errorf("error executing bulk insert for batch: %w", err)
	}

	return nil
}

func (p *PostgreSQLpgx) writeBlockIndexToDB(tx pgx.Tx, blockchain string, indexes []BlockIndex) error {
	tableName := BlocksTableName(blockchain)
	isBlockchainWithL1Chain := IsBlockchainWithL1Chain(blockchain)
	columns := []string{"block_number", "block_hash", "block_timestamp", "parent_hash", "row_id", "path"}

	valuesMap := make(map[string]UnnestInsertValueStruct)

	valuesMap["block_number"] = UnnestInsertValueStruct{
		Type:   "BIGINT",
		Values: []interface{}{},
	}

	valuesMap["block_hash"] = UnnestInsertValueStruct{
		Type:   "TEXT",
		Values: make([]interface{}, 0),
	}

	valuesMap["block_timestamp"] = UnnestInsertValueStruct{
		Type:   "BIGINT",
		Values: make([]interface{}, 0),
	}

	valuesMap["parent_hash"] = UnnestInsertValueStruct{
		Type:   "TEXT",
		Values: make([]interface{}, 0),
	}

	valuesMap["row_id"] = UnnestInsertValueStruct{
		Type:   "BIGINT",
		Values: make([]interface{}, 0),
	}

	valuesMap["path"] = UnnestInsertValueStruct{
		Type:   "TEXT",
		Values: make([]interface{}, 0),
	}

	if isBlockchainWithL1Chain {
		columns = append(columns, "l1_block_number")
		valuesMap["l1_block_number"] = UnnestInsertValueStruct{
			Type:   "BIGINT",
			Values: make([]interface{}, 0),
		}
	}

	for _, index := range indexes {

		updateValues(valuesMap, "block_number", index.BlockNumber)
		updateValues(valuesMap, "block_hash", index.BlockHash)
		updateValues(valuesMap, "block_timestamp", index.BlockTimestamp)
		updateValues(valuesMap, "parent_hash", index.ParentHash)
		updateValues(valuesMap, "row_id", index.RowID)
		updateValues(valuesMap, "path", index.Path)

		if isBlockchainWithL1Chain {
			updateValues(valuesMap, "l1_block_number", index.L1BlockNumber)
		}
	}

	ctx := context.Background()
	err = p.executeBatchInsert(tx, ctx, tableName, columns, valuesMap, "ON CONFLICT (block_number) DO NOTHING")

	if err != nil {
		return err
	}

	log.Printf("Add %d records into %s table", len(indexes), tableName)

	return nil
}

func (p *PostgreSQLpgx) writeTransactionIndexToDB(tx pgx.Tx, blockchain string, indexes []TransactionIndex) error {

	tableName := TransactionsTableName(blockchain)

	columns := []string{"block_number", "block_hash", "hash", "index", "type", "from_address", "to_address", "selector", "row_id", "path"}
	var valuesMap = make(map[string]UnnestInsertValueStruct)

	valuesMap["block_number"] = UnnestInsertValueStruct{
		Type:   "BIGINT",
		Values: make([]interface{}, 0),
	}

	valuesMap["block_hash"] = UnnestInsertValueStruct{
		Type:   "TEXT",
		Values: make([]interface{}, 0),
	}

	valuesMap["hash"] = UnnestInsertValueStruct{
		Type:   "TEXT",
		Values: make([]interface{}, 0),
	}

	valuesMap["index"] = UnnestInsertValueStruct{
		Type:   "BIGINT",
		Values: make([]interface{}, 0),
	}

	valuesMap["type"] = UnnestInsertValueStruct{
		Type:   "INT",
		Values: make([]interface{}, 0),
	}

	valuesMap["from_address"] = UnnestInsertValueStruct{
		Type:   "BYTEA",
		Values: make([]interface{}, 0),
	}

	valuesMap["to_address"] = UnnestInsertValueStruct{
		Type:   "BYTEA",
		Values: make([]interface{}, 0),
	}

	valuesMap["selector"] = UnnestInsertValueStruct{
		Type:   "TEXT",
		Values: make([]interface{}, 0),
	}

	valuesMap["row_id"] = UnnestInsertValueStruct{
		Type:   "BIGINT",
		Values: make([]interface{}, 0),
	}

	valuesMap["path"] = UnnestInsertValueStruct{
		Type:   "TEXT",
		Values: make([]interface{}, 0),
	}

	for _, index := range indexes {

		fromAddressBytes, err := decodeAddress(index.FromAddress)
		if err != nil {
			fmt.Println("Error decoding from address:", err, index)
			continue
		}

		toAddressBytes, err := decodeAddress(index.ToAddress)

		if err != nil {
			fmt.Println("Error decoding to address:", err, index)
			continue
		}

		updateValues(valuesMap, "block_number", index.BlockNumber)
		updateValues(valuesMap, "block_hash", index.BlockHash)
		updateValues(valuesMap, "hash", index.TransactionHash)
		updateValues(valuesMap, "index", index.TransactionIndex)
		updateValues(valuesMap, "type", index.Type)
		updateValues(valuesMap, "from_address", fromAddressBytes)
		updateValues(valuesMap, "to_address", toAddressBytes)
		updateValues(valuesMap, "selector", index.Selector)
		updateValues(valuesMap, "row_id", index.RowID)
		updateValues(valuesMap, "path", index.Path)

	}

	ctx := context.Background()

	err = p.executeBatchInsert(tx, ctx, tableName, columns, valuesMap, "ON CONFLICT (hash) DO NOTHING")

	if err != nil {
		return err
	}

	log.Printf("Add %d records into %s table", len(indexes), tableName)

	return nil
}

func (p *PostgreSQLpgx) writeLogIndexToDB(tx pgx.Tx, blockchain string, indexes []LogIndex) error {

	tableName := LogsTableName(blockchain)

	columns := []string{"transaction_hash", "block_hash", "address", "selector", "topic1", "topic2", "topic3", "row_id", "log_index", "path"}

	var valuesMap = make(map[string]UnnestInsertValueStruct)

	valuesMap["transaction_hash"] = UnnestInsertValueStruct{
		Type:   "TEXT",
		Values: make([]interface{}, 0),
	}

	valuesMap["block_hash"] = UnnestInsertValueStruct{
		Type:   "TEXT",
		Values: make([]interface{}, 0),
	}

	valuesMap["address"] = UnnestInsertValueStruct{
		Type:   "BYTEA",
		Values: make([]interface{}, 0),
	}

	valuesMap["selector"] = UnnestInsertValueStruct{
		Type:   "TEXT",
		Values: make([]interface{}, 0),
	}

	valuesMap["topic1"] = UnnestInsertValueStruct{
		Type:   "TEXT",
		Values: make([]interface{}, 0),
	}

	valuesMap["topic2"] = UnnestInsertValueStruct{
		Type:   "TEXT",
		Values: make([]interface{}, 0),
	}

	valuesMap["topic3"] = UnnestInsertValueStruct{
		Type:   "TEXT",
		Values: make([]interface{}, 0),
	}

	valuesMap["row_id"] = UnnestInsertValueStruct{
		Type:   "BIGINT",
		Values: make([]interface{}, 0),
	}

	valuesMap["log_index"] = UnnestInsertValueStruct{
		Type:   "BIGINT",
		Values: make([]interface{}, 0),
	}

	valuesMap["path"] = UnnestInsertValueStruct{
		Type:   "TEXT",
		Values: make([]interface{}, 0),
	}

	for _, index := range indexes {

		toAddressBytes, err := decodeAddress(index.Address)
		if err != nil {
			fmt.Println("Error decoding address:", err, index)
			continue
		}

		updateValues(valuesMap, "transaction_hash", index.TransactionHash)
		updateValues(valuesMap, "block_hash", index.BlockHash)
		updateValues(valuesMap, "address", toAddressBytes)
		updateValues(valuesMap, "selector", index.Selector)
		updateValues(valuesMap, "topic1", index.Topic1)
		updateValues(valuesMap, "topic2", index.Topic2)
		updateValues(valuesMap, "topic3", index.Topic3)
		updateValues(valuesMap, "row_id", index.RowID)
		updateValues(valuesMap, "log_index", index.LogIndex)
		updateValues(valuesMap, "path", index.Path)

	}

	ctx := context.Background()

	err = p.executeBatchInsert(tx, ctx, tableName, columns, valuesMap, "ON CONFLICT (transaction_hash, log_index) DO NOTHING")

	if err != nil {
		return err
	}

	log.Printf("Add %d records into %s table", len(indexes), tableName)

	return nil
}

// GetEdgeDBBlock fetch first or last block for specified blockchain
func (p *PostgreSQLpgx) GetEdgeDBBlock(ctx context.Context, blockchain, side string) (BlockIndex, error) {
	var blockIndex BlockIndex

	pool := p.GetPool()

	conn, acquireErr := pool.Acquire(ctx)
	if acquireErr != nil {
		return blockIndex, acquireErr
	}

	defer conn.Release()

	query := fmt.Sprintf("SELECT block_number,block_hash,block_timestamp,parent_hash,row_id,path,l1_block_number FROM %s ORDER BY block_number", BlocksTableName(blockchain))

	switch side {
	case "first":
		query = fmt.Sprintf("%s LIMIT 1", query)
	case "last":
		query = fmt.Sprintf("%s DESC LIMIT 1", query)
	default:
		return blockIndex, fmt.Errorf("not supported side, choose 'first' or 'last' block")
	}

	queryErr := conn.QueryRow(context.Background(), query).Scan(
		&blockIndex.BlockNumber,
		&blockIndex.BlockHash,
		&blockIndex.BlockTimestamp,
		&blockIndex.ParentHash,
		&blockIndex.RowID,
		&blockIndex.Path,
		&blockIndex.L1BlockNumber,
	)
	if queryErr != nil {
		return blockIndex, queryErr
	}

	blockIndex.chain = blockchain

	return blockIndex, nil
}

func (p *PostgreSQLpgx) GetLatestDBBlockNumber(blockchain string) (uint64, error) {

	pool := p.GetPool()

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return 0, err

	}

	defer conn.Release()

	var blockNumber uint64

	blocksTableName := BlocksTableName(blockchain)
	query := fmt.Sprintf("SELECT block_number FROM %s ORDER BY block_number DESC LIMIT 1", blocksTableName)

	err = conn.QueryRow(context.Background(), query).Scan(&blockNumber)
	if err != nil {
		log.Printf("No data found in %s table", blocksTableName)
		return 0, err
	}

	return blockNumber, nil

}

func (p *PostgreSQLpgx) ReadABIJobs(blockchain string) ([]AbiJob, error) {
	pool := p.GetPool()

	conn, err := pool.Acquire(context.Background())

	if err != nil {
		return nil, err
	}

	defer conn.Release()

	rows, err := conn.Query(context.Background(), "SELECT id, address, user_id, customer_id, abi_selector, chain, abi_name, status, historical_crawl_status, progress, moonworm_task_pickedup, abi, (abi::jsonb)->>'type' as abiType, created_at, updated_at FROM abi_jobs where chain=$1 ", blockchain)

	if err != nil {
		return nil, err
	}

	abiJobs, err := pgx.CollectRows(rows, pgx.RowToStructByName[AbiJob])
	if err != nil {
		return nil, err
	}

	// Check if we have at least one job before accessing
	if len(abiJobs) == 0 {
		return nil, nil // or return an appropriate error if this is considered an error state
	}

	log.Println("Parsed abiJobs:", len(abiJobs), "for blockchain:", blockchain)
	// If you need to process or log the first ABI job separately, do it here

	return abiJobs, nil
}

func (p *PostgreSQLpgx) GetCustomersIDs(blockchain string) ([]string, error) {
	pool := p.GetPool()

	conn, err := pool.Acquire(context.Background())

	if err != nil {
		return nil, err
	}

	defer conn.Release()

	rows, err := conn.Query(context.Background(), "SELECT DISTINCT customer_id FROM abi_jobs where customer_id is not null and blockchain=$1", blockchain)

	if err != nil {
		return nil, err
	}

	var customerIds []string

	for rows.Next() {

		var customerId string

		err = rows.Scan(&customerId)

		if err != nil {
			return nil, err
		}

		customerIds = append(customerIds, customerId)

	}

	return customerIds, nil
}

func (p *PostgreSQLpgx) ReadUpdates(blockchain string, fromBlock uint64, toBlock uint64, customerIds []string) ([]CustomerUpdates, error) {

	pool := p.GetPool()

	conn, err := pool.Acquire(context.Background())

	if err != nil {
		return nil, err
	}

	defer conn.Release()

	transactionsTableName := TransactionsTableName(blockchain)

	logsTableName := LogsTableName(blockchain)

	blocksTableName := BlocksTableName(blockchain)

	query := fmt.Sprintf(`WITH blocks as (
        SELECT
            block_number,
            block_hash,
            block_timestamp
        from
            %s
        WHERE
            block_number >= $1
            and block_number <= $2
    ),
    transactions AS (
        SELECT
            bk.block_number,
            bk.block_hash,
            bk.block_timestamp,
            tx.hash AS transaction_hash,
            tx.to_address AS transaction_address,
            tx.selector AS transaction_selector,
            tx.row_id AS transaction_row_id,
            tx.path AS transaction_path
        FROM
            blocks bk
            LEFT JOIN %s tx ON tx.block_hash = bk.block_hash
    ),
    events AS (
        SELECT
            bk.block_number,
            bk.block_hash,
            bk.block_timestamp,
            logs.transaction_hash AS transaction_hash,
            logs.address AS event_address,
            logs.selector AS event_selector,
            logs.row_id AS event_row_id,
            logs.path AS event_path
        FROM
            blocks bk
            LEFT JOIN %s logs ON logs.block_hash = bk.block_hash
    ),
    jobs AS (
        SELECT
            address as address,
            '0x' || encode(address, 'hex') as address_str,
            customer_id,
            abi_selector,
            abi_name,
            abi,
			(abi)::jsonb ->> 'type' as abi_type,
        	(abi)::jsonb ->> 'stateMutability' as abi_stateMutability
        FROM
            abi_jobs
        WHERE
            chain = $3
    ),
    address_abis AS (
        SELECT
            address_str,
            customer_id,
            json_object_agg(
                abi_selector,
                json_build_object(
                    'abi',
                    '[' || abi || ']',
                    'abi_name',
                    abi_name
                )
            ) AS abis_per_address
        FROM
            jobs
        GROUP BY
            address_str,
            customer_id
    ),
    reformatted_jobs AS (
        SELECT
            customer_id,
            json_object_agg(address_str, abis_per_address) AS abis
        FROM
            address_abis
        GROUP BY
            customer_id
    ),
    abi_transactions AS (
        SELECT
            transactions.block_number,
            transactions.block_timestamp,
            jobs.customer_id,
            jobs.abi_name,
            jobs.address_str,
            transactions.transaction_hash,
            transactions.transaction_address,
            transactions.transaction_selector,
            transactions.transaction_row_id,
            transactions.transaction_path
        FROM
            transactions
            inner JOIN jobs ON abi_type = 'function'
			AND abi_stateMutability != 'view'
			AND transactions.transaction_address = jobs.address
            AND transactions.transaction_selector = jobs.abi_selector
    ),
    abi_events AS (
        SELECT
            events.block_number,
            events.block_timestamp,
            jobs.customer_id,
            jobs.abi_name,
            jobs.address_str,
            events.transaction_hash,
            events.block_hash,
            events.event_address,
            events.event_selector,
            events.event_row_id,
            events.event_path
        FROM
            events
            inner JOIN jobs ON abi_type = 'event' 
			AND events.event_address = jobs.address
            AND events.event_selector = jobs.abi_selector
    ),
    combined AS (
        SELECT
            block_number,
            block_timestamp,
            customer_id,
            'transaction' AS type,
            abi_name,
            address_str,
            transaction_hash AS hash,
            transaction_address AS address,
            transaction_selector AS selector,
            transaction_row_id AS row_id,
            transaction_path AS path
        FROM
            abi_transactions
        UNION
        ALL
        SELECT
            block_number,
            block_timestamp,
            customer_id,
            'event' AS type,
            abi_name,
            address_str,
            transaction_hash AS hash,
            event_address AS address,
            event_selector AS selector,
            event_row_id AS row_id,
            event_path AS path
        FROM
            abi_events
    )
    SELECT
        customer_id,
        (
            SELECT
                abis
            from
                reformatted_jobs
            where
                reformatted_jobs.customer_id = combined.customer_id
        ) as abis,
        json_object_agg(block_number, block_timestamp) AS blocks_cache,
        json_build_object(
            'transactions',
            COALESCE(
                json_agg(
                    json_build_object(
                        'hash',
                        hash,
                        'address',
                        address_str,
                        'selector',
                        selector,
                        'row_id',
                        row_id,
                        'path',
                        path
                    )
                ) FILTER (
                    WHERE
                        type = 'transaction'
                ),
                '[]'
            ),
            'events',
            COALESCE(
                json_agg(
                    CASE
                        WHEN type = 'event' THEN json_build_object(
                            'hash',
                            hash,
                            'address',
                            address_str,
                            'selector',
                            selector,
                            'row_id',
                            row_id,
                            'path',
                            path
                        )
                    END
                ) FILTER (
                    WHERE
                        type = 'event'
                ),
                '[]'
            )
        ) AS data
    FROM
        combined
    GROUP BY
        customer_id`, blocksTableName, transactionsTableName, logsTableName)

	rows, err := conn.Query(context.Background(), query, fromBlock, toBlock, blockchain)

	if err != nil {
		log.Println("Error querying abi jobs from database", err)
		return nil, err
	}

	var result []CustomerUpdates

	for rows.Next() {
		var customerId string
		var abisJSON, blocksCacheJSON, dataJSON []byte

		// Scan the current row's columns into the variables
		err = rows.Scan(&customerId, &abisJSON, &blocksCacheJSON, &dataJSON)

		var abis map[string]map[string]map[string]string
		if err := json.Unmarshal(abisJSON, &abis); err != nil {
			log.Println("Error unmarshalling abis:", err)
			continue
		}

		var blocksCache map[string]uint64
		if err := json.Unmarshal(blocksCacheJSON, &blocksCache); err != nil {
			log.Println("Error unmarshalling blocks cache:", err)
			continue
		}

		var data RawChainData
		if err := json.Unmarshal(dataJSON, &data); err != nil {
			log.Println("Error unmarshalling data:", err)
			continue
		}

		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		transformedBlocksCache := make(map[uint64]uint64)
		for key, value := range blocksCache {
			uintKey, err := strconv.ParseUint(key, 10, 64)
			if err != nil {
				fmt.Println("Error converting key:", err)
				continue
			}
			transformedBlocksCache[uintKey] = value
		}

		// Append the JSON data to the slice
		result = append(result, CustomerUpdates{
			CustomerID:  customerId,
			Abis:        abis,
			BlocksCache: transformedBlocksCache,
			Data:        data,
		})
	}

	return result, nil

}

func (p *PostgreSQLpgx) EnsureCorrectSelectors(blockchain string, WriteToDB bool, outputFilePath string) error {

	pool := p.GetPool()

	conn, err := pool.Acquire(context.Background())

	if err != nil {
		return err
	}

	defer conn.Release()

	// Get all the ABI jobs for the blockchain

	abiJobs, err := p.ReadABIJobs(blockchain)

	if err != nil {
		return err
	}

	log.Println("Found", len(abiJobs), "ABI jobs for blockchain:", blockchain)

	// for each ABI job, check if the selector is correct

	f, err := os.OpenFile(outputFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		log.Println("Error opening file:", err)
		return err
	}

	writer := bufio.NewWriter(f)

	writer.WriteString(fmt.Sprintf("ABI jobs for blockchain: %s runned as WriteToDB: %v recorded at %s\n", blockchain, WriteToDB, time.Now().String()))

	for _, abiJob := range abiJobs {

		// Get the correct selector for the ABI
		abiObj, err := abi.JSON(strings.NewReader("[" + abiJob.Abi + "]"))

		if err != nil {
			log.Println("Error parsing ABI for ABI job:", abiJob.ID, err)
			return err
		}

		var selector string

		if abiJob.AbiType == "event" {
			selector = abiObj.Events[abiJob.AbiName].ID.String()
		} else {
			selectorRaw := abiObj.Methods[abiJob.AbiName].ID
			selector = fmt.Sprintf("0x%x", selectorRaw)
		}

		if err != nil {
			log.Println("Error getting selector for ABI job:", abiJob.ID, err)
			continue
		}

		// Check if the selector is correct

		if abiJob.AbiSelector != selector {

			if WriteToDB {
				// Update the selector in the database

				_, err := conn.Exec(context.Background(), "UPDATE abi_jobs SET abi_selector = $1 WHERE id = $2", selector, abiJob.ID)

				if err != nil {
					log.Println("Error updating selector for ABI job:", abiJob.ID, err)
					continue
				}

				log.Println("Updated selector:", abiJob.AbiSelector, " for ABI job:", abiJob.ID, " to new selector:", selector)

			}

			_, err = writer.WriteString(fmt.Sprintf("ABI job ID: %s, Name: %s, Address: %x, Selector: %s, Correct Selector: %s\n", abiJob.ID, abiJob.AbiName, abiJob.Address, abiJob.AbiSelector, selector))
			if err != nil {
				log.Println("Error writing to file:", err)
				continue
			}

		}

	}
	writer.Flush()

	f.Close()
	return nil
}

func (p *PostgreSQLpgx) WriteLabes(
	blockchain string,
	transactions []TransactionLabel,
	events []EventLabel,
) error {

	pool := p.GetPool()

	conn, err := pool.Acquire(context.Background())

	if err != nil {
		return err
	}

	defer conn.Release()

	tx, err := conn.Begin(context.Background())

	if err != nil {
		return err
	}

	defer func() {
		if err := recover(); err != nil {
			tx.Rollback(context.Background())
			panic(err)
		} else if err != nil {
			tx.Rollback(context.Background())
		} else {
			err = tx.Commit(context.Background())
		}
	}()

	if len(transactions) > 0 {
		err := p.WriteTransactions(tx, blockchain, transactions)
		if err != nil {
			return err
		}
	}

	if len(events) > 0 {
		err := p.WriteEvents(tx, blockchain, events)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PostgreSQLpgx) WriteEvents(tx pgx.Tx, blockchain string, events []EventLabel) error {

	tableName := LabelsTableName(blockchain)
	columns := []string{"id", "label", "transaction_hash", "log_index", "block_number", "block_hash", "block_timestamp", "caller_address", "origin_address", "address", "label_name", "label_type", "label_data"}
	var valuesMap = make(map[string]UnnestInsertValueStruct)

	valuesMap["id"] = UnnestInsertValueStruct{
		Type:   "UUID",
		Values: make([]interface{}, 0),
	}

	valuesMap["label"] = UnnestInsertValueStruct{
		Type:   "TEXT",
		Values: make([]interface{}, 0),
	}

	valuesMap["transaction_hash"] = UnnestInsertValueStruct{
		Type:   "TEXT",
		Values: make([]interface{}, 0),
	}

	valuesMap["log_index"] = UnnestInsertValueStruct{
		Type:   "BIGINT",
		Values: make([]interface{}, 0),
	}

	valuesMap["block_number"] = UnnestInsertValueStruct{
		Type:   "BIGINT",
		Values: make([]interface{}, 0),
	}

	valuesMap["block_hash"] = UnnestInsertValueStruct{
		Type:   "TEXT",
		Values: make([]interface{}, 0),
	}

	valuesMap["block_timestamp"] = UnnestInsertValueStruct{
		Type:   "BIGINT",
		Values: make([]interface{}, 0),
	}

	valuesMap["caller_address"] = UnnestInsertValueStruct{
		Type:   "BYTEA",
		Values: make([]interface{}, 0),
	}

	valuesMap["origin_address"] = UnnestInsertValueStruct{
		Type:   "BYTEA",
		Values: make([]interface{}, 0),
	}

	valuesMap["address"] = UnnestInsertValueStruct{
		Type:   "BYTEA",
		Values: make([]interface{}, 0),
	}

	valuesMap["label_name"] = UnnestInsertValueStruct{
		Type:   "TEXT",
		Values: make([]interface{}, 0),
	}

	valuesMap["label_type"] = UnnestInsertValueStruct{
		Type:   "TEXT",
		Values: make([]interface{}, 0),
	}

	valuesMap["label_data"] = UnnestInsertValueStruct{
		Type:   "jsonb",
		Values: make([]interface{}, 0),
	}

	for _, event := range events {

		id := uuid.New()

		callerAddressBytes, err := decodeAddress(event.CallerAddress)
		if err != nil {
			fmt.Println("Error decoding caller address:", err, event)
			continue
		}

		originAddressBytes, err := decodeAddress(event.OriginAddress)
		if err != nil {
			fmt.Println("Error decoding origin address:", err, event)
			continue
		}

		addressBytes, err := decodeAddress(event.Address)
		if err != nil {
			fmt.Println("Error decoding address:", err, event)
			continue
		}

		updateValues(valuesMap, "id", id)
		updateValues(valuesMap, "label", event.Label)
		updateValues(valuesMap, "transaction_hash", event.TransactionHash)
		updateValues(valuesMap, "log_index", event.LogIndex)
		updateValues(valuesMap, "block_number", event.BlockNumber)
		updateValues(valuesMap, "block_hash", event.BlockHash)
		updateValues(valuesMap, "block_timestamp", event.BlockTimestamp)
		updateValues(valuesMap, "caller_address", callerAddressBytes)
		updateValues(valuesMap, "origin_address", originAddressBytes)
		updateValues(valuesMap, "address", addressBytes)
		updateValues(valuesMap, "label_name", event.LabelName)
		updateValues(valuesMap, "label_type", event.LabelType)
		updateValues(valuesMap, "label_data", event.LabelData)

	}

	ctx := context.Background()

	err := p.executeBatchInsert(tx, ctx, tableName, columns, valuesMap, "ON CONFLICT DO NOTHING")

	if err != nil {
		return err
	}

	log.Printf("Saved %d events records into %s table", len(events), tableName)

	return nil
}

func (p *PostgreSQLpgx) WriteTransactions(tx pgx.Tx, blockchain string, transactions []TransactionLabel) error {
	tableName := LabelsTableName(blockchain)
	columns := []string{"id", "address", "block_number", "block_hash", "caller_address", "label_name", "label_type", "origin_address", "label", "transaction_hash", "label_data", "block_timestamp"}

	var valuesMap = make(map[string]UnnestInsertValueStruct)

	valuesMap["id"] = UnnestInsertValueStruct{
		Type:   "UUID",
		Values: make([]interface{}, 0),
	}

	valuesMap["address"] = UnnestInsertValueStruct{
		Type:   "BYTEA",
		Values: make([]interface{}, 0),
	}

	valuesMap["block_number"] = UnnestInsertValueStruct{
		Type:   "BIGINT",
		Values: make([]interface{}, 0),
	}

	valuesMap["block_hash"] = UnnestInsertValueStruct{
		Type:   "TEXT",
		Values: make([]interface{}, 0),
	}

	valuesMap["caller_address"] = UnnestInsertValueStruct{
		Type:   "BYTEA",
		Values: make([]interface{}, 0),
	}

	valuesMap["label_name"] = UnnestInsertValueStruct{
		Type:   "TEXT",
		Values: make([]interface{}, 0),
	}

	valuesMap["label_type"] = UnnestInsertValueStruct{
		Type:   "TEXT",
		Values: make([]interface{}, 0),
	}

	valuesMap["origin_address"] = UnnestInsertValueStruct{
		Type:   "BYTEA",
		Values: make([]interface{}, 0),
	}

	valuesMap["label"] = UnnestInsertValueStruct{
		Type:   "TEXT",
		Values: make([]interface{}, 0),
	}

	valuesMap["transaction_hash"] = UnnestInsertValueStruct{
		Type:   "TEXT",
		Values: make([]interface{}, 0),
	}

	valuesMap["label_data"] = UnnestInsertValueStruct{
		Type:   "jsonb",
		Values: make([]interface{}, 0),
	}

	valuesMap["block_timestamp"] = UnnestInsertValueStruct{
		Type:   "BIGINT",
		Values: make([]interface{}, 0),
	}

	for _, transaction := range transactions {

		id := uuid.New()

		addressBytes, err := decodeAddress(transaction.Address)
		if err != nil {
			fmt.Println("Error decoding address:", err, transaction)
			continue
		}

		callerAddressBytes, err := decodeAddress(transaction.CallerAddress)
		if err != nil {
			fmt.Println("Error decoding caller address:", err, transaction)
			continue
		}

		originAddressBytes, err := decodeAddress(transaction.OriginAddress)
		if err != nil {
			fmt.Println("Error decoding origin address:", err, transaction)
			continue
		}

		updateValues(valuesMap, "id", id)
		updateValues(valuesMap, "address", addressBytes)
		updateValues(valuesMap, "block_number", transaction.BlockNumber)
		updateValues(valuesMap, "block_hash", transaction.BlockHash)
		updateValues(valuesMap, "caller_address", callerAddressBytes)
		updateValues(valuesMap, "label_name", transaction.LabelName)
		updateValues(valuesMap, "label_type", transaction.LabelType)
		updateValues(valuesMap, "origin_address", originAddressBytes)
		updateValues(valuesMap, "label", transaction.Label)
		updateValues(valuesMap, "transaction_hash", transaction.TransactionHash)
		updateValues(valuesMap, "label_data", transaction.LabelData)
		updateValues(valuesMap, "block_timestamp", transaction.BlockTimestamp)

	}

	ctx := context.Background()

	err := p.executeBatchInsert(tx, ctx, tableName, columns, valuesMap, "ON CONFLICT DO NOTHING")

	if err != nil {
		return err
	}

	log.Printf("Saved %d transactions records into %s table", len(transactions), tableName)

	return nil
}
