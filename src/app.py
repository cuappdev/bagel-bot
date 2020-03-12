import sys

import dao
import divvy
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
            dao.chat_create(instance, group, slack_id)


if __name__ == '__main__':
    if len(sys.argv) <= 1:
        print('Provide a command')
    
    elif sys.argv[1] == 'db-contents':
        dao.user_print_all()
        dao.chat_print_all()
        dao.bagel_instance_print_all()

    elif sys.argv[1] == 'print':
        create_bagel_instance(False)
    
    elif sys.argv[1] == 'make': 
        create_bagel_instance(True)

    else:
        print('Invalid command: ' + sys.argv[1])
