
# l4d2_pug_system
## The backend layer of the ranked PUG system for Left 4 Dead 2

### POST /shutdown
##### Send the shutdown command to the program. It will wait until all lobbies end, and then exit the process.
Request parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>backend_auth</strong> | _string_ | Auth key |

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if command accepted, "false" otherwise |
| <strong>error</strong> | _int_ | Outputs the reason if the request is rejected.<br>1 - bad auth, 2 - already shutting down |

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
| <strong>mmr_uncertainty</strong> | _float32_ | How uncertain the system is about the player's rating |
| <strong>mmr_certain</strong> | _bool_ | Is the system certain about the player's rating |
| <strong>access</strong> | _int_ | Player's access level. -2 - completely banned, -1 - chat banned, 0 - regular player, 1 - behaviour moderator, 2 - cheat moderator, 3 - behaviour+cheat moderator, 4 - full admin access |
| <strong>profile_validated</strong> | _bool_ | New players must validate their profiles before playing |
| <strong>rules_accepted</strong> | _bool_ | New players must accept the rules before playing |
| <strong>is_online</strong> | _bool_ | Is player online right now |
| <strong>is_ingame</strong> | _bool_ | Is player in game right now |
| <strong>is_inlobby</strong> | _bool_ | Is player in lobby right now |
