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
