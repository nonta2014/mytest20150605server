package main

import (
	"fmt"
)

func main() {

	a := &TestObject{Name: "123", Text: "456"}

	fmt.Printf(a.Name + "\n")
	_, b := a.TestFunc()
	fmt.Printf(b + "\n")

	fmt.Printf("サンドボックス単体実行完了\n")
}

type TestObject struct {
	Name string
	Text string
}

func (to *TestObject) TestFunc() (a, b string) {
	return "abc", "defg"
}
