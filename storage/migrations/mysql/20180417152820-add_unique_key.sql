
-- +migrate Up
ALTER TABLE `node` ADD COLUMN `key` varchar(255) NOT NULL AFTER `id`;
ALTER TABLE `capabilities` CHANGE COLUMN `nodeAddress` `nodeKey` VARCHAR(255) NOT NULL;
CREATE UNIQUE INDEX `key` ON `node` (`key`);

-- +migrate Down
SIGNAL SQLSTATE '45000' SET message_text = 'Impossible down this migration';