package cache

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-keg/keg/contrib/config"
)

var redisCache Cache

func init() {
	config.LoadEnv()
	redisCache = NewRedis(
		os.Getenv("REDIS_ADDR"),
		SetPassword(os.Getenv("REDIS_PASSWORD")),
	)
}

func TestHas(t *testing.T) {
	has, err := redisCache.Has(context.Background(), "test_string")
	if err != nil {
		return
	}
	fmt.Println("test_string", has, err)
}

func TestSet(t *testing.T) {
	err := redisCache.Set(context.Background(), "test_string1", time.Now().String(), time.Minute*5)
	fmt.Println("Set test_string", err)
}

func TestGet(t *testing.T) {
	result, err := redisCache.Get(context.Background(), "test_string1")
	fmt.Println("Get test_string", string(result), err)
}

func TestForget(t *testing.T) {
	err := redisCache.Forget(context.Background(), "test_string")
	fmt.Println("Forget test_set", err)
}

func TestRemember(t *testing.T) {
	result, err := redisCache.Remember(context.Background(), "test_string", 20*time.Second, func() ([]byte, error) {
		return []byte("123456"), nil
	})
	fmt.Println("result", string(result), "err", err)
}
