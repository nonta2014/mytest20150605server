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
	// slData "NonchangGameServer/model/game/sampleland/_generatedDatas"
	"time"
)

//メモ：goonを使うと構造体名がKindになります。ご注意。
//（※指定方法がないわけではないけど）
// また、JSON名の指定は無意味に感じたのでスルーすることにしました。
type SampleLandPlayer struct {
	UUID string `json:"-" datastore:"-" goon:"id" endpoints="req"`
	//parentKeyとgoon idになっている値が被っているのは、なんかおかしい気がしないでもない。。
	ParentKey *datastore.Key `json:"-" datastore:"-" goon:"parent" endpoints="req"`
	//名前はベースと重複で保存
	Name string `datastore:",noindex"`
	//一応作成日時を保存。ユーザの総プレイ時間などに必要か？
	CreatedAt time.Time `json:"-"`

	//メモ:LevelはExperienceから算出可能だからいらないと判断。一旦おいとく
	// Level                 int            `json:"level" endpoints="req,d=1"`
	Experience int `datastore:",noindex" endpoints="d=1"`

	//スタミナが「最後に0になった日時」。これを保存することで、スタミナ回復判定を一元化できる。
	StaminaFlushAt time.Time `json:"-" datastore:",noindex" endpoints="-"`

	//マッチング用（？備忘録がてらの空実装。マッチングはどう実装したもんかなぁ。）
	PartyAttackTotal int

	//ゲーム情報
	CurrentFloor  int
	PossibleFloor int     //これより上のフロアには進めない。フロアボスを倒すと先に進める。
	Coin          int     //ゲーム内通貨
	FloorProgress float32 //探索率。0.0f-1.0f

	//ゲーム設定情報を含めてみる。
	//※これはAPIにチャンクを実装したら、その仕様で分離したい。Playerが持つべき情報じゃないし。。
	// Settings *sampleLandSetting `datastore:"-"`

	//応答用など
	//メモ：timeオブジェクトはjsに返しても`2015-06-15T14:30:14.767669Z`と処理しにくい形式になるので、PHP時代同様にUnixtimeで処理することにしています。（ベストプラクティスかどうかは微妙なところ）
	Stamina            int  `datastore:"-" endpoints="-"`
	StaminaNextHealSec int  `datastore:"-" endpoints="-"`
	IsSuccess          bool `datastore:"-" endpoints="-`
}

//
//
//
//
//

// ■＿＿＿■＿■＿＿■■■＿＿■■＿＿
// ■■＿■■＿■＿■＿＿＿＿■＿＿■＿
// ■＿■＿■＿■＿＿■■＿＿■＿＿＿＿
// ■＿＿＿■＿■＿＿＿＿■＿■＿＿■＿
// ■＿＿＿■＿■＿■■■＿＿＿■■＿＿

//以下、ファイル分けようか検討中。
//複合応答系。

type SampleLandPlayerAndGameSetting struct {
	//応答用
	//ゲーム設定とプレイヤーステータスを同時に返します。
	PlayerData *SampleLandPlayer
	Settings   *sampleLandSetting
}

type SampleLandExploreResult struct {
	//応答用
	//探検結果とプレイヤーステータスを返します。

	//状況に応じたresultを返します。
	PlayerData    *SampleLandPlayer
	ExploreResult *ExploreResult

	//以下は没メモ。Endpointsで「状況により構造を変えるJSON」の提供方法が不明。
	// ExploreResultCoinGet ExploreResultCoinGet `json:"ResultCoinGet" endpoints="d=null"`
	// ExploreResultBattle  ExploreResultBattle  `json:"ResultBattle" endpoints="d=null"`
}

type ExploreResult struct {
	ResultType string

	//コインアップリザルト用
	Coin int

	//敵登場リザルト用
	// EnemyData  slData.Enemy  //これやるとパニクった。i18n含んだからかな……。
	EnemyImageName string
	EnemyName      string
	Experience     int
}

//Result種別をまとめるインターフェイス
//※これやるとendpointがトラブった様子。「状況により構造を変えるJSON」を返す方法があるのかどうかは後日また別口で調査したいところ。
// type ExploreResult interface {
// 	GetResultType() string
// }

const (
	ExploreResultStaminaShortStr = "ExploreResultStaminaShort"
	ExploreResultCoinGetStr      = "ExploreResultCoinGet"
	ExploreResultBattleStr       = "ExploreResultBattle"
)

//コインゲットリザルト
// type ExploreResultCoinGet struct {
// 	Coin int
// }

//冒険 - バトルリザルト（検討中）
// type ExploreResultBattle struct {
// 	// Prize      sampleLandData.Prize
// 	Experience int
// }
