package main

import (
	wclient "github.com/voidshard/wysteria/client"
	wyc "github.com/voidshard/wysteria/common"
	wym "github.com/voidshard/wysteria/common/middleware"
	"errors"
	"fmt"
)

const (
	defaultWysteriaUrl = "nats://localhost:4222"
	defaultWysteriaDriver = wym.DriverNats
)

// Wrapper for wysteria connection implements assetDB
type wysteria struct {
	url  string
	driver string
	conn *wclient.Client
}

type WysteriaOpt func(*wysteria)

// Set connection url to server
func Host(url string) WysteriaOpt {
	return func(w *wysteria) {
		w.url = url
	}
}

// Set the middleware driver
func Driver(name string) WysteriaOpt {
	return func(w *wysteria) {
		w.driver = name
	}
}

// Create new wysteria connection (implements assetDB)
func NewWysteriaDB(opts ...WysteriaOpt) assetDB {
	db := &wysteria{
		url: defaultWysteriaUrl,
		driver: defaultWysteriaDriver,
	}

	for _, opt := range opts {
		opt(db)
	}

	return db
}

// connect to wysteria server
func (w *wysteria) Connect() error {
	client, err := wclient.New(wclient.Host(w.url), wclient.Driver(w.driver))
	if err != nil {
		return nil
	}

	w.conn = client
	return nil
}

// Return some options for search suggestions to the client
func (w *wysteria) SuggestInitial() (results []string, err error) {
	collections, err := w.conn.Search().FindCollections(wclient.Limit(1))
	if err != nil {
		return
	}
	for _, c := range collections {
		results = append(results, c.Name())
	}
	return
}

// Convert a wysteria Item into our internal asset description
func itemToAssetDescription(i *wclient.Item) *AssetDescription {
	value, _ := i.Facet(wyc.FacetCollection)
	return &AssetDescription{
		Id: i.Id(),
		Name: value,
		Class: i.Type(),
		Subclass: i.Variant(),
	}
}

// Convert a wysteria Version into our internal asset description
func versionToAssetDescription(v *wclient.Version) *AssetDescription {
	collectionName, _ := v.Facet(wyc.FacetCollection)
	itemType, _ := v.Facet(wyc.FacetItemType)
	itemVariant, _ := v.Facet(wyc.FacetItemVariant)
	return &AssetDescription{
		Id: v.Id(),
		Name: collectionName,
		Class: itemType,
		Subclass: itemVariant,
		Description: fmt.Sprintf("v%d", v.Version()),
	}
}

// Convert wysteria resource to resource internal resource desc
func resourceToDescription(r *wclient.Resource) *ResourceDescription {
	return &ResourceDescription{
		Name:  r.Name(),
		Class: r.Type(),
		URI:   r.Location(),
	}
}

// Given the source item and it's detailed data, add in data about the latest published version (if any)
func appendVersionData(i *wclient.Item, data *AssetData) (*AssetData, error) {
	published, err := i.PublishedVersion()
	if err != nil {
		return data, err
	}

	if published == nil {
		return data, nil
	}

	data.Version = int(published.Version())
	data.Attributes = i.Facets()

	linkedVersions, err := published.Linked()
	if err != nil {
		return data, err
	}

	for _, versions := range linkedVersions {
		for _, version := range versions {
			data.Linked = append(data.Linked, versionToAssetDescription(version))
		}
	}

	resources, err := published.Resources()
	for _, resource := range resources {
		data.Resources = append(data.Resources, resourceToDescription(resource))
	}
	return data, nil
}

// Extract detailed information about the given item & it's published version (if any)
func itemToAssetData(i *wclient.Item) (*AssetData, error) {
	allLinked := []*AssetDescription{}

	linkedItems, err := i.Linked()
	if err != nil {
		return nil, err
	}

	for _, items := range linkedItems {
		for _, item := range items {
			allLinked = append(allLinked, itemToAssetDescription(item))
		}
	}

	data := &AssetData{
		Data: itemToAssetDescription(i),
		Linked: allLinked,
	}
	return appendVersionData(i, data)
}

// Given a lose description (w/ partial or null fields), find possible matching assets
func (w *wysteria) Matches(q *AssetDescription, page int) ([]*AssetDescription, error) {
	if q.Name == "" {
		// With no collection name, we could search for all assets matching in all collections
		// but that could be expensive. Best wait till we're given a Collection name.
		return nil, nil
	}

	items, err := w.conn.Search(
		wclient.HasFacets(map[string]string{
			wyc.FacetCollection: q.Name,
		}),
		wclient.ItemType(q.Class),
		wclient.ItemVariant(q.Subclass),
	).FindItems(wclient.Limit(PageSize), wclient.Offset(page * PageSize))

	if err != nil {
		return nil, nil
	}

	result := []*AssetDescription{}
	for _, item := range items {
		result = append(result, itemToAssetDescription(item))
	}
	return result, nil
}

// Retrieve detailed asset data given the short description
func (w *wysteria) AssetData(q *AssetDescription) (*AssetData, error) {
	if q.Name == "" || q.Class == "" || q.Subclass == "" {
		return nil, errors.New("Required Description to include non null Name, Class, Subclass fields")
	}

	items, err := w.conn.Search(
		wclient.HasFacets(map[string]string{
			wyc.FacetCollection: q.Name,
		}),
		wclient.ItemType(q.Class),
		wclient.ItemVariant(q.Subclass),
	).FindItems(wclient.Limit(1))

	if err != nil {
		return nil, err
	}

	if len(items) == 1 {
		return itemToAssetData(items[0])
	}
	return nil, errors.New(fmt.Sprintf("Expected 1 result got %d", len(items)))
}

// Close connection to server
func (w *wysteria) Close() error {
	w.conn.Close()
	return nil
}

