#include <QtCore/QCoreApplication>
#include "socket.h"
#include <iostream>


#include <QJsonObject>
#include <QtDebug>

#include <unistd.h>
#include <thread>
#include <chrono>
using namespace std;


int main(){
    Socket sock;
    sock.init();


    const char ip[] = "127.0.0.1";

    sock.connect(ip, 1234);
    QJsonObject json;
    json.insert("test",1.0);
    json.insert("調子はどう",4);


    sock.push(&json);

    QJsonObject response = *sock.pop();

    qDebug() << response;
}