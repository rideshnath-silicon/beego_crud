package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type AddCarTable_20231206_181332 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &AddCarTable_20231206_181332{}
	m.Created = "20231206_181332"

	migration.Register("AddCarTable_20231206_181332", m)
}

// Run the migrations
func (m *AddCarTable_20231206_181332) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE IF NOT EXISTS "car" (
        "id" bigserial NOT NULL PRIMARY KEY,
        "car_name" text NOT NULL DEFAULT '' ,
        "car_image" text,
        "modified_by" text NOT NULL DEFAULT '' ,
        "model" text NOT NULL DEFAULT '' ,
        "car_type" text NOT NULL DEFAULT '' ,
        "ctreated_date" timestamp with time zone,
        "updated_date" timestamp with time zone
    );`)

}

// Reverse the migrations
func (m *AddCarTable_20231206_181332) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
