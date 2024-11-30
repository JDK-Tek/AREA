package area

import (
	"database/sql"
	"net/http"
	"encoding/json"
	"fmt"
)

type Area struct {
	Database *sql.DB
	Key string
}

type AreaRequest struct {
	Area *Area
	Writter http.ResponseWriter
	Request *http.Request
}

func (it *AreaRequest) Error(err error, code int) {
	it.ErrorStr(err.Error(), code)
}

func (it *AreaRequest) ErrorStr(str string, code int) {
	var data []byte
	var err error

	data, err = json.Marshal(map[string]string{
		"error": str,
	})
	if err != nil {
		http.Error(it.Writter, "marshal failed", code)
		return
	}
	http.Error(it.Writter, string(data), code)
}

func (it *AreaRequest) Reply(object any, code int) {
	var data []byte
	var err error

	it.Writter.WriteHeader(http.StatusOK)
	data, err = json.Marshal(object)
	if err != nil {
		it.Error(err, http.StatusInternalServerError)
		return
	}
    fmt.Fprintln(it.Writter, string(data))
}
