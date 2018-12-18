package serializer

// WishResponse 许愿序列化器
type WishResponse struct {
	Success bool   `json:"success"`
	Hard    int64  `json:"hard"`
	Type    string `json:"type"`
	Amount  int64  `json:"amount"`
	Stock   string `json:"stock"`
	Msg     string `json:"msg"`
}
