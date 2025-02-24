package server

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/G7DAO/seer/indexer"
)

const SERVER_API_VERSION = "0.0.1"

type Server struct {
	Host          string
	Port          int
	CORSWhitelist map[string]bool
	DbPool        *indexer.PostgreSQLpgx
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

	txs, txsErr := server.DbPool.GetTransactionsV2(blockchain, sourceAddress)
	if txsErr != nil {
		log.Printf("Unable to query rows, err: %v", txsErr)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
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
