package api

/*

チャットルームAPI

- account.goでUUIDログインしたユーザのみ発言可能です。

*/

import (
	"appengine/datastore"
	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
	"time"

	"NonchangGameServer/api/common"
	"NonchangGameServer/model"
)

func RegisterChatroomService() (*endpoints.RPCService, error) {
	chatroomService := &ChatroomService{}
	rpcService, err := endpoints.RegisterService(
		chatroomService,
		"chatroom",
		"v1",
		"チャットルーム機能を提供します。",
		true)
	if err != nil {
		return nil, err
	}
	register := func(orig, name, method, path, desc string) {
		m := rpcService.MethodByName(orig)
		if m == nil {
			panic("Missing method" + orig)
		}
		i := m.Info()
		i.Name, i.HTTPMethod, i.Path, i.Desc = name, method, path, desc
	}

	//メモ：項目の二つ目は動作には影響がない模様。しかし、同じ名前があるとエラーになる。。
	// 一つ目が対象メソッド。四つ目はブラウザにどうマップされるか。
	register("Ping", "ping", "GET", "chatroom/ping",
		"サーバ応答を確認します。常にresult{success:true}を返します。")

	// register("SignUp", "signup", "GET", "account/signup",
	// 	"UUIDと名前を受け取って新規ユーザとして登録します。")
	// register("Login", "login", "GET", "account/login",
	// 	"UUIDでログインします。")

	// //初期のコピペ分。countはとりあえず使えそう。listはoffsetがないと使い物にならないような……？
	register("List", "list", "GET", "chatroom/list",
		"チャット一覧を取得します。")
	register("Add", "add", "PUT", "chatroom/add",
		"メッセージを追加します。")
	register("Count", "count", "GET", "chatroom/count",
		"チャットメッセージの総数を取得します。")

	return rpcService, nil
}

//サービス定義開始
type ChatroomService struct{}

func (sv *ChatroomService) Ping(c endpoints.Context) (*common.SimpleResult, error) {
	return &common.SimpleResult{IsSuccess: true}, nil
}

//TODO - カーソル化
func (sv *ChatroomService) List(c endpoints.Context, r *common.LimitAndCursorRequest) (*model.ChatMessageList, error) {
	// result, err := model.AllChatMessages(c, "-PostedAt", r.Limit)
	result, err := model.AllChatMessages2(c, "-PostedAt", r.Limit, r.Cursor)
	if err != nil {
		c.Infof("\n======== AllChatMessages2でエラーです。 %v\n", err)
		return nil, err
	}
	return result, nil
}

func (sv *ChatroomService) Count(c endpoints.Context) (*common.CountResult, error) {
	n, err := datastore.NewQuery("ChatMessage").Count(c)
	if err != nil {
		return nil, err
	}
	return &common.CountResult{n}, nil
}

func (ps *ChatroomService) Add(ctx endpoints.Context, addData *model.ChatMessage) (*model.ChatMessage, error) {
	k := datastore.NewIncompleteKey(ctx, "ChatMessage", nil)

	// ctx.Infof("addDataチェック : ", addData.Name, 123)

	//初期値代入（リクエストで渡されても無視）
	addData.RoomKey = "default"
	addData.PostedAt = time.Now()

	_, err := datastore.Put(ctx, k, addData)

	if err == nil {
		addData.IsSuccess = true
		return addData, err
	} else {
		addData.IsSuccess = false
		return addData, err
	}
}
