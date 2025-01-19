package arearoute

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"area-backend/area"
	"os"
	"strings"
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
	Id int `json:"userid"`
}

func createActionReaction(a area.AreaRequest, bridge Bridge, userid int) int {
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
		userid,
	).Scan(&bridgeid)
	if err != nil {
		a.Error(err, http.StatusInternalServerError)
		return -1
	}
	return bridgeid
}

func find[T any](vec []T, f func(T) bool) *T {
	for _, v := range vec {
		if f(v) {
			return &v
		}
	}
	return nil
}

func removePerioooodSlaaay(s string) string {
	slen := len(s)
	last := slen - 1
	if slen > 0 && s[last] == '.' {
		return s[:last]
	}
	return s
}

func applyFuncToFirstLetter(s string, f func(string) string) string {
	if len(s) == 0 {
		return s
	}
	return f(string(s[0])) + s[1:]
}


func generateName(a area.AreaRequest, bridge Bridge) string {
	aServ := find(a.Area.About.Server.Services, func(s area.AboutSevice) bool {
		return s.Name == bridge.Action.Service
	})
	rServ := find(a.Area.About.Server.Services, func(s area.AboutSevice) bool {
		return s.Name == bridge.Reaction.Service
	})
	if aServ == nil {
		return "Unknown service for action."
	}
	if rServ == nil {
		return "Unknown service for reaction."
	}
	aDesc := find(aServ.Actions, func(x area.AboutSomething) bool {
		return x.Name == bridge.Action.Name
	})
	rDesc := find(rServ.Reactions, func(x area.AboutSomething) bool {
		return x.Name == bridge.Reaction.Name
	})
	if aServ == nil {
		return "Unknown action name."
	}
	if rServ == nil {
		return "Unknown reaction name."
	}
	return fmt.Sprintf("In %s service, %s, %s on %s",
		applyFuncToFirstLetter(aServ.Name, strings.ToUpper),
		removePerioooodSlaaay(aDesc.Description),
		applyFuncToFirstLetter(
			removePerioooodSlaaay(rDesc.Description),
			strings.ToLower,
		),
		applyFuncToFirstLetter(rServ.Name, strings.ToUpper),
	)
}

func appletUpdate(a area.AreaRequest, bridge Bridge) error {
	querry := `
		insert into areaapplets (name, users, action, reaction)
		values ($1, $2, $3, $4)
		on conflict (action, reaction)
		do update set users = areaapplets.users + 1
	`
	_, err := a.Area.Database.Exec(
		querry,
		generateName(a, bridge),
		1,
		bridge.Action.Service,
		bridge.Reaction.Service,
	)
	if err != nil {
		return err
	}
	return nil
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
	tosend.Bridge = createActionReaction(a, bridge, id)
	if tosend.Bridge == -1 {
		return
	}
	err = appletUpdate(a, bridge)
	if err != nil {
		a.Error(err, http.StatusInternalServerError)
		return
	}
	tosend.Spices = bridge.Action.Spices
	tosend.Id = id
	url := fmt.Sprintf("http://reverse-proxy:%s/service/%s/%s",
		os.Getenv("REVERSEPROXY_PORT"),
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
	body, err := io.ReadAll(rep.Body)
	if err != nil {
		a.Error(err, http.StatusInternalServerError)
		return
	}
	a.Reply(map[string]any{
		"status": "ok",
		"body": string(body),
	}, http.StatusOK)
}