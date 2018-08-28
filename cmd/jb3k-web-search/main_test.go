package main

import (
	"fmt"
	"testing"
)

func TestSum(t *testing.T) {
	fmt.Println("dsf")
	total := Sum(5, 5)
	if total != 10 {
		t.Errorf("Sum was incorrect, got: %d, want: %d.", total, 10)
	}
}
