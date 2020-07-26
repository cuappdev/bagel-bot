package main

import (
	"github.com/jinzhu/gorm"
	"io"
)

type Invocation string

func Run(input string, stdout io.Writer, stderr io.Writer, db *gorm.DB, s *Slack) (err error) {
	_, context, err := Parse(input, stdout, stderr)
	if err != nil {
		if _, err = io.WriteString(stderr, err.Error()+"\n"); err != nil {
			return err
		}
		return nil
	}
	if context == nil {
		return
	}

	invocation := Invocation(input)

	context.Bind(db)
	context.Bind(s)
	context.Bind(&invocation)

	return context.Run()
}
