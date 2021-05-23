#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#
ADDRESS="localhost:8883"

API1=$(curl -s -X POST ${ADDRESS}"/api/v1/test-dubbo/UserService/com.dubbogo.pixiu.UserService?group=test&version=1.0.0&method=GetUserByName" -d '{"types":["string"],"values":"tc"}' --header 'Content-Type: application/json')
API2=$(curl -s -X POST ${ADDRESS}"/api/v1/test-dubbo/UserService/com.dubbogo.pixiu.UserService?group=test&version=1.0.0&method=UpdateUserByName" -d '{"types":["string","body"],"values":["tc",{"id":"0001","code":1,"name":"tc","age":15}]}' --header "Content-Type: application/json")
API3=$(curl -s -X POST ${ADDRESS}"/api/v1/test-dubbo/UserService/com.dubbogo.pixiu.UserService?group=test&version=1.0.0&method=GetUserByCode" -d '{"types":["int"],"values":1}' --header "Content-Type: application/json")

ARRAY_API=(${API1} ${API2} ${API3})

for element in ${ARRAY_API[@]}
do
echo ${element}
done