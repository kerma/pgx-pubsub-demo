# Just use Postgres

Demo pub/sub implementation using 

* [NOTIFY](https://www.postgresql.org/docs/15/sql-notify.html)
* [LISTEN](https://www.postgresql.org/docs/15/sql-listen.html)
* [pg\_advisory\_lock](https://www.postgresql.org/docs/15/functions-admin.html#FUNCTIONS-ADVISORY-LOCKS) for running single subscriber at a time 
* [pgx](https://github.com/jackc/pgx) - PostgreSQL driver 
* [migrate](https://github.com/golang-migrate/migrate) - database migrations
