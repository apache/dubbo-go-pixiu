ADDRESS="localhost:8882"

API1=$(curl -s -X GET ${ADDRESS}"/api/v1/test-dubbo/user/tc?age=99")
API2=$(curl -s -X PUT ${ADDRESS}"/api/v1/test-dubbo/user/tc" -d '{"id":"0001","code":1,"name":"tc","age":66}' --header "Content-Type: application/json")
API3=$(curl -s -X PUT ${ADDRESS}"/api/v1/test-dubbo/user?name=tc" -d '{"id":"0001","code":1,"name":"tc","age":55}' --header "Content-Type: application/json")

ARRAY_API=(${API1} ${API2} ${API3})

for element in ${ARRAY_API[@]}
do
echo ${element}
done