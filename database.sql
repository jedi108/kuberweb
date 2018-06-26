
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL
);

CREATE UNIQUE INDEX idx_users_email on users (email);


create table if not exists history_action
(
	id serial not null
		constraint history_action_pkey
			primary key,
	action varchar(128) not null
)
;

create unique index if not exists history_action_id_uindex
	on history_action (id)
;

create unique index if not exists history_action_action_uindex
	on history_action (action)
;



create table if not exists history_user_actions
(
	id bigserial not null
		constraint history_user_actions_pkey
			primary key,
	timeinsert timestamp default timezone('utc'::text, now()) not null,
	user_id bigserial not null
		constraint history_user_actions_users_id_fk
			references users,
	action_id integer not null
		constraint history_user_actions_history_action_id_fk
			references history_action,
	isdone boolean not null,
	result varchar(255),
	reqparam varchar(128)
)
;

create unique index if not exists history_user_actions_id_uindex
	on history_user_actions (id)
;

INSERT INTO public.history_action (id, action) VALUES (1, 'setScale');
INSERT INTO public.history_action (id, action) VALUES (2, 'flushRedis');