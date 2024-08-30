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

type GSContext struct {
	GSConnect *GSConnect
	Response  http.ResponseWriter
	Request   *http.Request
}

func (gsc *GSConnect) Serve() {
	bind := fmt.Sprintf("%s:%d", gsc.Host, gsc.Port)
	gsc.Logger.Info(fmt.Sprintf("Listening on %s", bind))

	http.HandleFunc("/gsinit.php", gsc.withGSContext(GSInitRoute))
	http.ListenAndServe(bind, nil)
}

// Add "GSContext" for the handler function
func (gsc *GSConnect) withGSContext(handler func(*GSContext)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(&GSContext{
			GSConnect: gsc,
			Response:  w,
			Request:   r,
		})
	}
}
