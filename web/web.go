package web

import (
	"log"
	"net/http"
	"os"
	"strings"
	"uc_framework/web/api"
)

func parseBoolEnv(name string) bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv(name)))
	return v == "1" || v == "true" || v == "yes" || v == "on"
}

func buildHTTPSRedirectURL(r *http.Request, httpsAddr string) string {
	host := r.Host
	httpsHost := host
	if strings.TrimSpace(httpsAddr) != "" {
		if strings.HasPrefix(httpsAddr, ":") {
			port := strings.TrimPrefix(httpsAddr, ":")
			if h, _, found := strings.Cut(host, ":"); found {
				httpsHost = h + ":" + port
			} else {
				httpsHost = host + ":" + port
			}
		} else {
			httpsHost = httpsAddr
		}
	}
	return "https://" + httpsHost + r.URL.RequestURI()
}

func withSecurityHeaders(next http.Handler, tlsEnabled bool, enforceHTTPS bool, httpsAddr string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if enforceHTTPS && tlsEnabled && r.TLS == nil {
			http.Redirect(w, r, buildHTTPSRedirectURL(r, httpsAddr), http.StatusPermanentRedirect)
			return
		}

		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; connect-src 'self'; object-src 'none'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'")
		if r.TLS != nil {
			w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		}
		next.ServeHTTP(w, r)
	})
}

func envOrDefault(name, fallback string) string {
	v := strings.TrimSpace(os.Getenv(name))
	if v == "" {
		return fallback
	}
	return v
}

func fileExists(path string) bool {
	if strings.TrimSpace(path) == "" {
		return false
	}
	if st, err := os.Stat(path); err == nil && !st.IsDir() {
		return true
	}
	return false
}

func StartWebServer() {
	mux := http.NewServeMux()

	api.RegisterRoutes(mux)

	// Statische Dateien (HTML, CSS, JS, etc.) aus web/static bereitstellen
	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/", fs)

	httpAddr := envOrDefault("UC_FRAMEWORK_HTTP_ADDR", ":8080")
	httpsAddr := envOrDefault("UC_FRAMEWORK_HTTPS_ADDR", ":8443")
	tlsCertFile := strings.TrimSpace(os.Getenv("UC_FRAMEWORK_TLS_CERT_FILE"))
	tlsKeyFile := strings.TrimSpace(os.Getenv("UC_FRAMEWORK_TLS_KEY_FILE"))
	enforceHTTPS := parseBoolEnv("UC_FRAMEWORK_ENFORCE_HTTPS")
	internetMode := parseBoolEnv("UC_FRAMEWORK_INTERNET_MODE")

	tlsAvailable := fileExists(tlsCertFile) && fileExists(tlsKeyFile)
	handler := withSecurityHeaders(mux, tlsAvailable, enforceHTTPS, httpsAddr)

	if internetMode {
		log.Printf("[SECURITY] [STARTUP] Internet mode enabled - web transport checks active")
		if !tlsAvailable {
			log.Printf("[SECURITY] [CRITICAL] TLS certificate/key missing while UC_FRAMEWORK_INTERNET_MODE=true")
		}
		if !enforceHTTPS {
			log.Printf("[SECURITY] [CRITICAL] UC_FRAMEWORK_ENFORCE_HTTPS is disabled while internet mode is enabled")
		}
	}

	if tlsAvailable {
		go func() {
			log.Printf("Web API & Interface (HTTPS) started on %s", httpsAddr)
			if err := http.ListenAndServeTLS(httpsAddr, tlsCertFile, tlsKeyFile, handler); err != nil {
				log.Printf("HTTPS webserver stopped: %v", err)
			}
		}()
	} else {
		log.Printf("HTTPS disabled: TLS cert/key not configured or not found")
		if enforceHTTPS {
			log.Printf("UC_FRAMEWORK_ENFORCE_HTTPS is enabled but TLS is unavailable; serving HTTP without redirect")
		}
	}

	if enforceHTTPS && tlsAvailable {
		log.Printf("HTTP listener on %s redirects to HTTPS", httpAddr)
	} else {
		log.Printf("Web API & Interface (HTTP) started on %s", httpAddr)
		if internetMode {
			log.Printf("[SECURITY] [WARN] HTTP serves content without HTTPS redirect in internet mode")
		}
	}
	if err := http.ListenAndServe(httpAddr, handler); err != nil {
		log.Printf("Webserver stopped: %v", err)
	}
}
