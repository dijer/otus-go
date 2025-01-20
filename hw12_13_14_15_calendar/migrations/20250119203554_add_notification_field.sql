-- +goose Up
-- +goose StatementBegin
alter table events add column notification_sent boolean default FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table events drop column notification_sent;
-- +goose StatementEnd
