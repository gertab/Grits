# Phi

## Code Structure

### General Folder Structure

```text
phi/
├─ cmd/
├─ parser/
├─ process/
```

There are three main components.

- The `cmd` folder contains the entry point to either execute code from a file/string (`main.go`), or initiate a web-server (`web.go`) to compile and execute a program using an external interface.  
- The `parser` folder contains the parser *package* which processes a string and outputs a list of processes, ready to be executed.
- The `process` folder contains the *process* package which executes some given processes. It has several components:
  - `process/runtime.go`: Entry point for the runtime environment. Sets up the processes, channels and monitors and initiates the execution.
  - `process/form.go`: Contains the different forms that a process can take. Referred to as the abstract syntax tree of the processes.
  - `process/transition.go`: Defines how each form should execute.
  - `process/servers.go`: Sets up the concurrent servers (e.g. a monitor) that monitor or control the execution of the processes. Used to inform the web-server about the state of each process.

### Some examples

There are some examples in the `cmd/examples` folder.

### Other information

To slow down execution speed:
Set the delay property of the `RuntimeEnvironment` to a longer duration, e.g. `1000 * time.Millisecond`.

Features not currently implemented:

- [ ] Controller server to choose which transition do do

## Developing the interface (Matthieu)

### Web-server

The web-server accepts connections via websockets on `ws://localhost:8081/ws`. It accepts the following request message in JSON format to compile and execute a *program*: `{"type": "compile_program","program_to_compile": program}`. The following is an example where two processes will be spawned and interact with each other:

```json
{
    "type": "compile_program",
    "program_to_compile": 
        "prc[pid1]: send self<pid3, self>
         prc[pid2]: <a, b> <- recv pid1; close self"
}
```

After the request to compile, the web-server sends replies that indicate an *error*, an updated process configuration or an updated list of transitions.

#### Type "error"

Contains an error message.

```json
{
    "type": "error",
    "error_message": "Syntax error"
}
```

#### Type "processes_updated"

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

#### Type "rules_updated"

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

### Useful Links

- D3.js: [http://d3js.org/](http://d3js.org/)

## Grammar

```text
<prog> ::= <statement>*

<statement> ::= type <label> = <type>
              | let <label> ( [<parameters>] ) : <type> = <term>
              | prc '[' <name> ']' : <type> = <term>

<parameters> ::= <name> : <type> [ , <parameters> ]

<type> ::= <label>
         | 1
         | + { <branch_type> }
         | & { <branch_type> }
         | <type> * <type>
         | <type> -o <type>
         | ( <type> )

<branch_type> ::= <label> : <type> [ , <branch_type> ]

<term> ::= send <name> '<' <name> , <name> '>' 
        | '<' <name> , <name> '>' <- recv <name> ; <term>
        | <name> . <label> '<' <name> '>' 
        | case <name> ( <branches> )
        | <name> <- <polarity> new ( <term> ) ; <term>
        | <polarity> <label> ( [<names>] )
        | <polarity> fwd <names> <names>
        | '<' <name> , <name> '>' <- <polarity> split <name> ; <term>
        | close <name>
        | wait <name> ; term
        | cast <name> '<' <name> '>'
        | <name> <- shift <name> ; <term>

<polarity> ::= +
             | -

<branches> ::= <label> '<' <name> => term [ '|' <branches> ]

<names> ::= <name> [ ',' <names> ]

Others:
    <name> and <label> are any alpha-numeric label, where name represents a channel and label represents a labelled choice
    // Single line comments
    /* Multi line comments */
    whitespace is ignored

```

## Goals

Similar projects:

[http://www.emanueledosualdo.com/stargazer/?gist=6e54093b297c0f9df01d0c82f65b89f6](http://www.emanueledosualdo.com/stargazer/?gist=6e54093b297c0f9df01d0c82f65b89f6)

[https://www.nomos-lang.org/admin/interface](https://www.nomos-lang.org/admin/interface)

[https://github.com/bzhan/mars/tree/master/hhlpy](https://github.com/bzhan/mars/tree/master/hhlpy)

Useful features to have:

- [ ] Use network/hierarchy in D3 to display all the processes and the links between them
- [ ] Update the nodes as they execute and change
- [ ] Show the list of rules that have executed
- [ ] Implement a code editor that offers some existing examples (similar to [this](http://www.emanueledosualdo.com/stargazer/?gist=6e54093b297c0f9df01d0c82f65b89f6))
- [ ] Keep order of tree (e.g. if child -> parent -> grand-parent)
- [ ] Handle a forest of trees (instead of just a single tree)
- [ ] Offer different options, e.g. hide bodies and name (just show dot) -- to be able to scale for a large amount of processes -- (what if there a 1000 processes?)
- [ ] (?) process highlighting

Make sure to handle new/dead processes appropriately, rather than redrawing the whole canvas at each request (checkout enter/exit features).

Possible D3 features to consider: Force simulation/draggable processes
