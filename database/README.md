# swing-ranger Database

## Setup

### Users

The first script to run is `1_create_users.sql` - this file creates the two users used by the swing-ranger apps.
One user is a read-only user, and the other is for executing commands.

In your PostgreSQL instance, run the following script from the command line.
The users (roles) are created across your instance, so you only have to run this once.

```psql
> \i 1_create_users.sql
```

### Databases

I like to set up three databases, one for test, one for dev, and one for production.
These correspond to the connection strings for the infrastructure tests, my local dev environment, and my "production" environment.

I connect to the database I want to manage, and from a PSQL prompt, I run:

```psql
> \i 2_create_tables.sql
> \i 5_grant_permissions.sql
```

### Full process

1. Install PostgreSQL if not already done.
2. Navigate to the `database` directory in your terminal.
3. Log into `psql`. `sudo -u postgres psql`
4. Create a database: `CREATE DATABASE sr_test;`
5. Connect to the db: `\c sr_test`
6. Run `\i 1_create_users.sql`
7. Run `\i 2_create_tables.sql`
8. Run `\i 5_grant_permissions.sql`

You should be able to complete these steps for any db - just change the name of the database in the `CREATE DATABASE` command and then navigate to it using the `\c` command before running the scripts with `\i`.

## Useful SQL

### Find symbols in watchlists not in eod_prices

```sql
SELECT DISTINCT symbol
FROM public.watchlists w
WHERE NOT EXISTS (
    SELECT 1
    FROM public.eod_prices p
    WHERE p.symbol = w.symbol
)
ORDER BY symbol;
```

## Other notes

### Making a copy of the database.

The following sequence of commands will migrate your `sr_dev` database to `sr_prod`.

```
pg_dump -Fc -v -U postgres -h localhost sr_dev > sr_dev.dump
dropdb -U postgres -h localhost sr_prod
createdb -U postgres -h localhost sr_prod
pg_restore -v -j 8 --clean --if-exists -U postgres -h localhost -d sr_prod sr_dev.dump
```