package biz

import goauth2 "google.golang.org/api/oauth2/v2"

// 事件配置
type EnvConf struct {
	Enable bool   `json:"enable"`
	ID     string `json:"id"`
	Token  string `json:"token"`
}

type UploadResult struct {
	FileUrl      string `json:"file_url"`
	ThumbnailURL string `json:"thumbnail_url"`
	FileId       int    `json:"file_id"`
	FileName     string `json:"file_name"`
}

type GoogleClaims struct {
	Sub string `json:"sub"`
	*goauth2.Userinfo
}