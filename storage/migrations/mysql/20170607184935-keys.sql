
-- +migrate Up

CREATE INDEX `status_of_updated` ON `node` (`status`,`updated`);
CREATE INDEX `sessionId` ON `node` (`sessionId`);
ALTER TABLE `node` ADD PRIMARY KEY (`address`);
DROP INDEX `address` ON `node`;
CREATE INDEX `address` ON `capabilities` (`nodeAddress`);

-- +migrate Down
SIGNAL SQLSTATE '45000' SET message_text = 'Impossible down this migration';