module github.com/lekuruu/ubisoft-game-service

go 1.22.6

require (
	github.com/lekuruu/ubisoft-game-service/cdkey v0.0.0-20240831105814-85a1b7e9b455
	github.com/lekuruu/ubisoft-game-service/common v0.0.0-20240831105814-85a1b7e9b455
	github.com/lekuruu/ubisoft-game-service/gsconnect v0.0.0-20240831105814-85a1b7e9b455
	github.com/lekuruu/ubisoft-game-service/gsnat v0.0.0-20240831105814-85a1b7e9b455
	github.com/lekuruu/ubisoft-game-service/router v0.0.0-20240831105814-85a1b7e9b455
)

require golang.org/x/exp v0.0.0-20240823005443-9b4947da3948 // indirect

replace github.com/lekuruu/ubisoft-game-service/cdkey => ./cdkey

replace github.com/lekuruu/ubisoft-game-service/common => ./common

replace github.com/lekuruu/ubisoft-game-service/gsconnect => ./gsconnect

replace github.com/lekuruu/ubisoft-game-service/gsnat => ./gsnat

replace github.com/lekuruu/ubisoft-game-service/router => ./router

replace github.com/lekuruu/ubisoft-game-service/irc => ./irc
