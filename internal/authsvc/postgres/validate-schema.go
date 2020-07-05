package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// a query to create tables and indices if they don't already exist.
// Can be extracted into a separate sql file and used in conjunction with migration tools
const schema = `
-- create role enum type
do
$$
    begin
        if not exists(select 1 from pg_type where typname = 'role') then
            create type role as enum ('admin', 'member');
        end if;
    end
$$;

-- create users table
create table if not exists users
(
    id               varchar(50) primary key,
    active           bool        not null,
    activation_token text        not null,
    email            varchar(50) not null unique,
    first_name       varchar(35) not null,
    last_name        varchar(35) not null,
    password_hash    text        not null,
    role             role        not null,
    text_search      tsvector generated always as (to_tsvector('simple', first_name || ' ' || last_name || ' ' ||
                                                                         email)) stored
);

create index if not exists users_text_search_idx on users using gin (text_search);
create index if not exists users_active_idx on users (active);
create index if not exists users_activation_token_idx on users (activation_token);

-- create clients table. Clients in a sense of browsers, apps etc
create table if not exists clients
(
    id            varchar(25) primary key,
    refresh_token text not null
);

-- create user_client table
create table if not exists user_client
(
    user_id   varchar(25) not null,
    client_id varchar(25) not null,
    primary key (user_id, client_id),
    foreign key (client_id) references clients (id)
        on delete cascade
        on update cascade,
    foreign key (user_id) references users (id)
        on delete cascade
        on update cascade
);
`

// ValidateSchema creates tables and indices if they don't already exist
func ValidateSchema(ctx context.Context, db *pgxpool.Pool) error {
	if _, err := db.Exec(ctx, schema); err != nil {
		return err
	}
	return nil
}
