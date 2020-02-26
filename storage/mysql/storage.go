package mysql

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/qa-dev/jsonwire-grid/pool"
	"github.com/qa-dev/jsonwire-grid/pool/capabilities"
	"github.com/qa-dev/jsonwire-grid/storage"
)

type MysqlNodeModel struct {
	ID         string `db:"id"`
	Key        string `db:"key"`
	Type       string `db:"type"`
	Address    string `db:"address"`
	Status     string `db:"status"`
	SessionID  string `db:"sessionId"`
	Updated    int64  `db:"updated"`
	Registered int64  `db:"registred"`
}

type MysqlCapabilitiesModel struct {
	ID      int    `db:"id"`
	NodeKey string `db:"nodeKey"`
	SetID   string `db:"setId"`
	Name    string `db:"name"`
	Value   string `db:"value"`
}

type MysqlStorage struct {
	db *sqlx.DB
}

func NewMysqlStorage(db *sqlx.DB) *MysqlStorage {
	return &MysqlStorage{db}
}

func (s *MysqlStorage) Add(node pool.Node, limit int) error {
	tx, err := s.db.Beginx()
	if err != nil {
		err = errors.New("[MysqlStorage/Add] start transaction: " + err.Error())
		return err
	}
	// black magic, but it works
	result, err := tx.NamedExec(
		"INSERT INTO node (`key`, type, address, status, sessionId, updated, registred) "+
			"SELECT :key, :type, :address, :status, :sessionId, :updated, :registred "+
			"FROM DUAL "+
			"WHERE 0 = :limit OR EXISTS (SELECT TRUE FROM node WHERE type = :type HAVING count(*) < :limit)"+
			"ON DUPLICATE KEY UPDATE "+
			"type = :type, address = :address, status = :status, sessionId = :sessionId, updated = :updated, registred = :registred",
		map[string]interface{}{
			"key":       node.Key,
			"type":      string(node.Type),
			"address":   node.Address,
			"sessionId": node.SessionID,
			"status":    string(node.Status),
			"updated":   node.Updated,
			"registred": node.Registered,
			"limit":     limit,
		},
	)
	if err != nil {
		_ = tx.Rollback()
		return errors.New("[MysqlStorage/Add] insert entry in `node` table, " + err.Error())
	}

	countAffectedRows, err := result.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return errors.New("[MysqlStorage/Add] get affected rows, " + err.Error())
	}

	if countAffectedRows == 0 {
		_ = tx.Rollback()
		return errors.New("[MysqlStorage/Add] No rows was affected (may be limit reached)")
	}

	_, err = tx.Exec("DELETE FROM capabilities WHERE nodeKey = ?", node.Key)
	if err != nil {
		_ = tx.Rollback()
		return errors.New("[MysqlStorage/Add] delete old capabilities: " + err.Error())
	}
	var preparedCapabilities []map[string]interface{}
	var preparedCapability map[string]interface{}
	for i, caps := range node.CapabilitiesList {
		for name, value := range caps {
			preparedCapability = map[string]interface{}{
				"nodeKey": node.Key,
				"setId":   node.Key + "|" + strconv.Itoa(i), // просто уникальное значение для сета
				"name":    name,
				"value":   value,
			}
			preparedCapabilities = append(preparedCapabilities, preparedCapability)
		}
	}

	for _, preparedCapability := range preparedCapabilities {
		_, err = tx.NamedExec(
			"INSERT INTO capabilities (nodeKey, setId, name, value) "+
				"VALUES (:nodeKey, :setId, :name, :value)",
			preparedCapability,
		)
		if err != nil {
			_ = tx.Rollback()
			return errors.New("[MysqlStorage/Add] insert new capabilities: " + err.Error())

		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.New("[MysqlStorage/Add] commit transaction: " + err.Error())
	}
	return nil
}

