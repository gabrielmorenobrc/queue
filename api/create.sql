create schema queue;

set SEARCH_PATH = 'queue';

create table entry
(
  id         bigint                   not null,
  context    varchar(255)             not null,
  date       timestamp with time zone not null,
  data       jsonb                    not null,
  errorcount integer                  not null,
  primary key (id)
);

create index i_entry_1 on entry (context);
create index i_entry_2 on entry (context, id);
create index i_entry_3 on entry (context, errorcount);

create table success
(
  id       bigint                   not null,
  entry_id bigint                   not null,
  date     timestamp with time zone not null,
  primary key (id),
  foreign key (entry_id) references entry (id)
);

create unique index u_success_1 on success (entry_id);

create table error
(
  id       bigint                   not null,
  entry_id bigint                   not null,
  date     timestamp with time zone not null,
  data     jsonb                    not null,
  primary key (id),
  foreign key (entry_id) references entry (id)
);

create index i_error_1 on error (entry_id);