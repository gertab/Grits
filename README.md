# Phi/Harruba

[*Phi is currently a placeholder for the language name*]

Type system and interpreter for intuitionistic session types written in Go, based on the semi-axiomatic sequent calculus.

## How to use

You need to [install the Go language](https://go.dev/doc/install).
Then, build the project as follows:

```bash
go build .
```

This produces the executable `phi` file, which can be used to typecheck and run programs `./phi path/to/file.phi`.
You can find some examples in the [`examples`](/examples/) directory.

```bash
./phi examples/nat_double.phi
```

### Tool Flags

The tool supports various flags for customization:

- `--notypecheck`: skip typechecking
- `--noexecute`: skip execution
- `--verbosity <level>`: control verbosity (1 is the least verbose, 3 is the most)

### Benchmarking

Use the `--benchmark` flag to evaluate performance.
Optional flags include `--maxcores <number of cores>` and `--repeat <number of times>` for fine-tuning tests.

```bash
./phi --benchmark examples/nat_double.phi
./phi --benchmark --maxcores 1 --repeat 5 examples/nat_double.phi
```

To run all pre-configured benchmarks: `./phi --benchmarks`.
For detailed information on performance evaluation, refer to [benchmarks/readme](benchmarks/README.md).

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

The full grammar accepted by our compiler is as follows:

```text
<prog> ::= <statement>*

<statement> ::= type <label> = <type>                           // labelled session type       
              | let <label> ( [<param>] ) : <type> = <term>     // function declaration
              | let <label> '[' <param> ']' = <term>            // function declaration with explicit provider name
              | assuming <param>                                // add name type assumptions
              | prc '[' <name> ']' : <type> = <term>            // create processes
              | main <label> ( )                                // execute function

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

<modality> ::= r | rep | replicable                             // replicable mode
             | m | mul | multicast                              // multicast mode
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

### Building and Testing

To get dependencies, run `go get .`, and then build the project using `go build .`. To execute tests, run `go test ./...`.

<!-- might be useful to include a makefile -->

### Code Structure

There are three main components.

- **Parser** - Parses a program into a list of types, functions and processes (refer to [`parser/parser.go`](parser/parser.go)).
- **Typechecker** - Typechecks a programs using intuitionistic session types (refer to [`process/typechecker.go`](process/typechecker.go) and [`types/types.go`](types/types.go)).
- **Interpreter** - Programs are executed using either the non-polarized transition semantics (v1, [`process/transition_np.go`](process/transition_np.go)) or the (async) polarized version (v2, [`process/transition.go`](process/transition.go)).

Some other notable parts:

- The entry point can be found in [`main.go`](/main.go). Cli commands are parsed in [`cmd/cli.go`](cmd/cli.go).
- [`process/runtime.go`](/process/runtime.go): Entry point for the interpreter. Sets up the processes, channels and monitor before initiating execution.
- [`process/form.go`](/process/form.go): contains the different forms that a process can take. They are used to create the AST of processes.
- [`webserver/web_server.go`](/webserver/web_server.go): provides an external interface to compile and execute a program via a webserver (refer to the [docs](/webserver/web_server.md)) (wip)
