--
-- PostgreSQL database dump
--

-- Dumped from database version 9.4.7
-- Dumped by pg_dump version 9.5.2

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner:
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner:
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: clusters; Type: TABLE; Schema: public; Owner: dbuser
--

CREATE TABLE IF NOT EXISTS clusters (
    cluster_id uuid NOT NULL,
    data json
);


ALTER TABLE clusters OWNER TO dbuser;

--
-- Name: clusters_checkins; Type: TABLE; Schema: public; Owner: dbuser
--

CREATE TABLE IF NOT EXISTS clusters_checkins (
    checkins_id bigint NOT NULL,
    cluster_id uuid,
    created_at timestamp without time zone,
    data json
);


ALTER TABLE clusters_checkins OWNER TO dbuser;

--
-- Name: clusters_checkins_checkins_id_seq; Type: SEQUENCE; Schema: public; Owner: dbuser
--

CREATE SEQUENCE clusters_checkins_checkins_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE clusters_checkins_checkins_id_seq OWNER TO dbuser;

--
-- Name: clusters_checkins_checkins_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dbuser
--

ALTER SEQUENCE clusters_checkins_checkins_id_seq OWNED BY clusters_checkins.checkins_id;


--
-- Name: versions; Type: TABLE; Schema: public; Owner: dbuser
--

CREATE TABLE IF NOT EXISTS versions (
    version_id bigint NOT NULL,
    component_name character varying(32),
    train character varying(24),
    version character varying(32),
    release_timestamp timestamp without time zone,
    data json
);


ALTER TABLE versions OWNER TO dbuser;

--
-- Name: versions_version_id_seq; Type: SEQUENCE; Schema: public; Owner: dbuser
--

CREATE SEQUENCE versions_version_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE versions_version_id_seq OWNER TO dbuser;

--
-- Name: versions_version_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dbuser
--

ALTER SEQUENCE versions_version_id_seq OWNED BY versions.version_id;


--
-- Name: checkins_id; Type: DEFAULT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY clusters_checkins ALTER COLUMN checkins_id SET DEFAULT nextval('clusters_checkins_checkins_id_seq'::regclass);


--
-- Name: version_id; Type: DEFAULT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY versions ALTER COLUMN version_id SET DEFAULT nextval('versions_version_id_seq'::regclass);


--
-- Data for Name: clusters; Type: TABLE DATA; Schema: public; Owner: dbuser
--

