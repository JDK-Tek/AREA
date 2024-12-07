package arearoute

import (
	"area-backend/area"
	"encoding/json"
	"net/http"
	"fmt"
)

type AreaObject struct {
	Service string `json:"service"`
	Name string `json:"name"`
	Spices json.RawMessage `json:"spices"`
}

type Bridge struct {
	Action AreaObject `json:"action"`
	Reaction AreaObject `json:"reaction"`
}

func createActionReaction(a area.AreaRequest, bridge Bridge) int {
	var actionid, reactionid, bridgeid int
	querry := `select id from actions where service = $1
		and name = $2 and spices = $3`
	err := a.Area.Database.QueryRow(querry,
			bridge.Action.Name,
			bridge.Action.Service,
			bridge.Action.Spices,
		).Scan(&actionid)
	if err != nil {
		querry = `insert into actions (service, name, spices)
		values ($1, $2, $3) returning id`
		err = a.Area.Database.QueryRow(querry,
				bridge.Action.Name,
				bridge.Action.Service,
				bridge.Action.Spices,
			).Scan(&actionid)
	}
	if err != nil {
		a.Error(err, http.StatusInternalServerError)
		return -1
	}
	querry = `select id from reactions where service = $1
	and name = $2 and spices = $3`
	err = a.Area.Database.QueryRow(querry,
			bridge.Reaction.Name,
			bridge.Reaction.Service,
			bridge.Reaction.Spices,
		).Scan(&reactionid)
	if err != nil {
		querry = `insert into reactions (service, name, spices)
		values ($1, $2, $3) returning id`
		err = a.Area.Database.QueryRow(querry,
				bridge.Reaction.Name,
				bridge.Reaction.Service,
				bridge.Reaction.Spices,
			).Scan(&reactionid)
	}
	if err != nil {
		a.Error(err, http.StatusInternalServerError)
		return -1
	}
	querry = `insert into bridge (action, reaction, userid)
	values ($1, $2, $3) returning id`
	err = a.Area.Database.QueryRow(querry,
		actionid,
		reactionid,
		42,
	).Scan(&bridgeid)
	if err != nil {
		a.Error(err, http.StatusInternalServerError)
		return -1
	}
	return bridgeid
}

func NewArea(a area.AreaRequest) {
	var email, err = a.AssertToken()
	var bridge Bridge

	if err != nil {
		a.Error(err, http.StatusBadRequest)
		return
	}
	err = json.NewDecoder(a.Request.Body).Decode(&bridge)
	if err != nil {
		a.Error(err, http.StatusBadRequest)
		return
	}
	n := createActionReaction(a, bridge)
	if n == -1 {
		return
	}
	fmt.Fprintf(a.Writter, "Your email is %s, bridge of action is %d", email, n)
}