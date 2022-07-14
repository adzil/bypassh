package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTranslatePaths(t *testing.T) {
	input := []string{
		`-F\\wsl$\Ubuntu\home\john\.ssh\config`,
		`C:\Users\John\.ssh\config`,
		`\\wsl$\Ubuntu\home\doe\.ssh\config,C:\Users\Doe\.ssh\config`,
	}
	expected := []string{
		`-F/home/john/.ssh/config`,
		`/mnt/c/Users/John/.ssh/config`,
		`/home/doe/.ssh/config,/mnt/c/Users/Doe/.ssh/config`,
	}

	actual := translatePaths(input, "Ubuntu")

	assert.Equal(t, expected, actual)
}
