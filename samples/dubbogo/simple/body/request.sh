ADDRESS="localhost:8881"

API1=$(curl -s -X POST ${ADDRESS}"/api/v1/test-dubbo/user" -d '{"id":"0003","code":3,"name":"dubbogo","age":99}' --header 'Content-Type: application/json')
API2=$(curl -s -X PUT ${ADDRESS}"/api/v1/test-dubbo/user" -d '{"id":"0003","code":3,"name":"dubbogo","age":77}' --header "Content-Type: application/json")
API3=$(curl -s -X PUT ${ADDRESS}"/api/v1/test-dubbo/user2" -d '{"name":"dubbogo","user":{"id":"0003","code":3,"name":"dubbogo","age":88}}' --header "Content-Type: application/json")

ARRAY_API=(${API1} ${API2} ${API3})

for element in ${ARRAY_API[@]}
do
echo ${element}
#QUERY=$(curl -s -X GET ${ADDRESS}"/api/v1/test-dubbo/userByName?name=dubbogo")
#echo ${QUERY}
done