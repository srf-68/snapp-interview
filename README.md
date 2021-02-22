# snapp-interview
This is my solution for what Snapp team defined in this [task](https://github.com/AliKarami/interview-tasks/tree/master/term-frequency).

I developed this by Go language and using PostgreSQL as database.

To run this you should create a table in the PostgreSQL using following script:
```
CREATE SEQUENCE queries_id_seq;

CREATE TABLE queries (
    id integer NOT NULL DEFAULT nextval('queries_id_seq'),
	term character varying(1000),
	millitime bigint
);

ALTER SEQUENCE queries_id_seq
OWNED BY queries.id;
```
## How to use it
Send a search query to this service using the following REST API:
```
http://localhost:8585/index-search-query?query=.Please,%20email%20john.doe@foo.com%20by%2003-09,%20re:%20m37-xq.
```
and use the following REST API for get the N Top queries in the h hours ago:
```
localhost:8585/return-queries?hour=1&count=100
```