# Grits

Grits is a type-checker and interpreter for intuitionistic session types written in Go, based on the semi-axiomatic sequent calculus.

## How to use

You need to [install the Go language](https://go.dev/doc/install).
Then, build the project as follows:

```bash
go build .
```

This produces an executable `grits` file, which can be used to typecheck and run programs `./grits path/to/file.grits`.
You can find some examples in the [`examples`](/examples/) directory.

```bash
./grits examples/nat_double.grits
```

<details>
  <summary>Build in Windows</summary>
  
  The `go build .` command works similarly on Windows. In this case, a `grits.exe` executable file is produced, which can be used as `grits.exe <flags> <file>`. Example usage:

   ```bash
   grits.exe examples/hello.grits
   ```

</details>

<details>
  <summary>Run using docker</summary>
  
  An alternative way to build and run Grits is via Docker.

  1. You should have a Docker runtime installed. Installation instructions are available from [https://www.docker.com](https://www.docker.com). Ensure that the Docker daemon is running.
  2. Create a docker image tagged `grits` using `docker-build.sh`:
  
     ```bash
     chmod +x docker-build.sh
     ./docker-build.sh
     ```
  
     Or use `docker build -t grits:latest .`   directly.
  3. To run the docker image tagged `grits`, use the command `./docker-run.sh <flags> <file>` (might need to run `chmod +x   docker-run.sh`).
  4. For instance, to typecheck and execute `hello.grits`, use:
  
     ```bash
     ./docker-run.sh examples/hello.grits
     ```

</details>

### Tool Flags

The tool supports various flags for customization:

- `--notypecheck`: skip typechecking
- `--noexecute`: skip execution
- `--verbosity <level>`: control verbosity (1 is the least verbose, 3 is the most)

### Benchmarking

To benchmark a specific program, use the `--benchmark` flag.  Optional flags include `--maxcores <number of cores>` and `--repeat <number of times>` for fine-tuning tests. All results are stored in the `benchmark-results/` directory created during benchmarking. Example usage:

```bash
./grits --benchmark examples/nat_double.grits
./grits --benchmark --maxcores 1 --repeat 5 examples/nat_double.grits
```

To run all pre-configured benchmarks: `./grits --benchmarks`.
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

To build the project use `go build .`. To execute tests, run `go test ./...`.

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


<!-- markdownlint-configure-file {
  "no-inline-html": {
    "allowed_elements": [
      "details",
      "summary"
    ]
  }
} -->