package main

import (
	"github.com/alecthomas/kong"
	"github.com/jinzhu/gorm"
	"io"
)

type CmdSync struct {
}

func (cmd *CmdSync) Run(ctx *kong.Context, db *gorm.DB, s *Slack) (err error) {
	if err = SyncUsers(db, s); err != nil {
		if _, err = io.WriteString(ctx.Stderr, err.Error()+"\n"); err != nil {
			return err
		}
		return nil
	}
	if _, err = io.WriteString(ctx.Stdout, "sync successful\n"); err != nil {
		return err
	}
	return nil
}
