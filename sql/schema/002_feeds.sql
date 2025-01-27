-- +goose Up
create table feeds(
	id uuid primary key,
	created_at timestamp not null,
	updated_at timestamp not null,
	name varchar not null,
	url varchar unique not null,
	user_id uuid not null,
	foreign key (user_id) references users(id) on delete cascade
);

-- +goose Down
drop table feeds;
