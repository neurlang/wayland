package wl

import "testing"

func TestUserDataIsPerConnection(t *testing.T) {
	a, b := &Context{objects: map[ProxyId]Proxy{}}, &Context{objects: map[ProxyId]Proxy{}}
	ra, rb := &Registry{}, &Registry{}
	a.RegisterMapped(ra, 2)
	b.RegisterMapped(rb, 2)
	SetUserData(ra, &a)
	SetUserData(rb, &b)
	defer DeleteUserData(ra)
	defer DeleteUserData(rb)
	if got, _ := GetUserData[Context](ra); got != a {
		t.Errorf("registry 2 on the first connection got the second connection's context")
	}
}
