package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// a query to create tables and indices if they don't already exist.
// Can be extracted into a separate sql file and used in conjunction with migration tools.
// But I don't think there's a need for such complexity in this project
const schema = `
-- create resource enum type
do $$
    begin
        if not exists (select 1 from pg_type where typname = 'resource') then
            create type resource as enum ('paper', 'plastic', 'gadgets');
        end if;
    end
$$;

-- create event table
create table if not exists event (
    id varchar(25) primary key,
    date date not null,
    name varchar(25) not null,
    resources_allowed resource[] not null check (cardinality(resources_allowed) >= 1)
);

create index if not exists event_resources_allowed on event using gin(resources_allowed);
create index if not exists event_id_date_name_idx on event(id, date, name);
create index if not exists event_id_name_date_idx on event(id, name, date);

-- create pupil table
create table if not exists pupil (
    id varchar(25) not null primary key,
    class_letter char not null,
    class_year_formed integer not null,
    first_name varchar(25) not null,
    last_name varchar(25) not null,
    text_search tsvector generated always as (to_tsvector('simple', first_name || ' ' || last_name || ' ' ||
                                                                    class_year_formed::text || class_letter || ' ' ||
                                                                    class_letter || ' ' || class_year_formed::text))
        stored
);

create index if not exists pupil_text_search_idx on pupil using gin(text_search);
create index if not exists pupil_class_name_idx on pupil(class_year_formed, class_letter, last_name, first_name);

-- create resources table
create table if not exists resources (
    pupil_id varchar(25) not null,
    event_id varchar(25) not null,
    paper integer not null default 0,
    plastic integer not null default 0,
    gadgets integer not null default 0,
    primary key (pupil_id, event_id),
    foreign key (event_id) references event(id)
        on delete cascade
        on update cascade,
    foreign key (pupil_id) references pupil(id)
        on delete cascade
        on update cascade
);

create index if not exists resources_event_id_gadgets_idx on resources(event_id, gadgets desc nulls last);
create index if not exists resources_event_id_paper_idx on resources(event_id, paper desc nulls last);
create index if not exists resources_event_id_plastic_idx on resources(event_id, plastic desc nulls last);
`

// ValidateSchema creates tables and indices if they don't already exist
func ValidateSchema(ctx context.Context, db *sqlx.DB) error {
	if _, err := db.ExecContext(ctx, schema); err != nil {
		return err
	}
	return nil
}
