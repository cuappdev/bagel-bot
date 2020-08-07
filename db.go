package main

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"strconv"
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
	FeedbackDate        uint
	IsPlanned           bool
	IsCompleted         bool
	FeedbackMsgs        []FeedbackMsg
}

type BagelLog struct {
	gorm.Model
	Bagels     []Bagel
	Date       int64
	Invocation string
}

func BagelLog_Fetch(db *gorm.DB, id string) (log BagelLog, err error) {
	if id == "last" {
		db.Order("date desc").First(&log)
	} else {
		id, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return BagelLog{}, err
		}
		db.Where("id = ?", id).First(&log)
	}

	if log.ID == 0 {
		return BagelLog{}, errors.New("no such bagel log with id " + id)
	}

	return log, nil
}

type FeedbackMsg struct {
	gorm.Model
	BagelID            uint
	IncompleteActionID string
	PlannedActionID    string
	CompletedActionID  string
}

func FeedbackMsg_Create(db *gorm.DB) FeedbackMsg {
	feedbackMsg := FeedbackMsg{
		IncompleteActionID: "feedback_msg:" + uuid.New().String(),
		PlannedActionID:    "feedback_msg:" + uuid.New().String(),
		CompletedActionID:  "feedback_msg:" + uuid.New().String(),
	}

	db.Create(&feedbackMsg)

	return feedbackMsg
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
	db.AutoMigrate(&User{}, &Tag{}, &Bagel{}, &BagelLog{}, &FeedbackMsg{})
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
	bagelChats, err := s.FindChannel("bagel-testing", "")
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

func SyncBagels(db *gorm.DB, s *Slack) (err error) {
	var bagels []Bagel
	db.Find(&bagels)

	for _, bagel := range bagels {
		slackUserIds, err := s.ConversationsMembers(bagel.SlackConversationID)
		if err != nil {
			return err
		}

		db.Model(&bagel).Association("Users").Clear()

		for _, slackId := range slackUserIds {
			var dbUser User
			db.Where("slack_id = ?", slackId).Find(&dbUser)
			if dbUser.ID == 0 {
				log.Errorf("no such user with slack id %s. Did we perform a SyncUser(db, s) before this?", slackId)
				dbUser.SlackID = slackId
			}

			db.Model(&bagel).Association("Users").Append(&dbUser)
		}
	}

	return nil
}
