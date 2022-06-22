package db

import (
	"database/sql"
	"fmt"
	"github.com/galenliu/gateway/pkg/rules_engine"
	json "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
)

func (s *Storage) CreateRule(description rules_engine.Rule) (int64, error) {
	bytes, err := json.Marshal(description)
	if err != nil {
		return 0, err
	}
	stmt, err := s.db.Prepare("INSERT INTO rules (description) VALUES (?)")
	if err != nil {
		return 0, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Error("stmt close err: %s", err.Error())
		}
	}(stmt)
	res, ee := stmt.Exec(bytes)
	if ee != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Storage) UpdateRule(ruleId int64, r rules_engine.Rule) error {
	bytes, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("thing marshal error: %s", err.Error())
	}
	_, err = s.db.Exec(`update rules set id=@id where description=@description`, sql.Named("id", ruleId), sql.Named("description", bytes))
	return err
}

func (s *Storage) DeleteRule(ruleId int64) error {
	stmt, err := s.db.Prepare(`delete from rules where id = ?`)
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Error("stmt close err: %s", err.Error())
		}
	}(stmt)
	_, err = stmt.Exec(ruleId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetRules() map[int64]rules_engine.RuleDescription {
	rules := make(map[int64]rules_engine.RuleDescription)
	rows, err := s.db.Query("SELECT id, description FROM rules")
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Printf(err.Error())
		}
	}(rows)
	if err != nil {
		return nil
	}
	for rows.Next() {
		var id int64
		var description string
		err = rows.Scan(&id, &description)
		if err != nil {
			continue
		}
		var rule rules_engine.RuleDescription
		err := json.Unmarshal([]byte(description), &rule)
		if err != nil {
			continue
		}
		rules[id] = rule
	}
	return rules
}
