WITH deleted AS (
    DELETE FROM %s WHERE created_at < $1 RETURNING *
) SELECT count(*) from deleted;
