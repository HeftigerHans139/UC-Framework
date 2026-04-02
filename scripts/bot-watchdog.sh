#!/usr/bin/env bash
set -euo pipefail

SUPERVISOR_SCRIPT=""
BOT_PATH=""
BOT_ARGS=""
WORK_DIR="."
STATE_FILE=""
PID_FILE=""
WATCHDOG_PID_FILE=""
LOG_FILE=""
MIN_INTERVAL=60
MAX_INTERVAL=120

while [[ $# -gt 0 ]]; do
  case "$1" in
    --supervisor-script) SUPERVISOR_SCRIPT="$2"; shift 2 ;;
    --bot-path) BOT_PATH="$2"; shift 2 ;;
    --bot-args) BOT_ARGS="$2"; shift 2 ;;
    --work-dir) WORK_DIR="$2"; shift 2 ;;
    --state-file) STATE_FILE="$2"; shift 2 ;;
    --pid-file) PID_FILE="$2"; shift 2 ;;
    --watchdog-pid-file) WATCHDOG_PID_FILE="$2"; shift 2 ;;
    --log-file) LOG_FILE="$2"; shift 2 ;;
    --min-interval-sec) MIN_INTERVAL="$2"; shift 2 ;;
    --max-interval-sec) MAX_INTERVAL="$2"; shift 2 ;;
    *) shift ;;
  esac
done

if [[ -z "$SUPERVISOR_SCRIPT" || -z "$BOT_PATH" || -z "$STATE_FILE" || -z "$PID_FILE" || -z "$WATCHDOG_PID_FILE" ]]; then
  exit 1
fi

if (( MIN_INTERVAL < 60 )); then MIN_INTERVAL=60; fi
if (( MAX_INTERVAL < MIN_INTERVAL )); then MAX_INTERVAL=120; fi

ensure_parent_dir() {
  local p="$1"
  local d
  d="$(dirname "$p")"
  mkdir -p "$d"
}

log_line() {
  if [[ -n "$LOG_FILE" ]]; then
    ensure_parent_dir "$LOG_FILE"
    printf '%s %s\n' "$(date +%FT%T)" "$1" >> "$LOG_FILE"
  fi
}

read_desired() {
  local desired=0
  if [[ -f "$STATE_FILE" ]]; then
    # shellcheck disable=SC1090
    source "$STATE_FILE" || true
    desired="${DESIRED_RUNNING:-0}"
  fi
  echo "$desired"
}

ensure_parent_dir "$WATCHDOG_PID_FILE"
echo "$$" > "$WATCHDOG_PID_FILE"
log_line "watchdog started pid=$$"

while true; do
  desired="$(read_desired)"

  if [[ "$desired" == "1" ]]; then
    status_json="$(bash "$SUPERVISOR_SCRIPT" --action status --bot-path "$BOT_PATH" --bot-args "$BOT_ARGS" --work-dir "$WORK_DIR" --state-file "$STATE_FILE" --pid-file "$PID_FILE" --log-file "$LOG_FILE" 2>/dev/null || true)"

    if [[ "$status_json" != *'"running":true'* ]]; then
      log_line "bot not running, attempting restart"
      bash "$SUPERVISOR_SCRIPT" --action start --bot-path "$BOT_PATH" --bot-args "$BOT_ARGS" --work-dir "$WORK_DIR" --state-file "$STATE_FILE" --pid-file "$PID_FILE" --log-file "$LOG_FILE" >/dev/null 2>&1 || true
      sleep_time=$(( RANDOM % (MAX_INTERVAL - MIN_INTERVAL + 1) + MIN_INTERVAL ))
      sleep "$sleep_time"
      continue
    fi
  fi

  sleep 10
done
