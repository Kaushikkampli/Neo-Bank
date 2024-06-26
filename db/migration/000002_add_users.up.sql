CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "full_name" varchar NOT NULL,
  "hashed_passwd" varchar NOT NULL,
  "email_id" varchar NOT NULL,
  "passwd_updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

--CREATE UNIQUE INDEX ON "accounts" ("owner", "currency");
ALTER TABLE accounts ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency")