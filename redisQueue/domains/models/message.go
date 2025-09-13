package models

type Msg struct {
	UserID    string `json:"user_id"`
	Message   string `json:"message"`
	RepeatNum int    `json:"repeat_num"`
}
