CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table driver (
  id uuid not null primary key DEFAULT uuid_generate_v4(),
  created timestamp with time zone not null,
  name varchar(200)
);

create table driver_trackings (
  id uuid not null primary key DEFAULT uuid_generate_v4(),
  driver_id uuid not null,
  created timestamp with time zone not null,
  publisher_time  timestamp,
  latitude double precision,
  longitude double precision
);

insert into driver values
(uuid_generate_v4(), now() - interval '1 hour', 'teste 1'),
(uuid_generate_v4(), now() - interval '2 hour', 'teste 2'),
(uuid_generate_v4(), now() - interval '3 hour', 'teste 3'),
(uuid_generate_v4(), now() - interval '4 hour', 'teste 4'),
(uuid_generate_v4(), now() - interval '5 hour', 'teste 5');
