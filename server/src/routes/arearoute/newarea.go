package arearoute

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"area-backend/area"
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

type ToSend struct {
	Spices json.RawMessage `json:"spices"`
	Bridge int `json:"bridge"`
}

func createActionReaction(a area.AreaRequest, bridge Bridge) int {
	var actionid, reactionid, bridgeid int
	querry := `select id from actions where service = $1
		and name = $2 and spices = $3`
	err := a.Area.Database.QueryRow(querry,
		bridge.Action.Service,
		bridge.Action.Name,
		bridge.Action.Spices,
		).Scan(&actionid)
	if err != nil {
		querry = `insert into actions (service, name, spices)
		values ($1, $2, $3) returning id`
		err = a.Area.Database.QueryRow(querry,
			bridge.Action.Service,
			bridge.Action.Name,
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
		bridge.Reaction.Service,
		bridge.Reaction.Name,
		bridge.Reaction.Spices,
		).Scan(&reactionid)
	if err != nil {
		querry = `insert into reactions (service, name, spices)
		values ($1, $2, $3) returning id`
		err = a.Area.Database.QueryRow(querry,
			bridge.Reaction.Service,
			bridge.Reaction.Name,
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
	var id, err = a.AssertToken()
	var bridge Bridge
	var tosend ToSend

	if err != nil {
		a.Error(err, http.StatusBadRequest)
		return
	}
	err = json.NewDecoder(a.Request.Body).Decode(&bridge)
	if err != nil {
		a.Error(err, http.StatusBadRequest)
		return
	}
	tosend.Bridge = createActionReaction(a, bridge)
	if tosend.Bridge == -1 {
		return
	}
	tosend.Spices = bridge.Action.Spices
	url := fmt.Sprintf("http://reverse-proxy:42002/service/%s/%s",
		bridge.Action.Service,
		bridge.Action.Name)
	obj, err := json.Marshal(tosend)
	if err != nil {
		a.Error(err, http.StatusInternalServerError)
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(obj))
	if err != nil {
		a.Error(err, http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := http.Client{}
	rep, err := client.Do(req)
	if err != nil {
		a.Error(err, http.StatusBadGateway)
		return
	}
	defer rep.Body.Close()
	body, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		a.Error(err, http.StatusInternalServerError)
		return
	}
	a.Reply(map[string]any{
		"res": "Your email is send" + string(body) + " " + string(id),
	}, http.StatusOK)
	// fmt.Fprintf(a.Writter, "Your email is %d, Awnser is %s", id, string(body))
}