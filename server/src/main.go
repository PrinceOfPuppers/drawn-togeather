package main

import ()


func main() {
	serv := newServer("127.0.0.1","1234",3,3)
	serv.mainLoop()
}

