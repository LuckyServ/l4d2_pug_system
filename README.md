
# l4d2center.com
## The backend layer of the ranked PUG system for Left 4 Dead 2

### GET /ws
##### Websocket endpoint
Request parameters: None

Response parameters: None

<br/><br/>

### GET /shutdown
##### An admin command (access = 4) to shutdown the program
Request parameters: None

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if command accepted, "false" otherwise |
| <strong>error</strong> | _string_ | Outputs the reason if the request is rejected |

<br/><br/>

### GET /blocknewgames
##### An admin command (access = 4) to block new games(queue) from creation
Request parameters: None

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if command accepted, "false" otherwise |
| <strong>error</strong> | _string_ | Outputs the reason if the request is rejected |

<br/><br/>

### GET /setadmin
##### An admin command (access = 4) to add or remove admins
Request parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>steamid64</strong> | _string_ | SteamID 64 |
| <strong>access</strong> | _int_ | Access |

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if command accepted, "false" otherwise |
| <strong>error</strong> | _string_ | Outputs the reason if the request is rejected |

<br/><br/>

### GET /addban
##### Add ban record
Request parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>steamid64</strong> | _string_ | SteamID 64 |
| <strong>nickname</strong> | _string_ | Nickname |
| <strong>reason</strong> | _string_ | Ban reason |
| <strong>banlength</strong> | _int64_ | Ban length in milliseconds |
| <strong>bantype</strong> | _int_ | Ban type (-3 or -2 for now) |

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if command accepted, "false" otherwise |
| <strong>error</strong> | _string_ | Outputs the reason if the request is rejected |

<br/><br/>

### GET /unban
##### Unban a player
Request parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>steamid64</strong> | _string_ | SteamID 64 |

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if command accepted, "false" otherwise |
| <strong>error</strong> | _string_ | Outputs the reason if the request is rejected |

<br/><br/>

### GET /getknownaccs
##### Get known duplicate accounts (smurfs)
Request parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>steamid64</strong> | _string_ | SteamID 64 |

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if command accepted, "false" otherwise |
| <strong>accounts</strong> | _[]string_ | Array of Steam IDs if success == true |
| <strong>error</strong> | _string_ | Outputs the reason if the request is rejected |

<br/><br/>

### GET /overridevpn
##### Override VPN info
Request parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>ip</strong> | _string_ | IP |
| <strong>isvpn</strong> | _bool_ | VPN or not |

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if command accepted, "false" otherwise |
| <strong>error</strong> | _string_ | Outputs the reason if the request is rejected |

<br/><br/>

### GET /auth
##### Open this in browser to authorize
Request parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>home_page</strong> | _string_ | Redirect to this page after authorization |

Response parameters: None (303 redirect)

<br/><br/>

### GET /logout
##### Log out
Request parameters: None

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if logged out, "false" otherwise |
| <strong>error</strong> | _string_ | Outputs the reason if the request is rejected |

<br/><br/>

### GET /status
##### Get necessary info about program status, and signal about online status
Request parameters: None

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | Always true |
| <strong>no_new_games</strong> | _bool_ | Tells if creating new games is blocked |
| <strong>brokenmode</strong> | _bool_ | Tells if competitive plugins are broken by some L4D2 update. In this mode the gameservers are vanilla + Sourcemod. |
| <strong>time</strong> | _int64_ | System time in milliseconds |
| <strong>need_update_players</strong> | _bool_ | Should update players or not |
| <strong>need_update_queue</strong> | _bool_ | Should update queue or not |
| <strong>need_update_game</strong> | _bool_ | Should update game info or not |
| <strong>need_update_globalchat</strong> | _bool_ | Should update global chat or not |
| <strong>need_emit_readyup_sound</strong> | _bool_ | Should attract player attention or not |
| <strong>authorized</strong> | _bool_ | Authorized or not |

<br/><br/>

### GET /validateprofile
##### Ask to validate client profile
Request parameters: None

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if profile validated, "false" otherwise |
| <strong>error</strong> | _string_ | Outputs the reason if the request is rejected |

<br/><br/>

### GET /updatenameavatar, GET /updatenickname
##### Update nickname and avatar from Steam
Request parameters: None

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if nickname and avatar updated, "false" otherwise |
| <strong>error</strong> | _string_ | Outputs the reason if the request is rejected |

<br/><br/>

### GET /acceptrules
##### Accept rules of the project
Request parameters: None

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if rules accepted, "false" otherwise |
| <strong>error</strong> | _string_ | Outputs the reason if the request is rejected |

<br/><br/>

### GET /acceptban
##### Confirm that the ban reason is read
Request parameters: None

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if success, "false" otherwise |
| <strong>error</strong> | _string_ | Outputs the reason if the request is rejected |

