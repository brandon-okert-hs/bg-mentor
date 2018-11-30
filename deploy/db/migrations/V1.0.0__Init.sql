CREATE TABLE `members` (
    `id` int NOT NULL AUTO_INCREMENT,
    `name` varchar(100) NOT NULL,
    `email` varchar(200) NOT NULL,
    `avatarUrl` varchar(2048) DEFAULT NULL,
    `idtoken` varchar(2048) NOT NULL,
    `accessToken` varchar(2048) NOT NULL,
    `sub` varchar(2048) NOT NULL,
    PRIMARY KEY (`id`)
);

CREATE TABLE `entitlements` (
    `id` int NOT NULL AUTO_INCREMENT,
    `member_id` int NOT NULL,
    `create_tournaments` boolean NOT NULL DEFAULT 0,
    `manage_tournaments` boolean NOT NULL DEFAULT 0,
    `manage_members` boolean NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`)
);
