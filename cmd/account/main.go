package main

import "github.com/ArthurKnoep/moneway-challenge/pkg/cmd/account"

func main() {
	if err := cmd_account.RunServer(); err != nil {
		panic(err)
	}
}
