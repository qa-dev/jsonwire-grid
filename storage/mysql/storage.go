package mysql

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/qa-dev/jsonwire-grid/pool"
	"github.com/qa-dev/jsonwire-grid/pool/capabilities"
	"github.com/qa-dev/jsonwire-grid/storage"
	"strconv"
	"time"
)

type MysqlNodeModel struct {
	ID         string `db:"id"`
	Type       string `db:"type"`
	Address    string `db:"address"`
	Status     string `db:"status"`
	SessionID  string `db:"sessionId"`
	Updated    int64  `db:"updated"`
	Registered int64  `db:"registred"`
}

type MysqlCapabilitiesModel struct {
	ID          int    `db:"id"`
	NodeAddress string `db:"nodeAddress"`
	SetID       string `db:"setId"`
	Name        string `db:"name"`
	Value       string `db:"value"`
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
		err = errors.New("[MysqlStorage/Add] Can't start transaction: " + err.Error())
		return err
	}
	//todo: black magic, but it works
	result, err := tx.NamedExec(
		"INSERT INTO node (type, address, status, sessionId, updated, registred) "+
			"SELECT :type, :address, :status, :sessionId, :updated, :registred "+
			"FROM DUAL "+
			"WHERE 0 = :limit OR EXISTS (SELECT TRUE FROM node WHERE type = :type HAVING count(*) < :limit)"+
			"ON DUPLICATE KEY UPDATE "+
			"type = :type, status = :status, sessionId = :sessionId, updated = :updated, registred = :registred",
		map[string]interface{}{
			"type":      string(node.Type),
			"address":   node.Address,
			"sessionId": node.SessionID,
			"status":    string(node.Status),
			"updated":   node.Updated,
			"registred": node.Registered,
			"limit":     limit,
		},
	)

	countAffectedRows, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		err = errors.New("[MysqlStorage/Add] Can't get affected rows, " + err.Error())
		return err
	}

	if countAffectedRows == 0 {
		tx.Rollback()
		err = errors.New("[MysqlStorage/Add] No rows was affected (may be limit reached)")
		return err
	}

	if err != nil {
		tx.Rollback()
		err = errors.New("[MysqlStorage/Add] Can't insert new node: " + err.Error())
		return err
	}
	_, err = tx.Exec("DELETE FROM capabilities WHERE nodeAddress = ?", node.Address)
	if err != nil {
		tx.Rollback()
		err = errors.New("[MysqlStorage/Add] Can't delete old capabilities: " + err.Error())
		return err
	}
	var preparedCapabilities []map[string]interface{}
	var preparedCapability map[string]interface{}
	for i, caps := range node.CapabilitiesList {
		for name, value := range caps {
			preparedCapability = map[string]interface{}{
				"nodeAddress": node.Address,
				"setId":       node.Address + "|" + strconv.Itoa(i), // просто уникальное значение для сета
				"name":        name,
				"value":       value,
			}
			preparedCapabilities = append(preparedCapabilities, preparedCapability)
		}
	}

	for _, preparedCapability := range preparedCapabilities {
		_, err = tx.NamedExec(
			"INSERT INTO capabilities (nodeAddress, setId, name, value) "+
				"VALUES (:nodeAddress, :setId, :name, :value)",
			preparedCapability,
		)
		if err != nil {
			tx.Rollback()
			err = errors.New("[MysqlStorage/Add] Can't insert new capabilities: " + err.Error())
			return err

		}
	}

	err = tx.Commit()
	if err != nil {
		err = errors.New("[MysqlStorage/Add] Can't commit transaction: " + err.Error())
		return err
	}
	return err
}

func (s *MysqlStorage) ReserveAvailable(nodeList []pool.Node) (node pool.Node, err error) {
	tx, err := s.db.Beginx()
	if err != nil {
		err = errors.New("[MysqlStorage/ReserveAvailable] Can't start transaction: " + err.Error())
		return
	}

	nodeAddressList := make([]string, 0, len(nodeList))
	for _, node := range nodeList {
		nodeAddressList = append(nodeAddressList, node.Address)
	}
	args := []interface{}{string(pool.NodeStatusAvailable), nodeAddressList}

	//var row *sqlx.Row
	query, args, err := sqlx.In(
		"SELECT n.* FROM node n WHERE n.status = ? AND n.address IN (?) ORDER BY n.updated ASC LIMIT 1 FOR UPDATE;", args...)
	if err != nil {
		err = errors.New("[MysqlStorage/ReserveAvailable] Can't make select query: " + err.Error())
		return
	}
	row := tx.QueryRowx(query, args...)

	nodeModel := new(MysqlNodeModel)
	err = row.StructScan(nodeModel)

	if err == sql.ErrNoRows {
		tx.Rollback()
		err = storage.ErrNotFound
		return
	}

	if err != nil {
		tx.Rollback()
		err = errors.New("[MysqlStorage/ReserveAvailable] Can't select node: " + err.Error())
		return
	}

	nodeModel.Updated = time.Now().Unix()
	nodeModel.Status = string(pool.NodeStatusReserved)
	res, err := tx.Exec(
		"UPDATE node SET status = ?, updated = ? WHERE address = ?",
		nodeModel.Status,
		nodeModel.Updated,
		nodeModel.Address,
	)
	if err != nil {
		tx.Rollback()
		err = errors.New("[MysqlStorage/ReserveAvailable] Can't update status: " + err.Error())
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		err = errors.New("[MysqlStorage/ReserveAvailable] Can't update status: " + err.Error())
		return
	}
	if rowsAffected == 0 {
		tx.Rollback()
		err = errors.New("[MysqlStorage/ReserveAvailable] Can't update status: affected 0 rows")
		return
	}
	err = tx.Commit()
	if err != nil {
		err = errors.New("[MysqlStorage/ReserveAvailable] Can't commit transaction: " + err.Error())
		return
	}
	node = *mapper(nodeModel)
	return
}

