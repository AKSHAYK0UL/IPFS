package model

type Transaction struct {
	TxnId  string  `json:"txn_id" bson:"txn_id"`
	ToId   string  `json:"to_id" bson:"to_id"`
	FromId string  `json:"from_id" bson:"from_id"`
	Amount float64 `json:"amount" bson:"amount"`
	Nonce  int64   `json:"nonce" bson:"nonce"`
	Time   string  `json:"time" bson:"time"`
}
