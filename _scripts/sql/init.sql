CREATE TABLE Transactions
(
    ID         int auto_increment primary key,
    Account_ID int            null,
    Amount     decimal(20, 2) null,
    Currency   char(3)        null,
    OperationType varchar(20) null,
    CreatedAt  datetime       null
);

CREATE TABLE TransactionalOutbox
(
    ID        int auto_increment primary key,
    Body      json     null,
    Sent      bool     null,
    CreatedAt datetime null
);