module github.com/lekuruu/ubisoft-game-service/cdkey

go 1.22.6

require (
	github.com/lekuruu/ubisoft-game-service/common v0.0.0-20240831105814-85a1b7e9b455
	github.com/lekuruu/ubisoft-game-service/router v0.0.0-20240831105814-85a1b7e9b455
)

require (
	golang.org/x/exp v0.0.0-20240823005443-9b4947da3948 // indirect
	golang.org/x/time v0.0.0-20220722155302-e5dcc9cfc0b9 // indirect
	gopkg.in/irc.v4 v4.0.0 // indirect
)

replace github.com/lekuruu/ubisoft-game-service/common => ../common

replace github.com/lekuruu/ubisoft-game-service/router => ../router
