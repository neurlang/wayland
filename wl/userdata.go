package wl

import "sync"
import "weak"

var userDataMap = new(sync.Map)

type userDataKey1 weak.Pointer[Context]
type userDataKey2 ProxyId

func keyOf1(p Proxy) userDataKey1 { return userDataKey1(weak.Make(p.Context())) }
func keyOf2(p Proxy) userDataKey2 { return userDataKey2(p.Id()) }

func SetUserData[T any](key Proxy, value *T) {
	//println("set", key.Id())
	eface, is := userDataMap.Load(keyOf1(key))
	if is {
		eface.(*sync.Map).Store(keyOf2(key), *value)
	} else {
		var n = new(sync.Map)
		n.Store(keyOf2(key), *value)
		userDataMap.Store(keyOf1(key), n)

	}
}
func GetUserData[T any](key Proxy) (found *T, exists bool) {
	eface, is := userDataMap.Load(keyOf1(key))
	if !is {
		return
	}
	eface, is = eface.(*sync.Map).Load(keyOf2(key))
	if !is {
		return
	}

	found, exists = eface.(*T)
	return
}
func DeleteUserData(key Proxy) {
	//println("clr", key.Id())
	eface, is := userDataMap.Load(keyOf1(key))
	if !is {
		return
	}

	eface.(*sync.Map).Delete(keyOf2(key))
}

func CleanContextData(ctx *Context) {
	userDataMap.Delete(userDataKey1(weak.Make(ctx)))
}
