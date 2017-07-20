package main

import (
	"time"
	"fmt"
	"log"
)

const (
	// default page size for result sets
	PageSize = 100

	// default time a query will be cached internally
	DefaultCacheTTL = time.Second * 15
)


// Interface to some asset database system
type assetDB interface {
	// Connect to the asset db
	Connect() error

	// Return an initial list of suggestions
	SuggestInitial() ([]string, error)

	// Given a lose description (w/ partial or null fields), find possible matching assets
	Matches(*AssetDescription, int) ([]*AssetDescription, error)

	// Retrieve detailed asset data given the short description
	AssetData(*AssetDescription) (*AssetData, error)

	// Close connection(s)
	Close() error
}

// Interface to an asset cache -- this will be used in conjunction with the assetDB to
// cache information from the data source.
type assetCache interface {
	// connect to / setup the cache
	Connect() error

	// Retrieve AssetDescriptions from cache with the given key
	Matches(string) ([]*AssetDescription, error)

	// Retrieve AssetData from cache with the given key
	AssetData(string) (*AssetData, error)

	// Close connection(s)
	Close() error

	// Cache something with the given key
	Cache(string, interface{}) error
}

// AssetService is a union of some assetDB and an assetCache that can be used to find,
// retrieve and cache asset information
type AssetService struct {
	db assetDB
	cache assetCache
}

// Create a new AssetService that reads from the given DB
func NewAssetService(db assetDB) (*AssetService, error) {
	err := db.Connect()
	if err != nil {
		return nil, err
	}

	cache, err := newBoltCache()
	if err != nil {
		return nil, err
	}

	err = cache.Connect()
	if err != nil {
		return nil, err
	}

	service := &AssetService{
		db: db,
		cache: cache,
	}
	return service, nil
}

// Close connections to db and underlying cache
func (s *AssetService) Close() {
	s.db.Close()
	s.cache.Close()
}

// Return a list of initial suggestions
func (s *AssetService) SuggestInitial() ([]string, error) {
	return s.db.SuggestInitial()
}

// Try to find potential Matches for assets matching the given fields.
// Results are cached for some TTL
func (s *AssetService) Matches(name, class, subclass string, page int) ([]*AssetDescription, error) {
	cacheKey := fmt.Sprintf("match:%d:%s:%s:%s", page, name, class, subclass)
	cResult, err := s.cache.Matches(cacheKey)
	if err != nil {
		// The cache failing isn't fatal to the whole process, log as warning
		log.Println("Unable to read from cache:", err.Error())
	}
	if cResult != nil {
		// Cache hit!
		return cResult, nil
	}

	desired := &AssetDescription{
		Name: name,
		Class: class,
		Subclass: subclass,
	}

	// Look up in the database
	dResult, err := s.db.Matches(desired, page)
	if err != nil {
		return nil, err
	}

	// Be sure to cache the result
	s.cache.Cache(cacheKey, dResult)

	return dResult, nil
}

// Retrieve detailed asset data for a single asset matching the given fields.
// Results are cached for some TTL
func (s *AssetService) AssetData(name, class, subclass string) (*AssetData, error) {
	cacheKey := fmt.Sprintf("fetch:%s:%s:%s", name, class, subclass)
	cResult, err := s.cache.AssetData(cacheKey)
	if err != nil {
		// The cache failing isn't fatal to the whole process, log as warning
		log.Println("Unable to read from cache:", err.Error())
	}
	if cResult != nil {
		// Cache hit!
		return cResult, nil
	}

	desired := &AssetDescription{
		Name: name,
		Class: class,
		Subclass: subclass,
	}

	// Look up in the database
	dResult, err := s.db.AssetData(desired)
	if err != nil {
		return nil, err
	}

	// Be sure to cache the result
	s.cache.Cache(cacheKey, dResult)

	return dResult, nil
}