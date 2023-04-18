package main

import "testing"

func TestAddSuccess(t *testing.T) {
	result := Add(2, 3)

	expect := 5

	if result != expect {
		t.Errorf("Got %d, expected %d", expect, result)
	}
}
