CREATE TABLE user_known_parts_of_speech (
    id              bigserial PRIMARY KEY,
    "user_id"       bigint,
    part_of_speech  bigint,
    CONSTRAINT      fk_parts_of_speech
        FOREIGN KEY(part_of_speech)
                REFERENCES part_of_speech(id)
);