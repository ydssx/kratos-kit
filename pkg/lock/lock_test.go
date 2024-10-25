package lock

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/ydssx/kratos-kit/pkg/client/redis"

	goredis "github.com/redis/go-redis/v9"
)

func TestRedisLocker_Lock(t *testing.T) {
	cli, _ := redis.NewRedis(&goredis.Options{Addr: "localhost:6379"})
	r := NewLocker(cli)
	// err := r.Lock(context.Background(), "test", WithTries(0), WithDelay(time.Second), WithTTL(time.Second*3))
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// defer r.Unlock(context.Background(), "test")
	// log.Println("lock success")
	// err = r.Lock(context.Background(), "test", WithTries(1))
	// log.Println(err)
	a := 1
	f := func(r Locker, t int) {
		err := r.Lock(context.Background(), "test", WithTTL(time.Second*time.Duration(7)), WithDelay(time.Second))
		if err != nil {
			log.Println(err)
			return
		}
		defer func() {
			err = r.Unlock(context.Background(), "test")
			log.Println(err)
		}()
		log.Println("do something start,t=", t)
		a = a + t
		log.Println("do something success, a=", a)
		time.Sleep(time.Second * time.Duration(t))
	}
	go f(r, 5)
	go f(r, 3)
	time.Sleep(time.Second * 7)
	log.Println("a=", a)
}
