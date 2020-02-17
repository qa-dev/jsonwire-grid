package mongo

import (
	"context"
	"github.com/qa-dev/jsonwire-grid/pool/capabilities"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/qa-dev/jsonwire-grid/pool"
	"github.com/qa-dev/jsonwire-grid/storage/tests"
)

const (
	dbConnectionStringEnvVariable = "TEST_MONGO_CONNECTION"
)

var (
	MongoDriver PrepareMongo
)

type PrepareMongo struct {
	DbName           string
	Ctx              context.Context
	connectionClient *mongo.Client
}

func (p PrepareMongo) SetUp() {
	connectionString := os.Getenv(dbConnectionStringEnvVariable)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		panic("Database connection error: " + err.Error())
	}
	MongoDriver.Ctx = ctx
	MongoDriver.connectionClient = client
}

func (p PrepareMongo) TearDown() {
	_ = p.connectionClient.Disconnect(p.Ctx)
}

func (p PrepareMongo) CreateStorage() (pool.StorageInterface, func()) {

	dbName := tests.CreateDbName()
	db := MongoDriver.connectionClient.Database(dbName)
	storage := NewMongoStorage(db)
	deferFunc := func() {
		_ = storage.collection.Drop(MongoDriver.Ctx)
	}
	return storage, deferFunc
}
func TestMain(m *testing.M) {
	MongoDriver = PrepareMongo{}
	MongoDriver.SetUp()
	retCode := m.Run()
	MongoDriver.TearDown()
	os.Exit(retCode)
}

func TestMongo_Add(t *testing.T) {
	tests.TestStorage_Add(t, MongoDriver)
}
func TestMongo_Add_Repeat(t *testing.T) {
	tests.TestStorage_Add_Repeat(t, MongoDriver)

}
func TestStorage_Add_Limit_Overflow(t *testing.T) {
	t.Parallel()
	storage, deferFunc := MongoDriver.CreateStorage()
	defer deferFunc()
	node := pool.Node{
		Key:              "ololo",
		Address:          "ololo",
		CapabilitiesList: []capabilities.Capabilities{{"trololo": "lolo"}},
		Type:             pool.NodeTypePersistent,
	}
	limit := 1
	err := storage.Add(node, limit)
	assert.NotEmpty(t, err)
}

func TestMongo_GetAll(t *testing.T) {
	tests.TestStorage_GetAll(t, MongoDriver)
}

func TestMongo_GetByAddress(t *testing.T) {
	tests.TestStorage_GetByAddress(t, MongoDriver)
}

func TestMongo_GetBySession(t *testing.T) {
	tests.TestStorage_GetBySession(t, MongoDriver)
}

func TestMongo_GetCountWithStatus(t *testing.T) {
	tests.TestStorage_GetCountWithStatus(t, MongoDriver)
}
func TestMongo_Remove(t *testing.T) {
	tests.TestStorage_Remove(t, MongoDriver)
}

func TestMongo_ReserveAvailable_Positive(t *testing.T) {
	tests.TestStorage_ReserveAvailable_Positive(t, MongoDriver)
}

// TestMongo_ReserveAvailable_Negative see testStorage_ReserveAvailable_Negative
func TestMongo_ReserveAvailable_Negative(t *testing.T) {
	tests.TestStorage_ReserveAvailable_Negative(t, MongoDriver)
}

// TestMongo_SetAvailable see testStorage_SetAvailable
func TestMongo_SetAvailable(t *testing.T) {
	tests.TestStorage_SetAvailable(t, MongoDriver)
}

// TestMongo_SetBusy see testStorage_SetBusy
func TestMongo_SetBusy(t *testing.T) {
	tests.TestStorage_SetBusy(t, MongoDriver)
}

// TestMongo_UpdateAdderss_UpdatesValue see testStorage_UpdateAdderss_UpdatesValue
func TestMongo_UpdateAdderss_UpdatesValue(t *testing.T) {
	tests.TestStorage_UpdateAdderss_UpdatesValue(t, MongoDriver)
}

// TestMongo_UpdateAdderss_ReturnsErrNotFound see testStorage_UpdateAdderss_ReturnsErrNotFound
func TestMongo_UpdateAdderss_ReturnsErrNotFound(t *testing.T) {
	tests.TestStorage_UpdateAdderss_ReturnsErrNotFound(t, MongoDriver)
}
