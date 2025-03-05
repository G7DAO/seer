package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// CORS middleware
func (server *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var allowedOrigin string
		if server.CORSWhitelist["*"] {
			allowedOrigin = "*"
		} else {
			origin := r.Header.Get("Origin")
			if _, ok := server.CORSWhitelist[origin]; ok {
				allowedOrigin = origin
			}
		}

		if allowedOrigin != "" {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
			// Credentials are cookies, authorization headers, or TLS client certificates
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Log access requests in proper format
func (server *Server) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, readErr := io.ReadAll(r.Body)
		if readErr != nil {
			http.Error(w, "Unable to read body", http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		if len(body) > 0 {
			defer r.Body.Close()
		}

		next.ServeHTTP(w, r)

		var ip string
		var splitErr error
		realIp := r.Header["X-Real-Ip"]
		if len(realIp) == 0 {
			ip, _, splitErr = net.SplitHostPort(r.RemoteAddr)
			if splitErr != nil {
				http.Error(w, fmt.Sprintf("Unable to parse client IP: %s", r.RemoteAddr), http.StatusBadRequest)
				return
			}
		} else {
			ip = realIp[0]
		}

		fullURL := r.URL.Path
		if len(r.URL.RawQuery) > 0 {
			fullURL += "?" + r.URL.RawQuery
		}

		log.Printf("%s %s %s", ip, r.Method, fullURL)
	})
}

// Handle panic errors to prevent server shutdown
func (server *Server) panicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recoverErr := recover(); recoverErr != nil {
				log.Println("recovered", recoverErr)
				http.Error(w, "Internal server error", 500)
			}
		}()

		// There will be a defer with panic handler in each next function
		next.ServeHTTP(w, r)
	})
}

func (server *Server) accessMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Extract auth token
		var authToken string
		authHeaders := r.Header["Authorization"]
		for _, h := range authHeaders {
			if !strings.HasPrefix(h, "Bearer ") {
				http.Error(w, "Invalid authorization token provided", http.StatusForbidden)
				return
			}

			hSlice := strings.Split(h, " ")
			if len(hSlice) != 2 {
				http.Error(w, "Invalid authorization token provided", http.StatusForbidden)
				return
			}

			uuidErr := uuid.Validate(hSlice[1])
			if uuidErr != nil {
				http.Error(w, "Invalid authorization token provided", http.StatusForbidden)
				return
			}

			authToken = hSlice[1]
		}

		if authToken == "" {
			http.Error(w, "No authorization token was passed with the request", http.StatusForbidden)
			return
		}

		// Verify auth token
		userAuth, auErr := server.BugoutClient.Brood.Auth(authToken)
		if auErr != nil {
			if auErr.Error() == "Invalid status code in HTTP response: 404" {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
			log.Printf("Failed to fetch bugout user, err: %v", auErr)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if userAuth.ApplicationId != MOONSTREAM_APPLICATION_ID {
			http.Error(w, "Invalid user token provided. Use a valid Moonstream token", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
