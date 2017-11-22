package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	main()
}

func TestGenerateDelay(t *testing.T) {
	result := generateDelay()
	assert.Equal(t, (100 <= result) && (result <= 1500), true)
}
