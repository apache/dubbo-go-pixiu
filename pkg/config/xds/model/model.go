// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package model

//run to generate new model from ./proto/*.proto
//go:generate protoc -I=.  --go_opt=Madapter.proto=./model --go_opt=Maddress.proto=./model --go_opt=Mbootstrap.proto=./model --go_opt=Mcluster.proto=./model --go_opt=Mextension.proto=./model --go_opt=Mfilter.proto=./model --go_opt=Mlistener.proto=./model --go_opt=Mroute.proto=./model --go_opt=Mhealth_check.proto=./model --go_out=../../ ./*.proto
