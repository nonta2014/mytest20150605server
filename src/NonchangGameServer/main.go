package NonchangGameServer

import (
	"appengine"
	"fmt"
	"net/http"

	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
	// "./api" //goapp testでコンパイルテストする時用（メモ：もう使えない。。）
	"NonchangGameServer/api"
	"NonchangGameServer/api/game"

	// "appengine/memcache"
	// "math/rand"
	// "strconv"
	// "time"

	// "encoding/json"
)

func init() {

	//=========================================
	// テスト応答

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "わほーいせかい")
		c := appengine.NewContext(r)
		c.Infof("test log !! : ", "testtest")
	})

	//=========================================
	// memcacheテスト（書きこみExpirationのみ8000錠で実験完了）

	// http.HandleFunc("/memcacheTest01", func(w http.ResponseWriter, r *http.Request) {
	// 	c := appengine.NewContext(r)
	// 	token := "Oh, give me a home" + strconv.FormatUint(uint64(rand.Int()), 10)
	// 	item := &memcache.Item{
	// 		Key:   "testtesttest-prefix-lyric",
	// 		Value: []byte(token), //メモ:memcacheのvalueは、GAE/gではbyte arrayであることが要求される。
	// 		//Time: 30,
	// 		Expiration: time.Duration(5) * time.Second,
	// 		//Expiration: time.Now(),//time.After(time.Second*10),
	// 	}
	// 	if err := memcache.Add(c, item); err == memcache.ErrNotStored {
	// 		fmt.Fprint(w, "item with key "+item.Key+" already exists")
	// 	} else if err != nil {
	// 		fmt.Fprint(w, "error adding item: "+err.Error())
	// 	} else {
	// 		fmt.Fprint(w, `{"status":"succeed","data","`+token+`"}`)
	// 	}

	// 	//以下は取得方法など

	// 	// if item, err := memcache.Get(c, "lyric"); err == memcache.ErrCacheMiss {
	// 	// 	// c.Infof("item not in the cache")
	// 	// 	fmt.Fprint(w, `{"status":"fail","message":"item not in the cache"}`)
	// 	// } else if err != nil {
	// 	// 	c.Errorf("error getting item: %v", err)
	// 	// } else {
	// 	// 	c.Infof("the lyric is %q", item.Value)
	// 	// }
	// })

	//=========================================
	//cloud endpoint

	//初期テスト
	// if _, err := api.RegisterClickLogService(); err != nil {
	// 	panic(err.Error())
	// }

	if _, err := api.RegisterAccountService(); err != nil {
		panic(err.Error())
	}

	if _, err := api.RegisterChatroomService(); err != nil {
		panic(err.Error())
	}

	if _, err := game.RegisterSampleLandService(); err != nil {
		panic(err.Error())
	}

	endpoints.HandleHTTP()

}
