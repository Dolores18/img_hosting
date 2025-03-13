package models

// Response API通用响应结构
type Response struct {
	Status  string      `json:"status" example:"success"`
	Message string      `json:"message,omitempty" example:"操作成功"`
	Data    interface{} `json:"data,omitempty"`
}

// TokenVerifyResponse Token验证响应
type TokenVerifyResponse struct {
	Status string `json:"status" example:"valid"`
	UserID uint   `json:"user_id" example:"1"`
	Cached bool   `json:"cached" example:"false"`
}

// ImageUploadResponse 图片上传响应
type ImageUploadResponse struct {
	ImageID  uint   `json:"image_id" example:"1"`
	ImageURL string `json:"image_url" example:"http://example.com/images/1.jpg"`
}

// BatchUploadResponse 批量上传响应
type BatchUploadResponse struct {
	Message      string                `json:"message" example:"批量上传完成"`
	Total        int                   `json:"total" example:"5"`
	SuccessCount int                   `json:"success_count" example:"3"`
	Results      []ImageUploadResponse `json:"results"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Total    int64      `json:"total" example:"100"`
	Page     int        `json:"page" example:"1"`
	PageSize int        `json:"page_size" example:"10"`
	Users    []UserInfo `json:"users"`
}

// ImageListResponse 图片列表响应
type ImageListResponse struct {
	Total    int64   `json:"total" example:"100"`
	Page     int     `json:"page" example:"1"`
	PageSize int     `json:"page_size" example:"10"`
	Images   []Image `json:"images"`
}

// UserUpdateRequest 用户信息更新请求
type UserUpdateRequest struct {
	Name     string `json:"name" example:"张三"`
	Email    string `json:"email" example:"zhangsan@example.com"`
	Password string `json:"password,omitempty"`
	Age      int    `json:"age" example:"25"`
}

// StatusUpdateRequest 状态更新请求
type StatusUpdateRequest struct {
	Status string `json:"status" example:"active" enums:"active,inactive,banned"`
}

// RoleAssignRequest 角色分配请求
type RoleAssignRequest struct {
	Roles []string `json:"roles" example:"admin,editor"`
}

// FileListResponse 文件列表响应
type FileListResponse struct {
	Total    int64         `json:"total" example:"100"`
	Page     int           `json:"page" example:"1"`
	PageSize int           `json:"page_size" example:"10"`
	Files    []PrivateFile `json:"files"`
}

// FileUpdateRequest 文件更新请求
type FileUpdateRequest struct {
	FileName    string `json:"file_name" example:"document.pdf"`
	IsEncrypted bool   `json:"is_encrypted" example:"true"`
	Password    string `json:"password,omitempty" example:"your-password"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"your-password"`
	Psd      string `json:"psd,omitempty" example:"your-password"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token    string `json:"token" example:"eyJhbGciOiJIUzI1NiIs..."`
	UserID   uint   `json:"user_id" example:"1"`
	UserName string `json:"user_name" example:"张三"`
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	UserID uint   `json:"user_id" example:"1"`
	Status string `json:"status" example:"success"`
}

// PermissionResponse 权限响应
type PermissionResponse struct {
	Name        string `json:"name" example:"create_post"`
	Description string `json:"description" example:"允许创建新文章"`
}
