package main

import (
	"net/http"
	"os"
	"strings"
)

func (app *application) enableCORS(next http.Handler) http.Handler {
	// 1. Read the environment variable
	trustedOriginsEnv := os.Getenv("CORS_TRUSTED_ORIGINS")

	// 2. Provide a fallback just in case you forget to set it locally
	if trustedOriginsEnv == "" {
		trustedOriginsEnv = "http://localhost:5173,http://localhost:4173"
	}

	// 3. Split the comma-separated string into a Go slice
	allowedOrigins := strings.Split(trustedOriginsEnv, ",")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// 4. Check if the incoming origin is in our allowed list
		for _, allowed := range allowedOrigins {
			// strings.TrimSpace removes any accidental spaces in your env string
			if origin == strings.TrimSpace(allowed) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Add("Vary", "Origin")
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
