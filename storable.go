package storage

// StorableInfo provides the basic metadata every StorageObject will need.
// It is not necessary to embed this struct if the fields are implemented by
// the StorageObject, however, use is highly encouraged to support the addition
// of new metadata.
type StorableInfo struct {
	Created int64 `gorethink:"created,omitempty" json:"created"`
	Updated int64 `gorethink:"updated,omitempty" json:"updated"`
}

// Storable provides the basic metadata every StorageObject will need, it also
// provides a basic implementation of the StorageObject interface.
//
// It is not necessary to embed this struct, however, use is highly encouraged
// to support the addition of new metadata and CRUD operations.
type Storable struct {
	StorableInfo `gorethink:"info"`
}

// Created is a getter/setter for a StorageObjects Created property.
func (s *Storable) Created(t ...int64) int64 {
	if t != nil {
		s.StorableInfo.Created = t[0]
	}

	return s.StorableInfo.Created
}

// Updated is a getter/setter for a StorageObjects Updated property.
func (s *Storable) Updated(t ...int64) int64 {
	if t != nil {
		s.StorableInfo.Updated = t[0]
	}

	return s.StorableInfo.Updated
}
