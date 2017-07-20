package main

import (
	"time"
	"os"
	"path/filepath"
	"github.com/boltdb/bolt"
	"encoding/json"
	"log"
)

const (
	defaultCacheFolder = "asset_viewer_cache"
	defaultCacheBucket = "cache"
)


type boltCache struct {
	// path to root folder where boltdb files will be written to
	root string

	// name of primary cache bucket inside of boltdb
	cacheBucket []byte

	// boltdb connection
	db *bolt.DB

	// how long data will be cached for
	ttl time.Duration
}

func newBoltCache(opts ...boltOpt) (assetCache, error) {
	cache := &boltCache{
		root: filepath.Join(os.TempDir(), defaultCacheFolder),
		ttl: DefaultCacheTTL,
		cacheBucket: []byte(defaultCacheBucket),
	}

	for _, opt := range opts {
		opt(cache)
	}

	log.Println("Opening BoltDB cache", cache.root, "TTL:", DefaultCacheTTL)
	return cache, nil
}

type boltOpt func(*boltCache)

func rootFolder(location string) boltOpt {
	return func(b *boltCache) {
		b.root = location
	}
}

func ttl(t time.Duration) boltOpt {
	return func(b *boltCache) {
		b.ttl = t
	}
}

// Called on start up, this opens / creates a boltdb bucket for writing
func (b *boltCache) createBucket(name []byte) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(name)
		return err
	})
}

// connect to / setup the cache
func (b *boltCache) Connect() error {
	db, err := bolt.Open(b.root, 0600, nil)
	if err != nil {
		return err
	}
	b.db = db
	return b.createBucket(b.cacheBucket)
}

// Return cached match suggestions for search indicated by given key
func (b *boltCache) Matches(in string) ([]*AssetDescription, error) {
	data, err := b.lookupCache(in)
	if err != nil {
		return nil, err
	}

	if data == nil { // cache miss
		return nil, nil
	}

	result := []*AssetDescription{}
	fromCache := []AssetDescription{}

	err = json.Unmarshal(data, &fromCache)
	for i := range fromCache {
		result = append(result, &fromCache[i])
	}
	return result, nil
}

// Retrieve raw data from cache. Note this will not return data that has passed it's TTL
func (b *boltCache) lookupCache(in string) ([]byte, error) {
	key := []byte(in)
	rawCachedData := []byte{}

	b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(b.cacheBucket)
		rawCachedData = bucket.Get(key)
		return nil
	})

	if rawCachedData == nil {
		return nil, nil
	}

	cached := &cachedData{}
	err := json.Unmarshal(rawCachedData, cached)
	if err != nil {
		return nil, err
	}

	if time.Now().Unix() >= cached.Expire { // Cache hit, but data is expired
		// Remove expired data from cache
		go func() {
			b.db.Update(func(tx *bolt.Tx) error {
				return tx.Bucket(b.cacheBucket).Delete(key)
			})
		}()

		// Pretend we missed cache
		return nil, nil
	}

	return cached.Data, nil
}

// Return cached asset data for search indicated by given key
func (b *boltCache) AssetData(in string) (*AssetData, error) {
	data, err := b.lookupCache(in)
	if err != nil {
		return nil, err
	}

	if data == nil { // cache miss
		return nil, nil
	}

	result := &AssetData{}
	err = json.Unmarshal(data, result)
	return result, err
}

// Close connection(s)
func (b *boltCache) Close() error {
	// close bolt
	return b.db.Close()
}

// Internal boltdb struct that caches some data along with ttl
type cachedData struct {
	Data []byte
	Expire int64
}

// Cache something with the given key
func (b *boltCache) Cache(key string, data interface{}) error {
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}

	cached, err := json.Marshal(&cachedData{
		Data: raw,
		Expire: time.Now().Add(b.ttl).Unix(),
	})
	if err != nil {
		return err
	}

	return b.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(b.cacheBucket).Put([]byte(key), cached)
	})
}
