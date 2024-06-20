package indexer

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

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

func (p *PostgreSQLpgx) writeBlockIndexToDB(blockchain string, indexes []BlockIndex) error {
	tableName := blockchain + "_blocks"
	pool := p.GetPool()

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		fmt.Println("Connection error", err)
		return err
	}
	defer conn.Release()

	// Start building the bulk insert query
	var query string
	isBlockchainWithL1Chain := IsBlockchainWithL1Chain(blockchain)
	if isBlockchainWithL1Chain {
		query = fmt.Sprintf("INSERT INTO %s (block_number,block_hash,block_timestamp,parent_hash,row_id,path,l1_block_number) VALUES ", tableName)
	} else {
		query = fmt.Sprintf("INSERT INTO %s (block_number,block_hash,block_timestamp,parent_hash,row_id,path) VALUES ", tableName)
	}
	// Placeholder slice for query parameters
	var params []interface{}

	// Loop through indexes to append values and parameters
	var indexesLen int
	for i, index := range indexes {
		if isBlockchainWithL1Chain {
			query += fmt.Sprintf("( $%d, $%d, $%d, $%d, $%d, $%d, $%d),", i*7+1, i*7+2, i*7+3, i*7+4, i*7+5, i*7+6, i*7+7)
			params = append(params, index.BlockNumber, index.BlockHash, index.BlockTimestamp, index.ParentHash, index.RowID, index.Path, index.L1BlockNumber)
		} else {
			query += fmt.Sprintf("( $%d, $%d, $%d, $%d, $%d, $%d),", i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6)
			params = append(params, index.BlockNumber, index.BlockHash, index.BlockTimestamp, index.ParentHash, index.RowID, index.Path)
		}
		indexesLen++
	}

	// Remove the last comma from the query
	query = query[:len(query)-1]

	// Add the ON CONFLICT clause - adjust based on your conflict resolution strategy
	// For example, to do nothing on conflict with the 'id' column
	query += " ON CONFLICT (block_number) DO NOTHING"

	// Execute the query
	_, err = conn.Exec(context.Background(), query, params...)
	if err != nil {
		fmt.Println("Error executing bulk insert", err)
		return err
	}

	log.Printf("Saved %d records into %s table", indexesLen, tableName)
	return nil
}

func (p *PostgreSQLpgx) writeTransactionIndexToDB(tableName string, indexes []TransactionIndex) error {

	pool := p.GetPool()

	conn, err := pool.Acquire(context.Background())

	if err != nil {

		return err
	}

	defer conn.Release()

	// Start building the bulk insert query
	query := fmt.Sprintf("INSERT INTO %s (block_number, block_hash, hash, index, type, from_address, to_address, selector, row_id, path) VALUES ", tableName)

	// Placeholder slice for query parameters
	var params []interface{}

	// Loop through indexes to append values and parameters
	var toAddressBytes, fromAddressBytes []byte

	var indexesLen int
	for i, index := range indexes {
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d),", i*10+1, i*10+2, i*10+3, i*10+4, i*10+5, i*10+6, i*10+7, i*10+8, i*10+9, i*10+10)

		// Decode to_address
		if len(index.ToAddress) < 2 {
			// /x00 is the null byte
			toAddressBytes = []byte{0x00}
		} else {
			toAddressBytes, err = hex.DecodeString(index.ToAddress[2:]) // Remove the '0x' prefix before conversion
			if err != nil {
				fmt.Println("Error decoding to_address:", err, index)
				continue
			}
		}

		// Decode from_address
		if len(index.FromAddress) < 2 {
			fromAddressBytes = []byte{0x00}
		} else {
			fromAddressBytes, err = hex.DecodeString(index.FromAddress[2:])
			if err != nil {
				fmt.Println("Error decoding from_address:", err, index)
				continue
			}
		}

		// Append the parameters for this record
		params = append(params, index.BlockNumber, index.BlockHash, index.TransactionHash, index.TransactionIndex, index.Type, fromAddressBytes, toAddressBytes, index.Selector, index.RowID, index.Path)
		indexesLen++
	}

	// Remove the last comma from the query
	query = query[:len(query)-1]

	// Add the ON CONFLICT clause - adjust based on your conflict resolution strategy

	query += " ON CONFLICT (hash) DO NOTHING"

	// Execute the query

	_, err = conn.Exec(context.Background(), query, params...)

	if err != nil {

		fmt.Println("Error executing bulk insert", err)

		return err

	}

	log.Printf("Saved %d records into %s table", indexesLen, tableName)

	return nil

}

