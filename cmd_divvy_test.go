package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func numGroupsPerSize(t *testing.T, divvied [][]User) map[int]int {
	t.Helper()

	numGroupsPerSize := map[int]int{}
	for _, group := range divvied {
		numGroupsPerSize[len(group)]++
	}
	return numGroupsPerSize
}

func TestDivvy_Pair(t *testing.T) {
	users := []User {
		{Name: "1"},
		{Name: "2"},
		{Name: "3"},
		{Name: "4"},
	}
	divvied := pair(users)
	gps := numGroupsPerSize(t, divvied)
	assert.Equal(t, 1, len(gps))
	assert.Equal(t, 2, gps[2])

	users = []User {
		{Name: "1"},
		{Name: "2"},
		{Name: "3"},
		{Name: "4"},
		{Name: "5"},
	}
	divvied = pair(users)
	gps = numGroupsPerSize(t, divvied)
	assert.Equal(t, 2, len(gps))
	assert.Equal(t, 1, gps[2])
	assert.Equal(t, 1, gps[3])
}

func TestDivvy_Divvy(t *testing.T) {
	users := []User {
		{Name: "1"},
		{Name: "2"},
		{Name: "3"},
		{Name: "4"},
		{Name: "5"},
		{Name: "6"},
	}
	divvied := divvy(users, 3)
	gps := numGroupsPerSize(t, divvied)
	assert.Equal(t, 1, len(gps))
	assert.Equal(t, 0, gps[2])
	assert.Equal(t, 2, gps[3])

	users = []User {
		{Name: "1"},
		{Name: "2"},
		{Name: "3"},
		{Name: "4"},
		{Name: "5"},
		{Name: "6"},
		{Name: "7"},
		{Name: "8"},
	}
	divvied = divvy(users, 3)
	gps = numGroupsPerSize(t, divvied)
	assert.Equal(t, 2, len(gps))
	assert.Equal(t, 1, gps[2])
	assert.Equal(t, 2, gps[3])
}
