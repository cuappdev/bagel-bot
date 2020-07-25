package main

import (
	"encoding/csv"
	"github.com/alecthomas/kong"
	"io"
	"os"
	"strings"
)

type CLI struct {
	Tag CLITagCommand `cmd`
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

	parser, err := kong.New(&cli, kong.Writers(stdout, os.Stderr), kong.Exit(func(i int) {}))
	if err != nil {
		return nil, nil, err
	}

	context, err := parser.Parse(fields)
	if err != nil {
		return nil, nil, err
	}

	return &cli, context, nil
}
