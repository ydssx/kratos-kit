package models

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// table users 用户表
type User struct {
	BaseModel
	Username               string `json:"username" gorm:"column:username;not null"`
	Email                  string `json:"email" gorm:"column:email;not null"`                 // 登录邮箱
	AvatarPath             string `json:"avatar_path" gorm:"column:avatar_path;default:NULL"` // 头像路径
	PasswordHash           string `json:"password_hash" gorm:"column:password_hash;not null"`
	GoogleId               string `json:"google_id" gorm:"column:google_id;default:NULL"`                           // google登录id
	UUID                   string `json:"uuid" gorm:"column:uuid;not null"`
	Type                   int    `json:"type" gorm:"column:type;not null;default:0"`
	IPAddress              string `json:"ip_address" gorm:"column:ip_address"`
	CountryCode            string `json:"country_code" gorm:"column:country_code"`
	CityCode               string `json:"city_code" gorm:"column:city_code"`
	CountryName            string `json:"country_name" gorm:"column:country_name"`
	ZipCode                string `json:"zip_code" gorm:"column:zip_code"`
	Platform               int    `json:"platform" gorm:"column:platform;default:1"`                                   // 注册来源平台,1:h5,2:pc
	FirstName              string `json:"first_name" gorm:"column:first_name"`                                         // 姓
	LastName               string `json:"last_name" gorm:"column:last_name"`                                           // 名
	BrowserFingerprint     string `json:"browser_fingerprint" gorm:"column:browser_fingerprint"`                       // 浏览器指纹
	Note                   string `json:"note" gorm:"column:note"`                                                     // 备注
	DeviceModel            string `json:"device_model" gorm:"column:device_model"`                                     // 设备型号
}

type UserType int

const (
	UserTypeNormal UserType = iota
	UserTypeAiNode
	UserTypeAdmin
	UserTypeLogout
)

type userModel DB

func NewUserModel(tx ...*gorm.DB) *userModel {
	db := getDB(tx...).Table("users").Model(&User{})
	return &userModel{db: db}
}

func (m *userModel) Clone() *userModel {
	return &userModel{db: cloneDB(m.db)}
}

func (m *userModel) SetIds(ids ...int) *userModel {
	m.db = m.db.Where("id IN (?)", ids)
	return m
}

func (m *userModel) SetUUIds(uuids ...string) *userModel {
	m.db = m.db.Where("uuid IN (?)", uuids)
	return m
}

func (m *userModel) SetUsername(username string) *userModel {
	m.db = m.db.Where("username = ?", username)
	return m
}

func (m *userModel) SetEmail(email string) *userModel {
	m.db = m.db.Where("email = ?", email)
	return m
}

func (m *userModel) Order(expr string) *userModel {
	m.db = m.db.Order(expr)
	return m
}

func (m *userModel) Select(fields ...string) *userModel {
	m.db = m.db.Select(fields)
	return m
}

func (m *userModel) WithContext(ctx context.Context) *userModel {
	m.db = m.db.WithContext(ctx)
	return m
}

func (u *userModel) Create(user User) (User, error) {
	err := u.db.Create(&user).Error
	return user, err
}

func (m *userModel) Updates(values interface{}) error {
	return m.db.Updates(values).Error
}

func (m *userModel) FirstOne() (data *User, err error) {
	err = m.db.Take(&data).Error
	return
}

func (m *userModel) LastOne() (data *User, err error) {
	err = m.db.Last(&data).Error
	return
}

func (m *userModel) DeleteByPrimKey(key interface{}) error {
	return m.db.Where("id IN (?)", key).Delete(&User{}).Error
}

func (m *userModel) List() (data []User, err error) {
	err = m.db.Find(&data).Error
	return
}

func (m *userModel) PageList(limit, offset int) (data []User, total int64, err error) {
	query := m.db
	err = query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = query.Limit(limit).Offset(offset).Find(&data).Error
	return
}

func (m *userModel) Delete() error {
	return m.db.Delete(&User{}).Error
}

func (u *userModel) ListAll() (users []User, err error) {
	err = u.db.Find(&users).Error
	return
}

// Xlock
func (u *userModel) XLock() *userModel {
	u.db = u.db.Clauses(clause.Locking{Strength: "UPDATE"})
	return u
}

// SetUserType 设置用户类型
func (u *userModel) SetUserType(userType UserType) *userModel {
	u.db = u.db.Where("type = ?", userType)
	return u
}

// PluckIds 批量获取id
func (u *userModel) PluckIds() (ids []int64, err error) {
	err = u.db.Pluck("id", &ids).Error
	return
}

// 积分数大于指定值
func (u *userModel) PointsGt(points int) *userModel {
	u.db = u.db.Where("subscription_points > ?", points)
	return u
}

// SetBrowserFingerprint 设置浏览器指纹
func (u *userModel) SetBrowserFingerprint(fingerprint string) *userModel {
	u.db = u.db.Where("browser_fingerprint = ?", fingerprint)
	return u
}

func (u *userModel) Where(query string, args ...interface{}) *userModel {
	u.db = u.db.Where(query, args...)
	return u
}

func (u *userModel) Count() (int64, error) {
	var count int64
	err := u.db.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SetSubscribeStatus 设置订阅状态
func (u *userModel) SetSubscribeStatus(status string) *userModel {
	u.db = u.db.Where("subscription_status = ?", status)
	return u
}

// SetGoogleID 设置Google ID
func (u *userModel) SetGoogleID(googleID string) *userModel {
	u.db = u.db.Where("google_id = ?", googleID)
	return u
}
