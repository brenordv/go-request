# go-request!
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=brenordv_go-request&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=brenordv_go-request)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=brenordv_go-request&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=brenordv_go-request)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=brenordv_go-request&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=brenordv_go-request)

[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=brenordv_go-request&metric=bugs)](https://sonarcloud.io/summary/new_code?id=brenordv_go-request)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=brenordv_go-request&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=brenordv_go-request)
[![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=brenordv_go-request&metric=duplicated_lines_density)](https://sonarcloud.io/summary/new_code?id=brenordv_go-request)
[![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=brenordv_go-request&metric=sqale_index)](https://sonarcloud.io/summary/new_code?id=brenordv_go-request)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=brenordv_go-request&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=brenordv_go-request)

This is a simple application to make batch GET and POST requests.

## How to use
To use this application, you must provide a JSON with the desired parameters.

### JSON example
```json
{
  "get": {
    "numRequests": 100000,
    "url": "http://localhost:5001/sensors/v1/",
    "headers": {
      "Authorization": "Bearer eyJhbGciOiJBMTI4S1ciLCJlbmMiOiJBMTI4R0NNIiwidHlwIjoiSldUIiwiY3R5IjoiSldUIiwiemlwIjoiREVGIiwia2lkIjoiMXYZ.oid0BipNH3grGLDO_s55aIF7D-aYG45X.NsHC-_c82lCSjFOY.oAxgn0qXuAy3JnGx1wEdILNX0iWSvASTg2on4Q5vsoTOB2ynQkfpMPKPJ98PjKebb4zdNsgCbIbUdSVpvboH7JGcF1fQ8YRLiiWYAvQHLy684juhuzaTNECxRfAUu_ESnKV7nOjnUuPl3LHLiQ4StzEwBndoQ9Qff98QLl3atRyEH06T-YeVDedeJ2RthC7d8LZ8U5WyRyAqgAFyLetM3XCFz5fBnNXKwT0EHYmhhINyhRUijPZiQmE--GihlyfZzs9ryw5aMVEdU8P0gMzDINugesMyzhZlscniZjee50VKgMIoAwV_K7IGYFDuLnPojw28z57W3KcloZyeH-Ph8X9qwKb347PXX8kPV5uGztLzjKJnMsYOfnp7OMgh2ZkC7Re05oCXf2AyyDzZBqH6bImX3yk_mBV0jIO7fVMVtWRjoi2zZe6eoESeahnc19zfBJiO90m7ZKJJV5KqXa0uI6hF2YI_q-tijIw_zLE98h2kL9OPS5W4VATNobRDGhol_nFueOiT7Io8Lk4gg07XaZrDa8qq1N--Uxk5pV7CeTuLQZVhLxb5S5rNkU78IEClpvHfA-rT17KRmhUguC9mN0xjui-7fOURjcWUS6fkgqRkjbWDcabAxZTRzaYlWhU-lonq6UWRN5U0_XtYS9ub8v8eoU0sqCDFMP5GkUVh5HmFQSi3vbXRhhIo5_JPhi5HKoY-8UFU1L7RgL04kKocFdsCn3fzsSSCg56vv_85PjBrCF2-aO9llzccccyyVTbe7S-YxIcRPtknff0SRcdxIqQbTivB_vZ6HsTCgFWA4ZQ5QJhb21DI1TbKj0VWhuCv0JBphRPcbjm9ByZXi3QkweqYH3O--EC8n8MKdLfNTxGxPTI4cBErOHPM8uYjdr6vVmijcqHJ65wLyu1bv8U2wWJlhwakgyjOur1WmnHi1ykcpBcpYECDqax_ZJuUEBTdGL0nm5Q4BAhefm5fPUh4FMChtXWPhIFVR9ops7ojNjSH3VsF46vbZkcDseJldQRptie7iMOAz_nqQnEyS9sN04sGXrwgbtqEde_hLhAL0CyaKg8na1E4vxoVR2RMNpvHBVennUP_cQyvHIUvV4auYxq2Wpet2Z242efxQahejo5dvWdPoDEctE-dOBDHB7TjFoABHsfbs84pW0zCZpoBI5s7NSrB0y6dZ5mCLn7uEbEv7RWxNcyXQG4tcMxTd7zSLE5uKyOBP5ZLFOFuti_WrgoYxSDXOnMvf_shO_K48s87WEWFIIvBabno17hTUGwmVHzkBb00TmGkA7Brb6Fs2Q-fg7UXj0qqUEZcLmwW8glhoJ-OafbasG7hXTnWPQdhl2Z6586ofh6TvMHFw35bFK9iL9V3BBgH-5I6gbu3lSXPozEwlU7y66Waxtk7L1TJJIeV2L394O2KhfgzzK7nz_uPVcnwuefBWGRA2oNyeX84bvtLKJu6PKks7Fhy5qa4yPj8Sr5I45e97E1lYkNUJp6D5FeqGGJ1uNjQKgUH9xbgn-rO4czk4Z2cGVhMLq94T7NuTDpWka8AQXZ504GLpwW0IrDWoVBkshX9W2jPfxhaeZZ6Kchv749mWXH001QX1XKW.nmrf9SZVXgvEtbjPI_219w"
    }
  },
  "post": {
    "numRequests": 20000,
    "aggressiveMode": true,
    "url": "http://localhost:5001/sensors/v1/",
    "headers": {
      "Content-Type": "application/json"
    },
    "body": {
      "sensor_id": "S1",
      "reading_location_id": "L1",
      "reading_value": 42.007,
      "read_at": "2021-12-20T15:42:00.000Z"
    }
  }
}
```

The root key is the desired request method (GET or POST) and inside you must provide:
1. `numRequests`: Number of requests that will be made
2. `aggressiveMode`: If `false`, will wait for the current request to be finished before starting the next one. If `true`, requests will be made simultaneously.
3. `url`: The target address for the requests.
4. `headers`: Whatever you put here, will be sent as headers for each request.
5. `body`: (Only used for POST requests.) This will be the data sent in the request. 

### Running the requests
For the examples bellow, consider a json file named `my-request-config.json` placed in the same folder as 
the application.

### Get 
```shell
go-get.exe ./my-request-config.json
```

Example output:
```text
go-Request!::GET
Your session id is: 98c1136b-3ada-4aec-b1f5-c2783f0b19dd
Making GET requests 100% |█████████████████████████████████████████████| (5000/5000, 650 it/s)
Done! Elapsed time: 8.6664993s

Process finished with the exit code 0
```


### Post
```shell
go-post.exe ./my-request-config.json
```

Example output:
```text
go-Request!::POST
Your session id is: 59028850-6238-4601-9a92-bd51a313c667
Making POST requests 100% |█████████████████████████████████████████████| (5000/5000, 578 it/s)
Done! Elapsed time: 8.6464703s

Process finished with the exit code 0
```