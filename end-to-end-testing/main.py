from test_data import *
from queries import *

def go():
    log = logger()
    local_data = {}
    
    for p in palindromes:
        status, data, hint = add_message(BASE_URL, p, 201)
        local_data[data["id"]] = { "id": data["id"], "text": p, "is_palindrome": True }
        log(hint)
    
    for np in non_palindromes:
        status, data, hint = add_message(BASE_URL, np, 201)
        local_data[data["id"]] = { "id": data["id"], "text": np, "is_palindrome": False }
        log(hint)

    some_palindrome = [v for v in local_data.values() if v["is_palindrome"]][0]
    some_non_palindrome = [v for v in local_data.values() if not v["is_palindrome"]][0]

    status, data, hint = get_message(BASE_URL, some_palindrome["id"], 200, some_palindrome["text"], True)
    log(hint)

    status, data, hint = get_message(BASE_URL, some_non_palindrome["id"], 200, some_non_palindrome["text"], False)
    log(hint)

    status, data, hint = get_all_messages(BASE_URL, 200, local_data)
    log(hint)

    status, data, hint = delete_message(BASE_URL, some_palindrome["id"], 204)
    del local_data[some_palindrome["id"]]
    log(hint)

    status, data, hint = delete_message(BASE_URL, some_non_palindrome["id"], 204)
    del local_data[some_non_palindrome["id"]]
    log(hint)

    status, data, hint = delete_message(BASE_URL, some_palindrome["id"], 404)
    log(hint)

    status, data, hint = delete_message(BASE_URL, some_non_palindrome["id"], 404)
    log(hint)

    status, data, hint = update_message(BASE_URL, some_palindrome["id"], "abc", 404)
    log(hint)

    status, data, hint = update_message(BASE_URL, some_non_palindrome["id"], "abc", 404)
    log(hint)

    some_palindrome = [v for v in local_data.values() if v["is_palindrome"]][0]
    some_non_palindrome = [v for v in local_data.values() if not v["is_palindrome"]][0]

    status, data, hint = update_message(BASE_URL, some_palindrome["id"], some_non_palindrome["text"], 200)
    log(hint)

    status, data, hint = update_message(BASE_URL, some_non_palindrome["id"], some_palindrome["text"], 200)
    log(hint)

    status, data, hint = get_message(BASE_URL, some_palindrome["id"], 200, some_non_palindrome["text"], False)
    log(hint)

    status, data, hint = get_message(BASE_URL, some_non_palindrome["id"], 200, some_palindrome["text"], True)
    log(hint)

    status, data, hint = update_message(BASE_URL, some_palindrome["id"], some_palindrome["text"], 200)
    log(hint)

    status, data, hint = update_message(BASE_URL, some_non_palindrome["id"], some_non_palindrome["text"], 200)
    log(hint)

    status, data, hint = get_all_messages(BASE_URL, 200, local_data)
    log(hint)

    status, data, hint = add_message(BASE_URL, "", 201)
    none_palindrome_id = data["id"]
    local_data[none_palindrome_id] = { "id": none_palindrome_id, "text": "", "is_palindrome": None }
    log(hint)

    status, data, hint = get_message(BASE_URL, none_palindrome_id, 200, "", None)
    log(hint)

    status, data, hint = update_message(BASE_URL, none_palindrome_id, "abcba", 200)
    log(hint)

    status, data, hint = get_message(BASE_URL, none_palindrome_id, 200, "abcba", True)
    log(hint)

    status, data, hint = delete_all_messages(BASE_URL, 204)
    log(hint)

    local_data = {}
    status, data, hint = get_all_messages(BASE_URL, 200, local_data)
    log(hint)

if __name__ == '__main__':
    go()