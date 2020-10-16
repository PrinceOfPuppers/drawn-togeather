#include <QtCore/QCoreApplication>
#include "socket.h"
#include <iostream>


#include <QJsonObject>
#include <QtDebug>

#include <unistd.h>
#include <thread>
#include <chrono>
using namespace std;

int readInput(Socket *s){
    while (true){
        QJsonObject response = *s->pop();
        
        cout << "> " <<response["message"].toString().toUtf8().constData();
    }

}


int main(){
    Socket sock;


    const char ip[] = "127.0.0.1";

    sock.connect(ip, 1234);


    thread t = thread(readInput, &sock);
    while (true){
        char* message = new(char);
        cin >> message;

        QJsonObject* reply = new(QJsonObject);
        //cout<<"testing...\n";
        reply->insert("message",message);
        //cout<<"...123\n";
        sock.push(reply);
    }
}