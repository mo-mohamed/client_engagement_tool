
-- +migrate Up

CREATE TABLE IF NOT EXISTS `mobile_terminated` (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    profile_id INTEGER,
    group_id INTEGER,
    group_message_uuid VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,

    CONSTRAINT fk_mt_profile FOREIGN KEY (`profile_id`) REFERENCES `profile`(`id`),
    CONSTRAINT fk_mt_group FOREIGN KEY (`group_id`) REFERENCES `group`(`id`)
   
);
-- +migrate Down

DROP TABLE `mobile_terminated`;
