package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

type BotControlConfig struct {
	Enabled                bool
	SupervisorScript       string
	WatchdogScript         string
	BotExecutable          string
	WorkingDir             string
	BotArgs                []string
	StateFile              string
	PidFile                string
	WatchdogPidFile        string
	LogFile                string
	WatchdogMinIntervalSec int
	WatchdogMaxIntervalSec int
}

var botControlCfg BotControlConfig
var (
	BotStatusProvider      func() (map[string]interface{}, error)
	BotActionExecutor      func(action string) (map[string]interface{}, error)
	WatchdogStatusProvider func() (map[string]interface{}, error)
	WatchdogActionExecutor func(action string) (map[string]interface{}, error)
)

func ConfigureBotControl(cfg BotControlConfig) {
	botControlCfg = cfg
	if botControlCfg.WorkingDir == "" {
		botControlCfg.WorkingDir = "."
	}
	if botControlCfg.WatchdogMinIntervalSec < 60 {
		botControlCfg.WatchdogMinIntervalSec = 60
	}
	if botControlCfg.WatchdogMaxIntervalSec < botControlCfg.WatchdogMinIntervalSec {
		botControlCfg.WatchdogMaxIntervalSec = 120
	}
}

type botActionRequest struct {
	Action string `json:"action"`
}

