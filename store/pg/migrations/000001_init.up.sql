CREATE TABLE history(
    id serial,
    boardid VARCHAR(12),
    tradedate VARCHAR(10),
    shortname varchar(255),
    secid VARCHAR(50),
    numtrades INTEGER,
    value DECIMAL,
    open DECIMAL,
    low DECIMAL,
    high DECIMAL,
    legalcloseprice DECIMAL,
    waprice DECIMAL,
    close DECIMAL,
    volume INTEGER,
    marketprice2 DECIMAL,
    marketprice3 DECIMAL,
    admittedquote DECIMAL,
    mp2valtrd DECIMAL,
    marketprice3tradesvalue DECIMAL,
    admittedvalue DECIMAL,
    waval INTEGER,
    tradingsession INTEGER
);