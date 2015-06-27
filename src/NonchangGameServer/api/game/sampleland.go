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
	"math/rand"
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
	slData "NonchangGameServer/model/game/sampleland/_generatedDatas"
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

	register("GetPlayerAndSetting", "player", "GET", "game/sampleland/playerAndSetting",
		"ゲーム用プレイヤー情報取得")

	//以下は、通常はゲーム中に公開しないAPIです。
	//- スタミナ操作はあくまでゲームルール側に応じる形で行われるものです。
	//- ポイント消費でスタミナ回復などの場合も、回復API側で処理します。
	register("Dev_DecrementStamina", "dev_decrementStamina", "GET", "game/sampleland/dev_decrementStamina",
		"開発テスト用：行動力を1減らします。")
	register("Dev_IncrementStamina", "dev_incrementStamina", "GET", "game/sampleland/dev_incrementStamina",
		"開発テスト用：行動力を1増やします。")
	register("Dev_ResetStamina", "dev_resetStamina", "GET", "game/sampleland/dev_resetStamina",
		"開発テスト用：行動力を完全回復させます。")

	register("Explore", "explore", "GET", "game/sampleland/explore",
		"現在フロアを探検します。")

	return rpc, nil
}

//サービス定義開始
//メモ：エンドポイント公開分は、大文字構造体・大文字メソッド（公開）である必要があります。
type SampleLandService struct{}

//Pingサービス

func (sv *SampleLandService) Ping(c endpoints.Context) (*common.SimpleResult, error) {
	return &common.SimpleResult{IsSuccess: true}, nil
}

func (sv *SampleLandService) GetPlayerAndSetting(c endpoints.Context, req *common.UUIDRequest) (*samplelandmodel.SampleLandPlayerAndGameSetting, error) {

	if req.UUID == "" {
		return nil, errors.New("======== GetPlayer : リクエストにUUIDが含まれていません。jsonマップ名は`uuid`となっていますか？")
	}

	goon := goon.FromContext(c)
	settings := samplelandmodel.GetSettingInstance()

	data, _, err := sv.GetPlayerRecordByUUID(c, req.UUID)
	if err == nil {
		// l("======== DEBUG - 取得できました＾＾ 1 ")
		// return data, nil
		return &samplelandmodel.SampleLandPlayerAndGameSetting{PlayerData: data, Settings: settings}, nil
	} else if err != UUIDNotFoundError {
		//もし「見つからなかった」エラーじゃなければハンドリングしようがないのでerrを返します。
		return nil, err
	}

	//見つからなかった==新規ユーザアクセス。ここで新しいレコードを作ります

	// l("======== TODO - 新規レコード作って返す")
	data.ParentKey = datastore.NewKey(c, "Player", req.UUID, 0, nil)
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
	migrationPlayerData(data)

	//goonで保存
	if _, err := goon.Put(data); err != nil {
		//もし保存に失敗したら、念のためデータは返さない
		l("======== ゲーム初回アクセスエラー：新規データのgoon保存に失敗しました。。。 %v", data)
		return nil, err
	}

	//保存完了、作成値に応答値を加えてそのまま返す。
	data.Stamina = settings.MaxStamina
	data.StaminaNextHealSec = -1
	data.IsSuccess = true
	l("======== DEBUG 初回ユーザ作成 : %v ", data)
	return &samplelandmodel.SampleLandPlayerAndGameSetting{PlayerData: data, Settings: settings}, nil

}

func (sv *SampleLandService) Explore(c endpoints.Context, req *common.UUIDRequest) (*samplelandmodel.SampleLandExploreResult, error) {

	playerData, staminaObj, err := sv.GetPlayerRecordByUUID(c, req.UUID)
	if err != nil {
		return nil, err
	}

	//スタミナ確認
	staminaPoint, _, err := staminaObj.Get()
	if err != nil {
		return nil, err
	}

	//スタミナは3消費 TODOTODOTODO フロア定義に含めたい
	needStamina := 3

	if staminaPoint < needStamina {
		//スタミナ足りない応答
		return &samplelandmodel.SampleLandExploreResult{
				PlayerData: playerData,
				ExploreResult: &samplelandmodel.ExploreResult{
					ResultType: samplelandmodel.ExploreResultStaminaShortStr,
				},
			},
			nil
	}

	//スタミナを減らします。（goon格納はイベント分岐後）
	staminaObj.Add(-needStamina)
	playerData.StaminaFlushAt = staminaObj.FlushAt

	//探検APIは3種類のランダムイベントを返します。
	rand.Seed(time.Now().UnixNano())
	resultTypes := []string{
		samplelandmodel.ExploreResultBattleStr,
		samplelandmodel.ExploreResultCoinGetStr,
	}
	resultType := resultTypes[rand.Intn(len(resultTypes))]

	if resultType == samplelandmodel.ExploreResultBattleStr {
		//TODO - ここでDatastoregに経験値アップ処理

		//TODO - フロア情報からのランダム取得

		//take1 - ゴースト決め打ち
		// enemy := slData.Enemies[slData.ENEMY_KEY_GHOST]

		//1階決め打ち`TODOTODOTODO`で情報を取得して、そこの定義からランダム
		enemies := slData.Floors[1].Enemies
		enemy := enemies[rand.Intn(len(enemies))]

		//敵の名前をja決め打ち`TODOTODOTODO`で取得
		enemyName, err := enemy.Name.Get("ja")
		if err != nil {
			return nil, err
		}

		//TODO - とりあえず一律で倒したことにします（某リランドで負けたことないし）

		//経験値アップ
		exp := 100 //TODO - 敵定義から取得したい。
		playerData.Experience += exp

		//Datastoreに保存
		goon := goon.FromContext(c)
		if _, err := goon.Put(playerData); err != nil {
			l("======== 探検・バトルで保存エラー %v", err)
			return nil, err
		}

		//結果返納
		return &samplelandmodel.SampleLandExploreResult{
				PlayerData: playerData,
				ExploreResult: &samplelandmodel.ExploreResult{
					ResultType: resultType,
					// EnemyData:  enemy,
					EnemyImageName: enemy.ImageName,
					EnemyName:      enemyName,
					Experience:     exp,
				},
			},
			nil

	} else if resultType == samplelandmodel.ExploreResultCoinGetStr {
		//コインアップ
		coin := 100 //TODO - 敵定義から取得したい。
		playerData.Coin += coin

		//Datastoreに保存
		goon := goon.FromContext(c)
		if _, err := goon.Put(playerData); err != nil {
			l("======== 探検・コインゲットで保存エラー %v", err)
			return nil, err
		}

		//結果返納
		return &samplelandmodel.SampleLandExploreResult{
				PlayerData: playerData,
				ExploreResult: &samplelandmodel.ExploreResult{
					ResultType: resultType,
					Coin:       coin,
				},
			},
			nil

	} else {
		return nil, errors.New("不明なresultTypeです。" + resultType)
	}
}

