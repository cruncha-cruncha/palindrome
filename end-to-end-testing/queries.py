import requests

def logger():
    step = 1
    def log(msg):
        nonlocal step
        print(f"{step}. {msg}")
        step += 1
    return log

def get_message(base_url, id, expect_status, expect_text, expect_is_palindrome):
    hint = f"GET /messages/{id}"
    resp = requests.get(f"{base_url}/messages/{id}")
    if resp.status_code != expect_status:
        raise Exception(f"{hint}: Expected status {expect_status} but got {resp.status_code}")
    if resp.status_code != 200:
        return resp.status_code, {}, f"{hint}: OK, got status {resp.status_code}, as expected"

    data = resp.json()
    if data.get("text") != expect_text:
        raise Exception(f"{hint}: Expected text {expect_text} but got {data.get('text')}")
    if data.get("is_palindrome") != expect_is_palindrome:
        raise Exception(f"{hint}: Expected is_palindrome {expect_is_palindrome} but got {data.get('is_palindrome')}")

    return resp.status_code, data, f"{hint}: OK, is_palindrome {data.get('is_palindrome')}"

def add_message(base_url, text, expect_status):
    hint = f"POST /messages ({text})"
    resp = requests.post(f"{base_url}/messages", json={"text": text})
    if resp.status_code != expect_status:
        raise Exception(f"{hint}: Expected status {expect_status} but got {resp.status_code}")
    if resp.status_code != 201:
        return resp.status_code, {}, f"{hint}: OK, got status {resp.status_code}, as expected"

    data = resp.json()
    if not isinstance(data.get("id"), int):
        raise Exception(f"{hint}: Expected an integer id but got {data.get('id')}")

    return resp.status_code, data, f"{hint}: OK, got id {data.get('id')}"

def update_message(base_url, id, text, expect_status):
    hint = f"PUT /messages/{id} ({text})"
    resp = requests.put(f"{base_url}/messages/{id}", json={"text": text})
    if resp.status_code != expect_status:
        raise Exception(f"{hint}: Expected status {expect_status} but got {resp.status_code}")
    if resp.status_code != 200:
        return resp.status_code, {}, f"{hint}: OK, got status {resp.status_code}, as expected"
    
    # no data returned

    return resp.status_code, {}, f"{hint}: OK"

def delete_message(base_url, id, expect_status):
    hint = f"DELETE /messages/{id}"
    resp = requests.delete(f"{base_url}/messages/{id}")
    if resp.status_code != expect_status:
        raise Exception({}, f"{hint}: Expected status {expect_status} but got {resp.status_code}")
    if resp.status_code != 204:
        return resp.status_code, {}, f"{hint}: OK, got status {resp.status_code}, as expected"

    # no data returned

    return resp.status_code, {}, f"{hint}: OK"

def get_all_messages(base_url, expect_status, expect_messages):
    hint = f"GET /messages"
    resp = requests.get(f"{base_url}/messages")
    if resp.status_code != expect_status:
        raise Exception(f"{hint}: Expected status {expect_status} but got {resp.status_code}")
    if resp.status_code != 200:
        return resp.status_code, {}, f"{hint}: OK, got status {resp.status_code}, as expected"

    data = resp.json()
    if data.get("messages") is None or not isinstance(data.get("messages"), list):
        raise Exception(f"{hint}: Expected messages to be a list, but got {data.get('messages')}")

    messages = data.get("messages")
    for message in messages:
        if not isinstance(message.get("id"), int):
            raise Exception(f"{hint}: Expected id to be an int, but got {message.get('id')}")
        
        id = message.get("id")
        expected_message = expect_messages.get(id)
        if expected_message is None:
            raise Exception(f"{hint}: Unexpected message with id {id}")

        if message.get("text") != expected_message.get("text"):
            raise Exception(f"{hint}: Expected text {expected_message.get('text')} but got {message.get('text')}, id {id}")
        if message.get("is_palindrome") != expected_message.get("is_palindrome"):
            raise Exception(f"{hint}: Expected is_palindrome {expected_message.get('is_palindrome')} but got {message.get('is_palindrome')}, id {id}")
    
    if len(messages) != len(expect_messages):
        raise Exception(f"{hint}: Expected {len(expect_messages)} messages but got {len(messages)}")
    ids = [message.get("id") for message in messages]
    if len(ids) != len(set(ids)):
        raise Exception(f"{hint}: Expected unique ids, but got duplicates")
    if ids != sorted(ids):
        raise Exception(f"{hint}: Expected sorted ids ({sorted(ids)}), but got {ids}")

    return resp.status_code, data, f"{hint}: OK"

def delete_all_messages(base_url, expect_status):
    hint = f"DELETE /messages"
    resp = requests.delete(f"{base_url}/messages")
    if resp.status_code != expect_status:
        raise Exception(f"{hint}: Expected status {expect_status} but got {resp.status_code}")
    if resp.status_code != 204:
        return resp.status_code, {}, f"{hint}: OK, got status {resp.status_code}, as expected"

    # no data returned

    return resp.status_code, {}, f"{hint}: OK"
    

