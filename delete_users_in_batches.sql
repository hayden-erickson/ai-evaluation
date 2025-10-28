DELIMITER $$

DROP PROCEDURE IF EXISTS `delete_users_in_batches`$$

CREATE PROCEDURE `delete_users_in_batches`()
BEGIN
    -- Variable declarations
    DECLARE done INT DEFAULT FALSE;
    DECLARE num_deleted INT DEFAULT 0;
    DECLARE v_batch_id INT;
    DECLARE v_note VARCHAR(255);
    DECLARE v_error_info VARCHAR(300);
    DECLARE v_rec_cnt INT DEFAULT 0;
    DECLARE v_tab_name VARCHAR(50);
    -- Environment-specific user ID. Modify as per environment.
    -- Production: 2920483, Pre-production: 2387943, Beta: 4365, Alpha: 1033278
    DECLARE v_user INT DEFAULT 2920483; 

    -- Error handler
    DECLARE CONTINUE HANDLER FOR SQLEXCEPTION
    BEGIN
        GET DIAGNOSTICS CONDITION 1
        @p1 = RETURNED_SQLSTATE, @p2 = MYSQL_ERRNO, @p3 = MESSAGE_TEXT;
        SET v_error_info = CONCAT('ERROR ', @p2, ' (', @p1, '): ', @p3);
        ROLLBACK;
        INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
        VALUES (v_batch_id, v_tab_name, 'failed', v_error_info, v_user);
        COMMIT;
        SET done = TRUE;
    END;

    -- Drop temporary table if it exists
    DROP TEMPORARY TABLE IF EXISTS temp_user_ids;

    Outer_check: LOOP
        -- Batch ID Management
        SELECT IFNULL(MAX(batch_id), 0) + 1 INTO v_batch_id FROM deleted_users_info;

        -- Update records with NULL batch_id to new batch_id for eligible users
        UPDATE deleted_users_info dui
        JOIN users u ON dui.user_id = u.id
        JOIN users_roles ur ON u.id = ur.user_id
        JOIN roles r ON ur.role_id = r.id
        SET dui.batch_id = v_batch_id
        WHERE dui.batch_id IS NULL
          AND dui.deletion_status = 'pending'
          AND u.last_login_date < DATE_SUB(CURRENT_TIMESTAMP(), INTERVAL 1 YEAR)
          AND r.name IN ('Tenant', 'Tenant 24/7');
        COMMIT;

        -- User Status Updates for ineligible users
        -- Mark as canceled if user is an owner or shared user in v2_shares
        UPDATE deleted_users_info dui
        SET dui.deletion_status = 'canceled',
            dui.note = 'The user_id is either an Owner or Shared user, so cannot delete'
        WHERE dui.batch_id = v_batch_id AND (
            EXISTS (SELECT 1 FROM v2_shares vs WHERE vs.owner_id = dui.user_id) OR
            EXISTS (SELECT 1 FROM v2_shares vs WHERE vs.recipient_id = dui.user_id)
        );
        COMMIT;

        -- Mark as canceled if user has associated units
        UPDATE deleted_users_info dui
        SET dui.deletion_status = 'canceled',
            dui.note = 'The user_id is having a unit, so cannot delete'
        WHERE dui.batch_id = v_batch_id AND EXISTS (SELECT 1 FROM v2_units vu WHERE vu.user_id = dui.user_id);
        COMMIT;

        -- Mark as canceled if user has associated auctions
        UPDATE deleted_users_info dui
        SET dui.deletion_status = 'canceled',
            dui.note = 'The user_id is having an auction, so cannot delete'
        WHERE dui.batch_id = v_batch_id AND EXISTS (SELECT 1 FROM auction a WHERE a.user_id = dui.user_id);
        COMMIT;

        -- Mark as canceled if user has associated prelets
        UPDATE deleted_users_info dui
        SET dui.deletion_status = 'canceled',
            dui.note = 'The user_id is having a prelet, so cannot delete'
        WHERE dui.batch_id = v_batch_id AND EXISTS (SELECT 1 FROM prelets p WHERE p.user_id = dui.user_id);
        COMMIT;

        -- Mark eligible users as 'in progress'
        UPDATE deleted_users_info
        SET deletion_status = 'in progress'
        WHERE batch_id = v_batch_id AND deletion_status = 'pending';
        COMMIT;

        -- Create and populate temporary table for the current batch
        CREATE TEMPORARY TABLE temp_user_ids (
            id INT NOT NULL PRIMARY KEY,
            user_id INT NOT NULL
        );

        INSERT INTO temp_user_ids (id, user_id)
        SELECT id, user_id FROM deleted_users_info
        WHERE batch_id = v_batch_id AND deletion_status = 'in progress';

        -- Check if batch has records
        SELECT COUNT(*) INTO num_deleted FROM temp_user_ids;
        IF num_deleted = 0 THEN
            LEAVE Outer_check;
        END IF;

        -- Deletion process for each table
        BEGIN
            DECLARE table_list VARCHAR(1024) DEFAULT 'unit_status_updates,unit_overrides,transfer,support_activity,unlock_override_pins,user_acknowledgement,user_login_history,user_site_bookmarks,users_dashboard,users_notifications_settings,users_zone_triggers,v2_app_details,v2_tracking_ids,watch_users,2_factor_auth_pin,access_codes,devices,digital_audits,entry_activity,permissions,pending_notifications,oauth_session_store,note_comments,invalid_entry_attempts,users_roles';
            DECLARE current_pos INT DEFAULT 1;
            DECLARE next_pos INT;
            DECLARE current_table VARCHAR(50);

            table_loop: LOOP
                SET next_pos = LOCATE(',', table_list, current_pos);
                IF next_pos = 0 THEN
                    SET current_table = SUBSTRING(table_list, current_pos);
                ELSE
                    SET current_table = SUBSTRING(table_list, current_pos, next_pos - current_pos);
                END IF;

                SET v_tab_name = current_table;

                SET @sql = CONCAT('SELECT COUNT(*) INTO @cnt FROM ', v_tab_name, ' WHERE user_id IN (SELECT user_id FROM temp_user_ids)');
                PREPARE stmt FROM @sql;
                EXECUTE stmt;
                DEALLOCATE PREPARE stmt;

                IF @cnt > 0 THEN
                    -- Archive data
                    SET @sql = CONCAT('INSERT INTO archived_data.', v_tab_name, '_archived SELECT * FROM ', v_tab_name, ' WHERE user_id IN (SELECT user_id FROM temp_user_ids)');
                    PREPARE stmt FROM @sql;
                    EXECUTE stmt;
                    SET v_rec_cnt = ROW_COUNT();
                    DEALLOCATE PREPARE stmt;
                    
                    SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
                    INSERT INTO user_deletion_detail (batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);

                    -- Delete data
                    SET @sql = CONCAT('DELETE FROM ', v_tab_name, ' WHERE user_id IN (SELECT user_id FROM temp_user_ids) AND EXISTS (SELECT 1 FROM archived_data.', v_tab_name, '_archived a WHERE a.user_id = ', v_tab_name, '.user_id)');
                    PREPARE stmt FROM @sql;
                    EXECUTE stmt;
                    SET v_rec_cnt = ROW_COUNT();
                    DEALLOCATE PREPARE stmt;

                    SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
                    INSERT INTO user_deletion_detail (batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
                    COMMIT;
                END IF;

                IF next_pos = 0 THEN
                    LEAVE table_loop;
                END IF;
                SET current_pos = next_pos + 1;
            END LOOP;
        END;

        -- Delete from users table last
        SET v_tab_name = 'users';
        INSERT INTO archived_data.users_archived SELECT * FROM users WHERE id IN (SELECT user_id FROM temp_user_ids);
        SET v_rec_cnt = ROW_COUNT();
        SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
        INSERT INTO user_deletion_detail (batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);

        DELETE FROM users WHERE id IN (SELECT user_id FROM temp_user_ids);
        SET v_rec_cnt = ROW_COUNT();
        SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
        INSERT INTO user_deletion_detail (batch_id, table_name, note, last_updated_by) VALUES (v_batch_id, v_tab_name, v_note, v_user);
        COMMIT;

        -- Update deleted_users_info status to 'completed'
        UPDATE deleted_users_info
        SET deletion_status = 'completed',
            note = 'the user_id is deleted'
        WHERE id IN (SELECT id FROM temp_user_ids);
        COMMIT;

        -- Clean up temporary table
        DROP TEMPORARY TABLE temp_user_ids;

    END LOOP Outer_check;

    -- Clean up temporary table if loop was exited early
    DROP TEMPORARY TABLE IF EXISTS temp_user_ids;

    SELECT 'User deletion process completed.' AS message;

END$$

DELIMITER ;$$
