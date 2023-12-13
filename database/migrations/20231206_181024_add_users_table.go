package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type AddUsersTable_20231206_181024 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &AddUsersTable_20231206_181024{}
	m.Created = "20231206_181024"

	migration.Register("AddUsersTable_20231206_181024", m)
}

// Run the migrations
func (m *AddUsersTable_20231206_181024) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE IF NOT EXISTS "users" (
		"id" bigserial NOT NULL PRIMARY KEY,
		"first_name" text,
		"last_name" text,
		"email" text NOT NULL DEFAULT ''  UNIQUE,
		"phone_number" text,
		"country" integer DEFAULT 0 ,
		"role" text NOT NULL DEFAULT '' ,
		"age" integer NOT NULL DEFAULT 0 ,
		"password" text NOT NULL DEFAULT '' ,
		"otp" text,
		"verified" text,
		"created_at" timestamp with time zone,
		"updated_at" timestamp with time zone,
		"deleted_at" timestamp with time zone
	);`)

}

// Reverse the migrations
func (m *AddUsersTable_20231206_181024) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
