package wl

import "sync"

var userDataMap = new(sync.Map)

type userDataKey struct {
	ctx *Context
	id  ProxyId
}

func keyOf(p Proxy) userDataKey { return userDataKey{p.Context(), p.Id()} }

func SetUserData[T any](key Proxy, value *T) {
	//println("set", key.Id())
	userDataMap.Store(keyOf(key), *value)

}
func GetUserData[T any](key Proxy) (found *T, exists bool) {
	eface, is := userDataMap.Load(keyOf(key))
	if !is {
		return
	}
	found, exists = eface.(*T)
	return
}
func DeleteUserData(key Proxy) {
	//println("clr", key.Id())
	userDataMap.Delete(keyOf(key))
}
