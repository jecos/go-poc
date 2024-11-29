package models

type Aggregation struct {
	Bucket string `json:"key"`
	Count  int64  `json:"count"`
}
