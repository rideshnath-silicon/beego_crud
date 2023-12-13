-- Create "car" table
CREATE TABLE "public"."car" (
  "id" bigserial NOT NULL,
  "car_name" text NOT NULL DEFAULT '',
  "car_image" text NULL,
  "modified_by" text NOT NULL DEFAULT '',
  "model" text NOT NULL DEFAULT '',
  "car_type" text NOT NULL DEFAULT '',
  "ctreated_date" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create "home_setting" table
CREATE TABLE "public"."home_setting" (
  "id" bigserial NOT NULL,
  "section" character varying(255) NOT NULL DEFAULT '',
  "type" character varying(255) NOT NULL DEFAULT '',
  "key" character varying(255) NOT NULL DEFAULT '',
  "value" character varying(255) NOT NULL DEFAULT '',
  "demo" text NOT NULL,
  "created_at" timestamptz NULL,
  "update_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create "users" table
CREATE TABLE "public"."users" (
  "id" bigserial NOT NULL,
  "first_name" text NULL,
  "last_name" text NULL,
  "email" text NOT NULL DEFAULT '',
  "phone_number" text NULL,
  "country" integer NOT NULL DEFAULT 0,
  "role" text NOT NULL DEFAULT '',
  "age" integer NOT NULL DEFAULT 0,
  "password" text NOT NULL DEFAULT '',
  "otp" text NULL,
  "verified" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "users_email_key" to table: "users"
CREATE UNIQUE INDEX "users_email_key" ON "public"."users" ("email");
