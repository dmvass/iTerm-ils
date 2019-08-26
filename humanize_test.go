package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytes(t *testing.T) {
	tables := []struct {
		actual   int64
		expected string
	}{
		{1, "1B"},
		{512, "512B"},
		{2048, "2.0K"},
		{9874321, "9.9M"},
		{82854982, "83M"},
		{10000000000, "10G"},
		{712893712304234, "713T"},
		{6212893712323224, "6.2P"},
	}
	for _, table := range tables {
		assert.Equal(t, table.expected, Bytes(table.actual), "Actual size %d", table.actual)
	}
}

func TestPermissions(t *testing.T) {
	tables := []struct {
		actual   os.FileMode
		expected string
	}{
		{os.FileMode(0644), "rw-r--r--"},
		{os.FileMode(0464), "r--rw-r--"},
		{os.FileMode(0566), "r-xrw-rw-"},
		{os.FileMode(0600), "rw-------"},
		{os.FileMode(0700), "rwx------"},
		{os.FileMode(0706), "rwx---rw-"},
		{os.FileMode(0622), "rw--w--w-"},
	}
	for _, table := range tables {
		assert.Equal(t, table.expected, Permissions(table.actual), "actual number: %d", table.actual)
	}
}
