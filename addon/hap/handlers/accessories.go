package handlers

import (
	"gateway/logger"
	"gateway/services/hap"
	"gateway/services/hap/controller"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type Accessories struct {
	http.Handler
	mutex                 *sync.Mutex
	accessoriesController *controller.AccessoriesController
}

func NewAccessories(mutex *sync.Mutex, accessoriesController *controller.AccessoriesController) *Accessories {
	return &Accessories{
		mutex:                 mutex,
		accessoriesController: accessoriesController,
	}
}

//handle request: GET /accessories
func (handler *Accessories) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	logger.Info.Printf("%v GET /accessories", request.RemoteAddr)
	responseWriter.Header().Set("Content-Type", hap.HTTPContentTypeHAPJson)
	handler.mutex.Lock()
	res, err := handler.accessoriesController.HandleGetAccessories(request.Body)
	handler.mutex.Unlock()
	if err != nil {
		logger.Error.Print(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
	} else {
		wr := hap.NewChunkedWriter(responseWriter, 2048)
		b, _ := ioutil.ReadAll(res)
		log.Println(string(b))
		_, err := wr.Write(b)
		if err != nil {
			log.Panic(err)
		}

	}
}
