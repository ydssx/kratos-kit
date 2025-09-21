package health

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// HealthChecker 健康检查接口
type HealthChecker interface {
	Check(ctx context.Context) error
	Name() string
}

// DatabaseHealthChecker 数据库健康检查
type DatabaseHealthChecker struct {
	db *gorm.DB
}

func NewDatabaseHealthChecker(db *gorm.DB) *DatabaseHealthChecker {
	return &DatabaseHealthChecker{db: db}
}

func (h *DatabaseHealthChecker) Check(ctx context.Context) error {
	sqlDB, err := h.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}
	
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	
	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}
	
	return nil
}

func (h *DatabaseHealthChecker) Name() string {
	return "database"
}

// RedisHealthChecker Redis健康检查
type RedisHealthChecker struct {
	client *redis.Client
}

func NewRedisHealthChecker(client *redis.Client) *RedisHealthChecker {
	return &RedisHealthChecker{client: client}
}

func (h *RedisHealthChecker) Check(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	
	if err := h.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}
	
	return nil
}

func (h *RedisHealthChecker) Name() string {
	return "redis"
}

// HealthService 健康检查服务
type HealthService struct {
	checkers []HealthChecker
	logger   log.Logger
}

func NewHealthService(logger log.Logger, checkers ...HealthChecker) *HealthService {
	return &HealthService{
		checkers: checkers,
		logger:   logger,
	}
}

// CheckAll 检查所有服务
func (h *HealthService) CheckAll(ctx context.Context) map[string]error {
	results := make(map[string]error)
	
	for _, checker := range h.checkers {
		if err := checker.Check(ctx); err != nil {
			results[checker.Name()] = err
			h.logger.Errorf(ctx, "health check failed for %s: %v", checker.Name(), err)
		} else {
			results[checker.Name()] = nil
		}
	}
	
	return results
}

// IsHealthy 检查整体健康状态
func (h *HealthService) IsHealthy(ctx context.Context) bool {
	results := h.CheckAll(ctx)
	for _, err := range results {
		if err != nil {
			return false
		}
	}
	return true
}

// HTTPHandler 返回HTTP健康检查处理器
func (h *HealthService) HTTPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		results := h.CheckAll(ctx)
		
		w.Header().Set("Content-Type", "application/json")
		
		healthy := true
		for _, err := range results {
			if err != nil {
				healthy = false
				break
			}
		}
		
		if healthy {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"status":"healthy","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, `{"status":"unhealthy","timestamp":"%s","errors":{`, time.Now().Format(time.RFC3339))
			
			first := true
			for name, err := range results {
				if err != nil {
					if !first {
						fmt.Fprint(w, ",")
					}
					fmt.Fprintf(w, `"%s":"%s"`, name, err.Error())
					first = false
				}
			}
			fmt.Fprint(w, "}}")
		}
	}
}
