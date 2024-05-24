WITH deleted AS (
    DELETE FROM %s WHERE created_at < NOW() - INTERVAL $1 RETURNING *
) SELECT count(*) from deleted;
