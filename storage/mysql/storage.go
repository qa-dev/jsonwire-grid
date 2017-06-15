package mysql

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/qa-dev/jsonwire-grid/pool"
	"sort"
	"strconv"
	"strings"
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

type MysqlStorage struct {
	db *sqlx.DB
}

func NewMysqlStorage(db *sqlx.DB) *MysqlStorage {
	return &MysqlStorage{db}
}

func (s *MysqlStorage) Add(node pool.Node) error {
	tx, err := s.db.Beginx()
	if err != nil {
		err = errors.New("[MysqlStorage/Add] Can't start transaction: " + err.Error())
		return err
	}
	_, err = tx.NamedExec(
		"INSERT INTO node (type, address, status, sessionId, updated, registred) "+
			"VALUES (:type, :address, :status, :sessionId, :updated, :registred) "+
			"ON DUPLICATE KEY UPDATE "+
			"type = :type, status = :status, sessionId = :sessionId, updated = :updated, registred = :registred",
		map[string]interface{}{
			"type":      string(node.Type),
			"address":   node.Address,
			"sessionId": node.SessionID,
			"status":    string(node.Status),
			"updated":   node.Updated,
			"registred": node.Registered,
		},
	)
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
	for i, capabilities := range node.CapabilitiesList {
		for name, value := range capabilities {
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

func (s *MysqlStorage) ReserveAvailable(capabilities pool.Capabilities) (node pool.Node, err error) {
	tx, err := s.db.Beginx()
	if err != nil {
		err = errors.New("[MysqlStorage/ReserveAvailable] Can't start transaction: " + err.Error())
		return
	}
	nodeModel := new(MysqlNodeModel)
	where := "n.status = ?"

	possibleCapabilities, err := s.filterPossibleCapabilities(capabilities)
	if err != nil {
		tx.Rollback()
		err = errors.New("[MysqlStorage/ReserveAvailable] " + err.Error())
		return
	}
	args := []interface{}{string(pool.NodeStatusAvailable)}

	switch {
	case len(possibleCapabilities) > 0:
		capsConditionList := []string{}
		var condition string
		for name, value := range possibleCapabilities {
			var castedValue string
			switch value := value.(type) {
			case int:
				castedValue = strconv.Itoa(value)
			case float32:
				castedValue = strconv.FormatFloat(float64(value), 'f', -1, 64)
			case float64:
				castedValue = strconv.FormatFloat(value, 'f', -1, 64)
			case bool:
				castedValue = strconv.FormatBool(value)
			case string:
				castedValue = string(value)
			default:
				tx.Rollback()
				err = errors.New(fmt.Sprintf("[MysqlStorage/ReserveAvailable] Invalid capability vaslue type: %T, %v", value, value))
				return
			}

			condition = "c.name = '" + name + "' AND c.value = '" + castedValue + "'"
			capsConditionList = append(capsConditionList, condition)
		}
		where += " AND (" + strings.Join(capsConditionList, " OR ") + ")"

		countCapabilities := strconv.Itoa(len(capsConditionList))
		err = tx.QueryRowx(
			`
				SELECT
					n.id,
					n.type,
					n.status,
					n.address,
					n.sessionId,
					n.updated,
					n.registred
				FROM node n
				LEFT JOIN capabilities c ON n.address = c.nodeAddress AND `+where+`
			GROUP BY c.setId
			HAVING count(c.setId) =  `+countCapabilities+`
			ORDER BY n.updated ASC
			LIMIT 1
			FOR UPDATE
		`,
			args...).
			StructScan(nodeModel)
	default:
		err = tx.QueryRowx(
			`SELECT n.* FROM node n WHERE `+where+` ORDER BY n.updated ASC LIMIT 1 FOR UPDATE`,
			args...).
			StructScan(nodeModel)
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
	defer rows.Close()
	if err != nil {
		err = errors.New("[MysqlStorage/GetAll] Select error: " + err.Error())
		return nil, err
	}
	for rows.Next() {
		node := new(MysqlNodeModel)
		err := rows.StructScan(node)
		if err != nil {
			err = errors.New("[MysqlStorage/GetAll] Error on iteration: " + err.Error())
			return nil, err
		}
		nodeList = append(nodeList, *mapper(node))
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
		[]pool.Capabilities{}, //todo: заглушка для capabilities, так как пока мапить это не надо
	)
	return node
}

func (s *MysqlStorage) filterPossibleCapabilities(capabilities pool.Capabilities) (pool.Capabilities, error) {
	var possibleCapNameList []string
	err := s.db.Select(&possibleCapNameList, "SELECT DISTINCT name FROM capabilities")
	if err != nil {
		err = errors.New("[MysqlStorage/possibleCapabilityNameList] " + err.Error())
		return nil, err
	}

	sort.Strings(possibleCapNameList)
	for name, _ := range capabilities {
		if i := sort.SearchStrings(possibleCapNameList, name); !(i < len(possibleCapNameList) && possibleCapNameList[i] == name) {
			delete(capabilities, name)
		}
	}
	return capabilities, nil
}
