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

- I no longer need [sqlx](https://github.com/jmoiron/sqlx) and thus limit reflection usage to [pgx](https://github.com/jackc/pgx)
- With a set of template I can generate the code I need and update it easily with a few template changes

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
      users:                  # table name
        template: example.tpl # template to use for this table
        gostruct: User        # go struct which will be used for this table
        tablecols:            # table column configuration
          id:
            access: s         # id can only be SELECTed
          uid:
            access: s         # uid can only be SELECTed
          password:
            access: "-"       # password is totally ignored
          creation_date:
            access: s         # creation_date can only be selected
          username:
            access: i         # username can only be INSERTed (not available in SELECT nor UPDATE)
          enabled:
            access: u         # enabled can only be UPDATEd (not availble in SELECT not INSERT)
          email:
            access: siu       # email can be SELECTed INSERTed UPDATEd

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

    PGColumnName      : id
    PGDataType        : bigint
    PGOrdinalPosition : 1
    IsPK              : True
    GoFieldName       : ID
    PgxType           : pgtype.Int8
    JSON              : custom_id_field

    PGColumnName      : uid
    PGDataType        : text
    PGOrdinalPosition : 2
    IsPK              : False
    GoFieldName       : UID
    PgxType           : pgtype.Text
    JSON              : uid

    PGColumnName      : email
    PGDataType        : text
    PGOrdinalPosition : 4
    IsPK              : False
    GoFieldName       : Email
    PgxType           : pgtype.Text
    JSON              : email

    PGColumnName      : email_confirmed
    PGDataType        : boolean
    PGOrdinalPosition : 7
    IsPK              : False
    GoFieldName       : EmailConfirmed
    PgxType           : pgtype.Bool
    JSON              : email_confirmed

    PGColumnName      : creation_date
    PGDataType        : timestamp with time zone
    PGOrdinalPosition : 8
    IsPK              : False
    GoFieldName       : CreationDate
    PgxType           : pgtype.Timestamptz
    JSON              : creation_date

    PGColumnName      : update_date
    PGDataType        : timestamp with time zone
    PGOrdinalPosition : 9
    IsPK              : False
    GoFieldName       : UpdateDate
    PgxType           : pgtype.Timestamptz
    JSON              : update_date

    PGColumnName      : avatar_file_name
    PGDataType        : text
    PGOrdinalPosition : 10
    IsPK              : False
    GoFieldName       : AvatarFileName
    PgxType           : pgtype.Text
    JSON              : avatar_file_name




    INSERTable columns
    ------------------

    PGColumnName      : email
    PGDataType        : text
    PGOrdinalPosition : 4
    IsPK              : False
    GoFieldName       : Email
    PgxType           : pgtype.Text
    JSON              : email

    PGColumnName      : username
    PGDataType        : text
    PGOrdinalPosition : 5
    IsPK              : False
    GoFieldName       : Username
    PgxType           : pgtype.Text
    JSON              : username

    PGColumnName      : email_confirmed
    PGDataType        : boolean
    PGOrdinalPosition : 7
    IsPK              : False
    GoFieldName       : EmailConfirmed
    PgxType           : pgtype.Bool
    JSON              : email_confirmed

    PGColumnName      : update_date
    PGDataType        : timestamp with time zone
    PGOrdinalPosition : 9
    IsPK              : False
    GoFieldName       : UpdateDate
    PgxType           : pgtype.Timestamptz
    JSON              : update_date

    PGColumnName      : avatar_file_name
    PGDataType        : text
    PGOrdinalPosition : 10
    IsPK              : False
    GoFieldName       : AvatarFileName
    PgxType           : pgtype.Text
    JSON              : avatar_file_name




    UPDATEable columns
    ------------------

    PGColumnName      : email
    PGDataType        : text
    PGOrdinalPosition : 4
    IsPK              : False
    GoFieldName       : Email
    PgxType           : pgtype.Text
    JSON              : email

    PGColumnName      : enabled
    PGDataType        : boolean
    PGOrdinalPosition : 6
    IsPK              : False
    GoFieldName       : Enabled
    PgxType           : pgtype.Bool
    JSON              : enabled

    PGColumnName      : email_confirmed
    PGDataType        : boolean
    PGOrdinalPosition : 7
    IsPK              : False
    GoFieldName       : EmailConfirmed
    PgxType           : pgtype.Bool
    JSON              : email_confirmed

    PGColumnName      : update_date
    PGDataType        : timestamp with time zone
    PGOrdinalPosition : 9
    IsPK              : False
    GoFieldName       : UpdateDate
    PgxType           : pgtype.Timestamptz
    JSON              : update_date

    PGColumnName      : avatar_file_name
    PGDataType        : text
    PGOrdinalPosition : 10
    IsPK              : False
    GoFieldName       : AvatarFileName
    PgxType           : pgtype.Text
    JSON              : avatar_file_name




    Primary keys
    ------------
    PGColumnName      : id
    PGDataType        : bigint
    PGOrdinalPosition : 1
    IsPK              : True
    GoFieldName       : ID
    PgxType           : pgtype.Int8
    JSON              : custom_id_field



    all columns
    -----------
    PGColumnName      : id
    PGDataType        : bigint
    PGOrdinalPosition : 1
    IsPK              : True
    GoFieldName       : ID
    PgxType           : pgtype.Int8
    JSON              : custom_id_field

    PGColumnName      : uid
    PGDataType        : text
    PGOrdinalPosition : 2
    IsPK              : False
    GoFieldName       : UID
    PgxType           : pgtype.Text
    JSON              : uid

    PGColumnName      : email
    PGDataType        : text
    PGOrdinalPosition : 4
    IsPK              : False
    GoFieldName       : Email
    PgxType           : pgtype.Text
    JSON              : email

    PGColumnName      : username
    PGDataType        : text
    PGOrdinalPosition : 5
    IsPK              : False
    GoFieldName       : Username
    PgxType           : pgtype.Text
    JSON              : username

    PGColumnName      : enabled
    PGDataType        : boolean
    PGOrdinalPosition : 6
    IsPK              : False
    GoFieldName       : Enabled
    PgxType           : pgtype.Bool
    JSON              : enabled

    PGColumnName      : email_confirmed
    PGDataType        : boolean
    PGOrdinalPosition : 7
    IsPK              : False
    GoFieldName       : EmailConfirmed
    PgxType           : pgtype.Bool
    JSON              : email_confirmed

    PGColumnName      : creation_date
    PGDataType        : timestamp with time zone
    PGOrdinalPosition : 8
    IsPK              : False
    GoFieldName       : CreationDate
    PgxType           : pgtype.Timestamptz
    JSON              : creation_date

    PGColumnName      : update_date
    PGDataType        : timestamp with time zone
    PGOrdinalPosition : 9
    IsPK              : False
    GoFieldName       : UpdateDate
    PgxType           : pgtype.Timestamptz
    JSON              : update_date

    PGColumnName      : avatar_file_name
    PGDataType        : text
    PGOrdinalPosition : 10
    IsPK              : False
    GoFieldName       : AvatarFileName
    PgxType           : pgtype.Text
    JSON              : avatar_file_name


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
