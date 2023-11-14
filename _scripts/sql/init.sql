CREATE TABLE Transactions
(
    ID            int auto_increment primary key,
    Account_ID    int            null,
    Amount        decimal(20, 2) null,
    Currency      char(3)        null,
    OperationType varchar(20)    null,
    CreatedAt     datetime       null
);

CREATE TABLE Accounts
(
    ID        int auto_increment primary key,
    Document  varchar(50) null,
    CreatedAt datetime    null
);

CREATE TABLE TransactionalOutbox
(
    ID          int auto_increment primary key,
    Domain      varchar(50) null,
    Type        varchar(50) null,
    Body        json        null,
    Sent        bool        null,
    SentAt      datetime    null,
    CreatedAt   datetime    null
);
