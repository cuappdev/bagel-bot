# Divvy


def form_groups(members, num_groups, group_size):
    return [[members.pop() for i in range(group_size)] for j in range(num_groups)]


def divvy_up(members, group_size):
    mbrs = members.copy()

    if len(members) % group_size == 0:
        small_groups = 0
        large_groups = len(members) // group_size
    else:
        small_groups = (group_size - (len(members) % group_size)) % group_size
        large_groups = (len(members) - (group_size - 1) * small_groups) // group_size

    print('small_groups' + str(small_groups))
    print('large_groups' + str(large_groups))

    groups = [
        *form_groups(mbrs, small_groups, group_size - 1),
        *form_groups(mbrs, large_groups, group_size)
    ]

    return groups
