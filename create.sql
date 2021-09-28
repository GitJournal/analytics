create table analytics_events (
  ts timestamp not null,
  event_name text not null,
  props jsonb not null,

  pseudoId text,
  userId text,
  user_props jsonb not null,

  session_id integer not null,
  location_id integer not null,
  device_id bigint not null,
  package_id bigint not null
);

CREATE TYPE valid_platform AS ENUM ('android', 'ios', 'linux', 'macos', 'windows', 'web');

create table analytics_device_info (
  id bigint primary key,

  platform valid_platform,

  android_info jsonb,
  ios_info jsonb,
  linux_info jsonb,
  macos_info jsonb,
  windows_info jsonb,
  web_info jsonb
);

create table analytics_package_info (
  id bigint primary key,

  appName text not null,
  packageName text not null,
  version text not null,
  buildNumber text not null,
  buildSignature text not null
);

create table analytics_location (
  city_geoname_id integer primary key,
  city_name_en text not null,
  country_code text not null
);

alter table analytics_package_info add installSource text not null;

