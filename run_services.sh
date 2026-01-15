#!/usr/bin/env bash

# Ova skripta pokrece servise unutar tmux-a
# Dakle, skines tmux (sudo pacman -S tmux) i obavezno
# das ovoj skripti chmod +x
#
# Nisam sve servise dodao vec samo one koji "rade"
#
# Da ga ugasis: Ctrl+B -> :kill-session

SESSION="Mercypher"

tmux new-session -d -s $SESSION

# --- GATEWAY ---
tmux send-keys -t $SESSION:0 \
    'cd api-gateway && go run cmd/gateway/main.go' C-m

# --- USER SERVICE ---
tmux split-window -v -t $SESSION:0
tmux send-keys -t $SESSION:0.1 \
    'cd user-service && go run cmd/server/main.go' C-m

# --- SESSION SERVICE ---
tmux split-window -h -t $SESSION:0
tmux send-keys -t $SESSION:0.2 \
    'cd session-service && go run cmd/server/main.go' C-m

tmux attach -t $SESSION


