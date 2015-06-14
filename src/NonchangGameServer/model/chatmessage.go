package model

/*

- プレイヤーのデータ、基本的操作、リスト表現など
- API側（account.go）でも5,6行で終わるDatastoreアクセスは直接実装してます。
	- どういう基準でmodelとapiコードを分離するかは、随時要検討。
	- いっそ、こっちはstructしか持たないくらいのほうがいいのかな……。

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
	Message    string         `json:"message"    endpoints="req,noindex"`
	PlayerUUID string         `json:"playerUUID" endpoints="req"`

	//reqにしたい後から追加分
	//初期になかったため、マイグレーションしないといけないのでreqを外してます。
	//ベストプラクティスを知りたいところ……。
	Name string `json:"name" endpoint="req,noindex"`

	//応答用
	IsSuccess bool `json:"success"       datastore:"-"`
}

type ChatMessageList struct {
	//応答用の構造体です。
	Messages   []*ChatMessage `json:"messages"`
	AllCount   int            `json:"allCount"`
	NextCursor string         `json:"nextCursor"`
}

//カーソル指定版
//TODO - この規模のコードでも、1文字変数は潰したくなるなぁ。3文字はないと検索めんどいw

func AllChatMessages2(c appengine.Context, order string, limit int, cursorString string) (*ChatMessageList, error) {

	//不正防止
	if limit <= 0 {
		limit = 10
	}
	if limit >= 10 {
		limit = 10
	}

	//カーソル位置を適用
	q := datastore.NewQuery("ChatMessage").Order(order)
	if cursorString != "" {
		cursor, err := datastore.DecodeCursor(cursorString)
		if err == nil {
			q = q.Start(cursor)
		}
	}

	//結果格納
	datas := []*ChatMessage{}

	//走査開始
	t := q.Run(c)
	for i := 0; i < limit; i++ {
		var message ChatMessage
		messageKey, err := t.Next(&message)
		// c.Infof("\n============p確認 %+v", p)

		if err == datastore.Done {
			break
		}
		if err != nil {
			c.Errorf("fetching next Person: %v", err)
			break
		}

		//マイグレーション用コード
		//もし読み出したデータにNameが入ってなかったら、それはチャット側にNameを追加するのを忘れていた初期のレコード。
		//見つけたら、その都度PlayerDataから読み出して、こちらにも書き込みます（冗長化）。
		//TODO - 10件程度ではあるけど、取得済みがあったらマップあたりにキャッシュしたほうがいいだろうなぁ。同じプレイヤー名をその都度クエリするのは違う。。
		if message.Name == "" {
			c.Infof("\n============名前がnullでした。マイグレーションを開始します。 %+v", messageKey)
			//Playerレコードから名前取得
			k := datastore.NewKey(c, "Player", message.PlayerUUID, 0, nil)
			playerData := new(Player)
			getError := datastore.Get(c, k, playerData)
			if getError != nil {
				c.Infof("\n======== マイグレーションエラー：getErrorです。。 %v\n", getError)
				return nil, getError
			}
			//ChatMessageにプレイヤー名を保存
			message.Name = playerData.Name
			_, putError := datastore.Put(c, messageKey, &message)
			if putError != nil {
				c.Infof("\n======== マイグレーションエラー：putErrorです。。 %v\n", putError)
				return nil, putError
			}
		}

		datas = append(datas, &message)
	}
	nextCursor, err := t.Cursor()
	if err != nil {
		return nil, err
	}

	//全件数も取得、Endpointの応答に含める
	allCount, err := datastore.NewQuery("ChatMessage").Count(c)
	if err != nil {
		return nil, err
	}
	return &ChatMessageList{Messages: datas, AllCount: allCount, NextCursor: nextCursor.String()}, nil
}

//TODO - カーソル指定ができないと、件数が増えた時に使い物にならない。

//- あれっ。これってもしかして、Endpointに任せるために
//	ChatMessageListクラスのfuncにしたほうが見通しよくない？
//- 違うか？ChatMessageListをサービスにしちゃいけないよな。。

// func AllChatMessages(c appengine.Context, order string, limit int) (*ChatMessageList, error) {
// 	if limit <= 0 {
// 		limit = 10
// 	}
// 	if limit >= 10 {
// 		limit = 10
// 	}
// 	// c.Infof("\n============ order=%s,limit=%s ===========\n", order, limit)
// 	q := datastore.NewQuery("ChatMessage").Order(order).Limit(limit)
// 	// q := datastore.NewQuery("ChatMessage").Limit(limit)
// 	datas := make([]*ChatMessage, 0, limit)
// 	keys, err := q.GetAll(c, &datas)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for i, k := range keys {
// 		datas[i].Key = k
// 	}

// 	// c.Infof("\n============ %+v", datas)

// 	//全件数も取得、Endpointの応答に含める
// 	n, err := datastore.NewQuery("ChatMessage").Count(c)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// return &ChatMessageList{Messages: datas, AllCount: n, Cursor: q.Cursor().String()}, nil
// 	return &ChatMessageList{Messages: datas, AllCount: n}, nil
// }
