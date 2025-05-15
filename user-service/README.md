# User service readme

### Database setup

Set up the MySQL (or MariaDB) database.
```sql
CREATE DATABASE users CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### Environment variables

In the `./user-service/` directory create `.env` file with the following parameters:
 
```
DB_USER=root
DB_PASSWORD=
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=users
```
> ⚠️ Remember to set your parameters accordingly, values given above are defaults.