# FastCoincidencesRestService

Run application with next command

```
docker-compose build
docker-compose up
```

Production mode
```
export APP_ENV=production 
docker-compose build
docker-compose up
```

App generating fake data automatically in first run. Size of load data 10 million records
Average speed response 10 - 15sec per request, but this can optimize with save data in memory

Testing computer: 
```
mac with Core i5, 8Gb DDR4, 128SSD
```

Get request: http://localhost:12345/1/2
Response:
```json
{ "dupes": true }
```

Get request: http://localhost:12345/1/3
Response:

```json
{ "dupes": false }
```

Get request: http://localhost:12345/2/1
Response:

```json
{ "dupes": true }
```

Get request: http://localhost:12345/2/3
Response:

```json
{ "dupes": true }
```

Get request: http://localhost:12345/3/2
Response:

```json
{ "dupes": true }
```

Get request: http://localhost:12345/1/4
Response:

```json
{ "dupes": false }
```

Get request: http://localhost:12345/3/1
Response:

```json
{ "dupes": false}
```

Get request: http://localhost:12345/1/1
Response:

```json
{ "dupes": true}
```
