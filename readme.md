# Simple Mastodon Bot to boost all mentions of a certain hashtag of the people it's following

## Idea

The original usecase is to collect all the toots mentioning one of a set of hashtags by a group of people and re-toot/boost their stati in a single account. The single account can be displayed on a kiosk screen and be the public face of a organisation/group/hackerspace. You can use MastodonBoostBot to emulate groups in Mastodon. Users can post as an organisation account without requiring login credentials and without loosing the information about the original poster.

### Features

- mostly 12-factor app
    - configuration with environment variables
    - but no package control until go1.12 is more widespread
- easily containerized static golang executable
- fail-fast, restart-fast
- start once, never change configuration again (administrate trusted followers in mastodon, not bot configuration)

### similar software

- [tootgroup](https://github.com/oe4dns/tootgroup.py)
    - only needs to run periodically from cron while MastodonBoostBot needs to run constantly as a well defined service, otherwise it will miss a post
    - hiding the original poster, by extracting content from direct messages, is currently not planned in MastodonBoostBot

- [boost-bot](https://github.com/Gargron/boost-bot)
    - Very similar, but does not filter by followers

## Installation


### Build

    go get .
    go build

### Configure

- In your Mastodon account-page, create a new application in the development app
    - with the following permissions: <tt>read:accounts</tt> <tt>read:statuses</tt> <tt>read:search</tt> <tt>read:follows</tt> <tt>write:conversations</tt> <tt>write:statuses</tt> <tt>push</tt>
- copy <tt>example.env</tt> and fill the copy with the appropriately app-credentials
  e.g. to <tt>mastodonBoostBot_config1.env</tt>
- configure the hashtags you want to follow (without the leading <tt>#</tt>)

### Test locally

    source <(sed -E -n 's/[^#]+/export &/ p' mastodonBoostBot_config1.env)
    ./MastodonBoostBot --debug=ALL

### Deploy

- deploy executable and config to server of your choice
  e.g. when using <tt>systemd --user</tt> context copy config file to <tt>~/.config/mastodonBoostBot_config1.env</tt> and executable to <tt>~/bin/</tt> on your server
- run it as service e.g. with provided systemd service file
  - adapt the service file to your paths and needs
  - if your user is allowed to linger (sudo loginctl enable-linger <username>) you can install and run the service-file without root-permissions for your user only in <tt>~/.local/share/systemd/user/</tt>
- you can easily run multiple instances with different configurations this way


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


