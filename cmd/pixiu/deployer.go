package main

type Deployer interface {
	initialize() error

	start() error

	stop() error
}
