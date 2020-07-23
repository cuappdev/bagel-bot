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

type Dao struct {
	filename string
}

func NewDao(filename string) Dao {
	return Dao{filename: filename}
}

func NewDaoInMemory() Dao {
	return Dao{filename: "file::memory:?cache=shared"}
}

func (d Dao) OpenDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", d.filename)
	if err != nil {
		panic(err)
	}
	return db
}

func MigrateDB(db *gorm.DB) {
	db.AutoMigrate(&User{}, &Tag{}, &Bagel{}, &BagelLog{})
}

func (d Dao) MigrateDB() {
	log.Info("Migrating database " + d.filename)
	d.withDB(func(db *gorm.DB) {
		MigrateDB(db)
	})
}

func (d Dao) withDB(action func(*gorm.DB)) {
	db := d.OpenDB()

	defer func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()

	action(db)
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
