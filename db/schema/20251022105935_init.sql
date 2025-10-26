-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE table house(
  house_id UUID default uuid_generate_v4() PRIMARY KEY,
  address TEXT NOT NULL
);

CREATE table flat(
  flat_id UUID default uuid_generate_v4() PRIMARY KEY,
  flat_number int NOT NULL,
  house_id UUID NOT NULL,
  area INT NOT NULL,

  CONSTRAINT fk_flat_house
    FOREIGN KEY (house_id)
    REFERENCES house(house_id)
);

CREATE table resident(
  resident_id UUID default uuid_generate_v4() PRIMARY KEY,
  flat_id UUID NOT NULL,
  phone_number VARCHAR(15) NOT NULL,
  fio TEXT NOT NULL,

  CONSTRAINT resident_phone_number_digits_only
    CHECK (phone_number ~ '^\d{1,15}$'),

  CONSTRAINT redisent_flat
    FOREIGN KEY (flat_id)
    REFERENCES flat(flat_id)
);

CREATE table flat_owner(
  flat_owner_id UUID default uuid_generate_v4() PRIMARY KEY,
  flat_id UUID,
  resident_id UUID NOT NULL, 

  CONSTRAINT flat_owner_flat
    FOREIGN KEY (flat_id)
    REFERENCES flat(flat_id),

  CONSTRAINT flat_owner_resident
    FOREIGN KEY (resident_id)
    REFERENCES resident(resident_id)
);

CREATE table account(
  account_id UUID default uuid_generate_v4() PRIMARY KEY,
  balance DECIMAL DEFAULT 0,
  code INT NOT NULL,
  name TEXT NOT NULL

);

create table flat_account(
  flat_account_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  flat_id UUID UNIQUE,
  account_id UUID UNIQUE,

  CONSTRAINT flat_account_flat
    FOREIGN KEY (flat_id)
    REFERENCES flat(flat_id),

  CONSTRAINT flat_account_account
    FOREIGN KEY (account_id)
    REFERENCES account(account_id)
);

CREATE table payment(
  payment_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  account_id UUID NOT NULL,
  amount DECIMAL NOT NULL,
  note TEXT,

    CONSTRAINT fk_payment_account_id
      FOREIGN KEY (account_id)
      REFERENCES account(account_id)
);

create table period(
  period_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  period TEXT UNIQUE NOT NULL
);

create table payment_plan(
  payment_plan_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  house_id UUID NOT NULL,
  period_id UUID NOT NULL,
  date_from TIMESTAMP NOT NULL,
  date_to TIMESTAMP,
  monthly_amount DECIMAL NOT NULL,

    CONSTRAINT fk_payment_plan_period
      FOREIGN KEY (period_id)
      REFERENCES period(period_id),

    CONSTRAINT fk_payment_plan_house_id
      FOREIGN KEY (house_id)
      REFERENCES house(house_id)
);
-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin
drop table payment_plan;
drop table period;
drop table flat_owner;
drop table resident;
drop table flat_account;
drop table flat;
drop table house;
drop table payment;
drop table account;
-- drop table table_mapping;
-- +goose StatementEnd
