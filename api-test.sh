#curl --location https://api.ionos.com/databases/quota \
#  -H "Authorization: Bearer $IONOS_TOKEN" | jq
#
#curl --location https://dns.de-fra.ionos.com/quota \
#  -H "Authorization: Bearer $IONOS_TOKEN" | jq


<<<<<<< HEAD
curl -v localhost:8080/quotas -H "Authorization: Bearer $IONOS_TOKEN" | jq

curl -v localhost:8080/health -H "Authorization: Bearer $IONOS_TOKEN" | jq
=======
curl -v localhost:8080/quotas -H "Authorization: Bearer $GO_TOKEN" | jq
>>>>>>> 6e102e20353e89c2decc2ca8a3db0f67b12005d2
