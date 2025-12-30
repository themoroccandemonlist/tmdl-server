package middleware

import "net/http"

func ContentSecurityPolicy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; script-src 'self'; style-src 'self' https://cdn.jsdelivr.net/npm/daisyui@5.5.14 https://cdn.jsdelivr.net/npm/daisyui@5/themes.css; img-src 'self' data:; font-src 'self'; connect-src 'self'; media-src 'self'; frame-src 'none'; frame-ancestors 'none'; object-src 'none'; base-uri 'self'; form-action 'self'; upgrade-insecure-requests; ")
		next.ServeHTTP(w, r)
	})
}
