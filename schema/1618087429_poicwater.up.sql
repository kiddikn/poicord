CREATE TABLE poic_water (
	id serial not null,
	user_id varchar(255) not null,
	started_at timestamp not null,
    finished_at timestamp,
    revoked_at timestamp,
    created_at timestamp not null default CURRENT_TIMESTAMP,
	primary key(id)
);