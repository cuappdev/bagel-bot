package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"strings"
)

type User struct {
	ID               uint `gorm:"primary_key:true"`
	SlackID          string
	Name             string
	Tags             []Tag   `gorm:"many2many:user_tags;"`
	Bagels           []Bagel `gorm:"many2many:user_bagels;"`
	IsBagelChatsUser bool
	Deleted          bool
	IsBot            bool
}

type Tag struct {
	ID    uint `gorm:"primary_key:true"`
	Name  string
	Users []User `gorm:"many2many:user_tags;"`
}

type Bagel struct {
	gorm.Model
	Users               []User `gorm:"many2many:user_bagels;"`
	SlackConversationID string
	BagelLogID          uint
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

func SyncUsers(db *gorm.DB, s *Slack) (err error) {
	log.Info("Syncing with slack")

	log.Debug("Retrieving users")
	slackUsers, err := s.UsersList()
	if err != nil {
		return err
	}

	for _, slackUser := range slackUsers {
		var dbUser User
		db.Where("slack_id = ?", slackUser.ID).FirstOrCreate(&dbUser, User{})
		dbUser.Name = slackUser.Profile.RealName
		dbUser.SlackID = slackUser.ID
		dbUser.Deleted = slackUser.Deleted
		dbUser.IsBot = slackUser.IsBot
		db.Save(&dbUser)
	}

	log.Debug("Syncing with bagel-chats channel")
	bagelChats, err := findChannel(s, "bagel-testing")
	if err != nil {
		return err
	}
	if bagelChats != nil {
		bagelChatUserIDs, err := s.ConversationsMembers(bagelChats.ID)
		if err != nil {
			return err
		}
		isBagelChatsUser := map[string]bool{}
		for _, userID := range bagelChatUserIDs {
			isBagelChatsUser[userID] = true
		}

		var dbUsers []User
		db.Find(&dbUsers)
		for _, dbUser := range dbUsers {
			dbUser.IsBagelChatsUser = !dbUser.Deleted && !dbUser.IsBot && isBagelChatsUser[dbUser.SlackID]
			db.Save(&dbUser)
		}
	} else {
		log.Warning("Unable to find bagel-chats channel")
	}

	return nil
}

func findChannel(s *Slack, name string) (channel *SlackChannel, err error) {
	channels, err := s.UsersConversations(true, []string{"public_channel"})
	if err != nil {
		return nil, err
	}

	for _, channel := range channels {
		if strings.EqualFold(name, channel.Name) {
			return &channel, nil
		}
	}
	return nil, nil
}
