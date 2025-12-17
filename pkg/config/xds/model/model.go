package model

//run to generate new model from ./proto/*.proto
//go:generate protoc -I=.  --go_opt=Madapter.proto=./model --go_opt=Maddress.proto=./model --go_opt=Mbootstrap.proto=./model --go_opt=Mcluster.proto=./model --go_opt=Mextension.proto=./model --go_opt=Mfilter.proto=./model --go_opt=Mlistener.proto=./model --go_opt=Mroute.proto=./model --go_opt=Mhealth_check.proto=./model --go_out=../../ ./*.proto
