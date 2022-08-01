-- Table: public.orders

-- DROP TABLE IF EXISTS public.orders;

CREATE TABLE IF NOT EXISTS public.orders
(
    id bigint NOT NULL DEFAULT nextval('orders_id_seq'::regclass),
    maker character varying COLLATE pg_catalog."default" NOT NULL,
    kind smallint NOT NULL,
    assets jsonb NOT NULL,
    expired_at bigint NOT NULL,
    token_payment character varying COLLATE pg_catalog."default" NOT NULL,
    started_at bigint NOT NULL,
    base_price numeric(78,0) NOT NULL,
    ended_at bigint NOT NULL,
    ended_price numeric(78,0) NOT NULL,
    expected_state character varying COLLATE pg_catalog."default" NOT NULL,
    nonce bigint NOT NULL,
    market_fee_percentage bigint NOT NULL,
    signature character varying COLLATE pg_catalog."default" NOT NULL,
    hash character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT orders_pkey PRIMARY KEY (id)
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.orders
    OWNER to axie;
-- Index: orders_hash_idx

-- DROP INDEX IF EXISTS public.orders_hash_idx;

CREATE INDEX IF NOT EXISTS orders_hash_idx
    ON public.orders USING btree
    (hash COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;