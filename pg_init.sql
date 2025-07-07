-- Создаем пользователя, если не существует
DO $$
BEGIN
  IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'st1txh') THEN
CREATE ROLE st1txh WITH LOGIN PASSWORD 'St1txh_0000';
END IF;
END $$;

-- Создаем БД, если не существует
SELECT 'CREATE DATABASE st1txh_db OWNER st1txh'
    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'st1txh_db')\gexec