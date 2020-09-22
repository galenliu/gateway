package handlers

import (
	"gateway/services/hap"
	"gateway/services/hap/controller"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type Characteristics struct {
	http.Handler
	mutex      *sync.Mutex
	context    hap.Context
	controller *controller.CharacteristicController
}

func NewCharacteristics(mutex *sync.Mutex, cc *controller.CharacteristicController, context hap.Context) *Characteristics {
	handler := Characteristics{
		mutex:      mutex,
		controller: cc,
		context:    context,
	}
	return &handler
}

func (handler *Characteristics) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	var res io.Reader
	var err error

	handler.mutex.Lock()
	switch request.Method {
	case hap.MethodGet:
		request.ParseForm()
		log.Println("%v GET /characteristics %v", request.RemoteAddr, request.Form)
		session := handler.context.GetSessionForRequest(request)
		conn := session.Connection()
		res, err = handler.controller.HandleGetCharacteristics(request.Form, conn)
	case hap.MethodPut:
		log.Println("%v PUT /characteristics %v", request.RemoteAddr, request.Form)
		session := handler.context.GetSessionForRequest(request)
		conn := session.Connection()
		err = handler.controller.HandleUpdateCharacteristics(request.Body, conn)
	default:
		log.Println("Cannot handle Http Method", request.Method)
	}
	handler.mutex.Unlock()

	if err != nil {
		log.Panic(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
	} else {
		if res != nil {
			responseWriter.Header().Set("Content-Type", hap.HTTPContentTypeHAPJson)
			wr := hap.NewChunkedWriter(responseWriter, 2048)
			b, _ := ioutil.ReadAll(res)
			wr.Write(b)
		} else {
		}
		responseWriter.WriteHeader(http.StatusNoContent)
	}
}
