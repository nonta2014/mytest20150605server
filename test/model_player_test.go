package NonchangGameServer

import (
	"../src/NonchangGameServer/model"
	// "../src/NonchangGameServer/api"
	"appengine/aetest"
	// "appengine/memcache"
	"testing"
)

func TestPlayerModel(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	// it := &memcache.Item{
	// 	Key:   "some-key",
	// 	Value: []byte("some-value"),
	// }
	// err = memcache.Set(c, it)
	// if err != nil {
	// 	t.Fatalf("Set err: %v", err)
	// }
	// it, err = memcache.Get(c, "some-key")
	// if err != nil {
	// 	t.Fatalf("Get err: %v; want no error", err)
	// }
	// if g, w := string(it.Value), "some-value"; g != w {
	// 	t.Errorf("retrieved Item.Value = %q, want %q", g, w)
	// }

	// if model.GetTestText() != "teststring" {
	// 	// if "GetTestText()"!="teststring" {
	// 	t.Fatalf("GetTestTextの値が「teststring」ではありませんでした。")
	// }

	player := &model.Player{Name: "テストプレイヤーレコード", Level: 1}

	if player.Name != "テストプレイヤーレコード" {
		// if "GetTestText()"!="teststring" {
		t.Fatalf(`player.name値が"テストプレイヤーレコード"" ではありませんでした。`)
	}

}
