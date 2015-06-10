package NonchangGameServer

import (
	// "../src/NonchangGameServer/model"
	"../src/NonchangGameServer/api"
	"appengine/aetest"
	// "appengine/memcache"
	"testing"
)

func Test(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	t.Fatalf(api.AccountTest())

}
