package gsconnect

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

func GSInitRoute(gsc *GSContext) {
	query := gsc.Request.URL.Query()
	user := query.Get("user")
	product := query.Get("dp")

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

	// We need to hijack the connection to force-close the
	// connection after we send the response
	conn, buf, err := hijackConnection(gsc.Response)

	if err != nil {
		// Use the default response writer if hijacking fails
		gsc.Response.Header().Set("Content-Type", "text/plain")
		gsc.Response.WriteHeader(http.StatusOK)
		gsc.Response.Write([]byte(game))
		return
	}

	// Send the game to the client
	buf.Write([]byte("HTTP/1.1 200 OK\r\n"))
	buf.Write([]byte("Content-Type: text/plain\r\n\r\n"))
	buf.Write([]byte(game))
	buf.Flush()
	conn.Close()
}

func hijackConnection(response http.ResponseWriter) (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := response.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("hijacking not supported")
	}

	conn, buf, err := hj.Hijack()
	if err != nil {
		return nil, nil, err
	}

	return conn, buf, nil
}
