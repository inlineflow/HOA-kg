-- +goose Up
-- +goose StatementBegin
-- SELECT 'up SQL query';
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE table house(
  house_id UUID default uuid_generate_v4() PRIMARY KEY,
  address TEXT
);

CREATE table flat(
  flat_id UUID default uuid_generate_v4() PRIMARY KEY,
  flat_number int NOT NULL,
  house_id UUID,

  CONSTRAINT fk_flat_house
    FOREIGN KEY (house_id)
    REFERENCES house(house_id)
);

CREATE table account(
  account_id UUID default uuid_generate_v4() PRIMARY KEY,
  balance NUMERIC(12,2),
  name TEXT
  -- flat_id UUID UNIQUE,


);

CREATE table payment(
  payment_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  account_id UUID,
  note TEXT,

    CONSTRAINT fk_payment_account_id
      FOREIGN KEY (account_id)
      REFERENCES account(account_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table flat;
drop table house;
-- +goose StatementEnd
