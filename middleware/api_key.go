package middleware

import "net/http"

// funcction yg akan menerima api key, reuturn function handler
// akan mengembalikan http.handler

func APIKey(validApiKey string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			// before function running

			apiKey := r.Header.Get("X-Api-Key")

			if apiKey == "" {
				http.Error(w, "Api key Required", http.StatusUnauthorized)
				return
			}

			if apiKey != validApiKey {
				http.Error(w, "Invalid API Key", http.StatusUnauthorized)
			}

			next(w, r)

			// after function running

		}
	}
}
