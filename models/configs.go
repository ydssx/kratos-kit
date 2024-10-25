package models

import (
	"context"

	"gorm.io/gorm"
)

// table sys_config 系统配置表
type SysConfig struct {
	BaseModelNoDelete
	Cate  string `json:"cate" gorm:"column:cate;not null;default:''"`
	Code  string `json:"code" gorm:"column:code;not null;default:''"`
	Value string `json:"value" gorm:"column:value;not null;default:''"`
	Type  int    `json:"type" gorm:"column:type;not null;default:1"`
	Desc  string `json:"desc" gorm:"column:desc;not null;default:''"`
}

func (m SysConfig) TableName() string {
	return "configs"
}

type CONF_CODE string

const (
	SWAP_TIME_RATIO                     CONF_CODE = "SWAP_TIME_RATIO"                 // 每处理一秒视频预估需要多少秒时间比例
	GOOGLE_INSTANCE_MAX                 CONF_CODE = "GOOGLE_INSTANCE_MAX"             // google实例最大数量
	GOOGLE_INSTANCE_MIN                 CONF_CODE = "GOOGLE_INSTANCE_MIN"             // google实例最少数量
	ENABLE_INSTANCE_ADJUST              CONF_CODE = "ENABLE_INSTANCE_ADJUST"          // 是否开启实例自动伸缩功能，1=开启，0=关闭
	FREE_SECONDS_RELEASE                CONF_CODE = "FREE_SECONDS_RELEASE"            // 空闲多少秒后释放
	FEE_PER_SECONDS                     CONF_CODE = "FEE_PER_SECONDS"                 // 每多少秒计费一次
	INSTANCE_REGION_LIST                CONF_CODE = "INSTANCE_REGION_LIST"            // 实例购买区域
	INSTANCE_CONF                       CONF_CODE = "INSTANCE_CONF"                   // 实例配置
	QUEUE_TIMEOUT_SECONDS               CONF_CODE = "QUEUE_TIMEOUT_SECONDS"           // 排队中任务超时时间，单位秒
	PROCESSING_TIMEOUT_SECONDS          CONF_CODE = "PROCESSING_TIMEOUT_SECONDS"      // 处理中任务超时时间，单位秒
	FREE_USER_MAX_TIME_RATIO            CONF_CODE = "FREE_USER_MAX_TIME_RATIO"        // 免费用户上传视频最大时长,单位秒
	FACE_DETECT_THRESHOLD               CONF_CODE = "FACE_DETECT_THRESHOLD"           // 人脸检测最低置信度，0-1
	GOOGLE_INSTANCE_PREFIX              CONF_CODE = "GOOGLE_INSTANCE_PREFIX"          // google实例前缀
	FACEFUSION_CONF                     CONF_CODE = "FACEFUSION_CONF"                 // facefusion配置
	DEMO_VIDEO_URL                      CONF_CODE = "DEMO_VIDEO_URL"                  // 演示视频地址
	SITE_CONFIG                         CONF_CODE = "SITE_CONFIG"                     // 默认站点配置
	DEF_FB                              CONF_CODE = "fb"                              // 默认FB配置
	DEF_GA                              CONF_CODE = "ga"                              // 默认GA配置
	DEF_ENV                             CONF_CODE = "sys_env"                         // 默认系统事件配置
	DEF_DEMO_VIDEO_URL                  CONF_CODE = "demo_video_url"                  // 默认换脸演示视频地址
	GOOGLE_INSTANCE_ADD_WHEN_QUEUED     CONF_CODE = "GOOGLE_INSTANCE_ADD_WHEN_QUEUED" // 排队时长达到多少秒，执行扩容
	GEN_COVER_SERVICE_URL               CONF_CODE = "GEN_COVER_SERVICE_URL"           // 生成动态封面服务地址
	AD_MATERIAL_INFO                    CONF_CODE = "AD_MATERIAL_INFO"                // 广告素材信息
	STATIC_SITE_CONFIG                  CONF_CODE = "STATIC_SITE_CONFIG"              // 静态站点配置
	CONF_CODE_HOME_PAGE_SEO_TITLE       CONF_CODE = "home_page_title"
	CONF_CODE_HOME_PAGE_SEO_KEYWORDS    CONF_CODE = "home_page_seo_keywords"
	CONF_CODE_HOME_PAGE_SEO_DESCRIPTION CONF_CODE = "home_page_seo_description"
	CONF_CODE_SITE_NAME                 CONF_CODE = "name"
)

type sysConfigModel DB

func NewSysConfigModel(tx *gorm.DB) *sysConfigModel {
	db := tx.Table("configs").Model(&SysConfig{})
	return &sysConfigModel{db: db}
}

func (m *sysConfigModel) Clone() *sysConfigModel {
	m.db = cloneDB(m.db)
	return m
}

func (m *sysConfigModel) SetIds(ids ...int64) *sysConfigModel {
	m.db = m.db.Where("id IN (?)", ids)
	return m
}

func (m *sysConfigModel) Order(expr string) *sysConfigModel {
	m.db = m.db.Order(expr)
	return m
}

func (m *sysConfigModel) Select(fields ...string) *sysConfigModel {
	m.db = m.db.Select(fields)
	return m
}

func (m *sysConfigModel) WithContext(ctx context.Context) *sysConfigModel {
	m.db = m.db.WithContext(ctx)
	return m
}

func (m *sysConfigModel) Create(sysConfig SysConfig) error {
	return m.db.Create(&sysConfig).Error
}

func (m *sysConfigModel) Updates(values interface{}) error {
	return m.db.Updates(values).Error
}

func (m *sysConfigModel) FirstOne() (data *SysConfig, err error) {
	err = m.db.Take(&data).Error
	return
}

func (m *sysConfigModel) LastOne() (data *SysConfig, err error) {
	err = m.db.Last(&data).Error
	return
}

func (m *sysConfigModel) DeleteByPrimKey(key interface{}) error {
	return m.db.Where(" IN (?)", key).Delete(&SysConfig{}).Error
}

func (m *sysConfigModel) List() (data []SysConfig) {
	m.db.Find(&data)
	return
}

func (m *sysConfigModel) PageList(limit, offset int) (data []SysConfig, total int64, err error) {
	query := m.db
	err = query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = query.Limit(limit).Offset(offset).Find(&data).Error
	return
}

func (m *sysConfigModel) Delete() error {
	return m.db.Delete(&SysConfig{}).Error
}
