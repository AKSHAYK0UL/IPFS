package model

type Transaction struct {
	TxnId  string  `json:"txn_id"`
	ToId   string  `json:"to_id"`
	FromId string  `json:"from_id"`
	Amount float64 `json:"amount"`
	Nonce  int64   `json:"nonce"`
	Time   string  `json:"time"`
}

type IPFSTransaction struct {
	CID       string  `json:"cid"`
	Index     int     `json:"index"`
	PoolIndex int     `json:"pool_index"`
	Hash      string  `json:"hash"`
	PrevHash  string  `json:"prev_hash"`
	TxnId     string  `json:"txn_id"`
	ToId      string  `json:"to_id"`
	FromId    string  `json:"from_id"`
	Amount    float64 `json:"amount" `
	Nonce     int64   `json:"nonce" `
	Time      string  `json:"time"`
}
