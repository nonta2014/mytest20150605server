package main

//実行方法メモ : go run ./sandboxes/sandbox.go

import (
	// "../src/NonchangGameServer/model/stamina"
	"fmt"
	// "math"
	"math/rand"
	"time"
)

func main() {

	l("\n\n========= サンドボックス開始\n\n")

	// a := &TestObject{Name: "123", Text: "456"}
	// l("test %+v", a)

	// fmt.Printf(a.Name + "\n")
	// _, b := a.TestFunc()
	// fmt.Printf(b + "\n")

	//randの使い方テスト
	// rand.Seed(time.Now(
).UnixNano())
	// l("random test: %+v\n", rand.Intn(100))
	// l("random test: %+v\n", rand.Float32())
	// for i := 0; i < 50; i++ {
	// 	s := rand.Float32()
	// 	if s < 0.333 {
	// 		l("rand is downer. %+v\n", s)
	// 	} else if s < 0.666 {
	// 		l("rand is upper.%+v\n", s)
	// 	} else {
	// 		l("rand is high upper.%+v\n", s)
	// 	}
	// }
	l("\n========= サンドボックス単体実行完了\n\n")
}

type TestObject struct {
	Name string
	Text string
}

// func (to *TestObject) TestFunc() (a, b string) {
// 	return "abc", "defg"
// }

//手抜き用 : fmt.Printfを毎回打つのがだるいので。
func l(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf(format, a...)
}
