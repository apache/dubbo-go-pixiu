ADDRESS="localhost:8884"

API1=$(curl -s -X GET ${ADDRESS}"/api/v1/test-dubbo/userByName?name=tc")
API2=$(curl -s -X GET ${ADDRESS}"/api/v1/test-dubbo/userByCode?code=1")
API3=$(curl -s -X GET ${ADDRESS}"/api/v1/test-dubbo/userByNameAndAge?name=tc&age=99")

ARRAY_API=(${API1} ${API2} ${API3})

for element in ${ARRAY_API[@]}
do
echo ${element}
done
