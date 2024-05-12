DELETE FROM %s WHERE created_at < NOW() - INTERVAL $1;
