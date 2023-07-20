
-- +migrate Up

CREATE TABLE IF NOT EXISTS `group_profile` (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    profile_id INTEGER,
    group_id INTEGER,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,

    
    CONSTRAINT fk_group_profile_profile FOREIGN KEY (`profile_id`) REFERENCES `profile`(`id`),
    CONSTRAINT fkgroup_profile_group FOREIGN KEY (`group_id`) REFERENCES `group`(`id`)
     
);

-- +migrate Down
DROP TABLE `group_profile`;
