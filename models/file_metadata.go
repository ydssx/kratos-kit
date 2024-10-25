package models

import (
	"context"
	"time"

	"github.com/Gre-Z/common/jtime"
	"gorm.io/gorm"
)

// table file_metadata 文件上传记录表
type FileMetadata struct {
	BaseModel
	UserId        int            `json:"user_id" gorm:"column:user_id;default:0"`
	Filename      string         `json:"filename" gorm:"column:filename;not null"`
	FileUrl       string         `json:"file_url" gorm:"column:file_url;not null"` // 文件存储路径
	UploadTime    jtime.JsonTime `json:"upload_time" gorm:"column:upload_time;default:CURRENT_TIMESTAMP"`
	FileSize      int            `json:"file_size" gorm:"column:file_size;default:NULL"`              // 文件大小(字节)
	FileType      FileType       `json:"file_type" gorm:"column:file_type;default:NULL"`              // 文件类型
	FileMd5       string         `json:"file_md5" gorm:"column:file_md5;default:NULL"`                // 文件MD5值
	VideoDuration float64        `json:"video_duration" gorm:"column:video_duration;default:NULL"`    // 视频时长(秒)
	CoverUrl      string         `json:"cover_url" gorm:"column:cover_url;default:NULL"`              // 视频封面图路径
	Width         int            `json:"width" gorm:"column:width;not null;default:0"`                // 视频宽度，单位像素
	Height        int            `json:"height" gorm:"column:height;not null;default:0"`              // 视频高度，单位像素
	Fps           float64        `json:"fps" gorm:"column:fps;not null;default:0"`                    // 视频帧率
	Encoding      string         `json:"encoding" gorm:"column:encoding;type:VARCHAR(50);default:''"` // 编码
}

type fileMetadataModel DB

type FileType string

const (
	FileTypeVideo FileType = "video"
	FileTypeImage FileType = "image"
	FileTypeAudio FileType = "audio"
)

func NewFileMetadataModel(tx ...*gorm.DB) *fileMetadataModel {
	db := getDB(tx...).Table("file_metadata").Model(&FileMetadata{})
	return &fileMetadataModel{db: db}
}

func (m *fileMetadataModel) Clone() *fileMetadataModel {
	m.db = cloneDB(m.db)
	return m
}

func (m *fileMetadataModel) SetIds(ids ...int64) *fileMetadataModel {
	m.db = m.db.Where("id IN (?)", ids)
	return m
}

func (m *fileMetadataModel) SetMd5(md5 string) *fileMetadataModel {
	m.db = m.db.Where("file_md5 = ?", md5)
	return m
}

// SetUserId 设置用户ID
func (m *fileMetadataModel) SetUserId(userId ...int64) *fileMetadataModel {
	m.db = m.db.Where("user_id IN (?)", userId)
	return m
}

func (m *fileMetadataModel) Order(expr string) *fileMetadataModel {
	m.db = m.db.Order(expr)
	return m
}

func (m *fileMetadataModel) Select(fields ...string) *fileMetadataModel {
	m.db = m.db.Select(fields)
	return m
}

func (m *fileMetadataModel) WithContext(ctx context.Context) *fileMetadataModel {
	m.db = m.db.WithContext(ctx)
	return m
}

func (m *fileMetadataModel) Create(fileMetadata *FileMetadata) (*FileMetadata, error) {
	err := m.db.Create(&fileMetadata).Error
	return fileMetadata, err
}

func (m *fileMetadataModel) Updates(values interface{}) error {
	return m.db.Updates(values).Error
}

func (m *fileMetadataModel) FirstOne() (data FileMetadata, err error) {
	err = m.db.Take(&data).Error
	return
}

func (m *fileMetadataModel) LastOne() (data *FileMetadata, err error) {
	err = m.db.Last(&data).Error
	return
}

func (m *fileMetadataModel) DeleteByPrimKey(key interface{}) error {
	return m.db.Where("id IN (?)", key).Delete(&FileMetadata{}).Error
}

func (m *fileMetadataModel) List() (data []FileMetadata) {
	m.db.Find(&data)
	return
}

func (m *fileMetadataModel) PageList(limit, offset int) (data []FileMetadata, total int64, err error) {
	query := m.db
	err = query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = query.Limit(limit).Offset(offset).Find(&data).Error
	return
}

func (m *fileMetadataModel) Delete() error {
	return m.db.Delete(&FileMetadata{}).Error
}

// created_at小于
func (m *fileMetadataModel) CreatedAtLT(t time.Time) *fileMetadataModel {
	m.db = m.db.Where("created_at < ?", t)
	return m
}
