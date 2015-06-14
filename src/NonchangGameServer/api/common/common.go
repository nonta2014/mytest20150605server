package common

//endpoints汎用

//応答や受信で必要なJSONデータ構造のうち、典型的なものをここにまとめておきます。

type CountResult struct {
	N int `json:"count"`
}

type LimitRequest struct {
	Limit int `json:"limit" endpoints="d=10"`
}

//これは使わない。offsetアクセスはDatastoreにとって最悪
// type PagingRequest struct {
// 	Limit  int `json:"limit" endpoints="d=10"`
// 	Offset int `json:"offset" endpoints="d=0"`
// }

type LimitAndCursorRequest struct {
	Limit  int    `json:"limit"  endpoints="d=10"`
	Cursor string `json:"cursor" endpoints="d=none"`
}

type SimpleResult struct {
	IsSuccess bool `json:"success"`
}

type SimpleResultWithMessage struct {
	IsSuccess         bool   `json:"success"`
	AdditionalMessage string `json:"additional"`
}

type UUIDRequest struct {
	UUID string `json:"uuid"`
}
