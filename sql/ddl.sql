CREATE TABLE object_type (
  id SMALLINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(64) NOT NULL,
  description VARCHAR(512),
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  delete_time timestamp NULL DEFAULT NULL,
  PRIMARY KEY(id),
  UNIQUE KEY(name)
) ENGINE=InnoDB CHARSET utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE object_status (
  id SMALLINT UNSIGNED NOT NULL AUTO_INCREMENT,
  type_id SMALLINT UNSIGNED NOT NULL,
  name VARCHAR(64),
  description VARCHAR(512),
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  delete_time timestamp NULL DEFAULT NULL,
  PRIMARY KEY(id),
  UNIQUE KEY(type_id, name)
) ENGINE=InnoDB CHARSET utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE object_state (
  id SMALLINT UNSIGNED NOT NULL AUTO_INCREMENT,
  status_id SMALLINT UNSIGNED NOT NULL,
  name VARCHAR(64),
  description VARCHAR(512),
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  delete_time timestamp NULL DEFAULT NULL,
  PRIMARY KEY(id),
  UNIQUE KEY(status_id, name)
) ENGINE=InnoDB CHARSET utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE object (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  type_id SMALLINT UNSIGNED NOT NULL,
  name VARCHAR(256),
  version BIGINT NOT NULL DEFAULT 0,
  relation_version BIGINT NOT NULL DEFAULT 0,
  description VARCHAR(1024),
  status_id SMALLINT UNSIGNED NOT NULL,
  state_id SMALLINT UNSIGNED NOT NULL,
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp NULL DEFAULT NULL,
  delete_time timestamp NULL DEFAULT NULL,
  PRIMARY KEY(id),
  UNIQUE KEY(type_id, name),
  INDEX(delete_time, id)
) ENGINE=InnoDB CHARSET utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE deleted_object (
  id BIGINT UNSIGNED NOT NULL,
  type_id SMALLINT UNSIGNED NOT NULL,
  name VARCHAR(256),
  version BIGINT NOT NULL,
  relation_version BIGINT NOT NULL DEFAULT 0,
  description VARCHAR(1024),
  status_id SMALLINT UNSIGNED NOT NULL,
  state_id SMALLINT UNSIGNED NOT NULL,
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp NULL DEFAULT NULL,
  delete_time timestamp NULL DEFAULT NULL,
  PRIMARY KEY(id),
  INDEX(type_id, delete_time)
) ENGINE=InnoDB CHARSET utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE object_meta (
  id MEDIUMINT UNSIGNED NOT NULL,
  type_id SMALLINT UNSIGNED NOT NULL,
  name VARCHAR(256),
  value_type TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '1: STRING 2: INTEGER, 3: DOUBLE, 4: BOOLEAN',
  description VARCHAR(1024),
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  delete_time timestamp NULL DEFAULT NULL,
  PRIMARY KEY(id),
  UNIQUE KEY(type_id, name)
) ENGINE=InnoDB CHARSET utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE object_meta_value (
  object_id BIGINT UNSIGNED NOT NULL,
  meta_id MEDIUMINT UNSIGNED NOT NULL,
  value TEXT,
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp NULL DEFAULT NULL,
  delete_time timestamp NULL DEFAULT NULL,
  PRIMARY KEY(object_id, meta_id),
  INDEX(delete_time, object_id, meta_id)
) ENGINE=InnoDB CHARSET utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE deleted_object_meta_value (
  object_id BIGINT UNSIGNED NOT NULL,
  meta_id MEDIUMINT UNSIGNED NOT NULL,
  value TEXT,
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp NULL DEFAULT NULL,
  delete_time timestamp NULL DEFAULT NULL,
  PRIMARY KEY(object_id, meta_id),
  INDEX(delete_time, object_id)
) ENGINE=InnoDB CHARSET utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE object_log (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  object_id BIGINT UNSIGNED NOT NULL,
  level TINYINT NOT NULL DEFAULT 0 COMMENT '0: EMERGENCY 1: ALERT 2: CRITICAL 3: ERROR 4: WARNING 5: NOTICE 6: INFORMATIONAL 7: DEBUG 8: NOTE',
  format TINYINT NOT NULL DEFAULT 0 COMMENT '0: text/plain 1: application/json',
  source TINYINT NOT NULL DEFAULT 0 COMMENT '0: INTERNAL 1: API 2: USER 3: SYSTEM',
  message TEXT,
  create_time_by varchar(255) NOT NULL DEFAULT '',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  delete_time timestamp NULL DEFAULT NULL,
  PRIMARY KEY(id),
  INDEX(object_id, create_time),
  INDEX(delete_time, id)
) ENGINE=InnoDB CHARSET utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE deleted_object_log (
  id BIGINT UNSIGNED NOT NULL,
  object_id BIGINT UNSIGNED NOT NULL,
  level TINYINT NOT NULL DEFAULT 0 COMMENT '0: EMERGENCY 1: ALERT 2: CRITICAL 3: ERROR 4: WARNING 5: NOTICE 6: INFORMATIONAL 7: DEBUG 8: NOTE',
  format TINYINT NOT NULL DEFAULT 0 COMMENT '0: text/plain 1: application/json',
  source TINYINT NOT NULL DEFAULT 0 COMMENT '0: INTERNAL 1: API 2: USER 3: SYSTEM',
  message TEXT,
  create_time_by varchar(255) NOT NULL DEFAULT '',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  delete_time timestamp NULL DEFAULT NULL,
  PRIMARY KEY(id),
  INDEX(object_id, delete_time, id)
) ENGINE=InnoDB CHARSET utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE object_relation_type (
  id SMALLINT UNSIGNED NOT NULL AUTO_INCREMENT,
  from_type_id SMALLINT UNSIGNED NOT NULL,
  to_type_id SMALLINT UNSIGNED NOT NULL,
  name VARCHAR(64) NOT NULL,
  description VARCHAR(512),
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp NULL DEFAULT NULL,
  delete_time timestamp NULL DEFAULT NULL,
  PRIMARY KEY(id),
  UNIQUE KEY(from_type_id, to_type_id, name),
  INDEX(delete_time, id)
) ENGINE=InnoDB CHARSET utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE deleted_object_relation_type (
  id SMALLINT UNSIGNED NOT NULL,
  from_type_id SMALLINT UNSIGNED NOT NULL,
  to_type_id SMALLINT UNSIGNED NOT NULL,
  name VARCHAR(64) NOT NULL,
  description VARCHAR(512),
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp NULL DEFAULT NULL,
  delete_time timestamp NULL DEFAULT NULL,
  PRIMARY KEY(id),
  INDEX(from_type_id, delete_time)
) ENGINE=InnoDB CHARSET utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE object_relation (
  from_object_id BIGINT UNSIGNED NOT NULL,
  relation_type_id SMALLINT UNSIGNED NOT NULL,
  to_object_id BIGINT UNSIGNED NOT NULL,
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  delete_time timestamp NULL DEFAULT NULL,
  PRIMARY KEY(from_object_id, relation_type_id, to_object_id),
  INDEX(from_object_id, relation_type_id, update_time),
  INDEX(from_object_id, update_time),
  INDEX(delete_time, from_object_id, relation_type_id, to_object_id)
) ENGINE=InnoDB CHARSET utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE deleted_object_relation (
  from_object_id BIGINT UNSIGNED NOT NULL,
  relation_type_id SMALLINT UNSIGNED NOT NULL,
  to_object_id BIGINT UNSIGNED NOT NULL,
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  delete_time timestamp NULL DEFAULT NULL,
  PRIMARY KEY(from_object_id, relation_type_id, to_object_id),
  INDEX(from_object_id, delete_time)
) ENGINE=InnoDB CHARSET utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE object_relation_meta (
  id MEDIUMINT UNSIGNED NOT NULL AUTO_INCREMENT,
  type_id SMALLINT UNSIGNED NOT NULL,
  name VARCHAR(256),
  value_type TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '1: STRING 2: INTEGER, 3: DOUBLE, 4: BOOLEAN',
  description VARCHAR(1024),
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  delete_time timestamp NULL DEFAULT NULL,
  PRIMARY KEY(id),
  UNIQUE KEY(type_id, name)
) ENGINE=InnoDB CHARSET utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE object_relation_meta_value (
  from_object_id BIGINT UNSIGNED NOT NULL,
  relation_id BIGINT UNSIGNED NOT NULL,
  to_object_id BIGINT UNSIGNED NOT NULL,
  meta_id MEDIUMINT UNSIGNED NOT NULL,
  value TEXT,
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp NULL DEFAULT NULL,
  delete_time timestamp NULL DEFAULT NULL,
  PRIMARY KEY(from_object_id, relation_id, to_object_id, meta_id),
  INDEX(delete_time, from_object_id, relation_id, to_object_id, meta_id)
) ENGINE=InnoDB CHARSET utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE deleted_object_relation_meta_value (
  from_object_id BIGINT UNSIGNED NOT NULL,
  relation_id BIGINT UNSIGNED NOT NULL,
  to_object_id BIGINT UNSIGNED NOT NULL,
  meta_id MEDIUMINT UNSIGNED NOT NULL,
  value TEXT,
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp NULL DEFAULT NULL,
  delete_time timestamp NULL DEFAULT NULL,
  PRIMARY KEY(from_object_id, relation_id, to_object_id, meta_id),
  INDEX(delete_time, from_object_id, relation_id, to_object_id)
) ENGINE=InnoDB CHARSET utf8mb4 COLLATE=utf8mb4_bin;
