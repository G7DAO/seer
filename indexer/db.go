package indexer

import (
	"fmt"
	"net/url"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB is a global variable to hold the GORM database connection.
var DB *gorm.DB

// InitializeDB initializes the database connection using GORM with the given connection URI.
func InitializeDBConnection() error {

	// read the connection string from the environment

	connectionString := os.Getenv("MOONSTREAM_INDEX_URI")

	fmt.Println("connectionString", connectionString)

	uri, err := url.Parse(connectionString)
	if err != nil {
		return fmt.Errorf("parsing connection string: %v", err)
	}
	fmt.Println("uri", uri.Scheme)

	var gormDialector gorm.Dialector
	switch uri.Scheme {
	case "postgresql":
		gormDialector = postgres.Open(connectionString)
	case "sqlite":
		// For SQLite, the connection string is the path after the scheme
		path := connectionString[len("sqlite://"):]
		fmt.Println("sqlite file path", path)
		gormDialector = sqlite.Open(path)
	default:
		fmt.Println("unsupported database type")
		return fmt.Errorf("unsupported database type: %s", uri.Scheme)
	}

	// Open the database connection using GORM
	db, err := gorm.Open(gormDialector, &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return fmt.Errorf("opening database connection: %v", err)
	}

	// Store the GORM DB object in the global variable for later use
	DB = db
	return nil
}

// Example function to add a new BlockIndex
func AddBlockIndex(index BlockIndex) error {
	result := DB.Create(&index) // Pass a pointer of data to Create
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Example function to query BlockIndexes
func GetBlockIndexes() ([]BlockIndex, error) {
	var indexes []BlockIndex
	result := DB.Find(&indexes)
	if result.Error != nil {
		return nil, result.Error
	}
	return indexes, nil
}

// InitializeDB initializes tables in the database
func InitializeDB() error {

	// Auto-migrate block indexes
	// per each chain type

	fmt.Println("Auto-migrating tables")
	// block indexes

	var chains = []string{"ethereum", "polygon"}

	for _, chain := range chains {
		tableName := chain + "_block_index"
		// Manually specify the table name for migration
		DB.Table(tableName).AutoMigrate(&BlockIndex{})
		// add unique constraint
		if DB.Name() != "sqlite" {
			constraintName := fmt.Sprintf("unique_idx_%s_block_number", chain)
			if err := DB.Table(tableName).Exec(fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s UNIQUE (block_number)", tableName, constraintName)).Error; err != nil {
				fmt.Printf("Error adding unique constraint to table %s: %v\n", tableName, err)
				// Handle the error appropriately
			}
		} else {
			fmt.Println("Ignoring block unique constraint for SQLite on chain: ", chain)
		}
	}

	// transaction indexes

	for _, chain := range chains {
		tableName := chain + "_transaction_index"
		// Manually specify the table name for migration
		DB.Table(tableName).AutoMigrate(&TransactionIndex{})
		if DB.Name() != "sqlite" {
			// add unique constraint
			constraintName := fmt.Sprintf("unique_idx_%s_block_number_transaction_index", chain)
			if err := DB.Table(tableName).Exec(fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s UNIQUE (block_number, transaction_hash)", tableName, constraintName)).Error; err != nil {
				fmt.Printf("Error adding unique constraint to table %s: %v\n", tableName, err)
				// Handle the error appropriately
			}
		} else {
			fmt.Println("Ignoring transaction unique constraint for SQLite on chain: ", chain)
		}
	}

	// log indexes

	for _, chain := range chains {
		tableName := chain + "_log_index"
		// Manually specify the table name for migration
		DB.Table(tableName).AutoMigrate(&LogIndex{})

		if DB.Name() != "sqlite" {
			// add unique constraint
			constraintName := fmt.Sprintf("unique_idx_%s_block_number_log_index", chain)
			if err := DB.Table(tableName).Exec(fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s UNIQUE (block_number, transaction_hash, log_index)", tableName, constraintName)).Error; err != nil {
				fmt.Printf("Error adding unique constraint to table %s: %v\n", tableName, err)
				// Handle the error appropriately
			}
		} else {
			fmt.Println("Ignoring log unique constraint for SQLite on chain: ", chain)
		}
	}

	fmt.Println("Tables auto-migrated")

	return nil
}
