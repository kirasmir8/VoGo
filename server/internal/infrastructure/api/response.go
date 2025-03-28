package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func StatusMessageResponse(w http.ResponseWriter, status int, msg interface{}) {
	w.WriteHeader(status)
	if msg != nil {
		err := json.NewEncoder(w).Encode(msg)
		if err != nil {
			fmt.Println("Ошибка отправки данных: ", err)
		}
	}
}
