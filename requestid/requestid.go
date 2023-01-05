// Package requestid for requestid
package requestid

import (
	"github.com/v2pro/plz/gls"
	"sync"
	"time"
)

var (
	requestIDs = map[int64]interface{}{}
	rwm        sync.RWMutex
)

//Set 设置一个 RequestID
func Set(id interface{}) {
	goID := getGoID()
	rwm.Lock()
	defer rwm.Unlock()

	requestIDs[goID] = id
	go func() {
		time.Sleep(10 * time.Second)
		Delete()
	}()
}

//Get 返回设置的 RequestID
func Get() interface{} {
	goID := getGoID()
	rwm.RLock()
	defer rwm.RUnlock()

	return requestIDs[goID]
}

//Delete 删除设置的 RequestID
func Delete() {
	goID := getGoID()
	rwm.Lock()
	defer rwm.Unlock()

	delete(requestIDs, goID)
}

func getGoID() int64 {
	return int64(Goid())
}

//Goid func
func Goid() int64 {
	return gls.GoID()
}
