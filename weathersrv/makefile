#
# The include should be a single file that contains:
# export APIKEY := {APIKEY}
# export SECRET := {SECRET}
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

build: ## Build
	go build *.go

test: 
	gow -e=go,mod,html,js,css run .
	
run:
	gin

docker-buildM1:
	docker build --platform linux/amd64 -t donwb/weathersrv:0.1 .


docker-build:
	docker build -t donwb/videowall:0.1 .

docker-run:
	docker run --env-file=docker-env -p 80:1323 -it donwb/weathersrv:0.1

.DEFAULT_GOAL := all