# Testing

This directory contains some python scripts for testing sequential calls to the REST API. The `queries.py` file defines six functions, corresponding to the server's endpoints. Each function makes a request and verifies the result, then returns the status code, response payload (JSON), and a 'hint' (a string describing what happened).

- add_message(base_url, text, expect_status): POST /messages, expects 201, and an integer 'id' in the response payload
- get_message(base_url, id, expect_status, expect_text, expect_is_palindrome): GET /messages/{id}, expects 200, string 'text', boolean (but can be None) 'is_palindrome'
- update_message(base_url, id, text, expect_status): PUT /messages/{id}, expects 200, no data
- delete_message(base_url, id, expect_status): DELETE /messages/{id}, expects 204, no data
- get_all_messages(base_url, expect_status, expect_messages): GET /messages, expects 200, 'messages' array, with each item having 'id', 'text', and 'is_palindrome'
- delete_all_messages(base_url, expect_status): DELETE /messages, expects 204, no data

If the expect_status is different from what is normally considered successful, then the response payload is neither parsed nor validated. For example, attempting to get a message which doesn't exist should correctly return 404 with no payload.

## Environment

I wrote and ran these files on MacOs using Python 3.9.6 and the [requests](https://requests.readthedocs.io/en/latest/) package. Here are some commands that I found useful:

```shell
# establish a new python virtual env
python3 -m venv .venv 
# activate it
source .venv/bin/activate
# verify that it's running the correct python
which python 
# install the requests package
pip install requests 
# run a test (make sure go is running first)
python delay_zero.py 
# exit the virtual environment
deactivate
```

## Tests

_Caution_: The last step of both tests removes all messages (so the tests can be run repeatedly), so don't run them if there are messages you'd like to keep. However the first step of both tests verifies that there are no messages, so how we ended up here I'm not sure.

There are two test files: `delay_zero.py` and `delay_five.py`. The go server / REST API can be run with an artificial delay when determining if some text is a palindrome, by setting the environment variable S_DELAY to some number of seconds. `delay_zero.py` expects the go server to have no delay (aka `S_DELAY=0` or just don't set it). `delay_five.py` expects the go server to have a five second delay (aka `S_DELAY=5`): in most cases, is_palindrome will be None until five seconds after the initial add_message POST request.

Both tests raise an exception at the first sign of trouble, then immediately bail out.

## Output

From `delay_zero.py`:

1. GET /messages: OK
2. POST /messages (racecar): OK, got id 1
3. POST /messages (radar): OK, got id 2
4. POST /messages (A man, a plan, a canal, Panama): OK, got id 3
5. POST /messages (A Toyota's a Toyota): OK, got id 4
6. POST /messages (Drab as a fool, as aloof as a bard): OK, got id 5
7. POST /messages (Won't lovers revolt now?): OK, got id 6
8. POST /messages (Race fast, safe car): OK, got id 7
9. POST /messages (hello, world): OK, got id 8
10. POST /messages (racecars): OK, got id 9
11. POST /messages (Fear is the microbrewer): OK, got id 10
12. GET /messages/1: OK, is_palindrome True
13. GET /messages/8: OK, is_palindrome False
14. GET /messages: OK
15. DELETE /messages/1: OK
16. DELETE /messages/8: OK
17. DELETE /messages/1: OK, got status 404, as expected
18. DELETE /messages/8: OK, got status 404, as expected
19. PUT /messages/1 (abc): OK, got status 404, as expected
20. PUT /messages/8 (abc): OK, got status 404, as expected
21. GET /messages/1: OK, got status 404, as expected
22. GET /messages/8: OK, got status 404, as expected
23. PUT /messages/2 (racecars): OK
24. PUT /messages/9 (radar): OK
25. GET /messages/2: OK, is_palindrome False
26. GET /messages/9: OK, is_palindrome True
27. PUT /messages/2 (radar): OK
28. PUT /messages/9 (racecars): OK
29. GET /messages: OK
30. POST /messages (): OK, got id 11
31. GET /messages/11: OK, is_palindrome None
32. PUT /messages/11 (abcba): OK
33. GET /messages/11: OK, is_palindrome True
34. DELETE /messages: OK
35. GET /messages: OK

From `delay_five.py` (and `S_DELAY=5`):

1. GET /messages: OK
2. POST /messages (racecar): OK, got id 1
3. GET /messages/1: OK, is_palindrome None
4. Wait for 3 seconds
5. POST /messages (racecar): OK, got id 2
6. DELETE /messages/1: OK
7. GET /messages/1: OK, got status 404, as expected
8. Wait for 3 seconds
9. GET /messages/2: OK, is_palindrome True
10. PUT /messages/2 (hello, world): OK
11. Wait 1 second
12. GET /messages/2: OK, is_palindrome None
13. PUT /messages/2 (racecar): OK
14. Wait 1 second
15. PUT /messages/2 (racecars): OK
16. GET /messages: OK
17. Wait 4 seconds
18. GET /messages/2: OK, is_palindrome None
19. Wait 2 seconds
20. GET /messages/2: OK, is_palindrome False
21. DELETE /messages: OK
22. GET /messages: OK
