package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type User struct {
	ID     uint `gorm:"primary_key:true"`
	Name   string
	Tags   []Tag   `gorm:"many2many:user_tags;"`
	Bagels []Bagel `gorm:"many2many:user_bagels;"`
}

type Tag struct {
	ID    uint `gorm:"primary_key:true"`
	Name  string
	Users []User `gorm:"many2many:user_tags;"`
}

type Bagel struct {
	gorm.Model
	Users      []*User `gorm:"many2many:user_bagel;"`
	BagelLogID uint
}

type BagelLog struct {
	gorm.Model
	Bagels []Bagel
	Date   int64
}

func OpenDB(filename string) *gorm.DB {
	db, err := gorm.Open("sqlite3", filename)
	if err != nil {
		panic(err)
	}
	return db
}

func OpenDBInMemory() *gorm.DB {
	return OpenDB("file::memory:?cache=shared")
}

func MigrateDB(db *gorm.DB) {
	db.AutoMigrate(&User{}, &Tag{}, &Bagel{}, &BagelLog{})
}

func DBDump(db *gorm.DB) *gorm.DB {
	users := []User{}
	db.Preload("Tags").Preload("Bagels").Find(&users)
	fmt.Println("Users: (", len(users), ")")
	for _, user := range users {
		fmt.Println("|", user)
	}

	tags := []Tag{}
	db.Preload("Users").Find(&tags)
	fmt.Println("Tags: (", len(tags), ")")
	for _, tag := range tags {
		fmt.Println("|", tag)
	}

	return db
}
