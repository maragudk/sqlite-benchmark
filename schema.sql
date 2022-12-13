create table posts (
  id integer primary key,
  title text not null,
  content text not null,
  created text not null default (strftime('%Y-%m-%dT%H:%M:%fZ'))
);

create table comments (
  id integer primary key,
  post_id int not null references posts (id),
  name text not null,
  content text not null,
  created text not null default (strftime('%Y-%m-%dT%H:%M:%fZ'))
);

create index comment_created_idx on comments (created);
