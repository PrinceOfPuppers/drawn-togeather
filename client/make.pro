CONFIG += debug


INCLUDEPATH += include

OBJECTS_DIR = tmp
MOC_DIR = tmp

SOURCES += src/*.cpp
HEADERS += include/*.h

QT += core
QT += network

TARGET = build/client