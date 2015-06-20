package game

/*

ゲームAPI「sampleland」
	- 某xxランドを元に、ダンジョンを探索するだけのシンプルなゲームサービスです。
	-

*/

import (
	"errors"
	"log"
	// "fmt"
	"strings"
	"time"

	// "appengine"
	"appengine/datastore"
	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
	"github.com/mjibson/goon"

	// "NonchangGameServer/utils"
	// "NonchangGameServer/api"
	"NonchangGameServer/api/common"
	"NonchangGameServer/model"
	samplelandmodel "NonchangGameServer/model/game/sampleland"
	"NonchangGameServer/model/stamina"
)

var (
	UUIDNotFoundError = errors.New("no user found.")
)

func l(format string, a ...interface{}) {
	//TODO - 可変長引数の数に合わせて、自動でformatの中で「%+v,」として追加してやりたいなぁ。JavaScriptのconsole.log同等にしたい。
	log.Printf("==[at sampleland]== "+format, a...)
}

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
		mtd := rpc.MethodByName(orig)
		if mtd == nil {
			panic("Missing method " + orig)
		}
		i := mtd.Info()
		i.Name, i.HTTPMethod, i.Path, i.Desc = name, method, path, desc
		// あれ？このメソッドは結局なにをやっているんだ？
		// mもiもどこにも行ってないよなー……。
		// i(model.Info())に代入することで初期化が完了しているのかな？
	}
	//メモ：項目の二つ目は動作には影響がない模様。しかし、同一RPC上で同じ名前があるとエラーになる。。
	// 一つ目が対象メソッド。四つ目はブラウザにどうマップされるか。
	register("Ping", "ping", "GET", "game/sampleland/ping",
		"サーバ応答を確認します。常にresult:trueを返します。")

	register("GetPlayer", "player", "GET", "game/sampleland/player",
		"ゲーム用プレイヤー情報取得")

	//以下は、通常はゲーム中に公開しないAPIです。
	//- スタミナ操作はあくまでゲームルール側に応じる形で行われるものです。
	//- ポイント消費でスタミナ回復などの場合も、回復API側で処理します。
	register("Dev_DecrementStamina", "dev_decrementStamina", "GET", "game/sampleland/dev_decrementStamina",
		"開発テスト用：行動力を1減らします。")
	// register("Dev_IncrementStamina", "incrementStamina", "GET", "game/sampleland/dev_incrementStamina",
	// 	"開発テスト用：行動力を1増やします。")
	// register("Dev_ResetStamina", "incrementStamina", "GET", "game/sampleland/dev_resetStamina",
	// 	"開発テスト用：行動力を完全回復させます。")

	return rpc, nil
}

//サービス定義開始
//メモ：エンドポイント公開分は、大文字構造体・大文字メソッド（公開）である必要があります。
type SampleLandService struct{}

//Pingサービス

func (sv *SampleLandService) Ping(c endpoints.Context) (*common.SimpleResult, error) {
	return &common.SimpleResult{IsSuccess: true}, nil
}

func (sv *SampleLandService) GetPlayer(c endpoints.Context, req *common.UUIDRequest) (*samplelandmodel.SampleLandPlayer, error) {

	if req.UUID == "" {
		return nil, errors.New("======== GetPlayer : リクエストにUUIDが含まれていません。jsonマップ名は`uuid`となっていますか？")
	}

	goon := goon.FromContext(c)
	settings := samplelandmodel.GetSettingInstance()
	parentKey := datastore.NewKey(c, "Player", req.UUID, 0, nil)

	data, _, err := sv.GetPlayerRecordByUUID(c, req.UUID, *parentKey)
	if err == nil {
		l("======== DEBUG - 取得できました＾＾ 1 ")
		return data, nil
	} else if err != UUIDNotFoundError {
		//もし「見つからなかった」エラーじゃなければハンドリングしようがないのでerrを返します。
		return nil, err
	}

	//見つからなかった==新規ユーザアクセス。ここで新しいレコードを作ります

	// l("======== TODO - 新規レコード作って返す")
	data.ParentKey = parentKey
	data.CreatedAt = time.Now()
	// data.Experience = 0

	//プレイヤー名を取得
	playerData := new(model.Player)
	getError := datastore.Get(c, data.ParentKey, playerData)
	if getError != nil {
		l("======== ゲーム初回アクセスエラー：プレイヤーがマスタの方にいませんでした。。 UUID=%v", req.UUID)
		return nil, getError
	}
	data.Name = playerData.Name

	//goonで保存
	if _, err := goon.Put(data); err != nil {
		//もし保存に失敗したら、念のためデータは返さない
		l("======== ゲーム初回アクセスエラー：新規データのgoon保存に失敗しました。。。 %v", data)
		return nil, err
	}

	//保存完了、作成値に応答値を加えてそのまま返す。
	data.Stamina = settings.MaxStamina
	data.StaminaNextHealSec = -1
	data.Settings = settings
	data.IsSuccess = true
	l("======== DEBUG 初回ユーザ作成 : %v ", data)
	return data, nil

}

