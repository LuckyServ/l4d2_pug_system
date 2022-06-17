
# l4d2_pug_system
## The backend layer of the ranked PUG system for Left 4 Dead 2

### POST /status
##### Get necessary info about program status
Request parameters: none

Response parameters:
| Key | Type | Description
| ------ | ------ | ------ |
| <strong>success</strong> | _bool_ | Always returns "true" |
| <strong>shutdown</strong> | _bool_ | Tells if the program is goind to be shutdown soon (no new lobbies allowed, etc) |
| <strong>time</strong> | _int64_ | System time in milliseconds |

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
| <strong>error</strong> | _string_ | Outputs the reason if the request is rejected |

<br/><br/>

### POST /addauth
##### Add authorization record
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
| <strong>error</strong> | _string_ | Outputs the reason if the authorization got declined |