func (s *MysqlStorage) SetBusy(node pool.Node, sessionId string) error {
	res, err := s.db.Exec(
		"UPDATE node SET sessionId = ?, updated = ?, status = ? WHERE address = ?",
		sessionId,
		time.Now().Unix(),
		string(pool.NodeStatusBusy),
		node.Address,
	)
	if err != nil {
		err = errors.New("[MysqlStorage/SetBusy] Can't update status: " + err.Error())
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		err = errors.New("[MysqlStorage/SetBusy] Can't update status: " + err.Error())
		return err
	}
	if rowsAffected == 0 {
		err = errors.New("[MysqlStorage/SetBusy] Can't update status: affected 0 rows")
	}
	return err
}

func (s *MysqlStorage) SetAvailable(node pool.Node) error {
	res, err := s.db.Exec(
		"UPDATE node SET status = ?, updated = ?  WHERE address = ?",
		string(pool.NodeStatusAvailable),
		time.Now().Unix(),
		node.Address,
	)
	if err != nil {
		err = errors.New("[MysqlStorage/SetAvailable] Can't update status: " + err.Error())
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		err = errors.New("[MysqlStorage/SetAvailable] Can't update status: " + err.Error())
		return err
	}
	if rowsAffected == 0 {
		err = errors.New("[MysqlStorage/SetAvailable] Can't update status: affected 0 rows")
	}
	return err
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
		err = errors.New("[MysqlStorage/GetCountWithStatus] " + err.Error())
	}
	return count, err
}

func (s *MysqlStorage) GetBySession(sessionId string) (pool.Node, error) {
	node := new(MysqlNodeModel)
	err := s.db.QueryRowx("SELECT * FROM node WHERE sessionId = ?", sessionId).StructScan(node)
	if err != nil {
		err = errors.New("[MysqlStorage/GetBySession] " + err.Error())
	}
	return *mapper(node), err
}

func (s *MysqlStorage) GetByAddress(address string) (pool.Node, error) {
	node := new(MysqlNodeModel)
	queryString := "SELECT * FROM node WHERE address = ?"
	err := s.db.QueryRowx(queryString, address).StructScan(node)
	if err != nil {
		err = errors.New("[MysqlStorage/GetByAddress] " + err.Error())
	}
	return *mapper(node), err
}

func (s *MysqlStorage) GetAll() ([]pool.Node, error) {
	nodeList := make([]pool.Node, 0)
	rows, err := s.db.Queryx("SELECT * FROM node")
	if err != nil {
		err = errors.New("[MysqlStorage/GetAll] Select error: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		node := new(MysqlNodeModel)
		err := rows.StructScan(node)
		if err != nil {
			err = errors.New("[MysqlStorage/GetAll] Error on iteration: " + err.Error())
			return nil, err
		}
		nodeList = append(nodeList, *mapper(node))
	}

	var rowsCaps []MysqlCapabilitiesModel
	capsMap := map[string]map[string]capabilities.Capabilities{}
	err = s.db.Select(&rowsCaps, "SELECT * FROM capabilities")
	if err != nil {
		err = errors.New("[MysqlStorage/GetAll] Can't get all capabilities from db, " + err.Error())
		return nil, err
	}
	for _, row := range rowsCaps {
		_, ok := capsMap[row.NodeAddress]
		if !ok {
			capsMap[row.NodeAddress] = map[string]capabilities.Capabilities{}
		}
		currCaps, ok := capsMap[row.NodeAddress][row.SetID]
		if !ok {
			currCaps = capabilities.Capabilities{}
		}
		currCaps[row.Name] = row.Value
		capsMap[row.NodeAddress][row.SetID] = currCaps
	}

	for i, node := range nodeList {
		capsList := make([]capabilities.Capabilities, 0, len(capsMap[node.Address]))
		for _, currCaps := range capsMap[node.Address] {
			capsList = append(capsList, currCaps)
		}
		nodeList[i].CapabilitiesList = capsList
	}

	return nodeList, err
}

func (s *MysqlStorage) Remove(node pool.Node) error {
	res, err := s.db.Exec("DELETE FROM node WHERE address = ?", node.Address)
	if err != nil {
		err = errors.New("[MysqlStorage/Remove] Can't delete from node: " + err.Error())
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		err = errors.New("[MysqlStorage/Remove] Can't delete from node: " + err.Error())
		return err
	}
	if rowsAffected == 0 {
		err = errors.New("[MysqlStorage/Remove] Can't delete from node: affected 0 rows")
		return err
	}
	res, err = s.db.Exec("DELETE FROM capabilities WHERE nodeAddress = ?", node.Address)
	if err != nil {
		err = errors.New("[MysqlStorage/Remove] Can't delete from capabilities: " + err.Error())
		return err
	}
	rowsAffected, err = res.RowsAffected()
	if err != nil {
		err = errors.New("[MysqlStorage/Remove] Can't delete from capabilities: " + err.Error())
		return err
	}
	if rowsAffected == 0 {
		err = errors.New("[MysqlStorage/Remove] Can't delete from capabilities: affected 0 rows")
		return err
	}
	return err
}

func mapper(model *MysqlNodeModel) *pool.Node {
	node := pool.NewNode(
		pool.NodeType(model.Type),
		model.Address,
		pool.NodeStatus(model.Status),
		model.SessionID,
		model.Updated,
		model.Registered,
		[]capabilities.Capabilities{}, //todo: заглушка для capabilities, так как пока мапить это не надо
	)
	return node
}
