--
-- PostgreSQL database dump
--

-- Dumped from database version 16.9 (Ubuntu 16.9-0ubuntu0.24.04.1)
-- Dumped by pg_dump version 16.9 (Ubuntu 16.9-0ubuntu0.24.04.1)

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
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner:
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: films; Type: TABLE; Schema: public; Owner: st1txh
--

CREATE TABLE public.films (
    film_id uuid DEFAULT gen_random_uuid() NOT NULL,
    title character varying(255) NOT NULL,
    description text,
    rating numeric(3,1),
    release_date timestamp without time zone NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT films_rating_check CHECK (((rating >= (0)::numeric) AND (rating <= (10)::numeric)))
);


ALTER TABLE public.films OWNER TO st1txh;

--
-- Name: user_film; Type: TABLE; Schema: public; Owner: st1txh
--

CREATE TABLE public.user_film (
    film_id uuid NOT NULL,
    user_id uuid NOT NULL
);


ALTER TABLE public.user_film OWNER TO st1txh;

--
-- Name: users; Type: TABLE; Schema: public; Owner: st1txh
--

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    date_of_birth timestamp without time zone NOT NULL,
    gender character varying(1) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT users_date_of_birth_check CHECK ((date_of_birth <= CURRENT_TIMESTAMP)),
    CONSTRAINT users_gender_check CHECK (((gender)::text = ANY ((ARRAY['М'::character varying, 'Ж'::character varying])::text[])))
);


ALTER TABLE public.users OWNER TO st1txh;

--
-- Data for Name: films; Type: TABLE DATA; Schema: public; Owner: st1txh
--

COPY public.films (film_id, title, description, rating, release_date, created_at, updated_at) FROM stdin;
7cc0535a-c5f3-4098-a097-1b03c1c2d7e4	2+1	Негр кайфовали жил жизнь, но лафа закончилась мать кукушка подкинула залетыша и теперь ему приходится как-то выживать и учиться жить эту блядскую жизнь	6.5	2017-06-01 08:53:16	2025-06-19 20:52:36.853014+06	2025-06-19 20:52:36.853014+06
7cade65e-4b39-44f6-b9df-a53b3ebb2fb7	Загадочная жизнь Алкаша Владика	Эта история начинается с далекого 2002 года, когда не было Iphone, выебона, денег, счастья у Владика. Но зато у него было пиво....	10.0	2025-06-10 08:53:16	2025-06-19 20:59:09.544764+06	2025-06-19 20:59:09.544764+06
1ac8ba9d-cf82-4edf-8257-ab67d480f205	КРУТОЙ ФИЛЬМ КОТОРЫЙ НУЖНО ИЗМЕНИТЬ ПРЯМО СЕЙЧАС БЛЯТЬ	Негр крадется у миллиардера яйцо фаберже. потом отрабатывает косяк, еще там ягуар был клевый и прикольный фильм в целом	5.4	2002-02-11 08:53:16	2025-06-19 20:51:07.070621+06	2025-06-19 23:21:46.862339+06
\.


--
-- Data for Name: user_film; Type: TABLE DATA; Schema: public; Owner: st1txh
--

COPY public.user_film (film_id, user_id) FROM stdin;
7cade65e-4b39-44f6-b9df-a53b3ebb2fb7	e1063aa0-136f-4f49-9aea-e064ad430291
1ac8ba9d-cf82-4edf-8257-ab67d480f205	729e952b-d959-46ac-81d1-441246b2062b
1ac8ba9d-cf82-4edf-8257-ab67d480f205	95d0a7e4-c8fb-43eb-a10c-f99949ad1959
1ac8ba9d-cf82-4edf-8257-ab67d480f205	7f714ba3-54f3-4f24-bd2a-b2011ffcf46b
7cc0535a-c5f3-4098-a097-1b03c1c2d7e4	7f714ba3-54f3-4f24-bd2a-b2011ffcf46b
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: st1txh
--

COPY public.users (id, name, email, date_of_birth, gender, created_at, updated_at) FROM stdin;
9a27752e-7c48-4c43-af90-8b0a48ea99f4	Кожкенов Данияр	st1txh03@example.com	2003-06-05 01:00:16	Ж	2025-06-19 15:55:42.901748+06	2025-06-19 20:38:25.947407+06
e1063aa0-136f-4f49-9aea-e064ad430291	Сосков Владик	sos@example.com	2002-06-23 01:00:16	М	2025-06-19 21:00:18.450313+06	2025-06-19 21:00:18.450313+06
c431d22f-782f-441b-a87d-08cecdcb77be	Мезин Сеёза	kartavyi@mail.com	2004-02-22 01:00:16	М	2025-06-19 21:01:02.416498+06	2025-06-19 21:01:02.416498+06
729e952b-d959-46ac-81d1-441246b2062b	Михеева Анна	cyganenok@mail.com	2004-01-01 01:00:16	Ж	2025-06-19 21:02:05.107225+06	2025-06-19 21:02:05.107225+06
95d0a7e4-c8fb-43eb-a10c-f99949ad1959	Мигель Анаров	chernysh@mail.com	2000-01-01 01:00:16	Ж	2025-06-20 13:52:49.55157+06	2025-06-20 13:52:49.55157+06
7f714ba3-54f3-4f24-bd2a-b2011ffcf46b	Супер Дупер	sxh@example.com	2003-06-05 01:00:16	М	2025-06-20 19:43:22.602589+06	2025-06-20 19:43:22.602589+06
\.


--
-- Name: films films_pkey; Type: CONSTRAINT; Schema: public; Owner: st1txh
--

ALTER TABLE ONLY public.films
    ADD CONSTRAINT films_pkey PRIMARY KEY (film_id);


--
-- Name: films films_title_key; Type: CONSTRAINT; Schema: public; Owner: st1txh
--

ALTER TABLE ONLY public.films
    ADD CONSTRAINT films_title_key UNIQUE (title);


--
-- Name: user_film user_film_pkey; Type: CONSTRAINT; Schema: public; Owner: st1txh
--

ALTER TABLE ONLY public.user_film
    ADD CONSTRAINT user_film_pkey PRIMARY KEY (film_id, user_id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: st1txh
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: st1txh
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: user_film user_film_film_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: st1txh
--

ALTER TABLE ONLY public.user_film
    ADD CONSTRAINT user_film_film_id_fkey FOREIGN KEY (film_id) REFERENCES public.films(film_id) ON DELETE CASCADE;


--
-- Name: user_film user_film_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: st1txh
--

ALTER TABLE ONLY public.user_film
    ADD CONSTRAINT user_film_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

