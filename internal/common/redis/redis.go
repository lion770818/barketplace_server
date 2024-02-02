package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
	redis "github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}
type RedisParameter struct {
	Network      string
	Address      string
	Password     string
	DB           int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PoolSize     int
}

func NewRedis(param *RedisParameter) (*Redis, error) {

	if param == nil {
		return nil, fmt.Errorf("param is nil")
	}

	client := redis.NewClient(&redis.Options{
		Network:      param.Network,
		Addr:         param.Address,
		Password:     param.Password,
		DB:           param.DB,
		DialTimeout:  param.DialTimeout,
		ReadTimeout:  param.ReadTimeout,
		WriteTimeout: param.WriteTimeout,
		PoolSize:     param.PoolSize,
	})

	_, err := client.Ping(context.TODO()).Result()

	return &Redis{client: client}, err

}

func (rds *Redis) Set(key string, value interface{}) error {
	err := rds.client.Set(context.TODO(), key, value, 0).Err()

	if err != nil {
		return err
	}

	return nil
}

func (rds *Redis) GetSting(key string, defaultValue string) (string, error) {
	value, err := rds.client.Get(context.TODO(), key).Result()

	if err != nil {
		return defaultValue, err
	}

	return value, nil
}

func (rds *Redis) GetInt(key string, defaultValue int) (int, error) {
	value, err := rds.client.Get(context.TODO(), key).Result()

	if err != nil {
		return defaultValue, err
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue, err
	}

	return intValue, nil
}

func (rds *Redis) GetInt64(key string, defaultValue int64) (int64, error) {
	value, err := rds.client.Get(context.TODO(), key).Result()

	if err != nil {
		return defaultValue, err
	}

	int64Value, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue, err
	}

	return int64Value, nil
}

func (rds *Redis) GetFloat64(key string, defaultValue float64) (float64, error) {
	value, err := rds.client.Get(context.TODO(), key).Result()

	if err != nil {
		return defaultValue, err
	}

	float64Value, err := strconv.ParseFloat(value, 10)
	if err != nil {
		return defaultValue, err
	}

	return float64Value, nil
}

func (rds *Redis) HGetSting(key, field string, defaultValue string) (string, error) {
	value, err := rds.client.HGet(context.TODO(), key, field).Result()

	if err != nil {
		return defaultValue, err
	}

	return value, nil
}

func (rds *Redis) HGetInt(key, field string, defaultValue int) (int, error) {
	value, err := rds.client.HGet(context.TODO(), key, field).Result()

	if err != nil {
		return defaultValue, err
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue, err
	}

	return intValue, nil
}

func (rds *Redis) HGetInt64(key, field string, defaultValue int64) (int64, error) {
	value, err := rds.client.HGet(context.TODO(), key, field).Result()

	if err != nil {
		return defaultValue, err
	}

	int64Value, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue, err
	}

	return int64Value, nil
}

func (rds *Redis) HGetFloat64(key, field string, defaultValue float64) (float64, error) {
	value, err := rds.client.HGet(context.TODO(), key, field).Result()

	if err != nil {
		return defaultValue, err
	}

	float64Value, err := strconv.ParseFloat(value, 10)
	if err != nil {
		return defaultValue, err
	}

	return float64Value, nil
}

func (rds *Redis) HGetAll(key string) (map[string]string, error) {

	result, err := rds.client.HGetAll(context.TODO(), key).Result()

	return result, err
}

func (rds *Redis) HMSet(key string, data map[string]interface{}) error {
	err := rds.client.HMSet(context.TODO(), key, data).Err()
	if err != nil {
		return err
	}

	return nil
}

func (rds *Redis) HMGet(key string, fields []string) (map[string]interface{}, error) {

	m := make(map[string]interface{})

	for _, f := range fields {
		if r, err := rds.client.HMGet(context.TODO(), key, f).Result(); err == nil {
			m[f] = r[0]
		}
	}

	return m, nil
}

func (rds *Redis) HIncrBy(key, field string, incr int64) (int64, error) {

	val, err := rds.client.HIncrBy(context.TODO(), key, field, incr).Result()
	if err != nil {
		return 0, err
	}

	return val, nil
}

