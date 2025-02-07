I'm using Python 3.9.6. In the last part of this test, it removes all messages, so be careful if you want to keep messages.

```
python3 -m venv .venv
source .venv/bin/activate
which python
python main.py
deactivate
```

Output:
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