# Phi

[*Phi is currently a placeholder for the language name - other potential names: Harruba*]

Type system and interpreter for intuitionistic session types written in Go, based on the semi-axiomatic sequent calculus.

## How to use

You need to [install the Go language](https://go.dev/doc/install). Then, build the project as follows:

```bash
go build .
```

This produces the executable `phi` file, which can be used to typecheck and run programs `./phi path/to/file.phi`.
You can find some examples in the [`/examples`](/examples/) directory.

```bash
./phi examples/nat_double.phi
```

The tool supports flags like `--notypecheck` to skip typechecking, `--noexecute` to skip execution and `--verbosity <level>` (where 1 is the least verbose and 3 is the most) to control verbosity.

The `--benchmark` flag is used to evaluate the performance.
It can be used with the optional flags `--maxcores <number of cores>` and `--repeat <number of times>` to fine tune the tests.

```bash
./phi --benchmark examples/nat_double.phi
./phi --benchmark --maxcores 1 --repeat 5 ./examples/nat_double.phi
```

To run all pre-configured benchmarks, use `./phi --benchmarks`.
Further details about the performance evaluation can be found [here](benchmarks/README.md).

## Sample program

The following program defines a type `nat` (representing natural number), a function `double` and initiates two processes.

```text
type nat = +{zero : 1, succ : nat}

let double(x : nat) : nat =
    case x (
          zero<x'> => self.zero<x'>
        | succ<x'> => h <- new double(x');
                      d : nat <- new d.succ<h>;
                      self.succ<d>
    )

// Initiate execution
prc[d0] : nat =
    t : 1 <- new close t;
    z  : nat <- new z.zero<t>;
    self.succ<z>
prc[b] : nat = 
    d1 <- new double(d0);
    d2 <- new double(d1);
    fwd self d2
```

## Grammar

The full grammar accepted by out compiler is the following:

```text
<prog> ::= <statement>*

<statement> ::= type <label> = <type>                           // labelled session type       
              | let <label> ( [<param>] ) : <type> = <term>     // function declaration
              | let <label> '[' <param> ']' = <term>            // function declaration with explicit provider name
              | assuming <param>                                // add name type assumptions
              | prc '[' <name> ']' : <type> = <term>            // create processes
              | exec <label> ( )                                // execute function

<param> ::= <name> : <type> [ , <param> ]                       // typed variable names

<type> ::= [<modality>] <type_i>                                // session type with optional modality

<type_i> ::= <label>                                            // session type label
           | 1                                                  // unit type
           | + { <branch_type> }                                // internal choice
           | & { <branch_type> }                                // external choice
           | <type_i> * <type_i>                                // send
           | <type_i> -* <type_i>                               // receive
           | <modality> /\ <modality> <type_i>                  // upshift
           | <modality> \/ <modality> <type_i>                  // downshift
           | ( <type_i> ) 

<branch_type> ::= <label> : <type_i> [ , <branch_type> ]        // labelled branches

<modality> ::= u | unr | unrestricted                           // unresticted mode
             | r | rep | replicable                             // replicable mode
             | a | aff | affine                                 // affine mode
             | l | lin | linear                                 // linear mode

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
        | print <label> ; <term>                                // output label
        | ( <term> ) 

<branches> ::= <label> '<' <name> '>' => <term> [ '|' <branches> ] // term branches

<names> ::= <name> [ ',' <names> ]                              // list of names

<name> ::= 'self'                                               // provider channel[s]
         | <channel_name>                                       // channel name
         | <polarity> <channel_name>                            // channel with explicit polarity

<polarity> ::= +                                                // positive polarity
             | -                                                // negative polarity

Others:
    <label> is an alpha-numeric combination (e.g. used to represent a choice option)
    // Single line comments
    /* Multi line comments */
    whitespace is ignored
```

## Development Details

This project requires [Go](https://go.dev/doc/install) version 1.20 (or later).

To get dependencies, run `go get .`, and the you can build the project using `go  build .`.

## Code Structure

```text
phi/
├─ cmd/
├─ parser/
├─ process/
```

There are three main components.

1. *Parser*
2. *Typechecker*
3. *Execution*

- The `cmd` folder contains the entry point to either execute code from a file/string (`main.go`), or initiate a web-server (`web.go`) to compile and execute a program using an external interface.  
- The `parser` folder contains the parser *package* which processes a string and outputs a list of processes, ready to be executed.
- The `process` folder contains the *process* package which executes some given processes. It has several components:
  - `process/runtime.go`: Entry point for the runtime environment. Sets up the processes, channels and monitors and initiates the execution.
  - `process/form.go`: Contains the different forms that a process can take. Referred to as the abstract syntax tree of the processes.
  - `process/transition.go`: Defines how each form should execute.
  - `process/servers.go`: Sets up the concurrent servers (e.g. a monitor) that monitor or control the execution of the processes. Used to inform the web-server about the state of each process.
