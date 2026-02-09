-- +goose Up
-- +goose StatementBegin
CREATE TABLE attachments (
	id UUID PRIMARY KEY DEFAULT uuidv7(),
	object_name TEXT NOT NULL UNIQUE,
	file_name TEXT NOT NULL,
	byte_size BIGINT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE admins (
	id UUID PRIMARY KEY DEFAULT uuidv7(),
	name TEXT NOT NULL,
	email_address TEXT NOT NULL UNIQUE,
	password_digest TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE admin_sessions (
	id UUID PRIMARY KEY DEFAULT uuidv7(),
	admin_id UUID NOT NULL REFERENCES admins(id),
	token TEXT NOT NULL UNIQUE,
	ip_address TEXT NOT NULL,
	user_agent TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE users (
	id UUID PRIMARY KEY DEFAULT uuidv7(),
	name TEXT NOT NULL,
	email_address TEXT NOT NULL UNIQUE,
	password_digest TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE user_sessions (
	id UUID PRIMARY KEY DEFAULT uuidv7(),
	user_id UUID NOT NULL REFERENCES users(id),
	token TEXT NOT NULL UNIQUE,
	ip_address TEXT NOT NULL,
	user_agent TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_sessions;

DROP TABLE users;

DROP TABLE admin_sessions;

DROP TABLE admins;

DROP TABLE attachments;
-- +goose StatementEnd
