package main

import (
	"github.com/jinzhu/gorm"
	"io"
)

func Run(input string, stdout io.Writer, stderr io.Writer, db *gorm.DB) (err error) {
	_, context, err := Parse(input, stdout, stderr)
	if err != nil {
		return err
	}
	if context.Error != nil {
		return context.Error
	}
	context.Bind(db)
	return context.Run()
}
