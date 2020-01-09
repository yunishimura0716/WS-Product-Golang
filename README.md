# ws-product-golang

## Problems
- support counters by content selection and time, example counter Key `"sports:2020-01-08 22:01"`, Value `{views: 100, clicks: 4}`
- create go routine to upload counters to mock store every 5 sec
- global rate limit for stats handler
- try to leverage Golang's `channels` and/or `sync`
