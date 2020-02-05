import json
import os
import random
import requests
import sys
import time

bearer_token = 'Bearer ' + os.getenv('API_KEY')
GROUP_SIZE = 4
SLACK_API = 'https://slack.com/api/'


def slack_pr(endpoint, data):
    return requests.post(
        SLACK_API + endpoint,
        data = json.dumps(data),
        headers = {
            'Authorization': bearer_token,
            'Content-type': 'application/json'
        }
    ).json()


def slack_gr(endpoint, params=None):
    return requests.get(
        SLACK_API + endpoint,
        params = params,
        headers = {
            'Authorization': bearer_token,
            'Content-type': 'application/x-www-form-urlencoded'
        }
    ).json()


profile_cache = {}

def slack_get_profile(user):
    global profile_cache
    if user in profile_cache:
        return profile_cache[user]

    res = slack_gr('users.profile.get', { 'user': user })['profile']
    profile_cache[user] = res
    return res


def slack_get_name(user):
    return slack_get_profile(user)['real_name']

    
def slack_get_bagelers():
    res = slack_gr('users.conversations')
    bagelers = set()
    for channel in res['channels']:
        channel_id = channel['id']

        members = slack_gr('conversations.members',
                           { 'channel': channel_id })['members']

        for user in members:
            profile = slack_get_profile(user)
            if profile['real_name'] != 'bagel':
                bagelers.add(user)

    return list(bagelers)


def slack_create_group(users):
    res = slack_pr('conversations.open', { 'users': users })
    channel_id = res['channel']['id']

    res = slack_pr('chat.postMessage', {
        'text': 'Hello, is this thing on?',
        'channel': channel_id
    })
    res = slack_pr('chat.postMessage', {
        'text': 'Well anyways, enjoy your group coffee chat!',
        'channel': channel_id
    })


def form_groups(members, num_groups, group_size):
    return [[members.pop() for i in range(group_size)] for j in range(num_groups)]


def divvy_up(members, group_size):
    mbrs = members.copy()

    small_groups = (group_size - (len(members) % group_size)) % group_size
    large_groups = (len(members) - (group_size - 1) * small_groups) // group_size

    groups = [
        *form_groups(mbrs, small_groups, group_size - 1),
        *form_groups(mbrs, large_groups, group_size)
    ]

    return groups


def slack_message(message):
    res = slack_gr('users.conversations')
    for channel in res['channels']:
        channel_id = channel['id']
        res = slack_pr('chat.postMessage', {
            'text': message,
            'channel': channel_id
        })


if __name__ == '__main__':
    if sys.argv[1] == 'printm':
        print('Printing a message')
        print(sys.argv[2])
        slack_message(sys.argv[2])

    else:
        bagelers = slack_get_bagelers()
        random.shuffle(bagelers)
        groups = divvy_up(bagelers, GROUP_SIZE)

        if sys.argv[1] == 'print':
            print('Printing slack groups')
            for group in groups:
                time.sleep(5)
                print([slack_get_name(member) for member in group])

        elif sys.argv[1] == 'make':
            print('Making slack groups for real')
            for group in groups:
                time.sleep(5)
                print([slack_get_name(member) for member in group])
                slack_create_group(group)

