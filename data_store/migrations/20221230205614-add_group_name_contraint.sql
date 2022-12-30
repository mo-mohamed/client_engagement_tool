
-- +migrate Up

ALTER TABLE `group`
ADD CONSTRAINT `UC_Name` UNIQUE (`name`);
-- +migrate Down

ALTER TABLE `group`
DROP INDEX `UC_Name`;
