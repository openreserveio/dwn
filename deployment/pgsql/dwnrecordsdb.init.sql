DROP TABLE IF EXISTS record;
DROP TABLE IF EXISTS message_entry;

CREATE TABLE record (
    id VARCHAR PRIMARY KEY NOT NULL,
    record_id VARCHAR NOT NULL,
    creator_did VARCHAR NOT NULL,
    owner_did VARCHAR NOT NULL,
    writer_dids VARCHAR[],
    reader_dids VARCHAR[],
    initial_entry_id VARCHAR NOT NULL,
    latest_entry_id VARCHAR NOT NULL,
    latest_checkpoint_entry_id VARCHAR NOT NULL,
    create_date TIMESTAMP NOT NULL
);
CREATE INDEX idx_record_record_id ON record(record_id);
CREATE INDEX idx_record_creator_did ON record(creator_did);
CREATE INDEX idx_record_owner_did ON record(owner_did);


CREATE TABLE message_entry (
    id VARCHAR PRIMARY KEY NOT NULL,
    message_entry_id VARCHAR NOT NULL,
    previous_message_entry_id VARCHAR,
    record_id VARCHAR NOT NULL,
    message BYTEA NOT NULL,
    create_date TIMESTAMP NOT NULL
);
CREATE INDEX idx_message_entry_message_entry_id ON message_entry(message_entry_id);
CREATE INDEX idx_message_entry_previous_message_entry_id ON message_entry(previous_message_entry_id);
CREATE INDEX idx_message_entry_record_id ON message_entry(record_id);
