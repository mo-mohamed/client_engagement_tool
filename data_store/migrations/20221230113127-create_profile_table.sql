
-- +migrate Up

CREATE TABLE IF NOT EXISTS `profile` (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    mdn VARCHAR(100) NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,

    group_id INTEGER,
    CONSTRAINT fk_group_id FOREIGN KEY (`group_id`) REFERENCES `group`(`id`)

);

-- +migrate Down

DROP TABLE `profile`;
