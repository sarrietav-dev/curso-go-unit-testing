package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddSuccess(t *testing.T) {
	c := require.New(t)

	result := Add(2, 3)

	expect := 5

	c.Equal(expect, result)
}
