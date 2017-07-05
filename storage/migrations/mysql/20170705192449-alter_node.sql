
-- +migrate Up
ALTER TABLE `node` MODIFY COLUMN `address` varchar(255) NOT NULL DEFAULT '';

-- +migrate Down
ALTER TABLE `node` MODIFY COLUMN `address` varchar(20) NOT NULL DEFAULT '';