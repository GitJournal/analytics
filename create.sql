/* CREATE DATABASE IF NOT EXISTS gitjournal */
CREATE TABLE events
(
  `timestamp` DateTime('UTC'),
  `name` String,
  `params` Array(Tuple(String, String)),
  `previous_timestamp` DateTime('UTC'),
  `bundle_sequence_id` Int32,

  `user_pseudo_id` UUID,
  `user_properties` Array(Tuple(String, String)),

  `device` Nested(
    `category` LowCardinality(String),
    `mobile_brand_name` LowCardinality(String),
    `mobile_model_name` LowCardinality(String),
    `mobile_marketing_name` LowCardinality(String),
    `mobile_os_hardware_model` LowCardinality(String),
    `operating_system` LowCardinality(String),
    `operating_system_version` LowCardinality(String),
    `language` LowCardinality(String),
    `time_zone_offset_seconds` Int16
  ),

  `geo` Nested(
    `continent` LowCardinality(String),
    `country` LowCardinality(String),
    `region` LowCardinality(String),
    `city` LowCardinality(String),
    `sub_continent` LowCardinality(String),
    `metro` LowCardinality(String)
  ),

  `app_info` Nested(
    `id` LowCardinality(String),
    `version` LowCardinality(String)
  ),

  `stream_id` Int64,
  `platform` LowCardinality(String)
)

ENGINE = MergeTree()
PARTITION BY toYYYYMM(`timestamp`)
ORDER BY (`timestamp`)
