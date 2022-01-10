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

/*
 TODO: сделать процедуры, обновляющие таблицу forum_user и голоса на ветках
 */

CREATE UNLOGGED TABLE "forum_user"
(
    id          serial primary key,
    forum       citext references "forum"(slug) on delete cascade not null,
    "user"      citext references "user"(nickname) on delete cascade not null,
    UNIQUE (forum, "user")
);