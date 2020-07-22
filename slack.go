package main

import (
    "errors"
    "io/ioutil"
    "net/http"
    "encoding/json"
    "github.com/mitchellh/mapstructure"
)

type Slack struct {
    Token string
    Client http.Client
}

type SlackUser struct {
    ID string `mapstructure:"id"`
    TeamID string `mapstructure:"team_id"`
    Name string `mapstructure:"name"`
    Deleted bool `mapstructure:"deleted"`
    Profile SlackUserProfile `mapstructure:"profile"`
    IsAdmin bool `mapstructure:"is_admin"`
    IsOwner bool `mapstructure:"is_owner"`
    IsBot bool `mapstructure:"is_bot"`
    IsAppUser bool `mapstructure:"is_app_user"`
}

type SlackUserProfile struct {
    RealName string `mapstructure:"real_name"`
    DisplayName string `mapstructure:"display_name"`
}

func (s Slack) get(endpoint string, params map[string]string) ([]byte, error) {
    req, err := http.NewRequest("GET", "https://slack.com/api/" + endpoint, nil)
    if err != nil {
        return nil, err
    }

    q := req.URL.Query()
    q.Add("token", s.Token)
    for key, value := range params {
    	q.Add(key, value)
    }
    req.URL.RawQuery = q.Encode()

    resp, err := s.Client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    bytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    return bytes, nil
}

func (s Slack) usersList() ([]SlackUser, error) {
    resp, err := s.get("users.list", nil)
    if err != nil {
        return nil, err
    }

    var topLevelJsonResp map[string]interface{}
    err = json.Unmarshal(resp, &topLevelJsonResp)
    if err != nil {
        return nil, err
    }

    ok, typeConversionSucceeded := topLevelJsonResp["ok"].(bool)
    if !typeConversionSucceeded {
        log.Critical("Could not cast \"ok\" to type bool")
        return nil, nil
    }

    if !ok {
         errStr, typeConversionSucceeded := topLevelJsonResp["error"].(string)
         if !typeConversionSucceeded {
             log.Critical("Could not cast \"error\" to type string")
             return nil, nil
         }

         return nil, errors.New(errStr)
    }

    membersJson, present := topLevelJsonResp["members"]
    if !present {
        return nil, errors.New("response is missing members")
    }

    var users []SlackUser
    err = mapstructure.Decode(membersJson, &users)
    if err != nil {
        return nil, err
    }

    return users, nil
}

