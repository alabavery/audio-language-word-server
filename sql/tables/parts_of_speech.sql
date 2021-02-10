CREATE TABLE parts_of_speech (
    word            bigint,
    part_of_speech  varchar(30) NOT NULL,
    CONSTRAINT      fk_word
        FOREIGN KEY (word)
            REFERENCES word(id)
);
