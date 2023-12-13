package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type Procedure_20231213_102443 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Procedure_20231213_102443{}
	m.Created = "20231213_102443"

	migration.Register("Procedure_20231213_102443", m)
}

// Run the migrations
func (m *Procedure_20231213_102443) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`
	CREATE PROCEDURE "lastinsertuser" () LANGUAGE plpgsql AS '
	DECLARE
		last_user users%ROWTYPE;  -- Declaring a variable to hold the last user
	BEGIN
		SELECT * INTO last_user FROM users ORDER BY id DESC LIMIT 1;
		RAISE NOTICE ''Last User: %'', last_user;
		COMMIT;
	END;
	';`)

}

// Reverse the migrations
func (m *Procedure_20231213_102443) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
