package main

import (
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
)

type User struct {
    gorm.Model
    Name string
    Tags []*Tag `gorm:"many2many:users_tags;"`
    Bagels []*Bagel `gorm:"many2many:users_bagels;"`
}

type Tag struct {
    gorm.Model
    Name string
    Users []*User `gorm:many2many:users_tags;"`
}

type Bagel struct {
    gorm.Model
    Users []*User `gorm:many2many:users_bagel;"`
    BagelLogID uint
}

type BagelLog struct {
    gorm.Model
    Bagels []Bagel
    Date int64
}

func initializeDB() {
    log.Info("Performing DB migration")

    db, err := gorm.Open("sqlite3", "data.db")
    if err != nil {
        log.Critical(err.Error())
        panic("")
    }
    defer db.Close()

    db.AutoMigrate(&User{}, &Tag{}, &Bagel{}, &BagelLog{})
}

