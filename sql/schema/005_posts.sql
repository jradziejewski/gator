-- +goose Up
create table posts(
	id uuid primary key,
	created_at timestamp not null,
	updated_at timestamp not null,
	title varchar not null,
	url varchar unique not null,
	description varchar not null,
	published_at timestamp not null,
	feed_id uuid not null,
	foreign key (feed_id) references feeds(id) on delete cascade
);

-- +goose Down
drop table posts;
