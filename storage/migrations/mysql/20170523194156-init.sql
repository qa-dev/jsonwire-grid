
-- +migrate Up
CREATE TABLE `node` (
  `type` varchar(20) NOT NULL DEFAULT '',
  `status` varchar(20) NOT NULL DEFAULT '',
  `address` varchar(20) NOT NULL DEFAULT '',
  `sessionId` varchar(255) NOT NULL DEFAULT '',
  `updated` bigint(20) NOT NULL,
  `registred` bigint(20) NOT NULL,
  UNIQUE KEY `address` (`address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `capabilities` (
  `nodeAddress` varchar(20) NOT NULL DEFAULT '',
  `browserName` varchar(50) NOT NULL DEFAULT '',
  `browserVersion` varchar(50) NOT NULL DEFAULT '',
  `platformName` varchar(50) NOT NULL DEFAULT '',
  `pageLoadStrategy` varchar(50) NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- +migrate Down
DROP TABLE `capabilities`;
DROP TABLE `node`;
