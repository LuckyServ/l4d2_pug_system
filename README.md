
# l4d2_pug_system
## The backend layer of the ranked PUG system for Left 4 Dead 2

### GET /shutdown
##### An admin command (access = 4) to shutdown the program
Request parameters: None

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if command accepted, "false" otherwise |

<br/><br/>

### GET /status
##### Get necessary info about program status, and signal about online status
Request parameters: None

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | Always true |
| <strong>no_new_lobbies</strong> | _bool_ | Tells if creating new lobbies is blocked |
| <strong>brokenmode</strong> | _bool_ | Tells if competitive plugins are broken by some L4D2 update. In this mode the gameservers are vanilla + Sourcemod. |
| <strong>time</strong> | _int64_ | System time in milliseconds |
| <strong>need_update_players</strong> | _bool_ | Should update players or not |
| <strong>need_update_lobbies</strong> | _bool_ | Should update lobbies or not |
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

### GET /acceptrules
##### Accept rules of the project
Request parameters: None

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if rules accepted, "false" otherwise |
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
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>mmr</strong> | _int_ | Player's rating |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>mmr_certain</strong> | _bool_ | Is the system certain about the player's rating |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>access</strong> | _int_ | Player's access level<br>-2 - completely banned, -1 - chat banned, 0 - regular player, 1 - behaviour moderator, 2 - cheat moderator, 3 - behaviour+cheat moderator, 4 - full admin access |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>is_ingame</strong> | _bool_ | Is player in game right now |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>is_inlobby</strong> | _bool_ | Is player in lobby right now |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>is_idle</strong> | _bool_ | Is player online, but not doing anything |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>profile_validated</strong> | _bool_ | New players must validate their profiles before playing |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>rules_accepted</strong> | _bool_ | New players must accept the rules before playing |
| <strong>count</strong> |  | Numbers |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>online</strong> | _int_ | Number of online players |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>in_game</strong> | _int_ | Number of players in games |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>in_lobby</strong> | _int_ | Number of players in lobbies |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>idle</strong> | _int_ | Number of idle players |
| <strong>list</strong> | _[]_ | Array of players |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>steamid64</strong> | _string_ | Player Steam ID 64 |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>nickname_base64</strong> | _string_ | Base64 encoded nickname |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>mmr</strong> | _int_ | Player's rating |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>mmr_certain</strong> | _bool_ | Is the system certain about the player's rating |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>access</strong> | _int_ | Player's access level<br>-2 - completely banned, -1 - chat banned, 0 - regular player, 1 - behaviour moderator, 2 - cheat moderator, 3 - behaviour+cheat moderator, 4 - full admin access |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>is_ingame</strong> | _bool_ | Is player in game right now |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>is_inlobby</strong> | _bool_ | Is player in lobby right now |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>is_idle</strong> | _bool_ | Is player online, but not doing anything |

<br/><br/>

### GET /joinanylobby
##### Join any lobby, or create new one, if can't join
Request parameters: None

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if joined/created the lobby, "false" otherwise |
| <strong>error</strong> | _string_ | Outputs the reason if the request is rejected |

<br/><br/>

### GET /createlobby
##### Create lobby and join it
Request parameters: None

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if lobby created, "false" otherwise |
| <strong>error</strong> | _string_ | Outputs the reason if the request is rejected |

<br/><br/>

### GET /joinlobby
##### Join a specific lobby
Request parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>lobby_id</strong> | _string_ | ID of the lobby you want to join |

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if joined the lobby, "false" otherwise |
| <strong>error</strong> | _string_ | Outputs the reason if the request is rejected |

<br/><br/>

### GET /leavelobby
##### Leave lobby
Request parameters: None

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if left from lobby, "false" otherwise |
| <strong>error</strong> | _string_ | Outputs the reason if the request is rejected |

<br/><br/>

### GET /readyup
##### Ready up, when your lobby is in ReadyUp state
Request parameters: None

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if command accepted, "false" otherwise |
| <strong>error</strong> | _string_ | Outputs the reason if the request is rejected |

<br/><br/>

