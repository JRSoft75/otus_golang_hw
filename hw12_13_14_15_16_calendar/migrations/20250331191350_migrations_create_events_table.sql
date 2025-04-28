-- +goose Up
-- +goose StatementBegin
create table events (
    id serial primary key,
    title varchar(255) not null,
    description text  null,
    userId integer not null,
    startAt timestamp not null,
    endAt timestamp not null,
    notifyBefore integer null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table events;
-- +goose StatementEnd
