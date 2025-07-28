#!/bin/bash

# Custom Go command wrapper
# Source this file in your ~/.zshrc or ~/.bashrc with: source /path/to/go-aliases.sh

go() {
    if [ "$1" = "run" ] && [ $# -eq 1 ]; then
        # "go run" with no args → "go run main.go"
        echo "Running: go run main.go"
        command go run main.go
    elif [ "$1" = "fresh" ]; then
        # "go fresh" → "go run main.go fresh"
        echo "Running: go run main.go fresh"
        command go run main.go fresh
    else
        # Everything else passes through to normal go command
        command go "$@"
    fi
}

# You can add more custom go commands here
# Example:
# go-build() {
#     go build -o bin/app main.go
# }
