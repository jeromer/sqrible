CREATE TABLE users (
    id bigint NOT NULL,
    uid text NOT NULL,
    password text NOT NULL,
    email text NOT NULL,
    username text DEFAULT ''::text NOT NULL,
    enabled boolean DEFAULT false NOT NULL,
    email_confirmed boolean DEFAULT false NOT NULL,
    creation_date timestamp with time zone DEFAULT now() NOT NULL,
    update_date timestamp with time zone DEFAULT now() NOT NULL,
    avatar_file_name text NULL,
    twitter text NULL,
    facebook text NULL,
    gplus text NULL
);

CREATE SEQUENCE users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE users_id_seq OWNED BY users.id;

ALTER TABLE ONLY users ALTER COLUMN id SET DEFAULT nextval('users_id_seq'::regclass);

ALTER TABLE ONLY users
    ADD CONSTRAINT users_email_key UNIQUE (email);

ALTER TABLE ONLY users
    ADD CONSTRAINT users_uid_key UNIQUE (uid);

ALTER TABLE ONLY users
    ADD CONSTRAINT users_pk PRIMARY KEY (id);

ALTER TABLE ONLY users
    ADD CONSTRAINT users_username_key UNIQUE (username);
