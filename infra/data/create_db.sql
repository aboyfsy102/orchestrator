-- drop database if exists
DROP DATABASE IF EXISTS "dev3";

-- create database
CREATE DATABASE "dev3";

-- use database
\c "dev3";

-- create schema for FleetOrder
CREATE TABLE "orders" (
    "id" serial PRIMARY KEY,
    "order_id" varchar(255) NOT NULL,
    "status" varchar(255) NOT NULL,
    "created_at" timestamp DEFAULT now(),
    "deleted_at" timestamp DEFAULT now(),
    "updated_at" timestamp DEFAULT now(),
    "extra_tags" jsonb,
    "profile_name" integer NOT NULL,
    "storage" integer DEFAULT 0,
    "user_data" text,
    "requirements" text,
    "fleet_id" varchar(255) NOT NULL,
    "fleet_response" jsonb,
    "fleet_instances_ips" jsonb
) PARTITION BY RANGE (created_at);

-- create function to create new monthly partition
CREATE OR REPLACE FUNCTION create_monthly_orders_partition() RETURNS void AS $$
DECLARE
    partition_name text;
    partition_start date;
    partition_end date;
BEGIN
    partition_start := date_trunc('month', current_date);
    partition_end := partition_start + interval '1 month';
    partition_name := 'orders_' || to_char(partition_start, 'YYYY_MM');
    
    EXECUTE format('CREATE TABLE IF NOT EXISTS %I PARTITION OF orders FOR VALUES FROM (%L) TO (%L)', 
                   partition_name, partition_start, partition_end);
END;
$$ LANGUAGE plpgsql;

-- create initial partition for the current month
SELECT create_monthly_orders_partition();

-- schedule the function to run monthly
SELECT cron.schedule('create_monthly_orders_partition', '0 0 1 * *', 'SELECT create_monthly_orders_partition()');

-- create profiles table
CREATE TABLE "profiles" (
    "id" serial PRIMARY KEY,
    "name" varchar(255) NOT NULL,
    "ami_id" varchar(255) NOT NULL,
    "storage" integer DEFAULT 0,
    "access_role" varchar(255) NOT NULL,
    "comments" text,
    "user_data_prefix" text,
    "user_data" text,
    "user_data_post" text,
    "tags" jsonb,
);

-- create two users. 'j5v3-admin' and 'j5v3-report'
CREATE USER "j5v3_admin" WITH PASSWORD 'default_t0_change_pl3ase!';
CREATE USER "j5v3_report" WITH PASSWORD 'default_t0_change_pl3ase!';

-- grant permissions to the users
GRANT SELECT, INSERT, UPDATE ON "dev3"."orders" TO "j5v3_admin";
GRANT SELECT ON "dev3"."orders" TO "j5v3-report";
GRANT SELECT, INSERT, UPDATE ON "dev3"."profiles" TO "j5v3_report";
