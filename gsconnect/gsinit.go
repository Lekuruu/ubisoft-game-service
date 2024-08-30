package gsconnect

import (
	"fmt"
)

func GSInitRoute(gsc *GSContext) {
	fmt.Fprintf(gsc.Response, "GSInitRoute") // TODO
}