<br/><br/>

### GET /getonlineplayers
##### Get list of online players
Request parameters: None

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | Always "true" |
| <strong>authorized</strong> | _bool_ | Authorized or not |
| <strong>me</strong> |  | Info about an authorized player (only present if authorized) |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>steamid64</strong> | _string_ | Player Steam ID 64 |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>nickname_base64</strong> | _string_ | Base64 encoded nickname |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>avatar_small</strong> | _string_ | Small avatar |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>avatar_big</strong> | _string_ | Big avatar |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>mmr</strong> | _int_ | Player's rating |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>mmr_grade</strong> | _int_ | 1-11 if mmr valid and stable, 0 otherwise |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>access</strong> | _int_ | Player's access level<br>-2 - completely banned, -1 - chat banned, 0 - regular player, 1 - behaviour moderator, 2 - cheat moderator, 3 - behaviour+cheat moderator, 4 - full admin access |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>banreason</strong> | _string_ | Ban reason (base64 encoded), or empty if not banned |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>ban_accepted_at</strong> | _int64_ | When did the player confirm that he read the ban reason. 0 - if didnt confirm or not banned |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>ban_length</strong> | _int64_ | Ban length since the moment of ban confirmation, or 0 if not banned |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>is_ingame</strong> | _bool_ | Is player in game right now |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>is_inqueue</strong> | _bool_ | Is player in queue right now |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>profile_validated</strong> | _bool_ | New players must validate their profiles before playing |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>rules_accepted</strong> | _bool_ | New players must accept the rules before playing |
| <strong>count</strong> |  | Numbers |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>online</strong> | _int_ | Number of online players |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>in_game</strong> | _int_ | Number of players in games |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>in_queue</strong> | _int_ | Number of players in queue |
| <strong>list</strong> | _[]_ | Array of players |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>steamid64</strong> | _string_ | Player Steam ID 64 |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>nickname_base64</strong> | _string_ | Base64 encoded nickname |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>avatar_small</strong> | _string_ | Small avatar |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>mmr</strong> | _int_ | Player's rating |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>mmr_grade</strong> | _int_ | 1-11 if mmr valid and stable, 0 otherwise |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>access</strong> | _int_ | Player's access level<br>-2 - completely banned, -1 - chat banned, 0 - regular player, 1 - behaviour moderator, 2 - cheat moderator, 3 - behaviour+cheat moderator, 4 - full admin access |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>is_ingame</strong> | _bool_ | Is player in game right now |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>is_inqueue</strong> | _bool_ | Is player in queue right now |

<br/><br/>

### GET /getqueue
##### Get queue info
Request parameters: None

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | Always "true" |
| <strong>authorized</strong> | _bool_ | Authorized or not |
| <strong>is_inqueue</strong> | _bool_ | In queue or not (if authorized) |
| <strong>ready_state</strong> | _bool_ | Is queue in ReadyUp state |
| <strong>need_readyup</strong> | _bool_ | ReadyUp requested or not |
| <strong>player_count</strong> | _int_ | Number of players in queue |
| <strong>waiting_since</strong> | _int64_ | Since when this queue is waiting for players |
| <strong>ready_players</strong> | _int_ | Number of ready players |
| <strong>finishing_game</strong> | _int_ | Number of players finishing game soon |

<br/><br/>

### GET /joinqueue
##### Join queue
Request parameters: None

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if command accepted, "false" otherwise |
| <strong>error</strong> | _string_ | Outputs the reason if the request is rejected |

<br/><br/>

### GET /leavequeue
##### Leave queue
Request parameters: None

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if command accepted, "false" otherwise |
| <strong>error</strong> | _string_ | Outputs the reason if the request is rejected |

<br/><br/>

### GET /readyup
##### Ready up, when your queue is in ReadyUp state
Request parameters: None

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if command accepted, "false" otherwise |
| <strong>error</strong> | _string_ | Outputs the reason if the request is rejected |

<br/><br/>

### GET /getgame
##### Get the current game info of an authorized player
Request parameters: None

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if success, "false" if not authorized or not in game |
| <strong>game</strong> |  | Info about game, if success == true |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>id</strong> | _string_ | Game ID |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>status</strong> | _string_ | Current state of the game, as readable text |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>players_a</strong> | _[]_ | Array of players of team A |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>players_b</strong> | _[]_ | Array of players of team B |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>game_config</strong> | _string_ | Confogl config being played |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>campaign_name</strong> | _string_ | Campaign being played |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>pings_requested</strong> | _bool_ | In this state the pug system requires the player to send info about his ping for all gameservers |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>server_ip</strong> | _string_ | Gameserver IP |

<br/><br/>

