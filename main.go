package main

import (
	"log"
	"phi/cmd"
	"phi/parser"
	"phi/process"
	"time"
)

const development = false

func main() {
	if development {
		p := `
			prc[a] : 1 = wait b; print ok; close self
			prc[b] : 1 = close self
			`
		dev(p)
	} else {
		cmd.Cli()
	}
}

func dev(program string) {
	// For DEVELOPMENT only: we can run programs directly, bypassing the CLI version
	const (
		executionVersion = process.NORMAL_ASYNC
		typecheck        = true
		execute          = true
		delay            = 0 * time.Millisecond
	)
	var processes []*process.Process
	var assumedFreeNames []process.Name
	var globalEnv *process.GlobalEnvironment
	var err error

	processes, assumedFreeNames, globalEnv, err = parser.ParseString(program)

	if err != nil {
		log.Fatal(err)
		return
	}

	// globalEnv.LogLevels = []process.LogLevel{}
	globalEnv.LogLevels = []process.LogLevel{
		process.LOGINFO,
		process.LOGRULE,
		process.LOGPROCESSING,
		process.LOGRULEDETAILS,
		process.LOGMONITOR,
	}

	if typecheck {
		err = process.Typecheck(processes, assumedFreeNames, globalEnv)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	if execute {
		re, _, _ := process.NewRuntimeEnvironment()
		re.GlobalEnvironment = globalEnv
		re.Typechecked = typecheck
		re.Delay = delay

		process.InitializeProcesses(processes, nil, nil, re)
	}
}

/* ignore sample programs -- used for development*/

const program = `

// prc[a] : lin /\ rep 1 = x <- shift self; close x
// prc[b] : lin 1 = cast a<self>

prc[a] : lin 1 = x <- shift b; wait x; close self
prc[b] : rep \/ lin 1 = cast self<c>
prc[c] : rep 1 = close self


// type bin = +{b0 : bin, b1 : bin, e : 1}
// type list = +{cons : bin * list, nil : 1}

// // proc append (R : list) (L : list) (K : list) =
// //   recv L ( 'cons(x,L') => R' <- call append R' L' K ;
// //                           send R 'cons(x,R')
// //          | 'nil() => fwd R K )


// let append(L : list, K : list) : list =
//   case L ( cons<c> => 
//                 <x, L2> <- recv c;
//                 R <- new append(L2, K);
//                 zz : bin * list <- new (send self<x, R>);
//                 self.cons<zz>
//          | nil<c> => wait c; 
//                      fwd self K)

// let append2(L : list, K : list) : list =
//   case L ( cons<p> => 
//                 <x, L2> <- recv p;
//                 R2 <- new append2(L2, K);
//                 p2 : bin * list <- new send p2<x, R2>;
//                 self.cons<p2>  
//          | nil<u> => 
//                 wait u; 
//                 fwd self K)

// // Nil list
// let nilList() : list = 
//   c : 1 <- new close self;
//   self.nil<c>

// let append123() : list =
	
//   print step1;
  
//   n12 : 1 <- new close self;
//   n11 : bin <- new self.e<n12>;
//   n1 : bin <- new self.b1<n11> ;   // n1: 'b1 'e

//   print step2;

//   n23 : 1 <- new close self;
//   n22 : bin <- new self.e<n23>;
//   n21 : bin <- new self.b1<n22>;
//   n2 : bin <- new self.b0<n21>;             // n2: 'b0 'b1 'e

//   print step3;

//   lnil11 <- new nilList();
//   l11 : bin * list <- new send self<n1, lnil11>;
//   l1 : list <- new self.cons<l11>;          // l1 : cons(n1, nil)    
 
//   print step4;

//   lnil21 <- new nilList();
//   l21 : bin * list <- new send self<n2, lnil21>;
//   l2 : list <- new self.cons<l21>;          // l2 : cons(n2, nil)

//   print step5;

//   x <- new append(l1, l2);        // result : con(n1, cons(n2, nil))

//   print step6;
  
//   <x1, x2> <- split x;

//   print step7;
  
//   y <- new append(x1, x2);
  
//   print step8;
//   <y', y''> <- split y;
  
//   print step9;
//   drop y';
  
//   print step10;
//   fwd self y''







// let getBin() : bin =
//    n13 : 1 <- new close self;
//    n12 : bin <- new self.e<n13>;
//    n11 : bin <- new self.b1<n12>;
//    n1 : bin <- new self.b0<n11>;            // 'b0 'b1 'e
//    fwd self n1
 
// let getOtherBin() : bin =
//    n13 : 1 <- new close self;
//    n11 : bin <- new self.e<n13>;
//    n1 : bin <- new self.b1<n11>;            // 'b1 'e
//    fwd self n1
 
// let getList() : list =
//    n1 <- new getOtherBin();
//    n2 <- new getBin();
//    lnil <- new nilList();
//    lres21 : bin * list <- new send self<n2, lnil>;
//    lres2 : list <- new self.cons<lres21>;            // lres2 : cons(n2, nil)
//    lres1 : bin * list <- new send self<n1, lres2>;
//    lres : list <- new self.cons<lres1>;              // lres : cons(n1, cons(n2, nil))
//    fwd self lres
 
 
//  // cons(n1, nil)
//  let getListN1() : list =
//    n1 <- new getOtherBin();
//    lnil <- new nilList();
//    lres1 : bin * list <- new send self<n1, lnil>;
//    self.cons<lres1>
 
//  // cons(n2, nil)
//  let getListN2() : list =
//    n2 <- new getBin();
//    lnil <- new nilList();
//    lres21 : bin * list <- new send self<n2, lnil>;
//    self.cons<lres21> 
 
//  let appendN1N2() : list =
// 	 n1 <- new getListN1();          // cons('b1 'e, nil)
// 	 n2 <- new getListN2();          // cons('b0 'b1 'e, nil)
// 	 nappend <- new append2(n1, n2);
// 	 fwd self nappend
 
 
//  let consumeBin(b : bin) : 1 = 
// 		 case b ( b0<c> => print b0; consumeBin(c)
// 				| b1<c> => print b1; consumeBin(c)
// 				| e<c>  => print e; wait c; close self)
 
//  let consumeList(l : list) : 1 = 
// 		 case l ( cons<c> => print cons;
// 							 <b, L2> <- recv c;
// 							 bConsume <- new consumeBin(b);
// 							 wait bConsume;
// 							 consumeList(L2)
// 				| nil<c>  => print nil;
// 							 wait c;
// 							 close self)
 
//  prc[a] : list = append123()
//  prc[b] : 1 = print startinggg;
//  			 yy <- new consumeList(a); 
// 			 wait yy;
// 			 print okkkkkk;
// 			 close self
`

const program55 = `

// prc[a] : 1 -* 1 = <x, y> <- recv self; wait +u; wait +v; drop +x; close y
// prc[b] : 1 = drop -a; close self
// prc[u] : 1 = close self
// prc[v] : 1 = close self

// prc[pid1] = <a, b> <- recv pid2; drop +a; drop +b; close self
// prc[pid2] = fwd self +pid3
// prc[pid3] = fwd self +pid4
// prc[pid4] = fwd self +pid5
// prc[pid5] = send self<pid6, pid7>




// type A = 1 -* 1
// prc[pid1] : 1 = send pid2<pid6, self>
// prc[pid2] : A = fwd self pid3
// prc[pid3] : A = fwd self pid4
// prc[pid4] : A = fwd self pid5
// prc[pid5] : A = 
//             <a, b> <- recv self; 
// 			<a', a''> <- split a;
// 			<a''', a''''> <- split a';
// 			drop a'';
// 			drop a''';
// 			drop a'''';
// 			close b
// prc[pid6] : 1 = close self




// prc[pid1] : 1 = drop a; print okk2; close self
// prc[a] : 1 * 1 = send self <x, y>
// prc[x] : 1 = close self
// prc[y] : 1 = close self



// prc[pid1] : 1 = <a2, b2> <- recv a; print okk; drop a2; wait b2; print okk2; close self
// prc[a] : 1 * 1 = send self <x, y>
// prc[x] : 1 = close self
// prc[y] : 1 = close self


// type A5 = affine 1
// type B5 = affine 1

// let eg5(f : affine \/ linear (A5 -* B5)) : (affine \/ linear A5) -* (affine \/ linear B5) = 
//     <x, y> <- recv self;
//     w <- shift f;
//     v <- shift x;
//     z : affine \/ linear B5 <- new cast y<self>;
//     send w<v, z>

// // type A = linear 1

// // prc[a] : linear /\ affine A  = y <- shift self; close y
// // prc[b] : linear A = cast a<self>

// type firstNegType = &{'first : secondNegType -o 1}
// type secondNegType = &{'second : (1 -o 1)}

// // Double lolli
// main f4()

// let f1[w : (1 -* 1) -* 1, z : 1] = <x, y> <- recv w; print ok; send x<z, y>
// let f2[w : 1 -* 1] = <x, y> <- recv w; wait x; print ok2; close y
// let f3[w : 1, b : (1 -* 1) -* 1, u : 1 -* 1] = send b<u, w>

// let f4[w : 1] = 
// 		z : 1 <- new close self;
// 		b <- new f1(b, z);
// 		u <- new f2(u);
// 		f3(w, b, u)

`

const p = `
type bin = +{b0 : bin, b1 : bin, e : 1}
type list = +{cons : bin * list, nil : 1}

let append(L : list, K : list) : list =
  case L ( cons<c> => 
                <x, L2> <- recv c;
                R <- new append(L2, K);
                zz : bin * list <- new (send zz<x, R>);
                self.cons<zz>
         | nil<c> => wait c; 
                     fwd self K)

// Nil list
let nilList() : list = 
  c : 1 <- new close self;
  self.nil<c>

  let getBin() : bin =
  n13 : 1 <- new close self;
  n12 : bin <- new self.e<n13>;
  n11 : bin <- new self.b1<n12>;
  n1 : bin <- new self.b0<n11>;            // 'b0 'b1 'e
  fwd self n1

let getOtherBin() : bin =
  n13 : 1 <- new close self;
  n11 : bin <- new self.e<n13>;
  n1 : bin <- new self.b1<n11>;            // 'b1 'e
  fwd self n1

let getList() : list =
  n1 <- new getOtherBin();
  n2 <- new getBin();
  lnil <- new nilList();
  lres21 : bin * list <- new send self<n2, lnil>;
  lres2 : list <- new self.cons<lres21>;            // lres2 : cons(n2, nil)
  lres1 : bin * list <- new send self<n1, lres2>;
  lres : list <- new self.cons<lres1>;              // lres : cons(n1, cons(n2, nil))
  fwd self lres


// cons(n1, nil)
let getListN1() : list =
  n1 <- new getOtherBin();
  lnil <- new nilList();
  lres1 : bin * list <- new send self<n1, lnil>;
  self.cons<lres1>

// cons(n2, nil)
let getListN2() : list =
  n2 <- new getBin();
  lnil <- new nilList();
  lres21 : bin * list <- new send self<n2, lnil>;
  self.cons<lres21> 

let appendN1N2() : list =
    n1 <- new getListN1();          // cons('b1 'e, nil)
    n2 <- new getListN2();          // cons('b0 'b1 'e, nil)
    nappend <- new append(n1, n2);
    fwd self nappend

let consumeBin(b : bin) : 1 = 
        case b ( b0<c> => print b0; consumeBin(c)
               | b1<c> => print b1; consumeBin(c)
               | e<c>  => print e; wait c; close self)

let consumeList(l : list) : 1 = 
        case l ( cons<c> => print cons;
                            <b, L2> <- recv c;
                            bConsume <- new consumeBin(b);
                            wait bConsume;
							print end1;
                            consumeList(L2)
               | nil<c>  => print nil;
                            wait c;
							print end2;
                            close self)



let appendBySplit() : list =
	n1 <- new getListN1();          // cons('b1 'e, nil)
	<x1, x2> <- split n1;       
	nappend <- new append(x1, x2);

	<y1, y2> <- split nappend;       
	// drop y1;
	nappend2 <- new append(y1, y2);
	fwd self nappend2
						
// prc[a] : list = appendBySplit()
prc[a1, a2, a3] : list = getListN1()
prc[a] : list = 
		<s1, s2> <- split a1;
		append1 <- new append(s1, s2);
		<s3, s4> <- split a2;
		append2 <- new append(a3, s3);
		drop append2;
		// drop s4;
		append3 <- new append(append1, s4);
		// drop append3;
		fwd self append3


// consume result/list
prc[b] : 1 = yy <- new consumeList(a); 
			wait yy;
			print okkkkkk;
			close self

`

const pold = `

type A = affine 1 * 1

assuming x : A

prc[a] : 1 = y <- shift b; drop y; close self
prc[b] : affine \/ linear A = cast self<x>



// Drop
// assuming a : linear 1 * 1
// prc[b] : 1 = drop a; close self

// Split
// prc[pid0] : 1 = <u, v> <- split x; wait u; wait v; close self
// prc[x] : linear 1 = close self
		

// // Double lolli
// prc[a] : 1             = send b<u, self>
// prc[b] : (1 -* 1) -* 1 = <x, y> <- recv self; send x<z, y>
// prc[u] : 1 -* 1        = <x, y> <- recv self; wait x; close y
// prc[z] : 1             = close self


// type A = affine 1
// assuming x : A
// prc[a] : affine A = y <- shift b; drop y; close self
// prc[b] : affine \/ linear A = cast self<x>

// type A = replicable \/ affine B
// type B = replicable 1 * 1
// assuming y : B
// let f(b : A) : 1 = y <- shift b; drop y; close self
// prc[a] : 1 = x : A <- new cast self<y>; f(x)

// let f() : affine \/ linear 1 = x : affine 1 <- new (close x); cast self<x>
// let f2[w : affine \/ linear 1] = x : affine 1 <- new (close x); cast w<x>
// prc[a] : affine \/ linear 1 = x : affine 1 <- new (close x); cast self<x>

// type A = linear &{a : B, b : C}
// type B = 1 * (affine\/linear 1 -* 1)
// type C = ((replicable\/linear replicable\/replicable 1) -* 1) -* 1


// // This is not allowed:
// type A = (1 -* 1) -* 1
// assuming z : 1, u : 1 -* 1
// prc[a] : (1 -* 1) -* 1 = send b<u, self>
// prc[b] : A = <x, y> <- recv self; send x<z, y>

// prc[z] : 1 = close self
// prc[u] : 1 -* 1 = <x, y> <- recv self; wait x; close y
// prc[v] : 1 = close self


// type A = 1 -* (1 -* (1 * 1))
// prc[x1] : 1 -* (1 * 1) = send z<yy, self>
// prc[x2] : 1 * 1 = send x1<xx, self>
// prc[z] : A = <x, yy> <- recv self; 
// 			 <xx, yy> <- recv yy; 
// 			 send yy<x, xx>
// prc[xx] : 1 = close self
// prc[yy] : 1 = close self

// prc[final] : 1 = <g1, g2> <- recv x2;
// 				 drop g1;
// 				 drop g2;
// 				 close self


// assuming pid3 : 1, pid4 : 1

// prc[pid1] : 1 = <pid2_first, pid2_second> <- split pid2; /* split gets its polarity from the types */
// 				k : 1 <- new send pid2_first<pid3, self>;
// 				wait k;
// 				send pid2_second<pid4, self>
// prc[pid2] : 1 -* 1 = <a, b> <- recv self; 
// 					 drop a; 
// 					 close self

// // Positive fwd
// type A = +{label1 : B}
// type B = 1
// prc[y] : 1 = case ff (label1<cont> => print cont; wait cont; close self)
// prc[ff] : A = +fwd self z
// prc[z] : A = self.label1<x>
// prc[x] : B = close self

// // Positive fwd
// type A = +{label1 : B}
// type B = 1
// prc[y1] : 1 = case z1 (label1<cont> => print cont; wait cont; close self)
// prc[y2] : 1 = case z2 (label1<cont> => print cont; wait cont; close self)
// prc[z1, z2] : A = z : A <- new (self.label1<x>); +fwd self z
// prc[x] : B = close self

// // Positive fwd
// type A = &{label1 : B}
// type B = 1
// prc[y1] : 1 = z1.label1<self>
// prc[y2] : 1 = z2.label1<self>
// prc[z1, z2] : A = z : A <- new (case self (label1<cont> => print cont; close self)); +fwd self z
// prc[x] : B = close self





// type A = &{label : +{next : 1}}

// let f1(x : A) : +{next : 1} = x.label<self>
// let f2(y : 1) : A = case self (label<zz> => zz.next<y> )

// prc[x] : +{next : 1} = f1(z)

// prc[z] : A = f2(y)
// prc[y] : 1 = close self
// prc[final] : 1 = case x (next<z> => print z; drop z; close self)





// type A = 1 -* (1 -* (1 * 1))
// prc[x1] : 1 -* (1 * 1) = send z<yy, self>
// prc[x2] : 1 * 1 = send x1<xx, self>
// prc[z] : A = <x, y> <- recv self; 
// 			 <xx, y> <- recv y; 
// 			 send y<x, xx>
// prc[xx] : 1 = close self
// prc[yy] : 1 = close self

// prc[final] : 1 = <g1, g2> <- recv x2;
// 			     print g1;
// 			     print g2;
// 			     drop g1;
// 			     drop g2;
// 			     close self



// type A = 1

// let f() : A = x : A <- new close x; 
// 			wait x; 
// 			close self

// main f()






// prc[a] : 1 = close self
// prc[b] : 1 = -fwd self a 
// prc[c] : 1 = wait b; close self

// type A = &{label : 1}
// type B = 1 -* 1
// let f(a : A, b : B) : A * B = send self<a, b>
// prc[pid1] : 1 = x <- new f(a, b); 
// 				<u, v> <- recv x;  
// 				drop u; 
// 				drop v; 
// 				close self 			% a : A, b : B







// type A = 1 * 1

// prc[pid1] : 1 = 
// 		<a, b> <- +split pid2; 
// 		<a2, b2> <- recv a; 
// 		<a3, b3> <- recv b; 
// 		wait a2; 
// 		wait b2; 
// 		wait a3; 
// 		wait b3; 
// 		close self   % pid2 : A
// prc[pid2] : A = send self<pid3, pid4>	% pid3 : 1, pid4 : 1
// prc[pid3, pid4] : 1 = close self










// let f() : 1 = close self
// prc[pid1] : 1 = x : 1 <- new f(); wait x; close self

// let f2[w : 1] = close w
// prc[pid2] : 1 = x : 1 <- new f2(); wait x; close self

//////
// type A = 1
// type B = 1
// prc[pid1] : 1 = <a, b> <- +split pid2; <a2, b2> <- recv a; <a3, b3> <- recv b; close self	
// 												% pid2 : A * B
// prc[pid2] : A * B = send self<pid3, pid4>		% pid3 : A, pid4 : B
// prc[pid3] : A = close self
// prc[pid4] : B = close self



// type A = &{label1 : 1, label2 : 1, label3 : 1}
// let f2() : A = 
// 			case self (label1<a> => close a
// 					  |label2<a> => close a
// 					  |label3<a> => close a) 

// let f3(x : &{label1 : 1}) : 1 = x.label1<self>

// prc[b] : A = f2()
// prc[dd , aa] : 1 = send a<b, self>   % a : 1 -* 1, b : 1
// prc[c] : 1 = send a<b, self>   		 % a : 1 -* 1, b : 1  
`

const program_no_errors = `

let f1(a : 1, b : 1) : 1 * 1 = send self<a, b>

type A = +{l : 1, r : 1}
type B = 1 * A
let f2(a : A, b : 1 * A) : A * B = send self<a, b>
`

const program_with_errors = `
let f1(a : 1) : 1 * 1 = send self<a, b>
let f2(a : 1, b : 1, c : 1) : 1 * 1 = send self<a, b>
`

// let f2(x : 1, y : 1) : 1 * 1 = send x<y, self>

// prc[pid1] = <a, b> <- +split pid2; <a2, b2> <- recv a; <a3, b3> <- recv b; close self
// prc[pid2] = send self<pid3, self>

// type C = 1 * 1
// type D = 1 -* 1

// let func3(next_pid : D) : C = send self< next_pid, self>
// let func2(next_pid : &{a : 1, c : 1}) : &{a : 1, c : 1} = send self< next_pid, self>

// prc[pid1] : &{a : 1, c : 1} = <a, b> <- recv pid2; wait a; close self

// undefined label reference
// type B = &{a : unknownlabel, c : 1}
// type E = 1 * X
// type A = +{a : (1 -* &{a : FF * 1}), c : 1}
// let func2(next_pid : B) : &{a : ssss, c : 1} = send self< next_pid, self>
// prc[pid1] : &{a : ssss, c : 1} = <a, b> <- recv pid2; wait a; close self

// contractive
// type A = B

// multiple types with the same name
// type A = 1
// type B = 1
// type A = 1

// const program = `
// prc[pid1]: case a
//
//	( label1<b> => wait a; close self
//	| label2<b> => close self )
//
// prc[a]: self.label1<self>
//
//	`

const program22 = `
type Receive = 1bc
type Label = label
type Unit = 1
type Select = +{a : b}
type Select2 = +{a : b, c : d}
type Branch = &{a : b}
type Branch2 = &{a : b, c : d}
type Send = a * b
type Receive = c -* b
type Brack = (a)
type Complex = +{a : (x -* &{a : f * g}), c : d}
`
const program33 = `
type Unit = 1
type Select = +{a : b}

let func3(next_pid : B) : A = send self< next_pid, self>
let func1(next_pid : a * b) : a * b = send self< next_pid, self>

let func2(next_pid : s) = send self< next_pid, self>

prc[pid1] = <a, b> <- recv pid2; wait a; close self
prc[pid2] = +fwd self pid3
prc[pid3] = +fwd self pid4
prc[pid4] = +func1(pid5)
prc[pid5] = close self
`

const program1 = `
prc[a]: f.label2<self>
prc[f]: case self
		( label1<b> => close self
		| label2<b> => close self )
`

const program2 = `
prc[a]: <a, b> <- recv c; close self
prc[c]: send self<d, self>
`

const program3 = `
let
	false() = send self.false<self>
in
  prc[pid1]: x <- -new (send pid2 <a, x>); close self
  prc[pid2]: -fwd self pid3
  prc[pid3]: <x, y> <- recv self; close y
end		
`

const program4 = `
let
	false() = self.false<self>
	true() = self.true<self>
	neg(a) = case a ( true<b> => self.false<self>
					| false<b> => self.true<self> )
in
	prc[pid0]: +true()
    prc[pid1]: +neg(pid0)

	prc[result]: case pid1 ( true<b> => wait res_true; close self
						  | false<b> => wait res_false; close self )

   	prc[res_true]: close self
   	prc[res_false]: close self
end
  `
const program5 = `
let false(): A = self.false<self>
let true(): B = self.true<self>
let neg(a): C = case a ( true<b> => self.false<self>
					| false<b> => self.true<self> )
prc[pid0]: D = +true()
prc[pid1]: E = +neg(pid0)

prc[result]: case pid1 ( true<b> => wait res_true; close self
						| false<b> => wait res_false; close self )

prc[res_true]: close self
prc[res_false]: close self
  `

// const program = `
// prc[pid1]: <a, b> <- recv pid2; close self
// prc[pid2]: send self<pid3, self>
// 	`

// const program = `
// prc[pid0]: <x, y> <- +split pid2; <a, b> <- recv x; <c, d> <- recv y; close d
// prc[pid2]: send self <xx, self>
// prc[xx]: close self
//     `

// const program = ` 	/* FWD + RCV rule  -- ok with the original scenario */
// 	let
// 	in
// 	prc[pid1]: send pid2<pidother, self>
//  	prc[pid2]: -fwd self pid3
//  	prc[pid3]: -fwd self pid4
// 	prc[pid4]: <a, b> <- recv self; close a
// 	end`

// const program = ` 	/* FWD + RCV rule  -- ok with the original scenario */
// 	let
// 	in
// prc[pid1]: send pid2<pid5, self>
// prc[pid2]: -fwd self pid3
// prc[pid3]: -fwd self pid4
// prc[pid4]: <a, b> <- recv self; close a
// 	end`

// const program = ` 	/* FWD + SND rule -- the problematic scenario*/
// 	let
// 	in
// 	prc[pid1]: <a, b> <- recv pid2; close a
// 	prc[pid2]: +fwd self pid3
// 	// prc[pid3]: +fwd self pid4
// 	prc[pid3]: send self<pid5, self>
// 	end`

// const program = ` 	/* CLS rule*/
// 	let
// 	in
// 	prc[pid1]: wait pid2; close a
// 	prc[pid2]: close self
// 	end`

// const program = ` 	/* CLS + FWD rule - problematic*/
// 	let
// 	in
// 	prc[pid1]: wait pid2; close a
//  	prc[pid2]: +fwd self pid3
// 	prc[pid3]: close self
// 	end`

// const program = ` /* CLS rule */
// 		prc[pid1]: wait pid2; close a
// 		prc[pid2]: close self
// `

// const program = ` /* CLS rule */
// prc[pid1]: <a, b> <- recv pid2; wait a; close self
// prc[pid2]: send self<pid3, self>
// prc[pid3]: close self
// `
