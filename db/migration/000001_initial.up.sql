CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE "user" (
    "id" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    "username" varchar(255) NOT NULL UNIQUE,
    "email" varchar(255) NOT NULL UNIQUE,
    "password" varchar(255) NOT NULL,
    "phone_number" varchar(20) UNIQUE,
    "full_name" varchar(255),
    "avatar" varchar(255),
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "password_change_at" timestamptz NOT NULL DEFAULT now()
)