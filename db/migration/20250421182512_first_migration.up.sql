CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "phone" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "financials" (
  "id" bigserial PRIMARY KEY,
  "user_id" varchar NOT NULL,
  "amount" bigserial NOT NULL,
  "direction" varchar NOT NULL,
  "type_id" bigserial NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "financial_types" (
  "id" bigserial PRIMARY KEY,
  "type" varchar UNIQUE NOT NULL
);

CREATE INDEX ON "financial_types" ("type");

ALTER TABLE "financials" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("username");

ALTER TABLE "financials" ADD FOREIGN KEY ("type_id") REFERENCES "financial_types" ("id");
