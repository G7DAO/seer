package server

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	bugout "github.com/bugout-dev/bugout-go/pkg"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"

	"github.com/G7DAO/seer/indexer"
)

const SERVER_API_VERSION = "0.0.1"

type Server struct {
	Host          string
	Port          int
	CORSWhitelist map[string]bool
	DbPool        *indexer.PostgreSQLpgx
	BugoutClient  *bugout.BugoutClient
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
	SourceAddress     string `json:"source_address"`
	LowestBlockNumber uint64 `json:"lowest_block_number"`

	Nodes []GraphNode  `json:"nodes"`
	Links []GraphLinks `json:"links"`
}

func (server *Server) nowRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	serverTime := time.Now().Format("2006-01-02 15:04:05.999999-07")
	response := NowResponse{ServerTime: serverTime}
	json.NewEncoder(w).Encode(response)
}

func weiToEther(wei *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(params.Ether))
}

func ProcessGraphNodes(txs []indexer.Transaction, graphNodes map[string]uint64, graphLinks map[string]*big.Int) (uint64, []string) {
	var subAddressSls []string

	lowestBlockNum := txs[0].BlockNumber
	for _, tx := range txs {
		// Keep track of the lowest block number
		if tx.BlockNumber < lowestBlockNum {
			lowestBlockNum = tx.BlockNumber
		}

		// Initialize the source node if not already present in set
		if _, exists := graphNodes[tx.FromAddress]; !exists {
			graphNodes[tx.FromAddress] = 0
		}

		graphNodes[tx.FromAddress]++

		// Initialize the target node if not already present and increment subnodes for the source
		if _, exists := graphNodes[tx.ToAddress]; !exists {
			graphNodes[tx.ToAddress] = 0

			subAddressSls = append(subAddressSls, tx.ToAddress)
		}

		// Any bidirectional transfers are calculated under one link (link both from->to and to->from)
		if _, exists := graphLinks[fmt.Sprintf("%s-%s", tx.ToAddress, tx.FromAddress)]; !exists {
			if _, reversExists := graphLinks[fmt.Sprintf("%s-%s", tx.FromAddress, tx.ToAddress)]; !reversExists {
				graphLinks[fmt.Sprintf("%s-%s", tx.FromAddress, tx.ToAddress)] = tx.Value
			}
		}
	}

	return lowestBlockNum, subAddressSls
}

type TransactionsVolumeResponse struct {
	FromAddress    string `json:"from_address"`
	ToAddress      string `json:"to_address"`
	MinBlockNumber uint64 `json:"min_block_number"`
	MaxBlockNumber uint64 `json:"max_block_number"`
	Volume         string `json:"volume"`
	TxsCount       uint64 `json:"txs_count"`
}

