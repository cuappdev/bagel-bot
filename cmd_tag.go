package main

import (
	"github.com/alecthomas/kong"
	"gorm.io/gorm"
	"io"
	"strings"
)

type CmdTag struct {
	Create bool `help:"Create the tag if it doesn't already exist" short:"c"`
	Delete bool `help:"Delete the tag" short:"d"`
	Remove bool `help:"Remove users from tag" short:"rm"`

	Tag   string   `arg optional help:"Tag to view or append to"`
	Users []string `arg optional help:"Users to append"`
}

func (cmd *CmdTag) Run(ctx *kong.Context, db *gorm.DB) (err error) {
	switch ctx.Command() {
	case "tag":
		tags := []Tag{}
		db.Find(&tags)

		for _, tag := range tags {
			if _, err = io.WriteString(ctx.Stdout, tag.Name+"\n"); err != nil {
				return err
			}
		}

	case "tag <tag>":
		if cmd.Create && cmd.Delete {
			_, err = io.WriteString(ctx.Stderr, "conflicting options: --create and --delete")
			if err != nil {
				return err
			}
		}

		var tag Tag
		txn := db.Where("name = ?", cmd.Tag)
		if cmd.Create {
			txn.FirstOrCreate(&tag, Tag{Name: cmd.Tag})
		} else {
			txn.First(&tag)
		}

		if tag.ID == 0 {
			_, err = io.WriteString(ctx.Stderr, "no such tag "+cmd.Tag)
			return err
		}

		if cmd.Delete {
			db.Delete(&tag)
			return
		}

		var users []User
		db.Model(&tag).Association("Users").Find(&users)

		for _, user := range users {
			if _, err = io.WriteString(ctx.Stdout, user.Name+"\n"); err != nil {
				return err
			}
		}

	case "tag <tag> <users>":
		if cmd.Delete {
			_, err := io.WriteString(ctx.Stderr, "use `tag -d <tag>` to delete a tag")
			return err
		}

		var tag Tag
		txn := db.Where("name = ?", cmd.Tag)
		if cmd.Create {
			txn.FirstOrCreate(&tag, Tag{Name: cmd.Tag})
		} else {
			txn.First(&tag)
		}

		if tag.ID == 0 {
			_, err = io.WriteString(ctx.Stdout, "no such tag "+cmd.Tag)
			return err
		}

		var users []User
		db.Find(&users)

		nameToUser := map[string]User{}
		for _, user := range users {
			normalized := strings.ToLower(user.Name)
			nameToUser[normalized] = user
		}

		var usersToAddRemove []User
		for _, username := range cmd.Users {
			user, ok := nameToUser[strings.ToLower(username)]
			if !ok {
				if _, err = io.WriteString(ctx.Stderr, "no such user "+username); err != nil {
					return err
				}
			}

			usersToAddRemove = append(usersToAddRemove, user)
		}

		txn2 := db.Model(&tag).Association("Users")
		if cmd.Remove {
			txn2.Delete(usersToAddRemove)
		} else {
			txn2.Append(usersToAddRemove)
		}

	default:
		panic("unrecognized command " + ctx.Command())
	}

	return nil
}
