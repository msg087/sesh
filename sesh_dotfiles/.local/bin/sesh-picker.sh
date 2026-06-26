#!/usr/bin/env bash

set -uo pipefail

LOGFILE="${HOME}/.local/state/sesh-picker.log"
mkdir -p "$(dirname "$LOGFILE")"

log() {
    printf '[%s] %s\n' "$(date '+%Y-%m-%d %H:%M:%S')" "$*" >>"$LOGFILE"
}

log "started"

if [[ -n "${TMUX:-}" ]]; then
    FZF_CMD="fzf-tmux -p 90%,80%"
else
    FZF_CMD="fzf --height=80%"
fi

log "fzf command: $FZF_CMD"

session="$(
    sesh list --icons --hide-attached --hide-duplicates 2>>"$LOGFILE" | eval "$FZF_CMD \
        --no-sort \
        --ansi \
        --border \
        --border-label ' sesh ' \
        --prompt '⚡  ' \
        --header '  ^a all ^t tmux ^g configs ^x zoxide ^d tmux kill ^f find' \
        --bind 'tab:down,btab:up' \
        --bind 'ctrl-a:change-prompt(⚡  )+reload(sesh list --icons 2>>$LOGFILE)' \
        --bind 'ctrl-t:change-prompt(🪟  )+reload(sesh list -t --icons 2>>$LOGFILE)' \
        --bind 'ctrl-g:change-prompt(⚙️  )+reload(sesh list -c --icons 2>>$LOGFILE)' \
        --bind 'ctrl-x:change-prompt(📁  )+reload(sesh list -z --icons 2>>$LOGFILE)' \
        --bind 'ctrl-f:change-prompt(🔎  )+reload(fd -H -d 2 -t d -E .Trash . ~ 2>>$LOGFILE)' \
        --bind 'ctrl-d:execute-silent(tmux kill-session -t {2..} 2>>$LOGFILE)+reload(sesh list --icons 2>>$LOGFILE)' \
        --preview-window 'right:55%' \
        --preview 'sesh preview {} 2>>$LOGFILE'
    " 2>>"$LOGFILE"
)"
status=$?

log "fzf exited with status: $status"
log "selected session: ${session:-<none>}"

if [[ $status -ne 0 ]]; then
    log "cancelled or fzf failed"
    exit 0
fi

if [[ -z "$session" ]]; then
    log "no session selected"
    exit 0
fi

log "connecting to: $session"

if ! output="$(sesh connect "$session" 2>&1)"; then
    rc=$?
    log "sesh connect failed with status: $rc"
    log "$output"
    exit "$rc"
fi

log "sesh connect succeeded"
exit 0

