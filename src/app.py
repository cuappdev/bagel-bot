import json
import requests
import os
import random
import sys

bearer_token = 'Bearer ' + os.getenv('API_KEY')

def slack_pr(endpoint, data): 
    return requests.post(
        endpoint,
        data = json.dumps(data),
        headers = { 
            'Authorization': bearer_token,
            'Content-type': 'application/json'
        }
    ).json()

def slack_gr(endpoint, params):
    return requests.get(
        endpoint,
        params = params,
        headers = { 
            'Authorization': bearer_token,
            'Content-type': 'application/x-www-form-urlencoded'
        }
    ).json()

def slack_get_bagelers():
    res = slack_gr('https://slack.com/api/users.conversations', None) 
    bagelers = set()
    for channel in res['channels']:
        channel_id = channel['id']

        members = slack_gr('https://slack.com/api/conversations.members',
                           { 'channel': channel_id })['members']
        
        for user in members:
            bagelers.add(user)
    
    return list(bagelers)

def slack_create_group(users):
    res = slack_pr('https://slack.com/api/conversations.open', { 'users': users })
    print(res)
    channel_id = res['channel']['id']

    res = slack_pr('https://slack.com/api/chat.postMessage', { 
        'text': 'Hello, is this thing on?',
        'channel': channel_id
    })
    print(res)
    res = slack_pr('https://slack.com/api/chat.postMessage', { 
        'text': 'Well anyways, enjoy your first three person coffee chat!',
        'channel': channel_id
    })
    print(res)

def divvy_up(members, group_size):
    mbrs = members.copy()

    smallGroups = group_size - (len(members) % group_size)
    largeGroups = (len(members) - (group_size - 1) * smallGroups) // group_size
    groups = []

    for i in range(0, smallGroups):
        group = []
        for j in range(0, group_size - 1):
            group.append(mbrs.pop()) 

        groups.append(group)

    for i in range(0, largeGroups):
        group = []
        for j in range(0, group_size):
            group.append(mbrs.pop())

        groups.append(group)

    return groups
     

if __name__ == '__main__':
    bagelers = slack_get_bagelers()
    random.shuffle(bagelers)
    groups = divvy_up(bagelers, 4)

    for arg in sys.argv[1:]:
        if arg == 'print': 
            for group in groups:
                print(group)
        if arg == 'make':
            for group in groups:
                slack_create_group(group)

            