func (rds *Redis) HMGetByFields(key string, fields ...string) (map[string]interface{}, error) {

	if result, err := rds.client.HMGet(context.TODO(), key, fields...).Result(); err != nil {
		return nil, err

	} else {
		m := make(map[string]interface{})

		for i, r := range result {
			m[fields[i]] = r
		}

		return m, nil
	}
}

func (rds *Redis) Exist(key string) int64 {
	result, _ := rds.client.Exists(context.TODO(), key).Result()

	return result
}

func (rds *Redis) HExistAndGetString(key, fields string) (string, bool, error) {
	isExist, _ := rds.client.HExists(context.TODO(), key, fields).Result()
	if isExist {
		result, err := rds.client.HGet(context.TODO(), key, fields).Result()
		if err != nil {
			return "", false, err
		} else {
			return result, true, nil
		}
	} else {
		return "", false, nil
	}
}

// ListsLength
// intervalMillisecond 每次 scan 間隔時間，大於0才啟用。
func (rds *Redis) ListsLength(keys []string, intervalMillisecond int64) map[string]int {
	m := map[string]int{}

	for _, v := range keys {
		n, _ := rds.LLen(v)
		m[v] = int(n)
		if intervalMillisecond > 0 {
			time.Sleep(time.Millisecond * time.Duration(intervalMillisecond))
		}
	}
	return m
}

func (rds *Redis) LLen(key string) (int64, error) {
	counts, err := rds.client.LLen(context.TODO(), key).Result()
	if err != nil {
		return 0, err
	}

	return counts, nil
}

func (rds *Redis) LRange(key string, start, stop int64, defaultValue []string) ([]string, error) {
	listLength, err := rds.client.LLen(context.TODO(), key).Result()
	if err != nil {
		return []string{}, errors.New("is not list")
	}

	if listLength == 0 {
		return defaultValue, errors.New("key is not exist")
	}

	if start > (listLength - 1) {
		return defaultValue, errors.New("index out of range")
	}

	if start >= 0 && stop >= 0 && start > stop {
		return defaultValue, errors.New("illegal index")
	} else if start < 0 && stop >= 0 {
		return defaultValue, errors.New("illegal index")
	} else if start < 0 && stop < 0 && start > stop {
		return defaultValue, errors.New("illegal index")
	}

	total, err := rds.client.LRange(context.TODO(), key, start, stop).Result()
	if err != nil {
		return defaultValue, err
	}

	return total, nil
}

func (rds *Redis) Expire(key string, expire time.Duration) bool {
	result, _ := rds.client.Expire(context.TODO(), key, expire).Result()

	return result
}

func (rds *Redis) Delete(key ...string) (int64, error) {
	numberOfKeyRemove, err := rds.client.Del(context.TODO(), key...).Result()

	if err != nil {
		return 0, err
	}

	return numberOfKeyRemove, nil
}

func (rds *Redis) Publish(channel string, message interface{}) error {

	_, err := rds.client.Publish(context.TODO(), channel, message).Result()
	if err != nil {
		return err
	}

	return nil
}

func (rds *Redis) Subscribe(channel ...string) *redis.PubSub {

	subscriber := rds.client.Subscribe(context.TODO(), channel...)

	return subscriber
}

func (rds *Redis) ZRevRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	val, err := rds.client.ZRevRangeWithScores(context.TODO(), key, start, stop).Result()
	if err != nil {
		return []redis.Z{}, err
	}

	return val, nil
}

func (rds *Redis) ZCard(key string) (int64, error) {
	counts, err := rds.client.ZCard(context.TODO(), key).Result()
	if err != nil {
		return 0, err
	}

	return counts, nil
}

func (rds *Redis) Scan(cursor uint64, match string, count int64) ([]string, uint64, error) {
	keys, newCursor, err := rds.client.Scan(context.TODO(), cursor, match, count).Result()

	if err != nil {
		return []string{}, 0, err
	}

	return keys, newCursor, nil
}

// TTLSeconds return ttl in seconds
// if key expired returns 0
// if key not exist returns -2
// if error returns -3
func (rds *Redis) TTLSeconds(key string) int64 {
	t, err := rds.client.TTL(context.TODO(), key).Result()
	if err != nil {
		return -3
	}
	return int64(t.Seconds())
}

