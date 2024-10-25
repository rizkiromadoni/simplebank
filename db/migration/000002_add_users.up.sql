CREATE TABLE users (
    "username" varchar PRIMARY KEY,
    "email" varchar UNIQUE NOT NULL,
    "full_name" varchar NOT NULL,
    "password" varchar NOT NULL,
    "password_changed_at" timestamptz NOT NULL DEFAULT '001-01-01 00:00:00Z',
    "created_at" timestamptz NOT NULL DEFAULT(now())
);

ALTER TABLE accounts ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");
ALTER TABLE accounts ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency")