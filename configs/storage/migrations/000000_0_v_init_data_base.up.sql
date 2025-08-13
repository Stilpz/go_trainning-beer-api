CREATE TABLE IF NOT EXISTS "beers"(
    id bigint NOT NULL,
    name character varying(150) NOT NULL,
    brewery character varying(150) NOT NULL,
    country_code character varying(2) NOT NULL,
    price numeric(15, 6) NOT NULL,
    currency character varying(3) NOT NULL,
    created_at timestamp with time zone NOT NULL,
    PRIMARY KEY (id)
);

ALTER TABLE "beers" OWNER to postgres;