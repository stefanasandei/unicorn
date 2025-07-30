import requests
import json
import concurrent.futures
import time


def payload(msg):
    return json.dumps({
        "runtime": {
            "name": "go",
            "version": "3.12"
        },
        "project": {
            "entry": f'package main\nimport "fmt"\n\nfunc main() {{\n\tfmt.Print("{msg}")\n}}'
        },
        "process": {
            "time": "2s",
            "permissions": {
                "read": True
            }
        }
    })


def execute_request(request_num, msg):
    start_time = time.time()
    headers = {'Content-Type': 'application/json'}
    url = "http://localhost:3000/api/v1/execute"

    try:
        response = requests.post(url, headers=headers, data=payload(msg))
        data = json.loads(response.text)

        elapsed_time = time.time() - start_time
        out = data['output']['run']['stdout']
        prefix = ""
        if out != str(msg):
            prefix = "[!] "

        print(
            f"{prefix}Request {request_num}: status: {data['status']}, output: {out} (expected: {msg}), time: {elapsed_time:.2f} seconds")

        return elapsed_time
    except Exception as e:
        elapsed_time = time.time() - start_time
        print(f"Request {request_num}: Error executing request: {str(e)}, time: {elapsed_time:.2f} seconds")
        return elapsed_time


def run_in_parallel(num, ts):
    start_time = time.time()
    with concurrent.futures.ThreadPoolExecutor(max_workers=ts) as executor:
        futures = [executor.submit(execute_request, i, i) for i in range(1, num + 1)]

        # Wait for all threads to complete
        concurrent.futures.wait(futures)

    elapsed_times = [future.result() for future in futures]
    total_elapsed_time = time.time() - start_time

    mean_time = sum(elapsed_times) / num
    max_time = max(elapsed_times)
    min_time = min(elapsed_times)

    print(f"\nTotal elapsed time for {num} requests with {ts} threads: {total_elapsed_time:.2f} seconds")
    print(f"Mean time: {mean_time:.2f} seconds, Max time: {max_time:.2f} seconds, Min time: {min_time:.2f} seconds")


if __name__ == "__main__":
    num = 50
    _threads = 5
    run_in_parallel(num, _threads)

    # stats atm: 5s for 50 reqs => one worker can do 10 at the same time
    # about 0.5s / req
