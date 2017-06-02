
-- +migrate Up

CREATE TABLE `capabilities_ref` (
  `nodeAddress` VARCHAR(20) NOT NULL DEFAULT '',
  `setId`       VARCHAR(50) NOT NULL DEFAULT '',
  `name`        VARCHAR(50) NOT NULL DEFAULT '',
  `value`       VARCHAR(50) NOT NULL DEFAULT '',
  UNIQUE KEY `addressSetName` (`nodeAddress`, `setId`, `name`)
) ENGINE=InnoDB CHARSET=utf8;

INSERT INTO `capabilities_ref` (`nodeAddress`, `setId`, `name`, `value`)
  (SELECT `nodeAddress`, CONCAT(`nodeAddress`, `browserName`, `browserVersion`, `platformName`, `pageLoadStrategy`), "browserName", `browserName` as `value`  FROM capabilities HAVING `value` <> '')
  UNION
  (SELECT `nodeAddress`, CONCAT(`nodeAddress`, `browserName`, `browserVersion`, `platformName`, `pageLoadStrategy`), "browserVersion", `browserVersion` as `value`  FROM capabilities HAVING `value` <> '')
  UNION
  (SELECT `nodeAddress`, CONCAT(`nodeAddress`, `browserName`, `browserVersion`, `platformName`, `pageLoadStrategy`), "platformName", `platformName` as `value`  FROM capabilities HAVING `value` <> '')
  UNION
  (SELECT `nodeAddress`, CONCAT(`nodeAddress`, `browserName`, `browserVersion`, `platformName`, `pageLoadStrategy`), "pageLoadStrategy", `pageLoadStrategy` as `value`  FROM capabilities HAVING `value` <> '');

DROP TABLE `capabilities`;
RENAME TABLE `capabilities_ref` TO `capabilities`;

UNLOCK TABLES;

-- +migrate Down
SIGNAL SQLSTATE '45000' SET message_text = 'Impossible down this migration';