DELIMITER $$

CREATE PROCEDURE update_column_batch(
    IN tbl_name VARCHAR(255),
    IN col_name VARCHAR(255),
    IN old_value VARCHAR(255),
    IN new_value VARCHAR(255),
    IN batch_size INT
)
BEGIN
    DECLARE rows_affected INT DEFAULT 1;

    SET @sql_update = CONCAT(
            'UPDATE LOW_PRIORITY ', tbl_name,
            ' SET ', col_name, ' = ?',
            ' WHERE ', col_name, ' = ?',
            ' LIMIT ', batch_size
                      );
    PREPARE stmt FROM @sql_update;
    SET @old_val = old_value;
    SET @new_val = new_value;

    WHILE rows_affected > 0
        DO
            EXECUTE stmt USING @new_val, @old_val;
            SET rows_affected = ROW_COUNT();
        END WHILE;

    DEALLOCATE PREPARE stmt;
END$$

DELIMITER ;

CREATE INDEX idx_account_schedule_owner_type
    ON account_schedule (owner_type);

CREATE INDEX idx_operator_group_limits_owner_type
    ON operator_group_limits (owner_type);

START TRANSACTION;
CALL update_column_batch(
        'account_schedule',
        'owner_type',
        'account_groups',
        'operator_groups',
        100000
     );
COMMIT;

START TRANSACTION;
CALL update_column_batch(
        'operator_group_limits',
        'owner_type',
        'account_groups',
        'operator_groups',
        100000
     );
COMMIT;

DROP INDEX idx_account_schedule_owner_type
    ON account_schedule;

DROP INDEX idx_operator_group_limits_owner_type
    ON operator_group_limits;

DROP PROCEDURE update_column_batch;