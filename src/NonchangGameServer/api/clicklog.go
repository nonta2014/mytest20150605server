package api

import (
	// "appengine"
	"appengine/datastore"
	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
	"time"
)

// クリックログendpoint
// - アクセス例 : <http://localhost:8080/_ah/api/clicklog/v1/clicklogs>

func RegisterClickLogService() (*endpoints.RPCService, error) {
	// rpcService, err := endpoints.RegisterService(
	// 	&ClickLogService{},
	// 	"clicklog", "v1", "クリック座標を取るサンプルです。", true)
	// if err != nil {
	// 	return nil, err
	// }

	// info := rpcService.MethodByName("BoardGetMove").Info()
	// info.Path, info.HTTPMethod, info.Name = "board", "POST", "board.getmove"
	// info.Scopes, info.ClientIds, info.Audiences = scopes, clientIds, audiences

	// info = rpcService.MethodByName("ScoresList").Info()
	// info.Path, info.HTTPMethod, info.Name = "scores", "GET", "scores.list"
	// info.Scopes, info.ClientIds, info.Audiences = scopes, clientIds, audiences

	// info = rpcService.MethodByName("ScoresInsert").Info()
	// info.Path, info.HTTPMethod, info.Name = "scores", "POST", "scores.insert"
	// info.Scopes, info.ClientIds, info.Audiences = scopes, clientIds, audiences

	// return rpcService, nil

	clickLogService := &ClickLogService{}
	rpcService, err := endpoints.RegisterService(
		clickLogService,
		"clicklog",
		"v1",
		"クリック座標を取るサンプルです。",
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
	register("List", "clicklogs.list", "GET", "clicklogs", "List most recent.")
	register("Add", "clicklogs.add", "PUT", "clicklogs", "Add item.")
	register("Count", "clicklogs.count", "GET", "clicklogs/count", "Count all.")

	return rpcService, nil

}

//ClickLog Endpoints

//ログレコード本体
type ClickLog struct {
	Key  *datastore.Key `json:"id" datastore:"-"`
	Time time.Time      `json:"time"`
	X    float64        `json:"x" endpoints:"req"`
	Y    float64        `json:"y" endpoints:"req"`
	// TestName string `json:"testname" endpoints:"req"`
}
type ClickLogList struct {
	Items []*ClickLog `json:"items"`
}
type ClickLogListReq struct {
	Limit int `json:"limit" endpoints="d=10"` //d=default
}

type ClickLogService struct{}

//メモ：以下はClickLogServiceのメソッドとなるみたい。言語仕様がまだふわふわしてる。
func (s *ClickLogService) List(c endpoints.Context, r *ClickLogListReq) (*ClickLogList, error) {
	if r.Limit <= 0 {
		r.Limit = 10
	}
	q := datastore.NewQuery("ClickLog").Order("-Time").Limit(r.Limit)
	logs := make([]*ClickLog, 0, r.Limit)
	keys, err := q.GetAll(c, &logs)
	if err != nil {
		return nil, err
	}
	for i, k := range keys {
		logs[i].Key = k
	}
	return &ClickLogList{logs}, nil
}
func (s *ClickLogService) Add(ctx endpoints.Context, clicklog *ClickLog) error {
	k := datastore.NewIncompleteKey(ctx, "ClickLog", nil)
	_, err := datastore.Put(ctx, k, clicklog)
	return err
}

func (s *ClickLogService) Count(ctx endpoints.Context) (*Count, error) {
	n, err := datastore.NewQuery("ClickLog").Count(ctx)
	if err != nil {
		return nil, err
	}
	return &Count{n}, nil
}

//
//
//
//
//
//
//
//
//
//

//以下は、Goon使おうとしたときのテスト。ちょっと基礎と情報が足りないので後回し。
// 参考 : <http://qiita.com/soundTricker/items/194d4067b0e145544b56>

// type Worktime struct {
// 	Ymd       string         `datastore:"-" goon:"id"`
// 	MemberKey *datastore.Key `datastore:"-" goon:"parent"`

// 	StartTime     string     `datastore:"startTime,noindex"`
// 	EndTime       string     `datastore:"endTime,noindex"`
// 	RestTime      string     `datastore:"lateRestTime,noindex"`
// 	Etc           string     `datastore:"etc,noindex"`
// 	Worktime      int        `datastore:"worktime,noindex"`
// }

// func makeGoon(r *http.Request) goon.Goon {
// 	g := goon.NewGoon(r)
// 	return g
// }

// func PutWorktime(r *http.Request, memberKey *datastore.Key) (*Worktime, error) {
// 	g := goon.NewGoon(r)

// 	//goonタグの値以外を作る ※サンプルコードのため省略
// 	// w, err := createWorktime(r)
// 	w :=new Worktime

// 	if err != nil {
// 		return nil, err
// 	}

//     //Parent Keyを詰める
// 	w.MemberKey = memberKey

//     //Datastore KeyのNameを詰める
// 	w.Ymd = "2014-12-24"

//     //g.PutでPut KeyやKind名はGoon側で解決
//     //またPut時にMemcacheやIn-Memory Cacheしてくれる
//     //なお返却値は *datastore.Key, error
// 	if _, err := g.Put(w); err != nil {
// 		return w, err
// 	}

// 	return w, nil
// }
