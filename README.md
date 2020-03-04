![](bagel-logo.png =200x200)

# bagel-bot

*to help slack members group up and get food*

### steps:
1. clone the repository
2. create bagel bot
   - a) get a slack bot api key. you'll need to give it some permissions. probably `channels:read`, `chat:write`, `groups:write`, `mpim:write`, `im:write`, `users:read`, etc.
   - b) add the api key to your shell environment
   - c) add the slack bot to your slack workspace
   - d) add the bot to the channels you want to divvy up
3. run the app
   - a) install requirements
   - b) `python3 src/app.py <actions>` where `<actions>` is one of `print` or `make`.
        Using `print` will simulate a group matching, but will not create matches in your channel.
        Using  `make` will execute a real group matching process. 

### divvy algorithm:
let `n` be the number of people and `k` be the target group size.
the algorithm requires more than `n(n - 1)` people to be in the channel.
it generates the maximum number of `k` sized groups with `k-1` sized groups if necessary to fill the gaps.
