-- Stored Procedure: delete_users_in_batches
-- Description: Deletes users in batches, archives their data, and logs the process.
-- Author: Cascade

DELIMITER $$

DROP PROCEDURE IF EXISTS `delete_users_in_batches`$$

CREATE PROCEDURE `delete_users_in_batches`()
BEGIN
    -- Variable declarations
    DECLARE done INT DEFAULT FALSE;
    DECLARE num_deleted INT;
    DECLARE v_batch_id INT;
    DECLARE v_note VARCHAR(255);
    DECLARE v_error_info VARCHAR(300);
    DECLARE v_rec_cnt INT DEFAULT 0;
    DECLARE v_tab_name VARCHAR(50);
    -- TODO: Set v_user based on the environment (Prod: 2920483, Pre-prod: 2387943, Beta: 4365, Alpha: 1033278)
    DECLARE v_user INT DEFAULT 2920483;

    -- Error handler
    DECLARE CONTINUE HANDLER FOR SQLEXCEPTION
    BEGIN
        GET DIAGNOSTICS CONDITION 1
        @sqlstate = RETURNED_SQLSTATE, @errno = MYSQL_ERRNO, @text = MESSAGE_TEXT;
        SET v_error_info = CONCAT('ERROR ', @errno, ' (', @sqlstate, '): ', @text);
        ROLLBACK;
        INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `error_info`, `last_updated_by`)
        VALUES (v_batch_id, v_tab_name, 'failed', v_error_info, v_user);
        COMMIT;
        SET done = TRUE;
    END;

    -- Main processing loop
    Outer_check: LOOP
        -- Batch ID Management
        SELECT IFNULL(MAX(batch_id), 0) + 1 INTO v_batch_id FROM `deleted_users_info`;

        UPDATE `deleted_users_info` SET `batch_id` = v_batch_id WHERE `batch_id` IS NULL AND `deletion_status` = 'pending';
        COMMIT;

        -- Mark ineligible users as 'canceled'
        -- Ineligible due to v2_shares (owner)
        UPDATE `deleted_users_info` dui
        JOIN `v2_shares` vs ON dui.user_id = vs.owner_id
        SET dui.deletion_status = 'canceled', dui.note = 'The user_id is either an Owner or Shared user, so cannot delete'
        WHERE dui.batch_id = v_batch_id AND dui.deletion_status = 'pending';

        -- Ineligible due to v2_shares (recipient)
        UPDATE `deleted_users_info` dui
        JOIN `v2_shares` vs ON dui.user_id = vs.recipient_id
        SET dui.deletion_status = 'canceled', dui.note = 'The user_id is either an Owner or Shared user, so cannot delete'
        WHERE dui.batch_id = v_batch_id AND dui.deletion_status = 'pending';

        -- Ineligible due to v2_units
        UPDATE `deleted_users_info` dui
        JOIN `v2_units` vu ON dui.user_id = vu.user_id
        SET dui.deletion_status = 'canceled', dui.note = 'The user_id is having a unit, so cannot delete'
        WHERE dui.batch_id = v_batch_id AND dui.deletion_status = 'pending';

        -- Ineligible due to auction
        UPDATE `deleted_users_info` dui
        JOIN `auction` a ON dui.user_id = a.user_id
        SET dui.deletion_status = 'canceled', dui.note = 'The user_id is having an auction, so cannot delete'
        WHERE dui.batch_id = v_batch_id AND dui.deletion_status = 'pending';

        -- Ineligible due to prelets
        UPDATE `deleted_users_info` dui
        JOIN `prelets` p ON dui.user_id = p.user_id
        SET dui.deletion_status = 'canceled', dui.note = 'The user_id is having a prelet, so cannot delete'
        WHERE dui.batch_id = v_batch_id AND dui.deletion_status = 'pending';

        -- Ineligible due to last login date or role
        UPDATE `deleted_users_info` dui
        JOIN `users` u ON dui.user_id = u.id
        LEFT JOIN `users_roles` ur ON u.id = ur.user_id
        LEFT JOIN `roles` r ON ur.role_id = r.id
        SET dui.deletion_status = 'canceled', dui.note = 'User does not meet login or role criteria for deletion'
        WHERE dui.batch_id = v_batch_id
          AND dui.deletion_status = 'pending'
          AND (u.last_login_date >= DATE_SUB(CURRENT_TIMESTAMP(), INTERVAL 1 YEAR) OR r.name NOT IN ('Tenant', 'Tenant 24/7'));

        -- Mark eligible users as 'in progress'
        UPDATE `deleted_users_info`
        SET `deletion_status` = 'in progress'
        WHERE `batch_id` = v_batch_id AND `deletion_status` = 'pending';
        COMMIT;

        -- Temporary table for the current batch
        DROP TEMPORARY TABLE IF EXISTS `temp_user_ids`;
        CREATE TEMPORARY TABLE `temp_user_ids` (
            `id` INT NOT NULL PRIMARY KEY,
            `user_id` INT NOT NULL
        );

        INSERT INTO `temp_user_ids` (`id`, `user_id`)
        SELECT `id`, `user_id` FROM `deleted_users_info` WHERE `batch_id` = v_batch_id AND `deletion_status` = 'in progress';

        SELECT COUNT(*) INTO num_deleted FROM `temp_user_ids`;

        IF num_deleted = 0 THEN
            LEAVE Outer_check;
        END IF;

        -- Deletion process for each table
        -- Deletion process for each table

        SET v_tab_name = 'unit_status_updates';
        IF EXISTS (SELECT 1 FROM `unit_status_updates` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id) THEN
            START TRANSACTION;
            INSERT INTO `archived_data`.`unit_status_updates_archived` SELECT t.* FROM `unit_status_updates` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id;
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), v_user);
            DELETE t FROM `unit_status_updates` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id WHERE EXISTS (SELECT 1 FROM `archived_data`.`unit_status_updates_archived` a WHERE a.id = t.id);
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), v_user);
            COMMIT;
        END IF;

        SET v_tab_name = 'unit_overrides';
        IF EXISTS (SELECT 1 FROM `unit_overrides` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id) THEN
            START TRANSACTION;
            INSERT INTO `archived_data`.`unit_overrides_archived` SELECT t.* FROM `unit_overrides` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id;
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), v_user);
            DELETE t FROM `unit_overrides` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id WHERE EXISTS (SELECT 1 FROM `archived_data`.`unit_overrides_archived` a WHERE a.id = t.id);
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), v_user);
            COMMIT;
        END IF;

        SET v_tab_name = 'transfer';
        IF EXISTS (SELECT 1 FROM `transfer` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id) THEN
            START TRANSACTION;
            INSERT INTO `archived_data`.`transfer_archived` SELECT t.* FROM `transfer` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id;
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), v_user);
            DELETE t FROM `transfer` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id WHERE EXISTS (SELECT 1 FROM `archived_data`.`transfer_archived` a WHERE a.id = t.id);
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), v_user);
            COMMIT;
        END IF;

        SET v_tab_name = 'support_activity';
        IF EXISTS (SELECT 1 FROM `support_activity` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id) THEN
            START TRANSACTION;
            INSERT INTO `archived_data`.`support_activity_archived` SELECT t.* FROM `support_activity` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id;
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), v_user);
            DELETE t FROM `support_activity` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id WHERE EXISTS (SELECT 1 FROM `archived_data`.`support_activity_archived` a WHERE a.id = t.id);
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), v_user);
            COMMIT;
        END IF;

        SET v_tab_name = 'unlock_override_pins';
        IF EXISTS (SELECT 1 FROM `unlock_override_pins` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id) THEN
            START TRANSACTION;
            INSERT INTO `archived_data`.`unlock_override_pins_archived` SELECT t.* FROM `unlock_override_pins` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id;
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), v_user);
            DELETE t FROM `unlock_override_pins` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id WHERE EXISTS (SELECT 1 FROM `archived_data`.`unlock_override_pins_archived` a WHERE a.id = t.id);
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), v_user);
            COMMIT;
        END IF;

        SET v_tab_name = 'user_acknowledgement';
        IF EXISTS (SELECT 1 FROM `user_acknowledgement` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id) THEN
            START TRANSACTION;
            INSERT INTO `archived_data`.`user_acknowledgement_archived` SELECT t.* FROM `user_acknowledgement` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id;
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), v_user);
            DELETE t FROM `user_acknowledgement` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id WHERE EXISTS (SELECT 1 FROM `archived_data`.`user_acknowledgement_archived` a WHERE a.id = t.id);
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), v_user);
            COMMIT;
        END IF;

        SET v_tab_name = 'user_login_history';
        IF EXISTS (SELECT 1 FROM `user_login_history` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id) THEN
            START TRANSACTION;
            INSERT INTO `archived_data`.`user_login_history_archived` SELECT t.* FROM `user_login_history` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id;
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), v_user);
            DELETE t FROM `user_login_history` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id WHERE EXISTS (SELECT 1 FROM `archived_data`.`user_login_history_archived` a WHERE a.id = t.id);
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), v_user);
            COMMIT;
        END IF;

        SET v_tab_name = 'user_site_bookmarks';
        IF EXISTS (SELECT 1 FROM `user_site_bookmarks` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id) THEN
            START TRANSACTION;
            INSERT INTO `archived_data`.`user_site_bookmarks_archived` SELECT t.* FROM `user_site_bookmarks` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id;
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), v_user);
            DELETE t FROM `user_site_bookmarks` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id WHERE EXISTS (SELECT 1 FROM `archived_data`.`user_site_bookmarks_archived` a WHERE a.id = t.id);
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), v_user);
            COMMIT;
        END IF;

        SET v_tab_name = 'users_dashboard';
        IF EXISTS (SELECT 1 FROM `users_dashboard` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id) THEN
            START TRANSACTION;
            INSERT INTO `archived_data`.`users_dashboard_archived` SELECT t.* FROM `users_dashboard` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id;
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), v_user);
            DELETE t FROM `users_dashboard` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id WHERE EXISTS (SELECT 1 FROM `archived_data`.`users_dashboard_archived` a WHERE a.id = t.id);
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), v_user);
            COMMIT;
        END IF;

        SET v_tab_name = 'users_notifications_settings';
        IF EXISTS (SELECT 1 FROM `users_notifications_settings` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id) THEN
            START TRANSACTION;
            INSERT INTO `archived_data`.`users_notifications_settings_archived` SELECT t.* FROM `users_notifications_settings` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id;
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), v_user);
            DELETE t FROM `users_notifications_settings` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id WHERE EXISTS (SELECT 1 FROM `archived_data`.`users_notifications_settings_archived` a WHERE a.id = t.id);
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), v_user);
            COMMIT;
        END IF;

        SET v_tab_name = 'users_zone_triggers';
        IF EXISTS (SELECT 1 FROM `users_zone_triggers` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id) THEN
            START TRANSACTION;
            INSERT INTO `archived_data`.`users_zone_triggers_archived` SELECT t.* FROM `users_zone_triggers` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id;
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), v_user);
            DELETE t FROM `users_zone_triggers` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id WHERE EXISTS (SELECT 1 FROM `archived_data`.`users_zone_triggers_archived` a WHERE a.id = t.id);
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), v_user);
            COMMIT;
        END IF;

        SET v_tab_name = 'v2_app_details';
        IF EXISTS (SELECT 1 FROM `v2_app_details` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id) THEN
            START TRANSACTION;
            INSERT INTO `archived_data`.`v2_app_details_archived` SELECT t.* FROM `v2_app_details` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id;
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), v_user);
            DELETE t FROM `v2_app_details` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id WHERE EXISTS (SELECT 1 FROM `archived_data`.`v2_app_details_archived` a WHERE a.id = t.id);
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), v_user);
            COMMIT;
        END IF;

        SET v_tab_name = 'v2_tracking_ids';
        IF EXISTS (SELECT 1 FROM `v2_tracking_ids` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id) THEN
            START TRANSACTION;
            INSERT INTO `archived_data`.`v2_tracking_ids_archived` SELECT t.* FROM `v2_tracking_ids` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id;
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), v_user);
            DELETE t FROM `v2_tracking_ids` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id WHERE EXISTS (SELECT 1 FROM `archived_data`.`v2_tracking_ids_archived` a WHERE a.id = t.id);
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), v_user);
            COMMIT;
        END IF;

        SET v_tab_name = 'watch_users';
        IF EXISTS (SELECT 1 FROM `watch_users` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id) THEN
            START TRANSACTION;
            INSERT INTO `archived_data`.`watch_users_archived` SELECT t.* FROM `watch_users` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id;
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), v_user);
            DELETE t FROM `watch_users` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id WHERE EXISTS (SELECT 1 FROM `archived_data`.`watch_users_archived` a WHERE a.id = t.id);
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), v_user);
            COMMIT;
        END IF;

        SET v_tab_name = '2_factor_auth_pin';
        IF EXISTS (SELECT 1 FROM `2_factor_auth_pin` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id) THEN
            START TRANSACTION;
            INSERT INTO `archived_data`.`2_factor_auth_pin_archived` SELECT t.* FROM `2_factor_auth_pin` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id;
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), v_user);
            DELETE t FROM `2_factor_auth_pin` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id WHERE EXISTS (SELECT 1 FROM `archived_data`.`2_factor_auth_pin_archived` a WHERE a.id = t.id);
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), v_user);
            COMMIT;
        END IF;

        SET v_tab_name = 'access_codes';
        IF EXISTS (SELECT 1 FROM `access_codes` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id) THEN
            START TRANSACTION;
            INSERT INTO `archived_data`.`access_codes_archived` SELECT t.* FROM `access_codes` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id;
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), v_user);
            DELETE t FROM `access_codes` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id WHERE EXISTS (SELECT 1 FROM `archived_data`.`access_codes_archived` a WHERE a.id = t.id);
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), v_user);
            COMMIT;
        END IF;

        SET v_tab_name = 'devices';
        IF EXISTS (SELECT 1 FROM `devices` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id) THEN
            START TRANSACTION;
            INSERT INTO `archived_data`.`devices_archived` SELECT t.* FROM `devices` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id;
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), v_user);
            DELETE t FROM `devices` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id WHERE EXISTS (SELECT 1 FROM `archived_data`.`devices_archived` a WHERE a.id = t.id);
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), v_user);
            COMMIT;
        END IF;

        SET v_tab_name = 'digital_audits';
        IF EXISTS (SELECT 1 FROM `digital_audits` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id) THEN
            START TRANSACTION;
            INSERT INTO `archived_data`.`digital_audits_archived` SELECT t.* FROM `digital_audits` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id;
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), v_user);
            DELETE t FROM `digital_audits` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id WHERE EXISTS (SELECT 1 FROM `archived_data`.`digital_audits_archived` a WHERE a.id = t.id);
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), v_user);
            COMMIT;
        END IF;

        SET v_tab_name = 'entry_activity';
        IF EXISTS (SELECT 1 FROM `entry_activity` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id) THEN
            START TRANSACTION;
            INSERT INTO `archived_data`.`entry_activity_archived` SELECT t.* FROM `entry_activity` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id;
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), v_user);
            DELETE t FROM `entry_activity` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id WHERE EXISTS (SELECT 1 FROM `archived_data`.`entry_activity_archived` a WHERE a.id = t.id);
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), v_user);
            COMMIT;
        END IF;

        SET v_tab_name = 'permissions';
        IF EXISTS (SELECT 1 FROM `permissions` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id) THEN
            START TRANSACTION;
            INSERT INTO `archived_data`.`permissions_archived` SELECT t.* FROM `permissions` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id;
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), v_user);
            DELETE t FROM `permissions` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id WHERE EXISTS (SELECT 1 FROM `archived_data`.`permissions_archived` a WHERE a.id = t.id);
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), v_user);
            COMMIT;
        END IF;

        SET v_tab_name = 'pending_notifications';
        IF EXISTS (SELECT 1 FROM `pending_notifications` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id) THEN
            START TRANSACTION;
            INSERT INTO `archived_data`.`pending_notifications_archived` SELECT t.* FROM `pending_notifications` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id;
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), v_user);
            DELETE t FROM `pending_notifications` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id WHERE EXISTS (SELECT 1 FROM `archived_data`.`pending_notifications_archived` a WHERE a.id = t.id);
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), v_user);
            COMMIT;
        END IF;

        SET v_tab_name = 'oauth_session_store';
        IF EXISTS (SELECT 1 FROM `oauth_session_store` t JOIN `temp_user_ids` tmp ON t.owner_id = tmp.user_id) THEN
            START TRANSACTION;
            INSERT INTO `archived_data`.`oauth_session_store_archived` SELECT t.* FROM `oauth_session_store` t JOIN `temp_user_ids` tmp ON t.owner_id = tmp.user_id;
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), v_user);
            DELETE t FROM `oauth_session_store` t JOIN `temp_user_ids` tmp ON t.owner_id = tmp.user_id WHERE EXISTS (SELECT 1 FROM `archived_data`.`oauth_session_store_archived` a WHERE a.id = t.id);
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), v_user);
            COMMIT;
        END IF;

        SET v_tab_name = 'note_comments';
        IF EXISTS (SELECT 1 FROM `note_comments` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id) THEN
            START TRANSACTION;
            INSERT INTO `archived_data`.`note_comments_archived` SELECT t.* FROM `note_comments` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id;
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), v_user);
            DELETE t FROM `note_comments` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id WHERE EXISTS (SELECT 1 FROM `archived_data`.`note_comments_archived` a WHERE a.id = t.id);
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), v_user);
            COMMIT;
        END IF;

        SET v_tab_name = 'invalid_entry_attempts';
        IF EXISTS (SELECT 1 FROM `invalid_entry_attempts` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id) THEN
            START TRANSACTION;
            INSERT INTO `archived_data`.`invalid_entry_attempts_archived` SELECT t.* FROM `invalid_entry_attempts` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id;
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), v_user);
            DELETE t FROM `invalid_entry_attempts` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id WHERE EXISTS (SELECT 1 FROM `archived_data`.`invalid_entry_attempts_archived` a WHERE a.id = t.id);
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), v_user);
            COMMIT;
        END IF;

        SET v_tab_name = 'users_roles';
        IF EXISTS (SELECT 1 FROM `users_roles` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id) THEN
            START TRANSACTION;
            INSERT INTO `archived_data`.`users_roles_archived` SELECT t.* FROM `users_roles` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id;
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), v_user);
            DELETE t FROM `users_roles` t JOIN `temp_user_ids` tmp ON t.user_id = tmp.user_id WHERE EXISTS (SELECT 1 FROM `archived_data`.`users_roles_archived` a WHERE a.user_id = t.user_id AND a.role_id = t.role_id);
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), v_user);
            COMMIT;
        END IF;

        SET v_tab_name = 'users';
        IF EXISTS (SELECT 1 FROM `users` t JOIN `temp_user_ids` tmp ON t.id = tmp.user_id) THEN
            START TRANSACTION;
            INSERT INTO `archived_data`.`users_archived` SELECT t.* FROM `users` t JOIN `temp_user_ids` tmp ON t.id = tmp.user_id;
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), v_user);
            DELETE t FROM `users` t JOIN `temp_user_ids` tmp ON t.id = tmp.user_id WHERE EXISTS (SELECT 1 FROM `archived_data`.`users_archived` a WHERE a.id = t.id);
            SET v_rec_cnt = ROW_COUNT();
            INSERT INTO `user_deletion_detail` (`batch_id`, `table_name`, `note`, `last_updated_by`) VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), v_user);
            COMMIT;
        END IF;

        -- Update status to 'completed'
        UPDATE `deleted_users_info`
        SET `deletion_status` = 'completed', `note` = 'the user_id is deleted'
        WHERE `id` IN (SELECT `id` FROM `temp_user_ids`);
        COMMIT;

        DROP TEMPORARY TABLE IF EXISTS `temp_user_ids`;

    END LOOP Outer_check;

    SELECT 'User deletion process completed.' AS message;

END$$

DELIMITER ;
