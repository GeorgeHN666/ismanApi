SHELL=cmd
API_PORT=8080 ## production mode port = 443
ENV=development
## build: builds all binaries
build: clean build_front build_back
	@echo All binaries built!

## clean: cleans all binaries and runs go clean
clean:
	@echo Cleaning...
	@echo y | DEL /S dist
	@go clean
	@echo Cleaned and deleted binaries
## build_api: builds the API
build_api:
	@echo Building API...
	@go build -o ./dist/api.exe ./cmd/
	@echo API built!
## start: starts api
start: start_api

## start_api: start the api
start_api: build_api
	@echo Starting API...
	start .\dist\api.exe 
	@echo °API running in background

## stop: stops the API
stop: stop_api
	@echo All applications stopped

## stop_api: stops the API
stop_api:
	@echo Stopping API...
	@taskkill /IM api.exe /F
	@echo Stopped API