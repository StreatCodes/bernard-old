# Bernard

Easily monitor your infrastructure.

## TODO
 - Extend Server
   - Attach a bbolt DB
   - Create http server
   - create auth endpoint
   - create websocket
   - send all info on websocket connect
   - support creating/updating hosts, users
   - add roles to users

  - Frontend
    - Render hosts,
    - render checks
    - more info on hosts
    - more info on checks
    - login/auth
    - user management

 - Clients should connect to parent nodes using TLS
 - Clients should be able to configure a public key associated with a self signed cert when connecting to a parent node
 - Client should attempt to reconnect after disconnecting from parent node
 - Document a check