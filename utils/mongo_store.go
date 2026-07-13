package goutils

type MongoConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

type MongoStore struct {
}

type MongoStoreApi interface {
	Connect() error
	Close() error
	IsReady() bool
	IsCollectionExists(collectionName string) (bool, error)
}

func NewMongoStore(config *MongoConfig) *MongoStore {
	return &MongoStore{}
}

func (m *MongoStore) IsCollectionExists(collectionName string) (bool, error) {
	// Implement the logic to check if the collection exists in the MongoDB database
	// You can use the MongoDB Go driver to perform this operation
	// For example, you can use the ListCollectionNames method to get a list of collections and check if the specified collection exists
	return false, nil
}
