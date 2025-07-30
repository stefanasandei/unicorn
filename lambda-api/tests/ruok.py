import requests
import json


def payload(msg):
    return json.dumps({
        "runtime": {
            "name": "go",
            "version": "3.12"
        },
        "project": {
            "entry": f'package main\nimport "fmt"\n\nfunc main() {{\n\tfmt.Println("{msg}")\n}}'
        },
        "process": {
            "time": "2s",
            "permissions": {
                "read": True
            }
        }
    })


def send_req(msg):
    headers = {'Content-Type': 'application/json'}
    url = "http://localhost:3000/api/v1/execute"

    try:
        response = requests.post(url, headers=headers, data=payload(msg))
        data = json.loads(response.text)

        if response.status_code == 200 and data.get('status') == 'successful':
            print(f"Request was successful. Output: {data['output']['run']['stdout']}")
        else:
            print(f"Request failed. Status: {data['status']}")
    except Exception as e:
        print(f"Error executing request: {str(e)}")


if __name__ == "__main__":
    send_req("it works")
