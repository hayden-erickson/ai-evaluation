-- Stored Procedure: delete_users_in_batches
-- Purpose: Delete users in batches while archiving related data across multiple tables
-- Notes:
--   - Set v_user according to your environment before deploying:
--       Production:     2920483
--       Pre-production: 2387943
--       Beta:           4365
--       Alpha:          1033278
--   - Assumes existence of schema archived_data with mirror tables named <table>_archived
--   - Ensure proper indexes exist as per performance considerations in requirements

DROP PROCEDURE IF EXISTS delete_users_in_batches;
DELIMITER $$
CREATE PROCEDURE delete_users_in_batches()
BEGIN
    -- =============================
    -- Variable Declarations
    -- =============================
    DECLARE done INT DEFAULT FALSE;                -- loop control
    DECLARE num_deleted INT DEFAULT 0;             -- count of records in current batch
    DECLARE v_batch_id INT DEFAULT NULL;           -- current batch identifier
    DECLARE v_note VARCHAR(255) DEFAULT NULL;      -- operation notes
    DECLARE v_error_info VARCHAR(300) DEFAULT NULL;-- error info holder
    DECLARE v_rec_cnt INT DEFAULT 0;               -- record count tracker
    DECLARE v_tab_name VARCHAR(64) DEFAULT NULL;   -- current table name for logging
    DECLARE v_user INT DEFAULT 2920483;            -- CHANGE for environment (see header)

    -- Diagnostics variables for error handler (MySQL 8.0+)
    DECLARE v_sqlstate CHAR(5);
    DECLARE v_errno INT;
    DECLARE v_msg_text TEXT;

    -- =============================
    -- Error Handler
    -- =============================
    DECLARE CONTINUE HANDLER FOR SQLEXCEPTION
    BEGIN
        -- Capture error details
        GET DIAGNOSTICS CONDITION 1 v_errno = MYSQL_ERRNO, v_sqlstate = RETURNED_SQLSTATE, v_msg_text = MESSAGE_TEXT;
        SET v_error_info = CONCAT('ERROR ', v_errno, ' (', v_sqlstate, '): ', v_msg_text);
        -- Rollback whatever was in-flight
        ROLLBACK;
        -- Log error
        INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
        VALUES (v_batch_id, v_tab_name, 'failed', v_error_info, v_user);
        COMMIT; -- commit error log
        -- Stop further processing
        SET done = TRUE;
    END;

    -- =============================
    -- Outer Loop: process batches
    -- =============================
    Outer_check: WHILE NOT done DO

        -- 1) Batch ID Management
        --    - Assign a new batch_id to any NULL/pending rows
        SELECT COALESCE(MAX(batch_id), 0) INTO v_batch_id FROM deleted_users_info;
        SET v_batch_id = v_batch_id + 1;

        START TRANSACTION;
            UPDATE deleted_users_info
            SET batch_id = v_batch_id
            WHERE batch_id IS NULL
              AND deletion_status = 'pending';
        COMMIT;

        -- 2) User Status Updates (canceled reasons first)
        -- Cancel: user involved in shares (owner or recipient)
        SET v_tab_name = 'eligibility-checks';
        START TRANSACTION;
            UPDATE deleted_users_info dui
            JOIN users u ON u.id = dui.user_id
            SET dui.deletion_status = 'canceled',
                dui.note = 'The user_id is either an Owner or Shared user, so cannot delete'
            WHERE dui.batch_id = v_batch_id
              AND dui.deletion_status = 'pending'
              AND (
                    EXISTS (SELECT 1 FROM v2_shares s WHERE s.owner_id = u.id)
                 OR EXISTS (SELECT 1 FROM v2_shares s2 WHERE s2.recipient_id = u.id)
              );
        COMMIT;

        -- Cancel: user has associated units
        START TRANSACTION;
            UPDATE deleted_users_info dui
            JOIN users u ON u.id = dui.user_id
            SET dui.deletion_status = 'canceled',
                dui.note = 'The user_id is having a unit, so cannot delete'
            WHERE dui.batch_id = v_batch_id
              AND dui.deletion_status = 'pending'
              AND EXISTS (SELECT 1 FROM v2_units vu WHERE vu.user_id = u.id);
        COMMIT;

        -- Cancel: user has associated auctions
        START TRANSACTION;
            UPDATE deleted_users_info dui
            JOIN users u ON u.id = dui.user_id
            SET dui.deletion_status = 'canceled',
                dui.note = 'The user_id is having an auction, so cannot delete'
            WHERE dui.batch_id = v_batch_id
              AND dui.deletion_status = 'pending'
              AND EXISTS (SELECT 1 FROM auction a WHERE a.user_id = u.id);
        COMMIT;

        -- Cancel: user has associated prelets
        START TRANSACTION;
            UPDATE deleted_users_info dui
            JOIN users u ON u.id = dui.user_id
            SET dui.deletion_status = 'canceled',
                dui.note = 'The user_id is having a prelet, so cannot delete'
            WHERE dui.batch_id = v_batch_id
              AND dui.deletion_status = 'pending'
              AND EXISTS (SELECT 1 FROM prelets p WHERE p.user_id = u.id);
        COMMIT;

        -- Mark eligible users as in-progress
        START TRANSACTION;
            UPDATE deleted_users_info dui
            JOIN users u ON u.id = dui.user_id
            JOIN users_roles ur ON ur.user_id = u.id
            JOIN roles r ON r.id = ur.role_id
            SET dui.deletion_status = 'in progress'
            WHERE dui.batch_id = v_batch_id
              AND dui.deletion_status = 'pending'
              AND u.last_login_date < DATE_SUB(CURRENT_TIMESTAMP(), INTERVAL 1 YEAR)
              AND r.name IN ('Tenant', 'Tenant 24/7')
              AND NOT EXISTS (SELECT 1 FROM v2_shares s WHERE s.owner_id = u.id)
              AND NOT EXISTS (SELECT 1 FROM v2_shares s2 WHERE s2.recipient_id = u.id)
              AND NOT EXISTS (SELECT 1 FROM v2_units vu WHERE vu.user_id = u.id)
              AND NOT EXISTS (SELECT 1 FROM auction a WHERE a.user_id = u.id)
              AND NOT EXISTS (SELECT 1 FROM prelets p WHERE p.user_id = u.id);
        COMMIT;

        -- 3) Create/Populate temporary table for current batch
        SET v_tab_name = 'temp_user_ids';
        DROP TEMPORARY TABLE IF EXISTS temp_user_ids;
        CREATE TEMPORARY TABLE temp_user_ids (
            id INT NOT NULL PRIMARY KEY,
            user_id INT NOT NULL
        ) ENGINE=MEMORY;

        INSERT INTO temp_user_ids (id, user_id)
        SELECT id, user_id
        FROM deleted_users_info
        WHERE batch_id = v_batch_id
          AND deletion_status = 'in progress';

        SELECT COUNT(*) INTO num_deleted FROM temp_user_ids;
        IF num_deleted = 0 THEN
            -- No eligible users for this batch; exit loop
            SET done = TRUE;
            LEAVE Outer_check;
        END IF;

        -- Helper macro-like pattern per table:
        --   IF EXISTS rows for temp users
        --     INSERT into archived_data.<table>_archived SELECT * FROM <table> JOIN temp_user_ids
        --     Log inserted count
        --     DELETE FROM <table> JOIN temp_user_ids WHERE EXISTS matching archive row
        --     Log deleted count
        --     COMMIT

        -- 4) Process each related table in order

        -- unit_status_updates
        SET v_tab_name = 'unit_status_updates';
        IF EXISTS (SELECT 1 FROM `unit_status_updates` t JOIN temp_user_ids tu ON t.user_id = tu.user_id LIMIT 1) THEN
            START TRANSACTION;
                INSERT INTO archived_data.`unit_status_updates_archived`
                SELECT t.* FROM `unit_status_updates` t JOIN temp_user_ids tu ON t.user_id = tu.user_id;
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);

                DELETE t FROM `unit_status_updates` t
                JOIN temp_user_ids tu ON t.user_id = tu.user_id
                WHERE EXISTS (
                    SELECT 1 FROM archived_data.`unit_status_updates_archived` a
                    WHERE a.user_id = t.user_id
                );
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            COMMIT;
        END IF;

        -- unit_overrides
        SET v_tab_name = 'unit_overrides';
        IF EXISTS (SELECT 1 FROM `unit_overrides` t JOIN temp_user_ids tu ON t.user_id = tu.user_id LIMIT 1) THEN
            START TRANSACTION;
                INSERT INTO archived_data.`unit_overrides_archived`
                SELECT t.* FROM `unit_overrides` t JOIN temp_user_ids tu ON t.user_id = tu.user_id;
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);

                DELETE t FROM `unit_overrides` t
                JOIN temp_user_ids tu ON t.user_id = tu.user_id
                WHERE EXISTS (
                    SELECT 1 FROM archived_data.`unit_overrides_archived` a
                    WHERE a.user_id = t.user_id
                );
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            COMMIT;
        END IF;

        -- transfer
        SET v_tab_name = 'transfer';
        IF EXISTS (SELECT 1 FROM `transfer` t JOIN temp_user_ids tu ON t.user_id = tu.user_id LIMIT 1) THEN
            START TRANSACTION;
                INSERT INTO archived_data.`transfer_archived`
                SELECT t.* FROM `transfer` t JOIN temp_user_ids tu ON t.user_id = tu.user_id;
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);

                DELETE t FROM `transfer` t
                JOIN temp_user_ids tu ON t.user_id = tu.user_id
                WHERE EXISTS (
                    SELECT 1 FROM archived_data.`transfer_archived` a
                    WHERE a.user_id = t.user_id
                );
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            COMMIT;
        END IF;

        -- support_activity
        SET v_tab_name = 'support_activity';
        IF EXISTS (SELECT 1 FROM `support_activity` t JOIN temp_user_ids tu ON t.user_id = tu.user_id LIMIT 1) THEN
            START TRANSACTION;
                INSERT INTO archived_data.`support_activity_archived`
                SELECT t.* FROM `support_activity` t JOIN temp_user_ids tu ON t.user_id = tu.user_id;
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);

                DELETE t FROM `support_activity` t
                JOIN temp_user_ids tu ON t.user_id = tu.user_id
                WHERE EXISTS (
                    SELECT 1 FROM archived_data.`support_activity_archived` a
                    WHERE a.user_id = t.user_id
                );
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            COMMIT;
        END IF;

        -- unlock_override_pins
        SET v_tab_name = 'unlock_override_pins';
        IF EXISTS (SELECT 1 FROM `unlock_override_pins` t JOIN temp_user_ids tu ON t.user_id = tu.user_id LIMIT 1) THEN
            START TRANSACTION;
                INSERT INTO archived_data.`unlock_override_pins_archived`
                SELECT t.* FROM `unlock_override_pins` t JOIN temp_user_ids tu ON t.user_id = tu.user_id;
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);

                DELETE t FROM `unlock_override_pins` t
                JOIN temp_user_ids tu ON t.user_id = tu.user_id
                WHERE EXISTS (
                    SELECT 1 FROM archived_data.`unlock_override_pins_archived` a
                    WHERE a.user_id = t.user_id
                );
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            COMMIT;
        END IF;

        -- user_acknowledgement
        SET v_tab_name = 'user_acknowledgement';
        IF EXISTS (SELECT 1 FROM `user_acknowledgement` t JOIN temp_user_ids tu ON t.user_id = tu.user_id LIMIT 1) THEN
            START TRANSACTION;
                INSERT INTO archived_data.`user_acknowledgement_archived`
                SELECT t.* FROM `user_acknowledgement` t JOIN temp_user_ids tu ON t.user_id = tu.user_id;
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);

                DELETE t FROM `user_acknowledgement` t
                JOIN temp_user_ids tu ON t.user_id = tu.user_id
                WHERE EXISTS (
                    SELECT 1 FROM archived_data.`user_acknowledgement_archived` a
                    WHERE a.user_id = t.user_id
                );
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            COMMIT;
        END IF;

        -- user_login_history
        SET v_tab_name = 'user_login_history';
        IF EXISTS (SELECT 1 FROM `user_login_history` t JOIN temp_user_ids tu ON t.user_id = tu.user_id LIMIT 1) THEN
            START TRANSACTION;
                INSERT INTO archived_data.`user_login_history_archived`
                SELECT t.* FROM `user_login_history` t JOIN temp_user_ids tu ON t.user_id = tu.user_id;
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);

                DELETE t FROM `user_login_history` t
                JOIN temp_user_ids tu ON t.user_id = tu.user_id
                WHERE EXISTS (
                    SELECT 1 FROM archived_data.`user_login_history_archived` a
                    WHERE a.user_id = t.user_id
                );
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            COMMIT;
        END IF;

        -- user_site_bookmarks
        SET v_tab_name = 'user_site_bookmarks';
        IF EXISTS (SELECT 1 FROM `user_site_bookmarks` t JOIN temp_user_ids tu ON t.user_id = tu.user_id LIMIT 1) THEN
            START TRANSACTION;
                INSERT INTO archived_data.`user_site_bookmarks_archived`
                SELECT t.* FROM `user_site_bookmarks` t JOIN temp_user_ids tu ON t.user_id = tu.user_id;
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);

                DELETE t FROM `user_site_bookmarks` t
                JOIN temp_user_ids tu ON t.user_id = tu.user_id
                WHERE EXISTS (
                    SELECT 1 FROM archived_data.`user_site_bookmarks_archived` a
                    WHERE a.user_id = t.user_id
                );
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            COMMIT;
        END IF;

        -- users_dashboard
        SET v_tab_name = 'users_dashboard';
        IF EXISTS (SELECT 1 FROM `users_dashboard` t JOIN temp_user_ids tu ON t.user_id = tu.user_id LIMIT 1) THEN
            START TRANSACTION;
                INSERT INTO archived_data.`users_dashboard_archived`
                SELECT t.* FROM `users_dashboard` t JOIN temp_user_ids tu ON t.user_id = tu.user_id;
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);

                DELETE t FROM `users_dashboard` t
                JOIN temp_user_ids tu ON t.user_id = tu.user_id
                WHERE EXISTS (
                    SELECT 1 FROM archived_data.`users_dashboard_archived` a
                    WHERE a.user_id = t.user_id
                );
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            COMMIT;
        END IF;

        -- users_notifications_settings
        SET v_tab_name = 'users_notifications_settings';
        IF EXISTS (SELECT 1 FROM `users_notifications_settings` t JOIN temp_user_ids tu ON t.user_id = tu.user_id LIMIT 1) THEN
            START TRANSACTION;
                INSERT INTO archived_data.`users_notifications_settings_archived`
                SELECT t.* FROM `users_notifications_settings` t JOIN temp_user_ids tu ON t.user_id = tu.user_id;
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);

                DELETE t FROM `users_notifications_settings` t
                JOIN temp_user_ids tu ON t.user_id = tu.user_id
                WHERE EXISTS (
                    SELECT 1 FROM archived_data.`users_notifications_settings_archived` a
                    WHERE a.user_id = t.user_id
                );
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            COMMIT;
        END IF;

        -- users_zone_triggers
        SET v_tab_name = 'users_zone_triggers';
        IF EXISTS (SELECT 1 FROM `users_zone_triggers` t JOIN temp_user_ids tu ON t.user_id = tu.user_id LIMIT 1) THEN
            START TRANSACTION;
                INSERT INTO archived_data.`users_zone_triggers_archived`
                SELECT t.* FROM `users_zone_triggers` t JOIN temp_user_ids tu ON t.user_id = tu.user_id;
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);

                DELETE t FROM `users_zone_triggers` t
                JOIN temp_user_ids tu ON t.user_id = tu.user_id
                WHERE EXISTS (
                    SELECT 1 FROM archived_data.`users_zone_triggers_archived` a
                    WHERE a.user_id = t.user_id
                );
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            COMMIT;
        END IF;

        -- v2_app_details
        SET v_tab_name = 'v2_app_details';
        IF EXISTS (SELECT 1 FROM `v2_app_details` t JOIN temp_user_ids tu ON t.user_id = tu.user_id LIMIT 1) THEN
            START TRANSACTION;
                INSERT INTO archived_data.`v2_app_details_archived`
                SELECT t.* FROM `v2_app_details` t JOIN temp_user_ids tu ON t.user_id = tu.user_id;
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);

                DELETE t FROM `v2_app_details` t
                JOIN temp_user_ids tu ON t.user_id = tu.user_id
                WHERE EXISTS (
                    SELECT 1 FROM archived_data.`v2_app_details_archived` a
                    WHERE a.user_id = t.user_id
                );
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            COMMIT;
        END IF;

        -- v2_tracking_ids
        SET v_tab_name = 'v2_tracking_ids';
        IF EXISTS (SELECT 1 FROM `v2_tracking_ids` t JOIN temp_user_ids tu ON t.user_id = tu.user_id LIMIT 1) THEN
            START TRANSACTION;
                INSERT INTO archived_data.`v2_tracking_ids_archived`
                SELECT t.* FROM `v2_tracking_ids` t JOIN temp_user_ids tu ON t.user_id = tu.user_id;
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);

                DELETE t FROM `v2_tracking_ids` t
                JOIN temp_user_ids tu ON t.user_id = tu.user_id
                WHERE EXISTS (
                    SELECT 1 FROM archived_data.`v2_tracking_ids_archived` a
                    WHERE a.user_id = t.user_id
                );
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            COMMIT;
        END IF;

        -- watch_users
        SET v_tab_name = 'watch_users';
        IF EXISTS (SELECT 1 FROM `watch_users` t JOIN temp_user_ids tu ON t.user_id = tu.user_id LIMIT 1) THEN
            START TRANSACTION;
                INSERT INTO archived_data.`watch_users_archived`
                SELECT t.* FROM `watch_users` t JOIN temp_user_ids tu ON t.user_id = tu.user_id;
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);

                DELETE t FROM `watch_users` t
                JOIN temp_user_ids tu ON t.user_id = tu.user_id
                WHERE EXISTS (
                    SELECT 1 FROM archived_data.`watch_users_archived` a
                    WHERE a.user_id = t.user_id
                );
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            COMMIT;
        END IF;

        -- 2_factor_auth_pin
        SET v_tab_name = '2_factor_auth_pin';
        IF EXISTS (SELECT 1 FROM `2_factor_auth_pin` t JOIN temp_user_ids tu ON t.user_id = tu.user_id LIMIT 1) THEN
            START TRANSACTION;
                INSERT INTO archived_data.`2_factor_auth_pin_archived`
                SELECT t.* FROM `2_factor_auth_pin` t JOIN temp_user_ids tu ON t.user_id = tu.user_id;
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);

                DELETE t FROM `2_factor_auth_pin` t
                JOIN temp_user_ids tu ON t.user_id = tu.user_id
                WHERE EXISTS (
                    SELECT 1 FROM archived_data.`2_factor_auth_pin_archived` a
                    WHERE a.user_id = t.user_id
                );
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            COMMIT;
        END IF;

        -- access_codes
        SET v_tab_name = 'access_codes';
        IF EXISTS (SELECT 1 FROM `access_codes` t JOIN temp_user_ids tu ON t.user_id = tu.user_id LIMIT 1) THEN
            START TRANSACTION;
                INSERT INTO archived_data.`access_codes_archived`
                SELECT t.* FROM `access_codes` t JOIN temp_user_ids tu ON t.user_id = tu.user_id;
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);

                DELETE t FROM `access_codes` t
                JOIN temp_user_ids tu ON t.user_id = tu.user_id
                WHERE EXISTS (
                    SELECT 1 FROM archived_data.`access_codes_archived` a
                    WHERE a.user_id = t.user_id
                );
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            COMMIT;
        END IF;

        -- devices
        SET v_tab_name = 'devices';
        IF EXISTS (SELECT 1 FROM `devices` t JOIN temp_user_ids tu ON t.user_id = tu.user_id LIMIT 1) THEN
            START TRANSACTION;
                INSERT INTO archived_data.`devices_archived`
                SELECT t.* FROM `devices` t JOIN temp_user_ids tu ON t.user_id = tu.user_id;
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);

                DELETE t FROM `devices` t
                JOIN temp_user_ids tu ON t.user_id = tu.user_id
                WHERE EXISTS (
                    SELECT 1 FROM archived_data.`devices_archived` a
                    WHERE a.user_id = t.user_id
                );
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            COMMIT;
        END IF;

        -- digital_audits
        SET v_tab_name = 'digital_audits';
        IF EXISTS (SELECT 1 FROM `digital_audits` t JOIN temp_user_ids tu ON t.user_id = tu.user_id LIMIT 1) THEN
            START TRANSACTION;
                INSERT INTO archived_data.`digital_audits_archived`
                SELECT t.* FROM `digital_audits` t JOIN temp_user_ids tu ON t.user_id = tu.user_id;
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);

                DELETE t FROM `digital_audits` t
                JOIN temp_user_ids tu ON t.user_id = tu.user_id
                WHERE EXISTS (
                    SELECT 1 FROM archived_data.`digital_audits_archived` a
                    WHERE a.user_id = t.user_id
                );
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            COMMIT;
        END IF;

        -- entry_activity
        SET v_tab_name = 'entry_activity';
        IF EXISTS (SELECT 1 FROM `entry_activity` t JOIN temp_user_ids tu ON t.user_id = tu.user_id LIMIT 1) THEN
            START TRANSACTION;
                INSERT INTO archived_data.`entry_activity_archived`
                SELECT t.* FROM `entry_activity` t JOIN temp_user_ids tu ON t.user_id = tu.user_id;
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);

                DELETE t FROM `entry_activity` t
                JOIN temp_user_ids tu ON t.user_id = tu.user_id
                WHERE EXISTS (
                    SELECT 1 FROM archived_data.`entry_activity_archived` a
                    WHERE a.user_id = t.user_id
                );
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            COMMIT;
        END IF;

        -- permissions
        SET v_tab_name = 'permissions';
        IF EXISTS (SELECT 1 FROM `permissions` t JOIN temp_user_ids tu ON t.user_id = tu.user_id LIMIT 1) THEN
            START TRANSACTION;
                INSERT INTO archived_data.`permissions_archived`
                SELECT t.* FROM `permissions` t JOIN temp_user_ids tu ON t.user_id = tu.user_id;
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);

                DELETE t FROM `permissions` t
                JOIN temp_user_ids tu ON t.user_id = tu.user_id
                WHERE EXISTS (
                    SELECT 1 FROM archived_data.`permissions_archived` a
                    WHERE a.user_id = t.user_id
                );
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            COMMIT;
        END IF;

        -- pending_notifications
        SET v_tab_name = 'pending_notifications';
        IF EXISTS (SELECT 1 FROM `pending_notifications` t JOIN temp_user_ids tu ON t.user_id = tu.user_id LIMIT 1) THEN
            START TRANSACTION;
                INSERT INTO archived_data.`pending_notifications_archived`
                SELECT t.* FROM `pending_notifications` t JOIN temp_user_ids tu ON t.user_id = tu.user_id;
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);

                DELETE t FROM `pending_notifications` t
                JOIN temp_user_ids tu ON t.user_id = tu.user_id
                WHERE EXISTS (
                    SELECT 1 FROM archived_data.`pending_notifications_archived` a
                    WHERE a.user_id = t.user_id
                );
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            COMMIT;
        END IF;

        -- oauth_session_store
        SET v_tab_name = 'oauth_session_store';
        IF EXISTS (SELECT 1 FROM `oauth_session_store` t JOIN temp_user_ids tu ON t.user_id = tu.user_id LIMIT 1) THEN
            START TRANSACTION;
                INSERT INTO archived_data.`oauth_session_store_archived`
                SELECT t.* FROM `oauth_session_store` t JOIN temp_user_ids tu ON t.user_id = tu.user_id;
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);

                DELETE t FROM `oauth_session_store` t
                JOIN temp_user_ids tu ON t.user_id = tu.user_id
                WHERE EXISTS (
                    SELECT 1 FROM archived_data.`oauth_session_store_archived` a
                    WHERE a.user_id = t.user_id
                );
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            COMMIT;
        END IF;

        -- note_comments
        SET v_tab_name = 'note_comments';
        IF EXISTS (SELECT 1 FROM `note_comments` t JOIN temp_user_ids tu ON t.user_id = tu.user_id LIMIT 1) THEN
            START TRANSACTION;
                INSERT INTO archived_data.`note_comments_archived`
                SELECT t.* FROM `note_comments` t JOIN temp_user_ids tu ON t.user_id = tu.user_id;
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);

                DELETE t FROM `note_comments` t
                JOIN temp_user_ids tu ON t.user_id = tu.user_id
                WHERE EXISTS (
                    SELECT 1 FROM archived_data.`note_comments_archived` a
                    WHERE a.user_id = t.user_id
                );
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            COMMIT;
        END IF;

        -- invalid_entry_attempts
        SET v_tab_name = 'invalid_entry_attempts';
        IF EXISTS (SELECT 1 FROM `invalid_entry_attempts` t JOIN temp_user_ids tu ON t.user_id = tu.user_id LIMIT 1) THEN
            START TRANSACTION;
                INSERT INTO archived_data.`invalid_entry_attempts_archived`
                SELECT t.* FROM `invalid_entry_attempts` t JOIN temp_user_ids tu ON t.user_id = tu.user_id;
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);

                DELETE t FROM `invalid_entry_attempts` t
                JOIN temp_user_ids tu ON t.user_id = tu.user_id
                WHERE EXISTS (
                    SELECT 1 FROM archived_data.`invalid_entry_attempts_archived` a
                    WHERE a.user_id = t.user_id
                );
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            COMMIT;
        END IF;

        -- users_roles
        SET v_tab_name = 'users_roles';
        IF EXISTS (SELECT 1 FROM `users_roles` t JOIN temp_user_ids tu ON t.user_id = tu.user_id LIMIT 1) THEN
            START TRANSACTION;
                INSERT INTO archived_data.`users_roles_archived`
                SELECT t.* FROM `users_roles` t JOIN temp_user_ids tu ON t.user_id = tu.user_id;
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);

                DELETE t FROM `users_roles` t
                JOIN temp_user_ids tu ON t.user_id = tu.user_id
                WHERE EXISTS (
                    SELECT 1 FROM archived_data.`users_roles_archived` a
                    WHERE a.user_id = t.user_id
                );
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            COMMIT;
        END IF;

        -- Additional tables from requirements list

        -- users_notifications_settings already handled
        -- Add any missing tables here following the same pattern

        -- 5) Finally, delete from users table (archive then delete)
        SET v_tab_name = 'users';
        IF EXISTS (SELECT 1 FROM `users` t JOIN temp_user_ids tu ON t.id = tu.user_id LIMIT 1) THEN
            START TRANSACTION;
                INSERT INTO archived_data.`users_archived`
                SELECT t.* FROM `users` t JOIN temp_user_ids tu ON t.id = tu.user_id;
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);

                DELETE t FROM `users` t
                JOIN temp_user_ids tu ON t.id = tu.user_id
                WHERE EXISTS (
                    SELECT 1 FROM archived_data.`users_archived` a
                    WHERE a.id = t.id
                );
                SET v_rec_cnt = ROW_COUNT();
                INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
                VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            COMMIT;
        END IF;

        -- 6) Mark deleted_users_info as completed for this batch
        SET v_tab_name = 'deleted_users_info';
        START TRANSACTION;
            UPDATE deleted_users_info d
            JOIN temp_user_ids tu ON tu.id = d.id
            SET d.deletion_status = 'completed',
                d.note = 'the user_id is deleted'
            WHERE d.batch_id = v_batch_id
              AND d.deletion_status = 'in progress';
        COMMIT;

        -- 7) Cleanup temporary table
        DROP TEMPORARY TABLE IF EXISTS temp_user_ids;

        -- Continue outer loop to check for another batch
    END WHILE Outer_check;

    -- Completion message
    SELECT 'User deletion process completed.' AS message;
END $$
DELIMITER ;
