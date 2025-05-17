# User service readme

### Database setup

Set up the PostgreSQL database.
```sql
CREATE DATABASE users;
```

### Environment variables

In the `./user-service/` directory create `.env` file with the following parameters:
 
```
DB_USER=postgres
DB_PASSWORD=
DB_HOST=127.0.0.1
DB_PORT=5432
DB_NAME=users
DB_SSLMODE=disable
DB_TIMEZONE=UTC
```
> ⚠️ Remember to set your parameters accordingly, values given above are defaults.