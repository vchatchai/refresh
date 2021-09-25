package model

import (
	"database/sql"
	"encoding/json"
)

/**
 */
type MyNullString struct {
	sql.NullString
}

func (ns MyNullString) MarshalJSON() ([]byte, error) {

	if "NULL" == ns.String {
		// fmt.Println("MyNullString NULL")
		ns.Valid = false
	} else {
		// fmt.Println("MyNullString NULL")
		if ns.Valid {
			return json.Marshal(ns.String)
		}
	}

	return []byte(`null`), nil
}

func (ns *MyNullString) UnmarshalJSON(b []byte) error {

	err := json.Unmarshal(b, &ns.String)

	if err != nil {
		return err
	} else {
		ns.Valid = true
	}

	return nil

}

func (ns *MyNullString) Value() interface{} {

	if ns.Valid {
		return ns.String
	} else {
		return ns.NullString
	}
}

type MyNullTime struct {
	sql.NullTime
}

func (n MyNullTime) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Time)
	}
	return []byte(`null`), nil

}

func (ns *MyNullTime) UnmarshalJSON(b []byte) error {

	err := json.Unmarshal(b, &ns.Time)

	if err != nil {
		return err
	} else {
		ns.Valid = true
	}

	return nil

}

func (ns *MyNullTime) Value() interface{} {

	if ns.Valid {
		return ns.Time
	} else {
		return ns.NullTime
	}

}
