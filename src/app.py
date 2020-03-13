import sys

import dao
import divvy
import messages
import slack


GROUP_SIZE = 4


def create_bagel_instance(for_real):
    bagelers = slack.bagelers_slack_ids()
    users = [dao.user_fetch(slack_id, slack.name_get(slack_id)) for slack_id in bagelers]

    groups = divvy.divvy_up(users, GROUP_SIZE)
    print([[user.name for user in group] for group in groups])

    if for_real:
        instance = dao.bagel_instance_create()
        for group in groups:
            slack_group = [user.slack_id for user in group]
            slack_id = slack.mpim_create(slack_group)
            slack.message_channel(slack_id, messages.introduction())
            dao.chat_create(instance, group, slack_id)


if __name__ == '__main__':
    if len(sys.argv) <= 1:
        print('Provide a command')

    command = sys.argv[1]

    if command == 'db-contents':
        dao.user_print_all()
        dao.chat_print_all()
        dao.bagel_instance_print_all()

    elif command == 'print':
        create_bagel_instance(False)

    elif command == 'make':
        create_bagel_instance(True)

    elif command == 'remind':
        dao.chat_message_all(messages.reminder())

    elif command == 'remind-final':
        dao.chat_message_all(messages.final_reminder())

    else:
        print('Invalid command: ' + sys.argv[1])
