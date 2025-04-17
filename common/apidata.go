package common

import (
	"encoding/json"
	"log"

	"sync"
)

const LCDRowLength = 16

var MutexRW sync.RWMutex

type Apidata struct {
	Phenomena string `json:"phenomena"`
	Main      string `json:"main"`
	Timestamp int64  `json:"timestamp"`
}

func (a *Apidata) GetJson() []byte {
	MutexRW.RLock()
	defer MutexRW.RUnlock()

	jsonEncoded, err := json.Marshal(a)

	if err != nil {
		log.Printf("JSON encode error! %s\n", err.Error())
		return nil
	}

	return jsonEncoded
}

func (a *Apidata) SetData(adata Apidata) {
	MutexRW.Lock()
	defer MutexRW.Unlock()

	*a = adata
}
