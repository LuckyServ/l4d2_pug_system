
# l4d2_pug_system
## The backend layer of the ranked PUG system for Left 4 Dead 2

### POST /shutdown
##### Send the shutdown command to the program. It will properly exit the process.
Request parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>backend_auth</strong> | _string_ | Auth key |

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
| <strong>need_update_player</strong> | _bool_ | Should update player or not (only present if authorized) |

<br/><br/>

### GET /getme
##### Get info about an authorized player
Request parameters: None

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if info available, "false" otherwise (not authorized) |
| <strong>steamid64</strong> | _string_ | Steam ID 64 |
| <strong>nickname_base64</strong> | _string_ | Base64 encoded nickname |
| <strong>mmr</strong> | _int_ | Player's rating |
| <strong>mmr_certain</strong> | _bool_ | Is the system certain about the player's rating |
| <strong>access</strong> | _int_ | Player's access level<br>-2 - completely banned, -1 - chat banned, 0 - regular player, 1 - behaviour moderator, 2 - cheat moderator, 3 - behaviour+cheat moderator, 4 - full admin access |
| <strong>profile_validated</strong> | _bool_ | New players must validate their profiles before playing |
| <strong>rules_accepted</strong> | _bool_ | New players must accept the rules before playing |
| <strong>is_online</strong> | _bool_ | Is player online right now |
| <strong>is_ingame</strong> | _bool_ | Is player in game right now |
| <strong>is_inlobby</strong> | _bool_ | Is player in lobby right now |

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
| <strong>count</strong> |  | Numbers |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>online</strong> | _int_ | Number of online players |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>in_game</strong> | _int_ | Number of players in games |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>in_lobby</strong> | _int_ | Number of players in lobbies |
| <strong>list</strong> | _[]_ | Array of players |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>steamid64</strong> | _string_ | Player Steam ID 64 |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>nickname_base64</strong> | _string_ | Base64 encoded nickname |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>mmr</strong> | _int_ | Player's rating |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>mmr_certain</strong> | _bool_ | Is the system certain about the player's rating |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>access</strong> | _int_ | Player's access level<br>-2 - completely banned, -1 - chat banned, 0 - regular player, 1 - behaviour moderator, 2 - cheat moderator, 3 - behaviour+cheat moderator, 4 - full admin access |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>is_ingame</strong> | _bool_ | Is player in game right now |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>is_inlobby</strong> | _bool_ | Is player in lobby right now |

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

### GET /getlobbies
##### Get lobby list
Request parameters: None

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | Always "true" |
| <strong>count</strong> | _int_ | Number of lobbies |
| <strong>mylobby</strong> |  | The lobby the player participates in (only present if authorized and in lobby) |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>id</strong> | _string_ | ID |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>mmr_min</strong> | _int_ | Lowest allowed mmr, -2000000000 if unbounded |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>mmr_max</strong> | _int_ | Highest allowed mmr, 2000000000 if unbounded |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>created_at</strong> | _int64_ | Time of the creation (unix timestamp in milliseconds) |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>confogl_config</strong> | _string_ | Confogl config to be played |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>player_count</strong> | _int_ | Number of players in the lobby |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>readyup_state</strong> | _bool_ | Is lobby in readyup state |
| <strong>lobbies</strong> | _[]_ | Array of lobbies |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>id</strong> | _string_ | ID |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>mmr_min</strong> | _int_ | Lowest allowed mmr, -2000000000 if unbounded |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>mmr_max</strong> | _int_ | Highest allowed mmr, 2000000000 if unbounded |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>created_at</strong> | _int64_ | Time of the creation (unix timestamp in milliseconds) |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>confogl_config</strong> | _string_ | Confogl config to be played |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>player_count</strong> | _int_ | Number of players in the lobby |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>readyup_state</strong> | _bool_ | Is lobby in readyup state |
