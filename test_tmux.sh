#!/bin/bash

SESSION=eval

API_DIR="$HOME/code_wsl/jsr-netsuite/go_jsr_api"
NS_DIR="$HOME/code_wsl/jsr-netsuite/gprs-jsr-api"
IAC_DIR="$HOME/code_wsl/jsr-netsuite/pulumi_iac"

if ! tmux has-session -t "$SESSION" 2>/dev/null; then
    if [[ -n "${TMUX:-}" ]]; then
        TMUX_WIDTH="$(tmux display-message -p '#{client_width}')"
        TMUX_HEIGHT="$(tmux display-message -p '#{client_height}')"
    else
        TMUX_WIDTH="$(tput cols)"
        TMUX_HEIGHT="$(tput lines)"
    fi

    tmux new-session -d \
        -x "$TMUX_WIDTH" \
        -y "$TMUX_HEIGHT" \
        -s "$SESSION" \
        -n api \
        -c "$API_DIR"

    tmux split-window -v -l 12 -t "$SESSION:api.0" -c "$API_DIR"

    tmux new-window -t "$SESSION" -n ns -c "$NS_DIR"
    tmux split-window -v -l 12 -t "$SESSION:ns.0" -c "$NS_DIR"

    tmux new-window -t "$SESSION" -n iac -c "$IAC_DIR"
    tmux split-window -v -l 12 -t "$SESSION:iac.0" -c "$IAC_DIR"

    tmux select-window -t "$SESSION:api"
fi

if [[ -n "${TMUX:-}" ]]; then
    tmux switch-client -t "$SESSION"
else
    tmux attach-session -t "$SESSION"
fi