### GET /getgameservers
##### Get the list of L4D2 servers applicable for the pug system
Request parameters: None

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | Always true |
| <strong>servers</strong> | _[]_ | Array of servers |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>ip</strong> | _string_ | IP without port |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>domain</strong> | _string_ | Domain name pointing to the IP |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>ports</strong> | _[]string_ | Array of ports |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>region</strong> | _string_ | Server region |

<br/><br/>

### GET /pingsreceiver
##### Tell the system about your ping to gameservers
Request parameters:
| Type | Description
| ------ | ------ |
| _map[string]int_ | Map (array with keys) of pings, where key is IP address without port, and value is ping in milliseconds |

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if command accepted, "false" otherwise |
| <strong>error</strong> | _string_ | Outputs the reason if the request is rejected |

<br/><br/>

### GET /sendglobalchat
##### Send a message to the global chat
Request parameters:
| Type | Description
| ------ | ------ |
| <strong>text</strong> | _string_ | The message |

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if command accepted, "false" otherwise |
| <strong>error</strong> | _string_ | Outputs the reason if the request is rejected |

<br/><br/>

### GET /getglobalchat
##### Get all messages of global chat
Request parameters: None

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | Always true |
| <strong>messages</strong> |  | Ordered array of messages |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>time_stamp</strong> | _int64_ | Time the message was created at |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>base64text</strong> | _string_ | Base 64 encoded message text |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>steamid64</strong> | _string_ | Steam ID of a player who sent the message |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>base64name</strong> | _string_ | Base 64 encoded nickname |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>avatar_small</strong> | _string_ | Small avatar |

<br/><br/>

### GET /getbanrecords
##### Get list of bans
Request parameters:
| Type | Description
| ------ | ------ |
| <strong>page</strong> | _int_ | Pagination, starting from 0 |

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | Always true |
| <strong>page</strong> | _int_ | Current page |
| <strong>bans</strong> |  | Ordered array of ban records |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>base64name</strong> | _string_ | Base 64 encoded nickname |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>steamid64</strong> | _string_ | Steam ID 64 |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>created_at</strong> | _int64_ | Banned at |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>accepted_at</strong> | _int64_ | When did the player confirm that he read the ban reason. 0 - if didnt confirm |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>ban_length</strong> | _int64_ | Ban length since the moment of ban confirmation |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>base64reason</strong> | _string_ | Ban reason (base64 encoded) |

<br/><br/>

### POST /ticketcreate
##### Create new ticket
Request parameters:
| Type | Description
| ------ | ------ |
| <strong>ticket_text</strong> | _string_ | Ticket text |
| <strong>redirect_to</strong> | _string_ | Redirect to this page after creating ticket |
| <strong>ticket_type</strong> | _int_ | Ticket type. 1 - behaviour report, 2 - cheater report, 3 - protest ban, 4 - other |

Response parameters: None (303 redirect)

<br/><br/>

### POST /ticketreply
##### Reply to ticket (only for admins)
Request parameters:
| Type | Description
| ------ | ------ |
| <strong>message_text</strong> | _string_ | Message text |
| <strong>redirect_to</strong> | _string_ | Redirect to this page after sending message |
| <strong>ticket_id</strong> | _string_ | Ticket ID |

Response parameters: None (303 redirect)

<br/><br/>

### GET /ticketlist
##### Get my tickets
Request parameters: None

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if command accepted, "false" otherwise |
| <strong>error</strong> | _string_ | Outputs the reason if the request is rejected |
| <strong>opened</strong> |  | Ordered array of opened tickets |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>id</strong> | _string_ | Ticket unique ID |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>type</strong> | _int_ | Ticket type. 1 - behaviour report, 2 - cheater report, 3 - protest ban, 4 - other |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>created_by</strong> | _string_ | Ticket creator |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>created_at</strong> | _int64_ | Time in milliseconds |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>is_closed</strong> | _bool_ | Is ticket closed |
| <strong>closed</strong> |  | Ordered array of last 10 closed tickets (same parameters as in "opened") |
| <strong>as_admin</strong> |  | Ordered array of tickets for reviewing by admins (same parameters as in "opened") |

<br/><br/>

### GET /ticketmessages
##### Get messages of the ticket
Request parameters:
| Type | Description
| ------ | ------ |
| <strong>ticket_id</strong> | _string_ | Ticket ID |

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if command accepted, "false" otherwise |
| <strong>error</strong> | _string_ | Outputs the reason if the request is rejected |
| <strong>ticket_id</strong> | _string_ | Ticket unique ID |
| <strong>messages</strong> |  | Ordered array of messages in the ticket |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>ticket_id</strong> | _string_ | Ticket unique ID |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>by</strong> | _string_ | Message by |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>at</strong> | _int64_ | Message sent at (time in milliseconds) |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>text</strong> | _string_ | Base 64 encoded message text |
