#include "socket.h"
#include <QFile>

void Socket::connect(const char ip[], int port){

    cout <<"connecting to: "<< ip <<":"<<1234<<"\n";
    this->conn.connectToHost(ip,port);

    int waittime = 30000;

    if(conn.waitForConnected(waittime)){

        cout <<"connected\n\n";
    }
    else{
        cout <<"not connected\n";
    }


    this->sending_thread = thread(&Socket::sending_loop, this);
    this->receiving_thread = thread(&Socket::receiving_loop, this);



}
bool Socket::empty(){
    return this->receiving.empty();
}
void Socket::push(QJsonObject* data){
    this->sending.push(data);
}
QJsonObject* Socket::pop(){
    return this->receiving.pop();
}
bool Socket::is_active(){
    return (this->conn.state() != 0);
}



// TODO check both of these for memeory leaks
void Socket::send_json(QJsonObject *json){

    QByteArray serialized = QJsonDocument(*json).toJson();
    conn.write(serialized);

    conn.waitForBytesWritten(); 
}

QJsonObject* Socket::recieve_json(){
    while(true){

        if (conn.waitForReadyRead(-1)){

            QByteArray data = conn.readAll();
            QJsonDocument response = QJsonDocument::fromJson(data);

            QJsonObject *obj = new QJsonObject(response.object());

            return obj; 
        }
    }
}

// sends and deallocates messages in the sending queue
void Socket::sending_loop(){
    while (this->is_active()){
        QJsonObject *json = this->sending.pop();
        this->send_json(json);
        delete json; // message is sent, deallocating queue item
    }
}
void Socket::receiving_loop(){
    while (this->is_active()){
        QJsonObject* a = this->recieve_json();

        this->receiving.push(a);
    }
}