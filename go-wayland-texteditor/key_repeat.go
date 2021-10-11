package main

import "time"
import "sync"

var RepeatedFunc KeyReloader
var RepeatedKey string
var RepeatedKeyNotUnicode uint32
var RepeatedKeyTime uint64
var RepeatedKeyTimeAbs uint64
var RepeatedKeyMutex sync.Mutex

type KeyReloader interface {
	KeyReload(key string, notUnicode, time uint32)
}

func KeyRepeatSubscribe(function KeyReloader, key string, notUnicode, t uint32) {
	RepeatedKeyMutex.Lock()
	RepeatedFunc = function
	RepeatedKey = key
	RepeatedKeyNotUnicode = notUnicode
	RepeatedKeyTime = uint64(t)
	RepeatedKeyTimeAbs = uint64(time.Now().UnixNano() / 1000000)
	RepeatedKeyMutex.Unlock()
}

func init() {
	go func() {
		KeyRepeat := time.NewTicker(75 * time.Millisecond)
		for {
			select {
			case tim := <-KeyRepeat.C:
				var t = uint64(tim.UnixNano() / 1000000)
				RepeatedKeyMutex.Lock()
				function := RepeatedFunc
				key := RepeatedKey
				notUnicode := RepeatedKeyNotUnicode
				minTime := RepeatedKeyTime
				absTime := RepeatedKeyTimeAbs
				RepeatedKeyMutex.Unlock()

				if function != nil || key != "" || notUnicode != 0 {

					if t-absTime > 300 {
						go function.KeyReload(key, notUnicode, uint32(t-minTime))
					}
				}
			}
		}
	}()
}