func (s *MysqlStorage) ReserveAvailable(nodeList []pool.Node) (pool.Node, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return pool.Node{}, errors.New("[MysqlStorage/ReserveAvailable] start transaction: " + err.Error())
	}

	nodeKeyList := make([]string, 0, len(nodeList))
	for _, node := range nodeList {
		nodeKeyList = append(nodeKeyList, node.Key)
	}
	args := []interface{}{string(pool.NodeStatusAvailable), nodeKeyList}

	//var row *sqlx.Row
	query, args, err := sqlx.In(
		"SELECT n.* FROM node n WHERE n.status = ? AND n.`key` IN (?) ORDER BY n.updated ASC LIMIT 1 FOR UPDATE;", args...)
	if err != nil {
		return pool.Node{}, errors.New("[MysqlStorage/ReserveAvailable] make select query: " + err.Error())
	}
	row := tx.QueryRowx(query, args...)

	nodeModel := new(MysqlNodeModel)
	err = row.StructScan(nodeModel)

	if err == sql.ErrNoRows {
		_ = tx.Rollback()
		return pool.Node{}, storage.ErrNotFound
	}

	if err != nil {
		_ = tx.Rollback()
		return pool.Node{}, errors.New("[MysqlStorage/ReserveAvailable] select node: " + err.Error())
	}

	nodeModel.Updated = time.Now().Unix()
	nodeModel.Status = string(pool.NodeStatusReserved)
	res, err := tx.Exec(
		"UPDATE node SET status = ?, updated = ? WHERE `key` = ?",
		nodeModel.Status,
		nodeModel.Updated,
		nodeModel.Key,
	)
	if err != nil {
		_ = tx.Rollback()
		return pool.Node{}, errors.New("[MysqlStorage/ReserveAvailable] update table `node`: " + err.Error())
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return pool.Node{}, errors.New("[MysqlStorage/ReserveAvailable] get affected rows: " + err.Error())
	}
	if rowsAffected == 0 {
		_ = tx.Rollback()
		return pool.Node{}, errors.New("[MysqlStorage/ReserveAvailable] affected 0 rows")
	}
	err = tx.Commit()
	if err != nil {
		return pool.Node{}, errors.New("[MysqlStorage/ReserveAvailable] commit transaction: " + err.Error())
	}

	return *mapper(nodeModel), nil
}

func (s *MysqlStorage) SetBusy(node pool.Node, sessionID string) error {
	res, err := s.db.Exec(
		"UPDATE node SET sessionID = ?, updated = ?, status = ? WHERE `key` = ?",
		sessionID,
		time.Now().Unix(),
		string(pool.NodeStatusBusy),
		node.Key,
	)
	if err != nil {
		return errors.New("[MysqlStorage/SetBusy] update table `node`, " + err.Error())
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.New("[MysqlStorage/SetBusy] get affected rows, " + err.Error())
	}
	if rowsAffected == 0 {
		return errors.New("[MysqlStorage/SetBusy] affected 0 rows")
	}
	return nil
}

func (s *MysqlStorage) SetAvailable(node pool.Node) error {
	res, err := s.db.Exec(
		"UPDATE node SET status = ?, updated = ?  WHERE `key` = ?",
		string(pool.NodeStatusAvailable),
		time.Now().Unix(),
		node.Key,
	)
	if err != nil {
		err = errors.New("[MysqlStorage/SetAvailable] update table `node`, " + err.Error())
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		err = errors.New("[MysqlStorage/SetAvailable] get affected rows, " + err.Error())
		return err
	}
	if rowsAffected == 0 {
		return errors.New("[MysqlStorage/SetAvailable] affected 0 rows")
	}
	return nil
}

func (s *MysqlStorage) GetCountWithStatus(status *pool.NodeStatus) (int, error) {
	var count int
	var err error
	if status != nil {
		err = s.db.
			QueryRowx("SELECT COUNT(*) FROM node WHERE status = ?", string(*status)).
			Scan(&count)
	} else {
		err = s.db.
			QueryRowx("SELECT COUNT(*) FROM node").
			Scan(&count)
	}
	if err != nil {
		return 0, errors.New("[MysqlStorage/GetCountWithStatus] select from table `node`, " + err.Error())
	}
	return count, nil
}

func (s *MysqlStorage) GetBySession(sessionID string) (pool.Node, error) {
	node := new(MysqlNodeModel)
	err := s.db.QueryRowx("SELECT * FROM node WHERE sessionID = ?", sessionID).StructScan(node)
	if err != nil {
		return pool.Node{}, errors.New("[MysqlStorage/GetBySession] select from table `node`," + err.Error())
	}
	return *mapper(node), nil
}

