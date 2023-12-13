package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type AddHomeSettingTable_20231206_182230 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &AddHomeSettingTable_20231206_182230{}
	m.Created = "20231206_182230"

	migration.Register("AddHomeSettingTable_20231206_182230", m)
}

// Run the migrations
func (m *AddHomeSettingTable_20231206_182230) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`	CREATE TABLE IF NOT EXISTS "home_setting" (
        "id" bigserial NOT NULL PRIMARY KEY,
        "section" varchar(255) NOT NULL DEFAULT '' ,
        "type" varchar(255) NOT NULL DEFAULT '' ,
        "key" varchar(255) NOT NULL DEFAULT '' ,
        "value" varchar(255) NOT NULL DEFAULT '' ,
        "demo" text NOT NULL,
        "created_date" timestamp with time zone,
        "update_date" timestamp with time zone
    );`)

}

// Reverse the migrations
func (m *AddHomeSettingTable_20231206_182230) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
