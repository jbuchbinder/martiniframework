package martiniframework

import (
	"gopkg.in/redis.v3"
	"log"
	"time"
)

type SessionDriverRedis struct {
	client *redis.Client
}

func (s *SessionDriverRedis) Connect(params map[string]interface{}) error {
	s.client = redis.NewClient(&redis.Options{
		Addr:     params["Address"].(string),
		Password: params["Password"].(string),
		DB:       params["DatabaseId"].(int64),
	})
	_, err := s.client.Ping().Result()
	return err
}

func (s *SessionDriverRedis) Get(sid string) (string, error) {
	val, err := s.client.Get(sid).Result()
	if err == redis.Nil {
		log.Printf("GetSession(): %s does not exist", sid)
		return "", err
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

func (s *SessionDriverRedis) Del(sid string) error {
	return s.client.Del(sid).Err()
}

func (s *SessionDriverRedis) Set(sid string, data string, expiry time.Duration) error {
	return s.client.Set(sid, data, expiry).Err()
}
