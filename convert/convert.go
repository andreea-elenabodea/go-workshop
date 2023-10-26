package convert

import (
	"strconv"
	"workshop_demo/client"
	"workshop_demo/model"
)

func ConvertModels(
	dns client.DNSResponse,
	dbaas client.DBaaSResponse,
) model.ServerResponse {
	return model.ServerResponse{
		DNSResponse: model.Quota[model.DNSQuota]{
			Limit: dns.QuotaLimits,
			Usage: dns.QuotaUsage,
		},
		DBaaSResponse: model.Quota[model.DatabaseQuota]{
			Limit: dbaas.QuotaLimits,
			Usage: dbaas.QuotaUsage,
		},
	}
}

func stringToInt64Ptr(s string) *int64 {
	integer, _ := strconv.ParseInt(s, 10, 64)
	return &integer
}

func intToInt64Ptr(i int) *int64 {
	integer := int64(i)
	return &integer
}

// func ConvertModels(
// 	dns client.DNSResponse,
// 	dbaas client.DBaaSResponse,
// ) server.Quotas {
// 	return server.Quotas{
// 		DBaaS: &struct {
// 			Limits *server.DBaaSQuota "json:\"Limits,omitempty\""
// 			Usage  *server.DBaaSQuota "json:\"Usage,omitempty\""
// 		}{
// 			Limits: &server.DBaaSQuota{
// 				CPU:              stringToInt64Ptr(dbaas.QuotaLimits.Cpu),
// 				Memory:           stringToInt64Ptr(dbaas.QuotaLimits.Memory),
// 				MongoClusters:    stringToInt64Ptr(dbaas.QuotaLimits.CountMongoclustersDbaasIonosCom),
// 				PostgresClusters: stringToInt64Ptr(dbaas.QuotaLimits.CountPostgresclustersDbaasIonosCom),
// 				Storage:          stringToInt64Ptr(dbaas.QuotaLimits.Storage)},
// 			Usage: &server.DBaaSQuota{
// 				CPU:              stringToInt64Ptr(dbaas.QuotaUsage.Cpu),
// 				Memory:           stringToInt64Ptr(dbaas.QuotaUsage.Memory),
// 				MongoClusters:    stringToInt64Ptr(dbaas.QuotaUsage.CountMongoclustersDbaasIonosCom),
// 				PostgresClusters: stringToInt64Ptr(dbaas.QuotaUsage.CountPostgresclustersDbaasIonosCom),
// 				Storage:          stringToInt64Ptr(dbaas.QuotaUsage.Storage)},
// 		},
// 		DNS: &struct {
// 			Limits *server.DNSQuota "json:\"Limits,omitempty\""
// 			Usage  *server.DNSQuota "json:\"Usage,omitempty\""
// 		}{
// 			Limits: &server.DNSQuota{
// 				Records:        intToInt64Ptr(dns.QuotaLimits.Records),
// 				Zones:          intToInt64Ptr(dns.QuotaLimits.Zones),
// 				SecondaryZones: intToInt64Ptr(dns.QuotaLimits.SecondaryZones),
// 			},
// 			Usage: &server.DNSQuota{
// 				Records:        intToInt64Ptr(dns.QuotaUsage.Records),
// 				Zones:          intToInt64Ptr(dns.QuotaUsage.Zones),
// 				SecondaryZones: intToInt64Ptr(dns.QuotaUsage.SecondaryZones),
// 			},
// 		},
// 	}
// }
