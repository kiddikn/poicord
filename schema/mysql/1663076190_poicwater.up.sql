CREATE TABLE poic_water (
	`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	`user_id` varchar(255) not null,
	`started_at` datetime not null,
    `finished_at` datetime,
    `revoked_at` datetime,
    `created_at` datetime not null default CURRENT_TIMESTAMP,
	primary key (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE=utf8mb4_general_ci;
