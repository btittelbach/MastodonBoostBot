# Simple Mastodon Bot to boost all mentions of a certain hashtag of the people it's following

## Idea

The original usecase is to collect all the toots mentioning one of a set of hashtags by a group of people and re-toot/boost their stati in a single account. The single account can be displayed on a kiosk screen and be the public face of a organisation/group/hackerspace. You can use MastodonBoostBot to emulate groups in Mastodon. Users can post as an organisation account without requiring login credentials and without loosing the information about the original poster.

Administer your group (people whose tagged posts will be boosted) in Mastodon by following/unfollowing users.

Alternatively, you can disable the follower-filter and just boost/reblog every post containing one of your hashtags.

### Features

- mostly 12-factor app
    - configuration with environment variables
    - but no package control until go1.12 is more widespread
- easily containerized static golang executable
- fail-fast, restart-fast
- configure once and forget (administrate trusted followers in mastodon, not bot configuration)

### similar software

- [tootgroup](https://github.com/oe4dns/tootgroup.py)
    - only needs to run periodically from cron while MastodonBoostBot needs to run constantly as a well defined service, otherwise it will miss a post
    - hiding the original poster, by extracting content from direct messages, is currently not planned in MastodonBoostBot

- [boost-bot](https://github.com/Gargron/boost-bot)
    - Very similar, but does not filter by followers

## Installation


### Build

    go build

### Configure

1. In your Mastodon account-page, create a new application in the development app
    - with the following permissions: ```read:accounts``` ```read:statuses``` ```read:search``` ```read:follows``` ```write:conversations``` ```write:statuses``` ```push```
2. copy ```example.env``` and fill the copy with the appropriately app-credentials
  e.g. to ```mastodonBoostBot_config1.env```
3. configure the hashtags you want to follow (without the leading ```#```)

### Carbon Copy to Twitter

1. You will need an App Consumer Key and App Consumer Secret from a Twitter Developer (or get a Twitter Developer to run MastodonBoostBot for you)
2. run
 ```MBB_CCTWEET_CONSUMER_KEY=___ MBB_CCTWEET_CONSUMER_SECRET=___ ./MastodonBoostBot --starttwitteroauth``` after replacing ```___```
 with your key and secret.
3. authorize or let authorize the twitter account you want to use
4. enter the PIN on the cmdline
5. copy&paste the resulting settings into your environment file for the service

### Test locally

    source <(sed '/^\s*#/d;s/\(.\+\)=\(.*\)/export \1="\2"/;' mastodonBoostBot_config1.env)
    ./MastodonBoostBot --debug=ALL

### Deploy

1. deploy executable and config file to server of your choice
2. adapt the service file to your paths and needs or otherwise ensure that environment variables are loaded from config-file.
3. enable and start
   you can easily run multiple instances with different configurations


#### Example 1 - as user on a systemd system

1. copy <tt>MastodonBoostBot</tt> executable to <code>~/bin/</code>
1. copy config to <code>~/.config/mastodonBoostBot_config1.env</code>
1. <code>chmod 400 ~/.config/mastodonBoostBot_config1.env</code>
1. copy <code>user_service/MastodonBoostBot.service</code> to <code>~/.local/share/systemd/user/</code>
1. <code>systemctl --user start MastodonBoostBot.service</code>
1. <code>systemctl --user enable MastodonBoostBot.service</code>


#### Example 2 - systemwide on a systemd system

1. <code>mkdir -p /opt/MastodonBoostBot/</code>
1. <code>adduser mastodonbot</code>
1. <code>chown mastodonbot -R /opt/MastodonBoostBot/</code>
1. copy <tt>MastodonBoostBot</tt> executable to <code>/opt/MastodonBoostBot/</code>
1. copy config to <code>/opt/MastodonBoostBot/mastodonBoostBot_config1.env</code>
1. <code>chmod 400 /opt/MastodonBoostBot/mastodonBoostBot_config1.env</code>
1. copy <code>user_service/MastodonBoostBot.service</code> to <code>/etc/systemd/system</code>
1. adapt paths and config in <tt>MastodonBoostBot.service</tt>
    a. add <code>User=mastodonbot</code>
    a. <code>ExecStart</code>
    a. <code>EnvironmentFile</code>
    a. change <code>WantedBy</code> to <code>multi-user.target</code> or whatever you start your services with
    a. comment in additional protections. possibly try chrooting the service in <code>/opt/MastodonBoostBot</code>
1. <code>systemctl start MastodonBoostBot.service</code>
1. <code>systemctl enable MastodonBoostBot.service</code>


## Design

### Algorithm

- follow stream of configured list of hashtags using [StreamingAPI](https://github.com/tootsuite/documentation/blob/master/Using-the-API/Streaming-API.md)
- for all received [status](https://github.com/tootsuite/documentation/blob/master/Using-the-API/API.md#status)
    - Boost Status if the following optional checks pass
        - is status public (```visibility==public```)
        - not already reblogged (```reblogged==false||reblogged==nil```)
        - filter for configured hashtags
        - conversation unblocked? (```blocked==false||blocked==nil```)
        - original status or already reblogged? (```reblog==nil```)
        - verify it is a user we follow.
            - get [relationship](https://github.com/tootsuite/documentation/blob/master/Using-the-API/API.md#relationship) with the posting [account](https://github.com/tootsuite/documentation/blob/master/Using-the-API/API.md#account) (```GET /api/v1/accounts/relationships {id:status.account.id}```)
            - verify we follow this account (<tt>following==true && blocking==false && muting==false</tt>)
        - considered but not implemented: check against configured list of users (though that would move user-management from website to daemon-configfile)
        - planned but not implemented: filter for addtional content in status-text


- considered but not implemented: read "private" messages and toot them


