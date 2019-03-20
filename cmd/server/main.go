package main

import (
	"github.com/ArthurKnoep/moneway-challenge/pkg/cmd/server"
)

func main() {
	if err := cmd.RunServer(); err != nil {
		panic(err)
	}
}
