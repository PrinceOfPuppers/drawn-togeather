package main

import (
	"fmt"
	"bufio"
	"os"
)


func (serv server) eventSwitch(event event){
	switch event.eType {
	case inJSON:

		data, ok := event.data.(*jsonSender)
		if !ok{
			fmt.Println("wrong data type in event")
		}
		serv.jsonSwitch(data)
	
	case sessionDisc:
		data, ok := event.data.(int)
		if !ok{
			fmt.Println("wrong data type in event")
		}
		delete(serv.sessions,data)
		fmt.Printf("session number %v disconnected\n",data)
	}
}

func (serv server) jsonSwitch(json *jsonSender){
	fmt.Printf("> %v\n",(*json.json)["message"])
	//serv.relayChan <- json
}


func (serv server) mainLoop(){
	for serv.isActive(){
		event := *<-serv.eventChan
		serv.eventSwitch(event)

	}
}


func main() {
	serv := newServer("127.0.0.1","1234",3,3)
	go serv.mainLoop()

	reader := bufio.NewReader(os.Stdin)

	for{
		text, _ := reader.ReadString('\n')
		for _, session := range serv.sessions {
			servMsg := map[string]interface{}{"message": text}
			session.push(&servMsg)
		}
	}




}

