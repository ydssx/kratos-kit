package data

import (
	"context"
	"fmt"
	"time"

	"github.com/ydssx/kratos-kit/internal/biz"
	"github.com/ydssx/kratos-kit/models"
	"github.com/ydssx/kratos-kit/pkg/cache"
	"github.com/ydssx/kratos-kit/pkg/logger"
	"github.com/ydssx/kratos-kit/pkg/util"
)

type UserRepoCacheDecorator struct {
	*userRepo
	cache.Cache
}

// NewUserRepoCacheDecorator creates a new UserRepoCacheDecorator, which wraps
// a userRepo with caching capabilities using the provided cache.Cache
// implementation.
func NewUserRepoCacheDecorator(repo *userRepo, cache cache.Cache) biz.UserRepo {
	return &UserRepoCacheDecorator{repo, cache}
}

// ListUser retrieves users from cache if available, otherwise from database.
// It caches the retrieved users for 1 hour.
func (u *UserRepoCacheDecorator) ListUser(ctx context.Context, cond *biz.ListUserCond) (data []models.User) {
	key := fmt.Sprintf("user.list:%v", util.CalculateChecksum(cond))
	data, err := cache.WithCache(u.Cache, ctx, key, time.Hour, func() ([]models.User, error) {
		return u.ListUser(ctx, cond), nil
	})
	if err != nil {
		return nil
	}

	return
}

// GetUserByID 从缓存中获取用户数据。
// 如果缓存中不存在,则从数据库中获取用户数据,并设置缓存。
// 缓存时间为 1 小时。
func (u *UserRepoCacheDecorator) GetUserByID(ctx context.Context, id uint) (data *models.User, err error) {
	key := fmt.Sprintf("user:%v", id)

	return cache.WithCache(u.Cache, ctx, key, time.Hour, func() (*models.User, error) {
		return u.GetUserByID(ctx, id)
	})
}

// UpdateUser 更新用户数据,并删除与该用户相关的缓存
func (u *UserRepoCacheDecorator) UpdateUser(ctx context.Context, userId int, user interface{}) error {
	if err := u.userRepo.UpdateUser(ctx, userId, user); err != nil {
		return err
	}

	key := fmt.Sprintf("user:%v", userId)
	if err := u.Cache.Delete(ctx, key); err != nil {
		logger.Errorf(ctx, "删除缓存失败：%v", err)
	}

	return nil
}
