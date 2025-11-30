DELIMITER $$

DROP PROCEDURE IF EXISTS delete_users_in_batches $$
CREATE PROCEDURE delete_users_in_batches()
BEGIN
  DECLARE done INT DEFAULT FALSE;
  DECLARE num_deleted INT DEFAULT 0;
  DECLARE v_batch_id INT DEFAULT NULL;
  DECLARE v_note VARCHAR(255);
  DECLARE v_error_info VARCHAR(300);
  DECLARE v_rec_cnt INT DEFAULT 0;
  DECLARE v_tab_name VARCHAR(50);
  DECLARE v_user INT DEFAULT 2920483; -- SET PER ENV: Prod=2920483, Preprod=2387943, Beta=4365, Alpha=1033278

  -- Error handler: log and exit loop
  DECLARE CONTINUE HANDLER FOR SQLEXCEPTION
  BEGIN
    GET DIAGNOSTICS CONDITION 1 @p_errno = MYSQL_ERRNO, @p_sqlstate = SQLSTATE, @p_message = MESSAGE_TEXT;
    SET v_error_info = CONCAT('ERROR ', @p_errno, ' (', @p_sqlstate, '): ', @p_message);
    ROLLBACK;
    INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
    VALUES (v_batch_id, v_tab_name, 'failed', v_error_info, v_user);
    COMMIT;
    SET done = TRUE;
  END;

  Outer_check: WHILE NOT done DO
    -- 1) Determine and assign new batch_id
    SELECT COALESCE(MAX(batch_id), 0) INTO v_batch_id FROM deleted_users_info;
    SET v_batch_id = v_batch_id + 1;

    START TRANSACTION;
    UPDATE deleted_users_info
      SET batch_id = v_batch_id
    WHERE batch_id IS NULL AND deletion_status = 'pending';
    COMMIT;

    -- 2) Update user statuses (cancellations first by reason-specific rules)
    START TRANSACTION;

    -- Cancel: user is owner or recipient in shares
    UPDATE deleted_users_info di
    JOIN v2_shares s1 ON s1.owner_id = di.user_id
    SET di.deletion_status = 'canceled',
        di.note = 'The user_id is either an Owner or Shared user, so cannot delete'
    WHERE di.batch_id = v_batch_id AND di.deletion_status = 'pending';

    UPDATE deleted_users_info di
    JOIN v2_shares s2 ON s2.recipient_id = di.user_id
    SET di.deletion_status = 'canceled',
        di.note = 'The user_id is either an Owner or Shared user, so cannot delete'
    WHERE di.batch_id = v_batch_id AND di.deletion_status = 'pending';

    -- Cancel: user has associated units
    UPDATE deleted_users_info di
    JOIN v2_units vu ON vu.user_id = di.user_id
    SET di.deletion_status = 'canceled',
        di.note = 'The user_id is having a unit, so cannot delete'
    WHERE di.batch_id = v_batch_id AND di.deletion_status = 'pending';

    -- Cancel: user has associated auctions
    UPDATE deleted_users_info di
    JOIN auction a ON a.user_id = di.user_id
    SET di.deletion_status = 'canceled',
        di.note = 'The user_id is having an auction, so cannot delete'
    WHERE di.batch_id = v_batch_id AND di.deletion_status = 'pending';

    -- Cancel: user has associated prelets
    UPDATE deleted_users_info di
    JOIN prelets p ON p.user_id = di.user_id
    SET di.deletion_status = 'canceled',
        di.note = 'The user_id is having a prelet, so cannot delete'
    WHERE di.batch_id = v_batch_id AND di.deletion_status = 'pending';

    -- Mark eligible users as in progress (role + last_login_date + no disqualifiers)
    UPDATE deleted_users_info di
    JOIN users u ON u.id = di.user_id
    JOIN users_roles ur ON ur.user_id = u.id
    JOIN roles r ON r.id = ur.role_id AND r.name IN ('Tenant', 'Tenant 24/7')
    LEFT JOIN v2_shares s_own ON s_own.owner_id = u.id
    LEFT JOIN v2_shares s_rec ON s_rec.recipient_id = u.id
    LEFT JOIN v2_units vuu ON vuu.user_id = u.id
    LEFT JOIN auction aa ON aa.user_id = u.id
    LEFT JOIN prelets pp ON pp.user_id = u.id
    SET di.deletion_status = 'in progress'
    WHERE di.batch_id = v_batch_id
      AND di.deletion_status = 'pending'
      AND u.last_login_date < DATE_SUB(CURRENT_TIMESTAMP(), INTERVAL 1 YEAR)
      AND s_own.owner_id IS NULL
      AND s_rec.recipient_id IS NULL
      AND vuu.user_id IS NULL
      AND aa.user_id IS NULL
      AND pp.user_id IS NULL;

    COMMIT;

    -- 3) Create and populate temporary table for this batch
    DROP TEMPORARY TABLE IF EXISTS temp_user_ids;
    CREATE TEMPORARY TABLE temp_user_ids (
      id INT NOT NULL PRIMARY KEY,
      user_id INT NOT NULL
    ) ENGINE=Memory;

    INSERT INTO temp_user_ids (id, user_id)
    SELECT di.id, di.user_id
    FROM deleted_users_info di
    WHERE di.batch_id = v_batch_id AND di.deletion_status = 'in progress';

    SELECT COUNT(*) INTO num_deleted FROM temp_user_ids;
    IF num_deleted = 0 THEN
      -- No users to process in this batch → exit the loop
      SET done = TRUE;
      LEAVE Outer_check;
    END IF;

    -- Helper macro-like block to archive/delete a related table by user_id
    -- For each table: set v_tab_name, then run the block.

    -- unit_status_updates
    SET v_tab_name = 'unit_status_updates';
    IF EXISTS (SELECT 1 FROM unit_status_updates t JOIN temp_user_ids b ON b.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.unit_status_updates_archived
      SELECT t.* FROM unit_status_updates t JOIN temp_user_ids b ON b.user_id = t.user_id;
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
      DELETE t FROM unit_status_updates t
      JOIN temp_user_ids b ON b.user_id = t.user_id
      WHERE EXISTS (
        SELECT 1 FROM archived_data.unit_status_updates_archived a WHERE a.user_id = t.user_id
      );
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
      COMMIT;
    END IF;

    -- unit_overrides
    SET v_tab_name = 'unit_overrides';
    IF EXISTS (SELECT 1 FROM unit_overrides t JOIN temp_user_ids b ON b.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.unit_overrides_archived
      SELECT t.* FROM unit_overrides t JOIN temp_user_ids b ON b.user_id = t.user_id;
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
      DELETE t FROM unit_overrides t
      JOIN temp_user_ids b ON b.user_id = t.user_id
      WHERE EXISTS (
        SELECT 1 FROM archived_data.unit_overrides_archived a WHERE a.user_id = t.user_id
      );
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
      COMMIT;
    END IF;

    -- transfer
    SET v_tab_name = 'transfer';
    IF EXISTS (SELECT 1 FROM `transfer` t JOIN temp_user_ids b ON b.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.transfer_archived
      SELECT t.* FROM `transfer` t JOIN temp_user_ids b ON b.user_id = t.user_id;
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
      DELETE t FROM `transfer` t
      JOIN temp_user_ids b ON b.user_id = t.user_id
      WHERE EXISTS (
        SELECT 1 FROM archived_data.transfer_archived a WHERE a.user_id = t.user_id
      );
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
      COMMIT;
    END IF;

    -- support_activity
    SET v_tab_name = 'support_activity';
    IF EXISTS (SELECT 1 FROM support_activity t JOIN temp_user_ids b ON b.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.support_activity_archived
      SELECT t.* FROM support_activity t JOIN temp_user_ids b ON b.user_id = t.user_id;
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
      DELETE t FROM support_activity t
      JOIN temp_user_ids b ON b.user_id = t.user_id
      WHERE EXISTS (
        SELECT 1 FROM archived_data.support_activity_archived a WHERE a.user_id = t.user_id
      );
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
      COMMIT;
    END IF;

    -- unlock_override_pins
    SET v_tab_name = 'unlock_override_pins';
    IF EXISTS (SELECT 1 FROM unlock_override_pins t JOIN temp_user_ids b ON b.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.unlock_override_pins_archived
      SELECT t.* FROM unlock_override_pins t JOIN temp_user_ids b ON b.user_id = t.user_id;
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
      DELETE t FROM unlock_override_pins t
      JOIN temp_user_ids b ON b.user_id = t.user_id
      WHERE EXISTS (
        SELECT 1 FROM archived_data.unlock_override_pins_archived a WHERE a.user_id = t.user_id
      );
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
      COMMIT;
    END IF;

    -- user_acknowledgement
    SET v_tab_name = 'user_acknowledgement';
    IF EXISTS (SELECT 1 FROM user_acknowledgement t JOIN temp_user_ids b ON b.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.user_acknowledgement_archived
      SELECT t.* FROM user_acknowledgement t JOIN temp_user_ids b ON b.user_id = t.user_id;
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
      DELETE t FROM user_acknowledgement t
      JOIN temp_user_ids b ON b.user_id = t.user_id
      WHERE EXISTS (
        SELECT 1 FROM archived_data.user_acknowledgement_archived a WHERE a.user_id = t.user_id
      );
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
      COMMIT;
    END IF;

    -- user_login_history
    SET v_tab_name = 'user_login_history';
    IF EXISTS (SELECT 1 FROM user_login_history t JOIN temp_user_ids b ON b.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.user_login_history_archived
      SELECT t.* FROM user_login_history t JOIN temp_user_ids b ON b.user_id = t.user_id;
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
      DELETE t FROM user_login_history t
      JOIN temp_user_ids b ON b.user_id = t.user_id
      WHERE EXISTS (
        SELECT 1 FROM archived_data.user_login_history_archived a WHERE a.user_id = t.user_id
      );
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
      COMMIT;
    END IF;

    -- user_site_bookmarks
    SET v_tab_name = 'user_site_bookmarks';
    IF EXISTS (SELECT 1 FROM user_site_bookmarks t JOIN temp_user_ids b ON b.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.user_site_bookmarks_archived
      SELECT t.* FROM user_site_bookmarks t JOIN temp_user_ids b ON b.user_id = t.user_id;
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
      DELETE t FROM user_site_bookmarks t
      JOIN temp_user_ids b ON b.user_id = t.user_id
      WHERE EXISTS (
        SELECT 1 FROM archived_data.user_site_bookmarks_archived a WHERE a.user_id = t.user_id
      );
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
      COMMIT;
    END IF;

    -- users_dashboard
    SET v_tab_name = 'users_dashboard';
    IF EXISTS (SELECT 1 FROM users_dashboard t JOIN temp_user_ids b ON b.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.users_dashboard_archived
      SELECT t.* FROM users_dashboard t JOIN temp_user_ids b ON b.user_id = t.user_id;
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
      DELETE t FROM users_dashboard t
      JOIN temp_user_ids b ON b.user_id = t.user_id
      WHERE EXISTS (
        SELECT 1 FROM archived_data.users_dashboard_archived a WHERE a.user_id = t.user_id
      );
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
      COMMIT;
    END IF;

    -- users_notifications_settings
    SET v_tab_name = 'users_notifications_settings';
    IF EXISTS (SELECT 1 FROM users_notifications_settings t JOIN temp_user_ids b ON b.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.users_notifications_settings_archived
      SELECT t.* FROM users_notifications_settings t JOIN temp_user_ids b ON b.user_id = t.user_id;
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
      DELETE t FROM users_notifications_settings t
      JOIN temp_user_ids b ON b.user_id = t.user_id
      WHERE EXISTS (
        SELECT 1 FROM archived_data.users_notifications_settings_archived a WHERE a.user_id = t.user_id
      );
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
      COMMIT;
    END IF;

    -- users_zone_triggers
    SET v_tab_name = 'users_zone_triggers';
    IF EXISTS (SELECT 1 FROM users_zone_triggers t JOIN temp_user_ids b ON b.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.users_zone_triggers_archived
      SELECT t.* FROM users_zone_triggers t JOIN temp_user_ids b ON b.user_id = t.user_id;
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
      DELETE t FROM users_zone_triggers t
      JOIN temp_user_ids b ON b.user_id = t.user_id
      WHERE EXISTS (
        SELECT 1 FROM archived_data.users_zone_triggers_archived a WHERE a.user_id = t.user_id
      );
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
      COMMIT;
    END IF;

    -- v2_app_details
    SET v_tab_name = 'v2_app_details';
    IF EXISTS (SELECT 1 FROM v2_app_details t JOIN temp_user_ids b ON b.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.v2_app_details_archived
      SELECT t.* FROM v2_app_details t JOIN temp_user_ids b ON b.user_id = t.user_id;
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
      DELETE t FROM v2_app_details t
      JOIN temp_user_ids b ON b.user_id = t.user_id
      WHERE EXISTS (
        SELECT 1 FROM archived_data.v2_app_details_archived a WHERE a.user_id = t.user_id
      );
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
      COMMIT;
    END IF;

    -- v2_tracking_ids
    SET v_tab_name = 'v2_tracking_ids';
    IF EXISTS (SELECT 1 FROM v2_tracking_ids t JOIN temp_user_ids b ON b.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.v2_tracking_ids_archived
      SELECT t.* FROM v2_tracking_ids t JOIN temp_user_ids b ON b.user_id = t.user_id;
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
      DELETE t FROM v2_tracking_ids t
      JOIN temp_user_ids b ON b.user_id = t.user_id
      WHERE EXISTS (
        SELECT 1 FROM archived_data.v2_tracking_ids_archived a WHERE a.user_id = t.user_id
      );
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
      COMMIT;
    END IF;

    -- watch_users
    SET v_tab_name = 'watch_users';
    IF EXISTS (SELECT 1 FROM watch_users t JOIN temp_user_ids b ON b.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.watch_users_archived
      SELECT t.* FROM watch_users t JOIN temp_user_ids b ON b.user_id = t.user_id;
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
      DELETE t FROM watch_users t
      JOIN temp_user_ids b ON b.user_id = t.user_id
      WHERE EXISTS (
        SELECT 1 FROM archived_data.watch_users_archived a WHERE a.user_id = t.user_id
      );
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
      COMMIT;
    END IF;

    -- 2_factor_auth_pin
    SET v_tab_name = '2_factor_auth_pin';
    IF EXISTS (SELECT 1 FROM `2_factor_auth_pin` t JOIN temp_user_ids b ON b.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.`2_factor_auth_pin_archived`
      SELECT t.* FROM `2_factor_auth_pin` t JOIN temp_user_ids b ON b.user_id = t.user_id;
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
      DELETE t FROM `2_factor_auth_pin` t
      JOIN temp_user_ids b ON b.user_id = t.user_id
      WHERE EXISTS (
        SELECT 1 FROM archived_data.`2_factor_auth_pin_archived` a WHERE a.user_id = t.user_id
      );
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
      COMMIT;
    END IF;

    -- access_codes
    SET v_tab_name = 'access_codes';
    IF EXISTS (SELECT 1 FROM access_codes t JOIN temp_user_ids b ON b.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.access_codes_archived
      SELECT t.* FROM access_codes t JOIN temp_user_ids b ON b.user_id = t.user_id;
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
      DELETE t FROM access_codes t
      JOIN temp_user_ids b ON b.user_id = t.user_id
      WHERE EXISTS (
        SELECT 1 FROM archived_data.access_codes_archived a WHERE a.user_id = t.user_id
      );
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
      COMMIT;
    END IF;

    -- devices
    SET v_tab_name = 'devices';
    IF EXISTS (SELECT 1 FROM devices t JOIN temp_user_ids b ON b.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.devices_archived
      SELECT t.* FROM devices t JOIN temp_user_ids b ON b.user_id = t.user_id;
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
      DELETE t FROM devices t
      JOIN temp_user_ids b ON b.user_id = t.user_id
      WHERE EXISTS (
        SELECT 1 FROM archived_data.devices_archived a WHERE a.user_id = t.user_id
      );
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
      COMMIT;
    END IF;

    -- digital_audits
    SET v_tab_name = 'digital_audits';
    IF EXISTS (SELECT 1 FROM digital_audits t JOIN temp_user_ids b ON b.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.digital_audits_archived
      SELECT t.* FROM digital_audits t JOIN temp_user_ids b ON b.user_id = t.user_id;
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
      DELETE t FROM digital_audits t
      JOIN temp_user_ids b ON b.user_id = t.user_id
      WHERE EXISTS (
        SELECT 1 FROM archived_data.digital_audits_archived a WHERE a.user_id = t.user_id
      );
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
      COMMIT;
    END IF;

    -- entry_activity
    SET v_tab_name = 'entry_activity';
    IF EXISTS (SELECT 1 FROM entry_activity t JOIN temp_user_ids b ON b.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.entry_activity_archived
      SELECT t.* FROM entry_activity t JOIN temp_user_ids b ON b.user_id = t.user_id;
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
      DELETE t FROM entry_activity t
      JOIN temp_user_ids b ON b.user_id = t.user_id
      WHERE EXISTS (
        SELECT 1 FROM archived_data.entry_activity_archived a WHERE a.user_id = t.user_id
      );
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
      COMMIT;
    END IF;

    -- permissions
    SET v_tab_name = 'permissions';
    IF EXISTS (SELECT 1 FROM permissions t JOIN temp_user_ids b ON b.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.permissions_archived
      SELECT t.* FROM permissions t JOIN temp_user_ids b ON b.user_id = t.user_id;
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
      DELETE t FROM permissions t
      JOIN temp_user_ids b ON b.user_id = t.user_id
      WHERE EXISTS (
        SELECT 1 FROM archived_data.permissions_archived a WHERE a.user_id = t.user_id
      );
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
      COMMIT;
    END IF;

    -- pending_notifications
    SET v_tab_name = 'pending_notifications';
    IF EXISTS (SELECT 1 FROM pending_notifications t JOIN temp_user_ids b ON b.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.pending_notifications_archived
      SELECT t.* FROM pending_notifications t JOIN temp_user_ids b ON b.user_id = t.user_id;
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
      DELETE t FROM pending_notifications t
      JOIN temp_user_ids b ON b.user_id = t.user_id
      WHERE EXISTS (
        SELECT 1 FROM archived_data.pending_notifications_archived a WHERE a.user_id = t.user_id
      );
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
      COMMIT;
    END IF;

    -- oauth_session_store
    SET v_tab_name = 'oauth_session_store';
    IF EXISTS (SELECT 1 FROM oauth_session_store t JOIN temp_user_ids b ON b.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.oauth_session_store_archived
      SELECT t.* FROM oauth_session_store t JOIN temp_user_ids b ON b.user_id = t.user_id;
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
      DELETE t FROM oauth_session_store t
      JOIN temp_user_ids b ON b.user_id = t.user_id
      WHERE EXISTS (
        SELECT 1 FROM archived_data.oauth_session_store_archived a WHERE a.user_id = t.user_id
      );
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
      COMMIT;
    END IF;

    -- note_comments
    SET v_tab_name = 'note_comments';
    IF EXISTS (SELECT 1 FROM note_comments t JOIN temp_user_ids b ON b.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.note_comments_archived
      SELECT t.* FROM note_comments t JOIN temp_user_ids b ON b.user_id = t.user_id;
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
      DELETE t FROM note_comments t
      JOIN temp_user_ids b ON b.user_id = t.user_id
      WHERE EXISTS (
        SELECT 1 FROM archived_data.note_comments_archived a WHERE a.user_id = t.user_id
      );
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
      COMMIT;
    END IF;

    -- invalid_entry_attempts
    SET v_tab_name = 'invalid_entry_attempts';
    IF EXISTS (SELECT 1 FROM invalid_entry_attempts t JOIN temp_user_ids b ON b.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.invalid_entry_attempts_archived
      SELECT t.* FROM invalid_entry_attempts t JOIN temp_user_ids b ON b.user_id = t.user_id;
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
      DELETE t FROM invalid_entry_attempts t
      JOIN temp_user_ids b ON b.user_id = t.user_id
      WHERE EXISTS (
        SELECT 1 FROM archived_data.invalid_entry_attempts_archived a WHERE a.user_id = t.user_id
      );
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
      COMMIT;
    END IF;

    -- users_roles
    SET v_tab_name = 'users_roles';
    IF EXISTS (SELECT 1 FROM users_roles t JOIN temp_user_ids b ON b.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.users_roles_archived
      SELECT t.* FROM users_roles t JOIN temp_user_ids b ON b.user_id = t.user_id;
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
      DELETE t FROM users_roles t
      JOIN temp_user_ids b ON b.user_id = t.user_id
      WHERE EXISTS (
        SELECT 1 FROM archived_data.users_roles_archived a WHERE a.user_id = t.user_id
      );
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
      COMMIT;
    END IF;

    -- user-specific tables not yet covered
    -- users_notifications_settings etc already handled; continue with remaining list

    -- users_site tables done earlier; continue

    -- users table (main) - process last
    SET v_tab_name = 'users';
    IF EXISTS (SELECT 1 FROM users u JOIN temp_user_ids b ON b.user_id = u.id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.users_archived
      SELECT u.* FROM users u JOIN temp_user_ids b ON b.user_id = u.id;
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records inserted : ', v_rec_cnt), NULL, v_user);
      DELETE u FROM users u
      JOIN temp_user_ids b ON b.user_id = u.id
      WHERE EXISTS (
        SELECT 1 FROM archived_data.users_archived a WHERE a.id = u.id
      );
      SET v_rec_cnt = ROW_COUNT();
      INSERT INTO user_deletion_detail (batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, CONCAT('Number of Records deleted : ', v_rec_cnt), NULL, v_user);
      COMMIT;
    END IF;

    -- 4) Mark completion and clean up
    START TRANSACTION;
    UPDATE deleted_users_info di
    JOIN temp_user_ids b ON b.id = di.id
    SET di.deletion_status = 'completed',
        di.note = 'the user_id is deleted';
    COMMIT;

    DROP TEMPORARY TABLE IF EXISTS temp_user_ids;

    -- Finished one batch, loop will continue to look for next batch
  END WHILE Outer_check;

  -- Completion message
  SELECT 'User deletion process completed.' AS message;
END $$

DELIMITER ;


