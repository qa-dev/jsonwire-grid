package mongo

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/qa-dev/jsonwire-grid/pool"
	"github.com/qa-dev/jsonwire-grid/storage"
)

type Storage struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewMongoStorage(db *mongo.Database) *Storage {
	return &Storage{
		collection: db.Collection("grid"),
		ctx:        context.Background(),
	}
}

func (m *Storage) Add(node pool.Node, limit int) error {
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"key": node.Key}
	update := bson.D{{Key: "$set", Value: node}}

	if limit > 0 {
		return errors.New("[Mongo/Add] limit strategy temporary not supported ")
	}

	result, err := m.collection.UpdateOne(m.ctx, filter, update, opts)
	if err != nil {
		return errors.New("[Mongo/Add] update node, " + err.Error())
	}
	if result.MatchedCount == 0 && result.UpsertedCount == 0 {
		return errors.New("[Mongo/Add] No rows was affected ")
	}

	return nil
}

func (m *Storage) ReserveAvailable(nodeList []pool.Node) (pool.Node, error) {
	nodeKeyList := make([]string, 0, len(nodeList))
	node := pool.Node{}
	for _, node := range nodeList {
		nodeKeyList = append(nodeKeyList, node.Key)
	}

	filter := bson.M{"key": bson.M{"$in": nodeKeyList}, "status": pool.NodeStatusAvailable}

	update := bson.M{"$set": bson.M{
		"updated": time.Now().Unix(),
		"status":  string(pool.NodeStatusReserved),
	}}

	opts := options.
		FindOneAndUpdate().
		SetReturnDocument(options.After).
		SetSort(bson.M{"updated": 1})

	err := m.collection.FindOneAndUpdate(m.ctx, filter, update, opts).Decode(&node)
	if err != nil {
		err = errors.New("[Mongo/ReserveAvailable] find and update node, " + err.Error())
		return node, err
	}
	return node, nil
}

func (m *Storage) SetBusy(node pool.Node, sessionID string) error {
	filter := bson.M{"key": node.Key}
	update := bson.M{"$set": bson.M{
		"session_id": sessionID,
		"updated":    time.Now().Unix(),
		"status":     string(pool.NodeStatusBusy),
	}}
	result, err := m.collection.UpdateOne(m.ctx, filter, update)
	if err != nil {
		err = errors.New("[Mongo/SetBusy] update node in collection, " + err.Error())
		return err
	}
	if result.ModifiedCount == 0 {
		return storage.ErrNotFound
	}
	return nil
}

func (m *Storage) SetAvailable(node pool.Node) error {
	filter := bson.M{"key": node.Key}
	update := bson.M{"$set": bson.M{
		"updated": time.Now().Unix(),
		"status":  string(pool.NodeStatusAvailable),
	}}
	result, err := m.collection.UpdateOne(m.ctx, filter, update)
	if err != nil {
		err = errors.New("[Mongo/SetAvailable] update node in collection, " + err.Error())
		return err
	}
	if result.ModifiedCount == 0 {
		return storage.ErrNotFound
	}
	return nil
}

func (m *Storage) GetCountWithStatus(status *pool.NodeStatus) (int, error) {
	var count int
	var err error
	filter := bson.M{}
	if status != nil {
		filter = bson.M{"status": status}
	}
	cnt, err := m.collection.CountDocuments(m.ctx, filter)
	count = int(cnt)
	if err != nil {
		return 0, errors.New("[Mongo/GetCountWithStatus] count nodes in collection, " + err.Error())
	}
	return count, nil
}

func (m *Storage) GetBySession(session string) (pool.Node, error) {
	node, err := m.getByField("session_id", session)
	if err != nil {
		return pool.Node{}, errors.New("[Mongo/GetBySession] find node in collection, " + err.Error())
	}
	return node, nil
}

func (m *Storage) GetByAddress(address string) (pool.Node, error) {
	node, err := m.getByField("address", address)
	if err != nil {
		return pool.Node{}, errors.New("[Mongo/GetByAddress] find node in collection, " + err.Error())
	}
	return node, nil
}

func (m *Storage) GetAll() ([]pool.Node, error) {
	nodeList := make([]pool.Node, 0)
	resultCursor, err := m.collection.Find(m.ctx, bson.M{})
	if err != nil {
		return nil, errors.New("[Mongo/GetAll] find node in collection, " + err.Error())
	}
	defer resultCursor.Close(m.ctx)
	for resultCursor.Next(m.ctx) {
		var result pool.Node
		err := resultCursor.Decode(&result)
		if err != nil {
			return nil, errors.New("[Mongo/GetAll] decode node data" + err.Error())
		}
		nodeList = append(nodeList, result)
	}
	if err := resultCursor.Err(); err != nil {
		return nil, errors.New("[Mongo/GetAll] iterate result" + err.Error())
	}
	return nodeList, nil
}

func (m *Storage) Remove(node pool.Node) error {
	rowsAffected, err := m.collection.DeleteOne(m.ctx, bson.M{"key": node.Key})
	if err != nil {
		return errors.New("[Mongo/Remove] delete from node collection, " + err.Error())
	}
	if rowsAffected.DeletedCount == 0 {
		return errors.New("[Mongo/Remove] delete from node collection: affected 0 rows")
	}
	return nil
}

func (m *Storage) UpdateAddress(node pool.Node, newAddress string) error {
	filter := bson.M{"key": node.Key}
	update := bson.M{"$set": bson.M{"address": newAddress}}
	result, err := m.collection.UpdateOne(m.ctx, filter, update)
	if err != nil {
		err = errors.New("[Mongo/UpdateAddress], " + err.Error())
		return err
	}
	if result.ModifiedCount == 0 {
		return storage.ErrNotFound
	}
	return nil
}

func (m *Storage) getByField(key, value string) (pool.Node, error) {
	node := pool.Node{}
	filter := bson.M{
		key: bson.M{"$eq": value},
	}
	err := m.collection.FindOne(m.ctx, filter).Decode(&node)
	return node, err
}
