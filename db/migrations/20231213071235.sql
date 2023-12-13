-- Create "email_logs" table
CREATE TABLE "public"."email_logs" (
  "LogId" bigserial NOT NULL,
  "emailTo" text NOT NULL DEFAULT '',
  "name" text NOT NULL DEFAULT '',
  "subject" text NOT NULL DEFAULT '',
  "body" text NOT NULL DEFAULT '',
  "status" text NOT NULL DEFAULT '',
  PRIMARY KEY ("LogId")
);
