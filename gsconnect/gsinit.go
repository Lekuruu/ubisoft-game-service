package gsconnect

import (
	"fmt"
	"net/http"
)

func GSInitRoute(gsc *GSContext) {
	query := gsc.Request.URL.Query()
	user := query.Get("user")
	product := query.Get("dp")
	gsc.Response.Header().Add("Content-Type", "text/plain")

	if product == "" {
		gsc.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	game, ok := gsc.Server.Games[product]

	if !ok {
		gsc.Response.WriteHeader(http.StatusNotFound)
		return
	}

	log := fmt.Sprintf("'%s' is connecting to '%s'", user, product)
	gsc.Server.Logger.Info(log)

	gsc.Response.WriteHeader(http.StatusOK)
	gsc.Response.Write([]byte(game))
}
