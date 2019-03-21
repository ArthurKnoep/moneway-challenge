#!/usr/bin/env bash
cd `git rev-parse --show-toplevel`
go build -o balance ./cmd/balance/main.go
go build -o transaction ./cmd/transaction/main.go
go build -o account ./cmd/account/main.go
