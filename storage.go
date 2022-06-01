package dataset

// StorageSystem describes the functions required to implement
// a dataset storage system. Currently two types of storage systems
// are supported -- pairtree and sql storage (via MySQL 8 and JSON columns)
// If the funcs describe are not supported by the storage system they
// must return a "Not Implemented" error value.
type StorageSystem interface {

	// Open opens the storage system and returns an storage struct and error
	// It is passed either a filename. For a Pairtree the would be the
	// path to collection.json and for a sql store file holding a DSN
	//
	// ```
	//  store, err := c.Store.Open(c.Access)
	//  if err != nil {
	//     ...
	//  }
	// ```
	//
	Open(name string, dsnURI string) (*StorageSystem, error)

	// Close closes the storage system freeing resources as needed.
	//
	// ```
	//   if err := storage.Close(); err != nil {
	//      ...
	//   }
	// ```
	//
	Close() error

	// Create stores a new JSON object in the collection
	// It takes a string as a key and a byte slice of encoded JSON
	//
	//   err := storage.Create("123", []byte(`{"one": 1}`))
	//   if err != nil {
	//      ...
	//   }
	//
	Create(string, []byte) error

	// Read retrieves takes a string as a key and returns the encoded
	// JSON document from the collection
	//
	//   src, err := storage.Read("123")
	//   if err != nil {
	//      ...
	//   }
	//   obj := map[string]interface{}{}
	//   if err := json.Unmarshal(src, &obj); err != nil {
	//      ...
	//   }
	Read(string) ([]byte, error)

	// Versions returns a list of semver formatted version strings avialable for an JSON object
	Versions(string) ([]string, error)

	// ReadVersion takes a key and semver version string and return that version of the
	// JSON object.
	ReadVersion(string, string) ([]byte, error)

	// Update takes a key and encoded JSON object and updates a
	// JSON document in the collection.
	//
	//   key := "123"
	//   src := []byte(`{"one": 1, "two": 2}`)
	//   if err := storage.Update(key, src); err != nil {
	//      ...
	//   }
	//
	Update(string, []byte) error

	// Delete removes all versions and attachments of a JSON document.
	//
	//   key := "123"
	//   if err := storage.Delete(key); err != nil {
	//      ...
	//   }
	//
	Delete(string) error

	// Keys returns all keys in a collection as a slice of strings.
	//
	//   var keys []string
	//   keys, _ = storage.List()
	//   /* iterate over the keys retrieved */
	//   for _, key := range keys {
	//      ...
	//   }
	//
	Keys() ([]string, error)

	// HasKey returns true if collection is open and key exists,
	// false otherwise.
	HasKey(string) bool

	// Length returns the number of records in the collection
	Length() int64
}
