from test_data import *
from queries import *
import time

# run with server like: S_DELAY=5 go run .

def go():
    log = logger()
    local_data = {}
    
    status, data, hint = get_all_messages(BASE_URL, 200, {})
    log(hint)

    first_palindrome_text = palindromes[0]
    status, data_1, hint = add_message(BASE_URL, first_palindrome_text, 201)
    log(hint)

    status, data_2, hint = get_message(BASE_URL, data_1["id"], 200, first_palindrome_text, None)
    log(hint)

    log("Wait for 3 seconds")
    time.sleep(3)

    status, data_3, hint = add_message(BASE_URL, first_palindrome_text, 201)
    log(hint)

    status, data_4, hint = delete_message(BASE_URL, data_1["id"], 204)
    log(hint)

    status, data_5, hint = get_message(BASE_URL, data_1["id"], 404, "", None)
    log(hint)

    log("Wait for 3 seconds")
    time.sleep(3)
    
    status, data, hint = get_message(BASE_URL, data_3["id"], 200, first_palindrome_text, True)
    log(hint)

    first_non_palindrome_text = non_palindromes[0]
    status, data, hint = update_message(BASE_URL, data_3["id"], first_non_palindrome_text, 200)
    log(hint)

    log("Wait 1 second")
    time.sleep(1)

    status, data, hint = get_message(BASE_URL, data_3["id"], 200, first_non_palindrome_text, None)
    log(hint)

    status, data, hint = update_message(BASE_URL, data_3["id"], first_palindrome_text, 200)
    log(hint)

    log("Wait 1 second")
    time.sleep(1)

    second_non_palindrome_text = non_palindromes[1]
    status, data, hint = update_message(BASE_URL, data_3["id"], second_non_palindrome_text, 200)
    log(hint)

    status, data, hint = get_all_messages(BASE_URL, 200, { data_3["id"]: { "id": data_3["id"], "text": second_non_palindrome_text, "is_palindrome": None } })
    log(hint)

    log("Wait 4 seconds")
    time.sleep(4)

    status, data, hint = get_message(BASE_URL, data_3["id"], 200, second_non_palindrome_text, None)
    log(hint)

    log("Wait 2 seconds")
    time.sleep(2)

    status, data, hint = get_message(BASE_URL, data_3["id"], 200, second_non_palindrome_text, False)
    log(hint)

    status, data, hint = delete_all_messages(BASE_URL, 204)
    log(hint)

    status, data, hint = get_all_messages(BASE_URL, 200, {})
    log(hint)

if __name__ == '__main__':
    go()