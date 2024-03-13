package indexer

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB is a global variable to hold the GORM database connection.

type PostgreSQLpgx struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLpgx() (*PostgreSQLpgx, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)

	defer cancel()

	connect := os.Getenv("MOONSTREAM_INDEX_URI")

	if connect == "" {
		panic("MOONSTREAM_INDEX_URI is not set")
	}

	pool, err := pgxpool.New(ctx, os.Getenv("MOONSTREAM_INDEX_URI"))
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

// write to database
func (p *PostgreSQLpgx) WriteIndex(index interface{}) error {
	// write to database
	return nil
}

// read from database

func (p *PostgreSQLpgx) ReadBlockIndex(startBlock uint64, endBlock uint64) ([]BlockIndex, error) {

	pool := p.GetPool()

	conn, err := pool.Acquire(context.Background())

	if err != nil {
		return nil, err
	}

	defer conn.Release()

	rows, err := conn.Query(context.Background(), "SELECT * FROM products WHERE id = $1", 1)

	blocksIndex, err := pgx.CollectRows(rows, pgx.RowToStructByName[BlockIndex])

	if err != nil {
		return nil, err
	}

	return blocksIndex, nil

}

func (p *PostgreSQLpgx) ReadTransactionIndex(startBlock uint64, endBlock uint64) ([]TransactionIndex, error) {
	// read from database
	return nil, nil
}

func (p *PostgreSQLpgx) ReadLogIndex(startBlock uint64, endBlock uint64) ([]LogIndex, error) {

	// read from database
	return nil, nil
}

func (p *PostgreSQLpgx) writeBlockIndexToDB(tableName string, indexes []BlockIndex) error {
	pool := p.GetPool()

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		fmt.Println("Connection error", err)
		return err
	}
	defer conn.Release()

	// Start building the bulk insert query
	query := fmt.Sprintf("INSERT INTO %s ( block_number, block_hash, block_timestamp, parent_hash, path) VALUES ", tableName)

	// Placeholder slice for query parameters
	var params []interface{}

	// Loop through indexes to append values and parameters
	for i, index := range indexes {

		query += fmt.Sprintf("( $%d, $%d, $%d, $%d, $%d),", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5)
		params = append(params, index.BlockNumber, index.BlockHash, index.BlockTimestamp, index.ParentHash, index.Path)
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

	fmt.Println("Records inserted into", tableName)
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
	query := fmt.Sprintf("INSERT INTO %s (block_number, block_hash,  hash, index, path) VALUES ", tableName)

	// Placeholder slice for query parameters
	var params []interface{}

	// Loop through indexes to append values and parameters

	for i, index := range indexes {

		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d),", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5)
		params = append(params, index.BlockNumber, index.BlockHash, index.TransactionHash, index.TransactionIndex, index.Path)
		//fmt.Println(uuid, index.BlockNumber, index.BlockHash, index.BlockTimestamp, index.TransactionHash, index.TransactionIndex, index.Path)
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

	fmt.Println("Records inserted into", tableName)

	return nil

}

func (p *PostgreSQLpgx) writeLogIndexToDB(tableName string, indexes []LogIndex) error {

	pool := p.GetPool()

	conn, err := pool.Acquire(context.Background())

	if err != nil {

		return err
	}

	defer conn.Release()

	// Start building the bulk insert query
	query := fmt.Sprintf("INSERT INTO %s (address, block_hash, selector, topic1, topic2, transaction_hash, log_index, path) VALUES ", tableName)

	// Placeholder slice for query parameters

	var params []interface{}

	// Loop through indexes to append values and parameters

	for i, index := range indexes {

		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d),", i*8+1, i*8+2, i*8+3, i*8+4, i*8+5, i*8+6, i*8+7, i*8+8)

		params = append(params, index.Address, index.BlockHash, index.Selector, index.Topic1, index.Topic2, index.TransactionHash, index.LogIndex, index.Path)

	}

	// Remove the last comma from the query

	query = query[:len(query)-1]

	// Add the ON CONFLICT clause - adjust based on your conflict resolution strategy

	query += " ON CONFLICT (transaction_hash, log_index) DO NOTHING"

	// Execute the query

	_, err = conn.Exec(context.Background(), query, params...)

	if err != nil {

		fmt.Println("Error executing bulk insert", err)

		return err

	}

	fmt.Println("Records inserted into", tableName)

	return nil

}
