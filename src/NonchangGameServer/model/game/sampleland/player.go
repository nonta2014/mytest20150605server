package sampleland

/*

- プレイヤーのデータ、基本的操作、リスト表現など
- API側（account.go）でも5,6行で終わるDatastoreアクセスは直接実装してます。
	- どういう基準でmodelとapiコードを分離するかは、随時要検討。

*/

import (
	// "appengine"
	"appengine/datastore"
	// "github.com/mjibson/goon"
	"time"
)

//goonを使うと構造体名がKindになります。ご注意。
type SampleLandPlayer struct {
	UUID      string         `json:"-" datastore:"-" endpoints="req" goon:"id"`
	ParentKey *datastore.Key `json:"-" datastore:"-" endpoints="req" goon:"parent"`
	CreatedAt time.Time      `json:"createdAt" `
	Level     int            `json:"level" endpoints="req,d=1"`
	Stamina   int            `json:"stamina" endpoints="d=30"`
	Name      string         `json:"name" `
	//DataStoreに入れない系
	//応答用
	IsSuccess bool `json:"success" datastore:"-"`
}

//プレイヤーの一覧を取得する動線は一旦なくなりました。
// type PlayerList struct {
// 	Players []*Player `json:"players"`
// }
