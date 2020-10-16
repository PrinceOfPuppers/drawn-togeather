package main



import (
	"fmt"
	"net"
	//"bytes"
	//"io"
	//"encoding/json"
	//"time"
)

// server event enum
const (
	inJSON = iota		// {inJSON, jsonSender{sessionNum, json}}
	sessionDisc = iota	// {sessionDisc, sessionNum}
)


type event struct{
	eType int
	data interface{} // ie session disconnect would have session number (int)
}
func newEvent(eventType int,data interface{}) *event{
	return &event{eventType, data}
}

type server struct{
	active *bool
	nextSessNum *int //counts up to give each session a unique number, used in listen loop

	li net.Listener
	eventChan chan *event // server is lone reader, sessions and such can write
	relayChan chan *jsonSender // server is lone writer, relay routine relays to all sessions
	sessions map[int]*session
}
func newServer(ip string, port string, eventChanSize int, relayChanSize int) *server{
	li, err := net.Listen("tcp",ip+":"+port)
	if err != nil{
		fmt.Println(err)
	}

	active := true
	nextSessNum := 1
	eventChan := make(chan *event,eventChanSize)
	relayChan := make(chan *jsonSender,relayChanSize)
	sessions := make(map[int]*session)

	serv := server{&active, &nextSessNum, li, eventChan, relayChan, sessions}

	go serv.relayLoop()
	go serv.listenLoop()
	return &serv
}
func (serv server) isActive() bool{
	return *serv.active
}
func (serv server) close(){
	for _, session := range serv.sessions {
		session.closeConnection()
	}
	*serv.active = false
}
// go routine for relaying messages in relayChan to all sessions except sender
func (serv server) relayLoop(){
	for serv.isActive(){
		data := <- serv.relayChan
		for _, session := range serv.sessions {
			//if sessionNum != data.sessionNum{
			session.push(data.json)
			//}

		}
	}
}
// go routine for adding new sessions to the server
func (serv server) listenLoop(){
	fmt.Println("listening...")
	for serv.isActive(){
		conn, err := serv.li.Accept()
		if err!=nil{
			fmt.Println(err)
		}

		s := newSession(*serv.nextSessNum, conn,&serv.eventChan, 3)
		serv.sessions[*serv.nextSessNum] = s
		
		*serv.nextSessNum++

		// add new session event here
	}

}

func (serv server) mainLoop(){

	for serv.isActive(){
		event := *<-serv.eventChan

		switch event.eType {
		case inJSON:
			data, ok := event.data.(*jsonSender)
			if !ok{
				fmt.Println("wrong data type in event")
			}
			serv.relayChan <- data
		
		case sessionDisc:
			data, ok := event.data.(int)
			if !ok{
				fmt.Println("wrong data type in event")
			}
			delete(serv.sessions,data)
			fmt.Printf("session number %v disconnected\n",data)
		}
	}
}
