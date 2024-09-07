module github.com/lekuruu/ubisoft-game-service

go 1.22.6

require (
	github.com/lekuruu/ubisoft-game-service/cdkey v0.0.0-20240831105814-85a1b7e9b455
	github.com/lekuruu/ubisoft-game-service/common v0.0.0-20240907193506-02b049cca13f
	github.com/lekuruu/ubisoft-game-service/gsconnect v0.0.0-20240831105814-85a1b7e9b455
	github.com/lekuruu/ubisoft-game-service/gsnat v0.0.0-20240831105814-85a1b7e9b455
	github.com/lekuruu/ubisoft-game-service/irc v0.0.0-00010101000000-000000000000
	github.com/lekuruu/ubisoft-game-service/router v0.0.0-20240831105814-85a1b7e9b455
)

require (
	github.com/lekuruu/ubisoft-game-service/proxy v0.0.0-00010101000000-000000000000
	golang.org/x/exp v0.0.0-20240823005443-9b4947da3948 // indirect
	golang.org/x/time v0.0.0-20220722155302-e5dcc9cfc0b9 // indirect
	gopkg.in/irc.v4 v4.0.0 // indirect
)

replace github.com/lekuruu/ubisoft-game-service/cdkey => ./cdkey

replace github.com/lekuruu/ubisoft-game-service/common => ./common

replace github.com/lekuruu/ubisoft-game-service/gsconnect => ./gsconnect

replace github.com/lekuruu/ubisoft-game-service/gsnat => ./gsnat

replace github.com/lekuruu/ubisoft-game-service/router => ./router

replace github.com/lekuruu/ubisoft-game-service/irc => ./irc

replace github.com/lekuruu/ubisoft-game-service/proxy => ./proxy
