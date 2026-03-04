package main

import "net/http"

func (app *application) secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Content-Security-Policy (CSP)
		// This is the most important one!
		// "default-src 'self'" tells the browser:
		// "Only run scripts/styles that come from MY domain. Block everything else."
		// This kills the <script>alert(1)</script> attack.
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; "+
				"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; "+
				"font-src 'self' https://fonts.gstatic.com")

		// 2. Referrer-Policy
		// Controls how much information is included in the Referer header
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")

		// 3. X-Content-Type-Options
		// Prevents the browser from "guessing" the content type (MIME sniffing)
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// 4. X-Frame-Options
		// Prevents your site from being put in an <iframe> (Clickjacking protection)
		w.Header().Set("X-Frame-Options", "deny")

		// 5. X-XSS-Protection
		// (Legacy) Basic XSS filter for older browsers
		w.Header().Set("X-XSS-Protection", "0")

		// Add these headers to prevent browser caching during development
		// TODO: add the following if the environment is Dev
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		next.ServeHTTP(w, r)
	})
}
