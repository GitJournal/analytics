/* CREATE DATABASE IF NOT EXISTS gitjournal */
CREATE TABLE events
(
    `timestamp` DateTime('UTC'),
    `name` STRING,
    `params` STRING, /* how to store this? */
    `previous_timestamp`: DateTime('UTC'),
    `bundle_sequence_id`: Int32,

    `user_pseudo_id` UUID,
    `user_properties` STRING, /* how to store this? */

    `device.category` LowCardinality(STRING),
    `device.mobile_brand_name` LowCardinality(STRING),
    `device.mobile_model_name` LowCardinality(STRING),
    `device.mobile_marketing_name` LowCardinality(STRING),
    `device.mobile_os_hardware_model` LowCardinality(STRING),
    `device.operating_system` LowCardinality(STRING),
    `device.operating_system_version` LowCardinality(STRING),
    `device.language` LowCardinality(STRING),
    `device.time_zone_offset_seconds` Int16,

    `geo.continent` LowCardinality(STRING),
    `geo.country` LowCardinality(STRING),
    `geo.region` LowCardinality(STRING),
    `geo.city` LowCardinality(STRING),
    `geo.sub_continent` LowCardinality(STRING),
    `geo.metro` LowCardinality(STRING),

    `app_info.id` LowCardinality(STRING),
    `app_info.version` LowCardinality(STRING),

    `stream_id` Int64,
    `platform` LowCardinality(STRING)
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(`timestamp`)
