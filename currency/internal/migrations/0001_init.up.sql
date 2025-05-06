CREATE TABLE currency_rates (
    date DATE    NOT NULL,
    code TEXT    NOT NULL,
    rate DOUBLE PRECISION NOT NULL,
    PRIMARY KEY (date, code)
);
