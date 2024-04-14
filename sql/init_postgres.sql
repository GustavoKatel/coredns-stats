-- create table
CREATE TABLE IF NOT EXISTS %[1]s (
    id SERIAL PRIMARY KEY,
    domain varchar,
    question_type varchar,
    metadata JSONB,
    created_at timestamp default (now() at time zone 'utc')
);

-- create index
CREATE INDEX IF NOT EXISTS domain_index ON %[1]s (domain);
CREATE INDEX IF NOT EXISTS question_type_index ON %[1]s (question_type);
CREATE INDEX IF NOT EXISTS created_at_index ON %[1]s (created_at);
