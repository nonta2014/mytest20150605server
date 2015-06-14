package game

/*

ゲームAPI「sampleland」
	- 某xxランドを元に、ダンジョンを探索するだけのシンプルなゲームサービスです。
	-

*/

import (
	// "appengine"
	// "appengine/datastore"
	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
	"github.com/mjibson/goon"
	// "time"

	// "NonchangGameServer/api"
	"NonchangGameServer/api/common"
	// "NonchangGameServer/model"
	mm "NonchangGameServer/model/game/sampleland"
)

func Test() string {
	return "sampleland test"
}

//endpoint登録

func RegisterSampleLandService() (*endpoints.RPCService, error) {
	serv := &SampleLandService{}
	rpc, err := endpoints.RegisterService(
		serv,
		"sampleland",
		"v1",
		"「サンプルランド」ゲームAPIサービスです。",
		true)
	if err != nil {
		return nil, err
	}
	register := func(orig, name, method, path, desc string) {
		m := rpc.MethodByName(orig)
		if m == nil {
			panic("Missing method " + orig)
		}
		i := m.Info()
		i.Name, i.HTTPMethod, i.Path, i.Desc = name, method, path, desc
		// あれ？このメソッドは結局なにをやっているんだ？
		// mもiもどこにも行ってないよなー……。
		// i(m.Info())に代入することで初期化が完了しているのかな？
	}
	//メモ：項目の二つ目は動作には影響がない模様。しかし、同一RPC上で同じ名前があるとエラーになる。。
	// 一つ目が対象メソッド。四つ目はブラウザにどうマップされるか。
	register("Ping", "ping", "GET", "game/sampleland/ping",
		"サーバ応答を確認します。常にresult:trueを返します。")
	register("GetPlayer", "player", "GET", "game/sampleland/player",
		"ゲーム用プレイヤー情報取得")

	return rpc, nil
}

//サービス定義開始
//メモ：エンドポイント公開分は、大文字構造体・大文字メソッド（公開）である必要があります。
type SampleLandService struct{}

//Pingサービス

func (sv *SampleLandService) Ping(c endpoints.Context) (*common.SimpleResult, error) {
	return &common.SimpleResult{IsSuccess: true}, nil
}

// func (sv *SampleLandService) Count(c endpoints.Context) (*common.CountResult, error) {
// 	// n, err := datastore.NewQuery("ChatMessage").Count(c)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	g := goon.FromContext(c)
// 	n := g.Count(datastore.NewQuery("ChatMessage"))
// 	return &common.CountResult{n}, nil
// }

func (sv *SampleLandService) GetPlayer(c endpoints.Context, req *common.UUIDRequest) (*mm.SampleLandPlayer, error) {
	goon := goon.FromContext(c)
	data := &mm.SampleLandPlayer{UUID: req.UUID}
	if err := goon.Get(data); err != nil {
		return nil, err
	}
	return data, nil
	// return nil, nil
}
