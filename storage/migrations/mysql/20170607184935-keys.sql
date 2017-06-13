
-- +migrate Up

ALTER TABLE `node` ADD COLUMN `id` BIGINT(20) NOT NULL AUTO_INCREMENT PRIMARY KEY FIRST;
ALTER TABLE `capabilities` ADD COLUMN `id` BIGINT(20) NOT NULL AUTO_INCREMENT PRIMARY KEY FIRST;
CREATE INDEX `status_of_updated` ON `node` (`status`,`updated`);
CREATE INDEX `address` ON `capabilities` (`nodeAddress`);
CREATE INDEX `sessionId` ON `node` (`sessionId`);
DROP INDEX `addressSetName` ON `capabilities`;
CREATE UNIQUE INDEX `addressSetName` ON `capabilities` (`name`, `nodeAddress`, `setId`);
# SLEEP 1;

-- +migrate Down
SIGNAL SQLSTATE '45000' SET message_text = 'Impossible down this migration';