//共通部分引き出し
func (sv *SampleLandService) GetPlayerRecordByUUID(c endpoints.Context, uuid string, parentKey datastore.Key) (*samplelandmodel.SampleLandPlayer, *stamina.Stamina, error) {

	//メモ:呼び元と同じことを多重処理してるけど大丈夫かしら？
	goon := goon.FromContext(c)
	settings := samplelandmodel.GetSettingInstance()
	data := &samplelandmodel.SampleLandPlayer{UUID: uuid, ParentKey: &parentKey}
	err := goon.Get(data)
	if err == nil {
		//取得できたらそのまま返す
		l("======== GetPlayerRecordByUUID DEBUG - 取得できました＾＾ 2 ")

		//StaminaFlushAtからスタミナ算出
		staminaObj := stamina.New(data.StaminaFlushAt, settings.MaxStamina, settings.StaminaHealSec)
		staminaPoint, nextHealSec := staminaObj.Get()
		l("======== GetPlayerRecordByUUID 側Get直後 : %+v,%+v", staminaPoint, nextHealSec)
		data.Stamina = staminaPoint
		data.StaminaNextHealSec = nextHealSec

		data.Settings = settings
		data.IsSuccess = true
		return data, staminaObj, nil
	}

	//もし「見つからなかった」エラーじゃなければハンドリングしようがないのでerrを返す
	if !strings.Contains(err.Error(), "goon: cannot find a key for struct") &&
		!strings.Contains(err.Error(), "datastore: no such entity") {
		l("======== goon.Getでエラーです。 %v", err.Error())
		return nil, nil, err
	}

	//見つからなかったエラー
	return nil, nil, UUIDNotFoundError
}

func (sv *SampleLandService) Dev_DecrementStamina(c endpoints.Context, req *common.UUIDRequest) (*samplelandmodel.SampleLandPlayer, error) {
	l("======== Dev_DecrementStaminaに来ました。")

	// l("======== Dev_DecrementStamina ************* テスト開始 *************** ")
	// l("======== Dev_DecrementStamina ************* 今何時？ *************** %+v", time.Now())
	// // settings := samplelandmodel.GetSettingInstance()
	// // data := &samplelandmodel.SampleLandPlayer{UUID: req.UUID, ParentKey: (datastore.NewKey(c, "Player", req.UUID, 0, nil))}
	// // staminaObj := stamina.New(data.StaminaFlushAt, settings.MaxStamina, settings.StaminaHealSec)
	// // l("======== Dev_DecrementStamina ************* staminaObj *************** %+v", staminaObj)

	// sss := stamina.NewUtil(20, 2)
	// sss2, sss3 := sss.GetCurrentStatuses(time.Now().Add(-time.Second * 18))
	// l("======== Dev_DecrementStamina ************* sss *************** %+v%+v", sss2, sss3)

	// l("======== Dev_DecrementStamina ************* テスト終了 *************** ")

	//以下、実際のコード
	// l("======== Dev_DecrementStamina ************* 以下、実際のコード *************** ")

	goon := goon.FromContext(c)
	// settings := samplelandmodel.GetSettingInstance()
	parentKey := datastore.NewKey(c, "Player", req.UUID, 0, nil)
	data, staminaObj, err := sv.GetPlayerRecordByUUID(c, req.UUID, *parentKey)
	if err != nil {
		return nil, err
	}
	// staminaObj := staminaObj.New(data.StaminaFlushAt, settings.MaxStamina, settings.StaminaHealSec)

	// l("======== Dev_DecrementStamina 1 %+v", staminaObj)
	a, b := staminaObj.Get()
	_ = a
	_ = b
	staminaObj.Add(-2)
	l("======== Dev_DecrementStamina 側add直後 : %+v,%+v", a, b)
	//応答フィールドに代入
	// l("======== Dev_DecrementStamina 側get直前")
	staminaPoint, nextHealSec := staminaObj.Get()
	l("======== Dev_DecrementStamina 側Get直後 : %+v,%+v", staminaPoint, nextHealSec)
	data.Stamina = staminaPoint
	data.StaminaNextHealSec = nextHealSec

	// l("======== Dev_DecrementStamina 1.1 %+v", staminaObj)//確かに増えてる
	data.StaminaFlushAt = staminaObj.FlushAt
	if _, err := goon.Put(data); err != nil {
		//もし保存に失敗したら、念のためデータは返さない
		l("======== スタミナ保存エラー：データのgoon保存に失敗しました。。。 %v", data)
		return nil, err
	}
	// l("======== Dev_DecrementStamina 2 (増えてる？) %+v", staminaObj)//増えてる

	return data, nil
}
