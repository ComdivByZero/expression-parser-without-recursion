/* 
 * expr		= adder { ('+' | '-') adder }. 
 * adder	= mult { ('*' | '/') mult }.
 * mult		= number | '(' expr ')'.
 * number	= digit { digit }.
 * digit	= '0' .. '9'.
 */
package main

import (
	"io"
	"os"
	"fmt"
)

type (
	scanner struct {
		r io.Reader
		buf [4096] byte
		i, len int
		
		l rune
	}
	stack struct {
		s [64] struct {
			function, state int
			
			v int
		}
		i int
	}
)

const (
	fExpr = iota
	fAdder
	fMult
)

func read(s *scanner) error {
	var e error
	s.len, e = s.r.Read(s.buf[:])
	for e == nil && s.len <= 0 {
		s.len, e = s.r.Read(s.buf[:])
	}
	return e
}

func scan(s *scanner) {
	for s.i >= s.len && read(s) == nil {
		s.i = 0
		for s.buf[s.i] == ' ' {
			s.i++
		}
	}
	s.l = rune(s.buf[s.i])
	s.i++
}

func step(st *stack) {
	st.s[st.i].state++
}

func state(st *stack) int {
	return st.s[st.i].state
}

func call(st *stack, fun int) {
	step(st)
	st.i++
	st.s[st.i].function = fun
	st.s[st.i].state = 0
}

func ret(st *stack) {
	st.i--;
}

func jump(st *stack, state int) {
	st.s[st.i].state = state
}

func getret(st *stack) int {
	step(st)
	return st.s[st.i + 1].v
}

func number(s *scanner) int {
	var v int
	if s.l < '0' || s.l > '9' {
		s.l = -1
		v = 0
	} else {
		for s.l >= '0' && s.l <= '9' {
			v = v * 10 + int(s.l - '0') 
			scan(s)
		}
	}
	return v
}

func mult(st *stack, s *scanner) {
	switch state(st) {
	case 0:
		if s.l == '(' {
			scan(s)
			call(st, fExpr)
		} else {
			st.s[st.i].v = number(s)
			ret(st)
		}
	case 1:
		if s.l == ')' {
			scan(s)
			st.s[st.i].v = getret(st)
			ret(st)
		} else {
			s.l = -1
			ret(st)
		}
	}
}

func adder(st *stack, s *scanner) {
	switch state(st) {
	case 0:
		call(st, fMult)
	case 1:
		st.s[st.i].v = getret(st)
		jump(st, 5)
	case 2:
		if s.l == '/' {
			step(st)
		}
		scan(s)
		call(st, fMult)
	case 3:
		st.s[st.i].v *= getret(st)
		step(st)
	case 4:
		st.s[st.i].v /= getret(st)
	case 5:
		if s.l == '*' || s.l == '/' {
			jump(st, 2)
		} else {
			ret(st)
		}
	}
}

func expr(st *stack, s *scanner) {
	switch state(st) {
	case 0:
		call(st, fAdder)
	case 1:
		st.s[st.i].v = getret(st)
		jump(st, 5)
	case 2:
		if s.l == '-' {
			step(st)
		}
		scan(s)
		call(st, fAdder)
	case 3:
		st.s[st.i].v += getret(st)
		step(st)
	case 4:
		st.s[st.i].v -= getret(st)
	case 5:
		if s.l == '+' || s.l == '-' {
			jump(st, 2)
		} else {
			ret(st)
		}
	} 
}

func calc(s *scanner) int {
	var st stack
	st.i = 0
	st.s[0].function = fExpr
	for st.i >= 0 {
		switch st.s[st.i].function {
		case fExpr:
			expr(&st, s)
		case fAdder:
			adder(&st, s)
		case fMult:
			mult(&st, s)
		}
	}
	return st.s[0].v
}

func main() {
	var s scanner
	s.r = os.Stdin
	scan(&s)
	v := calc(&s)
	if s.l > -1 {
		fmt.Printf("%v\n", v)
	} else {
		fmt.Printf("error\n")
	}
}
