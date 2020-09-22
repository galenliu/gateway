package handlers

import (
	"net/http"
)

type PairSetup struct {
	http.Handler
}

func NewPairSetup() *PairSetup {
	return &PairSetup{}

}

func (p *PairSetup) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
