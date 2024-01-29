CREATE TABLE "users_info" (
                              "id" bigserial PRIMARY KEY,
                              "username" varchar(20) UNIQUE NOT NULL,
                              "hashed_password" varchar(256) NOT NULL,
                              "full_name" varchar(30) NOT NULL,
                              "email" varchar(25) UNIQUE NOT NULL,
                              "phone" varchar(25) UNIQUE,
                              "gender" smallint NOT NULL DEFAULT 0,
                              "avatar" varchar(512),
                              "deleted_flag" smallint NOT NULL DEFAULT 0,
                              "password_changed_at" date NOT NULL DEFAULT '0001-01-01 00:00:00Z',
                              "created_at" date NOT NULL DEFAULT (now()),
                              "updated_at" date NOT NULL DEFAULT (now())
);

CREATE TABLE "post_info" (
                             "id" bigserial PRIMARY KEY,
                             "title" varchar(70) NOT NULL,
                             "content" varchar(2048) NOT NULL,
                             "total_price" varchar(10) NOT NULL DEFAULT '0',
                             "post_user_id" bigint NOT NULL,
                             "delivery_type" smallint NOT NULL,
                             "area" varchar(100),
                             "item_num" int NOT NULL DEFAULT 0,
                             "post_status" smallint NOT NULL DEFAULT 0,
                             "negotiable" smallint NOT NULL DEFAULT 1,
                             "images" varchar(512) NOT NULL,
                             "deleted_flag" smallint NOT NULL DEFAULT 0,
                             "created_at" date NOT NULL DEFAULT (now()),
                             "updated_at" date NOT NULL DEFAULT (now())
);

CREATE TABLE "trading_history" (
                                       "id" bigserial PRIMARY KEY,
                                       "post_id" bigint NOT NULL,
                                       "sold_to_user_id" bigint NOT NULL,
                                       "seller_id" bigint NOT NULL,
                                       "price" varchar(10) NOT NULL DEFAULT '0',
                                       "deleted_flag" smallint NOT NULL DEFAULT 0,
                                       "created_at" date NOT NULL DEFAULT (now()),
                                       "updated_at" date NOT NULL DEFAULT (now())
);

CREATE TABLE "interest_info" (
                                 "id" bigserial PRIMARY KEY,
                                 "post_id" bigint NOT NULL,
                                 "interested_user_id" bigint NOT NULL,
                                 "created_at" date NOT NULL DEFAULT (now()),
                                 "updated_at" date NOT NULL DEFAULT (now())
);

CREATE INDEX ON "post_info" ("post_user_id");

CREATE UNIQUE INDEX ON "trading_history" ("post_id", "sold_to_user_id");

CREATE INDEX ON "interest_info" ("post_id");

CREATE INDEX ON "interest_info" ("interested_user_id");

CREATE UNIQUE INDEX ON "interest_info" ("post_id", "interested_user_id");

COMMENT ON COLUMN "users_info"."username" IS 'unique username';

COMMENT ON COLUMN "users_info"."hashed_password" IS 'encrypted password';

COMMENT ON COLUMN "users_info"."full_name" IS 'lastname firstname';

COMMENT ON COLUMN "users_info"."email" IS 'unique email address';

COMMENT ON COLUMN "users_info"."phone" IS 'unique, including country code';

COMMENT ON COLUMN "users_info"."gender" IS '0: femail, 1: male, 2: other';

COMMENT ON COLUMN "users_info"."avatar" IS 'user icon address';

COMMENT ON COLUMN "users_info"."deleted_flag" IS '0: active, 1: deleted';

COMMENT ON COLUMN "post_info"."title" IS 'post title';

COMMENT ON COLUMN "post_info"."content" IS 'post content';

COMMENT ON COLUMN "post_info"."total_price" IS 'post price, accurate to cent';

COMMENT ON COLUMN "post_info"."post_user_id" IS 'user who posted this post';

COMMENT ON COLUMN "post_info"."delivery_type" IS '0: pick up. 1: mail';

COMMENT ON COLUMN "post_info"."area" IS 'the area that the seller wants to trade';

COMMENT ON COLUMN "post_info"."item_num" IS 'total items in this post';

COMMENT ON COLUMN "post_info"."post_status" IS '0: active, 1: sold, 2: inactive';

COMMENT ON COLUMN "post_info"."negotiable" IS '0: not negotiable, 1: negotiable';

COMMENT ON COLUMN "post_info"."images" IS 'post images, separated by comma';

COMMENT ON COLUMN "post_info"."deleted_flag" IS '0: active, 1: deleted';

COMMENT ON COLUMN "trading_history"."post_id" IS 'related post id';

COMMENT ON COLUMN "trading_history"."sold_to_user_id" IS 'user who bought at least one item in the post';

COMMENT ON COLUMN "trading_history"."seller_id" IS 'seller id save for later usage';

COMMENT ON COLUMN "trading_history"."price" IS 'true transaction price';

COMMENT ON COLUMN "interest_info"."post_id" IS 'related post id';

COMMENT ON COLUMN "interest_info"."interested_user_id" IS 'user interested to this post';

-- ALTER TABLE "post_info" ADD FOREIGN KEY ("post_user_id") REFERENCES "users_info" ("id");
--
-- ALTER TABLE "trading_history" ADD FOREIGN KEY ("post_id") REFERENCES "post_info" ("id");
--
-- ALTER TABLE "trading_history" ADD FOREIGN KEY ("sold_to_user_id") REFERENCES "users_info" ("id");
--
-- ALTER TABLE "interest_info" ADD FOREIGN KEY ("post_id") REFERENCES "post_info" ("id");
--
-- ALTER TABLE "interest_info" ADD FOREIGN KEY ("interested_user_id") REFERENCES "users_info" ("id");

-- DROP TABLE IF EXISTS users_info;
-- DROP TABLE IF EXISTS post_info;
-- DROP TABLE IF EXISTS trading_history;
-- DROP TABLE IF EXISTS interest_info;
