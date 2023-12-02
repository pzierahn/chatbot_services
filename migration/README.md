# Migration Guide From Supabase to Firebase

## Database

To run database migrations, you can use the following command:

```shell
docker run --rm -it postgres /bin/bash

pg_dump -Fc -v -d $SUPABASE_DB -f brainboost.dump
pg_restore -v -d $AWS_BRAINBOOST_DB brainboost.dump
```

```sql
ALTER TABLE collections
    ADD COLUMN new_user_id VARCHAR(36); -- Adjust the size based on your UUID representation (e.g., VARCHAR(36) for standard UUIDs)

UPDATE collections
    SET new_user_id = CAST(user_id AS CHAR(36));

ALTER TABLE collections
    DROP COLUMN user_id;

ALTER TABLE collections
    RENAME COLUMN new_user_id TO user_id;
```
