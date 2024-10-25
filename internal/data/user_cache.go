package data

import (
	"context"

	"github.com/ydssx/kratos-kit/internal/biz"
	"github.com/ydssx/kratos-kit/models"
	"github.com/ydssx/kratos-kit/pkg/cache"
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
	// key := fmt.Sprintf("user.list:%v", util.CalculateChecksum(cond))
	// err := u.Get(key, &data)
	// if err == nil {
	// 	return
	// }

	data = u.userRepo.ListUser(ctx, cond)

	// err = u.Set(key, data, time.Hour)
	// if err != nil {
	// 	logger.Errorf(ctx, "缓存用户列表失败：%v", err)
	// }

	return
}

// GetUserByID 从缓存中获取用户数据。
// 如果缓存中不存在,则从数据库中获取用户数据,并设置缓存。
// 缓存时间为 1 小时。
func (u *UserRepoCacheDecorator) GetUserByID(ctx context.Context, id uint) (data *models.User, err error) {
	//key := fmt.Sprintf("user:%v", id)
	//err = u.Get(key, &data)
	//if err == nil {
	//	return
	//}

	data, err = u.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	//err = u.Set(key, data, time.Hour)
	//if err != nil {
	//	return nil, err
	//}

	return
}

// UpdateUser 更新用户数据,并删除与该用户相关的缓存
// 如果更新用户失败,返回错误
// 如果删除缓存失败,记录错误日志但不影响函数返回
// 成功更新用户后返回 nil
func (u *UserRepoCacheDecorator) UpdateUser(ctx context.Context, userId int, user interface{}) error {
	if err := u.userRepo.UpdateUser(ctx, userId, user); err != nil {
		return err
	}

	// key := fmt.Sprintf("user:%v", updatedUser.ID)
	// err := u.Delete(key)
	// if err != nil {
	// 	logger.Errorf(ctx, "删除缓存失败：%v", err)
	// }

	return nil
}
