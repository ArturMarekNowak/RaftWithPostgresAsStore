-- yeah, I dont like SQL
DROP DATABASE "Raft1";
DROP DATABASE "Raft2";
DROP DATABASE "Raft3";
CREATE DATABASE "Raft1";
CREATE DATABASE "Raft2";
CREATE DATABASE "Raft3";

SELECT * FROM public.fsms;
SELECT * FROM public.snapshots;
SELECT * FROM public.logs;
SELECT * FROM public.stable_logs;

SHOW max_connections;
ALTER SYSTEM SET max_connections TO '10000';