#
#  Licensed to the Apache Software Foundation (ASF) under one or more
#  contributor license agreements.  See the NOTICE file distributed with
#  this work for additional information regarding copyright ownership.
#  The ASF licenses this file to You under the Apache License, Version 2.0
#  (the "License"); you may not use this file except in compliance with
#  the License.  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.

zkJarName="zookeeper-3.4.9-fatjar.jar"
remoteJarUrl="https://github.com/dubbogo/resources/raw/master/zookeeper-4unitest/contrib/fatjar/${zkJarName}"
zkJarPath="pkg/registry/zookeeper-4unittest/contrib/fatjar"
zkJar="${zkJarPath}/${zkJarName}"

if [ ! -f "${zkJar}" ]; then
    mkdir -p ${zkJarPath}
    wget -P "${zkJarPath}" ${remoteJarUrl}
fi

# download envoy
envoyPath="out/linux_amd64/"
sudo apt update
sudo apt install -y apt-transport-https gnupg2 curl lsb-release
curl -sL 'https://deb.dl.getenvoy.io/public/gpg.8115BA8E629CC074.key' | sudo gpg --dearmor --yes -o /usr/share/keyrings/getenvoy-keyring.gpg
echo a077cb587a1b622e03aa4bf2f3689de14658a9497a9af2c427bba5f4cc3c4723 /usr/share/keyrings/getenvoy-keyring.gpg | sha256sum --check
echo "deb [arch=amd64 signed-by=/usr/share/keyrings/getenvoy-keyring.gpg] https://deb.dl.getenvoy.io/public/deb/ubuntu $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/getenvoy.list
sudo apt update
sudo apt install -y getenvoy-envoy
mkdir -p ${envoyPath}
cp /usr/bin/envoy ${envoyPath}/envoy