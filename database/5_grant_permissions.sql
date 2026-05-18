DO $$
DECLARE
    sch text;
    schemas text[] := ARRAY['public'];
BEGIN
    FOREACH sch IN ARRAY schemas
    LOOP
        IF EXISTS (SELECT 1 FROM information_schema.schemata WHERE schema_name = sch) THEN
            -- Full control for admin
            EXECUTE 'GRANT USAGE, CREATE ON SCHEMA ' || sch || ' TO sr_admin;';
            EXECUTE 'GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA ' || sch || ' TO sr_admin;';
            EXECUTE 'GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA ' || sch || ' TO sr_admin;';
            EXECUTE 'ALTER DEFAULT PRIVILEGES IN SCHEMA ' || sch || ' GRANT ALL PRIVILEGES ON TABLES TO sr_admin;';
            EXECUTE 'ALTER DEFAULT PRIVILEGES IN SCHEMA ' || sch || ' GRANT ALL PRIVILEGES ON SEQUENCES TO sr_admin;';

            -- Read-only for reader
            EXECUTE 'GRANT USAGE ON SCHEMA ' || sch || ' TO sr_reader;';
            EXECUTE 'GRANT SELECT ON ALL TABLES IN SCHEMA ' || sch || ' TO sr_reader;';
            EXECUTE 'GRANT SELECT ON ALL SEQUENCES IN SCHEMA ' || sch || ' TO sr_reader;';
            EXECUTE 'ALTER DEFAULT PRIVILEGES IN SCHEMA ' || sch || ' GRANT SELECT ON TABLES TO sr_reader;';
            EXECUTE 'ALTER DEFAULT PRIVILEGES IN SCHEMA ' || sch || ' GRANT SELECT ON SEQUENCES TO sr_reader;';
        END IF;
    END LOOP;
END $$;