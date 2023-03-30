CREATE TABLE "word" (
  "id" serial PRIMARY KEY,
  "word" varchar UNIQUE NOT NULL,
  "description" varchar NOT NULL,
  "draw_no" int DEFAULT 0,
  "draw_success" int DEFAULT 0,
  "draw_fail" int DEFAULT 0,
  "created_at" timestamp DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" timestamp DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "word_migrations" (
  "id" serial PRIMARY KEY,
  "no" varchar,
  "updated_at" timestamp DEFAULT (CURRENT_TIMESTAMP)
);