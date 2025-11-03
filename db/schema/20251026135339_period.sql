-- +goose Up
-- +goose StatementBegin
create table period(
  period_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  tag TEXT UNIQUE NOT NULL,
  description TEXT NOT NULL,
  period_group TEXT NOT NULL UNIQUE
);


insert into period(tag, description, period_group)
VAlUES
    ('fdom', 'Первый день месяца', 'monthly'),
    ('ldom', 'Последний день месяца', 'monthly'),
    ('bi-weekly', 'Каждые две недели', 'bi-weekly');

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
DELETE FROM period
where tag in ('fdom', 'ldom');
drop table payment_plan;
drop table period;
-- +goose StatementEnd
