
-- +migrate Up
ALTER TABLE `node` ADD COLUMN `key` varchar(255) NOT NULL AFTER `id`;

-- +migrate Down
SIGNAL SQLSTATE '45000' SET message_text = 'Impossible down this migration';