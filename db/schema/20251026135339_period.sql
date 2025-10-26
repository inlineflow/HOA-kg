-- +goose Up
-- +goose StatementBegin
insert into period(period)
VAlUES
    ('first_day_of_the_month'),
    ('last_day_of_the_month');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM period
where TEXT in ('first_day_of_the_month', 'last_day_of_the_month');
-- +goose StatementEnd
