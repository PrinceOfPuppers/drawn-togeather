# pragma once

#include <QObject>
#include <QTcpSocket>
#include <QByteArray>
#include <QJsonObject>
#include <QJsonDocument>
#include <iostream>
#include <queue>

#include <thread>
#include <mutex>
#include <condition_variable>

using namespace std;

template <class T>
class BlockingQueue
{

    public:

        void push(T datum){
            unique_lock<mutex> lock(mut);

            this->data.push(datum);

            this->cond.notify_one();
            lock.unlock();
        }
        T pop(){
            std::unique_lock<std::mutex> lock(mut);
            this->cond.wait( lock, [&](){return !this->data.empty();} );
            T datum = this->data.front();
            this->data.pop();
            lock.unlock();
            return datum;
        }

        bool empty(){return this->data.empty();}

    private:
        queue<T> data;
        mutex mut;
        condition_variable cond;
};

//enum class sock_event_type{
//    in_json, 
//    serv_close
//};
//
//template <typename T>
//struct sock_event{
//    sock_event_type e_type;
//    T data;
//
//};


class Socket
{
    public: 
        void init();
        void connect(const char[], int);
        void push(QJsonObject*);
        QJsonObject* pop();
        bool empty(); // indicates if recieving is empty

        bool is_active();

        thread sending_thread;
        thread receiving_thread;
    private:
        QTcpSocket conn;


        BlockingQueue<QJsonObject*> sending;
        BlockingQueue<QJsonObject*> receiving;

        void sending_loop();
        void receiving_loop();

        void send_json(QJsonObject *json_obj);
        QJsonObject* recieve_json();

};