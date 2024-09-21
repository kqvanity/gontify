package main

import (
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {
	t.Run("Case 1", func(t *testing.T) {
		notfs, err := dunstHistory()
		if err != nil {
			t.Error(err)
		}
		fmt.Println(notfs)
	})
}
