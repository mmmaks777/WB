CREATE TABLE "orders" ("id" bigserial,"created_at" timestamptz,"updated_at" timestamptz,"deleted_at" timestamptz,"order_uid" text,"track_number" text,"entry" text,
"delivery_name" text,"delivery_phone" text,"delivery_zip" text,"delivery_city" text,"delivery_address" text,"delivery_region" text,"delivery_email" text,
"payment_transaction" text,"payment_request_id" text,"payment_currency" text,"payment_provider" text,"payment_amount" bigint,"payment_payment_dt" bigint,
"payment_bank" text,"payment_delivery_cost" bigint,"payment_goods_total" bigint,"payment_custom_fee" bigint,"items" json,"locale" text,"internal_signature" text,
"customer_id" text,"delivery_service" text,"shard_key" text,"sm_id" bigint,"date_created" text,"oof_shard" text,PRIMARY KEY ("id"))