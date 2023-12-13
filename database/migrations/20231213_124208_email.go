package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type Email_20231213_124208 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Email_20231213_124208{}
	m.Created = "20231213_124208"

	migration.Register("Email_20231213_124208", m)
}

// Run the migrations
func (m *Email_20231213_124208) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE "public"."email_logs" (
		"LogId" bigserial NOT NULL,
		"emailTo" text NOT NULL DEFAULT '',
		"name" text NOT NULL DEFAULT '',
		"subject" text NOT NULL DEFAULT '',
		"body" text NOT NULL DEFAULT '',
		"status" text NOT NULL DEFAULT '',
		PRIMARY KEY ("LogId")
	);`)

}

// Reverse the migrations
func (m *Email_20231213_124208) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
