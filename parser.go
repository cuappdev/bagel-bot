package main

import (
	"github.com/alecthomas/kong"
	"github.com/mattn/go-shellwords"
	"io"
)

type CLI struct {
	Tag   CmdTag   `cmd help:"Create/delete tags; View/add/remove users from tags;"`
	Divvy CmdDivvy `cmd help:"Divvy the users in the bagel-chats channel or with a specific tag"`
	Sync  CmdSync  `cmd help:"Sync data with slack"`
	Log   CmdLog   `cmd help:"View bagels that were previously made"`
	Msg   CmdMsg   `cmd help:"Read/send messages to groups"`
}

func Parse(input string, stdout io.Writer, stderr io.Writer) (cli *CLI, ctx *kong.Context, err error) {
	cli = &CLI{}
	args, err := shellwords.Parse(input)
	if len(args) > 0 && args[0] == "bagel" {
		args = args[1:]
	}

	exited := false
	parser, err := kong.New(cli, kong.Writers(stdout, stderr), kong.Exit(func(i int) {
		exited = true
	}))
	if err != nil {
		return nil, nil, err
	}

	ctx, err = parser.Parse(args)
	if exited {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}

	return cli, ctx, nil
}
