package main

import (
	"phi/parser"
	"phi/process"
	"testing"
)

// Checks related to the typechecker

// Parses the cases and expects each to pass the typechecking phase successfully (or not)
func runThroughTypechecker(t *testing.T, cases []string, pass bool) {
	for i, c := range cases {
		processes, globalEnv, err := parser.ParseString(c)

		if err != nil {
			t.Errorf("compilation error in case #%d: %s\n", i, err.Error())
		}

		err = process.Typecheck(processes, globalEnv)

		if pass {
			if err != nil {
				t.Errorf("expected no type errors in case #%d, but found %s\n", i, err.Error())
			}
		} else {
			if err == nil {
				t.Errorf("expected type error in case #%d, but didn't find any\n", i)
			}
		}
	}
}

// Typechecker -> these programs should pass the typechecker
func TestTypecheckCorrectPrograms(t *testing.T) {

	cases := []string{
		"type A = 1",
		"let f() : 1 = close self",
		// send
		// MulR
		"let f(a : 1, b : 1) : 1 * 1 = send self<a, b>",
		`type A = +{l1 : 1}
			type B = 1 -o 1
			let f(a : A, b : B) : A * B = send self<a, b>`,
		// ImpL
		"let f2(a : 1 -o 1, b : 1) : 1 = send a<b, self>",
		`type A = +{l1 : 1}
			type B = 1 * 1
			let f(a : A -o B, b : A) : B = send a<b, self>`,
		// receive
		// ImpR
		"let f1() : 1 -o 1 = <x, y> <- recv self; close y",
		"let f2(b : 1) : 1 -o (1 * 1) = <x, y> <- recv self; send y<x, b>",
		"let f1() : (1 -o 1) -o 1 = <x, y> <- recv self; <x2, y2> <- recv x; close y",
	}

	runThroughTypechecker(t, cases, true)
}

// Typechecker -> these programs should fail
func TestTypecheckIncorrectPrograms(t *testing.T) {
	cases := []string{
		"type A = B",
		"prc[a] : A = close self",
		"let f() : 1 -o A = close self",
		// MulL (extra non used names)
		"let f(c : 1, a : 1, b : 1) : 1 * 1 = send self<a, b>",
		// MulL (missing names)
		"let f(b : 1) : 1 * 1 = send self<a, b>",
		// MulL (incorrect self type)
		"let f(a : 1, b : 1) : &{a : 1} = send self<a, b>",
		"let f(a : 1, b : 1) : 1 * &{a : 1} = send self<a, b>",
		// MulL (incorrect self type)
		"let f(a : &{a : 1}, b : &{b : 1}) : &{a : 1} * 1 = send self<a, b>",
		// MulL (wrong types)
		"let f(a : 1 -o 1, b : 1) : 1 * 1 = send self<a, b>",
		"let f(a : 1, b : &{a: 1}) : 1 * 1 = send self<a, b>",
		"let f(a : 1, b : 1) : 1 * (1 * 1) = send self<a, b>",
		// ImpL
		"let f2(a : 1 -o 1, b : 1) : +{x : 1} = send a<b, self>",
		"let f2(a : 1 -o 1, b : +{x : 1}) : 1 = send a<b, self>",
		"let f2(a : 1 * 1, b : 1) : 1 = send a<b, self>",
		// MulR/ImpL
		"let f2(a : 1 -o 1, b : 1) : 1 = send a<b, c>",
		"let f2(a : 1 -o 1, c : 1) : 1 = send a<self, c>",
		// ImpR
		"let f2() : 1 * 1 = <x, y> <- recv self; close y",
		"let f2(b : 1) : 1 -o (1 * 1) = <x, y> <- recv self; send x<y, b>",
	}

	runThroughTypechecker(t, cases, false)
}
