package main

import "github.com/ArthurKnoep/moneway-challenge/pkg/cmd/transaction"

func main() {
	if err := cmd_transaction.RunServer(); err != nil {
		panic(err)
	}
}
