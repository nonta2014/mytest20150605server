package api

/*

アカウントAPI

- ping応答（Ping）
	- ※Signup画面表示前のために用意。本来はシステムAPIに分けるべきところを手抜き。
- 新規ユーザ登録（Signup）
- uuidによるログイン＋ping応答（Login）
	- 自分のステータスを取得


- APIリストメモ（TODO）

	- Pingサービス（成功）
		- <http://localhost:8080/_ah/api/account/v1/ping>


*/

import (
	// "appengine"
	"appengine/datastore"
	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
	"time"

	"NonchangGameServer/api/common"
	"NonchangGameServer/model"
)

func AccountTest() string {
	return "account test"
}

//endpoint登録

func RegisterAccountService() (*endpoints.RPCService, error) {
	accountService := &AccountService{}
	rpcService, err := endpoints.RegisterService(
		accountService,
		"account",
		"v1",
		"アカウント関係サービスです。",
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
		// あれ？このメソッドは結局なにをやっているんだ？
		// mもiもどこにも行ってないよなー……。
		// i(m.Info())に代入することで初期化が完了しているのかな？
	}
	//メモ：項目の二つ目は動作には影響がない模様。しかし、同じ名前があるとエラーになる。。
	// 一つ目が対象メソッド。四つ目はブラウザにどうマップされるか。
	register("Ping", "ping", "GET", "account/ping",
		"サーバ応答を確認します。常にresult:trueを返します。")

	register("SignUp", "signup", "GET", "account/signup",
		"UUIDと名前を受け取って新規ユーザとして登録します。")

	register("Login", "login", "GET", "account/login",
		"UUIDでログインします。")

	//あれ？これだとgapi.client.account.account.listになっちゃうか。
	// register("List", "account.list", "GET", "account",
	// 	"プレイヤー一覧を取得します。")

	register("Count", "count", "GET", "account/count",
		"プレイヤーの総数を取得します。")

	return rpcService, nil
}

//サービス定義
type AccountService struct{}

//Pingサービス

func (ps *AccountService) Ping(c endpoints.Context) (*common.SimpleResult, error) {
	return &common.SimpleResult{IsSuccess: true}, nil
}

// //以下は実際には不要。offsetできないのかな？ できれば管理用には十分なんだけども。
// func (ps *AccountService) List(ctx endpoints.Context, r *LimitRequest) (*model.PlayerList, error) {
// 	result, err := model.AllPlayers(ctx, "-Time", r.Limit)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return result, nil
// }

//管理系でプレイヤー総数は取得できる必要があるだろう
func (s *AccountService) Count(ctx endpoints.Context) (*common.CountResult, error) {
	n, err := datastore.NewQuery("Player").Count(ctx)
	if err != nil {
		return nil, err
	}
	return &common.CountResult{n}, nil
}

//
//
//

// いよいよ自作APIとして設計開始。とりあえずpingコピーしてみる。

func (ps *AccountService) SignUp(ctx endpoints.Context, addData *model.Player) (*model.Player, error) {
	// k := datastore.NewIncompleteKey(ctx, "Player", nil)
	k := datastore.NewKey(ctx, "Player", addData.UUID, 0, nil) //最後のはparentのKey。Ancestor？調べなきゃ。

	// ctx.Infof("addDataチェック : ", addData.Name, 123)

	//初期値代入（リクエストで渡されても無視）
	addData.Level = 1
	addData.CreatedAt = time.Now()
	addData.LastLoginAt = time.Now()

	_, err := datastore.Put(ctx, k, addData)

	if err == nil {
		addData.IsSuccess = true
		return addData, err
	} else {
		addData.IsSuccess = false
		return addData, err
	}
}

type PlayerGetReq struct {
	UUID string `json:"uuid"`
}

func (ps *AccountService) Login(c endpoints.Context, req *common.UUIDRequest) (*model.Player, error) {
	k := datastore.NewKey(c, "Player", req.UUID, 0, nil)
	playerData := new(model.Player)
	err := datastore.Get(c, k, playerData)
	if err != nil {
		// http.Error(w, err.Error(), 500)
		return &model.Player{IsSuccess: false}, err
	}
	playerData.IsSuccess = true
	playerData.UUID = req.UUID
	return playerData, nil
}
