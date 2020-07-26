package main

import (
	"github.com/alecthomas/kong"
	"github.com/jinzhu/gorm"
	"io"
)

type CmdSync struct {
	Data []string `arg help:"Sync the specified data" enum:"users,bagels"`
}

func (cmd *CmdSync) Run(ctx *kong.Context, db *gorm.DB, s *Slack) (err error) {
	for _, whatToSync := range cmd.Data {
		switch whatToSync {
		case "users":
			if err = SyncUsers(db, s); err != nil {
				if _, err = io.WriteString(ctx.Stderr, err.Error()+"\n"); err != nil {
					return err
				}
				return nil
			}
			if _, err = io.WriteString(ctx.Stdout, "user sync successful\n"); err != nil {
				return err
			}

		case "bagels":
			if err = SyncBagels(db, s); err != nil {
				if _, err = io.WriteString(ctx.Stderr, err.Error()+"\n"); err != nil {
					return err
				}
				return nil
			}
			if _, err = io.WriteString(ctx.Stdout, "bagels sync successful\n"); err != nil {
				return err
			}

		default:
			panic("unexpected sync: " + whatToSync)
		}
	}

	return nil
}
