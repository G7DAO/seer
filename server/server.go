package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/G7DAO/seer/indexer"
)

const SERVER_API_VERSION = "0.0.1"

type Server struct {
	Host          string
	Port          int
	CORSWhitelist map[string]bool
	DbPool        *pgxpool.Pool
}

type PingResponse struct {
	Status string `json:"status"`
}

type VersionResponse struct {
	Version string `json:"version"`
}

type NowResponse struct {
	ServerTime string `json:"server_time"`
}

func (server *Server) pingRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := PingResponse{Status: "ok"}
	json.NewEncoder(w).Encode(response)
}

func (server *Server) versionRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := VersionResponse{Version: SERVER_API_VERSION}
	json.NewEncoder(w).Encode(response)
}

type Transaction struct {
	BlockNumber uint64   `json:"block_number"`
	FromAddress string   `json:"from_address"`
	ToAddress   string   `json:"to_address"`
	Value       *big.Int `json:"value"`
	Depth       int      `json:"depth"`
}

type GraphNode struct {
	Id       string `json:"id"`
	SubNodes uint64 `json:"sub_nodes"`
}

type GraphLinks struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Value  string `json:"value"`
}

type GraphResponse struct {
	Nodes []GraphNode  `json:"nodes"`
	Links []GraphLinks `json:"links"`
}

func (server *Server) nowRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	serverTime := time.Now().Format("2006-01-02 15:04:05.999999-07")
	response := NowResponse{ServerTime: serverTime}
	json.NewEncoder(w).Encode(response)
}

func (server *Server) graphsV2TxsRoute(w http.ResponseWriter, r *http.Request) {
	sourceAddress := r.URL.Query().Get("source_address")
	if sourceAddress == "" {
		http.Error(w, "Address is required", http.StatusBadRequest)
		return
	}

	if !common.IsHexAddress(sourceAddress) {
		http.Error(w, "Incorrect address type", http.StatusBadRequest)
		return
	}

	blockchain := r.URL.Query().Get("blockchain")
	if blockchain == "" {
		http.Error(w, "Blockchain is required", http.StatusBadRequest)
		return
	}

	txTableName, txTableErr := indexer.TransactionsTableName(blockchain)
	if txTableErr != nil {
		http.Error(w, "Unsupported blockchain provided", http.StatusBadRequest)
		return
	}

	conn, err := server.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	defer conn.Release()

	// rows, err := conn.Query(context.Background(), "SELECT block_number, from_address, to_address, value, 1 AS depth FROM ethereum_transactions WHERE from_address = $1 ORDER BY block_number", sourceAddress)

	query := fmt.Sprintf(`WITH RECURSIVE address_chain AS (
		-- Base case starting with specified 'from_address'
		SELECT
			block_number,
			from_address,
			to_address,
			value,
			1 AS depth
		FROM %s
		WHERE from_address = $1
	
		UNION ALL
		
		-- Recursive case
		SELECT
			e.block_number,
			e.from_address,
			e.to_address,
			e.value,
			ac.depth + 1
		FROM %s e
		JOIN address_chain ac
			ON e.from_address = ac.to_address
		WHERE ac.depth < 3  -- Limit depth up to
	)
	-- Return the entire chain of addresses with depth
	SELECT * FROM address_chain;`, txTableName, txTableName)

	// TODO: Add block_number to lower then first transaction at source_address
	rows, err := conn.Query(context.Background(), query, sourceAddress)

	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var txs []Transaction
	for rows.Next() {
		var tx Transaction
		var valueStr string

		// Scan the values into the struct fields
		err = rows.Scan(&tx.BlockNumber, &tx.FromAddress, &tx.ToAddress, &valueStr, &tx.Depth)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		tx.Value = new(big.Int)
		tx.Value.SetString(valueStr, 10)

		// Append the transaction to the slice
		txs = append(txs, tx)
	}

	// Parse list of data to graph structure
	var graphData GraphResponse

	graphNodes := make(map[string]uint64)
	graphLinks := make(map[string]*big.Int)
	subnodesTracker := make(map[string]map[string]bool) // Tracks unique subnodes per FromAddress

	for _, tx := range txs {
		if _, exists := graphNodes[tx.FromAddress]; !exists {
			graphNodes[tx.FromAddress] = 0
		}
		if _, exists := graphNodes[tx.ToAddress]; !exists {
			graphNodes[tx.ToAddress] = 0
		}

		// Ensure unique subnodes are counted
		if _, exists := subnodesTracker[tx.FromAddress]; !exists {
			subnodesTracker[tx.FromAddress] = make(map[string]bool)
		}
		if !subnodesTracker[tx.FromAddress][tx.ToAddress] { // Only count if this ToAddress is new for FromAddress
			subnodesTracker[tx.FromAddress][tx.ToAddress] = true
			graphNodes[tx.FromAddress]++
		}

		// Any bidirectional transfers are calculated under one link
		if _, exists := graphLinks[fmt.Sprintf("%s-%s", tx.ToAddress, tx.FromAddress)]; !exists {
			graphLinks[fmt.Sprintf("%s-%s", tx.FromAddress, tx.ToAddress)] = tx.Value
		} else {
			fromToAddress := fmt.Sprintf("%s-%s", tx.FromAddress, tx.ToAddress)

			log.Printf("Appended to %s additional %s", fromToAddress, tx.Value.String()) // TODO: debug, delete

			val := new(big.Int).Add(graphLinks[fromToAddress], tx.Value)
			graphLinks[fromToAddress] = val
		}
	}

	for node, subNodes := range graphNodes {
		graphData.Nodes = append(graphData.Nodes, GraphNode{Id: node, SubNodes: subNodes})
	}
	for link, value := range graphLinks {
		linkSls := strings.Split(link, "-")
		graphData.Links = append(graphData.Links, GraphLinks{Source: linkSls[0], Target: linkSls[1], Value: value.String()})
	}

	json.NewEncoder(w).Encode(graphData)
}

func (server *Server) Run(host string, port int, corsWhitelist map[string]bool) {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/graphs/v2/txs", server.graphsV2TxsRoute)
	serveMux.HandleFunc("/now", server.nowRoute)
	serveMux.HandleFunc("/ping", server.pingRoute)
	serveMux.HandleFunc("/version", server.versionRoute)

	// Set list of common middleware, from bottom to top
	commonHandler := server.corsMiddleware(serveMux)
	commonHandler = server.logMiddleware(commonHandler)
	commonHandler = server.panicMiddleware(commonHandler)

	s := http.Server{
		Addr:         fmt.Sprintf("%s:%d", server.Host, server.Port),
		Handler:      commonHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	sErr := s.ListenAndServe()
	if sErr != nil {
		log.Printf("Failed to start server listener, err: %v", sErr)
		os.Exit(1)
	}
}
