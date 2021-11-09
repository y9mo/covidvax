# Database migration scripts

We use [migrate](https://github.com/golang-migrate/migrate)
to handle migrations.

### Create a new migration

 Run this command to create a new migration (up/down):

```
./vendor/migrate create -ext sql -dir db/migrations/ <MIGRATION_NAME>
```

### Apply migrations

Apply all migrations:

```
make migrations OPTS=up
```

#### More Commands

Remove all migrations:

```
make migrations OPTS=down
```

Apply only last migration:

```
make migrations OPTS=up 1
```

Remove only last migration changes:

```
make migrations OPTS="down 1"
```
