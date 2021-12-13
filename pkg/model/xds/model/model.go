package model

//run to generate new model from ./proto/*.proto
//go:generate protoc -I=./proto -I=PROTOC_LIB --go_out=./ ./proto/*.proto
