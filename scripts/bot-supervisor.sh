#!/usr/bin/env bash
set -euo pipefail

ACTION=""
BOT_PATH=""
BOT_ARGS=""
WORK_DIR="."
STATE_FILE=""
PID_FILE=""
LOG_FILE=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --action) ACTION="$2"; shift 2 ;;
    --bot-path) BOT_PATH="$2"; shift 2 ;;
    --bot-args) BOT_ARGS="$2"; shift 2 ;;
    --work-dir) WORK_DIR="$2"; shift 2 ;;
    --state-file) STATE_FILE="$2"; shift 2 ;;
    --pid-file) PID_FILE="$2"; shift 2 ;;
    --log-file) LOG_FILE="$2"; shift 2 ;;
    *) shift ;;
  esac
done

if [[ -z "$ACTION" || -z "$BOT_PATH" || -z "$STATE_FILE" || -z "$PID_FILE" ]]; then
  echo '{"ok":false,"message":"missing required arguments"}'
  exit 1
fi

ensure_parent_dir() {
  local p="$1"
  local d
  d="$(dirname "$p")"
  mkdir -p "$d"
}

ensure_parent_dir "$STATE_FILE"
ensure_parent_dir "$PID_FILE"
if [[ -n "$LOG_FILE" ]]; then
  ensure_parent_dir "$LOG_FILE"
fi

log_line() {
  if [[ -n "$LOG_FILE" ]]; then
    printf '%s %s\n' "$(date +%FT%T)" "$1" >> "$LOG_FILE"
  fi
}

read_state() {
  DESIRED_RUNNING=0
  LAST_ACTION=""
  LAST_ERROR=""
  if [[ -f "$STATE_FILE" ]]; then
    # shellcheck disable=SC1090
    source "$STATE_FILE" || true
    DESIRED_RUNNING="${DESIRED_RUNNING:-0}"
    LAST_ACTION="${LAST_ACTION:-}"
    LAST_ERROR="${LAST_ERROR:-}"
  fi
}

write_state() {
  cat > "$STATE_FILE" <<EOF
DESIRED_RUNNING=${DESIRED_RUNNING}
LAST_ACTION="${LAST_ACTION//\"/\"}"
LAST_ERROR="${LAST_ERROR//\"/\"}"
LAST_UPDATE="$(date +%FT%T)"
EOF
}

get_running_pid() {
  if [[ ! -f "$PID_FILE" ]]; then
    return 1
  fi
  local pid
  pid="$(tr -d '[:space:]' < "$PID_FILE")"
  if [[ -z "$pid" ]]; then
    rm -f "$PID_FILE"
    return 1
  fi
  if kill -0 "$pid" 2>/dev/null; then
    echo "$pid"
    return 0
  fi
  rm -f "$PID_FILE"
  return 1
}

resolve_bot_path() {
  if [[ "$BOT_PATH" = /* ]]; then
    echo "$BOT_PATH"
  else
    echo "$WORK_DIR/$BOT_PATH"
  fi
}

json_out() {
  echo "$1"
}

read_state

case "$ACTION" in
  status)
    if pid="$(get_running_pid)"; then
      json_out "{\"ok\":true,\"running\":true,\"pid\":$pid,\"desired_running\":$DESIRED_RUNNING,\"last_action\":\"$LAST_ACTION\",\"last_error\":\"$LAST_ERROR\"}"
    else
      json_out "{\"ok\":true,\"running\":false,\"pid\":0,\"desired_running\":$DESIRED_RUNNING,\"last_action\":\"$LAST_ACTION\",\"last_error\":\"$LAST_ERROR\"}"
    fi
    ;;

  start)
    DESIRED_RUNNING=1
    LAST_ACTION="start"

    if pid="$(get_running_pid)"; then
      write_state
      json_out "{\"ok\":true,\"running\":true,\"pid\":$pid,\"message\":\"already running\",\"desired_running\":true}"
      exit 0
    fi

    bot_bin="$(resolve_bot_path)"
    set +e
    nohup "$bot_bin" $BOT_ARGS >> "${LOG_FILE:-/dev/null}" 2>&1 &
    started_pid=$!
    set -e

    if [[ -z "$started_pid" ]]; then
      LAST_ERROR="start failed"
      write_state
      json_out "{\"ok\":false,\"running\":false,\"pid\":0,\"message\":\"start failed\",\"desired_running\":true}"
      exit 1
    fi

    echo "$started_pid" > "$PID_FILE"
    LAST_ERROR=""
    write_state
    json_out "{\"ok\":true,\"running\":true,\"pid\":$started_pid,\"message\":\"started\",\"desired_running\":true}"
    ;;

  stop)
    DESIRED_RUNNING=0
    LAST_ACTION="stop"

    if pid="$(get_running_pid)"; then
      kill "$pid" 2>/dev/null || true
      sleep 1
      kill -9 "$pid" 2>/dev/null || true
    fi

    rm -f "$PID_FILE"
    LAST_ERROR=""
    write_state
    json_out "{\"ok\":true,\"running\":false,\"pid\":0,\"message\":\"stopped\",\"desired_running\":false}"
    ;;

  restart)
    DESIRED_RUNNING=1
    LAST_ACTION="restart"

    if pid="$(get_running_pid)"; then
      kill "$pid" 2>/dev/null || true
      sleep 1
      kill -9 "$pid" 2>/dev/null || true
    fi

    bot_bin="$(resolve_bot_path)"
    set +e
    nohup "$bot_bin" $BOT_ARGS >> "${LOG_FILE:-/dev/null}" 2>&1 &
    started_pid=$!
    set -e

    if [[ -z "$started_pid" ]]; then
      LAST_ERROR="restart start failed"
      write_state
      json_out "{\"ok\":false,\"running\":false,\"pid\":0,\"message\":\"restart start failed\",\"desired_running\":true}"
      exit 1
    fi

    echo "$started_pid" > "$PID_FILE"
    LAST_ERROR=""
    write_state
    json_out "{\"ok\":true,\"running\":true,\"pid\":$started_pid,\"message\":\"restarted\",\"desired_running\":true}"
    ;;

  *)
    json_out '{"ok":false,"message":"invalid action"}'
    exit 1
    ;;
esac
