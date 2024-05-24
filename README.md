<h1 align="center">Transactional Outbox Pattern</h1>

## Description
- Hello 

## Architecture
- Hello

## Example

#### Creating account
`Request`
```bash
curl -i --request POST 'http://localhost:8080/v1/accounts' \
--header 'Content-Type: application/json' \
--data-raw '{
    "document": "07091058424"
}'
```

`Response`
```json
{
    "id":1
}
```

#### Creating transaction
`Request`
```bash
curl -i --request POST 'http://localhost:8080/v1/transactions' \
--header 'Content-Type: application/json' \
--data-raw '{
    "account_id": 1,
    "currency": "BRL",
    "operation_type": "CREDIT",
    "amount": 150.00
}'
```

`Response`
```json
{
    "id":1
}
```

## Author
- Gabriel Sabadini Facina - [GSabadini](https://github.com/GSabadini)

## License
Copyright © 2024 [GSabadini](https://github.com/GSabadini).
This project is [MIT](LICENSE) licensed.