func (p *PostgreSQLpgx) writeLogIndexToDB(tableName string, indexes []LogIndex) error {

	pool := p.GetPool()

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := pool.Acquire(ctx)
	if err != nil {
		fmt.Println("Connection error", err)
		return err
	}
	defer conn.Release()

	// Start a transaction
	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure the transaction is either committed or rolled back
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback(ctx)
			panic(err) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback(ctx) // err is non-nil; don't change it
		} else {
			err = tx.Commit(ctx) // err is nil; if Commit returns error update err
		}
	}()

	// Define the batch size

	var addressBytes []byte

	for i := 0; i < len(indexes); i += InsertBatchSize {
		// Determine the end of the current batch
		end := i + InsertBatchSize
		if end > len(indexes) {
			end = len(indexes)
		}

		// Start building the bulk insert query
		query := fmt.Sprintf("INSERT INTO %s (transaction_hash, block_hash, address, selector, topic1, topic2, row_id, log_index, path) VALUES ", tableName)

		// Placeholder slice for query parameters
		var params []interface{}

		// Loop through indexes to append values and parameters

		var indexesLen int
		for i, index := range indexes[i:end] {

			query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d),", i*9+1, i*9+2, i*9+3, i*9+4, i*9+5, i*9+6, i*9+7, i*9+8, i*9+9)

			if len(index.Address) < 2 {
				// /x00 is the null byte
				addressBytes = []byte{0x00}

			} else {

				addressBytes, err = hex.DecodeString(index.Address[2:]) // Remove the '0x' prefix before conversion

				if err != nil {
					fmt.Println("Error decoding address", err, index)
					continue
				}
			}

			params = append(params, index.TransactionHash, index.BlockHash, addressBytes, index.Selector, index.Topic1, index.Topic2, index.RowID, index.LogIndex, index.Path)
			indexesLen++
		}

		query = query[:len(query)-1]

		// Add the ON CONFLICT clause - adjust based on your conflict resolution strategy

		query += " ON CONFLICT (transaction_hash, log_index) DO NOTHING"

		if _, err := tx.Exec(context.Background(), query, params...); err != nil {
			fmt.Println("Error executing bulk insert", err)
			return fmt.Errorf("error executing bulk insert for batch: %w", err)
		}

		log.Printf("Saved %d records into %s table", indexesLen, tableName)

	}

	if err != nil {

		log.Println("Error writing log index to database", err)

		return err

	}

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

func (p *PostgreSQLpgx) GetListOfFiles(ctx context.Context, blockchain string, object_type string) ([]string, error) {
	pool := p.GetPool()

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	func_of_object_type := func() string {
		if object_type == "transactions" {
			return TransactionsTableName(blockchain)
		} else if object_type == "logs" {
			return LogsTableName(blockchain)
		} else if object_type == "blocks" {
			return BlocksTableName(blockchain)
		} else {
			return ""
		}
	}

	defer conn.Release()

	var files []string

	query := fmt.Sprintf("SELECT path FROM %s group by path order by block_number", func_of_object_type())

	rows, err := conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var path string
		err = rows.Scan(&path)
		if err != nil {
			return nil, err
		}
		files = append(files, path)
	}

	return files, nil
}

