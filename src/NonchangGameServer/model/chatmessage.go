package model

/*

- プレイヤーのデータ、基本的操作、リスト表現など
- API側（account.go）でも5,6行で終わるDatastoreアクセスは直接実装してます。
	- どういう基準でmodelとapiコードを分離するかは、随時要検討。

*/

import (
	"appengine"
	"appengine/datastore"
	"time"
)

type ChatMessage struct {
	Key        *datastore.Key `json:"-" datastore:"-"`
	RoomKey    string         `json:"roomKey"`
	PostedAt   time.Time      `json:"postedAt" `
	Message    string         `json:"message" endpoints="req,noindex"`
	PlayerUUID string         `json:"playerUUID" endpoints="req"`
	//応答用
	IsSuccess bool `json:"success" datastore:"-"`
}

type ChatMessageList struct {
	//応答用構造
	Messages []*ChatMessage `json:"messages"`
	AllCount int            `json:"allCount"`
}

//TODO - カーソル指定ができないと、件数が増えた時に使い物にならない。

func AllChatMessages(c appengine.Context, order string, limit int) (*ChatMessageList, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit >= 10 {
		limit = 10
	}
	// c.Infof("\n============ order=%s,limit=%s ===========\n", order, limit)
	q := datastore.NewQuery("ChatMessage").Order(order).Limit(limit)
	// q := datastore.NewQuery("ChatMessage").Limit(limit)
	datas := make([]*ChatMessage, 0, limit)
	keys, err := q.GetAll(c, &datas)
	if err != nil {
		return nil, err
	}
	for i, k := range keys {
		datas[i].Key = k
	}

	// c.Infof("\n============ %+v", datas)

	//全件数も取得
	n, err := datastore.NewQuery("ChatMessage").Count(c)
	if err != nil {
		return nil, err
	}

	return &ChatMessageList{Messages: datas, AllCount: n}, nil
}
