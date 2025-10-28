-- =====================================================================
-- Stored Procedure: delete_users_in_batches
-- Description: Deletes users in batches while archiving their data
--              across multiple related tables. Handles errors, logs
--              operations, and processes users based on specific
--              eligibility criteria.
-- =====================================================================

-- Drop the procedure if it already exists to ensure a clean re-creation.
DROP PROCEDURE IF EXISTS delete_users_in_batches;

-- Change the delimiter to $$ to allow for semicolons within the procedure body.
DELIMITER $$

CREATE PROCEDURE delete_users_in_batches()
BEGIN
    -- =====================================================================
    -- Variable Declarations
    -- =====================================================================
    DECLARE done INT DEFAULT FALSE;
    DECLARE num_to_process INT DEFAULT 0;
    DECLARE v_batch_id INT;
    DECLARE v_note VARCHAR(255);
    DECLARE v_error_info VARCHAR(300);
    DECLARE v_rec_cnt INT DEFAULT 0;
    DECLARE v_tab_name VARCHAR(50);

    -- Environment-specific user ID for logging/auditing.
    -- This value should be updated based on the target environment.
    -- Production: 2920483
    -- Pre-production: 2387943
    -- Beta: 4365
    -- Alpha: 1033278
    DECLARE v_user INT DEFAULT 2920483; -- SET FOR PRODUCTION

    -- =====================================================================
    -- Error Handler
    -- =====================================================================
    -- This handler will catch any SQL exceptions that occur during execution.
    DECLARE CONTINUE HANDLER FOR SQLEXCEPTION
    BEGIN
        GET DIAGNOSTICS CONDITION 1
            @sqlstate = RETURNED_SQLSTATE,
            @errno = MYSQL_ERRNO,
            @text = MESSAGE_TEXT;
        SET v_error_info = CONCAT('ERROR ', @errno, ' (', @sqlstate, '): ', @text);
        
        -- Rollback the failed transaction and log the error.
        ROLLBACK;
        INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
        VALUES (v_batch_id, v_tab_name, 'failed', v_error_info, v_user);
        COMMIT; -- Commit the error log entry.
        
        -- Set 'done' to TRUE to exit the processing loop on error.
        SET done = TRUE;
    END;

    -- =====================================================================
    -- Main Processing Loop
    -- =====================================================================
    -- This loop processes users in batches. Each iteration handles one batch.
    Outer_check: LOOP
        -- If the error handler was triggered, exit the loop.
        IF done THEN
            LEAVE Outer_check;
        END IF;

        -- Check for any users in 'pending' state without a batch_id.
        SELECT count(*) INTO num_to_process FROM deleted_users_info WHERE deletion_status = 'pending' AND batch_id IS NULL;
        IF num_to_process = 0 THEN
            LEAVE Outer_check; -- No more users to process, exit loop.
        END IF;

        -- =================================================================
        -- 1. Batch ID Management
        -- =================================================================
        -- Determine the next batch_id and assign it to pending users.
        SELECT IFNULL(MAX(batch_id), 0) + 1 INTO v_batch_id FROM deleted_users_info;
        
        UPDATE deleted_users_info SET batch_id = v_batch_id WHERE batch_id IS NULL AND deletion_status = 'pending';
        COMMIT;

        -- =================================================================
        -- 2. User Status Updates (Eligibility Checks)
        -- =================================================================
        -- Mark users as 'canceled' if they do not meet deletion criteria.
        -- Note: These checks are performed on users in the current batch with 'pending' status.

        -- Condition: User is an owner or shared user in v2_shares.
        UPDATE deleted_users_info d
        SET d.deletion_status = 'canceled', d.note = 'The user_id is either an Owner or Shared user, so cannot delete'
        WHERE d.batch_id = v_batch_id AND d.deletion_status = 'pending'
          AND EXISTS (SELECT 1 FROM v2_shares s WHERE s.owner_id = d.user_id OR s.recipient_id = d.user_id);

        -- Condition: User has associated units in v2_units.
        UPDATE deleted_users_info d
        SET d.deletion_status = 'canceled', d.note = 'The user_id is having a unit, so cannot delete'
        WHERE d.batch_id = v_batch_id AND d.deletion_status = 'pending'
          AND EXISTS (SELECT 1 FROM v2_units u WHERE u.user_id = d.user_id);

        -- Condition: User has associated auctions.
        UPDATE deleted_users_info d
        SET d.deletion_status = 'canceled', d.note = 'The user_id is having an auction, so cannot delete'
        WHERE d.batch_id = v_batch_id AND d.deletion_status = 'pending'
          AND EXISTS (SELECT 1 FROM auction a WHERE a.user_id = d.user_id);

        -- Condition: User has associated prelets.
        UPDATE deleted_users_info d
        SET d.deletion_status = 'canceled', d.note = 'The user_id is having a prelet, so cannot delete'
        WHERE d.batch_id = v_batch_id AND d.deletion_status = 'pending'
          AND EXISTS (SELECT 1 FROM prelets p WHERE p.user_id = d.user_id);

        -- Mark eligible users as 'in progress'.
        -- Eligibility: last login > 1 year ago AND role is 'Tenant' or 'Tenant 24/7'.
        UPDATE deleted_users_info d
        JOIN users u ON d.user_id = u.id
        JOIN users_roles ur ON u.id = ur.user_id
        JOIN roles r ON ur.role_id = r.id
        SET d.deletion_status = 'in progress'
        WHERE d.batch_id = v_batch_id
          AND d.deletion_status = 'pending'
          AND u.last_login_date < DATE_SUB(CURRENT_TIMESTAMP(), INTERVAL 1 YEAR)
          AND r.name IN ('Tenant', 'Tenant 24/7');

        -- Any remaining 'pending' users in this batch did not meet all criteria.
        UPDATE deleted_users_info
        SET deletion_status = 'canceled', note = 'User did not meet all eligibility criteria (e.g., role, last login).'
        WHERE batch_id = v_batch_id AND deletion_status = 'pending';
        
        COMMIT;

        -- =================================================================
        -- 3. Create and Populate Temporary Table
        -- =================================================================
        -- This table holds the IDs of users to be processed in this batch.
        DROP TEMPORARY TABLE IF EXISTS temp_user_ids;
        CREATE TEMPORARY TABLE temp_user_ids (
            id INT NOT NULL PRIMARY KEY,
            user_id INT NOT NULL,
            INDEX(user_id)
        );

        INSERT INTO temp_user_ids (id, user_id)
        SELECT id, user_id FROM deleted_users_info WHERE batch_id = v_batch_id AND deletion_status = 'in progress';
        
        -- =================================================================
        -- 4. Check if Batch Has Records and Process
        -- =================================================================
        SELECT COUNT(*) INTO v_rec_cnt FROM temp_user_ids;
        IF v_rec_cnt > 0 THEN
            -- Helper macro for archive and delete operations to reduce redundancy.
            -- This is a conceptual guide; MySQL procedures don't have macros.
            -- The logic is implemented manually for each table below.

            -- PROCESSING LOGIC FOR EACH TABLE:
            -- 1. SET v_tab_name = 'table_to_process';
            -- 2. IF EXISTS (SELECT 1 FROM `table_to_process` WHERE user_id IN (SELECT user_id FROM temp_user_ids)) THEN
            -- 3.   START TRANSACTION;
            -- 4.   INSERT INTO archived_data.`table_to_process_archived` SELECT * FROM `table_to_process` WHERE user_id IN ...;
            -- 5.   -- Log insertion
            -- 6.   DELETE FROM `table_to_process` WHERE user_id IN ...;
            -- 7.   -- Log deletion
            -- 8.   COMMIT;
            -- 9. END IF;

            -- Process each of the 29 related tables.
            
            SET v_tab_name = 'unit_status_updates';
            IF EXISTS (SELECT 1 FROM unit_status_updates t WHERE t.user_id IN (SELECT user_id FROM temp_user_ids)) THEN
                START TRANSACTION;
                INSERT INTO archived_data.unit_status_updates_archived SELECT * FROM unit_status_updates WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                DELETE FROM unit_status_updates WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                COMMIT;
            END IF;

            SET v_tab_name = 'unit_overrides';
            IF EXISTS (SELECT 1 FROM unit_overrides t WHERE t.user_id IN (SELECT user_id FROM temp_user_ids)) THEN
                START TRANSACTION;
                INSERT INTO archived_data.unit_overrides_archived SELECT * FROM unit_overrides WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                DELETE FROM unit_overrides WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                COMMIT;
            END IF;

            SET v_tab_name = 'transfer';
            IF EXISTS (SELECT 1 FROM transfer t WHERE t.user_id IN (SELECT user_id FROM temp_user_ids)) THEN
                START TRANSACTION;
                INSERT INTO archived_data.transfer_archived SELECT * FROM transfer WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                DELETE FROM transfer WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                COMMIT;
            END IF;

            SET v_tab_name = 'support_activity';
            IF EXISTS (SELECT 1 FROM support_activity t WHERE t.user_id IN (SELECT user_id FROM temp_user_ids)) THEN
                START TRANSACTION;
                INSERT INTO archived_data.support_activity_archived SELECT * FROM support_activity WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                DELETE FROM support_activity WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                COMMIT;
            END IF;

            SET v_tab_name = 'unlock_override_pins';
            IF EXISTS (SELECT 1 FROM unlock_override_pins t WHERE t.user_id IN (SELECT user_id FROM temp_user_ids)) THEN
                START TRANSACTION;
                INSERT INTO archived_data.unlock_override_pins_archived SELECT * FROM unlock_override_pins WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                DELETE FROM unlock_override_pins WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                COMMIT;
            END IF;

            SET v_tab_name = 'user_acknowledgement';
            IF EXISTS (SELECT 1 FROM user_acknowledgement t WHERE t.user_id IN (SELECT user_id FROM temp_user_ids)) THEN
                START TRANSACTION;
                INSERT INTO archived_data.user_acknowledgement_archived SELECT * FROM user_acknowledgement WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                DELETE FROM user_acknowledgement WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                COMMIT;
            END IF;

            SET v_tab_name = 'user_login_history';
            IF EXISTS (SELECT 1 FROM user_login_history t WHERE t.user_id IN (SELECT user_id FROM temp_user_ids)) THEN
                START TRANSACTION;
                INSERT INTO archived_data.user_login_history_archived SELECT * FROM user_login_history WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                DELETE FROM user_login_history WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                COMMIT;
            END IF;

            SET v_tab_name = 'user_site_bookmarks';
            IF EXISTS (SELECT 1 FROM user_site_bookmarks t WHERE t.user_id IN (SELECT user_id FROM temp_user_ids)) THEN
                START TRANSACTION;
                INSERT INTO archived_data.user_site_bookmarks_archived SELECT * FROM user_site_bookmarks WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                DELETE FROM user_site_bookmarks WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                COMMIT;
            END IF;

            SET v_tab_name = 'users_dashboard';
            IF EXISTS (SELECT 1 FROM users_dashboard t WHERE t.user_id IN (SELECT user_id FROM temp_user_ids)) THEN
                START TRANSACTION;
                INSERT INTO archived_data.users_dashboard_archived SELECT * FROM users_dashboard WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                DELETE FROM users_dashboard WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                COMMIT;
            END IF;

            SET v_tab_name = 'users_notifications_settings';
            IF EXISTS (SELECT 1 FROM users_notifications_settings t WHERE t.user_id IN (SELECT user_id FROM temp_user_ids)) THEN
                START TRANSACTION;
                INSERT INTO archived_data.users_notifications_settings_archived SELECT * FROM users_notifications_settings WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                DELETE FROM users_notifications_settings WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                COMMIT;
            END IF;

            SET v_tab_name = 'users_zone_triggers';
            IF EXISTS (SELECT 1 FROM users_zone_triggers t WHERE t.user_id IN (SELECT user_id FROM temp_user_ids)) THEN
                START TRANSACTION;
                INSERT INTO archived_data.users_zone_triggers_archived SELECT * FROM users_zone_triggers WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                DELETE FROM users_zone_triggers WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                COMMIT;
            END IF;

            SET v_tab_name = 'v2_app_details';
            IF EXISTS (SELECT 1 FROM v2_app_details t WHERE t.user_id IN (SELECT user_id FROM temp_user_ids)) THEN
                START TRANSACTION;
                INSERT INTO archived_data.v2_app_details_archived SELECT * FROM v2_app_details WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                DELETE FROM v2_app_details WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                COMMIT;
            END IF;

            SET v_tab_name = 'v2_tracking_ids';
            IF EXISTS (SELECT 1 FROM v2_tracking_ids t WHERE t.user_id IN (SELECT user_id FROM temp_user_ids)) THEN
                START TRANSACTION;
                INSERT INTO archived_data.v2_tracking_ids_archived SELECT * FROM v2_tracking_ids WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                DELETE FROM v2_tracking_ids WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                COMMIT;
            END IF;

            SET v_tab_name = 'watch_users';
            IF EXISTS (SELECT 1 FROM watch_users t WHERE t.user_id IN (SELECT user_id FROM temp_user_ids)) THEN
                START TRANSACTION;
                INSERT INTO archived_data.watch_users_archived SELECT * FROM watch_users WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                DELETE FROM watch_users WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                COMMIT;
            END IF;

            SET v_tab_name = '2_factor_auth_pin';
            IF EXISTS (SELECT 1 FROM `2_factor_auth_pin` t WHERE t.user_id IN (SELECT user_id FROM temp_user_ids)) THEN
                START TRANSACTION;
                INSERT INTO archived_data.`2_factor_auth_pin_archived` SELECT * FROM `2_factor_auth_pin` WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                DELETE FROM `2_factor_auth_pin` WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                COMMIT;
            END IF;

            SET v_tab_name = 'access_codes';
            IF EXISTS (SELECT 1 FROM access_codes t WHERE t.user_id IN (SELECT user_id FROM temp_user_ids)) THEN
                START TRANSACTION;
                INSERT INTO archived_data.access_codes_archived SELECT * FROM access_codes WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                DELETE FROM access_codes WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                COMMIT;
            END IF;

            SET v_tab_name = 'devices';
            IF EXISTS (SELECT 1 FROM devices t WHERE t.user_id IN (SELECT user_id FROM temp_user_ids)) THEN
                START TRANSACTION;
                INSERT INTO archived_data.devices_archived SELECT * FROM devices WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                DELETE FROM devices WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                COMMIT;
            END IF;

            SET v_tab_name = 'digital_audits';
            IF EXISTS (SELECT 1 FROM digital_audits t WHERE t.user_id IN (SELECT user_id FROM temp_user_ids)) THEN
                START TRANSACTION;
                INSERT INTO archived_data.digital_audits_archived SELECT * FROM digital_audits WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                DELETE FROM digital_audits WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                COMMIT;
            END IF;

            SET v_tab_name = 'entry_activity';
            IF EXISTS (SELECT 1 FROM entry_activity t WHERE t.user_id IN (SELECT user_id FROM temp_user_ids)) THEN
                START TRANSACTION;
                INSERT INTO archived_data.entry_activity_archived SELECT * FROM entry_activity WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                DELETE FROM entry_activity WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                COMMIT;
            END IF;

            SET v_tab_name = 'permissions';
            IF EXISTS (SELECT 1 FROM permissions t WHERE t.user_id IN (SELECT user_id FROM temp_user_ids)) THEN
                START TRANSACTION;
                INSERT INTO archived_data.permissions_archived SELECT * FROM permissions WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                DELETE FROM permissions WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                COMMIT;
            END IF;

            SET v_tab_name = 'pending_notifications';
            IF EXISTS (SELECT 1 FROM pending_notifications t WHERE t.user_id IN (SELECT user_id FROM temp_user_ids)) THEN
                START TRANSACTION;
                INSERT INTO archived_data.pending_notifications_archived SELECT * FROM pending_notifications WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                DELETE FROM pending_notifications WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                COMMIT;
            END IF;

            SET v_tab_name = 'oauth_session_store';
            IF EXISTS (SELECT 1 FROM oauth_session_store t WHERE t.user_id IN (SELECT user_id FROM temp_user_ids)) THEN
                START TRANSACTION;
                INSERT INTO archived_data.oauth_session_store_archived SELECT * FROM oauth_session_store WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                DELETE FROM oauth_session_store WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                COMMIT;
            END IF;

            SET v_tab_name = 'note_comments';
            IF EXISTS (SELECT 1 FROM note_comments t WHERE t.user_id IN (SELECT user_id FROM temp_user_ids)) THEN
                START TRANSACTION;
                INSERT INTO archived_data.note_comments_archived SELECT * FROM note_comments WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                DELETE FROM note_comments WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                COMMIT;
            END IF;

            SET v_tab_name = 'invalid_entry_attempts';
            IF EXISTS (SELECT 1 FROM invalid_entry_attempts t WHERE t.user_id IN (SELECT user_id FROM temp_user_ids)) THEN
                START TRANSACTION;
                INSERT INTO archived_data.invalid_entry_attempts_archived SELECT * FROM invalid_entry_attempts WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                DELETE FROM invalid_entry_attempts WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                COMMIT;
            END IF;

            SET v_tab_name = 'users_roles';
            IF EXISTS (SELECT 1 FROM users_roles t WHERE t.user_id IN (SELECT user_id FROM temp_user_ids)) THEN
                START TRANSACTION;
                INSERT INTO archived_data.users_roles_archived SELECT * FROM users_roles WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                DELETE FROM users_roles WHERE user_id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                COMMIT;
            END IF;
            
            -- =================================================================
            -- 5. Process the main `users` table last.
            -- =================================================================
            SET v_tab_name = 'users';
            IF EXISTS (SELECT 1 FROM users t WHERE t.id IN (SELECT user_id FROM temp_user_ids)) THEN
                START TRANSACTION;
                INSERT INTO archived_data.users_archived SELECT * FROM users WHERE id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                DELETE FROM users WHERE id IN (SELECT user_id FROM temp_user_ids);
                SET v_rec_cnt = ROW_COUNT(); SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
                INSERT INTO user_deletion_detail(batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                COMMIT;
            END IF;

            -- =================================================================
            -- 6. Update Final Status for Completed Users
            -- =================================================================
            UPDATE deleted_users_info
            SET deletion_status = 'completed', note = 'the user_id is deleted'
            WHERE id IN (SELECT id FROM temp_user_ids);
            COMMIT;

        END IF;

        -- =================================================================
        -- 7. Clean Up Temporary Table
        -- =================================================================
        DROP TEMPORARY TABLE IF EXISTS temp_user_ids;

    END LOOP Outer_check;

    -- =====================================================================
    -- Completion Message
    -- =====================================================================
    SELECT 'User deletion process completed.' AS message;

END$$

-- Reset the delimiter back to the default semicolon.
DELIMITER ;

