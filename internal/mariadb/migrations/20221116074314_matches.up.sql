CREATE TABLE swipes (
	user_id 		INT4 UNSIGNED NOT NULL REFERENCES users(id),
	profile_id 	INT4 UNSIGNED NOT NULL REFERENCES users(id),
	preference 	BOOLEAN NOT NULL,
	created_at	TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

	PRIMARY KEY(user_id, profile_id),
	CHECK(user_id != profile_id)
);

CREATE TABLE matches (
	id 						INT8 UNSIGNED AUTO_INCREMENT NOT NULL,
	user_id_low 	INT4 UNSIGNED NOT NULL REFERENCES users(id),
	user_id_high 	INT4 UNSIGNED NOT NULL REFERENCES users(id),
	created_at		TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

	PRIMARY KEY(id),
	CHECK(user_id_low < user_id_high),
	UNIQUE(user_id_low, user_id_high)
);
