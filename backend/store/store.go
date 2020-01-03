package store

import (
	"encoding/json"
	"github.com/dlyahov/startuplink-web-go/backend/model"
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
	"log"
)

const dbFileName = "app.db"

type Storage interface {
	SaveUser(user *model.User) error
	FindUser(id string) (*model.User, error)
}

var (
	ErrUserNotFound = errors.New("user not found")
)

type BoltDb struct {
	db *bolt.DB
}

const userBuckets = "users"

func NewStorage(options *bolt.Options) (*BoltDb, error) {
	log.Printf("Get bolt store with options %+v", options)
	boltDb, err := bolt.Open(dbFileName, 0666, options)
	if err != nil {
		log.Fatal("Could not open database file")
		return nil, err
	}
	err = boltDb.Update(func(tx *bolt.Tx) error {
		if _, e := tx.CreateBucketIfNotExists([]byte(userBuckets)); e != nil {
			return errors.Wrapf(e, "failed to create bucket %s", userBuckets)
		}
		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to create bucket)")
	}

	log.Printf("Bolt store created")
	return &BoltDb{boltDb}, nil
}

func (boltDb *BoltDb) SaveUser(user model.User) error {
	db := boltDb.db
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(userBuckets))
		userJson, err := json.Marshal(user)
		if err != nil {
			log.Fatal("Could not marshal user to json")
			return err
		}
		err = bucket.Put([]byte(user.Id), userJson)
		if err != nil {
			return errors.Wrapf(err, "failed to save user data. User id %s", user.Id)
		}
		return nil
	})

	if err != nil {
		log.Fatal("failed to save user")
	}
	return nil
}

func (boltDb *BoltDb) FindUser(id string) (model.User, error) {
	db := boltDb.db

	user := model.User{}
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(userBuckets))
		userBytes := bucket.Get([]byte(id))
		if userBytes == nil {
			return ErrUserNotFound
		}

		err := json.Unmarshal(userBytes, &user)
		if err != nil {
			return errors.Wrapf(err, "Could not marshal user with id %s", id)
		}
		return nil
	})

	if err != nil {
		log.Fatal("failed to save user")
		return user, err
	}
	return user, nil
}
