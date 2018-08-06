# ktoblzcheck-rest
Simple Rest Service to consume [ktoblzcheck](https://ktoblzcheck.sourceforge.net/)

### Build
- `make docker`
- `docker run -it -p 4040:4040 mllrsohn/ktoblzcheck-rest:latest`

### Example
`curl --data "iban=DE66200700240929139400&county=de" http://localhost:4040`
