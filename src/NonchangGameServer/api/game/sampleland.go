package game

/*

ゲームAPI「sampleland」
	- 某xxランドを元に、ダンジョンを探索するだけのシンプルなゲームサービスです。
	-

*/

import (
	"errors"
	// "fmt"
	"strings"
	"time"

	// "appengine"
	"appengine/datastore"
	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
	"github.com/mjibson/goon"

	"NonchangGameServer/utils"
	// "NonchangGameServer/api"
	"NonchangGameServer/api/common"
	m "NonchangGameServer/model"
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

	register("Dev_DecrementStamina", "decrementStamina", "GET", "game/sampleland/dev/decrementStamina",
		"開発テスト用：行動力を1減らします。")

	return rpc, nil
}

//サービス定義開始
//メモ：エンドポイント公開分は、大文字構造体・大文字メソッド（公開）である必要があります。
type SampleLandService struct{}

//Pingサービス

func (sv *SampleLandService) Ping(c endpoints.Context) (*common.SimpleResult, error) {
	return &common.SimpleResult{IsSuccess: true}, nil
}

// func isExistPlayer(c endpoints.Context, g *goon.Goon, uuid string) (bool, error) {

// 	n, err := g.Count(datastore.NewQuery("SampleLandPlayer"))
// 	if err != nil {
// 		c.Infof("\n======== g.Countでエラー : %+v\n", err)
// 		return false, err
// 	}
// 	c.Infof("\n======== isExistPlayer : %v\n", n)
// 	if n >= 2 {
// 		return false, errors.New(fmt.Sprintf("isExistPlayer - 該当プレイヤーレコードが複数ありました。。件数=%v", n))
// 	}
// 	return (n == 1), nil
// }

func (sv *SampleLandService) GetPlayer(c endpoints.Context, req *common.UUIDRequest) (*mm.SampleLandPlayer, error) {

	if req.UUID == "" {
		return nil, errors.New("\n\n======== GetPlayer : リクエストにUUIDが含まれていません。jsonマップ名は`uuid`となっていますか？\n\n")
	}

	goon := goon.FromContext(c)

	//存在確認
	// if exist, err := isExistPlayer(c, goon, req.UUID); err != nil {
	// 	c.Infof("\n======== isExistPlayerでエラーです。 %v\n", err)
	// 	return nil, err
	// } else {
	// 	c.Infof("\n======== GetPlayer 存在したー %+v", exist)
	// }

	settings := mm.GetSettingInstance()
	parentKey := datastore.NewKey(c, "Player", req.UUID, 0, nil)

	data := &mm.SampleLandPlayer{UUID: req.UUID, ParentKey: parentKey}
	err := goon.Get(data)
	if err == nil {
		//取得できたらそのまま返す
		c.Infof("\n\n======== DEBUG - 取得できました＾＾ \n\n")

		//スタミナ計算
		stamina, nextHealSec := stamina.New(settings.MaxStamina, settings.StaminaHealSec).GetCurrentStatuses(data.StaminaFlushAt.Unix())
		data.Stamina = stamina
		data.StaminaNextHealSec = nextHealSec

		data.Settings = settings
		data.IsSuccess = true
		return data, nil
	}

	//もし「見つからなかった」エラーじゃなければハンドリングしようがないのでerrを返す
	if !strings.Contains(err.Error(), "goon: cannot find a key for struct") &&
		!strings.Contains(err.Error(), "datastore: no such entity") {
		c.Infof("\n\n======== goon.Getでエラーです。 %v\n\n", err.Error())
		return nil, err
	}

	//見つからなかった==新規ユーザアクセス。ここで新しいレコードを作ります

	// c.Infof("\n======== TODO - 新規レコード作って返す")
	data.ParentKey = parentKey
	data.CreatedAt = time.Now()
	// data.Experience = 0

	//プレイヤー名を取得
	playerData := new(m.Player)
	getError := datastore.Get(c, data.ParentKey, playerData)
	if getError != nil {
		c.Infof("\n\n======== ゲーム初回アクセスエラー：プレイヤーがマスタの方にいませんでした。。 UUID=%v\n\n", req.UUID)
		return nil, getError
	}
	data.Name = playerData.Name

	//goonで保存
	if _, err := goon.Put(data); err != nil {
		//もし保存に失敗したら、念のためデータは返さない
		c.Infof("\n\n======== ゲーム初回アクセスエラー：新規データのgoon保存に失敗しました。。。 %v\n\n", data)
		return nil, err
	}

	//保存完了、作成値をそのまま返す。
	data.Settings = settings
	data.IsSuccess = true
	c.Infof("\n\n======== DEBUG 初回ユーザ作成 : %v\n\n", data)
	return data, nil

}

func (sv *SampleLandService) Dev_DecrementStamina(c endpoints.Context, req *common.UUIDRequest) (*common.SimpleResult, error) {
	//TODOTODOTODO
	return &common.SimpleResult{IsSuccess: true}, nil
}
