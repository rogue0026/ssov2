BEGIN;
    DROP INDEX IF EXISTS T_Users_user_id_idx;
    DROP TABLE IF EXISTS T_Users;
COMMIT;