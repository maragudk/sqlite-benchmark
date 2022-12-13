create table posts (
  id serial primary key,
  title text not null,
  content text not null,
  created text not null default now()
);

create table comments (
  id serial primary key,
  post_id int not null references posts (id),
  name text not null,
  content text not null,
  created text not null default now()
);

create index comment_created_idx on comments (created);
