# Webhook to Matrix Bot

![DevGuard Badge](https://api.main.devguard.org/api/v1/badges/cvss/05bcc0c3-98fc-4d7a-a438-0d1e404f62c3)

This little Golang bot/server listens to a webhook and sends the received data to a Matrix room.

![](https://matrix.org/images/matrix-logo-white.svg)

## How to setup:

- Clone repo: `git clone https://github.com/l3montree-dev/webhook-to-matrix.git`
- Create Matrix User if you don't have one already - e.g. using [Element](https://element.io/):
  - Logout of existing account if necessary and create a "normal" new account
  - Login to the new account and generate a recovery key (e.g. start a chat with another user -> element will ask to store the key)
  - Generate access token for user
    - `curl -X POST -H 'Content-Type: application/json' -d '{ "type":"m.login.password", "user":"username", "password":"password" }' "https://matrix.org/_matrix/client/r0/login"`
- Setup config file (`cp .env.example .env`) and adjust the variables in the `.env` file
- Run the server
  - Directly via `Go`: `go run main.go`
  - Docker: `docker build -t webhook-to-matrix . && docker run --rm -p 5001:5001 -v $(pwd)/.env:/app/.env:ro webhook-to-matrix`
  - Helm: TBD
- Setup apps (e.g. Glitchtip / Botkube / ...) to send data to the bot. E.g.:
  - `http://your-domain.com/webhook/my-webhook-secret/glitchtip?roomid=xyz` [Docs](https://glitchtip.com/documentation/error-tracking#turn-on-alerts)
  - `http://your-domain.com/webhook/my-webhook-secret/botkube?roomid=xyz` [Docs](https://docs.botkube.io/installation/webhook/)

## Architecture 

```mermaid
flowchart LR;
    botkube
    glitchtip
    Matrix

    subgraph Webhook to Matrix Bot
        WTM-HTTP(HTTP Endpoint)
        WTM-Delivery(Delivery Agent)

        WTM-HTTP -. transform message .-> WTM-Delivery
    end

    botkube -. send via webhook .-> WTM-HTTP
    glitchtip -. send via webhook .-> WTM-HTTP
    WTM-Delivery -. send to matrix .-> Matrix
```

## Supported Applications

<a href="https://glitchtip.com/"><img src="https://glitchtip.com/assets/logo-again.svg" width="200px"></a>
<a href="https://botkube.io/"><img src="https://github.com/kubeshop/botkube/raw/main/docs/assets/botkube-title.png" width="200px"></a>
<a href="https://matrix.org/"><img src="https://matrix.org/images/matrix-logo.svg" width="200px"></a>

