package logger

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/wazofski/store/utils"
)

type Logger interface {
	Printf(string, ...interface{})
	Fatalln(error)
}

type _Logger struct {
	Module string
}

func New(module string) Logger {
	return &_Logger{
		Module: module,
	}
}

type _Msg struct {
	Module    string
	Message   string
	Timestamp string
}

func jsonify(module, msg string) string {
	_msg := _Msg{
		Module:    module,
		Message:   msg,
		Timestamp: utils.Timestamp(),
	}

	data, _ := json.MarshalIndent(_msg, "", " ")

	return string(data)
}

func (l *_Logger) Printf(msg string, params ...interface{}) {
	fmt.Println(jsonify(l.Module, fmt.Sprintf(msg, params...)))
}

func (l *_Logger) Fatalln(msg error) {
	log.Panicf(jsonify(l.Module, msg.Error()))
}
