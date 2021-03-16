package main

func main() {
	initLog()
	initConfig()
	go connectToServer()
	startHub()
}
