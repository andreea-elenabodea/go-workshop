package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"workshop_demo/client"
	"workshop_demo/converter"
	"workshop_demo/server"

	"golang.org/x/exp/slog"
)

func main() {

	// http.HandleFunc("/quotas",
	// 	func(w http.ResponseWriter, r *http.Request) {
	// 		token := r.Header.Get("Authorization")

	// 		dbaasQuotas, err := client.DBaaSQuotas(token)
	// 		if err != nil {
	// 			w.WriteHeader(http.StatusInternalServerError)
	// 			slog.Error("Error while getting DBaaS quotas: ", err)
	// 			return
	// 		}convert_testrr := client.DNSQuotas(token)
	// 		if err != nil {
	// 			w.WriteHeader(http.StatusInternalServerError)
	// 			slog.Error("Error while getting DNS quotas: ", err)
	// 			return
	// 		}

	// 		serverResponse := convertModels(*dnsResponse, *dbaasQuotas)

	// 		jsonBody, err := json.Marshal(serverResponse)
	// 		if err != nil {
	// 			w.WriteHeader(http.StatusInternalServerError)
	// 			slog.Error("Error while marshaling response: ", err)
	// 			return
	// 		}

	// 		w.WriteHeader(http.StatusOK)
	// 		w.Write(jsonBody)
	// 	},
	// )

	serverImplementation := ServerImplementation{}
	handler := server.Handler(&serverImplementation)

	println("starting server")
	http.ListenAndServe(":8080", handler)
}

// func convertModels(
// 	dns client.DNSResponse,
// 	dbaas client.DBaaSResponse,
// ) model.ServerResponse {
// 	return model.ServerResponse{
// 		DNSResponse: model.Quota[model.DNSQuota]{
// 			Limit: dns.QuotaLimits,
// 			Usage: dns.QuotaUsage,
// 		},
// 		DBaaSResponse: model.Quota[model.DatabaseQuota]{
// 			Limit: dbaas.QuotaLimits,
// 			Usage: dbaas.QuotaUsage,
// 		},
// 	}
// }

// func convertModels(
//     dns model.Quota[model.DNSQuota],
//     dbaas model.Quota[model.DatabaseQuota],
// ) server.Quotas {
//     return server.Quotas{
//         DBaaS: &struct {
//             Limits *server.DBaaSQuota "json:\"Limits,omitempty\""
//             Usage  *server.DBaaSQuota "json:\"Usage,omitempty\""
//         }{
//             Limits: &server.DBaaSQuota{
//                 CPU:              stringToInt64Ptr(dbaas.Limit.Cpu),
//                 Memory:           stringToInt64Ptr(dbaas.Limit.Memory),
//                 MongoClusters:    stringToInt64Ptr(dbaas.Limit.CountMongoclustersDbaasIonosCom),
//                 PostgresClusters: stringToInt64Ptr(dbaas.Limit.CountPostgresclustersDbaasIonosCom),
//                 Storage:          stringToInt64Ptr(dbaas.Limit.Storage)},
//             Usage: &server.DBaaSQuota{
//                 CPU:              stringToInt64Ptr(dbaas.Usage.Cpu),
//                 Memory:           stringToInt64Ptr(dbaas.Usage.Memory),
//                 MongoClusters:    stringToInt64Ptr(dbaas.Usage.CountMongoclustersDbaasIonosCom),
//                 PostgresClusters: stringToInt64Ptr(dbaas.Usage.CountPostgresclustersDbaasIonosCom),
//                 Storage:          stringToInt64Ptr(dbaas.Usage.Storage)},
//         },
//         DNS: &struct {
//             Limits *server.DNSQuota "json:\"Limits,omitempty\""
//             Usage  *server.DNSQuota "json:\"Usage,omitempty\""
//         }{
//             Limits: &server.DNSQuota{
//                 Records:        intToInt64Ptr(dns.Limit.Records),
//                 Zones:          intToInt64Ptr(dns.Limit.Zones),
//                 SecondaryZones: intToInt64Ptr(dns.Limit.SecondaryZones),
//             },
//             Usage: &server.DNSQuota{
//                 Records:         intToInt64Ptr(dns.Usage.Records),
//                 Zones:           intToInt64Ptr(dns.Usage.Zones),
//                 SecondaryZones:  intToInt64Ptr(dns.Usage.SecondaryZones),
//             },
//         },
//     }
// }

func stringToInt64Ptr(s string) *int64 {
    integer, _ := strconv.ParseInt(s, 10, 64)
    return &integer
}

func intToInt64Ptr(i int) *int64 {
    integer := int64(i)
    return &integer
}

type ServerImplementation struct {

}

func( s *ServerImplementation) GetHealth(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func( s *ServerImplementation) GetQuotas(w http.ResponseWriter, r *http.Request){

	
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

	serverResponse := converter.ConvertModels(*dnsResponse, *dbaasQuotas)

	jsonBody, err := json.Marshal(serverResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Error while marshaling response: ", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBody)

}
