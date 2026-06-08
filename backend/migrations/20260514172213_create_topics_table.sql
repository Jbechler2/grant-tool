-- +goose Up
CREATE TABLE topics (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  grant_writer_id UUID NOT NULL,
  label VARCHAR(255) NOT NULL,
  CONSTRAINT unique_topics_labels
    UNIQUE (grant_writer_id, label),
  CONSTRAINT fk_users_topics
    FOREIGN KEY(grant_writer_id)
      REFERENCES users(id)
      ON DELETE CASCADE
);  

CREATE INDEX idx_topics_id ON topics(id);

CREATE TABLE grants_topics (
  topic_id UUID NOT NULL,
  grant_id UUID NOT NULL,
  PRIMARY KEY (topic_id, grant_id),
  CONSTRAINT fk_grants_topics
    FOREIGN KEY(grant_id)
      REFERENCES grants(id)
      ON DELETE CASCADE,
  CONSTRAINT fk_topics_grants
    FOREIGN KEY(topic_id)
      REFERENCES topics(id)
      ON DELETE CASCADE
);

CREATE INDEX idx_grants_topics_grant_id ON grants_topics(grant_id);

CREATE TABLE clients_topics (
  topic_id UUID NOT NULL,
  client_id UUID NOT NULL,
  PRIMARY KEY (topic_id, client_id),
  CONSTRAINT fk_clients_topics
    FOREIGN KEY(client_id)
      REFERENCES clients(id)
      ON DELETE CASCADE,
  CONSTRAINT fk_topics_clients
    FOREIGN KEY(topic_id)
      REFERENCES topics(id)
      ON DELETE CASCADE
);

CREATE INDEX idx_clients_topics_client_id ON clients_topics(client_id);


-- +goose Down
DROP TABLE IF EXISTS clients_topics;
DROP TABLE IF EXISTS grants_topics;
DROP TABLE IF EXISTS topics;


