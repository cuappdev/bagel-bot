package main

import (
	"encoding/json"
	"errors"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Slack struct {
	Token      string
	HttpClient http.Client
}

type SlackResponseMetadata struct {
	NextCursor string `mapstructure:"next_cursor"`
}

type SlackUser struct {
	ID        string           `mapstructure:"id"`
	TeamID    string           `mapstructure:"team_id"`
	Name      string           `mapstructure:"name"`
	Deleted   bool             `mapstructure:"deleted"`
	Profile   SlackUserProfile `mapstructure:"profile"`
	IsAdmin   bool             `mapstructure:"is_admin"`
	IsOwner   bool             `mapstructure:"is_owner"`
	IsBot     bool             `mapstructure:"is_bot"`
	IsAppUser bool             `mapstructure:"is_app_user"`
}

type SlackUserProfile struct {
	RealName    string `mapstructure:"real_name"`
	DisplayName string `mapstructure:"display_name"`
}

type SlackChannel struct {
	ID         string `mapstructure:"id"`
	Name       string `mapstructure:"name"`
	IsChannel  bool   `mapstructure:"is_channel"`
	IsGroup    bool   `mapstructure:"is_group"`
	IsIm       bool   `mapstructure:"is_im"`
	Created    int    `mapstructure:"created"`
	Creator    string `mapstructure:"creator"`
	IsArchived bool   `mapstructure:"is_archived"`
	IsGeneral  bool   `mapstructure:"is_general"`
	IsMember   bool   `mapstructure:"is_member"`
	IsPrivate  bool   `mapstructure:"is_private"`
	IsMpim     bool   `mapstructure:"id_mpim"`
	NumMembers int    `mapstructure:"num_members"`
}

func (s Slack) request(method string, endpoint string, params map[string]string, contentKey *string) (interface{}, *SlackResponseMetadata, error) {
	req, err := http.NewRequest(method, "https://slack.com/api/"+endpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	q := req.URL.Query()
	q.Add("token", s.Token)

	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := s.HttpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	var topLevelJsonResp map[string]interface{}
	err = json.Unmarshal(bytes, &topLevelJsonResp)
	if err != nil {
		return nil, nil, err
	}

	ok, typeConversionSucceeded := topLevelJsonResp["ok"].(bool)
	if !typeConversionSucceeded {
		log.Critical("Could not cast \"ok\" to type bool")
		return nil, nil, nil
	}

	if !ok {
		errStr, typeConversionSucceeded := topLevelJsonResp["error"].(string)
		if !typeConversionSucceeded {
			log.Critical("Could not cast \"error\" to type string")
			return nil, nil, nil
		}

		return nil, nil, errors.New(errStr)
	}

	var contentJson interface{}
	if contentKey != nil {
		var present bool
		contentJson, present = topLevelJsonResp[*contentKey]
		if !present {
			return nil, nil, errors.New("response is missing " + *contentKey)
		}
	}

	metadataJson, present := topLevelJsonResp["response_metadata"]
	var metadata *SlackResponseMetadata
	if present {
		err = mapstructure.Decode(metadataJson, &metadata)
		if err != nil {
			metadata = nil
		}
	}

	return contentJson, metadata, nil
}

func (s Slack) get(endpoint string, params map[string]string, contentKey *string) (interface{}, *SlackResponseMetadata, error) {
	return s.request("GET", endpoint, params, contentKey)
}

func (s Slack) post(endpoint string, params map[string]string, contentKey *string) (interface{}, *SlackResponseMetadata, error) {
	return s.request("POST", endpoint, params, contentKey)
}

func (s Slack) cursorCollect(endpoint string, params map[string]string, contentKey *string, initialValue interface{}, collect func(interface{}, interface{}) (interface{}, error)) (interface{}, error) {
	cursor := ""

	for {
		if cursor != "" {
			params["cursor"] = cursor
		}

		jsonContent, metadata, err := s.get(endpoint, params, contentKey)
		if err != nil {
			return nil, err
		}

		initialValue, err = collect(initialValue, jsonContent)
		if err != nil {
			return nil, err
		}

		if metadata == nil {
			break
		}

		cursor = metadata.NextCursor
		if cursor == "" {
			break
		}
	}

	return initialValue, nil
}

func (s Slack) getChannels(endpoint string, params map[string]string) ([]SlackChannel, error) {
	contentKey := "channels"
	collected, err := s.cursorCollect(
		endpoint,
		params,
		&contentKey,
		[]SlackChannel{},
		func(collected interface{}, json interface{}) (interface{}, error) {
			var partial []SlackChannel
			err := mapstructure.Decode(json, &partial)
			if err != nil {
				return nil, err
			}
			return append(collected.([]SlackChannel), partial...), nil
		})
	if err != nil {
		return nil, err
	}

	return collected.([]SlackChannel), nil
}

func (s Slack) getMembers(endpoint string, params map[string]string) ([]SlackUser, error) {
	contentKey := "members"
	collected, err := s.cursorCollect(
		endpoint,
		params,
		&contentKey,
		[]SlackUser{},
		func(collected interface{}, json interface{}) (interface{}, error) {
			var partial []SlackUser
			err := mapstructure.Decode(json, &partial)
			if err != nil {
				return nil, err
			}
			return append(collected.([]SlackUser), partial...), nil
		})
	if err != nil {
		return nil, err
	}
	return collected.([]SlackUser), nil
}

func (s Slack) getStrings(endpoint string, params map[string]string, contentKey string) ([]string, error) {
	collected, err := s.cursorCollect(
		endpoint,
		params,
		&contentKey,
		[]string{},
		func(collected interface{}, json interface{}) (interface{}, error) {
			var partial []string
			err := mapstructure.Decode(json, &partial)
			if err != nil {
				return nil, err
			}
			return append(collected.([]string), partial...), nil
		})
	if err != nil {
		return nil, err
	}
	return collected.([]string), nil
}

func (s Slack) ApiTest() error {
	_, _, err := s.get("api.test", map[string]string{"foo": "bar"}, "args")
	return err
}

func (s Slack) ChatPostMessage(channel string, text string) error {
	params := map[string]string{"channel": channel, "text": text}
	_, _, err := s.post("chat.postMessage", params, nil)
	return err
}

func (s Slack) ConversationsList(excludeArchived bool, types []string) ([]SlackChannel, error) {
	params := map[string]string{"exclude_archived": strconv.FormatBool(excludeArchived)}
	if len(types) != 0 {
		params["types"] = strings.Join(types, ",")
	}

	return s.getChannels("conversations.list", params)
}

func (s Slack) ConversationsMembers(channel string) ([]string, error) {
	params := map[string]string{"channel": channel}
	return s.getStrings("conversations.members", params, "members")
}

func (s Slack) ConversationsOpen(users []string) (string, error) {
	params := map[string]string{"users": strings.Join(users, ",")}
	channelJson, _, err := s.get("conversations.open", params, "channel")
	if err != nil {
		return "", err
	}
	return channelJson.(map[string]string)["id"], nil
}

func (s Slack) UsersConversations(excludeArchived bool, types []string) ([]SlackChannel, error) {
	params := map[string]string{"exclude_archived": strconv.FormatBool(excludeArchived)}
	if len(types) != 0 {
		params["types"] = strings.Join(types, ",")
	}

	return s.getChannels("users.conversations", params)
}

func (s Slack) UsersList() ([]SlackUser, error) {
	return s.getMembers("users.list", nil)
}
