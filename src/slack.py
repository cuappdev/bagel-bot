import json
import os
import random
import sys
import time

import requests


# Constants


bearer_token = 'Bearer ' + os.getenv('API_KEY')
BAGEL_ICON_URL = 'https://raw.githubusercontent.com/cuappdev/bagel-bot/master/bagel-logo.png'
SLACK_API = 'https://slack.com/api/'
TESTING_CHANNEL_ID = 'CTF81MFH6'


# Pull / Get Request


def pr(endpoint, data):
    try: requests.
    return requests.post(
        SLACK_API + endpoint,
        data = json.dumps(data),
        headers = {
            'Authorization': bearer_token,
            'Content-type': 'application/json'
        }
    ).json()


def gr(endpoint, params=None):
    return requests.get(
        SLACK_API + endpoint,
        params = params,
        headers = {
            'Authorization': bearer_token,
            'Content-type': 'application/x-www-form-urlencoded'
        }
    ).json()


# Profiles


profile_cache = {}


def profile_get(slack_id):
    global profile_cache
    if slack_id in profile_cache:
        return profile_cache[slack_id]
        
    res = gr('users.profile.get', { 'user': slack_id })
    profile_cache[slack_id] = res['profile']
    return res['profile']


def name_get(slack_id):
    return profile_get(slack_id)['real_name']


# Messaging


def channels_get_all():
    res = gr('users.conversations')
    return res['channels']


def channel_leave(channel_id):
    pr('conversations.leave', {
        'channel': channel_id
    })


def message_channel(message, channel_id):
    return pr('chat.postMessage', {
        'channel': channel_id,
        'icon_url': BAGEL_ICON_URL,
        'text': message,
        'username': 'bagel',
    })


def message_test_channel(message):
    message_channel(message, TESTING_CHANNEL_ID)


# Group Management


# Create a group chat (multi-person instant message)
def mpim_create(users):
    res = pr('conversations.open', { 'users': users })
    channel_id = res['channel']['id']

    res = message_channel('Welcome to bagel chats, with a new and improved icon :bagel-bot::bagel-bot::bagel-bot:', channel_id)
    return channel_id


def mpim_get_all():
    return gr('users.conversations', { 'types': 'mpim' })


def bagelers_slack_ids(): 
    res = gr('users.conversations')
    bagelers = set()
    for channel in res['channels']:
        channel_id = channel['id']

        members = gr('conversations.members',
                           { 'channel': channel_id })['members']

        for slack_id in members:
            if name_get(slack_id) == 'bagel':
                continue

            bagelers.add(slack_id)

    return list(bagelers)
