Problematic
-----------

Sqrible came out of the [frustation](https://github.com/jackc/pgx/issues/253) of working with
[pgx](https://github.com/jackc/pgx) and [sqlx](https://github.com/jmoiron/sqlx).

I had to buid a [crappy bridge](https://github.com/jeromer/pgx/commit/5097a8100cb350853e1aa8bcf05787ff41c69216)
which led to nowhere and this did not solve the problem that for every request
[sqlx](https://github.com/jmoiron/sqlx)has to use Golang's reflection to map
query results to struct fields.

This did not sound right.

The other option is to map each field manually but that's tedious and error prone.

So the solution is to generate go code that would do everything automatically.

The advantages are:

- I no longer need sqlx and thus limit reflection usage to pgx
- With a set of template I can generate the code I need and update it easily
  with a few template changes

How does it work ?
------------------

Sqrible scans informations about a given table, reads what you configured for
this table in a yaml config file and output the generated code you built with the
provided template.

Building
--------

    make deps build

Usage
-----

    ./bin/sqrible --help

    Usage of ./bin/sqrible:
      -c string
            /path/to/config/file.yaml
      -d string
            template dir
      -t string
            table name

Example
-------

Given the following table:

                                            Table "public.users"
    ┌──────────────────┬──────────────────────────┬────────────────────────────────────────────────────┐
    │      Column      │           Type           │                     Modifiers                      │
    ├──────────────────┼──────────────────────────┼────────────────────────────────────────────────────┤
    │ id               │ bigint                   │ not null default nextval('users_id_seq'::regclass) │
    │ uid              │ text                     │ not null                                           │
    │ password         │ text                     │ not null                                           │
    │ email            │ text                     │ not null                                           │
    │ username         │ text                     │ not null default ''::text                          │
    │ enabled          │ boolean                  │ not null default false                             │
    │ email_confirmed  │ boolean                  │ not null default false                             │
    │ creation_date    │ timestamp with time zone │ not null default now()                             │
    │ update_date      │ timestamp with time zone │ not null default now()                             │
    │ avatar_file_name │ text                     │ not null default ''::text                          │
    └──────────────────┴──────────────────────────┴────────────────────────────────────────────────────┘

Which can be created with the following SQL:

    CREATE TABLE public.users
    (
      id bigint NOT NULL DEFAULT nextval('users_id_seq'::regclass),
      uid text NOT NULL,
      password text NOT NULL,
      email text NOT NULL,
      username text NOT NULL DEFAULT ''::text,
      enabled boolean NOT NULL DEFAULT false,
      email_confirmed boolean NOT NULL DEFAULT false,
      creation_date timestamp with time zone NOT NULL DEFAULT now(),
      update_date timestamp with time zone NOT NULL DEFAULT now(),
      avatar_file_name text NOT NULL DEFAULT ''::text,
      CONSTRAINT users_pk PRIMARY KEY (id),
      CONSTRAINT users_email_key UNIQUE (email),
      CONSTRAINT users_uid_key UNIQUE (uid),
      CONSTRAINT users_username_key UNIQUE (username)
    )
    WITH (
      OIDS=FALSE
    );
    ALTER TABLE public.users
      OWNER TO postgres;

And a `sqrible.yml` configuration file with the following contents:

    tables:
      users:                     # table name
        template: example.tpl    # template to use for this table
        gostruct: User           # go struct which will be used for this table
        tablecols:               # table column configuration
          id: s                  # id can only be SELECTed
          uid: s                 # uid can only be SELECTed
          password: "-"          # password is totally ignored
          creation_date: s       # creation_date can only be selected
          username: i            # username can only be INSERTed (not available in SELECT nor UPDATE)
          enabled: u             # enabled can only be UPDATEd (not availble in SELECT not INSERT)
          email: siu             # email can be SELECTed INSERTed UPDATEd

          # no configuration means : SELECTable, INSERTable, UPDATEable

You can now run:

    make build && PGDATABASE=yourdb PGHOST=localhost ./bin/sqrible -c ./sqrible.yml -d . -t users

You will see the following output:

    Example template
    ----------------

    Table:
    - name: users
    - go struct name: User
    - template used: example.tpl





    SELECTable columns
    ------------------


    id bigint 1 True ID pgtype.Int8

    uid text 2 False UID pgtype.Text

    email text 4 False Email pgtype.Text

    email_confirmed boolean 7 False EmailConfirmed string

    creation_date timestamp with time zone 8 False CreationDate pgtype.Timestamptz

    update_date timestamp with time zone 9 False UpdateDate pgtype.Timestamptz

    avatar_file_name text 10 False AvatarFileName pgtype.Text



    INSERTable columns
    ------------------


    email text 4 False Email pgtype.Text

    username text 5 False Username pgtype.Text

    email_confirmed boolean 7 False EmailConfirmed string

    update_date timestamp with time zone 9 False UpdateDate pgtype.Timestamptz

    avatar_file_name text 10 False AvatarFileName pgtype.Text



    UPDATEable columns
    ------------------


    email text 4 False Email pgtype.Text

    enabled boolean 6 False Enabled string

    email_confirmed boolean 7 False EmailConfirmed string

    update_date timestamp with time zone 9 False UpdateDate pgtype.Timestamptz

    avatar_file_name text 10 False AvatarFileName pgtype.Text



    Primary keys
    ------------

    id bigint 1 True ID pgtype.Int8


    all columns
    -----------

    id bigint 1 True ID pgtype.Int8 True True False False

    uid text 2 False UID pgtype.Text True True False False

    email text 4 False Email pgtype.Text True True True True

    username text 5 False Username pgtype.Text True False True False

    enabled boolean 6 False Enabled string True False False True

    email_confirmed boolean 7 False EmailConfirmed string False True True True

    creation_date timestamp with time zone 8 False CreationDate pgtype.Timestamptz True True False False

    update_date timestamp with time zone 9 False UpdateDate pgtype.Timestamptz False True True True

    avatar_file_name text 10 False AvatarFileName pgtype.Text False True True True


Since you have all the informations you need about a specific table it becomes easy to
write some templates which will generate the Go code you need. The usage of
[pongo2](https://github.com/flosch/pongo2) makes it relatively easy.

Other projects in this field
----------------------------

- [pgxdata](https://github.com/jackc/pgxdata)  (which I used for inspiration)
- [xo](https://github.com/knq/xo)
- [genieql](https://bitbucket.org/jatone/genieql)

FAQ
===

Do you plan to support other DBMSes ?
-------------------------------------

No. I only need Postgres support and do not plan to provide any support for
MySQL, Oracle ... Have at look at [xo](https://github.com/knq/xo) ;)

Do you provide base templates ?
-------------------------------

No. This kind of code is too much project specific. I prefer providing enough
informations in the `Table` template variables, and an [example template](example.tpl)
so people can do whatever they want for their project.
