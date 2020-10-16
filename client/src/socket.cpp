#include "socket.h"
#include <QFile>

void Socket::init(){
    cout << "trest"<<"\n";
}

void Socket::connect(const char ip[], int port){


    cout <<"connecting to: "<< ip <<":"<<1234<<"\n";
    this->conn.connectToHost(ip,port);

    int waittime = 30000;

    if(conn.waitForConnected(waittime)){

        cout <<"connected"<<"\n";
    }
    else{
        std::cout <<"not connected"<<"\n";
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
void Socket::send_json(QJsonObject *json_obj){

    QByteArray serialized = QJsonDocument(*json_obj).toJson();
    qDebug() << serialized;
    conn.write(serialized);

    conn.waitForBytesWritten(); 
}

QJsonObject* Socket::recieve_json(){
    while(true){

        conn.waitForReadyRead(); 

        if (conn.bytesAvailable() > 0){
            QByteArray data = conn.readAll();
            QJsonDocument response = QJsonDocument::fromJson(data);

            qDebug() << response;

            QJsonObject *obj = new QJsonObject(response.object());
            return obj; 
        }
    }
}

void Socket::sending_loop(){
    while (this->is_active()){
        this->send_json(this->sending.pop());
    }
}
void Socket::receiving_loop(){
    while (this->is_active()){
        this->receiving.push(this->recieve_json());
    }
}