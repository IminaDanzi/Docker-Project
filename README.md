# Docker-Project

This is a simplified backend that listens to POST request. Once it receives such a request it will push it to a PostgreSQL database.

To connect to the database, the server will the following environment variables:

DB_USER - the username.
DB_PASSWORD - the password.
DB_NAME - the database name.
DB_HOST - the database hostname.
DB_PORT - the port to connect to.
DB_SSL - SSL mode (example: require / disable).
