This is a heavily modified version of sqrible - and a work in progress. 
Do not use unless you like the following changes and can customize to your needs.   

* Automatic detection of nullable fields
* Chooses different data types to be used for nullable fields. 
* Added "package" variable to the table yaml which I needed in my templates
* Changed transposition of data types away from pgx - this is a work in progress
* Detection of data types handles user defined types - this is a work in progress


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

Examples
-------

Have a look in the `examples` directory which contains one simple template
which showcases informations sqrible has about a table. The `advanced` example
showcases a more real life usage.

Demo
----

<a href="https://asciinema.org/a/UM4BUS1rL0EnT3OhcMLlQccVF"><img src="https://asciinema.org/a/UM4BUS1rL0EnT3OhcMLlQccVF.png" width="700"/></a>

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
