CREATE DATABASE calendar_db;

CREATE USER calendar_user WITH encrypted password 'pass4calendarusr';

GRANT ALL PRIVILEGES ON database calendar_db TO calendar_user;