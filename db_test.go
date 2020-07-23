package main

import (
	"fmt"
)

func ExampleDao_UsersTags() {
	d := NewDaoInMemory()
	db := d.OpenDB()
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
	fmt.Println(alice)
	fmt.Println(bob)
	fmt.Println(carol)
	fmt.Println(david)

	fmt.Println("Tags:")
	fmt.Println(frontend)
	fmt.Println(backend)
	fmt.Println(teamLead)

	var tags []Tag
	db.Model(&alice).Association("Tags").Find(&tags)
	fmt.Println("Before: ", tags)

	alice.Tags = nil
	db.Save(&alice)

	db.Model(&alice).Association("Tags").Find(&tags)
	fmt.Println("After: ", tags)

	_ = db.Close()

	// Output:
	// {1 Alice [] []}
	// {2 Bob [] []}
	// {3 Carol [] []}
	// {4 David [] []}
	// Tags:
	// {1 Frontend [{1 Alice [] []} {2 Bob [] []}]}
	// {2 Backend [{3 Carol [] []} {4 David [] []}]}
	// {3 Team Lead [{1 Alice [] []} {4 David [] []}]}
	// Before:  [{1 Frontend []} {3 Team Lead []}]
	// After:  [{1 Frontend []} {3 Team Lead []}]
}
