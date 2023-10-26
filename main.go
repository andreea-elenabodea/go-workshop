package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"workshop_demo/client"
	"workshop_demo/convert"
	"workshop_demo/server"

	"golang.org/x/exp/slog"
)

func main() {
	serverImplementation := ServerImplementation{}
	handler := server.Handler(&serverImplementation)

	println("starting server")
	http.ListenAndServe(":8080", handler)
}

type ServerImplementation struct {
}

func (s *ServerImplementation) GetHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (s *ServerImplementation) GetQuotas(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")

	dbaasQuotas, err := client.DBaaSQuotas(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Error while getting DBaaS quotas: ", err)
		return
	}

	dnsResponse, err := client.DNSQuotas(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Error while getting DNS quotas: ", err)
		return
	}

<<<<<<< HEAD
	serverResponse := convert.ConvertModels(*dnsResponse, *dbaasQuotas)

	jsonBody, err := json.Marshal(serverResponse)

=======
	serverResponse := converter.ConvertModels(*dnsResponse, *dbaasQuotas)

	jsonBody, err := json.Marshal(serverResponse)
>>>>>>> 6e102e20353e89c2decc2ca8a3db0f67b12005d2
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Error while marshaling response: ", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBody)

}
