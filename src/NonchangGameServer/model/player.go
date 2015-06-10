package model

/*

- プレイヤーのデータ、基本的操作、リスト表現など
- API側（account.go）でも5,6行で終わるDatastoreアクセスは直接実装してます。
	- どういう基準でmodelとapiコードを分離するかは、随時要検討。

*/

import (
	// "appengine"
	"appengine/datastore"
	"time"
)

type Player struct {
	//メモ：キーはNewKeyで渡すので、datastore指定自体は"-"としてスルーします。
	Key         *datastore.Key `json:"-" datastore:"-"`
	CreatedAt   time.Time      `json:"createdAt" `
	LastLoginAt time.Time      `json:"lastLoginAt" `
	Level       int            `json:"level" endpoints="req,d=1"`
	Name        string         `json:"name" `
	//DataStoreに入れない系
	UUID string `json:"uuid" datastore:"-" endpoints="req"`
	//応答用
	IsSuccess bool `json:"success" datastore:"-"`
}

//プレイヤーの一覧を取得する動線は一旦なくなりました。
// type PlayerList struct {
// 	Players []*Player `json:"players"`
// }
