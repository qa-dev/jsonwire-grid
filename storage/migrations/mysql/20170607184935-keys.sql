
-- +migrate Up
CREATE INDEX `part_of_name` ON `node` (`status`,`updated`);
CREATE INDEX `part_of_name` ON `node` (`sessionId`);
DROP INDEX `address` ON `node`;
ALTER TABLE `node` ADD PRIMARY KEY (`address`);
ALTER TABLE  `capabilities` ADD  CONSTRAINT `address` FOREIGN KEY (`nodeAddress`) REFERENCES `node` (`address`) ON DELETE CASCADE ON UPDATE NO ACTION

-- +migrate Down
