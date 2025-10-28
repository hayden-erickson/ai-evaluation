-- ============================================================================
-- Stored Procedure: delete_users_in_batches
-- Description: Deletes users in batches while archiving their data across
--              multiple related tables with comprehensive error handling
-- ============================================================================

DELIMITER $$

DROP PROCEDURE IF EXISTS delete_users_in_batches$$

CREATE PROCEDURE delete_users_in_batches()
BEGIN
    -- ========================================================================
    -- Variable Declarations
    -- ========================================================================
    DECLARE done INT DEFAULT FALSE;
    DECLARE num_deleted INT;
    DECLARE v_batch_id INT;
    DECLARE v_note VARCHAR(255);
    DECLARE v_error_info VARCHAR(300);
    DECLARE v_rec_cnt INT DEFAULT 0;
    DECLARE v_tab_name VARCHAR(50);
    DECLARE v_user INT DEFAULT 2920483;  -- Production: 2920483, Pre-prod: 2387943, Beta: 4365, Alpha: 1033278
    
    -- ========================================================================
    -- Error Handler
    -- ========================================================================
    DECLARE CONTINUE HANDLER FOR SQLEXCEPTION
    BEGIN
        GET DIAGNOSTICS CONDITION 1
            v_error_info = CONCAT('ERROR ', MYSQL_ERRNO, ' (', RETURNED_SQLSTATE, '): ', MESSAGE_TEXT);
        
        ROLLBACK;
        
        INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
        VALUES (v_batch_id, v_tab_name, 'failed', v_error_info, v_user);
        
        COMMIT;
        
        SET done = TRUE;
    END;
    
    -- ========================================================================
    -- Outer Loop: Process batches until no eligible users remain
    -- ========================================================================
    Outer_check: WHILE NOT done DO
        
        -- ====================================================================
        -- Step 1: Assign or increment batch_id
        -- ====================================================================
        SELECT COALESCE(MAX(batch_id), 0) + 1 INTO v_batch_id FROM deleted_users_info;
        
        UPDATE deleted_users_info
        SET batch_id = v_batch_id
        WHERE batch_id IS NULL;
        
        COMMIT;
        
        -- ====================================================================
        -- Step 2: Mark users as 'canceled' if they have v2_shares
        -- ====================================================================
        UPDATE deleted_users_info dui
        INNER JOIN users u ON dui.user_id = u.id
        SET dui.deletion_status = 'canceled',
            dui.note = 'The user_id is either an Owner or Shared user, so cannot delete'
        WHERE dui.batch_id = v_batch_id
            AND dui.deletion_status = 'pending'
            AND EXISTS (
                SELECT 1 FROM v2_shares 
                WHERE owner_id = u.id OR recipient_id = u.id
            );
        
        -- ====================================================================
        -- Step 3: Mark users as 'canceled' if they have v2_units
        -- ====================================================================
        UPDATE deleted_users_info dui
        INNER JOIN users u ON dui.user_id = u.id
        SET dui.deletion_status = 'canceled',
            dui.note = 'The user_id is having a unit, so cannot delete'
        WHERE dui.batch_id = v_batch_id
            AND dui.deletion_status = 'pending'
            AND EXISTS (
                SELECT 1 FROM v2_units WHERE user_id = u.id
            );
        
        -- ====================================================================
        -- Step 4: Mark users as 'canceled' if they have auction records
        -- ====================================================================
        UPDATE deleted_users_info dui
        INNER JOIN users u ON dui.user_id = u.id
        SET dui.deletion_status = 'canceled',
            dui.note = 'The user_id is having an auction, so cannot delete'
        WHERE dui.batch_id = v_batch_id
            AND dui.deletion_status = 'pending'
            AND EXISTS (
                SELECT 1 FROM auction WHERE user_id = u.id
            );
        
        -- ====================================================================
        -- Step 5: Mark users as 'canceled' if they have prelet records
        -- ====================================================================
        UPDATE deleted_users_info dui
        INNER JOIN users u ON dui.user_id = u.id
        SET dui.deletion_status = 'canceled',
            dui.note = 'The user_id is having a prelet, so cannot delete'
        WHERE dui.batch_id = v_batch_id
            AND dui.deletion_status = 'pending'
            AND EXISTS (
                SELECT 1 FROM prelets WHERE user_id = u.id
            );
        
        -- ====================================================================
        -- Step 6: Mark eligible users as 'in progress'
        -- ====================================================================
        UPDATE deleted_users_info dui
        INNER JOIN users u ON dui.user_id = u.id
        INNER JOIN users_roles ur ON u.id = ur.user_id
        INNER JOIN roles r ON ur.role_id = r.id
        SET dui.deletion_status = 'in progress'
        WHERE dui.batch_id = v_batch_id
            AND dui.deletion_status = 'pending'
            AND u.last_login_date < DATE_SUB(CURRENT_TIMESTAMP(), INTERVAL 1 YEAR)
            AND r.name IN ('Tenant', 'Tenant 24/7');
        
        COMMIT;
        
        -- ====================================================================
        -- Step 7: Create and populate temporary table for current batch
        -- ====================================================================
        DROP TEMPORARY TABLE IF EXISTS temp_user_ids;
        
        CREATE TEMPORARY TABLE temp_user_ids (
            id INT NOT NULL PRIMARY KEY,
            user_id INT NOT NULL
        );
        
        INSERT INTO temp_user_ids (id, user_id)
        SELECT id, user_id
        FROM deleted_users_info
        WHERE batch_id = v_batch_id
            AND deletion_status = 'in progress';
        
        -- Get count of users in current batch
        SELECT COUNT(*) INTO num_deleted FROM temp_user_ids;
        
        -- Exit if no eligible users found
        IF num_deleted = 0 THEN
            LEAVE Outer_check;
        END IF;
        
        -- ====================================================================
        -- Step 8: Process related tables - Archive and Delete
        -- ====================================================================
        
        -- Table 1: unit_status_updates
        SET v_tab_name = 'unit_status_updates';
        IF EXISTS (SELECT 1 FROM unit_status_updates usu INNER JOIN temp_user_ids t ON usu.user_id = t.user_id LIMIT 1) THEN
            -- Archive records
            INSERT INTO archived_data.unit_status_updates_archived
            SELECT usu.* FROM unit_status_updates usu
            INNER JOIN temp_user_ids t ON usu.user_id = t.user_id;
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
            
            -- Delete from source
            DELETE usu FROM unit_status_updates usu
            INNER JOIN temp_user_ids t ON usu.user_id = t.user_id
            WHERE EXISTS (
                SELECT 1 FROM archived_data.unit_status_updates_archived a 
                WHERE a.user_id = usu.user_id
            );
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            
            COMMIT;
        END IF;
        
        -- Table 2: unit_overrides
        SET v_tab_name = 'unit_overrides';
        IF EXISTS (SELECT 1 FROM unit_overrides uo INNER JOIN temp_user_ids t ON uo.user_id = t.user_id LIMIT 1) THEN
            INSERT INTO archived_data.unit_overrides_archived
            SELECT uo.* FROM unit_overrides uo
            INNER JOIN temp_user_ids t ON uo.user_id = t.user_id;
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
            
            DELETE uo FROM unit_overrides uo
            INNER JOIN temp_user_ids t ON uo.user_id = t.user_id
            WHERE EXISTS (SELECT 1 FROM archived_data.unit_overrides_archived a WHERE a.user_id = uo.user_id);
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            
            COMMIT;
        END IF;
        
        -- Table 3: transfer
        SET v_tab_name = 'transfer';
        IF EXISTS (SELECT 1 FROM transfer tr INNER JOIN temp_user_ids t ON tr.user_id = t.user_id LIMIT 1) THEN
            INSERT INTO archived_data.transfer_archived
            SELECT tr.* FROM transfer tr
            INNER JOIN temp_user_ids t ON tr.user_id = t.user_id;
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
            
            DELETE tr FROM transfer tr
            INNER JOIN temp_user_ids t ON tr.user_id = t.user_id
            WHERE EXISTS (SELECT 1 FROM archived_data.transfer_archived a WHERE a.user_id = tr.user_id);
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            
            COMMIT;
        END IF;
        
        -- Table 4: support_activity
        SET v_tab_name = 'support_activity';
        IF EXISTS (SELECT 1 FROM support_activity sa INNER JOIN temp_user_ids t ON sa.user_id = t.user_id LIMIT 1) THEN
            INSERT INTO archived_data.support_activity_archived
            SELECT sa.* FROM support_activity sa
            INNER JOIN temp_user_ids t ON sa.user_id = t.user_id;
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
            
            DELETE sa FROM support_activity sa
            INNER JOIN temp_user_ids t ON sa.user_id = t.user_id
            WHERE EXISTS (SELECT 1 FROM archived_data.support_activity_archived a WHERE a.user_id = sa.user_id);
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            
            COMMIT;
        END IF;
        
        -- Table 5: unlock_override_pins
        SET v_tab_name = 'unlock_override_pins';
        IF EXISTS (SELECT 1 FROM unlock_override_pins uop INNER JOIN temp_user_ids t ON uop.user_id = t.user_id LIMIT 1) THEN
            INSERT INTO archived_data.unlock_override_pins_archived
            SELECT uop.* FROM unlock_override_pins uop
            INNER JOIN temp_user_ids t ON uop.user_id = t.user_id;
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
            
            DELETE uop FROM unlock_override_pins uop
            INNER JOIN temp_user_ids t ON uop.user_id = t.user_id
            WHERE EXISTS (SELECT 1 FROM archived_data.unlock_override_pins_archived a WHERE a.user_id = uop.user_id);
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            
            COMMIT;
        END IF;
        
        -- Table 6: user_acknowledgement
        SET v_tab_name = 'user_acknowledgement';
        IF EXISTS (SELECT 1 FROM user_acknowledgement ua INNER JOIN temp_user_ids t ON ua.user_id = t.user_id LIMIT 1) THEN
            INSERT INTO archived_data.user_acknowledgement_archived
            SELECT ua.* FROM user_acknowledgement ua
            INNER JOIN temp_user_ids t ON ua.user_id = t.user_id;
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
            
            DELETE ua FROM user_acknowledgement ua
            INNER JOIN temp_user_ids t ON ua.user_id = t.user_id
            WHERE EXISTS (SELECT 1 FROM archived_data.user_acknowledgement_archived a WHERE a.user_id = ua.user_id);
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            
            COMMIT;
        END IF;
        
        -- Table 7: user_login_history
        SET v_tab_name = 'user_login_history';
        IF EXISTS (SELECT 1 FROM user_login_history ulh INNER JOIN temp_user_ids t ON ulh.user_id = t.user_id LIMIT 1) THEN
            INSERT INTO archived_data.user_login_history_archived
            SELECT ulh.* FROM user_login_history ulh
            INNER JOIN temp_user_ids t ON ulh.user_id = t.user_id;
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
            
            DELETE ulh FROM user_login_history ulh
            INNER JOIN temp_user_ids t ON ulh.user_id = t.user_id
            WHERE EXISTS (SELECT 1 FROM archived_data.user_login_history_archived a WHERE a.user_id = ulh.user_id);
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            
            COMMIT;
        END IF;
        
        -- Table 8: user_site_bookmarks
        SET v_tab_name = 'user_site_bookmarks';
        IF EXISTS (SELECT 1 FROM user_site_bookmarks usb INNER JOIN temp_user_ids t ON usb.user_id = t.user_id LIMIT 1) THEN
            INSERT INTO archived_data.user_site_bookmarks_archived
            SELECT usb.* FROM user_site_bookmarks usb
            INNER JOIN temp_user_ids t ON usb.user_id = t.user_id;
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
            
            DELETE usb FROM user_site_bookmarks usb
            INNER JOIN temp_user_ids t ON usb.user_id = t.user_id
            WHERE EXISTS (SELECT 1 FROM archived_data.user_site_bookmarks_archived a WHERE a.user_id = usb.user_id);
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            
            COMMIT;
        END IF;
        
        -- Table 9: users_dashboard
        SET v_tab_name = 'users_dashboard';
        IF EXISTS (SELECT 1 FROM users_dashboard ud INNER JOIN temp_user_ids t ON ud.user_id = t.user_id LIMIT 1) THEN
            INSERT INTO archived_data.users_dashboard_archived
            SELECT ud.* FROM users_dashboard ud
            INNER JOIN temp_user_ids t ON ud.user_id = t.user_id;
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
            
            DELETE ud FROM users_dashboard ud
            INNER JOIN temp_user_ids t ON ud.user_id = t.user_id
            WHERE EXISTS (SELECT 1 FROM archived_data.users_dashboard_archived a WHERE a.user_id = ud.user_id);
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            
            COMMIT;
        END IF;
        
        -- Table 10: users_notifications_settings
        SET v_tab_name = 'users_notifications_settings';
        IF EXISTS (SELECT 1 FROM users_notifications_settings uns INNER JOIN temp_user_ids t ON uns.user_id = t.user_id LIMIT 1) THEN
            INSERT INTO archived_data.users_notifications_settings_archived
            SELECT uns.* FROM users_notifications_settings uns
            INNER JOIN temp_user_ids t ON uns.user_id = t.user_id;
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
            
            DELETE uns FROM users_notifications_settings uns
            INNER JOIN temp_user_ids t ON uns.user_id = t.user_id
            WHERE EXISTS (SELECT 1 FROM archived_data.users_notifications_settings_archived a WHERE a.user_id = uns.user_id);
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            
            COMMIT;
        END IF;
        
        -- Table 11: users_zone_triggers
        SET v_tab_name = 'users_zone_triggers';
        IF EXISTS (SELECT 1 FROM users_zone_triggers uzt INNER JOIN temp_user_ids t ON uzt.user_id = t.user_id LIMIT 1) THEN
            INSERT INTO archived_data.users_zone_triggers_archived
            SELECT uzt.* FROM users_zone_triggers uzt
            INNER JOIN temp_user_ids t ON uzt.user_id = t.user_id;
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
            
            DELETE uzt FROM users_zone_triggers uzt
            INNER JOIN temp_user_ids t ON uzt.user_id = t.user_id
            WHERE EXISTS (SELECT 1 FROM archived_data.users_zone_triggers_archived a WHERE a.user_id = uzt.user_id);
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            
            COMMIT;
        END IF;
        
        -- Table 12: v2_app_details
        SET v_tab_name = 'v2_app_details';
        IF EXISTS (SELECT 1 FROM v2_app_details vad INNER JOIN temp_user_ids t ON vad.user_id = t.user_id LIMIT 1) THEN
            INSERT INTO archived_data.v2_app_details_archived
            SELECT vad.* FROM v2_app_details vad
            INNER JOIN temp_user_ids t ON vad.user_id = t.user_id;
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
            
            DELETE vad FROM v2_app_details vad
            INNER JOIN temp_user_ids t ON vad.user_id = t.user_id
            WHERE EXISTS (SELECT 1 FROM archived_data.v2_app_details_archived a WHERE a.user_id = vad.user_id);
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            
            COMMIT;
        END IF;
        
        -- Table 13: v2_tracking_ids
        SET v_tab_name = 'v2_tracking_ids';
        IF EXISTS (SELECT 1 FROM v2_tracking_ids vti INNER JOIN temp_user_ids t ON vti.user_id = t.user_id LIMIT 1) THEN
            INSERT INTO archived_data.v2_tracking_ids_archived
            SELECT vti.* FROM v2_tracking_ids vti
            INNER JOIN temp_user_ids t ON vti.user_id = t.user_id;
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
            
            DELETE vti FROM v2_tracking_ids vti
            INNER JOIN temp_user_ids t ON vti.user_id = t.user_id
            WHERE EXISTS (SELECT 1 FROM archived_data.v2_tracking_ids_archived a WHERE a.user_id = vti.user_id);
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            
            COMMIT;
        END IF;
        
        -- Table 14: watch_users
        SET v_tab_name = 'watch_users';
        IF EXISTS (SELECT 1 FROM watch_users wu INNER JOIN temp_user_ids t ON wu.user_id = t.user_id LIMIT 1) THEN
            INSERT INTO archived_data.watch_users_archived
            SELECT wu.* FROM watch_users wu
            INNER JOIN temp_user_ids t ON wu.user_id = t.user_id;
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
            
            DELETE wu FROM watch_users wu
            INNER JOIN temp_user_ids t ON wu.user_id = t.user_id
            WHERE EXISTS (SELECT 1 FROM archived_data.watch_users_archived a WHERE a.user_id = wu.user_id);
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            
            COMMIT;
        END IF;
        
        -- Table 15: 2_factor_auth_pin
        SET v_tab_name = '2_factor_auth_pin';
        IF EXISTS (SELECT 1 FROM 2_factor_auth_pin fap INNER JOIN temp_user_ids t ON fap.user_id = t.user_id LIMIT 1) THEN
            INSERT INTO archived_data.2_factor_auth_pin_archived
            SELECT fap.* FROM 2_factor_auth_pin fap
            INNER JOIN temp_user_ids t ON fap.user_id = t.user_id;
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
            
            DELETE fap FROM 2_factor_auth_pin fap
            INNER JOIN temp_user_ids t ON fap.user_id = t.user_id
            WHERE EXISTS (SELECT 1 FROM archived_data.2_factor_auth_pin_archived a WHERE a.user_id = fap.user_id);
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            
            COMMIT;
        END IF;
        
        -- Table 16: access_codes
        SET v_tab_name = 'access_codes';
        IF EXISTS (SELECT 1 FROM access_codes ac INNER JOIN temp_user_ids t ON ac.user_id = t.user_id LIMIT 1) THEN
            INSERT INTO archived_data.access_codes_archived
            SELECT ac.* FROM access_codes ac
            INNER JOIN temp_user_ids t ON ac.user_id = t.user_id;
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
            
            DELETE ac FROM access_codes ac
            INNER JOIN temp_user_ids t ON ac.user_id = t.user_id
            WHERE EXISTS (SELECT 1 FROM archived_data.access_codes_archived a WHERE a.user_id = ac.user_id);
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            
            COMMIT;
        END IF;
        
        -- Table 17: devices
        SET v_tab_name = 'devices';
        IF EXISTS (SELECT 1 FROM devices d INNER JOIN temp_user_ids t ON d.user_id = t.user_id LIMIT 1) THEN
            INSERT INTO archived_data.devices_archived
            SELECT d.* FROM devices d
            INNER JOIN temp_user_ids t ON d.user_id = t.user_id;
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
            
            DELETE d FROM devices d
            INNER JOIN temp_user_ids t ON d.user_id = t.user_id
            WHERE EXISTS (SELECT 1 FROM archived_data.devices_archived a WHERE a.user_id = d.user_id);
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            
            COMMIT;
        END IF;
        
        -- Table 18: digital_audits
        SET v_tab_name = 'digital_audits';
        IF EXISTS (SELECT 1 FROM digital_audits da INNER JOIN temp_user_ids t ON da.user_id = t.user_id LIMIT 1) THEN
            INSERT INTO archived_data.digital_audits_archived
            SELECT da.* FROM digital_audits da
            INNER JOIN temp_user_ids t ON da.user_id = t.user_id;
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
            
            DELETE da FROM digital_audits da
            INNER JOIN temp_user_ids t ON da.user_id = t.user_id
            WHERE EXISTS (SELECT 1 FROM archived_data.digital_audits_archived a WHERE a.user_id = da.user_id);
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            
            COMMIT;
        END IF;
        
        -- Table 19: entry_activity
        SET v_tab_name = 'entry_activity';
        IF EXISTS (SELECT 1 FROM entry_activity ea INNER JOIN temp_user_ids t ON ea.user_id = t.user_id LIMIT 1) THEN
            INSERT INTO archived_data.entry_activity_archived
            SELECT ea.* FROM entry_activity ea
            INNER JOIN temp_user_ids t ON ea.user_id = t.user_id;
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
            
            DELETE ea FROM entry_activity ea
            INNER JOIN temp_user_ids t ON ea.user_id = t.user_id
            WHERE EXISTS (SELECT 1 FROM archived_data.entry_activity_archived a WHERE a.user_id = ea.user_id);
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            
            COMMIT;
        END IF;
        
        -- Table 20: permissions
        SET v_tab_name = 'permissions';
        IF EXISTS (SELECT 1 FROM permissions p INNER JOIN temp_user_ids t ON p.user_id = t.user_id LIMIT 1) THEN
            INSERT INTO archived_data.permissions_archived
            SELECT p.* FROM permissions p
            INNER JOIN temp_user_ids t ON p.user_id = t.user_id;
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
            
            DELETE p FROM permissions p
            INNER JOIN temp_user_ids t ON p.user_id = t.user_id
            WHERE EXISTS (SELECT 1 FROM archived_data.permissions_archived a WHERE a.user_id = p.user_id);
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            
            COMMIT;
        END IF;
        
        -- Table 21: pending_notifications
        SET v_tab_name = 'pending_notifications';
        IF EXISTS (SELECT 1 FROM pending_notifications pn INNER JOIN temp_user_ids t ON pn.user_id = t.user_id LIMIT 1) THEN
            INSERT INTO archived_data.pending_notifications_archived
            SELECT pn.* FROM pending_notifications pn
            INNER JOIN temp_user_ids t ON pn.user_id = t.user_id;
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
            
            DELETE pn FROM pending_notifications pn
            INNER JOIN temp_user_ids t ON pn.user_id = t.user_id
            WHERE EXISTS (SELECT 1 FROM archived_data.pending_notifications_archived a WHERE a.user_id = pn.user_id);
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            
            COMMIT;
        END IF;
        
        -- Table 22: oauth_session_store
        SET v_tab_name = 'oauth_session_store';
        IF EXISTS (SELECT 1 FROM oauth_session_store oss INNER JOIN temp_user_ids t ON oss.user_id = t.user_id LIMIT 1) THEN
            INSERT INTO archived_data.oauth_session_store_archived
            SELECT oss.* FROM oauth_session_store oss
            INNER JOIN temp_user_ids t ON oss.user_id = t.user_id;
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
            
            DELETE oss FROM oauth_session_store oss
            INNER JOIN temp_user_ids t ON oss.user_id = t.user_id
            WHERE EXISTS (SELECT 1 FROM archived_data.oauth_session_store_archived a WHERE a.user_id = oss.user_id);
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            
            COMMIT;
        END IF;
        
        -- Table 23: note_comments
        SET v_tab_name = 'note_comments';
        IF EXISTS (SELECT 1 FROM note_comments nc INNER JOIN temp_user_ids t ON nc.user_id = t.user_id LIMIT 1) THEN
            INSERT INTO archived_data.note_comments_archived
            SELECT nc.* FROM note_comments nc
            INNER JOIN temp_user_ids t ON nc.user_id = t.user_id;
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
            
            DELETE nc FROM note_comments nc
            INNER JOIN temp_user_ids t ON nc.user_id = t.user_id
            WHERE EXISTS (SELECT 1 FROM archived_data.note_comments_archived a WHERE a.user_id = nc.user_id);
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            
            COMMIT;
        END IF;
        
        -- Table 24: invalid_entry_attempts
        SET v_tab_name = 'invalid_entry_attempts';
        IF EXISTS (SELECT 1 FROM invalid_entry_attempts iea INNER JOIN temp_user_ids t ON iea.user_id = t.user_id LIMIT 1) THEN
            INSERT INTO archived_data.invalid_entry_attempts_archived
            SELECT iea.* FROM invalid_entry_attempts iea
            INNER JOIN temp_user_ids t ON iea.user_id = t.user_id;
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
            
            DELETE iea FROM invalid_entry_attempts iea
            INNER JOIN temp_user_ids t ON iea.user_id = t.user_id
            WHERE EXISTS (SELECT 1 FROM archived_data.invalid_entry_attempts_archived a WHERE a.user_id = iea.user_id);
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            
            COMMIT;
        END IF;
        
        -- Table 25: users_roles
        SET v_tab_name = 'users_roles';
        IF EXISTS (SELECT 1 FROM users_roles ur INNER JOIN temp_user_ids t ON ur.user_id = t.user_id LIMIT 1) THEN
            INSERT INTO archived_data.users_roles_archived
            SELECT ur.* FROM users_roles ur
            INNER JOIN temp_user_ids t ON ur.user_id = t.user_id;
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
            
            DELETE ur FROM users_roles ur
            INNER JOIN temp_user_ids t ON ur.user_id = t.user_id
            WHERE EXISTS (SELECT 1 FROM archived_data.users_roles_archived a WHERE a.user_id = ur.user_id);
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            
            COMMIT;
        END IF;
        
        -- ====================================================================
        -- Step 9: Delete from users table (main table - deleted last)
        -- ====================================================================
        SET v_tab_name = 'users';
        IF EXISTS (SELECT 1 FROM users u INNER JOIN temp_user_ids t ON u.id = t.user_id LIMIT 1) THEN
            INSERT INTO archived_data.users_archived
            SELECT u.* FROM users u
            INNER JOIN temp_user_ids t ON u.id = t.user_id;
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
            
            DELETE u FROM users u
            INNER JOIN temp_user_ids t ON u.id = t.user_id
            WHERE EXISTS (SELECT 1 FROM archived_data.users_archived a WHERE a.id = u.id);
            
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
            VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
            
            COMMIT;
        END IF;
        
        -- ====================================================================
        -- Step 10: Update deleted_users_info status to 'completed'
        -- ====================================================================
        UPDATE deleted_users_info
        SET deletion_status = 'completed',
            note = 'the user_id is deleted'
        WHERE id IN (SELECT id FROM temp_user_ids);
        
        COMMIT;
        
        -- ====================================================================
        -- Step 11: Clean up temporary table
        -- ====================================================================
        DROP TEMPORARY TABLE IF EXISTS temp_user_ids;
        
    END WHILE Outer_check;
    
    -- ========================================================================
    -- Return completion message
    -- ========================================================================
    SELECT 'User deletion process completed.' AS message;
    
END$$

DELIMITER ;