func (p *PostgreSQLpgx) GetLatestDBBlockNumber(blockchain string) (uint64, error) {

	pool := p.GetPool()

	conn, err := pool.Acquire(context.Background())

	if err != nil {

		return 0, err

	}

	defer conn.Release()

	var blockNumber uint64

	query := fmt.Sprintf("SELECT block_number FROM %s ORDER BY block_number DESC LIMIT 1", BlocksTableName(blockchain))

	err = conn.QueryRow(context.Background(), query).Scan(&blockNumber)

	if err != nil {

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

	rows, err := conn.Query(context.Background(), "SELECT id, address, user_id, customer_id, abi_selector, chain, abi_name, status, historical_crawl_status, progress, moonworm_task_pickedup, abi, created_at, updated_at FROM abi_jobs where chain=$1 ", blockchain)

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

	log.Println("Parsed abiJobs:", len(abiJobs), " for blockchain:", blockchain)
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
            abi
        FROM
            abi_jobs
        WHERE
            chain = $3 and ((abi)::jsonb ->> 'type' = 'function' or (abi)::jsonb ->> 'type' = 'event')
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
            inner JOIN jobs ON transactions.transaction_address = jobs.address
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
            inner JOIN jobs ON events.event_address = jobs.address
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

func (p *PostgreSQLpgx) WriteEvents(blockchain string, events []EventLabel) error {
	pool := p.GetPool()

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := pool.Acquire(ctx)
	if err != nil {
		fmt.Println("Connection error", err)
		return err
	}
	defer conn.Release()

	// Start a transaction
	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure the transaction is either committed or rolled back
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback(ctx)
			panic(err) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback(ctx) // err is non-nil; don't change it
		} else {
			err = tx.Commit(ctx) // err is nil; if Commit returns error update err
		}
	}()

	tableName := LabelsTableName(blockchain)

	// Too many parameters error
	// Batch insert events calculated as parameters_amount_per_row*batch_size <= 65535 (max number of parameters in a single query)
	for i := 0; i < len(events); i += InsertBatchSize {
		// Determine the end of the current batch
		end := i + InsertBatchSize
		if end > len(events) {
			end = len(events)
		}

		// Start building the bulk insert query
		query := fmt.Sprintf("INSERT INTO %s (id, label, transaction_hash, log_index, block_number, block_hash, block_timestamp, caller_address, origin_address, address, label_name, label_type, label_data) VALUES ", tableName)

		var params []interface{}

		// Loop through labels to append values and parameters
		for j, label := range events[i:end] {
			id := uuid.New()
			query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d),",
				j*13+1, j*13+2, j*13+3, j*13+4, j*13+5, j*13+6, j*13+7, j*13+8, j*13+9, j*13+10, j*13+11, j*13+12, j*13+13)
			params = append(params, id, label.Label, label.TransactionHash, label.LogIndex, label.BlockNumber, label.BlockHash, label.BlockTimestamp, label.CallerAddress, label.OriginAddress, label.Address, label.LabelName, label.LabelType, label.LabelData)
		}

		// Remove the last comma
		query = query[:len(query)-1]

		// ON CONFLICT clause - skip duplicates
		query += " ON CONFLICT (transaction_hash, log_index) where label='seer' and label_type = 'event' DO NOTHING"

		if _, err := tx.Exec(ctx, query, params...); err != nil {
			fmt.Println("Error executing bulk insert", err)
			return fmt.Errorf("error executing bulk insert for batch: %w", err)
		}
	}

	log.Printf("Records %d events inserted into %s", len(events), tableName)

	return nil
}

func (p *PostgreSQLpgx) WriteTransactions(blockchain string, transactions []TransactionLabel) error {
	pool := p.GetPool()

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	// Start a transaction
	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure the transaction is either committed or rolled back
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback(ctx)
			panic(err) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback(ctx) // err is non-nil; don't change it
		} else {
			err = tx.Commit(ctx) // err is nil; if Commit returns error update err
		}
	}()

	tableName := LabelsTableName(blockchain)

	query := fmt.Sprintf("INSERT INTO %s (id, address, block_number, block_hash, caller_address, label_name, label_type, origin_address, label, transaction_hash, label_data, block_timestamp) VALUES ", tableName)

	// Placeholder slice for query parameters
	var params []interface{}

	for i := 0; i < len(transactions); i += InsertBatchSize {
		// Determine the end of the current batch
		end := i + InsertBatchSize
		if end > len(transactions) {
			end = len(transactions)
		}

		// Loop through transactions to append values and parameters
		for row, label := range transactions[i:end] {
			id := uuid.New()
			query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d),",
				row*12+1, row*12+2, row*12+3, row*12+4, row*12+5, row*12+6, row*12+7, row*12+8, row*12+9, row*12+10, row*12+11, row*12+12)
			params = append(params, id, label.Address, label.BlockNumber, label.BlockHash, label.CallerAddress, label.LabelName, label.LabelType, label.OriginAddress, label.Label, label.TransactionHash, label.LabelData, label.BlockTimestamp)
		}

		// Remove the last comma from the query
		query = query[:len(query)-1]

		// Add the ON CONFLICT clause - skip duplicates
		query += " ON CONFLICT (transaction_hash) where label='seer' and label_type = 'tx_call' DO NOTHING"

		// Execute the query
		_, err = conn.Exec(context.Background(), query, params...)
		if err != nil {
			log.Println("Error executing bulk insert", err)
			return err
		}
	}

	log.Printf("Records %d transactions inserted into %s", len(transactions), tableName)
	return nil
}