func (s *MysqlStorage) GetByAddress(address string) (pool.Node, error) {
	node := new(MysqlNodeModel)
	queryString := "SELECT * FROM node WHERE address = ?"
	err := s.db.QueryRowx(queryString, address).StructScan(node)
	if err != nil {
		return pool.Node{}, errors.New("[MysqlStorage/GetByAddress] select from table `node`," + err.Error())
	}
	return *mapper(node), nil
}

func (s *MysqlStorage) GetAll() ([]pool.Node, error) {
	nodeList := make([]pool.Node, 0)
	rows, err := s.db.Queryx("SELECT * FROM node")
	if err != nil {
		return nil, errors.New("[MysqlStorage/GetAll] select from table `node`," + err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		node := new(MysqlNodeModel)
		err := rows.StructScan(node)
		if err != nil {
			return nil, errors.New("[MysqlStorage/GetAll] iterate result rows" + err.Error())
		}
		nodeList = append(nodeList, *mapper(node))
	}

	var rowsCaps []MysqlCapabilitiesModel
	capsMap := map[string]map[string]capabilities.Capabilities{}
	err = s.db.Select(&rowsCaps, "SELECT * FROM capabilities")
	if err != nil {
		return nil, errors.New("[MysqlStorage/GetAll] get all capabilities from db, " + err.Error())
	}
	for _, row := range rowsCaps {
		_, ok := capsMap[row.NodeKey]
		if !ok {
			capsMap[row.NodeKey] = map[string]capabilities.Capabilities{}
		}
		currCaps, ok := capsMap[row.NodeKey][row.SetID]
		if !ok {
			currCaps = capabilities.Capabilities{}
		}
		currCaps[row.Name] = row.Value
		capsMap[row.NodeKey][row.SetID] = currCaps
	}

	for i, node := range nodeList {
		capsList := make([]capabilities.Capabilities, 0, len(capsMap[node.Key]))
		for _, currCaps := range capsMap[node.Key] {
			capsList = append(capsList, currCaps)
		}
		nodeList[i].CapabilitiesList = capsList
	}

	return nodeList, nil
}

func (s *MysqlStorage) Remove(node pool.Node) error {
	res, err := s.db.Exec("DELETE FROM node WHERE `key` = ?", node.Key)
	if err != nil {
		return errors.New("[MysqlStorage/Remove] delete from table `node`, " + err.Error())
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.New("[MysqlStorage/Remove] delete from node: get affected rows," + err.Error())
	}
	if rowsAffected == 0 {
		return errors.New("[MysqlStorage/Remove] delete from node: affected 0 rows")
	}
	res, err = s.db.Exec("DELETE FROM capabilities WHERE `nodeKey` = ?", node.Key)
	if err != nil {
		return errors.New("[MysqlStorage/Remove] delete from table `capabilities`: " + err.Error())
	}
	rowsAffected, err = res.RowsAffected()
	if err != nil {
		return errors.New("[MysqlStorage/Remove] delete from capabilities: " + err.Error())
	}
	if rowsAffected == 0 {
		return errors.New("[MysqlStorage/Remove] delete from capabilities: affected 0 rows")
	}
	return nil
}

func (s *MysqlStorage) UpdateAddress(node pool.Node, newAddress string) error {
	res, err := s.db.Exec(
		"UPDATE node SET address = ? WHERE `key` = ?",
		newAddress,
		node.Key,
	)
	if err != nil {
		err = errors.New("[MysqlStorage/UpdateAddress] update table `node`, " + err.Error())
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		err = errors.New("[MysqlStorage/UpdateAddress] get affected rows, " + err.Error())
		return err
	}
	if rowsAffected == 0 {
		return storage.ErrNotFound
	}
	return nil
}

func mapper(model *MysqlNodeModel) *pool.Node {
	node := pool.NewNode(
		model.Key,
		pool.NodeType(model.Type),
		model.Address,
		pool.NodeStatus(model.Status),
		model.SessionID,
		model.Updated,
		model.Registered,
		[]capabilities.Capabilities{}, //todo: stub, temporary
	)
	return node
}
