package main

import (
	"encoding/csv"
	"github.com/alecthomas/kong"
	"io"
	"strings"
)

type CLI struct {
	Tag   CmdTag   `cmd`
	Divvy CmdDivvy `cmd`
	Sync  CmdSync  `cmd`
}

func Parse(input string, stdout io.Writer, stderr io.Writer) (*CLI, *kong.Context, error) {
	cli := CLI{}

	reader := csv.NewReader(strings.NewReader(input))
	reader.Comma = ' '
	fields, err := reader.Read()
	if err != nil {
		return nil, nil, err
	}

	if len(fields) > 0 && fields[0] == "bagel" {
		fields = fields[1:]
	}

	exited := false
	parser, err := kong.New(&cli, kong.Writers(stdout, stderr), kong.Exit(func(i int) {
		exited = true
	}))
	if err != nil {
		return nil, nil, err
	}

	context, err := parser.Parse(fields)
	if exited {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}

	return &cli, context, nil
}
