DELIMITER $$

DROP PROCEDURE IF EXISTS delete_users_in_batches$$
CREATE PROCEDURE delete_users_in_batches()
BEGIN
  DECLARE done INT DEFAULT FALSE;
  DECLARE num_deleted INT DEFAULT 0;
  DECLARE v_batch_id INT;
  DECLARE v_note VARCHAR(255);
  DECLARE v_error_info VARCHAR(300);
  DECLARE v_rec_cnt INT DEFAULT 0;
  DECLARE v_tab_name VARCHAR(50);
  DECLARE v_user INT DEFAULT 2920483;

  DECLARE v_sqlstate CHAR(5) DEFAULT '00000';
  DECLARE v_errno INT DEFAULT 0;
  DECLARE v_message TEXT;

  DECLARE CONTINUE HANDLER FOR SQLEXCEPTION
  BEGIN
    GET DIAGNOSTICS CONDITION 1 v_sqlstate = RETURNED_SQLSTATE, v_errno = MYSQL_ERRNO, v_message = MESSAGE_TEXT;
    SET v_error_info = CONCAT('ERROR ', v_errno, ' (', v_sqlstate, '): ', v_message);
    ROLLBACK;
    INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
    VALUES (v_batch_id, v_tab_name, 'failed', v_error_info, v_user);
    COMMIT;
    SET done = TRUE;
  END;

  Outer_check: WHILE NOT done DO
    SELECT COALESCE(MAX(batch_id), 0) + 1 INTO v_batch_id FROM deleted_users_info;

    START TRANSACTION;
    UPDATE deleted_users_info
    SET batch_id = v_batch_id
    WHERE batch_id IS NULL;
    COMMIT;

    START TRANSACTION;
    UPDATE deleted_users_info d
      JOIN users u ON u.id = d.user_id
      JOIN users_roles ur ON ur.user_id = u.id
      JOIN roles r ON r.id = ur.role_id AND r.name IN ('Tenant','Tenant 24/7')
    SET d.deletion_status = 'in progress'
    WHERE d.batch_id = v_batch_id
      AND d.deletion_status = 'pending'
      AND u.last_login_date < DATE_SUB(CURRENT_TIMESTAMP(), INTERVAL 1 YEAR)
      AND NOT EXISTS (SELECT 1 FROM v2_shares s WHERE s.owner_id = d.user_id OR s.recipient_id = d.user_id)
      AND NOT EXISTS (SELECT 1 FROM v2_units vu WHERE vu.user_id = d.user_id)
      AND NOT EXISTS (SELECT 1 FROM auction a WHERE a.user_id = d.user_id)
      AND NOT EXISTS (SELECT 1 FROM prelets p WHERE p.user_id = d.user_id);

    UPDATE deleted_users_info d
    SET d.deletion_status = 'canceled', d.note = 'The user_id is either an Owner or Shared user, so cannot delete'
    WHERE d.batch_id = v_batch_id
      AND d.deletion_status = 'pending'
      AND EXISTS (SELECT 1 FROM v2_shares s WHERE s.owner_id = d.user_id OR s.recipient_id = d.user_id);

    UPDATE deleted_users_info d
    SET d.deletion_status = 'canceled', d.note = 'The user_id is having a unit, so cannot delete'
    WHERE d.batch_id = v_batch_id
      AND d.deletion_status = 'pending'
      AND EXISTS (SELECT 1 FROM v2_units vu WHERE vu.user_id = d.user_id);

    UPDATE deleted_users_info d
    SET d.deletion_status = 'canceled', d.note = 'The user_id is having an auction, so cannot delete'
    WHERE d.batch_id = v_batch_id
      AND d.deletion_status = 'pending'
      AND EXISTS (SELECT 1 FROM auction a WHERE a.user_id = d.user_id);

    UPDATE deleted_users_info d
    SET d.deletion_status = 'canceled', d.note = 'The user_id is having a prelet, so cannot delete'
    WHERE d.batch_id = v_batch_id
      AND d.deletion_status = 'pending'
      AND EXISTS (SELECT 1 FROM prelets p WHERE p.user_id = d.user_id);
    COMMIT;

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
      SET done = TRUE;
      LEAVE Outer_check;
    END IF;

    SET v_tab_name = 'unit_status_updates';
    IF EXISTS (SELECT 1 FROM unit_status_updates s JOIN temp_user_ids t ON s.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.unit_status_updates_archived
      SELECT s.* FROM unit_status_updates s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      DELETE s FROM unit_status_updates s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids)
        AND EXISTS (SELECT 1 FROM archived_data.unit_status_updates_archived a WHERE a.user_id = s.user_id);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      COMMIT;
    END IF;

    SET v_tab_name = 'unit_overrides';
    IF EXISTS (SELECT 1 FROM unit_overrides s JOIN temp_user_ids t ON s.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.unit_overrides_archived
      SELECT s.* FROM unit_overrides s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      DELETE s FROM unit_overrides s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids)
        AND EXISTS (SELECT 1 FROM archived_data.unit_overrides_archived a WHERE a.user_id = s.user_id);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      COMMIT;
    END IF;

    SET v_tab_name = 'transfer';
    IF EXISTS (SELECT 1 FROM transfer s JOIN temp_user_ids t ON s.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.transfer_archived
      SELECT s.* FROM transfer s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      DELETE s FROM transfer s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids)
        AND EXISTS (SELECT 1 FROM archived_data.transfer_archived a WHERE a.user_id = s.user_id);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      COMMIT;
    END IF;

    SET v_tab_name = 'support_activity';
    IF EXISTS (SELECT 1 FROM support_activity s JOIN temp_user_ids t ON s.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.support_activity_archived
      SELECT s.* FROM support_activity s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      DELETE s FROM support_activity s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids)
        AND EXISTS (SELECT 1 FROM archived_data.support_activity_archived a WHERE a.user_id = s.user_id);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      COMMIT;
    END IF;

    SET v_tab_name = 'unlock_override_pins';
    IF EXISTS (SELECT 1 FROM unlock_override_pins s JOIN temp_user_ids t ON s.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.unlock_override_pins_archived
      SELECT s.* FROM unlock_override_pins s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      DELETE s FROM unlock_override_pins s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids)
        AND EXISTS (SELECT 1 FROM archived_data.unlock_override_pins_archived a WHERE a.user_id = s.user_id);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      COMMIT;
    END IF;

    SET v_tab_name = 'user_acknowledgement';
    IF EXISTS (SELECT 1 FROM user_acknowledgement s JOIN temp_user_ids t ON s.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.user_acknowledgement_archived
      SELECT s.* FROM user_acknowledgement s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      DELETE s FROM user_acknowledgement s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids)
        AND EXISTS (SELECT 1 FROM archived_data.user_acknowledgement_archived a WHERE a.user_id = s.user_id);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      COMMIT;
    END IF;

    SET v_tab_name = 'user_login_history';
    IF EXISTS (SELECT 1 FROM user_login_history s JOIN temp_user_ids t ON s.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.user_login_history_archived
      SELECT s.* FROM user_login_history s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      DELETE s FROM user_login_history s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids)
        AND EXISTS (SELECT 1 FROM archived_data.user_login_history_archived a WHERE a.user_id = s.user_id);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      COMMIT;
    END IF;

    SET v_tab_name = 'user_site_bookmarks';
    IF EXISTS (SELECT 1 FROM user_site_bookmarks s JOIN temp_user_ids t ON s.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.user_site_bookmarks_archived
      SELECT s.* FROM user_site_bookmarks s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      DELETE s FROM user_site_bookmarks s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids)
        AND EXISTS (SELECT 1 FROM archived_data.user_site_bookmarks_archived a WHERE a.user_id = s.user_id);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      COMMIT;
    END IF;

    SET v_tab_name = 'users_dashboard';
    IF EXISTS (SELECT 1 FROM users_dashboard s JOIN temp_user_ids t ON s.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.users_dashboard_archived
      SELECT s.* FROM users_dashboard s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      DELETE s FROM users_dashboard s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids)
        AND EXISTS (SELECT 1 FROM archived_data.users_dashboard_archived a WHERE a.user_id = s.user_id);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      COMMIT;
    END IF;

    SET v_tab_name = 'users_notifications_settings';
    IF EXISTS (SELECT 1 FROM users_notifications_settings s JOIN temp_user_ids t ON s.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.users_notifications_settings_archived
      SELECT s.* FROM users_notifications_settings s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      DELETE s FROM users_notifications_settings s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids)
        AND EXISTS (SELECT 1 FROM archived_data.users_notifications_settings_archived a WHERE a.user_id = s.user_id);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      COMMIT;
    END IF;

    SET v_tab_name = 'users_zone_triggers';
    IF EXISTS (SELECT 1 FROM users_zone_triggers s JOIN temp_user_ids t ON s.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.users_zone_triggers_archived
      SELECT s.* FROM users_zone_triggers s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      DELETE s FROM users_zone_triggers s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids)
        AND EXISTS (SELECT 1 FROM archived_data.users_zone_triggers_archived a WHERE a.user_id = s.user_id);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      COMMIT;
    END IF;

    SET v_tab_name = 'v2_app_details';
    IF EXISTS (SELECT 1 FROM v2_app_details s JOIN temp_user_ids t ON s.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.v2_app_details_archived
      SELECT s.* FROM v2_app_details s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      DELETE s FROM v2_app_details s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids)
        AND EXISTS (SELECT 1 FROM archived_data.v2_app_details_archived a WHERE a.user_id = s.user_id);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      COMMIT;
    END IF;

    SET v_tab_name = 'v2_tracking_ids';
    IF EXISTS (SELECT 1 FROM v2_tracking_ids s JOIN temp_user_ids t ON s.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.v2_tracking_ids_archived
      SELECT s.* FROM v2_tracking_ids s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      DELETE s FROM v2_tracking_ids s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids)
        AND EXISTS (SELECT 1 FROM archived_data.v2_tracking_ids_archived a WHERE a.user_id = s.user_id);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      COMMIT;
    END IF;

    SET v_tab_name = 'watch_users';
    IF EXISTS (SELECT 1 FROM watch_users s JOIN temp_user_ids t ON s.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.watch_users_archived
      SELECT s.* FROM watch_users s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      DELETE s FROM watch_users s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids)
        AND EXISTS (SELECT 1 FROM archived_data.watch_users_archived a WHERE a.user_id = s.user_id);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      COMMIT;
    END IF;

    SET v_tab_name = '2_factor_auth_pin';
    IF EXISTS (SELECT 1 FROM `2_factor_auth_pin` s JOIN temp_user_ids t ON s.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.`2_factor_auth_pin_archived`
      SELECT s.* FROM `2_factor_auth_pin` s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      DELETE s FROM `2_factor_auth_pin` s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids)
        AND EXISTS (SELECT 1 FROM archived_data.`2_factor_auth_pin_archived` a WHERE a.user_id = s.user_id);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      COMMIT;
    END IF;

    SET v_tab_name = 'access_codes';
    IF EXISTS (SELECT 1 FROM access_codes s JOIN temp_user_ids t ON s.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.access_codes_archived
      SELECT s.* FROM access_codes s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      DELETE s FROM access_codes s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids)
        AND EXISTS (SELECT 1 FROM archived_data.access_codes_archived a WHERE a.user_id = s.user_id);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      COMMIT;
    END IF;

    SET v_tab_name = 'devices';
    IF EXISTS (SELECT 1 FROM devices s JOIN temp_user_ids t ON s.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.devices_archived
      SELECT s.* FROM devices s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      DELETE s FROM devices s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids)
        AND EXISTS (SELECT 1 FROM archived_data.devices_archived a WHERE a.user_id = s.user_id);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      COMMIT;
    END IF;

    SET v_tab_name = 'digital_audits';
    IF EXISTS (SELECT 1 FROM digital_audits s JOIN temp_user_ids t ON s.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.digital_audits_archived
      SELECT s.* FROM digital_audits s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      DELETE s FROM digital_audits s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids)
        AND EXISTS (SELECT 1 FROM archived_data.digital_audits_archived a WHERE a.user_id = s.user_id);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      COMMIT;
    END IF;

    SET v_tab_name = 'entry_activity';
    IF EXISTS (SELECT 1 FROM entry_activity s JOIN temp_user_ids t ON s.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.entry_activity_archived
      SELECT s.* FROM entry_activity s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      DELETE s FROM entry_activity s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids)
        AND EXISTS (SELECT 1 FROM archived_data.entry_activity_archived a WHERE a.user_id = s.user_id);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      COMMIT;
    END IF;

    SET v_tab_name = 'permissions';
    IF EXISTS (SELECT 1 FROM permissions s JOIN temp_user_ids t ON s.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.permissions_archived
      SELECT s.* FROM permissions s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      DELETE s FROM permissions s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids)
        AND EXISTS (SELECT 1 FROM archived_data.permissions_archived a WHERE a.user_id = s.user_id);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      COMMIT;
    END IF;

    SET v_tab_name = 'pending_notifications';
    IF EXISTS (SELECT 1 FROM pending_notifications s JOIN temp_user_ids t ON s.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.pending_notifications_archived
      SELECT s.* FROM pending_notifications s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      DELETE s FROM pending_notifications s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids)
        AND EXISTS (SELECT 1 FROM archived_data.pending_notifications_archived a WHERE a.user_id = s.user_id);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      COMMIT;
    END IF;

    SET v_tab_name = 'oauth_session_store';
    IF EXISTS (SELECT 1 FROM oauth_session_store s JOIN temp_user_ids t ON s.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.oauth_session_store_archived
      SELECT s.* FROM oauth_session_store s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      DELETE s FROM oauth_session_store s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids)
        AND EXISTS (SELECT 1 FROM archived_data.oauth_session_store_archived a WHERE a.user_id = s.user_id);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      COMMIT;
    END IF;

    SET v_tab_name = 'note_comments';
    IF EXISTS (SELECT 1 FROM note_comments s JOIN temp_user_ids t ON s.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.note_comments_archived
      SELECT s.* FROM note_comments s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      DELETE s FROM note_comments s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids)
        AND EXISTS (SELECT 1 FROM archived_data.note_comments_archived a WHERE a.user_id = s.user_id);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      COMMIT;
    END IF;

    SET v_tab_name = 'invalid_entry_attempts';
    IF EXISTS (SELECT 1 FROM invalid_entry_attempts s JOIN temp_user_ids t ON s.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.invalid_entry_attempts_archived
      SELECT s.* FROM invalid_entry_attempts s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      DELETE s FROM invalid_entry_attempts s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids)
        AND EXISTS (SELECT 1 FROM archived_data.invalid_entry_attempts_archived a WHERE a.user_id = s.user_id);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      COMMIT;
    END IF;

    SET v_tab_name = 'users_roles';
    IF EXISTS (SELECT 1 FROM users_roles s JOIN temp_user_ids t ON s.user_id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.users_roles_archived
      SELECT s.* FROM users_roles s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      DELETE s FROM users_roles s WHERE s.user_id IN (SELECT user_id FROM temp_user_ids)
        AND EXISTS (SELECT 1 FROM archived_data.users_roles_archived a WHERE a.user_id = s.user_id);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      COMMIT;
    END IF;

    SET v_tab_name = 'users';
    IF EXISTS (SELECT 1 FROM users s JOIN temp_user_ids t ON s.id = t.user_id LIMIT 1) THEN
      START TRANSACTION;
      INSERT INTO archived_data.users_archived
      SELECT s.* FROM users s WHERE s.id IN (SELECT user_id FROM temp_user_ids);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records inserted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      DELETE s FROM users s WHERE s.id IN (SELECT user_id FROM temp_user_ids)
        AND EXISTS (SELECT 1 FROM archived_data.users_archived a WHERE a.id = s.id);
      SET v_rec_cnt = ROW_COUNT();
      SET v_note = CONCAT('Number of Records deleted : ', v_rec_cnt);
      INSERT INTO user_deletion_detail(batch_id, table_name, note, error_info, last_updated_by)
      VALUES (v_batch_id, v_tab_name, v_note, NULL, v_user);
      COMMIT;
    END IF;

    START TRANSACTION;
    UPDATE deleted_users_info d
      JOIN temp_user_ids t ON t.id = d.id
    SET d.deletion_status = 'completed', d.note = 'the user_id is deleted'
    WHERE d.batch_id = v_batch_id AND d.deletion_status = 'in progress';
    COMMIT;

    DROP TEMPORARY TABLE IF EXISTS temp_user_ids;
  END WHILE;

  SELECT 'User deletion process completed.' AS message;
END$$

DELIMITER ;
