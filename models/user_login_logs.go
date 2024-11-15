package models

import (
	"context"
	"time"

	"github.com/Gre-Z/common/jtime"
	"gorm.io/gorm"
)

// table user_login_logs 用户上线日志表
type UserLoginLog struct {
	ID          uint           `json:"id" gorm:"primary_key"`
	UserId      int            `json:"user_id" gorm:"column:user_id;not null"`       // 用户ID
	LoginDate   jtime.JsonTime `json:"login_date" gorm:"column:login_date;not null"` // 登录日期
	IpAddress   string         `json:"ip_address" gorm:"column:ip_address"`          // IP地址
	CountryName string         `json:"country_name" gorm:"column:country_name"`      // 国家名
	CountryCode string         `json:"country_code" gorm:"column:country_code"`      // 国家编码
	CityName    string         `json:"city_name" gorm:"column:city_name"`            // 城市名
	Extra       string         `json:"extra" gorm:"column:extra"`                    // 附加信息
	CreatedAt   jtime.JsonTime `json:"created_at" gorm:"column:created_at"`          // 创建时间
}

type userLoginLogModel DB

func NewUserLoginLogModel(tx ...*gorm.DB) *userLoginLogModel {
	db := getDB(tx...).Table("user_login_logs").Model(&UserLoginLog{})
	return &userLoginLogModel{db: db}
}

func (m *userLoginLogModel) Clone() *userLoginLogModel {
	m.db = cloneDB(m.db)
	return m
}

func (m *userLoginLogModel) SetIds(ids ...int64) *userLoginLogModel {
	m.db = m.db.Where("user_login_logs.id IN (?)", ids)
	return m
}

func (m *userLoginLogModel) Order(expr string) *userLoginLogModel {
	m.db = m.db.Order(expr)
	return m
}

func (m *userLoginLogModel) Select(fields ...string) *userLoginLogModel {
	m.db = m.db.Select(fields)
	return m
}

func (m *userLoginLogModel) WithContext(ctx context.Context) *userLoginLogModel {
	m.db = m.db.WithContext(ctx)
	return m
}

func (m *userLoginLogModel) Create(userLoginLog UserLoginLog) error {
	return m.db.Create(&userLoginLog).Error
}

func (m *userLoginLogModel) Updates(values interface{}) error {
	return m.db.Updates(values).Error
}

func (m *userLoginLogModel) FirstOne() (data *UserLoginLog, err error) {
	err = m.db.Take(&data).Error
	return
}

func (m *userLoginLogModel) LastOne() (data *UserLoginLog, err error) {
	err = m.db.Last(&data).Error
	return
}

func (m *userLoginLogModel) DeleteByPrimKey(key interface{}) error {
	return m.db.Where(" IN (?)", key).Delete(&UserLoginLog{}).Error
}

func (m *userLoginLogModel) List() (data []UserLoginLog) {
	m.db.Find(&data)
	return
}

func (m *userLoginLogModel) PageList(limit, offset int) (data []UserLoginLog, total int64, err error) {
	query := m.db
	err = query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = query.Limit(limit).Offset(offset).Find(&data).Error
	return
}

func (m *userLoginLogModel) Delete() error {
	return m.db.Delete(&UserLoginLog{}).Error
}

// SetUserId 设置用户ID
func (m *userLoginLogModel) SetUserId(userId int) *userLoginLogModel {
	m.db = m.db.Where("user_login_logs.user_id = ?", userId)
	return m
}

// SetLoginDate 设置登录日期
func (m *userLoginLogModel) SetLoginDate(loginDate time.Time) *userLoginLogModel {
	m.db = m.db.Where("user_login_logs.login_date = ?", loginDate.Format("2006-01-02"))
	return m
}

func (m *userLoginLogModel) Count() (total int64, err error) {
	err = m.db.Select("DISTINCT user_login_logs.user_id").Count(&total).Error
	return
}

// 左连接用户表查询
func (m *userLoginLogModel) LeftJoinUser() *userLoginLogModel {
	m.db = m.db.Joins("LEFT JOIN users ON users.id = user_login_logs.user_id")
	return m
}

// SetDomain 设置域名
func (m *userLoginLogModel) SetDomain(domain string) *userLoginLogModel {
	m.db = m.db.Where("users.source_domain = ?", domain)
	return m
}

// LteLoginDate 小于等于登录日期
func (m *userLoginLogModel) LteLoginDate(loginDate time.Time) *userLoginLogModel {
	m.db = m.db.Where("user_login_logs.login_date <= ?", loginDate.Format("2006-01-02"))
	return m
}

// GteLoginDate 大于等于登录日期
func (m *userLoginLogModel) GteLoginDate(loginDate time.Time) *userLoginLogModel {
	m.db = m.db.Where("user_login_logs.login_date >= ?", loginDate.Format("2006-01-02"))
	return m
}
