#
# The include should be a single file that contains:
# export APIKEY := {APIKEY}
# export SECRET := {SECRET}
#
# LOCAL DEV: easiest auth is the Netatmo dev-console "Token generator"
# (check read_station -> Generate). Put the refresh token in env as RTOKEN,
# `rm -f tokens.json`, then `make`. Avoids the OAuth redirect / invalid_grant.
# See README.md for details.
#
include env

$(info $$DEVICEID is [${DEVICEID}])
$(info $$HOMEID is [${HOMEID}])
$(info $$BTOKEN is [${BTOKEN}])
$(info $$RTOKEN is [${RTOKEN}])
$(info $$CLIENTID is [${CLIENTID}])
$(info $$CLIENTSECRET is [${CLIENTSECRET}])

all:
	go run *.go

# Local dev: clear any stale persisted token so the server re-seeds from
# RTOKEN in the env file, then run. Use after pasting a fresh refresh token
# from the Netatmo dev-console Token generator.
local:
	@if [ -z "$(RTOKEN)" ]; then echo "RTOKEN is not set in env — paste a refresh token there first"; exit 1; fi
	rm -f tokens.json tokens.json.tmp
	go run *.go

build: ## Build
	go build *.go

test: 
	gow -e=go,mod,html,js,css run .
	
run:
	gin

docker-buildM1:
	docker build --platform linux/amd64 -t donwb/weathersrv:0.5 .


docker-run:
	docker run --platform linux/amd64 --env-file=docker-env -p 80:1323 -it donwb/weathersrv:0.5

.DEFAULT_GOAL := all