package webhooks

import (
	"net/http"

	hobware "github.com/leeola/hob/middleware"
)

func handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := hobware.GetLog(r)
		log.Debug("received webhook")
		http.Error(w, http.StatusText(501), 501)
	}
}
