CREATE TABLE `migration_meta` (
    `id` integer PRIMARY KEY AUTOINCREMENT NOT NULL,
    `version` integer NOT NULL,
    `applied_at` datetime NOT NULL
);

CREATE TABLE `pdf_summary` (
	`id` integer PRIMARY KEY AUTOINCREMENT NOT NULL,
	`file` text,
	`summary` text,
    `title` text,
	`intermediate_summary` text
);