func (server *Server) graphsVolumeRoute(w http.ResponseWriter, r *http.Request) {
	fromAddressQe := r.URL.Query().Get("from_address")
	if fromAddressQe == "" {
		http.Error(w, "from_address is required", http.StatusBadRequest)
		return
	}

	if !common.IsHexAddress(fromAddressQe) {
		http.Error(w, "Incorrect address type", http.StatusBadRequest)
		return
	}

	toAddressQe := r.URL.Query().Get("to_address")
	if toAddressQe == "" {
		http.Error(w, "to_address is required", http.StatusBadRequest)
		return
	}

	if !common.IsHexAddress(toAddressQe) {
		http.Error(w, "Incorrect address type", http.StatusBadRequest)
		return
	}

	blockchainQe := r.URL.Query().Get("blockchain")
	if blockchainQe == "" {
		http.Error(w, "Blockchain is required", http.StatusBadRequest)
		return
	}

	var lowestBlockNumQeUint uint64
	lowestBlockNumQe := r.URL.Query().Get("lowest_block_number")
	if lowestBlockNumQe != "" {
		var parseUintErr error
		lowestBlockNumQeUint, parseUintErr = strconv.ParseUint(lowestBlockNumQe, 10, 64)
		if parseUintErr != nil {
			http.Error(w, "lowest_block_number should be an integer", http.StatusBadRequest)
			return
		}
	}

	limitTxs := 1000000
	txsVol, txsErr := server.DbPool.GetTransactionsVolume(blockchainQe, fromAddressQe, toAddressQe, limitTxs, lowestBlockNumQeUint, false)
	if txsErr != nil {
		if txsErr.Error() == "not found" {
			http.Error(w, "No transactions found", http.StatusNotFound)
			return
		}
		log.Printf("Unable to query the row, err: %v", txsErr)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := TransactionsVolumeResponse{
		FromAddress:    fromAddressQe,
		ToAddress:      toAddressQe,
		MinBlockNumber: txsVol.MinBlockNumber,
		MaxBlockNumber: txsVol.MaxBlockNumber,
		Volume:         fmt.Sprintf("%.2f ETH", weiToEther(txsVol.Volume)),
		TxsCount:       txsVol.TxsCount,
	}

	json.NewEncoder(w).Encode(response)
}

func (server *Server) graphsTxsRoute(w http.ResponseWriter, r *http.Request) {
	sourceAddressQe := r.URL.Query().Get("source_address")
	if sourceAddressQe == "" {
		http.Error(w, "source_address is required", http.StatusBadRequest)
		return
	}

	if !common.IsHexAddress(sourceAddressQe) {
		http.Error(w, "Incorrect address type", http.StatusBadRequest)
		return
	}

	blockchainQe := r.URL.Query().Get("blockchain")
	if blockchainQe == "" {
		http.Error(w, "blockchain is required", http.StatusBadRequest)
		return
	}

	var lowestBlockNumQeUint uint64
	lowestBlockNumQe := r.URL.Query().Get("lowest_block_number")
	if lowestBlockNumQe != "" {
		var parseUintErr error
		lowestBlockNumQeUint, parseUintErr = strconv.ParseUint(lowestBlockNumQe, 10, 64)
		if parseUintErr != nil {
			http.Error(w, "lowest_block_number should be an integer", http.StatusBadRequest)
			return
		}
	}

	// TODO: decide, do we want to track where source_address appears in to_address or not

	limitTxs := 100
	// Fetch all unique nodes
	// It does not return all transactions between nodes, but only first
	// Also it gives us lowest block_number for this address, so we do not
	// query transactions for subnodes which were executed this address
	// appeared in blockchain
	txs, txsErr := server.DbPool.GetTransactions(blockchainQe, []string{sourceAddressQe}, limitTxs, lowestBlockNumQeUint, true)
	if txsErr != nil {
		log.Printf("Unable to query rows, err: %v", txsErr)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	graphResponse := GraphResponse{
		SourceAddress:     sourceAddressQe,
		LowestBlockNumber: lowestBlockNumQeUint,
		Nodes:             []GraphNode{},
		Links:             []GraphLinks{},
	}

	if len(txs) == 0 {
		json.NewEncoder(w).Encode(graphResponse)
		return
	}

	// First iteration of parse unique list of txs to nodes graph structure
	// According to DISTINCT ON (to_address) we assume there is only one
	// transaction for one from-to pair
	graphNodes := make(map[string]uint64)
	graphLinks := make(map[string]*big.Int)
	lowestBlockNum, subAddressSls := ProcessGraphNodes(txs, graphNodes, graphLinks)

	graphResponse.LowestBlockNumber = lowestBlockNum

	// Second iteration of parse depth equal 2
	// Query subnodes for source address with txs greater then first tx of source address
	subTxs, subTxsErr := server.DbPool.GetTransactions(blockchainQe, subAddressSls, limitTxs, lowestBlockNum, true)
	if subTxsErr != nil {
		log.Printf("Unable to query rows, err: %v", subTxsErr)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(subTxs) != 0 {
		ProcessGraphNodes(subTxs, graphNodes, graphLinks)
	}

	for node, subNodes := range graphNodes {
		graphResponse.Nodes = append(graphResponse.Nodes, GraphNode{Id: node, SubNodes: subNodes})
	}
	for link, value := range graphLinks {
		linkSls := strings.Split(link, "-")
		valueEth := weiToEther(value)
		graphResponse.Links = append(graphResponse.Links, GraphLinks{Source: linkSls[0], Target: linkSls[1], Value: fmt.Sprintf("%.2f ETH", valueEth)})
	}

	json.NewEncoder(w).Encode(graphResponse)
}

func (server *Server) Run(host string, port int, corsWhitelist map[string]bool) {
	serveMux := http.NewServeMux()
	serveMux.Handle("/graphs/txs", server.accessMiddleware(http.HandlerFunc(server.graphsTxsRoute)))
	serveMux.Handle("/graphs/volume", server.accessMiddleware(http.HandlerFunc(server.graphsVolumeRoute)))
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
