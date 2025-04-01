#!/usr/bin/env bash
echo "Go is required"
if [ "$1" == "dev" ]; then
    go install github.com/a-h/templ/cmd/templ@latest
    go install github.com/air-verse/air@latest
fi
