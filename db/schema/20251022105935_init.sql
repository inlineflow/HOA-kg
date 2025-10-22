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
  address TEXT,
  house_id UUID,

  CONSTRAINT fk_flat_house
    FOREIGN KEY (house_id)
    REFERENCES house(house_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table flat;
drop table house;
-- +goose StatementEnd
