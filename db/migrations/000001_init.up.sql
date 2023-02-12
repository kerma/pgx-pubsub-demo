CREATE TABLE IF NOT EXISTS events (
	id BIGSERIAL,
	created_at timestamp default now(),
	data jsonb
);

CREATE OR REPLACE FUNCTION notify_event() RETURNS trigger AS $$
DECLARE
BEGIN
	PERFORM pg_notify(NEW.data->>'kind', NEW.id::text);
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER notify_event AFTER INSERT ON events FOR EACH ROW EXECUTE PROCEDURE notify_event();

