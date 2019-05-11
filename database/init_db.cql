DROP KEYSPACE IF EXISTS watering;

CREATE KEYSPACE IF NOT EXISTS watering WITH REPLICATION={'class':'SimpleStrategy', 'replication_factor': 1};

USE watering;
CREATE TABLE IF NOT EXISTS watering.schedule (
  id UUID PRIMARY KEY ,

  end BIGINT,
  start BIGINT
);