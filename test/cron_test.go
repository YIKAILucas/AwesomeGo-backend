package main

import (
	"fmt"
	"testing"
)

func TestCron(t *testing.T) {
	x, err := getActiveDevice()
	if err != nil {

	}

	fmt.Println(x)
}
