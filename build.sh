#!/bin/bash

GOOS=darwin GOARCH=arm64 go build -o bin/darwin-arm64/localiser
GOOS=darwin GOARCH=amd64 go build -o bin/darwin-amd64/localiser
