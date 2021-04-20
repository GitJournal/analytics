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

  /* Should I use nested instead - Not supported by offical go driver */
  `device.category` LowCardinality(String),
  `device.mobile_brand_name` LowCardinality(String),
  `device.mobile_model_name` LowCardinality(String),
  `device.mobile_marketing_name` LowCardinality(String),
  `device.mobile_os_hardware_model` LowCardinality(String),
  `device.operating_system` LowCardinality(String),
  `device.operating_system_version` LowCardinality(String),
  `device.language` LowCardinality(String),
  `device.time_zone_offset_seconds` Int16,

  `geo.continent` LowCardinality(String),
  `geo.country` LowCardinality(String),
  `geo.region` LowCardinality(String),
  `geo.city` LowCardinality(String),
  `geo.sub_continent` LowCardinality(String),
  `geo.metro` LowCardinality(String),

  `app_info.id` LowCardinality(String),
  `app_info.version` LowCardinality(String),

  `stream_id` Int64,
  `platform` LowCardinality(String)
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(`timestamp`)
ORDER BY (`timestamp`)
