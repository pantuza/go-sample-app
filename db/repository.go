package db

import (
	"encoding/json"
	"errors"
	"github.com/pantuza/go-sample-app/music"
	"gopkg.in/redis.v5"
	"log"
)

var ErrMusicNotFound = errors.New("Music Not Found")
var ErrMusicDecode = errors.New("Music decode error")

type MusicRepository struct {
	session *redis.Client
}

func (mr *MusicRepository) Create(m *music.Music) {

	value, _ := json.Marshal(m)
	err := mr.session.Set(m.Name, string(value), 0).Err()
	if err != nil {
		log.Println("Failed to set music:", m.Name)
	}
}

func (mr *MusicRepository) Get(key string) (*music.Music, error) {
	result, err := mr.session.Get(key).Result()
	if err != nil {
		log.Println("Failed to get key: ", key)
		return nil, ErrMusicNotFound
	}
	m := &music.Music{}
	if err := json.Unmarshal([]byte(result), m); err != nil {
		log.Println("Failed to decode json", err)
		return nil, ErrMusicDecode
	}

	return m, nil
}

func (mr *MusicRepository) Update(m *music.Music) error {

	if _, err := mr.Get(m.Name); err != nil {
		return ErrMusicNotFound
	}
	value, _ := json.Marshal(m)
	err := mr.session.Set(m.Name, string(value), 0).Err()
	if err != nil {
		log.Println("Failed to update key: ", m.Name)
	}
	return nil
}

func (mr *MusicRepository) Delete(key string) error {

	if _, err := mr.Get(key); err != nil {
		return ErrMusicNotFound
	}

	err := mr.session.Del(key).Err()
	if err != nil {
		log.Println("Failed to delete:", key)
		return err
	}

	return nil
}

func (mr *MusicRepository) List() *[]music.Music {
	musics := make([]music.Music, 0, 10)
	var keys []string
	keys, _ = mr.session.Keys("*").Result()
	for _, v := range keys {
		m, _ := mr.Get(v)
		musics = append(musics, *m)
	}
	return &musics
}

func New() *MusicRepository {
	session := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	log.Println("Redis connection: OK")
	return &MusicRepository{session}
}
