package main

import "time"
import "sync"

var RepeatedFunc KeyReloader
var RepeatedKey string
var RepeatedKeyNotUnicode uint32
var RepeatedKeyTime uint32
var RepeatedKeyMutex sync.Mutex

type KeyReloader interface {
	KeyReload(key string, notUnicode, time uint32)
}

func KeyRepeatSubscribe(function KeyReloader, key string, notUnicode, time uint32) {
	RepeatedKeyMutex.Lock()
	RepeatedFunc = function
	RepeatedKey = key
	RepeatedKeyNotUnicode = notUnicode
	RepeatedKeyTime = time
	RepeatedKeyMutex.Unlock()
}

func init() {
	go func() {
		KeyRepeat := time.NewTicker(50 * time.Millisecond)
		for {
			select {
			case time := <-KeyRepeat.C:
				var t = uint32(time.UnixNano()/1000000)
				RepeatedKeyMutex.Lock()
				function := RepeatedFunc
				key := RepeatedKey
				notUnicode := RepeatedKeyNotUnicode
				minTime := RepeatedKeyTime
				RepeatedKeyMutex.Unlock()
				if (function != nil || key != "" || notUnicode != 0) {
					if (uint32(t) - uint32(minTime)) > 200 {
						go function.KeyReload(key, notUnicode, t)
					}
				}
			}
		}
	}()
}
