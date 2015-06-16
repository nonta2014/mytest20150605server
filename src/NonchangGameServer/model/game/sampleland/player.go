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
	"time"
)

//メモ：goonを使うと構造体名がKindになります。ご注意。
//（※指定方法がないわけではないけど）
// また、JSON名の指定は無意味に感じたのでスルーすることにしました。
type SampleLandPlayer struct {
	UUID      string         `json:"-" datastore:"-" goon:"id" endpoints="req"`
	ParentKey *datastore.Key `json:"-" datastore:"-" goon:"parent" endpoints="req"`
	CreatedAt time.Time      `json:"-"`
	//メモ:LevelはExperienceから算出可能だからいらないかも。一旦おいとく
	// Level                 int            `json:"level" endpoints="req,d=1"`
	Experience     int       `datastore:",noindex" endpoints="d=1"`
	StaminaFlushAt time.Time `json:"-" datastore:",noindex" endpoints="-"`
	Name           string    `datastore:",noindex"`

	//マッチング用（？備忘録がてらの空実装。マッチングはどう実装したもんかなぁ。）
	PartyAttackTotal int

	//設定含めてみる。
	//※これはAPIにチャンクを実装したら、その仕様で分離したい。Playerが持つような情報じゃないし。。
	Settings *sampleLandSetting `datastore:"-"`

	//応答用など
	//メモ：timeオブジェクトはjsに返しても`2015-06-15T14:30:14.767669Z`と処理しにくい形式になるので、PHP時代同様にUnixtimeで処理することにしています。（ベストプラクティスかどうかは微妙なところ）
	Stamina            int  `datastore:"-" endpoints="-"`
	StaminaNextHealSec int  `datastore:"-" endpoints="-"`
	IsSuccess          bool `datastore:"-" endpoints="-`
}
