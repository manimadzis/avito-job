CREATE DOMAIN MONEY_ AS NUMERIC(10, 2);

CREATE TYPE TRANSACTION_STATUS AS ENUM (
    'PENDING',
    'DONE'
);

CREATE TABLE IF NOT EXISTS "user" (
    id bigint PRIMARY KEY,
    balance MONEY_ DEFAULT 0,
    reserved_balance MONEY_ DEFAULT 0
);

CREATE TABLE IF NOT EXISTS "transaction" (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL,
    amount MONEY_ NOT NULL,
    "status" TRANSACTION_STATUS NOT NULL,
    service_id bigint,
    order_id bigint,
    "description" text,
    "timestamp" timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX transaction_user_service_order_amount_index ON "transaction" (user_id, amount, service_id, order_id);

CREATE TABLE IF NOT EXISTS "service" (
    id bigint PRIMARY KEY,
    "name" text
);

CREATE OR REPLACE PROCEDURE add_service (service_id bigint, "name" text)
LANGUAGE SQL
AS $$
    INSERT INTO "service" (id, "name")
        VALUES (service_id, "name")
    ON CONFLICT (id)
        DO UPDATE SET
            "name" = add_service."name";
$$;

CREATE OR REPLACE PROCEDURE replenish_balance (user_id bigint, amount MONEY_)
LANGUAGE plpgsql
AS $$
BEGIN
    BEGIN
        INSERT INTO "user" (id)
            VALUES (user_id);
    EXCEPTION
        WHEN unique_violation THEN
    END;
INSERT INTO "transaction" (user_id, amount, service_id, order_id, status)
    VALUES (user_id, amount, NULL, NULL, 'DONE');
            UPDATE
                "user"
            SET
                balance = balance + amount
            WHERE
                id = user_id;
END;
$$;

-- Raise no_data_found if user_id is unknown
CREATE OR REPLACE FUNCTION get_balance (user_id bigint)
    RETURNS MONEY_
    LANGUAGE plpgsql
    AS $$
DECLARE
    amount MONEY_;
BEGIN
    SELECT
        balance INTO STRICT amount
    FROM
        "user"
    WHERE
        id = user_id;
    RETURN amount;
END;
$$;

-- Raise exception with MESSAGE = NOT_ENOUGH_MONEY if amount greater than balance
CREATE OR REPLACE PROCEDURE reserve_money (user_id bigint, amount MONEY_, service_id bigint, order_id bigint, description text DEFAULT NULL)
LANGUAGE plpgsql
AS $$
BEGIN
    IF get_balance (user_id) < amount THEN
        RAISE EXCEPTION
            USING MESSAGE = 'NOT_ENOUGH_MONEY';
        END IF;
        UPDATE
            "user"
        SET
            reserved_balance = reserved_balance + amount,
            balance = balance - amount;
        INSERT INTO "transaction" (user_id, amount, service_id, order_id, "status", "description")
            VALUES (user_id, - amount, service_id, order_id, 'PENDING', "description");
END;
$$;

-- Raise exception with MESSAGE = RECOGNIZE_UNKNOWN_TRANSACTION if don't update any transaction
CREATE OR REPLACE PROCEDURE recognize_revenue (user_id bigint, amount MONEY_, service_id bigint, order_id bigint)
LANGUAGE plpgsql
AS $$
DECLARE
    affected_number int;
BEGIN
    WITH cte AS (
        UPDATE
            "transaction" t
        SET
            status = 'DONE'
        WHERE
            t.user_id = recognize_revenue.user_id
            AND t.service_id = recognize_revenue.service_id
            AND t.order_id = recognize_revenue.order_id
            AND t.amount = - recognize_revenue.amount
        RETURNING
            1
)
    SELECT
        count(*) INTO affected_number
    FROM
        cte;
    IF affected_number < 1 THEN
        RAISE EXCEPTION
            USING MESSAGE = 'RECOGNIZE_UNKNOWN_TRANSACTION';
        END IF;
        UPDATE
            "user"
        SET
            reserved_balance = reserved_balance - amount;
END;
$$;

CREATE OR REPLACE FUNCTION get_month_report (month int, year int)
    RETURNS TABLE (
        service_name text,
        revenue MONEY_)
    LANGUAGE SQL
    AS $$
    SELECT
        COALESCE(s.name, FORMAT('Service #%s', s.id)),
        t.amount
    FROM (
        SELECT
            service_id,
            - sum(amount) amount
        FROM
            "transaction"
        WHERE
            "status" = 'DONE'
            AND "timestamp" >= make_timestamp(year, month, 1, 0, 0, 0.0)
            AND "timestamp" < make_timestamp(year + (month + 1) / 12, (month + 1) % 12 + 1 * (month + 1) / 12, 1, 0, 0, 0.0)
        GROUP BY
            service_id) t
    LEFT JOIN "service" s ON t.service_id = s.id
WHERE
    service_id IS NOT NULL
$$;

CREATE OR REPLACE FUNCTION get_history_sorted_by_timestamp (user_id bigint, "offset" bigint, "limit" bigint, reverse boolean DEFAULT FALSE)
    RETURNS TABLE (
        "timestamp" timestamp,
        amount MONEY_,
        "description" text)
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF reverse THEN
        RETURN QUERY
        SELECT
            t."timestamp",
            t.amount,
            coalesce(t."description", 'No desctiption')
        FROM
            "transaction" t
        WHERE
            t.user_id = get_history_sorted_by_timestamp.user_id
        ORDER BY
            t."timestamp" ASC
        LIMIT "limit" OFFSET "offset";
    ELSE
        RETURN QUERY
        SELECT
            t."timestamp",
            t.amount,
            coalesce(t."description", 'No desctiption')
        FROM
            "transaction" t
        WHERE
            t.user_id = get_history_sorted_by_timestamp.user_id
        ORDER BY
            t."timestamp" DESC
        LIMIT "limit" OFFSET "offset";
    END IF;
END;
$$;