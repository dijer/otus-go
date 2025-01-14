-- +goose Up
-- +goose StatementBegin
create table events (
 id serial primary key,
 owner bigint,
 title text,
 description text,
 start_time timestamp not null,
 end_time timestamp not null
);
create index owner_idx on events (owner);
create index start_idx on events using btree (start_time);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table events;
-- +goose StatementEnd
