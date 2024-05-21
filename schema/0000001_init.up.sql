CREATE OR REPLACE FUNCTION make_uid() RETURNS text AS $$
DECLARE
new_uid text;
    done bool;
BEGIN
    done := false;
    WHILE NOT done LOOP
        new_uid := md5(''||now()::text||random()::text);
        done := NOT exists(SELECT 1 FROM wallets WHERE id=new_uid);
END LOOP;
RETURN new_uid;
END;
$$ LANGUAGE PLPGSQL VOLATILE;

CREATE TABLE IF NOT EXISTS wallets
(
    id TEXT DEFAULT make_uid()::text NOT NULL UNIQUE,
    balance INTEGER DEFAULT 0 CHECK (balance >= 0) NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions
(
    time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    from_wallet_id TEXT NOT NULL REFERENCES wallets(id),
    to_wallet_id TEXT NOT NULL REFERENCES wallets(id),
    amount INTEGER NOT NULL CHECK (amount > 0)
);