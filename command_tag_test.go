package main

import (
	"fmt"
	"os"
)

func ExampleTag_ListTags() {
	db := OpenDBInMemory()
	MigrateDB(db)

	tags := []*Tag{
		{Name: "frontend"},
		{Name: "backend"},
		{Name: "leads"},
	}
	for _, tag := range tags {
		db.Create(tag)
	}

	if err := Run("tag", os.Stdout, os.Stdout, db); err != nil {
		fmt.Println(err)
	}

	if err := db.Close(); err != nil {
		fmt.Println(err)
	}

	// Unordered output:
	// frontend
	// backend
	// leads
}

func ExampleTag_DeleteTag() {
	db := OpenDBInMemory()
	MigrateDB(db)

	db.Create(&Tag{Name: "backend"})

	if err := Run("tag -d backend", os.Stdout, os.Stdout, db); err != nil {
		fmt.Println(err)
	}

	var tags []Tag
	db.Find(&tags)
	for _, tag := range tags {
		fmt.Println(tag.Name)
	}

	if err := db.Close(); err != nil {
		fmt.Println(err)
	}

	// Output:
}

func ExampleTag_Create() {
	db := OpenDBInMemory()
	MigrateDB(db)

	db.Create(&Tag{Name: "backend"})

	if err := Run("tag -c frontend", os.Stdout, os.Stdout, db); err != nil {
		fmt.Println(err)
	}

	var tags []Tag
	db.Find(&tags)
	for _, tag := range tags {
		fmt.Println(tag.Name)
	}

	if err := db.Close(); err != nil {
		fmt.Println(err)
	}

	// Unordered output:
	// frontend
	// backend
}

func ExampleTag_CreateWithUsers() {
	db := OpenDBInMemory()
	MigrateDB(db)

	db.Create(&User{Name: "conner"})
	db.Create(&User{Name: "kevin chan"})
	db.Create(&User{Name: "megan"})
	db.Create(&Tag{Name: "backend"})

	if err := Run("tag -c frontend \"kevin chan\" megan", os.Stdout, os.Stdout, db); err != nil {
		fmt.Println(err)
	}

	var frontendTag Tag
	db.Where("name = 'frontend'").First(&frontendTag)

	var frontendUsers []User
	db.Model(&frontendTag).Association("Users").Find(&frontendUsers)
	for _, user := range frontendUsers {
		fmt.Println(user.Name)
	}

	if err := db.Close(); err != nil {
		fmt.Println(err)
	}

	// Unordered output:
	// kevin chan
	// megan
}

func ExampleTag_ListTaggedUsers() {
	db := OpenDBInMemory()
	MigrateDB(db)

	db.Create(&Tag{Name: "backend", Users: []User{
		{Name: "conner"},
		{Name: "kevin"},
		{Name: "megan"},
	}})

	if err := Run("tag backend", os.Stdout, os.Stdout, db); err != nil {
		fmt.Println(err)
	}

	if err := db.Close(); err != nil {
		fmt.Println(err)
	}

	// Unordered output:
	// conner
	// kevin
	// megan
}

func ExampleTag_AddUsersToTag() {
	db := OpenDBInMemory()
	MigrateDB(db)

	db.Create(&User{Name: "conner"})
	db.Create(&User{Name: "kevin chan"})
	db.Create(&User{Name: "megan"})
	db.Create(&Tag{Name: "frontend"})
	db.Create(&Tag{Name: "backend"})

	if err := Run("tag frontend \"kevin chan\" megan", os.Stdout, os.Stdout, db); err != nil {
		fmt.Println(err)
	}

	var frontendTag Tag
	db.Where("name = 'frontend'").First(&frontendTag)

	var frontendUsers []User
	db.Model(&frontendTag).Association("Users").Find(&frontendUsers)
	for _, user := range frontendUsers {
		fmt.Println(user.Name)
	}

	if err := db.Close(); err != nil {
		fmt.Println(err)
	}

	// Unordered output:
	// kevin chan
	// megan
}

func ExampleTag_RemoveUsersFromTag() {
	db := OpenDBInMemory()
	MigrateDB(db)

	db.Create(&Tag{
		Name: "frontend",
		Users: []User{
			{Name: "kevin chan"},
			{Name: "conner"},
		},
	})

	if err := Run("tag frontend --remove \"kevin chan\"", os.Stdout, os.Stdout, db); err != nil {
		fmt.Println(err)
	}

	var frontendTag Tag
	db.Where("name = 'frontend'").First(&frontendTag)

	var frontendUsers []User
	db.Model(&frontendTag).Association("Users").Find(&frontendUsers)
	for _, user := range frontendUsers {
		fmt.Println(user.Name)
	}

	if err := db.Close(); err != nil {
		fmt.Println(err)
	}

	// Unordered output:
	// conner
}
