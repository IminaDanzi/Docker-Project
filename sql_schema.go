package dat

func GetSchemaSQL() string {
	return `
	  
	  CREATE TABLE IF NOT EXISTS migrations (
		name text NOT NULL,
		"time" timestamp without time zone DEFAULT now()
	  );
	  
	  DO $$
	  DECLARE
		  migration_name text := 'init_requestbody';
	  BEGIN
		  IF EXISTS (
			  SELECT
				  1
			  FROM
				  migrations
			  WHERE
				  name = migration_name) THEN
		  RETURN;
	  END IF;

	  CREATE TABLE requestbody (
		id SERIAL PRIMARY KEY,
		created timestamp with time zone DEFAULT now( ) NOT NULL,
		body text NOT NULL
	  );

	  INSERT INTO migrations (name) VALUES (migration_name);
	  END 
	  $$;
	`
}
