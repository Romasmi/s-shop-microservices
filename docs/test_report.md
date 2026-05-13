# e2e test report

→ 1. Create User
POST http://arch.homework:8080/user  
200 OK ★ 849ms time ★ 450B↑ 415B↓ size ★ 9↑ 4↓ headers ★ 0 cookies
┌ ↑ raw ★ 162B
│ {
│   "username": "e2e_Matt_Carroll97",
│   "firstName": "Flavio",
│   "lastName": "Macejkovic",
│   "email": "Demetris.Gorczany@hotmail.com",
│   "phone": "201-423-2487"
│ }
└
┌ ↓ application/json ★ text ★ json ★ utf8 ★ 260B
│ {"id":"019e157e-5627-7c23-85eb-d2129b8aa142", "username":"e2e_Matt_Carroll97", "email":"Demetris.Gorczany@hotmail.com", "firstName":"Flavio", "lastName":"Macejkovic", "phone":"201-423-2487", "password":"7IvD3s!n.bdC", "createdAt":"2026-05-11T05:24:21.674441Z"}
└
prepare   wait   dns-lookup   tcp-handshake   transfer-start   download   process   total
26ms      11ms   15ms         813µs           800ms            19ms       496µs     875ms

✓  Status code is 200
POST http://arch.homework:8080/auth/login  
200 OK ★ 533ms time ★ 349B↑ 354B↓ size ★ 9↑ 3↓ headers ★ 0 cookies
┌ ↑ raw ★ 56B
│ {"login":"e2e_Matt_Carroll97","password":"7IvD3s!n.bdC"}
└
┌ ↓ application/json ★ text ★ json ★ utf8 ★ 245B
│ {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMDE5ZTE1N2UtNTYyNy03YzIzLTg1ZWItZDIxMjliOGFhMTQyIiwibG9naW4iOiJlMmVfTWF0dF9DYXJyb2xsOTciLCJleHAiOjE3Nzg0ODA2NjMsImlhdCI6MTc3ODQ3NzA2M30.w_m70rI7RCEh3TyymcGZNQyRRPuENaayaUiCYa7VL1k"}
└
prepare   wait   dns-lookup   tcp-handshake   transfer-start   download   process   total
1ms       2ms    (cache)      (cache)         526ms            3ms        95µs      534ms

┌
│ 'E2E Token obtained'
└

→ 2. Top Up Money
POST http://arch.homework:8080/billing/accounts/019e157e-5627-7c23-85eb-d2129b8aa142/topup  
200 OK ★ 29ms time ★ 620B↑ 221B↓ size ★ 10↑ 4↓ headers ★ 0 cookies
┌ ↑ raw ★ 22B
│ {
│   "amount": "1000"
│ }
└
┌ ↓ application/json ★ text ★ json ★ utf8 ★ 67B
│ {"userId":"019e157e-5627-7c23-85eb-d2129b8aa142", "balance":"1000"}
└
prepare   wait    dns-lookup   tcp-handshake   transfer-start   download   process   total
1ms       371µs   (cache)      (cache)         25ms             2ms        42µs      30ms

✓  Status code is 200
✓  Balance is 1000

→ 3. Make Success Order
POST http://arch.homework:8080/order  
200 OK ★ 176ms time ★ 616B↑ 281B↓ size ★ 10↑ 4↓ headers ★ 0 cookies
┌ ↑ raw ★ 72B
│ {
│   "userId": "019e157e-5627-7c23-85eb-d2129b8aa142",
│   "price": "600"
│ }
└
┌ ↓ application/json ★ text ★ json ★ utf8 ★ 126B
│ {"id":"2a66ee3c-589e-4dbf-8ee5-1fbc6b9398d0","userId":"019e157e-5627-7c23-85eb-d2129b8aa142","price":"600","status":"SUCCESS"}
└
prepare   wait    dns-lookup   tcp-handshake   transfer-start   download   process   total
1ms       481µs   (cache)      (cache)         171ms            3ms        53µs      176ms

✓  Status code is 200

→ 4. Check Balance (Decreased)
GET http://arch.homework:8080/billing/accounts/019e157e-5627-7c23-85eb-d2129b8aa142  
200 OK ★ 27ms time ★ 539B↑ 220B↓ size ★ 8↑ 4↓ headers ★ 0 cookies
┌ ↓ application/json ★ text ★ json ★ utf8 ★ 66B
│ {"userId":"019e157e-5627-7c23-85eb-d2129b8aa142", "balance":"400"}
└
prepare   wait    dns-lookup   tcp-handshake   transfer-start   download   process   total
2ms       438µs   (cache)      (cache)         23ms             2ms        316µs     28ms

✓  Status code is 200
✓  Balance is 400

→ 5. Check Notification (Success)
GET http://arch.homework:8080/notification/messages?userId=019e157e-5627-7c23-85eb-d2129b8aa142  
200 OK ★ 22ms time ★ 551B↑ 390B↓ size ★ 8↑ 4↓ headers ★ 0 cookies
┌ ↓ application/json ★ text ★ json ★ utf8 ★ 235B
│ {"messages":[{"id":"b303429d-28d8-4221-ad54-7f3e4106e1f5", "userId":"019e157e-5627-7c23-85eb-d2129b8aa142", "orderId":"2a66ee3c-589e-4dbf-8ee5-1fbc6b9398d0", "type":"ORDER_SUCCESS", "timestamp":"2026-05-11 05:24:23.394446 +0000 UTC"}]}
└
prepare   wait    dns-lookup   tcp-handshake   transfer-start   download   process   total
782µs     978µs   (cache)      (cache)         19ms             1ms        37µs      22ms

✓  Status code is 200
✓  Has ORDER_SUCCESS message

→ 6. Make Failed Order
POST http://arch.homework:8080/order  
200 OK ★ 101ms time ★ 617B↑ 281B↓ size ★ 10↑ 4↓ headers ★ 0 cookies
┌ ↑ raw ★ 73B
│ {
│   "userId": "019e157e-5627-7c23-85eb-d2129b8aa142",
│   "price": "1000"
│ }
└
┌ ↓ application/json ★ text ★ json ★ utf8 ★ 126B
│ {"id":"ab3a8c51-80bb-413f-b68e-d4427d980142","userId":"019e157e-5627-7c23-85eb-d2129b8aa142","price":"1000","status":"FAILED"}
└
prepare   wait    dns-lookup   tcp-handshake   transfer-start   download   process   total
1ms       268µs   (cache)      (cache)         97ms             2ms        588µs     102ms

✓  Status code is 200

→ 7. Check Balance (Unchanged)
GET http://arch.homework:8080/billing/accounts/019e157e-5627-7c23-85eb-d2129b8aa142  
200 OK ★ 26ms time ★ 539B↑ 220B↓ size ★ 8↑ 4↓ headers ★ 0 cookies
┌ ↓ application/json ★ text ★ json ★ utf8 ★ 66B
│ {"userId":"019e157e-5627-7c23-85eb-d2129b8aa142", "balance":"400"}
└
prepare   wait   dns-lookup   tcp-handshake   transfer-start   download   process   total
1ms       1ms    (cache)      (cache)         22ms             1ms        39µs      26ms

✓  Status code is 200
✓  Balance is still 400

→ 8. Check Notification (Failure)
GET http://arch.homework:8080/notification/messages?userId=019e157e-5627-7c23-85eb-d2129b8aa142  
200 OK ★ 14ms time ★ 551B↑ 610B↓ size ★ 8↑ 4↓ headers ★ 0 cookies
┌ ↓ application/json ★ text ★ json ★ utf8 ★ 455B
│ {"messages":[{"id":"b303429d-28d8-4221-ad54-7f3e4106e1f5", "userId":"019e157e-5627-7c23-85eb-d2129b8aa142", "orderId":"2a66ee3c-589e-4dbf-8ee5-1fbc6b9398d0", "type":"ORDER_SUCCESS", "timestamp":"2026-05-11 05:24:23.394446 +0000 UTC"}, {"id":"df003284-594c-4d55-8fe7-cb487d898e38", "userId":"019e157e
│ -5627-7c23-85eb-d2129b8aa142", "orderId":"ab3a8c51-80bb-413f-b68e-d4427d980142", "type":"ORDER_FAILED", "timestamp":"2026-05-11 05:24:23.62731 +0000 UTC"}]}
└
prepare   wait    dns-lookup   tcp-handshake   transfer-start   download   process   total
2ms       311µs   (cache)      (cache)         11ms             1ms        58µs      16ms

✓  Status code is 200
✓  Has ORDER_FAILED message

┌─────────────────────────┬─────────────────────┬─────────────────────┐
│                         │            executed │              failed │
├─────────────────────────┼─────────────────────┼─────────────────────┤
│              iterations │                   1 │                   0 │
├─────────────────────────┼─────────────────────┼─────────────────────┤
│                requests │                   9 │                   0 │
├─────────────────────────┼─────────────────────┼─────────────────────┤
│            test-scripts │                   8 │                   0 │
├─────────────────────────┼─────────────────────┼─────────────────────┤
│      prerequest-scripts │                   0 │                   0 │
├─────────────────────────┼─────────────────────┼─────────────────────┤
│              assertions │                  13 │                   0 │
├─────────────────────────┴─────────────────────┴─────────────────────┤
│ total run duration: 2.1s                                            │
├─────────────────────────────────────────────────────────────────────┤
│ total data received: 1.65kB (approx)                                │
├─────────────────────────────────────────────────────────────────────┤
│ average response time: 197ms [min: 14ms, max: 849ms, s.d.: 278ms]   │
├─────────────────────────────────────────────────────────────────────┤
│ average DNS lookup time: 15ms [min: 15ms, max: 15ms, s.d.: 0µs]     │
├─────────────────────────────────────────────────────────────────────┤
│ average first byte time: 188ms [min: 11ms, max: 800ms, s.d.: 266ms] │