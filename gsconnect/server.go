package gsconnect

import (
	"fmt"
	"net/http"

	"github.com/lekuruu/ubisoft-game-service/common"
)

type GSConnect struct {
	Host   string
	Port   int
	Logger common.Logger
	Games  map[string]string
}

func (gsc *GSConnect) Serve() {
	bind := fmt.Sprintf("%s:%d", gsc.Host, gsc.Port)
	gsc.Logger.Info(fmt.Sprintf("Listening on %s", bind))
	http.ListenAndServe(bind, nil)
}
