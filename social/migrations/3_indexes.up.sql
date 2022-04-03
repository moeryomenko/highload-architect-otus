CREATE INDEX profile_creation ON profiles(created_at) ALGORITHM=INPLACE LOCK=NONE;


CREATE INDEX profile_name ON profiles(first_name(16), last_name(16)) ALGORITHM=INPLACE LOCK=NONE;


CREATE UNIQUE INDEX user_nickname_idx USING HASH ON users(nickname(16)) ALGORITHM=INPLACE LOCK=NONE;
