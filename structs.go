package main

type AssetDescription struct {
	Id string `json:"-"`

	// Asset group / collection name
	Name string

	// Asset major class / type
	Class string

	// Asset subtype
	Subclass string

	// Optional string left to assetDB to set as desired
	Description string
}

// A loosely defined way of describing some 'resource' connected to an asset
type ResourceDescription struct {
	// Resource name indicating designation
	Name string

	// Resource type indicating intended use
	Class string

	// Some kind of resource URI
	URI string
}

type AssetData struct {
	// Asset description data
	Data *AssetDescription

	// Asset attributes (arbitrary key: value pairs)
	Attributes map[string]string

	// Asset revision
	Version int

	// Optional thumbnail that could be loaded
	Thumbnail string

	// Assets that this is related to
	Linked []*AssetDescription

	// Resources of this asset
	Resources []*ResourceDescription
}
