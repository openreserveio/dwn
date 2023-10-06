
DROP TABLE IF EXISTS hook_record;
DROP TABLE IF EXISTS hook_configuration_entry;

CREATE TABLE hook_record (
    id VARCHAR PRIMARY KEY NOT NULL,
    hook_record_id VARCHAR NOT NULL,
    creator_did VARCHAR NOT NULL,
    owner_did VARCHAR NOT NULL,
    initial_hook_config_entry_id VARCHAR,
    latest_hook_config_entry_id VARCHAR,
    create_date TIMESTAMP NOT NULL,
    filter_data_record_id VARCHAR,
    filter_schema VARCHAR,
    filter_protocol VARCHAR,
    filter_protocol_version VARCHAR,
);
CREATE INDEX idx_hook_record_hook_record_id ON hook_record(hook_record_id);
CREATE INDEX idx_hook_record_creator_did ON hook_record(creator_did);
CREATE INDEX idx_hook_record_owner_did ON hook_record(owner_did);
CREATE INDEX idx_hook_record_filter_data_record_id ON hook_record(filter_data_record_id);
CREATE INDEX idx_hook_record_filter_schema ON hook_record(filter_schema);
CREATE INDEX idx_hook_record_filter_protocol ON hook_record(filter_protocol);


CREATE TABLE hook_configuration_entry (
    id VARCHAR PRIMARY KEY NOT NULL,
    configuration_entry_id VARCHAR NOT NULL,
    hook_record_id VARCHAR NOT NULL,
    message BYTEA NOT NULL,
    create_date TIMESTAMP NOT NULL
);
CREATE INDEX idx_hook_configuration_entry_configuration_entry_id ON hook_configuration_entry(configuration_entry_id);
CREATE INDEX idx_hook_configuration_entry_hook_record_id ON hook_configuration_entry(hook_record_id);

