CREATE TABLE IF NOT EXISTS users (
	id         BINARY(16) NOT NULL,
	nickname   TEXT       NOT NULL,
	password   TEXT       NOT NULL,
	created_at DATETIME   NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;

CREATE TABLE IF NOT EXISTS `profiles` (
	id         BINARY(16) NOT NULL,
	first_name TEXT       NOT NULL,
	last_name  TEXT       NOT NULL,
	age        INT        NOT NULL,
	gender     TEXT       NOT NULL,
	interests  TEXT       NOT NULL, -- store as text, because MySQL don't support array.
	city       TEXT       NOT NULL,
	created_at DATETIME   NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;
