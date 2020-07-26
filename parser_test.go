package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestParse_Help(t *testing.T) {
	stdout := strings.Builder{}
	_, _, _ = Parse("bagel --help", &stdout, nil)
	assert.NotEmpty(t, stdout.String())
	t.Log("\n" + stdout.String())

	stdout = strings.Builder{}
	_, _, _ = Parse("bagel tag --help", &stdout, nil)
	assert.NotEmpty(t, stdout.String())
	t.Log("\n" + stdout.String())

	stdout = strings.Builder{}
	_, _, _ = Parse("bagel divvy --help", &stdout, nil)
	assert.NotEmpty(t, stdout.String())
	t.Log("\n" + stdout.String())
}

func TestParse_Tag(t *testing.T) {
	cli, context, err := Parse("bagel tag", nil, nil)
	assert.Nil(t, err, err)
	assert.Equal(t, "tag", context.Command())
	assert.False(t, cli.Tag.Create)
	assert.Equal(t, "", cli.Tag.Tag)
	assert.Equal(t, 0, len(cli.Tag.Users))

	cli, context, err = Parse("bagel tag frontend", nil, nil)
	assert.Nil(t, err, err)
	assert.Equal(t, "tag <tag>", context.Command())
	assert.False(t, cli.Tag.Create)
	assert.Equal(t, "frontend", cli.Tag.Tag)
	assert.Equal(t, 0, len(cli.Tag.Users))

	cli, context, err = Parse("bagel tag -c frontend", nil, nil)
	assert.Nil(t, err, err)
	assert.Equal(t, "tag <tag>", context.Command())
	assert.True(t, cli.Tag.Create)
	assert.Equal(t, "frontend", cli.Tag.Tag)
	assert.Equal(t, 0, len(cli.Tag.Users))

	cli, context, err = Parse("bagel tag backend megan", nil, nil)
	assert.Nil(t, err, err)
	assert.Equal(t, "tag <tag> <users>", context.Command())
	assert.False(t, cli.Tag.Create)
	assert.Equal(t, "backend", cli.Tag.Tag)
	assert.Equal(t, []string{"megan"}, cli.Tag.Users)

	cli, context, err = Parse("bagel tag backend megan \"kevin chan\"", nil, nil)
	assert.Nil(t, err, err)
	assert.Equal(t, "tag <tag> <users>", context.Command())
	assert.False(t, cli.Tag.Create)
	assert.Equal(t, "backend", cli.Tag.Tag)
	assert.Equal(t, []string{"megan", "kevin chan"}, cli.Tag.Users)
}

func TestParse_Divvy(t *testing.T) {
	cli, context, err := Parse("bagel divvy", nil, nil)
	assert.Nil(t, err, err)
	assert.Equal(t, "divvy", context.Command())
	assert.Equal(t, 2, cli.Divvy.Size)

	cli, context, err = Parse("bagel divvy all", nil, nil)
	assert.Nil(t, err, err)
	assert.Equal(t, "divvy <tag>", context.Command())
	assert.Equal(t, 2, cli.Divvy.Size)
	assert.Equal(t, "all", cli.Divvy.Tag)

	cli, context, err = Parse("bagel divvy --size 3 all", nil, nil)
	assert.Nil(t, err, err)
	assert.Equal(t, "divvy <tag>", context.Command())
	assert.Equal(t, 3, cli.Divvy.Size)
}
