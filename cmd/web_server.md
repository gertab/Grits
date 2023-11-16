# Web-server

todo give example to start interpreter as a web server

The web-server accepts connections via websockets on `ws://localhost:8081/ws`. It accepts the following request message in JSON format to compile and execute a *program*: `{"type": "compile_program","program_to_compile": program}`. The following is an example where two processes will be spawned and interact with each other:

```json
{
    "type": "compile_program",
    "program_to_compile": 
        "prc[pid1] = send self<pid3, self>
         prc[pid2] = <a, b> <- recv pid1; close self"
}
```

# Reply

After the request to compile, the web-server sends replies that indicate an *error*, an updated process configuration or an updated list of transitions.

## Type "error"

Contains an error message.

```json
{
    "type": "error",
    "error_message": "Syntax error"
}
```

## Type "processes_updated"

When the process configuration changes, the new list of process is sent, including the links between the different processes. The following is an example.

```json
{
    "type": "processes_updated",
    "payload": {
        "processes": [
            {
                "id": "1",
                "providers": [
                    "pid1[1]"
                ],
                "body": "send self<pid3,self>"
            },
            {
                "id": "2",
                "providers": [
                    "pid2[2]"
                ],
                "body": "<a,b> <- recv pid1[1]; close self"
            }
        ],
        "links": [
            {
                "source": "2",
                "destination": "1"
            }
        ]
    }
}
```

## Type "rules_updated"

When processes transition, the list of the transitions is sent as a message.

```json
{
    "type": "rules_updated",
    "rules": [
        {
            "id": "0",
            "providers": [
                "pid2[2]"
            ],
            "rule": "SND"
        }
    ]
}
```
