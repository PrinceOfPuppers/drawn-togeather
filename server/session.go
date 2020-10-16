package main

import (
	"fmt"
	"net"
	//"bytes"
	//"io"
	"encoding/json"
	//"time"
)


type jsonSender struct{
	// a pair for holding json and who sent it

	sessionNum int
	json *map[string]interface{}
}
func newJSONSender(sessionNum int, json *map[string]interface{}) *jsonSender{
	return &jsonSender{sessionNum, json}
}



type session struct {
	active *bool
	sessionNum int
	conn net.Conn
	
	encoder *json.Encoder
	decoder *json.Decoder
	
	sendChan chan *map[string]interface{} // what is pushed to for sending messages
	servEventChan *chan *event // pointer to server events, what inbound messages are put into
} 
////////////
// public //
////////////
func newSession(sessionNum int, conn net.Conn, servEventChan *chan *event, chanSize int) *session{

    sendChan := make(chan *map[string]interface{},chanSize)

	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)
	active := true
    sess := session{&active, sessionNum, conn,encoder,decoder,sendChan,servEventChan}
    go sess.sendLoop()
    go sess.recieveLoop()


    return &sess
}
func (sess session) closeConnection(){
	if sess.isActive(){ // ensures conn is not already closed
		*sess.servEventChan <- newEvent(sessionDisc,sess.sessionNum)
		*sess.active = false
        sess.conn.Close()
    }
}
func (sess session) push(data *map[string]interface{}){
	if sess.isActive(){
		select {
		case sess.sendChan <- data:
		default:
			fmt.Println("session output channel full")
		}
	}else{
		fmt.Println("session inactive")
	}

}
func (sess session) isActive() bool {
    return *sess.active
}


////////////
// private//
////////////
func (sess session) recieveJSON() *map[string]interface{}{
	data := new(map[string]interface{})
	err := sess.decoder.Decode(data)

    if err != nil {
		sess.closeConnection()
		return nil
    }

    return data
}
func (sess session) sendJSON(data *map[string]interface{}){
	err := sess.encoder.Encode(data)
    
    if err!= nil{
        fmt.Println(err)
        sess.closeConnection()
    }

}
// routine for adding to recieved messages to servChanI
func (sess session) recieveLoop(){
    for sess.isActive() {
		jsonPntr := sess.recieveJSON()
		if jsonPntr != nil{
			*sess.servEventChan <- newEvent(inJSON,newJSONSender(sess.sessionNum, jsonPntr))
		}
    }
}
// routine for sending messages in servChanO
func (sess session) sendLoop(){
    for sess.isActive() {
        jsonPntr := <- sess.sendChan
        sess.sendJSON(jsonPntr)
    }
}