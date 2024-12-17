package applet

import (
	"area-backend/area"
	"net/http"
	"fmt"
)

func GetApplets(a area.AreaRequest) {
	fmt.Println("get applets")
	rowsApplet, err := a.Area.Database.Query("SELECT id, name, link, users, service1, service2 FROM applets")
	if err != nil {
		a.Error(err, http.StatusInternalServerError)
		return
	}
	defer rowsApplet.Close()

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
		defer rowsService1.Close()
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
		defer rowsService2.Close()
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
