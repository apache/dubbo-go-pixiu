ADDRESS="localhost:8883"

API1=$(curl -s -X POST ${ADDRESS}"/api/v1/test-dubbo/UserService/com.dubbogo.pixiu.UserService?group=test&version=1.0.0&method=GetUserByName" -d '{"types":["string"],"values":"tc"}' --header 'Content-Type: application/json')
API2=$(curl -s -X POST ${ADDRESS}"/api/v1/test-dubbo/UserService/com.dubbogo.pixiu.UserService?group=test&version=1.0.0&method=UpdateUserByName" -d '{"types":["string","body"],"values":["tc",{"id":"0001","code":1,"name":"tc","age":15}]}' --header "Content-Type: application/json")
API3=$(curl -s -X POST ${ADDRESS}"/api/v1/test-dubbo/UserService/com.dubbogo.pixiu.UserService?group=test&version=1.0.0&method=GetUserByCode" -d '{"types":["int"],"values":1}' --header "Content-Type: application/json")

ARRAY_API=(${API1} ${API2} ${API3})

for element in ${ARRAY_API[@]}
do
echo ${element}
done