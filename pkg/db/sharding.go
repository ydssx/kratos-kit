package db

import (
	"fmt"
	"hash/crc32"
)

type ShardingStrategy int

const (
	// 分表策略
	HashSharding  ShardingStrategy = iota // 哈希分片
	RangeSharding                         // 范围分片
	TimeSharding                          // 时间分片
)

type TableSharding struct {
	TableName    string           // 逻辑表名
	Strategy     ShardingStrategy // 分表策略
	TableNum     int              // 分表数量
	RangeSize    int64            // 范围分片大小
	TimeInterval string           // 时间间隔 (day, month, year)
}

// 获取实际的表名
func (ts *TableSharding) GetTableName(shardKey any) string {
	switch ts.Strategy {
	case HashSharding:
		return ts.getHashTableName(shardKey)
	case RangeSharding:
		return ts.getRangeTableName(shardKey)
	case TimeSharding:
		return ts.getTimeTableName(shardKey)
	default:
		return ts.TableName
	}
}

// 哈希分片
func (ts *TableSharding) getHashTableName(shardKey any) string {
	hash := crc32.ChecksumIEEE([]byte(fmt.Sprintf("%v", shardKey)))
	index := hash % uint32(ts.TableNum)
	return fmt.Sprintf("%s_%d", ts.TableName, index)
}

// 范围分片
func (ts *TableSharding) getRangeTableName(shardKey any) string {
	id, ok := shardKey.(int64)
	if !ok {
		return ts.TableName
	}
	index := id / ts.RangeSize
	return fmt.Sprintf("%s_%d", ts.TableName, index)
}

// 时间分片
func (ts *TableSharding) getTimeTableName(shardKey any) string {
	timeStr, ok := shardKey.(string)
	if !ok {
		return ts.TableName
	}
	return fmt.Sprintf("%s_%s", ts.TableName, timeStr)
}

// 批量查询多个分表
func (sq *ShardingQuery) MultiFind(dest any, tableIndices []int, conds ...any) error {
	var results []any
	for _, idx := range tableIndices {
		var result any
		tableName := fmt.Sprintf("%s_%d", sq.Sharding.TableName, idx)
		if err := sq.DB.Table(tableName).Find(&result, conds...).Error; err != nil {
			return err
		}
		results = append(results, result)
	}
	// 合并结果
	return nil
}

// 自动创建分表
func (sq *ShardingQuery) CreateTable(tableIndex int) error {
	tableName := fmt.Sprintf("%s_%d", sq.Sharding.TableName, tableIndex)
	return sq.DB.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s LIKE %s`,
		tableName, sq.Sharding.TableName)).Error
}
