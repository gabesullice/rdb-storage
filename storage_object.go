package storage

// StorageObject defines an interface that allows data models to integrate with
// the storage package and utilize its helper methods.
type StorageObject interface {
	Create() error
	Validate() error

	Id(...string) string
	Created(...int64) int64
	Updated(...int64) int64
	ModelName() string
}