### GET /getlobbies
##### Get lobby list
Request parameters: None

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | Always "true" |
| <strong>authorized</strong> | _bool_ | Authorized or not |
| <strong>is_inlobby</strong> | _bool_ | If authorized, is player in lobby or not |
| <strong>count</strong> | _int_ | Number of lobbies |
| <strong>need_readyup</strong> | _bool_ | Should request readyup or not |
| <strong>mylobby</strong> |  | The lobby the player participates in (only present if authorized and in lobby) |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>id</strong> | _string_ | ID |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>mmr_min</strong> | _int_ | Lowest allowed mmr, -2000000000 if unbounded |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>mmr_max</strong> | _int_ | Highest allowed mmr, 2000000000 if unbounded |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>created_at</strong> | _int64_ | Time of the creation (unix timestamp in milliseconds) |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>confogl_config</strong> | _string_ | Confogl config to be played |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>player_count</strong> | _int_ | Number of players in the lobby |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>readyup_state</strong> | _bool_ | Is lobby in readyup state |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>ready_players</strong> | _int_ | Number of ready players |
| <strong>lobbies</strong> | _[]_ | Array of lobbies |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>id</strong> | _string_ | ID |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>mmr_min</strong> | _int_ | Lowest allowed mmr, -2000000000 if unbounded |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>mmr_max</strong> | _int_ | Highest allowed mmr, 2000000000 if unbounded |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>created_at</strong> | _int64_ | Time of the creation (unix timestamp in milliseconds) |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>confogl_config</strong> | _string_ | Confogl config to be played |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>player_count</strong> | _int_ | Number of players in the lobby |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>readyup_state</strong> | _bool_ | Is lobby in readyup state |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>ready_players</strong> | _int_ | Number of ready players |

<br/><br/>

### GET /getgame
##### Get the current game info of an authorized player
Request parameters: None

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if success, "false" if not authorized or not in game |
| <strong>game</strong> |  | Info about game, if success == true |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>players_a</strong> | _[]_ | Array of players of team A |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>players_b</strong> | _[]_ | Array of players of team B |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>game_config</strong> | _string_ | Confogl config being played |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>campaign_name</strong> | _string_ | Campaign being played |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>pings_requested</strong> | _bool_ | In this state the pug system requires the player to send info about his ping for all gameservers |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>server_ip</strong> | _string_ | Gameserver IP |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>mmr_min</strong> | _int_ | Lowest allowed mmr, -2000000000 if unbounded |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>mmr_max</strong> | _int_ | Highest allowed mmr, 2000000000 if unbounded |

<br/><br/>

### GET /getgameservers
##### Get the list of L4D2 servers applicable for the l4d2_pug_system
Request parameters: None

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | Always true |
| <strong>gameservers</strong> | _[]string_ | Array of IP:PORT |
| <strong>servers</strong> | _[]string_ | Array of IP |

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
| <strong>success</strong> | _bool_ | Always true |
| <strong>gameservers</strong> | _[]string_ | Array of IP:PORT |
| <strong>servers</strong> | _[]string_ | Array of IP |

<br/><br/><br/><br/><br/><br/>
## API for gameservers
<br>All responses are VDF (Valve Data Format)
<br>All responses are headed with "VDFresponse" key
<br>All request and response parameters are strings
<br>

### POST /gs/getgame
##### Get the current game info on the server
Request parameters:
| Key | Description
| ------ | ------ |
| <strong>auth_key</strong> | Backend auth key |
| <strong>ip</strong> | IP address of the server (with port) |

Response parameters:
| Key | Description
| ------ | ------ |
| <strong>success</strong> | "true" if game found on this IP |
| <strong>error</strong> | Error text if success == false |
| <strong>player_a</strong>N | Steam ID 64 of player N in team A (N is value from 0 to 3) |
| <strong>player_b</strong>N | Steam ID 64 of player N in team B (N is value from 0 to 3) |
| <strong>confogl</strong> | Confogl config |
| <strong>first_map</strong> | First map of the campaign |
| <strong>last_map</strong> | Last map of the campaign |
| <strong>mmr_min</strong> | Minimum mmr allowed in this game |
| <strong>mmr_max</strong> | Maximum mmr allowed in this game |
| <strong>game_state</strong> | Game state, possible values are "wait_readyup", "other" |
