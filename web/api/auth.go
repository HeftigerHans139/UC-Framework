package api

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token              string `json:"token"`
	ExpiresAt          int64  `json:"expires_at"`
	MustChangePassword bool   `json:"must_change_password"`
}

type authModeResponse struct {
	Enabled              bool   `json:"enabled"`
	Provider             string `json:"provider"`
	Mode                 string `json:"mode"`
	RanksystemConfigured bool   `json:"ranksystem_configured"`
	Username             string `json:"username"`
	ForcePasswordChange  bool   `json:"force_password_change"`
}

type changeAuthModeRequest struct {
	Mode string `json:"mode"`
}

type authHealthResponse struct {
	Provider   string `json:"provider"`
	Healthy    bool   `json:"healthy"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

type changePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type changeUsernameRequest struct {
	CurrentPassword string `json:"current_password"`
	NewUsername     string `json:"new_username"`
}

type sessionState struct {
	ExpiresAt          time.Time
	LastSeen           time.Time
	MustChangePassword bool
}

type loginAttemptState struct {
	FirstFailedAt time.Time
	FailedCount   int
	LockedUntil   time.Time
}

type AuthConfig struct {
	Enabled                 bool
	Provider                string
	AllowRanksystemFallback bool
	Username                string
	Password                string
	PasswordHash            string
	ForcePasswordChange     bool
	SessionTTLMinutes       int
	Ranksystem              RanksystemAuthConfig
}

type RanksystemAuthConfig struct {
	LoginURL      string
	UsernameField string
	PasswordField string
	APIKeyHeader  string
	APIKeyValue   string
	BearerToken   string
}

var (
	authEnabled                 = true
	authProvider                = "local"
	authAllowRanksystemFallback = false
	authUsername                = "admin"
	authPassword                = ""
	authPasswordHash            = ""
	authForcePasswordChange     = false
	authHTTPClient              = &http.Client{Timeout: 10 * time.Second}
	ranksystemCfg               RanksystemAuthConfig

	SaveWebAuthPasswordHashFunc            func(passwordHash string) error
	SaveWebAuthForcePasswordChangeFunc     func(force bool) error
	SaveWebAuthUsernameFunc                func(username string) error
	SaveWebAuthEnabledFunc                 func(enabled bool) error
	SaveWebAuthProviderFunc                func(provider string) error
	SaveWebAuthAllowRanksystemFallbackFunc func(enabled bool) error

	tokenStoreMu  sync.Mutex
	tokenStore    = map[string]sessionState{}
	sessionTTL    = 24 * time.Hour
	inactivityTTL = 30 * time.Minute

	loginAttemptMu    sync.Mutex
	loginAttempts     = map[string]*loginAttemptState{}
	loginMaxFails     = 5
	loginFailWindow   = 10 * time.Minute
	loginLockDuration = 15 * time.Minute
	trustProxyHeaders = false
)

func parseBoolEnv(name string) bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv(name)))
	return v == "1" || v == "true" || v == "yes" || v == "on"
}

func ConfigureAuth(cfg AuthConfig) {
	authEnabled = cfg.Enabled
	trustProxyHeaders = parseBoolEnv("UC_FRAMEWORK_TRUST_PROXY")

	provider := strings.TrimSpace(strings.ToLower(cfg.Provider))
	if provider == "" {
		provider = "local"
	}
	authProvider = provider
	authAllowRanksystemFallback = cfg.AllowRanksystemFallback
	if authProvider != "local" {
		authAllowRanksystemFallback = false
	}

	if strings.TrimSpace(cfg.Username) != "" {
		authUsername = strings.TrimSpace(cfg.Username)
	}
	authPassword = cfg.Password
	authPasswordHash = strings.TrimSpace(cfg.PasswordHash)
	authForcePasswordChange = cfg.ForcePasswordChange

	if envHash := strings.TrimSpace(os.Getenv("UC_FRAMEWORK_WEB_AUTH_PASSWORD_HASH")); envHash != "" {
		authPasswordHash = envHash
	}
	if envPass := os.Getenv("UC_FRAMEWORK_WEB_AUTH_PASSWORD"); envPass != "" {
		authPassword = envPass
	}

	if cfg.SessionTTLMinutes > 0 {
		sessionTTL = time.Duration(cfg.SessionTTLMinutes) * time.Minute
	}

	ranksystemCfg = cfg.Ranksystem
	if strings.TrimSpace(ranksystemCfg.UsernameField) == "" {
		ranksystemCfg.UsernameField = "username"
	}
	if strings.TrimSpace(ranksystemCfg.PasswordField) == "" {
		ranksystemCfg.PasswordField = "password"
	}
	if strings.TrimSpace(ranksystemCfg.APIKeyHeader) == "" {
		ranksystemCfg.APIKeyHeader = "X-API-Key"
	}

	if authProvider == "local" && authPasswordHash == "" && authPassword == "" {
		log.Println("[auth] WARN: local auth active, but no password hash/password configured")
	}
}

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !authEnabled {
			next(w, r)
			return
		}
		token := extractBearerToken(r.Header.Get("Authorization"))
		session, ok := getValidSession(token)
		if token == "" || !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if session.MustChangePassword && !isAllowedDuringForcedPasswordChange(r.URL.Path) {
			http.Error(w, "password change required", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("[LOGIN] POST required, got %s", r.Method)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !authEnabled {
		log.Printf("[LOGIN] Auth disabled, returning empty token")
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(loginResponse{Token: "", ExpiresAt: 0})
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	attemptKey := buildLoginAttemptKey(r, req.Username)
	locked, waitFor := isLoginLocked(attemptKey)
	if locked {
		w.Header().Set("Retry-After", fmt.Sprintf("%d", int(waitFor.Seconds())))
		http.Error(w, "too many failed attempts, try again later", http.StatusTooManyRequests)
		return
	}

	if err := authenticate(req.Username, req.Password); err != nil {
		recordLoginFailure(attemptKey)
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	clearLoginFailure(attemptKey)

	token, err := generateToken(32)
	if err != nil {
		http.Error(w, "token generation failed", http.StatusInternalServerError)
		return
	}
	expires := time.Now().Add(sessionTTL)

	tokenStoreMu.Lock()
	tokenStore[token] = sessionState{
		ExpiresAt:          expires,
		LastSeen:           time.Now(),
		MustChangePassword: authProvider == "local" && authForcePasswordChange,
	}
	tokenStoreMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(loginResponse{
		Token:              token,
		ExpiresAt:          expires.Unix(),
		MustChangePassword: authProvider == "local" && authForcePasswordChange,
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !authEnabled {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	token := extractBearerToken(r.Header.Get("Authorization"))
	if token != "" {
		tokenStoreMu.Lock()
		delete(tokenStore, token)
		tokenStoreMu.Unlock()
	}
	w.WriteHeader(http.StatusNoContent)
}

func AuthModeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(authModeResponse{
			Enabled:              authEnabled,
			Provider:             authProvider,
			Mode:                 currentAuthMode(),
			RanksystemConfigured: hasRanksystemLoginConfigured(),
			Username:             authUsername,
			ForcePasswordChange:  authForcePasswordChange,
		})
		return
	case http.MethodPost:
		var req changeAuthModeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		mode := strings.TrimSpace(strings.ToLower(req.Mode))
		enabled := true
		provider := "local"
		allowFallback := false

		switch mode {
		case "none":
			enabled = false
			provider = "local"
		case "local":
			enabled = true
			provider = "local"
		case "ranksystem":
			enabled = true
			provider = "ranksystem"
		case "local_ranksystem":
			enabled = true
			provider = "local"
			allowFallback = true
		default:
			http.Error(w, "invalid mode", http.StatusBadRequest)
			return
		}

		authEnabled = enabled
		authProvider = provider
		authAllowRanksystemFallback = allowFallback

		if SaveWebAuthEnabledFunc != nil {
			if err := SaveWebAuthEnabledFunc(enabled); err != nil {
				http.Error(w, "failed to persist auth mode", http.StatusInternalServerError)
				return
			}
		}
		if SaveWebAuthProviderFunc != nil {
			if err := SaveWebAuthProviderFunc(provider); err != nil {
				http.Error(w, "failed to persist auth mode", http.StatusInternalServerError)
				return
			}
		}
		if SaveWebAuthAllowRanksystemFallbackFunc != nil {
			if err := SaveWebAuthAllowRanksystemFallbackFunc(allowFallback); err != nil {
				http.Error(w, "failed to persist auth mode", http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(authModeResponse{
			Enabled:              authEnabled,
			Provider:             authProvider,
			Mode:                 currentAuthMode(),
			RanksystemConfigured: hasRanksystemLoginConfigured(),
			Username:             authUsername,
			ForcePasswordChange:  authForcePasswordChange,
		})
		return
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func AuthHealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	resp := authHealthResponse{
		Provider: authProvider,
		Healthy:  false,
		Message:  "ranksystem health check unavailable",
	}

	ranksystemModeActive := authProvider == "ranksystem" || (authProvider == "local" && authAllowRanksystemFallback)
	if !ranksystemModeActive {
		resp.Message = "ranksystem login mode is not active"
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	if !hasRanksystemLoginConfigured() {
		resp.Message = "ranksystem login_url missing"
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	req, err := http.NewRequest(http.MethodGet, ranksystemCfg.LoginURL, nil)
	if err != nil {
		resp.Healthy = false
		resp.Message = "invalid login_url"
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	if strings.TrimSpace(ranksystemCfg.APIKeyValue) != "" {
		req.Header.Set(ranksystemCfg.APIKeyHeader, ranksystemCfg.APIKeyValue)
	}
	if strings.TrimSpace(ranksystemCfg.BearerToken) != "" {
		req.Header.Set("Authorization", "Bearer "+ranksystemCfg.BearerToken)
	}

	extResp, err := authHTTPClient.Do(req)
	if err != nil {
		resp.Healthy = false
		resp.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	defer extResp.Body.Close()
	resp.StatusCode = extResp.StatusCode
	if extResp.StatusCode >= 500 {
		resp.Healthy = false
		resp.Message = "ranksystem endpoint returned server error"
	} else {
		resp.Healthy = true
		resp.Message = "ranksystem endpoint reachable"
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func ChangeLocalPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if authProvider != "local" {
		http.Error(w, "password change only available for local provider", http.StatusBadRequest)
		return
	}

	var req changePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if len(req.NewPassword) < 8 {
		http.Error(w, "new password too short (min 8)", http.StatusBadRequest)
		return
	}

	if err := authenticateLocal(authUsername, req.CurrentPassword); err != nil {
		http.Error(w, "current password invalid", http.StatusUnauthorized)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "hash generation failed", http.StatusInternalServerError)
		return
	}

	authPasswordHash = string(hash)
	authPassword = ""
	authForcePasswordChange = false

	if SaveWebAuthPasswordHashFunc != nil {
		if err := SaveWebAuthPasswordHashFunc(authPasswordHash); err != nil {
			http.Error(w, "failed to persist auth config", http.StatusInternalServerError)
			return
		}
	}
	if SaveWebAuthForcePasswordChangeFunc != nil {
		if err := SaveWebAuthForcePasswordChangeFunc(false); err != nil {
			http.Error(w, "failed to persist auth config", http.StatusInternalServerError)
			return
		}
	}

	clearMustChangeFlagForAllSessions()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"ok":      true,
		"message": "password updated",
	})
}

func ChangeLocalUsernameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if authProvider != "local" {
		http.Error(w, "username change only available for local provider", http.StatusBadRequest)
		return
	}

	var req changeUsernameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	newUsername := strings.TrimSpace(req.NewUsername)
	if len(newUsername) < 3 {
		http.Error(w, "new username too short (min 3)", http.StatusBadRequest)
		return
	}
	if len(newUsername) > 64 {
		http.Error(w, "new username too long (max 64)", http.StatusBadRequest)
		return
	}

	if err := authenticateLocal(authUsername, req.CurrentPassword); err != nil {
		http.Error(w, "current password invalid", http.StatusUnauthorized)
		return
	}

	authUsername = newUsername

	if SaveWebAuthUsernameFunc != nil {
		if err := SaveWebAuthUsernameFunc(authUsername); err != nil {
			http.Error(w, "failed to persist auth config", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"ok":       true,
		"message":  "username updated",
		"username": authUsername,
	})
}

func extractBearerToken(authHeader string) string {
	if authHeader == "" {
		return ""
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		return ""
	}
	if strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

func getValidSession(token string) (sessionState, bool) {
	now := time.Now()
	tokenStoreMu.Lock()
	defer tokenStoreMu.Unlock()

	for t, exp := range tokenStore {
		if now.After(exp.ExpiresAt) || now.Sub(exp.LastSeen) > inactivityTTL {
			delete(tokenStore, t)
		}
	}

	session, ok := tokenStore[token]
	if !ok {
		return sessionState{}, false
	}
	if now.After(session.ExpiresAt) || now.Sub(session.LastSeen) > inactivityTTL {
		delete(tokenStore, token)
		return sessionState{}, false
	}
	session.LastSeen = now
	tokenStore[token] = session
	return session, true
}

func isAllowedDuringForcedPasswordChange(path string) bool {
	switch path {
	case "/api/auth/password", "/api/auth/username", "/api/auth/mode", "/api/logout":
		return true
	default:
		return false
	}
}

func clearMustChangeFlagForAllSessions() {
	tokenStoreMu.Lock()
	defer tokenStoreMu.Unlock()
	for token, session := range tokenStore {
		session.MustChangePassword = false
		tokenStore[token] = session
	}
}

func generateToken(size int) (string, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func authenticate(username, password string) error {
	switch authProvider {
	case "local":
		if err := authenticateLocal(username, password); err == nil {
			return nil
		}

		// Optionaler Zweit-Login: Wenn Rank-System konfiguriert ist,
		// darf neben lokal auch gegen das Rank-System authentifiziert werden.
		if authAllowRanksystemFallback && hasRanksystemLoginConfigured() {
			if err := authenticateRanksystem(username, password); err == nil {
				return nil
			}
		}

		return fmt.Errorf("invalid credentials")
	case "ranksystem":
		return authenticateRanksystem(username, password)
	default:
		return fmt.Errorf("unknown auth provider")
	}
}

func hasRanksystemLoginConfigured() bool {
	return strings.TrimSpace(ranksystemCfg.LoginURL) != ""
}

func currentAuthMode() string {
	if !authEnabled {
		return "none"
	}
	if authProvider == "ranksystem" {
		return "ranksystem"
	}
	if authProvider == "local" && authAllowRanksystemFallback {
		return "local_ranksystem"
	}
	return "local"
}

func authenticateLocal(username, password string) error {
	if username != authUsername {
		return fmt.Errorf("invalid credentials")
	}

	if authPasswordHash != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(authPasswordHash), []byte(password)); err == nil {
			return nil
		}
		return fmt.Errorf("invalid credentials")
	}

	if authPassword != "" && password == authPassword {
		return nil
	}

	return fmt.Errorf("invalid credentials")
}

func authenticateRanksystem(username, password string) error {
	if strings.TrimSpace(ranksystemCfg.LoginURL) == "" {
		return fmt.Errorf("ranksystem login_url missing")
	}

	if err := authenticateRanksystemJSON(username, password); err == nil {
		return nil
	}

	if err := authenticateRanksystemForm(username, password); err == nil {
		return nil
	}

	return fmt.Errorf("external auth failed")
}

func authenticateRanksystemJSON(username, password string) error {

	payload := map[string]string{
		ranksystemCfg.UsernameField: username,
		ranksystemCfg.PasswordField: password,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, ranksystemCfg.LoginURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	setRanksystemAuthHeaders(req)

	ok, err := doRanksystemAuthRequest(req)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("external auth failed")
	}
	return nil
}

func authenticateRanksystemForm(username, password string) error {
	form := url.Values{}
	form.Set(ranksystemCfg.UsernameField, username)
	form.Set(ranksystemCfg.PasswordField, password)

	req, err := http.NewRequest(http.MethodPost, ranksystemCfg.LoginURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	setRanksystemAuthHeaders(req)

	ok, err := doRanksystemAuthRequest(req)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("external auth failed")
	}
	return nil
}

func setRanksystemAuthHeaders(req *http.Request) {
	if strings.TrimSpace(ranksystemCfg.APIKeyValue) != "" {
		req.Header.Set(ranksystemCfg.APIKeyHeader, ranksystemCfg.APIKeyValue)
	}
	if strings.TrimSpace(ranksystemCfg.BearerToken) != "" {
		req.Header.Set("Authorization", "Bearer "+ranksystemCfg.BearerToken)
	}
}

func doRanksystemAuthRequest(req *http.Request) (bool, error) {
	client := &http.Client{
		Timeout: authHTTPClient.Timeout,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return true, nil
	}

	if resp.StatusCode >= 300 && resp.StatusCode <= 399 {
		location := strings.ToLower(strings.TrimSpace(resp.Header.Get("Location")))
		if strings.Contains(location, "bot.php") {
			return true, nil
		}
	}

	return false, nil
}

func buildLoginAttemptKey(r *http.Request, username string) string {
	ip := strings.TrimSpace(r.RemoteAddr)
	if trustProxyHeaders {
		if forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); forwarded != "" {
			if comma := strings.Index(forwarded, ","); comma > 0 {
				ip = strings.TrimSpace(forwarded[:comma])
			} else {
				ip = forwarded
			}
		} else if realIP := strings.TrimSpace(r.Header.Get("X-Real-IP")); realIP != "" {
			ip = realIP
		}
	}

	if host, _, err := net.SplitHostPort(ip); err == nil {
		ip = host
	}
	if parsed := net.ParseIP(strings.TrimSpace(ip)); parsed == nil {
		ip = "unknown"
	}

	uname := strings.ToLower(strings.TrimSpace(username))
	if uname == "" {
		uname = "_"
	}
	return ip + "|" + uname
}

func isLoginLocked(key string) (bool, time.Duration) {
	now := time.Now()
	loginAttemptMu.Lock()
	defer loginAttemptMu.Unlock()
	state, ok := loginAttempts[key]
	if !ok {
		return false, 0
	}
	if !state.LockedUntil.IsZero() && now.Before(state.LockedUntil) {
		return true, time.Until(state.LockedUntil)
	}
	if !state.FirstFailedAt.IsZero() && now.Sub(state.FirstFailedAt) > loginFailWindow {
		delete(loginAttempts, key)
	}
	return false, 0
}

func recordLoginFailure(key string) {
	now := time.Now()
	loginAttemptMu.Lock()
	defer loginAttemptMu.Unlock()

	state, ok := loginAttempts[key]
	if !ok {
		loginAttempts[key] = &loginAttemptState{
			FirstFailedAt: now,
			FailedCount:   1,
		}
		return
	}

	if state.FirstFailedAt.IsZero() || now.Sub(state.FirstFailedAt) > loginFailWindow {
		state.FirstFailedAt = now
		state.FailedCount = 1
		state.LockedUntil = time.Time{}
		return
	}

	state.FailedCount++
	if state.FailedCount >= loginMaxFails {
		state.LockedUntil = now.Add(loginLockDuration)
	}
}

func clearLoginFailure(key string) {
	loginAttemptMu.Lock()
	delete(loginAttempts, key)
	loginAttemptMu.Unlock()
}
