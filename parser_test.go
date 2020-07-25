package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

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

	stdout := strings.Builder{}
	cli, context, err = Parse("bagel --help tag", &stdout, nil)
	assert.Nil(t, err)
	assert.NotEmpty(t, stdout.String())
}