// TTLMilliSeconds return ttl in milli seconds
// if key expired returns 0
// if key not exist returns -2
// if error returns -3
func (rds *Redis) TTLMilliSeconds(key string) int64 {
	t, err := rds.client.TTL(context.TODO(), key).Result()
	if err != nil {
		return -3
	}
	return int64(t.Milliseconds())

}

//ScanLiterally 循環 scan 直到找出所有 pattern
//func (rds *Redis) ScanLiterally(keys []string, cursor uint64, scanPattern string, scanCount int64) ([]string, error) {
//
//	gotKeys, gotCursor, err := rds.Scan(cursor, scanPattern, scanCount)
//	if err != nil {
//		return keys, fmt.Errorf("scan err %s", err)
//	}
//
//	keys = append(keys, gotKeys...)
//
//	if gotCursor != 0 {
//		return rds.ScanLiterally(keys, gotCursor, scanPattern, scanCount)
//	}
//	return keys, nil
//}

// ScanLiterally
// scanPattern , e.g. testPattern* 。scanCount 每次 scan 掃幾個 key。
// intervalMillisecond 每次 scan 間隔時間，大於0才啟用。
func (rds *Redis) ScanLiterally(scanPattern string, scanCount int64, intervalMillisecond int64) (retKeys []string, retErr error) {
	logPrefix := "ScanLiterally() "

	defer func() {
		if r := recover(); r != nil {
			retErr = fmt.Errorf("%s panic  %s", logPrefix, r)
		}
	}()

	retKeys = []string{}
	var scanCursor uint64 = 0
	for {
		keys, cursor, err := rds.Scan(scanCursor, scanPattern, scanCount)
		if err != nil {
			retErr = fmt.Errorf("%s redis scan() error %s", logPrefix, err)
			return
		}

		if len(keys) > 0 {
			retKeys = append(retKeys, keys...)
		}

		if cursor == 0 {
			break
		}
		scanCursor = cursor
		if intervalMillisecond > 0 {
			time.Sleep(time.Millisecond * time.Duration(intervalMillisecond))
		}
	}
	return
}

// 創建租約
//
//	租約成立: 返回 true , timestamp (id 到期時間, 以毫秒為單位)
//	租約失敗: 返回 false, timestamp
//	錯誤: 返回 false, 0
func (rds *Redis) LeaseID(key string, id string, ttl time.Duration) (bool, int64) {

	ctx := context.Background()
	pip := rds.client.Pipeline()
	now := float64(time.Now().UnixMilli()) / 1000

	// 先清除所有過期的 id, 再嘗試加入 id
	pip.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%f", now))
	add := pip.ZAddNX(ctx, key, redis.Z{
		Score:  now + ttl.Seconds(), // 到期時間(秒)
		Member: id,                  // 租用的id
	})
	res := pip.ZScore(ctx, key, id)
	_, err := pip.Exec(ctx)
	if err != nil {
		return false, 0
	}
	return add.Val() == 1, int64(res.Val() * 1000)
}

// 修改租約時效
//
//	成功: 返回 timestamp (id 到期時間, 以毫秒為單位)
//	失敗: 返回 0
func (rds *Redis) RenewLeaseID(key string, id string, ttl time.Duration) int64 {

	ctx := context.Background()
	pip := rds.client.Pipeline()
	now := float64(time.Now().UnixMilli()) / 1000

	// 先清除所有過期的 id 後, 嘗試延長租期, 再取得租約到期時間
	pip.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%f", now))
	pip.ZAddXX(ctx, key, redis.Z{
		Score:  now + ttl.Seconds(), // 到期時間(秒)
		Member: id,                  // 租用的id
	})
	res := pip.ZScore(ctx, key, id)
	_, err := pip.Exec(ctx)
	if err != nil {
		return 0
	}
	return int64(res.Val() * 1000)
}

func (rds *Redis) GetClient() *redis.Client {
	return rds.client
}

func (rds *Redis) Close() error {
	err := rds.client.Close()

	if err != nil {
		return err
	}

	return nil
}
