-- Create transactions table
CREATE TABLE transaction (
  customer_id INTEGER NOT NULL,
  amount INTEGER NOT NULL,
  type CHAR(1) NOT NULL,
  description TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create customer table
CREATE TABLE customer (
  id INTEGER NOT NULL,
  balance BIGINT NOT NULL,
  credit INTEGER NOT NULL,
  history JSONB NOT NULL DEFAULT jsonb_build_object()
);

CREATE INDEX idx_customer_id ON customer (id);

-- Create default customer balance
INSERT INTO "customer" ("id", "balance", "credit") VALUES
(1, 0, 100000);

INSERT INTO "customer" ("id", "balance", "credit") VALUES
(2, 0, 80000);

INSERT INTO "customer" ("id", "balance", "credit") VALUES
(3, 0, 1000000);

INSERT INTO "customer" ("id", "balance", "credit") VALUES
(4, 0, 10000000);

INSERT INTO "customer" ("id", "balance", "credit") VALUES
(5, 0, 500000);

-- Create reconcile balance trigger
CREATE
OR REPLACE FUNCTION reconcile_customer_balance() RETURNS TRIGGER LANGUAGE PLPGSQL AS $$ BEGIN
  UPDATE customer SET balance = balance + NEW.amount WHERE id = NEW.customer_id;

  RETURN NEW;
END;
$$;

CREATE TRIGGER reconcile_customer_balance
AFTER
INSERT
  ON transaction FOR EACH ROW EXECUTE PROCEDURE reconcile_customer_balance();