package main

import (
	"os/user"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCommand(t *testing.T) {
	cmd, err := NewCommand(&Theme{}, []string{"location"})
	assert.Nil(t, err)
	assert.Equal(t, "location", cmd.location)

	cmd, err = NewCommand(&Theme{}, []string{"location", "-l"})
	assert.Nil(t, err)
	assert.Equal(t, "location", cmd.location)
	assert.True(t, cmd.LongFormat)

	cmd, err = NewCommand(&Theme{}, []string{})
	assert.Nil(t, err)
	assert.Equal(t, defaultLocation, cmd.location)
}

func TestParseFlagMatched(t *testing.T) {
	cmd := new(Command)
	matched := cmd.parseFlag("-lfFaRth")
	assert.Equal(t, 7, matched)
	assert.True(t, cmd.LongFormat)
	assert.True(t, cmd.NotSort)
	assert.True(t, cmd.Revealing)
	assert.True(t, cmd.ListAll)
	assert.True(t, cmd.Recursion)
	assert.True(t, cmd.TimeSort)
	assert.True(t, cmd.HumanSize)
}

func TestParseFlagNotMatched(t *testing.T) {
	cmd := new(Command)
	matched := cmd.parseFlag("-")
	assert.Equal(t, 0, matched)
	assert.False(t, cmd.LongFormat)
	assert.False(t, cmd.NotSort)
	assert.False(t, cmd.Revealing)
	assert.False(t, cmd.ListAll)
	assert.False(t, cmd.Recursion)
	assert.False(t, cmd.TimeSort)
	assert.False(t, cmd.HumanSize)
}

func TestGetUser(t *testing.T) {
	assert.Equal(t, "", getUser(nil))
	assert.Equal(t, "", getUser(&syscall.Stat_t{Uid: 65536}))
	assert.Equal(t, "root", getUser(&syscall.Stat_t{Uid: 0}))
}

func TestGetGroup(t *testing.T) {
	assert.Equal(t, "", getGroup(nil))
	assert.Equal(t, "", getGroup(&syscall.Stat_t{Gid: 65536}))
	group, err := user.LookupGroupId("0")
	assert.Nil(t, err)
	assert.Equal(t, group.Name, getGroup(&syscall.Stat_t{Gid: 0}))
}

func TestGetNlink(t *testing.T) {
	assert.EqualValues(t, 0, getNlink(nil))
	assert.EqualValues(t, 1, getNlink(&syscall.Stat_t{Nlink: 1}))
}

func TestITermIcon(t *testing.T) {
	expected := " \033]1337;File=inline=1;height=1:icon\a"
	assert.Equal(t, expected, iTermIcon("icon"))
}
