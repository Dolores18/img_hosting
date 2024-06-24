package models

// User 用户模型
type UserInfo struct {
	Name string `json:"name" binding:"required,sign"`
	Age  int    `json:"age" binding:"required,gte=18"`
	Psd  string `json:"psd" binding:"required,password"`
}
