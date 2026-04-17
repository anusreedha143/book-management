package main

import "net/http"

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Allow your local Vue dev server to connect.
		// When you deploy Vue to the cloud later, you will change this to your real domain.
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")

		// 2. Allow all necessary HTTP methods for CRUD operations
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// 3. Allow headers that Axios uses to send JSON
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")

		// 4. Handle the browser's invisible "Preflight" security check
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 5. Move on to the actual route handler
		next.ServeHTTP(w, r)
	})
}
