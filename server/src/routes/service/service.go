package service

import (
	"area-backend/area"
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	"strconv"
)

func GetServices(a area.AreaRequest) {
	rows, err := a.Area.Database.Query("SELECT id, name, logo, link, colorN, colorH FROM services")
	if err != nil {
		a.Error(err, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var services []map[string]any

	for rows.Next() {
		var id int
		var name, logo, link, colorN, colorH string

		if err := rows.Scan(&id, &name, &logo, &link, &colorN, &colorH); err != nil {
			a.Error(err, http.StatusInternalServerError)
			return
		}

		services = append(services, map[string]any{
			"id":   id,
			"name": name,
			"logo": logo,
			"link": link,
			"color": map[string]any{
				"normal": colorN,
				"hover":  colorH,
			},
		})
	}

	if err := rows.Err(); err != nil {
		a.Error(err, http.StatusInternalServerError)
		return
	}

	a.Reply(map[string]any{
		"res": services,
	}, http.StatusOK)
}

func GetServiceID(a area.AreaRequest) (int, bool) {
	serviceIdStr := mux.Vars(a.Request)["id"]
	if serviceIdStr == "" {
		a.Error(fmt.Errorf("missing service ID"), http.StatusBadRequest)
		return (-1), false
	}

	serviceId, err := strconv.Atoi(serviceIdStr)
	if err != nil {
		a.Error(fmt.Errorf("invalid service ID: %v", err), http.StatusBadRequest)
		return (-1), false
	}
	return serviceId, true
}

func GetServiceApplets(a area.AreaRequest) {
	var serviceId int
	var convStatus bool

	serviceId, convStatus = GetServiceID(a)

	if !convStatus {
		a.Error(fmt.Errorf("invalid service ID"), http.StatusBadRequest)
		return
	}

	rowsApplet, err := a.Area.Database.Query("SELECT id, name, link, users, service1, service2 FROM applets WHERE service1 = $1 OR service2 = $1", serviceId)
	if err != nil {
		a.Error(err, http.StatusInternalServerError)
		return
	}

	var applets []map[string]any

	for rowsApplet.Next() {
		var appletName, appletLink string
		var appletId, appletUsers, appletService1, appletService2 int
		if err := rowsApplet.Scan(&appletId, &appletName, &appletLink, &appletUsers, &appletService1, &appletService2); err != nil {
			a.Error(err, http.StatusInternalServerError)
			return
		}

		var service1Name, service1Logo, service1ColorN, service1ColorH string
		rowsService1, err := a.Area.Database.Query("SELECT name, logo, colorN, colorH FROM services WHERE id = $1", appletService1)
		if err != nil {
			a.Error(err, http.StatusInternalServerError)
			return
		}
		rowsService1.Next()
		if err := rowsService1.Scan(&service1Name, &service1Logo, &service1ColorN, &service1ColorH); err != nil {
			a.Error(err, http.StatusInternalServerError)
			return
		}

		var service2Logo string
		rowsService2, err := a.Area.Database.Query("SELECT logo FROM services WHERE id = $1", appletService2)
		if err != nil {
			a.Error(err, http.StatusInternalServerError)
			return
		}
		rowsService2.Next()
		if err := rowsService2.Scan(&service2Logo); err != nil {
			a.Error(err, http.StatusInternalServerError)
			return
		}

		applets = append(applets, map[string]any{
			"id":    appletId,
			"name":  appletName,
			"link":  appletLink,
			"users": appletUsers,
			"service": map[string]any{
				"name":        service1Name,
				"logo":        service1Logo,
				"logopartner": service2Logo,
				"color": map[string]any{
					"normal": service1ColorN,
					"hover":  service1ColorH,
				},
			},
		})

	}

	if err := rowsApplet.Err(); err != nil {
		a.Error(err, http.StatusInternalServerError)
		return
	}

	a.Reply(map[string]any{
		"res": applets,
	}, http.StatusOK)
}
