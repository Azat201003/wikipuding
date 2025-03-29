--
-- PostgreSQL database dump
--

-- Dumped from database version 16.3
-- Dumped by pg_dump version 16.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: wiki; Type: SCHEMA; Schema: -; Owner: wiki
--

CREATE SCHEMA wiki;


ALTER SCHEMA wiki OWNER TO wiki;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: articles; Type: TABLE; Schema: wiki; Owner: wiki
--

CREATE TABLE wiki.articles (
    id bigint NOT NULL,
    content text,
    title text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    creator_id bigint,
    base_article_id bigint,
    is_base boolean DEFAULT true
);


ALTER TABLE wiki.articles OWNER TO wiki;

--
-- Name: articles_id_seq; Type: SEQUENCE; Schema: wiki; Owner: wiki
--

ALTER TABLE wiki.articles ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME wiki.articles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: users; Type: TABLE; Schema: wiki; Owner: wiki
--

CREATE TABLE wiki.users (
    id bigint NOT NULL,
    username text,
    level bigint,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    token text,
    password text,
    is_staff boolean
);


ALTER TABLE wiki.users OWNER TO wiki;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: wiki; Owner: wiki
--

ALTER TABLE wiki.users ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME wiki.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Data for Name: articles; Type: TABLE DATA; Schema: wiki; Owner: wiki
--

COPY wiki.articles (id, content, title, created_at, updated_at, deleted_at, creator_id, base_article_id, is_base) FROM stdin;
18	May be it something interesting.	Something	2025-03-24 20:09:45.367439+03	2025-03-24 20:09:45.367439+03	\N	11	0	t
19	It is very <b>cool</b> and <b>interesting</b> page!	Cool article	2025-03-24 20:10:44.106973+03	2025-03-24 20:10:44.106973+03	\N	10	0	t
20	Much of people thinks, that python is programming language, but it <b>isn't</b>. Python is very big <b>snake</b>. It's green color and long length.	Python isn't snake	2025-03-24 21:20:08.985253+03	2025-03-24 21:20:08.985253+03	\N	10	0	t
26	Apple is circleful fruit. It may be red, yellow, green. It is sweety	Apple	2025-03-25 19:05:07.168412+03	2025-03-25 19:05:07.168412+03	\N	10	0	t
27	Orange is orange and it is beutiful. Orange is circleful fruit.	Orange	2025-03-25 19:07:21.19485+03	2025-03-25 19:07:21.19485+03	\N	11	0	t
23	<b>May be</b> it isnt something interesting.	Something	\N	\N	\N	10	21	f
24	May be it is interesting.	Something	\N	\N	\N	10	18	f
25	<i>May be</i> it is interesting.	Something	\N	\N	\N	11	24	f
21	<b>May be</b> it is something interesting.	Something	\N	\N	\N	11	18	f
0			\N	\N	\N	\N	0	f
30	I don't know what to writ, but i need some articles.	I don't know	2025-03-25 19:15:59.022365+03	2025-03-25 19:15:59.022365+03	\N	11	0	t
33	Apple is circle-shaped fruit. It may be red, yellow, green. It's sweety.	Apple	2025-03-28 01:14:29.421015+03	2025-03-28 01:14:29.421015+03	\N	10	26	f
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: wiki; Owner: wiki
--

COPY wiki.users (id, username, level, created_at, updated_at, deleted_at, token, password, is_staff) FROM stdin;
10	abeme	0	2025-03-24 18:44:18.913732+03	2025-03-24 18:44:18.913732+03	\N	token2	1234	\N
11	Coolmen	0	2025-03-24 19:25:20.096796+03	2025-03-24 19:25:20.096796+03	\N	token3	4321	\N
\.


--
-- Name: articles_id_seq; Type: SEQUENCE SET; Schema: wiki; Owner: wiki
--

SELECT pg_catalog.setval('wiki.articles_id_seq', 33, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: wiki; Owner: wiki
--

SELECT pg_catalog.setval('wiki.users_id_seq', 11, true);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: wiki; Owner: wiki
--

ALTER TABLE ONLY wiki.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_articles_deleted_at; Type: INDEX; Schema: wiki; Owner: wiki
--

CREATE INDEX idx_articles_deleted_at ON wiki.articles USING btree (deleted_at);


--
-- Name: idx_users_deleted_at; Type: INDEX; Schema: wiki; Owner: wiki
--

CREATE INDEX idx_users_deleted_at ON wiki.users USING btree (deleted_at);


--
-- PostgreSQL database dump complete
--

