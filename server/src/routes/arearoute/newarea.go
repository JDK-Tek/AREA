package arearoute

import (
	"area-backend/area"
	"net/http"
)

func NewArea(a area.AreaRequest) {
	var email, err = a.AssertToken()

	if err != nil {
		a.Error(err, http.StatusBadRequest)
		return
	}
	a.Reply("coucou " + email, http.StatusOK)
}