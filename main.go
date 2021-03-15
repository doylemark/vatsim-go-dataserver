package main

func main() {
	initConfig()
	go connectToServer()
	startHub()
}
