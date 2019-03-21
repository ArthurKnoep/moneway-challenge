package main

import "github.com/ArthurKnoep/moneway-challenge/pkg/cmd/balance"

func main() {
	if err := cmd_balance.RunServer(); err != nil {
		panic(err)
	}
}
