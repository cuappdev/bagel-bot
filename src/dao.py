import time

from sqlalchemy import desc

from db import session, BagelInstance, Chat, ChatStatus, User
from messages import final_reminder, introduction, reminder


def print_thick_spacer():
    print('====================')
    print('\n')

def print_thin_spacer():
    print('--------------------')
    print('\n')


# User


def user_print_all():
    print_thick_spacer()
    print('Users')

    users = User.query().all()
    for user in users:
        print(user.name)

    print_thin_spacer()


def user_fetch(slack_id, name):
    user = User.query().filter_by(slack_id=slack_id).first()
    if user:
        return user
    else:
        return user_create(slack_id, name)


def user_create(slack_id, name):
    user = User(slack_id, name, True)
    session.add(user)
    session.commit()
    return user



# Bagel Instance


def bagel_instance_print_all():
    print_thick_spacer()
    print('Bagel Instances')

    bagel_instances = BagelInstance.query().all()
    for instance in bagel_instances:
        print(str(instance.id) + ' ' +  str(instance.bagel_date))

    print_thin_spacer()


def bagel_instance_create():
    bagel_instance = BagelInstance(int(time.time()))
    session.add(bagel_instance)
    session.commit()
    return bagel_instance

def bagel_instance_get_current():
    return BagelInstance.query().order_by(desc(BagelInstance.bagel_date)).first()


# Chat


def chat_print_all():
    print_thick_spacer()
    print('Chats')

    chats = Chat.query().all()
    for chat in chats:
        print('' + str(chat.bagel_instance_id) + ' ' + str([user.name for user in chat.users]))

    print_thin_spacer()


def chat_create(bagel_instance, users, slack_id):
    chat = Chat(slack_id)

    chat.slack_id = slack_id
    chat.bagel_instance_id = bagel_instance.id
    chat.users.extend(users)

    session.add(chat)
    session.commit()
    return chat

def chat_message_all(message):
    bagel_instance = bagel_instance_get_current()

    for chat in bagel_instance.chats:
        slack.message_channel(chat.slack_id, message)
