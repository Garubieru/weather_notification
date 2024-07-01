CREATE TABLE account (
  id VARCHAR(24) PRIMARY KEY NOT NULL,
  name VARCHAR(100) NOT NULL,
  username VARCHAR(30) NOT NULL UNIQUE,
  email VARCHAR(320) NOT NULL UNIQUE,
  phone VARCHAR(20) NOT NULL,
  password CHAR(60) NOT NULL
);

CREATE TABLE session (
  id VARCHAR(24) PRIMARY KEY,
  account_id VARCHAR(24),
  expire_date DATETIME,
  INDEX idx_account_id (account_id),
  FOREIGN KEY (account_id) REFERENCES account(id) ON DELETE CASCADE
);

CREATE TABLE city (
  id VARCHAR(24) PRIMARY KEY,
  external_id VARCHAR(20) NOT NULL,
  state_code VARCHAR(2) NOT NULL,
  name VARCHAR(100) NOT NULL
);

CREATE TABLE notification_schedule (
  id VARCHAR(24) PRIMARY KEY,
  scheduled_date DATETIME NOT NULL,
  interval_in_days TINYINT UNSIGNED NOT NULL,
  hour TINYINT UNSIGNED NOT NULL,
  city_id VARCHAR(24) NOT NULL,
  is_coastal_city BOOL,
  active BOOL NOT NULL,
  method VARCHAR(20) NOT NULl,
  account_id VARCHAR(24) NOT NULL,
  INDEX idx_scheduled_date (scheduled_date),
  INDEX idx_account_id (account_id),
  FOREIGN KEY (account_id) REFERENCES account(id) ON DELETE CASCADE,
  FOREIGN KEY (city_id) REFERENCES city(id) ON DELETE CASCADE
);

CREATE TABLE topic_partition_offset (
  id INT AUTO_INCREMENT PRIMARY KEY,
  topic VARCHAR(100),
  offset INT,
  last_offset INT
);

CREATE TABLE notification (
  id VARCHAR(24) PRIMARY KEY,
  payload JSON NOT NULL,
  account_id VARCHAR(24) NOT NULL,
  seen BOOL NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_account_id (account_id),
  FOREIGN KEY (account_id) REFERENCES account(id) ON DELETE CASCADE
);

