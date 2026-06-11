# About

This is my general backend server with administration tool.

The goal is to provide user CRUD in the admin tool that can then be authenticated/operated on by the server
via JWTs.

# WIP -- This is the idea design

# Admin

```bash
./bin/admin database new /database/location/databaseName.db 

./bin/admin /database/location/databaseName.db user add <email> <password>

# Outputs user id in the database

# Supplying --is-admin should make the new user an admin

./bin/admin user del <email | user-id>

./bin/admin user password <email> <new-password> 
# Resets the user's password to whatever was given - errors if the user isn't found

./bin/admin user role [admin|standard] <email | user-id>

# Users are NOT marked as verified on creation, even if admin 
# to mark them as valid

./bin/admin user verify <email | user-id>
```

