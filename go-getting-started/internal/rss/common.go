package rss

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"

	"log"
	"time"

	"github.com/adjust/rmq/v4"

	m "go-getting-started/internal/common"
)

const (
	NUM_ELEMENTS_IN_REDIS_KEY_RSS = 1000
	//consumer
	prefetchLimit = 1000
	pollDuration  = 100 * time.Millisecond
	numConsumers  = 5

	reportBatchSize = 10000
	consumeDuration = time.Millisecond
	shouldLog       = false

	//producer
	numDeliveries = 100000000
	batchSize     = 10000
)

var (
	ErrNil = errors.New("no matching record found in redis database")
	ctx    = context.Background()
	Ctx    = context.TODO()
)

//initiation
func (rss *Rss) Init(address string, password string) error {
	rss.Redis = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0, //2
	})
	pong, err := rss.Redis.Ping(Ctx).Result()
	fmt.Println(pong, err) //todo remove
	if err != nil {
		fmt.Println(err)
		return err
	}
	rss.IsInitialized = true
	return nil
}

func (rss *Rss) IsFollowerExistsInRedis(uid uint64) (bool, error) {
	return rss.IsFollowerExistsInRedis_s(strconv.FormatUint(uid, 10))
}

func (rss *Rss) IsFollowerExistsInRedis_s(followee_id string) (bool, error) {

	key := "followee_" + followee_id
	numInRedis, err := rss.Redis.ZCard(ctx, key).Uint64()
	if err != nil {
		return false, err
	}

	if numInRedis > 0 {
		return true, nil
	}
	return false, nil
}

func (rss *Rss) PushFollower(follower_id string, followee_id string) error {
	println("PushFollower")
	key := "followee_" + followee_id

	i, err := strconv.ParseInt(follower_id, 10, 64)
	if err != nil {
		return err
	}
	_, err = rss.Redis.ZAdd(ctx, key, &redis.Z{
		Score:  float64(i),
		Member: follower_id}).Result()

	return err
}

func (rss *Rss) PushFollower_int64(follower_id int64, followee_id int64) error {

	key := "followee_" + strconv.FormatInt(followee_id, 10)

	_, err := rss.Redis.ZAdd(ctx, key, &redis.Z{
		Score:  float64(follower_id),
		Member: strconv.FormatInt(follower_id, 10)}).Result()

	return err
}

func (rss *Rss) PushPost(post *m.Post) error {

	p, err := json.Marshal(post)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return err //todo set continue
	}

	key_f := "followee_" + strconv.Itoa(post.AuthorId) //todo int64
	results_f, err := rss.Redis.ZRevRange(ctx, key_f, 0, -1).Result()
	if err != nil {
		return err
	}

	for i := 0; i < len(results_f); i++ {

		follower := results_f[i]
		key := "all_posts_for_user_id_" + follower
		exist, _ := rss.IsPostsExistInRedis_s(follower) //TODO:? add ONLY if REDIS Record exists
		if exist {
			// println("exist")
			_, err = rss.Redis.ZAdd(ctx, key, &redis.Z{
				Score:  float64(post.CreatedAt.Unix()), //Todo add something to score
				Member: p}).Result()
			if err != nil {
				fmt.Printf("Error: %s", err)
				return err //todo set continue
			}
		}
	}

	return nil
}

func (rss *Rss) PushPosts(uid int64, posts []m.Post) error {

	key := "all_posts_for_user_id_" + strconv.FormatInt(uid, 10)

	for i := 0; i < len(posts); i++ {
		post := posts[i]
		p, err := json.Marshal(post)
		if err != nil {
			fmt.Printf("Error: %s", err)
			return err //todo set continue
		}

		res, err := rss.Redis.ZAdd(ctx, key, &redis.Z{
			Score:  float64(post.CreatedAt.Unix()), //Todo add something to score
			Member: p}).Result()
		if err != nil {
			fmt.Printf("Error: %s", err)
			return err //todo set continue
		}
		println(res)
	}

	return nil
}

func (rss *Rss) IsPostsExistInRedis(uid int64) (bool, error) {
	return rss.IsPostsExistInRedis_s(strconv.FormatInt(uid, 10))
}

func (rss *Rss) IsPostsExistInRedis_s(uid string) (bool, error) {

	key := "all_posts_for_user_id_" + uid
	numInRedis, err := rss.Redis.ZCard(ctx, key).Uint64()
	if err != nil {
		return false, err
	}

	if numInRedis > 0 {
		return true, nil
	}
	return false, nil
}

func (rss *Rss) GetPosts(uid int64) ([]m.Post, error) {
	return rss.GetPosts_s(strconv.FormatInt(uid, 10))
}

func (rss *Rss) GetPosts_s(uid string) ([]m.Post, error) {

	key := "all_posts_for_user_id_" + uid
	println(key)

	numInRedis, err := rss.Redis.ZCard(ctx, key).Uint64()
	if err != nil {
		return nil, err
	}

	if numInRedis > NUM_ELEMENTS_IN_REDIS_KEY_RSS {
		rss.Redis.ZRemRangeByRank(ctx, key, 0, int64(numInRedis)-int64(NUM_ELEMENTS_IN_REDIS_KEY_RSS))
	}

	results, err := rss.Redis.ZRevRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var posts []m.Post
	for i := 0; i < len(results); i++ {
		var post m.Post
		err := json.Unmarshal([]byte(results[i]), &post)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	println("num posts")
	println(len(posts))

	return posts, nil
}

//we do not use this function...............
func (rss *Rss) CreateQueue(address string) {
	// connection, err := rmq.OpenConnection("producer", "tcp", address, 2, nil)
	// if err != nil {
	// 	panic(err)
	// }

	connection, err := rmq.OpenConnectionWithRedisClient("producer", rss.Redis, nil)
	if err != nil {
		panic(err)
	}

	things, err := connection.OpenQueue("things")
	if err != nil {
		panic(err)
	}
	foobars, err := connection.OpenQueue("foobars")
	if err != nil {
		panic(err)
	}

	var before time.Time
	for i := 0; i < numDeliveries; i++ {
		delivery := fmt.Sprintf("delivery %d", i)
		if err := things.Publish(delivery); err != nil {
			log.Printf("failed to publish: %s", err)
		}

		if i%batchSize == 0 {
			duration := time.Now().Sub(before)
			before = time.Now()
			perSecond := time.Second / (duration / batchSize)
			log.Printf("produced %d %s %d", i, delivery, perSecond)
			if err := foobars.Publish("foo"); err != nil {
				log.Printf("failed to publish: %s", err)
			}
		}
	}
}
