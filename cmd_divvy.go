package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"gorm.io/gorm"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type CmdDivvy struct {
	DryRun bool   `help:"Divvy users and print the result without actually making bagel chats"`
	Size   int    `help:"Number of people per bagelchat" default:"2"`
	Tag    string `arg optional help:"Divvy the tagged users. By default, bagel divvies the users in the #bagel-chats channel"`
}

func (cmd *CmdDivvy) Run(ctx *kong.Context, db *gorm.DB, s *Slack, invocation *Invocation) (err error) {
	if cmd.Size < 2 {
		if _, err = io.WriteString(ctx.Stderr, "size must be >= 2"); err != nil {
			return err
		}
	}

	if err = SyncUsers(db, s); err != nil {
		return err
	}

	var users []User
	if cmd.Tag == "" {
		db.Where("is_bagel_chats_user <> 0").Find(&users)
	} else {
		var tag Tag
		db.Where("name = ?", cmd.Tag).First(&tag)
		if tag.ID == 0 {
			_, err = io.WriteString(ctx.Stderr, "no such tag "+cmd.Tag+"\n")
			return err
		}

		db.Model(&tag).Association("Users").Find(&users)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(users), func(i, j int) {
		temp := users[i]
		users[i] = users[j]
		users[j] = temp
	})

	var divvied [][]User
	if cmd.Size == 2 {
		divvied = pair(users)
	} else {
		minUsers := cmd.Size * (cmd.Size - 1)
		if len(users) < minUsers {
			msg := fmt.Sprintf("need at least %d users to divvy with size %d; %s only has %d users.", minUsers, cmd.Size, cmd.Tag, len(users))
			if _, err = io.WriteString(ctx.Stderr, msg); err != nil {
				return err
			}
		}
		divvied = divvy(users, cmd.Size)
	}

	if err = printDivvied(divvied, ctx.Stdout); err != nil {
		return err
	}

	if cmd.DryRun {
		return
	}

	bagelLog := BagelLog{
		Date:       time.Now().Unix(),
		Invocation: string(*invocation),
	}
	db.Create(&bagelLog)

	for _, group := range divvied {
		var userIds []string
		for _, user := range group {
			userIds = append(userIds, user.SlackID)
		}

		slackConvId, err := s.ConversationsOpen(userIds)
		if err != nil {
			return err
		}

		db.Model(&bagelLog).Association("Bagel")
		bagel := Bagel{
			SlackConversationID: slackConvId,
		}
		db.Create(&bagel)
		db.Model(&bagel).Association("Users").Append(&users)
		db.Model(&bagelLog).Association("Bagels").Append(&bagel)

		if err = addIntroduction(s, slackConvId, userIds); err != nil {
			return err
		}
	}

	return nil
}

func pair(users []User) [][]User {
	if len(users) <= 3 {
		return [][]User{users}
	}

	numGroups := len(users) / 2
	divvied := make([][]User, numGroups)
	for i := 0; i < numGroups-1; i++ {
		divvied[i] = users[2*i : 2*i+2]
	}
	divvied[numGroups-1] = users[2*(numGroups-1):]

	return divvied
}

func divvy(users []User, size int) [][]User {
	if size < 3 {
		panic("size must be >= 3")
	}
	minUsers := size * (size - 1)
	if len(users) < minUsers {
		panic("there must be at least " + strconv.Itoa(minUsers) + " users n*(n-1)")
	}

	var numSmall int
	var numLarge int
	if len(users)%size == 0 {
		numSmall = 0
		numLarge = len(users) / size
	} else {
		numSmall = (size - (len(users) % size)) % size
		numLarge = (len(users) - (size-1)*numSmall) / size
	}

	divvied := make([][]User, numSmall+numLarge)
	for i := 0; i < numSmall; i++ {
		divvied[i] = users[i*(size-1) : (i+1)*(size-1)]
	}
	for i := 0; i < numLarge; i++ {
		divvied[numSmall+i] = users[numSmall*(size-1)+i*size : numSmall*(size-1)+(i+1)*size]
	}
	return divvied
}

func printDivvied(divvied [][]User, stdout io.Writer) error {
	for i, group := range divvied {
		var usernames []string
		for _, user := range group {
			usernames = append(usernames, user.Name)
		}

		msg := fmt.Sprintf("Group %d: %s\n", i+1, strings.Join(usernames, ", "))
		if _, err := io.WriteString(stdout, msg); err != nil {
			return err
		}
	}
	return nil
}

func addIntroduction(s *Slack, channelId string, userIds []string) (err error) {
	introductions := [][]string{
		{
			"Welcome to your bagel partner for this week! For the first month of chats, you will be randomly paired with people in your subteam :).",
		},
	}

	rand.Seed(time.Now().UnixNano())
	lines := introductions[rand.Intn(len(introductions))]
	for _, line := range lines {
		if err = s.ChatPostMessage(channelId, line, nil); err != nil {
			return err
		}
	}
	return nil
}
