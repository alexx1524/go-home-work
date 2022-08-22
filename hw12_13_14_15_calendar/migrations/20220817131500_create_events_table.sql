-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.events
(
    id uuid NOT NULL,
    title VARCHAR not null,
    description TEXT null,
    user_id uuid NOT NULL,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    CONSTRAINT "PK_id" PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE if exists public.events
-- +goose StatementEnd
