package db

import (
	"gorm.io/gorm"
)

type ShardingQuery struct {
	DB       *gorm.DB
	Sharding *TableSharding
	ShardKey any
}

func NewShardingQuery(db *gorm.DB, sharding *TableSharding) *ShardingQuery {
	return &ShardingQuery{
		DB:       db,
		Sharding: sharding,
	}
}

func (sq *ShardingQuery) SetShardKey(key any) *ShardingQuery {
	sq.ShardKey = key
	return sq
}

// 查询方法
func (sq *ShardingQuery) Find(dest any, conds ...any) error {
	tableName := sq.Sharding.GetTableName(sq.ShardKey)
	return sq.DB.Table(tableName).Find(dest, conds...).Error
}

// 插入方法
func (sq *ShardingQuery) Create(value any) error {
	tableName := sq.Sharding.GetTableName(sq.ShardKey)
	return sq.DB.Table(tableName).Create(value).Error
}

// 更新方法
func (sq *ShardingQuery) Update(column string, value any) error {
	tableName := sq.Sharding.GetTableName(sq.ShardKey)
	return sq.DB.Table(tableName).Update(column, value).Error
}

// 删除方法
func (sq *ShardingQuery) Delete(conds ...any) error {
	tableName := sq.Sharding.GetTableName(sq.ShardKey)
	return sq.DB.Table(tableName).Delete(conds).Error
}
