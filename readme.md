# Simple Mastodon Bot to boost all mentions of a certain hashtag of the people it's following

### Design

- follow configured list of hashtags using [StreamingAPI](https://github.com/tootsuite/documentation/blob/master/Using-the-API/Streaming-API.md)
- for all received [status](https://github.com/tootsuite/documentation/blob/master/Using-the-API/API.md#status)
  - check if status is public (```visibility==public```)
  - check if we have not reblogged status already (```reblogged==false||reblogged==nil```)
  - verify that tags contains one of our configured hashtags (```tags```)
  - check we have not blocked the conversation (```blocked==false||blocked==nil```)
  - check that it is an original status (```reblog==nil```)
  - get [relationship](https://github.com/tootsuite/documentation/blob/master/Using-the-API/API.md#relationship) with the posting [account](https://github.com/tootsuite/documentation/blob/master/Using-the-API/API.md#account) (```GET /api/v1/accounts/relationships {id:status.account.id}```)
    - verify that we follow this account (<tt>following==true && blocking==false && muting==false</tt>)
  - Boost Status if all checks passs