type botAPIResponse struct {
	OK      bool                   `json:"ok"`
	Message string                 `json:"message,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

func BotStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !botControlCfg.Enabled {
		http.Error(w, "bot control disabled", http.StatusServiceUnavailable)
		return
	}

	if BotStatusProvider != nil {
		status, err := BotStatusProvider()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(status)
		return
	}

	status, err := runSupervisorAction("status")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(status)
}

func BotActionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !botControlCfg.Enabled {
		http.Error(w, "bot control disabled", http.StatusServiceUnavailable)
		return
	}

	var req botActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	action := strings.ToLower(strings.TrimSpace(req.Action))
	if action != "start" && action != "stop" && action != "restart" {
		http.Error(w, "invalid action", http.StatusBadRequest)
		return
	}

	if BotActionExecutor != nil {
		status, err := BotActionExecutor(action)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(status)
		return
	}

	status, err := runSupervisorAction(action)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(status)
}

func BotWatchdogStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !botControlCfg.Enabled {
		http.Error(w, "bot control disabled", http.StatusServiceUnavailable)
		return
	}

	if WatchdogStatusProvider != nil {
		status, err := WatchdogStatusProvider()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(status)
		return
	}

	running, pid := isWatchdogRunning()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":      true,
		"running": running,
		"pid":     pid,
	})
}

func BotWatchdogActionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !botControlCfg.Enabled {
		http.Error(w, "bot control disabled", http.StatusServiceUnavailable)
		return
	}

	var req botActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	action := strings.ToLower(strings.TrimSpace(req.Action))
	if WatchdogActionExecutor != nil {
		status, err := WatchdogActionExecutor(action)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(status)
		return
	}

	switch action {
	case "start":
		if err := StartBotWatchdog(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "stop":
		if err := StopBotWatchdog(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "invalid action", http.StatusBadRequest)
		return
	}

	running, pid := isWatchdogRunning()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":      true,
		"running": running,
		"pid":     pid,
	})
}

func BotSystemInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !botControlCfg.Enabled {
		http.Error(w, "bot control disabled", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":                true,
		"os":                runtime.GOOS,
		"arch":              runtime.GOARCH,
		"working_dir":       botControlCfg.WorkingDir,
		"bot_executable":    botControlCfg.BotExecutable,
		"supervisor_script": botControlCfg.SupervisorScript,
		"watchdog_script":   botControlCfg.WatchdogScript,
		"log_file":          botControlCfg.LogFile,
	})
}

func BotLogsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !botControlCfg.Enabled {
		http.Error(w, "bot control disabled", http.StatusServiceUnavailable)
		return
	}

	lines := parseLogLinesParam(r.URL.Query().Get("lines"))
	level := parseLogLevelParam(r.URL.Query().Get("level"))
	tail, err := readLastLinesFiltered(botControlCfg.LogFile, lines, level)
	if err != nil {
		http.Error(w, fmt.Sprintf("log read failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":       true,
		"file":     botControlCfg.LogFile,
		"lines":    lines,
		"level":    level,
		"content":  tail,
		"has_data": strings.TrimSpace(tail) != "",
	})
}

func StartBotWatchdog() error {
	if !botControlCfg.Enabled {
		return fmt.Errorf("bot control disabled")
	}
	running, _ := isWatchdogRunning()
	if running {
		return nil
	}

	cmd := watchdogStartCommand()
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start watchdog: %w", err)
	}
	return nil
}

func StopBotWatchdog() error {
	running, pid := isWatchdogRunning()
	if !running || pid <= 0 {
		_ = os.Remove(botControlCfg.WatchdogPidFile)
		return nil
	}

	if err := killProcess(pid); err != nil {
		return fmt.Errorf("stop watchdog failed: %w", err)
	}
	_ = os.Remove(botControlCfg.WatchdogPidFile)
	return nil
}

func runSupervisorAction(action string) (map[string]interface{}, error) {
	cmd := supervisorActionCommand(action)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("supervisor action failed: %v (%s)", err, strings.TrimSpace(string(out)))
	}

	result := map[string]interface{}{}
	trimmed := strings.TrimSpace(string(out))
	if trimmed == "" {
		return map[string]interface{}{"ok": true, "message": "ok"}, nil
	}
	if err := json.Unmarshal([]byte(trimmed), &result); err != nil {
		return map[string]interface{}{"ok": true, "message": trimmed}, nil
	}
	return result, nil
}

func parseLogLinesParam(v string) int {
	n, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil {
		return 100
	}
	switch n {
	case 25, 100, 250, 500:
		return n
	default:
		return 100
	}
}

func parseLogLevelParam(v string) string {
	level := strings.ToLower(strings.TrimSpace(v))
	switch level {
	case "all", "error", "warn", "info", "debug":
		return level
	default:
		return "all"
	}
}

func lineMatchesLevel(line, level string) bool {
	if level == "all" {
		return true
	}
	u := strings.ToUpper(line)
	switch level {
	case "error":
		return strings.Contains(u, "ERROR") || strings.Contains(u, " ERR ") || strings.HasPrefix(u, "ERR ")
	case "warn":
		return strings.Contains(u, "WARN") || strings.Contains(u, "WARNING")
	case "info":
		return strings.Contains(u, "INFO")
	case "debug":
		return strings.Contains(u, "DEBUG") || strings.Contains(u, " DBG ") || strings.HasPrefix(u, "DBG ")
	default:
		return true
	}
}

func readLastLinesFiltered(path string, maxLines int, level string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("log file path missing")
	}
	if maxLines <= 0 {
		maxLines = 100
	}

	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	bufSize := 1024 * 1024
	scanner.Buffer(make([]byte, 0, 64*1024), bufSize)

	ring := make([]string, maxLines)
	count := 0
	idx := 0
	for scanner.Scan() {
		line := scanner.Text()
		if !lineMatchesLevel(line, level) {
			continue
		}
		ring[idx] = line
		idx = (idx + 1) % maxLines
		if count < maxLines {
			count++
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	if count == 0 {
		return "", nil
	}

	start := 0
	if count == maxLines {
		start = idx
	}
	out := make([]string, 0, count)
	for i := 0; i < count; i++ {
		pos := (start + i) % maxLines
		out = append(out, ring[pos])
	}
	return strings.Join(out, "\n"), nil
}

func isWatchdogRunning() (bool, int) {
	pid, err := readPIDFromFile(botControlCfg.WatchdogPidFile)
	if err != nil || pid <= 0 {
		return false, 0
	}
	if processExists(pid) {
		return true, pid
	}
	_ = os.Remove(botControlCfg.WatchdogPidFile)
	return false, 0
}

func readPIDFromFile(path string) (int, error) {
	if strings.TrimSpace(path) == "" {
		return 0, fmt.Errorf("pid file path missing")
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(b)))
	if err != nil {
		return 0, err
	}
	return pid, nil
}

func processExists(pid int) bool {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("powershell", "-NoProfile", "-Command", fmt.Sprintf("if (Get-Process -Id %d -ErrorAction SilentlyContinue) { '1' } else { '0' }", pid))
		out, err := cmd.CombinedOutput()
		if err != nil {
			return false
		}
		return strings.TrimSpace(string(out)) == "1"
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = proc.Signal(syscall.Signal(0))
	return err == nil
}

func watchdogStartCommand() *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command(
			"powershell",
			"-NoProfile",
			"-ExecutionPolicy", "Bypass",
			"-File", botControlCfg.WatchdogScript,
			"-SupervisorScript", botControlCfg.SupervisorScript,
			"-BotPath", botControlCfg.BotExecutable,
			"-BotArgs", strings.Join(botControlCfg.BotArgs, " "),
			"-WorkDir", botControlCfg.WorkingDir,
			"-StateFile", botControlCfg.StateFile,
			"-PidFile", botControlCfg.PidFile,
			"-WatchdogPidFile", botControlCfg.WatchdogPidFile,
			"-LogFile", botControlCfg.LogFile,
			"-MinIntervalSec", strconv.Itoa(botControlCfg.WatchdogMinIntervalSec),
			"-MaxIntervalSec", strconv.Itoa(botControlCfg.WatchdogMaxIntervalSec),
		)
	}

	return exec.Command(
		"bash",
		botControlCfg.WatchdogScript,
		"--supervisor-script", botControlCfg.SupervisorScript,
		"--bot-path", botControlCfg.BotExecutable,
		"--bot-args", strings.Join(botControlCfg.BotArgs, " "),
		"--work-dir", botControlCfg.WorkingDir,
		"--state-file", botControlCfg.StateFile,
		"--pid-file", botControlCfg.PidFile,
		"--watchdog-pid-file", botControlCfg.WatchdogPidFile,
		"--log-file", botControlCfg.LogFile,
		"--min-interval-sec", strconv.Itoa(botControlCfg.WatchdogMinIntervalSec),
		"--max-interval-sec", strconv.Itoa(botControlCfg.WatchdogMaxIntervalSec),
	)
}

func supervisorActionCommand(action string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		args := []string{
			"-NoProfile",
			"-ExecutionPolicy", "Bypass",
			"-File", botControlCfg.SupervisorScript,
			"-Action", action,
			"-BotPath", botControlCfg.BotExecutable,
			"-BotArgs", strings.Join(botControlCfg.BotArgs, " "),
			"-WorkDir", botControlCfg.WorkingDir,
			"-StateFile", botControlCfg.StateFile,
			"-PidFile", botControlCfg.PidFile,
			"-LogFile", botControlCfg.LogFile,
		}
		return exec.Command("powershell", args...)
	}

	return exec.Command(
		"bash",
		botControlCfg.SupervisorScript,
		"--action", action,
		"--bot-path", botControlCfg.BotExecutable,
		"--bot-args", strings.Join(botControlCfg.BotArgs, " "),
		"--work-dir", botControlCfg.WorkingDir,
		"--state-file", botControlCfg.StateFile,
		"--pid-file", botControlCfg.PidFile,
		"--log-file", botControlCfg.LogFile,
	)
}

func killProcess(pid int) error {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("powershell", "-NoProfile", "-Command", fmt.Sprintf("Stop-Process -Id %d -Force", pid))
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("%v (%s)", err, strings.TrimSpace(string(out)))
		}
		return nil
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	if err := proc.Signal(syscall.SIGTERM); err != nil {
		return err
	}
	return nil
}

func init() {
	cwd, _ := os.Getwd()
	supervisorScript := filepath.Join(cwd, "scripts", "bot-supervisor.sh")
	watchdogScript := filepath.Join(cwd, "scripts", "bot-watchdog.sh")
	if runtime.GOOS == "windows" {
		supervisorScript = filepath.Join(cwd, "scripts", "bot-supervisor.ps1")
		watchdogScript = filepath.Join(cwd, "scripts", "bot-watchdog.ps1")
	}
	botControlCfg = BotControlConfig{
		Enabled:          false,
		SupervisorScript: supervisorScript,
		WatchdogScript:   watchdogScript,
		WorkingDir:       cwd,
	}
}
