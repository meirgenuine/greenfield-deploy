package middleware

import (
	"fmt"
	"log"
	"net/http"
)

func HandlePanics(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			r := recover()
			if r != nil {
				log.Println("recovered", r)
				http.Error(w, fmt.Sprintf(`{"msg": "%v"}`, r), http.StatusInternalServerError)
			}
		}()
		f(w, r)
	}
}
