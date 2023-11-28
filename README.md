# Phi

[*Phi is currently a placeholder for the name*]

Type system and interpreter for intuitionistic session types written in Go, based on the semi-axiomatic sequent calculus.

## type system todo

- [ ] split rule
- [ ] add modalities (then check up/down shifting)

## How to build project

You need to [install the Go language](https://go.dev/doc/install).

Then you can build the project using the following command:

To get dependencies:

`go get .`

To build project:

`go  build .`

To run the project directly, use the following command:

`go  run .`

## How to use

After building the project (using `go  build .`), you can use the CLI version...

<!-- show how the cli version works -->

## Sample programs

## Grammar

```text
<prog> ::= <statement>*

<statement> ::= type <label> = <type>                           // labelled session type       
              | let <label> ( [<param>] ) : <type> = <term>     // function declaration
              | let <label> '[' <param> ']' = <term>            // function declaration with explicit provider name
              | assuming <param>                                // add name type assumptions
              | prc '[' <name> ']' : <type> = <term>            // create processes
              | exec <label> ( )                                // execute function

<param> ::= <name> : <type> [ , <param> ]                       // typed variable names

<type> ::= '[' <modality> '] <type_i>                           // session type with modality
         | <type_i>                                             // session type

<type_i> ::= <label>                                            // session type label
           | 1                                                  // unit type
           | + { <branch_type> }                                // internal choice
           | & { <branch_type> }                                // external choice
           | <type_i> * <type_i>                                // send
           | <type_i> -* <type_i>                               // receive
           | <madality> /\ <madality> <type_i>                  // upshift
           | <madality> \/ <madality> <type_i>                  // downshift
           | ( <type_i> ) 

<branch_type> ::= <label> : <type_i> [ , <branch_type> ]        // labelled branches

<modality> ::= u | unrestricted                                 // unresticted mode
             | r | replicable                                   // replicable mode
             | a | affine                                       // affine mode
             | l | linear                                       // linear mode

<term> ::= send <name> '<' <name> , <name> '>'                  // send names
        | '<' <name> , <name> '>' <- recv <name> ; <term>       // receive names
        | <name> . <label> '<' <name> '>'                       // send label
        | case <name> ( <branches> )                            // receive label
        | <name> [ : <type> ] <- new <term>; <term>             // spawn new process
        | <label> ( [<names>] )                                 // function call
        | fwd <name> <name>                                     // forward name
        | '<' <name> , <name> '>' <- split <name> ; <term>      // split name
        | close <name>                                          // close name
        | wait <name> ; term                                    // wait for name to close
        | cast <name> '<' <name> '>'                            // send shift
        | <name> <- shift <name> ; <term>                       // receive shift
        | ( <term> ) 

<branches> ::= <label> '<' <name> '> => term [ '|' <branches> ] // term branches

<names> ::= <name> [ ',' <names> ]                              // list of names

<name> ::= 'self'                                               // provider channel[s]
         | <label>                                              // channel name
         | <polarity> label                                     // channel with explicit polarity

<polarity> ::= +                                                // positive polarity
             | -                                                // negative polarity

Others:
    <name> refers to a channel
    <label> ia an alpha-numeric combination, representing a choice option
    // Single line comments
    /* Multi line comments */
    whitespace is ignored
```

## Some examples

There are some examples in the `cmd/examples` folder.

## Code Structure

## General Folder Structure

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

### Other information

To slow down execution speed:
Set the delay property of the `RuntimeEnvironment` to a longer duration, e.g. `1000 * time.Millisecond`.

Features not currently implemented:

- [ ] Controller server to choose which transition do do
