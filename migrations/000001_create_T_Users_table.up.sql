BEGIN;
    CREATE TABLE IF NOT EXISTS T_Users (
        user_id SERIAL PRIMARY KEY,
        login varchar(50) NOT NULL UNIQUE,
        password_hash varchar(255) NOT NULL,
        email varchar(50) NOT NULL
    );
    CREATE INDEX IF NOT EXISTS T_Users_user_id_idx on T_Users(user_id);
COMMIT;