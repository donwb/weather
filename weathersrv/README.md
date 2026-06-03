# weathersrv

Go service that reads the Netatmo weather station and exposes `/current` as JSON
(consumed by the Tidbyt app in `../tidbyt`).

## Endpoints

- `GET /current` — current weather as JSON. Returns `503`/`502` with `{"error": ...}`
  if Netatmo can't be reached or there's no valid token (never crashes the server).
- `GET /` — login page; links to the Netatmo OAuth consent screen.
- `GET /auth_redirect` — OAuth callback; exchanges the code and persists tokens.

## Config (env vars)

| Var            | Purpose                                                            |
|----------------|-------------------------------------------------------------------|
| `HOMEID`       | Netatmo home id                                                   |
| `DEVICEID`     | Netatmo base-station MAC                                           |
| `CLIENTID`     | Netatmo app client id                                             |
| `CLIENTSECRET` | Netatmo app client secret                                         |
| `AUTHREDIRECT` | OAuth redirect URI — must match the Netatmo dev-console exactly    |
| `RTOKEN`       | Refresh token used to bootstrap on first run (see below)          |
| `BTOKEN`       | Optional initial access token                                     |
| `TOKENFILE`    | Where tokens are persisted (default `tokens.json`)                |

Tokens are persisted to `TOKENFILE` and refreshed automatically. `RTOKEN`/`BTOKEN`
only seed the store when there's no token file yet — after that the file is the
source of truth (Netatmo rotates the refresh token on every use).

## Local development

The OAuth redirect flow needs `http://localhost:1323/auth_redirect` registered in
the Netatmo dev console, and it's easy to trip `invalid_grant`. The painless path:

1. dev.netatmo.com → your app → **Token generator** → check `read_station` → Generate.
2. Copy the generated **refresh token** into `env` as `RTOKEN`.
3. `rm -f tokens.json` (so bootstrap re-reads `RTOKEN`), then `make` to run.

`/current` should return data immediately — no redirect flow needed.

If you do use the redirect flow instead: start from `/`, complete consent, and let
it redirect **once**. Never reload `/auth_redirect` — that reuses a consumed,
single-use code and always returns `invalid_grant`.

## Deploy notes

Runs on DigitalOcean App Platform, whose filesystem is **ephemeral**. `tokens.json`
survives in-instance restarts but is wiped on every redeploy, so after a redeploy
you'll need to re-authorize (or re-seed `RTOKEN`). For redeploy-durable auth, move
token storage off local disk (managed Postgres, Spaces/S3, or a KV).
