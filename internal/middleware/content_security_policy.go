package middleware

import "net/http"

func ContentSecurityPolicy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; script-src 'self' 'unsafe-eval'; style-src 'self' 'unsafe-hashes' 'sha256-SgZQWsfLqFIbXUavZS4pxgi9Pr0JFuIh5pAp0LdrHPU=' 'sha256-faU7yAF8NxuMTNEwVmBz+VcYeIoBQ2EMHW3WaVxCvnk='; img-src 'self' data:; font-src 'self'; connect-src 'self'; media-src 'self'; frame-src 'none'; frame-ancestors 'none'; object-src 'none'; base-uri 'self'; form-action 'self'; upgrade-insecure-requests; ")
		next.ServeHTTP(w, r)
	})
}
