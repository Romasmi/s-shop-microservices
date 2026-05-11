# Test Report

make test-postman
newman run docs/postman.json --verbose
newman

S Shop System

→ 1. Create User
POST http://arch.homework:8080/user  
200 OK ★ 725ms time ★ 437B↑ 402B↓ size ★ 9↑ 4↓ headers ★ 0 cookies
┌ ↑ raw ★ 149B
│ {
│   "username": "e2e_Juwan_Kozey",
│   "firstName": "Yolanda",
│   "lastName": "Becker",
│   "email": "Anya.Olson@hotmail.com",
│   "phone": "360-626-1031"
│ }
└
┌ ↓ application/json ★ text ★ json ★ utf8 ★ 247B
│ {"id":"019e16e5-0931-7106-9786-eab159a669fa", "username":"e2e_Juwan_Kozey", "email":"Anya.Olson@hotmail.com", "firstName":"Yolanda", "lastName":"Becker", "phone":"360-626-1031", "password":"@@vYi5ww!.Gf", "createdAt":"2026-05-11T11:56:09.394420Z"}
└
prepare   wait   dns-lookup   tcp-handshake   transfer-start   download   process   total
20ms      14ms   12ms         919µs           680ms            15ms       504µs     745ms

✓  Status code is 200
POST http://arch.homework:8080/auth/login  
200 OK ★ 468ms time ★ 346B↑ 350B↓ size ★ 9↑ 3↓ headers ★ 0 cookies
┌ ↑ raw ★ 53B
│ {"login":"e2e_Juwan_Kozey","password":"@@vYi5ww!.Gf"}
└
┌ ↓ application/json ★ text ★ json ★ utf8 ★ 241B
│ {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMDE5ZTE2ZTUtMDkzMS03MTA2LTk3ODYtZWFiMTU5YTY2OWZhIiwibG9naW4iOiJlMmVfSnV3YW5fS296ZXkiLCJleHAiOjE3Nzg1MDQxNzAsImlhdCI6MTc3ODUwMDU3MH0.3m7KI-lpU9ttX4mINXftYg1IcfVqedeXNjVRuHG_wDo"}
└
prepare   wait   dns-lookup   tcp-handshake   transfer-start   download   process   total
1ms       1ms    (cache)      (cache)         460ms            5ms        119µs     469ms

┌
│ 'E2E Token obtained'
└

→ 1.1 Setup Warehouse
POST http://arch.homework:8080/warehouse/products  
200 OK ★ 242ms time ★ 626B↑ 219B↓ size ★ 10↑ 4↓ headers ★ 0 cookies
┌ ↑ raw ★ 73B
│ {
│   "productId": "550e8400-e29b-41d4-a716-446655440000",
│   "count": 100
│ }
└
┌ ↓ application/json ★ text ★ json ★ utf8 ★ 65B
│ {"productId":"550e8400-e29b-41d4-a716-446655440000", "count":100}
└
prepare   wait    dns-lookup   tcp-handshake   transfer-start   download   process   total
1ms       617µs   (cache)      (cache)         238ms            2ms        113µs     243ms


→ 1.2 Setup Delivery
POST http://arch.homework:8080/delivery/couriers  
200 OK ★ 22ms time ★ 580B↑ 229B↓ size ★ 10↑ 4↓ headers ★ 0 cookies
┌ ↑ raw ★ 28B
│ {
│   "name": "Courier John"
│ }
└
┌ ↓ application/json ★ text ★ json ★ utf8 ★ 75B
│ {"courierId":"34e3e947-6acd-4e5d-bc21-e9ca0e547cce", "name":"Courier John"}
└
prepare   wait    dns-lookup   tcp-handshake   transfer-start   download   process   total
554µs     512µs   (cache)      (cache)         18ms             2ms        41µs      21ms


→ 2. Top Up Money
POST http://arch.homework:8080/billing/accounts/019e16e5-0931-7106-9786-eab159a669fa/topup  
200 OK ★ 19ms time ★ 616B↑ 221B↓ size ★ 10↑ 4↓ headers ★ 0 cookies
┌ ↑ raw ★ 22B
│ {
│   "amount": "1000"
│ }
└
┌ ↓ application/json ★ text ★ json ★ utf8 ★ 67B
│ {"userId":"019e16e5-0931-7106-9786-eab159a669fa", "balance":"1000"}
└
prepare   wait    dns-lookup   tcp-handshake   transfer-start   download   process   total
1ms       324µs   (cache)      (cache)         16ms             2ms        253µs     20ms

✓  Status code is 200
✓  Balance is 1000

→ 3. Make Success Order
POST http://arch.homework:8080/order  
200 OK ★ 510ms time ★ 682B↑ 281B↓ size ★ 10↑ 4↓ headers ★ 0 cookies
┌ ↑ raw ★ 141B
│ {
│   "userId": "019e16e5-0931-7106-9786-eab159a669fa",
│   "price": "600",
│   "productId": "550e8400-e29b-41d4-a716-446655440000",
│   "count": 1
│ }
└
┌ ↓ application/json ★ text ★ json ★ utf8 ★ 126B
│ {"id":"e398a19f-6bae-4b02-87ac-29a36380af43","userId":"019e16e5-0931-7106-9786-eab159a669fa","price":"600","status":"SUCCESS"}
└
prepare   wait   dns-lookup   tcp-handshake   transfer-start   download   process   total
1ms       1ms    (cache)      (cache)         507ms            1ms        364µs     511ms

✓  Status code is 200

→ 4. Check Balance (Decreased)
GET http://arch.homework:8080/billing/accounts/019e16e5-0931-7106-9786-eab159a669fa  
200 OK ★ 11ms time ★ 535B↑ 220B↓ size ★ 8↑ 4↓ headers ★ 0 cookies
┌ ↓ application/json ★ text ★ json ★ utf8 ★ 66B
│ {"userId":"019e16e5-0931-7106-9786-eab159a669fa", "balance":"400"}
└
prepare   wait    dns-lookup   tcp-handshake   transfer-start   download   process   total
1ms       319µs   (cache)      (cache)         8ms              1ms        299µs     12ms

✓  Status code is 200
✓  Balance is 400

→ 5. Check Notification (Success)
GET http://arch.homework:8080/notification/messages?userId=019e16e5-0931-7106-9786-eab159a669fa  
200 OK ★ 19ms time ★ 547B↑ 390B↓ size ★ 8↑ 4↓ headers ★ 0 cookies
┌ ↓ application/json ★ text ★ json ★ utf8 ★ 235B
│ {"messages":[{"id":"0ea21df4-229a-4c89-8fbf-d927ca0b6c99", "userId":"019e16e5-0931-7106-9786-eab159a669fa", "orderId":"e398a19f-6bae-4b02-87ac-29a36380af43", "type":"ORDER_SUCCESS", "timestamp":"2026-05-11 11:56:11.514986 +0000 UTC"}]}
└
prepare   wait   dns-lookup   tcp-handshake   transfer-start   download   process   total
1ms       1ms    (cache)      (cache)         15ms             1ms        39µs      19ms

✓  Status code is 200
✓  Has ORDER_SUCCESS message

→ 6. Make Failed Order
POST http://arch.homework:8080/order  
200 OK ★ 41ms time ★ 683B↑ 281B↓ size ★ 10↑ 4↓ headers ★ 0 cookies
┌ ↑ raw ★ 142B
│ {
│   "userId": "019e16e5-0931-7106-9786-eab159a669fa",
│   "price": "1000",
│   "productId": "550e8400-e29b-41d4-a716-446655440000",
│   "count": 1
│ }
└
┌ ↓ application/json ★ text ★ json ★ utf8 ★ 126B
│ {"id":"3f52307d-2eba-4bf4-b386-b27c998dc1be","userId":"019e16e5-0931-7106-9786-eab159a669fa","price":"1000","status":"FAILED"}
└
prepare   wait    dns-lookup   tcp-handshake   transfer-start   download   process   total
2ms       336µs   (cache)      (cache)         39ms             1ms        44µs      43ms

✓  Status code is 200

→ 7. Check Balance (Unchanged)
GET http://arch.homework:8080/billing/accounts/019e16e5-0931-7106-9786-eab159a669fa  
200 OK ★ 12ms time ★ 535B↑ 220B↓ size ★ 8↑ 4↓ headers ★ 0 cookies
┌ ↓ application/json ★ text ★ json ★ utf8 ★ 66B
│ {"userId":"019e16e5-0931-7106-9786-eab159a669fa", "balance":"400"}
└
prepare   wait   dns-lookup   tcp-handshake   transfer-start   download   process   total
2ms       2ms    (cache)      (cache)         5ms              3ms        702µs     14ms

✓  Status code is 200
✓  Balance is still 400

→ 8. Check Notification (Failure)
GET http://arch.homework:8080/notification/messages?userId=019e16e5-0931-7106-9786-eab159a669fa  
200 OK ★ 6ms time ★ 547B↑ 611B↓ size ★ 8↑ 4↓ headers ★ 0 cookies
┌ ↓ application/json ★ text ★ json ★ utf8 ★ 456B
│ {"messages":[{"id":"0ea21df4-229a-4c89-8fbf-d927ca0b6c99", "userId":"019e16e5-0931-7106-9786-eab159a669fa", "orderId":"e398a19f-6bae-4b02-87ac-29a36380af43", "type":"ORDER_SUCCESS", "timestamp":"2026-05-11 11:56:11.514986 +0000 UTC"}, {"id":"1a87958d-a076-4bb7-8ef0-5473c8cd4ecf", "userId":"019e16e5-0931-7106-97
│ 86-eab159a669fa", "orderId":"3f52307d-2eba-4bf4-b386-b27c998dc1be", "type":"ORDER_FAILED", "timestamp":"2026-05-11 11:56:11.650326 +0000 UTC"}]}
└
prepare   wait    dns-lookup   tcp-handshake   transfer-start   download   process   total
582µs     540µs   (cache)      (cache)         3ms              994µs      31µs      6ms

✓  Status code is 200
✓  Has ORDER_FAILED message

┌─────────────────────────┬─────────────────────┬────────────────────┐
│                         │            executed │             failed │
├─────────────────────────┼─────────────────────┼────────────────────┤
│              iterations │                   1 │                  0 │
├─────────────────────────┼─────────────────────┼────────────────────┤
│                requests │                  11 │                  0 │
├─────────────────────────┼─────────────────────┼────────────────────┤
│            test-scripts │                   8 │                  0 │
├─────────────────────────┼─────────────────────┼────────────────────┤
│      prerequest-scripts │                   0 │                  0 │
├─────────────────────────┼─────────────────────┼────────────────────┤
│              assertions │                  13 │                  0 │
├─────────────────────────┴─────────────────────┴────────────────────┤
│ total run duration: 2.4s                                           │
├────────────────────────────────────────────────────────────────────┤
│ total data received: 1.77kB (approx)                               │
├────────────────────────────────────────────────────────────────────┤
│ average response time: 188ms [min: 6ms, max: 725ms, s.d.: 247ms]   │
├────────────────────────────────────────────────────────────────────┤
│ average DNS lookup time: 12ms [min: 12ms, max: 12ms, s.d.: 0µs]    │
├────────────────────────────────────────────────────────────────────┤
│ average first byte time: 181ms [min: 3ms, max: 680ms, s.d.: 239ms] │
└────────────────────────────────────────────────────────────────────┘
romasmi@Anastasiias-MacBook-Pro s-shop-microservices % git push --force
Enumerating objects: 19, done.
Counting objects: 100% (19/19), done.
Delta compression using up to 12 threads
Compressing objects: 100% (10/10), done.
Writing objects: 100% (10/10), 805 bytes | 805.00 KiB/s, done.
Total 10 (delta 8), reused 0 (delta 0), pack-reused 0
remote: Resolving deltas: 100% (8/8), completed with 8 local objects.
To https://github.com/Romasmi/s-shop-microservices
+ 64f5c2a...bdd1732 HW8 -> HW8 (forced update)