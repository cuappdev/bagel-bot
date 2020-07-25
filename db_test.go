package main

import (
	"fmt"
	"strconv"
)

func ExampleDB_UsersTags() {
	db := OpenDBInMemory()
	MigrateDB(db)

	alice := User{Name: "Alice"}
	bob := User{Name: "Bob"}
	carol := User{Name: "Carol"}
	david := User{Name: "David"}

	db.Create(&alice)
	db.Create(&bob)
	db.Create(&carol)
	db.Create(&david)

	frontend := Tag{Name: "Frontend", Users: []User{alice, bob}}
	backend := Tag{Name: "Backend", Users: []User{carol, david}}
	teamLead := Tag{Name: "Team Lead", Users: []User{alice, david}}

	db.Create(&frontend)
	db.Create(&backend)
	db.Create(&teamLead)

	fmt.Println("Users:")
	fmt.Println(alice.Name)
	fmt.Println(bob.Name)
	fmt.Println(carol.Name)
	fmt.Println(david.Name)

	fmt.Println("Tags:")
	fmt.Println(frontend.Name)
	fmt.Println(backend.Name)
	fmt.Println(teamLead.Name)

	var tags []Tag
	db.Model(&alice).Association("Tags").Find(&tags)
	fmt.Print("Before: ")
	for i, tag := range tags {
		if i < len(tags) - 1 {
			fmt.Print(tag.Name + " ")
		} else {
			fmt.Print(tag.Name)
		}
	}
	fmt.Println()

	alice.Tags = nil
	db.Save(&alice)

	db.Model(&alice).Association("Tags").Find(&tags)
	fmt.Print("After: ")
	for i, tag := range tags {
		if i < len(tags) - 1 {
			fmt.Print(tag.Name + " ")
		} else {
			fmt.Print(tag.Name)
		}
	}
	fmt.Println()

	if err := db.Close(); err != nil {
		fmt.Println(err)
	}

	// Output:
	// Users:
	// Alice
	// Bob
	// Carol
	// David
	// Tags:
	// Frontend
	// Backend
	// Team Lead
	// Before: Frontend Team Lead
	// After: Frontend Team Lead
}

func ExampleDB_BagelLogs() {
	db := OpenDBInMemory()
	MigrateDB(db)

	log := BagelLog{Date: 123456}
	db.Create(&log)
	for i := 0; i < 5; i++ {
		db.Create(&Bagel{
			SlackConversationID: strconv.Itoa(1000 + i),
			BagelLogID: log.ID,
		})
	}

	var bagels []Bagel
	db.Model(&log).Association("Bagels").Find(&bagels)
	for _, bagel := range bagels {
		fmt.Println(bagel.SlackConversationID, bagel.BagelLogID)
	}

	if err := db.Close(); err != nil {
		fmt.Println(err)
	}

	// Output:
	// 1000 1
	// 1001 1
	// 1002 1
	// 1003 1
	// 1004 1
}
