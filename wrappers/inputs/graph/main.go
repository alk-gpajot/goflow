// Code generated by lib-go/goflow DO NOT EDIT.

// +build !codeanalysis

package main

import (
	"context"
	"fmt"
)

func assert(b bool) {
	if !b {
		panic("assertion failed")
	}
}

func main() {
	var g Inputs

	g = NewInputs()
	g.Run(context.Background(), 1, []int{2, 3}, false)
	fmt.Println("Test done.")

	fmt.Println("All tests done.")
}
