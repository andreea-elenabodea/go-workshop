package convert_test

import (
	"testing"
	"workshop_demo/client"
	"workshop_demo/convert"
	"workshop_demo/model"

	"github.com/sagikazarmark/slog-shim"
	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	token := "Bearer eyJ0eXAiOiJKV1QiLCJraWQiOiJjMzJmYjk5NS1hYWE3LTQ3MmUtYjQxNi0yOGM5MGNmNjA3MzIiLCJhbGciOiJSUzI1NiJ9.eyJpc3MiOiJpb25vc2Nsb3VkIiwiaWF0IjoxNjk4MjM3MzI4LCJjbGllbnQiOiJVU0VSIiwiaWRlbnRpdHkiOnsiY29udHJhY3ROdW1iZXIiOjMxOTUxMDU4LCJpc1BhcmVudCI6ZmFsc2UsInByaXZpbGVnZXMiOlsiTUFOQUdFX0RCQUFTIiwiQUNDRVNTX0FORF9NQU5BR0VfRE5TIl0sInV1aWQiOiI1MzgyZTlkYS00YzIwLTRmOTAtYjZjYy0zMjRkZjVkNDQzZGQiLCJyZXNlbGxlcklkIjoxLCJyZWdEb21haW4iOiJpb25vcy5kZSIsInJvbGUiOiJ1c2VyIn0sImV4cCI6MTcyOTc5NDkyOH0.2Ph45TTBBFbkWcXplvvuv95luug6IR2EIZqC68mBOYHUcq0lXHFM-e2OtwNy-33UkJgfV02XRFtVl5sO9mPCjeK5a5GVrsvrSJ_uVWVIGHNp1tfJDBXOD_JhUnJTcIGVF5wsjaH8Gh1vizmKtAZFvNL1prW-zvo2et2Q1lr1XK0gAwqwq7_PPIu3SOp8UfeE19_eqKpF_9IQNNBfNYeZc99CUXYK42AON9pApyOdqeL3nf7wiU1dxM-7aI6meGv6cvokuAW__8k7FJELgX_QFwWYJr1Sn9Gdnb5DYAKHMcvzL3w415bx7zh9soNwCHzSop0rpjAlrFYlBv0Rz8dtuA"

	dbaasQuotas, err := client.DBaaSQuotas(token)
	if err != nil {
		slog.Error("Error while getting DBaaS quotas: ", err)
		return
	}

	dnsResponse, err := client.DNSQuotas(token)
	if err != nil {
		slog.Error("Error while getting DNS quotas: ", err)
		return
	}
	actual := model.ServerResponse{
		DNSResponse: model.Quota[model.DNSQuota]{
			Limit: model.DNSQuota{
				Records:        100000,
				SecondaryZones: 100000,
				Zones:          50000,
			},
			Usage: model.DNSQuota{
				Records:        0,
				SecondaryZones: 0,
				Zones:          0,
			},
		},
		DBaaSResponse: model.Quota[model.DatabaseQuota]{
			Limit: model.DatabaseQuota{
				CountMongoclustersDbaasIonosCom:    "5",
				CountPostgresclustersDbaasIonosCom: "10",
				Cpu:                                "16",
				Memory:                             "32768",
				Storage:                            "1536000",
			},
			Usage: model.DatabaseQuota{
				CountMongoclustersDbaasIonosCom:    "",
				CountPostgresclustersDbaasIonosCom: "1",
				Cpu:                                "1",
				Memory:                             "2048",
				Storage:                            "2048",
			},
		},
	}
	assert.Equal(t, actual, convert.ConvertModels(*dnsResponse, *dbaasQuotas))
}
