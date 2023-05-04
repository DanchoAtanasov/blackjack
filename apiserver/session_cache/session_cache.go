package sessioncache

import (
	env "apiserver/environment"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis"
)

type sessionCacheInterface interface {
	StoreSession()
	GetSession(string)
	DeleteSession()
}

type SessionCache struct {
	client *redis.Client
}

func NewSessionCache() *SessionCache {
	sessionCache := SessionCache{
		client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:6379", env.REDIS_HOST),
			Password: "",
			DB:       0,
		}),
	}
	return &sessionCache
}

type PlayerSessionInformation struct {
	Name    string
	BuyIn   int
	CurrBet int
}

func (sessionCache SessionCache) StoreSession(token string, psi PlayerSessionInformation) {
	json, err := json.Marshal(psi)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = sessionCache.client.Set(token, json, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
}

// Redis functions to be moved to a separate file
func (sessionCache SessionCache) GetSession(token string) PlayerSessionInformation {
	val, err := sessionCache.client.Get(token).Result()
	if err != nil {
		fmt.Println(err)
		return PlayerSessionInformation{}
	}
	fmt.Printf("got from redis %s\n", val)

	var playerSession PlayerSessionInformation
	err = json.Unmarshal([]byte(val), &playerSession)
	if err != nil {
		fmt.Println("failed to unmarshal")
		return PlayerSessionInformation{}
	}

	return playerSession
}

func (sessionCache SessionCache) DeleteSession(sessionId string) error {
	err := sessionCache.client.Del(sessionId).Err()
	if err != nil {
		fmt.Println(err)
		return errors.New("Could not delete session")
	}

	return nil
}
