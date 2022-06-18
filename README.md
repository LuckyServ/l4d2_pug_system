
# l4d2_pug_system
## The backend layer of the ranked PUG system for Left 4 Dead 2

### POST /status
##### Get necessary info about program status
Request parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>backend_auth</strong> | _string_ | Auth key |
| <strong>steamid64</strong> (optional) | _string_ | Steam ID 64 (Profile ID) of the authorized player |

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | True if request succeeded, false otherwise |
| <strong>shutdown</strong> | _bool_ | Tells if the program is goind to be shutdown soon (no new lobbies allowed, etc) |
| <strong>time</strong> | _int64_ | System time in milliseconds |
| <strong>players_updated</strong> | _int64_ | Last time players list was updated (in milliseconds) |
| <strong>error</strong> | _int_ | Outputs the reason if the request is rejected.<br>1 - bad auth |

<br/><br/>

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

### POST /addauth
##### Add player authorization record
Request parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>backend_auth</strong> | _string_ | Auth key |
| <strong>steamid64</strong> | _string_ | Steam ID 64 (Profile ID) |
| <strong>nickname_base64</strong> | _string_ | Base64 encoded nickname |

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if authorization added, "false" otherwise |
| <strong>session_id</strong> | _string_ | Returns the session id if the authorization got accepted |
| <strong>error</strong> | _int_ | Outputs the reason if the authorization got declined.<br>1 - bad auth, 2 - bad parameters |

<br/><br/>

### POST /removeauth
##### Remove player authorization record
Request parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>backend_auth</strong> | _string_ | Auth key |
| <strong>session_id</strong> | _string_ | Session ID which should be removed |

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if authorization removed, "false" otherwise |
| <strong>error</strong> | _int_ | Outputs the reason if the operation fails.<br>1 - bad auth, 2 - bad parameters, 3 - no such auth key |

<br/><br/>

### POST /updateactivity
##### Update player's last activity. This is needed to keep record of online players.
Request parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>backend_auth</strong> | _string_ | Auth key |
| <strong>steamid64</strong> | _string_ | Steam ID 64 (Profile ID) |

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if activity timestamp updated, "false" otherwise |
| <strong>error</strong> | _int_ | Outputs the reason if the operation fails.<br>1 - bad auth, 2 - bad parameters, 3 - no such player |

<br/><br/>

### POST /getplayer
##### Get info about a player
Request parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>backend_auth</strong> | _string_ | Auth key |
| <strong>steamid64</strong> | _string_ | Steam ID 64 (Profile ID) |

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | "true" if info available, "false" otherwise |
| <strong>error</strong> | _int_ | Outputs the reason if the operation fails.<br>1 - bad auth, 2 - bad parameters, 3 - no such player |
| <strong>player</strong> | _array_ | Array of the player's info |
| <strong>player["SteamID64"]</strong> | _string_ | Steam ID 64 (Profile ID) |
| <strong>player["NicknameBase64"]</strong> | _string_ | Base64 encoded nickname |
| <strong>player["Mmr"]</strong> | _int_ | Player's rating |
| <strong>player["MmrUncertainty"]</strong> | _int_ | How uncertain the system is about the player's rating |
| <strong>player["Access"]</strong> | _int_ | Player's access level. -2 - completely banned, -1 - chat banned, 0 - regular player, 1 - behaviour moderator, 2 - cheat moderator, 3 - behaviour+cheat moderator, 4 - full admin access |
| <strong>player["ProfValidated"]</strong> | _bool_ | New players must validate their profiles before playing |
| <strong>player["Pings"]</strong> | _array_ | Array of player's pings to every gameserver. Used to choose the best gameserver |
| <strong>player["Pings"]["</strong>127.0.0.1<strong>"]</strong> | _int_ | The ping |
| <strong>player["PingsUpdated"]</strong> | _int64_ | Last time player's pings were updated (in milliseconds) |
| <strong>player["LastActivity"]</strong> | _int64_ | Last time player showed any activity |
| <strong>player["IsOnline"]</strong> | _bool_ | Is player online right now |
| <strong>player["IsInGame"]</strong> | _bool_ | Is player in game right now |
| <strong>player["IsInLobby"]</strong> | _bool_ | Is player in lobby right now |
| <strong>player["LastUpdated"]</strong> | _int64_ | Last time player's info was changed (except LastActivity) (in milliseconds) |
