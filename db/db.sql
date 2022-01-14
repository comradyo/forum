CREATE EXTENSION IF NOT EXISTS citext;

DROP TABLE IF EXISTS "user" CASCADE;
DROP TABLE IF EXISTS "forum" CASCADE;
DROP TABLE IF EXISTS "thread" CASCADE;
DROP TABLE IF EXISTS "post" CASCADE;
DROP TABLE IF EXISTS "vote" CASCADE;
DROP TABLE IF EXISTS "forum_user" CASCADE;

CREATE UNLOGGED TABLE "user"
(
    id          serial primary key,
    nickname    citext collate "C" not null unique,
    fullname    text not null,
    about       text default '',
    email       citext collate "C" not null unique
);

CREATE UNLOGGED TABLE "forum"
(
    id          serial primary key,
    title       text not null,
    "user"      citext references "user"(nickname) on delete cascade not null,
    slug        citext not null unique,
    posts       int default 0,
    threads     int default 0
);

CREATE UNLOGGED TABLE "thread"
(
    id          serial primary key,
    title       text not null,
    author      citext references "user"(nickname) on delete cascade not null,
    forum       citext references "forum"(slug) on delete cascade not null,
    message     text not null,
    votes       int default 0,
    slug        citext,
    created     timestamp with time zone default now()
);

CREATE UNLOGGED TABLE "post"
(
    id          serial primary key,
    parent      int default 0,
    author      citext references "user"(nickname) on delete cascade not null,
    message     text not null,
    is_edited   bool not null default false,
    forum       citext references "forum"(slug) on delete cascade not null,
    thread      int references "thread"(id) on delete cascade not null,
    created     timestamp with time zone default now(),
    path        int[]
);

CREATE UNLOGGED TABLE "vote"
(
    id          serial primary key,
    thread      int references "thread"(id) on delete cascade not null,
    "user"      citext references "user"(nickname) on delete cascade not null,
    voice       int not null,
    UNIQUE (thread, "user")
);

CREATE UNLOGGED TABLE "forum_user"
(
    id          serial primary key,
    forum       citext references "forum"(slug) on delete cascade not null,
    "user"      citext references "user"(nickname) on delete cascade not null,
    UNIQUE (forum, "user")
);

-----------------------------------------
create or replace function increment_forum_threads() returns trigger as $$
begin
    update forum set threads = threads + 1 where forum.slug = new.forum;
    insert into forum_user (forum, "user") values (new.forum, new.Author)
    on conflict do nothing;
    return null;
end
$$ language 'plpgsql';

create trigger trigger_thread_create
    after insert on thread
    for each row execute procedure increment_forum_threads();
-----------------------------------------

-----------------------------------------
create or replace function update_thread_votes() returns trigger as $$
begin
    update thread set votes = (votes + new.voice) where id = new.thread;
    return new;
end
$$ language 'plpgsql';

create trigger trigger_vote_create
    after insert on vote
    for each row execute procedure update_thread_votes();
-----------------------------------------

-----------------------------------------
create or replace function change_thread_votes() returns trigger as $$
begin
    update thread set votes = (votes - old.voice + new.voice) where id = new.thread;
    return new;
end
$$ language 'plpgsql';

create trigger trigger_vote_change
    after update on vote
    for each row execute procedure change_thread_votes();
-----------------------------------------

-----------------------------------------
create or replace function create_post() returns trigger as $$
declare
    parent_path int[];
begin
    update forum set posts = posts + 1 where slug = new.forum;
    insert into forum_user (forum, "user") values (new.forum, new.Author)
    on conflict do nothing;
    parent_path = (select path from post where id = new.parent limit 1);
    new.path = parent_path || NEW.id;
    return new;
end
$$ language 'plpgsql';

create trigger trigger_create_post
    before insert on post
    for each row execute procedure create_post();
-----------------------------------------

-----------------------------------------

drop index if exists index_user_on_nickname;
drop index if exists index_user_on_email;

drop index if exists index_thread_on_slug;
drop index if exists index_thread_on_id;
drop index if exists index_thread_on_forum;
drop index if exists index_thread_on_forum_and_created;

drop index if exists index_post_on_thread_and_path_and_id;
drop index if exists index_post_on_parent_path_and_id;
drop index if exists index_post_on_thread_and_id;

-----------------------------------------
create index if not exists index_user_on_nickname on "user" using hash(nickname);
create index if not exists index_user_on_email on "user" using hash(email);
-----------------------------------------
create index if not exists index_forum_slug on forum using hash(slug);
-----------------------------------------
create index if not exists index_thread_on_created on thread(created);
create index if not exists index_thread_on_slug on thread using hash(slug);
create index if not exists index_thread_on_id on thread(id);
create index if not exists index_thread_on_forum on thread using hash(forum);

--- +- 2150 rpc (вместе с index_post_on_parent_path_and_id)
drop index if exists index_post_on_thread_and_path_and_id;
create index if not exists index_thread_on_forum_and_created on "thread"(forum, created);

--- +- 2150 rpc (вместе с index_post_on_parent_path_and_id)
drop index if exists index_post_on_thread_and_path_and_id;
create index if not exists index_post_on_thread_and_path_and_id on "post"(thread, path, id);
--- +- 200 rpc
drop index if exists index_post_on_parent;
create index if not exists index_post_on_parent on "post"(parent, id);
--- +- 2150 rpc
drop index if exists index_post_on_parent_path_and_id;
create index if not exists index_post_on_parent_path_and_id on "post"((path[1]), path);

VACUUM;
VACUUM ANALYSE;