// ＿■■■＿■＿＿■＿■■■＿＿
// ■＿＿＿＿■＿＿■＿■＿＿■＿
// ＿■■＿＿■＿＿■＿■■■■＿
// ＿＿＿■＿■＿＿■＿■＿＿■＿
// ■■■＿＿＿■■＿＿■■■＿＿

//共通部分引き出し
func (sv *SampleLandService) GetPlayerRecordByUUID(c endpoints.Context, uuid string) (*samplelandmodel.SampleLandPlayer, *stamina.Stamina, error) {

	parentKey := datastore.NewKey(c, "Player", uuid, 0, nil)

	//メモ:goon取得、呼び元と同じことを多重処理してるけど負荷大丈夫かしら？
	goon := goon.FromContext(c)
	settings := samplelandmodel.GetSettingInstance()
	data := &samplelandmodel.SampleLandPlayer{UUID: uuid, ParentKey: parentKey}
	err := goon.Get(data)
	if err == nil {
		//取得できたらそのまま返す
		// l("======== GetPlayerRecordByUUID DEBUG - 取得できました＾＾ 2 ")

		//StaminaFlushAtからスタミナ算出
		staminaObj := stamina.New(data.StaminaFlushAt, settings.MaxStamina, settings.StaminaHealSec)
		staminaPoint, nextHealSec, err := staminaObj.Get()
		if err != nil {
			return nil, nil, err
		}
		data.Stamina = staminaPoint
		data.StaminaNextHealSec = nextHealSec

		// data.Settings = settings
		data.IsSuccess = true
		migrationPlayerData(data)
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

func staminaUpdateSub(sv *SampleLandService, c endpoints.Context, uuid string, addStamina int) (*samplelandmodel.SampleLandPlayer, error) {

	goon := goon.FromContext(c)
	data, staminaObj, err := sv.GetPlayerRecordByUUID(c, uuid)
	if err != nil {
		return nil, err
	}
	staminaObj.Add(addStamina)
	staminaPoint, nextHealSec, err := staminaObj.Get()
	if err != nil {
		return nil, err
	}
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

func migrationPlayerData(data *samplelandmodel.SampleLandPlayer) error {
	//プレイヤーレコードのマイグレーション
	// - あとから追加したmodelデータのうち、初期値をここで判定します。
	// - 取得部分の末尾で呼ぶことで、更新時も一元的に反映させることができます。
	if data.CurrentFloor == 0 {
		data.CurrentFloor = 1
	}
	if data.FloorProgress == 0 {
		data.FloorProgress = 1
	}
	if data.PossibleFloor == 0 {
		data.PossibleFloor = 1
	}
	return nil
}

// ■■■＿＿■■■■＿■＿＿＿■＿
// ■＿＿■＿■＿＿＿＿■＿＿＿■＿
// ■＿＿■＿■■■■＿＿■＿■＿＿
// ■＿＿■＿■＿＿＿＿＿■＿■＿＿
// ■■■＿＿■■■■＿＿＿■＿＿＿

//以下は開発用の直変更APIです。本番時に省く方法を検討中。。

func (sv *SampleLandService) Dev_DecrementStamina(c endpoints.Context, req *common.UUIDRequest) (*samplelandmodel.SampleLandPlayer, error) {

	data, err := staminaUpdateSub(sv, c, req.UUID, -1) //スタミナ減らしてupdate
	if err != nil {
		return nil, err
	}
	return data, nil

}

func (sv *SampleLandService) Dev_IncrementStamina(c endpoints.Context, req *common.UUIDRequest) (*samplelandmodel.SampleLandPlayer, error) {
	data, err := staminaUpdateSub(sv, c, req.UUID, 1) //スタミナ増やしてupdate
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (sv *SampleLandService) Dev_ResetStamina(c endpoints.Context, req *common.UUIDRequest) (*samplelandmodel.SampleLandPlayer, error) {
	settings := samplelandmodel.GetSettingInstance()
	data, err := staminaUpdateSub(sv, c, req.UUID, settings.MaxStamina)
	if err != nil {
		return nil, err
	}
	return data, nil
}
