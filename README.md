# ktoblzcheck-rest
Simple Rest Service to consume [ktoblzcheck](https://ktoblzcheck.sourceforge.net/)

### Example
- `docker run -it -p 4040:4040 steffenmllr/ktoblzcheck-rest:latest`
- `curl --data "iban=DE66200700240929139400&county=de" http://localhost:4040`