COPY clusters (cluster_id, data) FROM stdin;
4f0d9118-cfaa-4265-9335-cc31c6c7f15f	{"components":[{"component":{"description":"For testing only!","name":"deis-builder"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-controller"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-database"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-logger"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-minio"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-monitor-grafana"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-monitor-influxdb"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-monitor-stdout"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-registry"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-router"},"version":{"train":"","version":"v2-beta"}},{"component":{"name":"deis-workflow-manager"},"version":{"train":""}}],"id":"4f0d9118-cfaa-4265-9335-cc31c6c7f15f"}
\.


--
-- Data for Name: clusters_checkins; Type: TABLE DATA; Schema: public; Owner: dbuser
--

COPY clusters_checkins (checkins_id, cluster_id, created_at, data) FROM stdin;
1	4f0d9118-cfaa-4265-9335-cc31c6c7f15f	2016-05-31 23:16:25	{"components":[{"component":{"description":"For testing only!","name":"deis-builder"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-controller"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-database"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-logger"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-minio"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-monitor-grafana"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-monitor-influxdb"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-monitor-stdout"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-registry"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-router"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-workflow-manager"},"version":{"train":"","version":"v2-beta"}}],"id":"4f0d9118-cfaa-4265-9335-cc31c6c7f15f"}
2	4f0d9118-cfaa-4265-9335-cc31c6c7f15f	2016-05-31 23:25:51	{"components":[{"component":{"description":"For testing only!","name":"deis-builder"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-controller"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-database"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-logger"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-minio"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-monitor-grafana"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-monitor-influxdb"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-monitor-stdout"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-registry"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-router"},"version":{"train":"","version":"v2-beta"}},{"component":{"name":"deis-workflow-manager"},"version":{"train":""}}],"id":"4f0d9118-cfaa-4265-9335-cc31c6c7f15f"}
3	4f0d9118-cfaa-4265-9335-cc31c6c7f15f	2016-05-31 23:27:48	{"components":[{"component":{"description":"For testing only!","name":"deis-builder"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-controller"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-database"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-logger"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-minio"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-monitor-grafana"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-monitor-influxdb"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-monitor-stdout"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-registry"},"version":{"train":"","version":"v2-beta"}},{"component":{"description":"For testing only!","name":"deis-router"},"version":{"train":"","version":"v2-beta"}},{"component":{"name":"deis-workflow-manager"},"version":{"train":""}}],"id":"4f0d9118-cfaa-4265-9335-cc31c6c7f15f"}
\.


--
-- Name: clusters_checkins_checkins_id_seq; Type: SEQUENCE SET; Schema: public; Owner: dbuser
--

SELECT pg_catalog.setval('clusters_checkins_checkins_id_seq', 3, true);


--
-- Data for Name: versions; Type: TABLE DATA; Schema: public; Owner: dbuser
--

COPY versions (version_id, component_name, train, version, release_timestamp, data) FROM stdin;
45	deis-builder	beta	2.0.0-rc1	2016-04-30 00:00:00	{"notes":"rc 1 release"}
46	deis-controller	beta	2.0.0-rc1	2016-04-30 00:00:00	{"notes":"rc 1 release"}
47	deis-database	beta	2.0.0-rc1	2016-04-30 00:00:00	{"notes":"rc 1 release"}
48	deis-logger	beta	2.0.0-rc1	2016-04-30 00:00:00	{"notes":"rc 1 release"}
49	deis-logger-fluentd	beta	2.0.0-rc1	2016-04-30 00:00:00	{"notes":"rc 1 release"}
50	deis-registry	beta	2.0.0-rc1	2016-04-30 00:00:00	{"notes":"rc 1 release"}
51	deis-router	beta	2.0.0-rc1	2016-04-30 00:00:00	{"notes":"rc 1 release"}
52	deis-workflow-manager	beta	2.0.0-rc1	2016-04-30 00:00:00	{"notes":"rc 1 release"}
53	deis-minio	beta	2.0.0-rc1	2016-04-30 00:00:00	{"notes":"rc 1 release"}
54	deis-monitor-grafana	beta	2.0.0-rc1	2016-04-30 00:00:00	{"notes":"rc 1 release"}
55	deis-monitor-influxdb	beta	2.0.0-rc1	2016-04-30 00:00:00	{"notes":"rc 1 release"}
56	deis-monitor-stdout	beta	2.0.0-rc1	2016-04-30 00:00:00	{"notes":"rc 1 release"}
57	deis-monitor-telegraf	beta	2.0.0-rc1	2016-04-30 00:00:00	{"notes":"rc 1 release"}
58	deis-builder	beta	2.0.0-beta3	2016-04-28 00:00:00	{"notes":"beta 3 release"}
59	deis-controller	beta	2.0.0-beta3	2016-04-28 00:00:00	{"notes":"beta 3 release"}
60	deis-database	beta	2.0.0-beta3	2016-04-28 00:00:00	{"notes":"beta 3 release"}
61	deis-logger	beta	2.0.0-beta3	2016-04-28 00:00:00	{"notes":"beta 3 release"}
62	deis-logger-fluentd	beta	2.0.0-beta3	2016-04-28 00:00:00	{"notes":"beta 3 release"}
63	deis-registry	beta	2.0.0-beta3	2016-04-28 00:00:00	{"notes":"beta 3 release"}
64	deis-router	beta	2.0.0-beta3	2016-04-28 00:00:00	{"notes":"beta 3 release"}
65	deis-workflow-manager	beta	2.0.0-beta3	2016-04-28 00:00:00	{"notes":"beta 3 release"}
66	deis-minio	beta	2.0.0-beta3	2016-04-28 00:00:00	{"notes":"beta 3 release"}
67	deis-monitor-grafana	beta	2.0.0-beta3	2016-04-28 00:00:00	{"notes":"beta 3 release"}
68	deis-monitor-influxdb	beta	2.0.0-beta3	2016-04-28 00:00:00	{"notes":"beta 3 release"}
69	deis-monitor-stdout	beta	2.0.0-beta3	2016-04-28 00:00:00	{"notes":"beta 3 release"}
71	deis-monitor-telegraf	beta	2.0.0-beta3	2016-04-28 00:00:00	{"notes":"beta 3 release"}
\.


--
-- Name: versions_version_id_seq; Type: SEQUENCE SET; Schema: public; Owner: dbuser
--

SELECT pg_catalog.setval('versions_version_id_seq', 71, true);


--
-- Name: clusters_checkins_pkey; Type: CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY clusters_checkins
    ADD CONSTRAINT clusters_checkins_pkey PRIMARY KEY (checkins_id);


--
-- Name: clusters_pkey; Type: CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY clusters
    ADD CONSTRAINT clusters_pkey PRIMARY KEY (cluster_id);


--
-- Name: versions_component_name_train_version_key; Type: CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY versions
    ADD CONSTRAINT versions_component_name_train_version_key UNIQUE (component_name, train, version);


--
-- Name: versions_pkey; Type: CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY versions
    ADD CONSTRAINT versions_pkey PRIMARY KEY (version_id);


--
-- Name: public; Type: ACL; Schema: -; Owner: dbuser
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM dbuser;
GRANT ALL ON SCHEMA public TO dbuser;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--
