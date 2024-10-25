package models

import (
	"context"

	"gorm.io/gorm"
)

// table admin_users 系统管理员表
type AdminUser struct {
	BaseModelNoDelete
	Account  string `json:"account" gorm:"column:account;not null;default:''"`   // 账号
	Password string `json:"password" gorm:"column:password;not null;default:''"` // 密码
	Name     string `json:"name" gorm:"column:name;not null;default:''"`         // 姓名
	IsAdmin  int    `json:"is_admin" gorm:"column:is_admin;not null;default:0"`  // 0=普通账号，1=管理员
	Status   int    `json:"status" gorm:"column:status;not null;default:1"`      // 状态(0禁用,1启用)
}

func (m AdminUser) TableName() string {
	return "admin_users"
}

type adminUserModel DB

func NewAdminUserModel(tx *gorm.DB) *adminUserModel {
	db := tx.Model(&AdminUser{})
	return &adminUserModel{db: db}
}

func (m *adminUserModel) Clone() *adminUserModel {
	m.db = cloneDB(m.db)
	return m
}

func (m *adminUserModel) SetIds(ids ...int64) *adminUserModel {
	m.db = m.db.Where("id IN (?)", ids)
	return m
}

func (m *adminUserModel) Order(expr string) *adminUserModel {
	m.db = m.db.Order(expr)
	return m
}

func (m *adminUserModel) Select(fields ...string) *adminUserModel {
	m.db = m.db.Select(fields)
	return m
}

func (m *adminUserModel) WithContext(ctx context.Context) *adminUserModel {
	m.db = m.db.WithContext(ctx)
	return m
}

func (m *adminUserModel) Create(adminUser AdminUser) error {
	return m.db.Create(&adminUser).Error
}

func (m *adminUserModel) Updates(values interface{}) error {
	return m.db.Updates(values).Error
}

func (m *adminUserModel) FirstOne() (data *AdminUser, err error) {
	err = m.db.Take(&data).Error
	return
}

func (m *adminUserModel) LastOne() (data *AdminUser, err error) {
	err = m.db.Last(&data).Error
	return
}

func (m *adminUserModel) DeleteByPrimKey(key interface{}) error {
	return m.db.Where(" IN (?)", key).Delete(&AdminUser{}).Error
}

func (m *adminUserModel) List() (data []AdminUser, err error) {
	err = m.db.Find(&data).Error
	return
}

func (m *adminUserModel) PageList(pageNum, pageSize int) (data []AdminUser, total int64, err error) {
	query := m.db
	err = query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = query.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&data).Error
	return
}

func (m *adminUserModel) Delete() error {
	return m.db.Delete(&AdminUser{}).Error
}
