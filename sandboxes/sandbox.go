package main

import (
	"../src/NonchangGameServer/model/stamina"
	"fmt"
	// "math"
	"time"
)

func main() {

	l("\n\n========= サンドボックス開始\n\n")

	// a := &TestObject{Name: "123", Text: "456"}

	// fmt.Printf(a.Name + "\n")
	// _, b := a.TestFunc()
	// fmt.Printf(b + "\n")

	//===========================
	//スタミナ実装テスト

	l("- スタミナ実装テスト\n")

	ut := stamina.NewUtil(30, 5)

	{

		savedFlushUnixtime := time.Now().Add(-4 * time.Second) //n秒前としてみる
		currentStamina, nextHealSec := ut.GetCurrentStatuses(savedFlushUnixtime)

		l("\t Stamina=%+v\n", currentStamina) //今のスタミナ表示
		l("\t あと%+v秒で回復します。\n", nextHealSec)  //今のスタミナ表示
	}

	{
		maxtime := ut.GetJustMaxFlushAt()
		savedFlushUnixtime := ut.GetStaminaDownedFlushTime(maxtime, 5)
		savedFlushUnixtime = savedFlushUnixtime.Add(-1 * time.Second)
		currentStamina, nextHealSec := ut.GetCurrentStatuses(savedFlushUnixtime)
		l("\t Stamina=%+v\n", currentStamina)
		l("\t あと%+v秒で回復します。\n", nextHealSec) //あれっ。おかしい。
	}

	l("\n========= サンドボックス単体実行完了\n\n")
}

// type TestObject struct {
// 	Name string
// 	Text string
// }

// func (to *TestObject) TestFunc() (a, b string) {
// 	return "abc", "defg"
// }

//手抜き用 : fmt.Printfを毎回打つのがだるいので。
func l(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf(format, a...)
}
