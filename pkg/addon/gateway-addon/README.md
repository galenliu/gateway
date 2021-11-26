```go
package main

import (
	addon "github.com/galenliu/gateway-addon"
	"log"
)

type ExampleAdapter struct {
	*addon.Adapter
}

// StartPairing 客户端配对时，回调的方法
func (a *ExampleAdapter) StartPairing(timeout float64) {
	log.Printf("adapter:(%s)- StartPairing() not implemented", a.GetID())
}

// CancelPairing 客户端取消配对时回调的方法
func (a *ExampleAdapter) CancelPairing() {

	log.Printf("adapter:(%s)- CancelPairing() not implemented", a.GetID())

}

type ExampleDevice struct {
	*addon.Device
}

type ExampleProperty struct {
	*addon.Property
}
```