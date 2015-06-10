package api

//endpoints汎用

//応答や受信で必要なデータ構造のうち、典型的なものをここにまとめておきます。

type Count struct {
	N int `json:"count"`
}

type LimitRequest struct {
	Limit int `json:"limit" endpoints="d=10"`
}
type PagingRequest struct {
	Limit  int `json:"limit" endpoints="d=10"`
	Offset int `json:"offset" endpoints="d=0"`
}

type SimpleResult struct {
	IsSuccess bool `json:"success"`
}

type SimpleResultWithMessage struct {
	IsSuccess         bool   `json:"success"`
	AdditionalMessage string `json:"additional"`
}
