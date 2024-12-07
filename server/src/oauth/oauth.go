package oauth

import (
	"area-backend/area"
	"net/http"

	"github.com/gorilla/mux"
)

func discord(a area.AreaRequest) {
	
}

func Router(a area.AreaRequest) {
	v := mux.Vars(a.Request)
	switch (v["redirection"]) {
	case "discord":
		discord(a)
	default:
		a.Reply(map[string]string{
			"error": v["redirection"] + " not found.",
		}, http.StatusNotFound)
	}
}