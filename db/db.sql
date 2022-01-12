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
    slug        citext unique,
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
    update thread set votes = votes + new.voice where id = new.thread;
    return null;
end
$$ language 'plpgsql';

create trigger trigger_vote_create
    after insert on vote
    for each row execute procedure update_thread_votes();
-----------------------------------------

-----------------------------------------
create or replace function change_thread_votes() returns trigger as $$
begin
    if new.voice <> old.voice then
        if new.voice > 0 then
            update thread set votes = (votes + 2) where id = new.thread;
        else
            update thread set votes = (votes - 2) where id = new.thread;
        end if;
    end if;
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