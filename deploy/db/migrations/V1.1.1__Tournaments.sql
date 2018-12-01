ALTER TABLE `tournaments` DROP COLUMN `startDate`, DROP COLUMN `checkinDate`;

ALTER TABLE `tournaments` ADD COLUMN `startDate` timestamp NOT NULL;
ALTER TABLE `tournaments` ADD COLUMN `checkinDate` timestamp NOT NULL;
