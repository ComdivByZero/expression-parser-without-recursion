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

func mult(s *scanner) int {
	var v int
	if s.l == '(' {
		scan(s)
		v = expr(s)
		if s.l == ')' {
			scan(s)
		} else {
			s.l = -1
		}
	} else {
		v = number(s)
	}
	return v
}

func adder(s *scanner) int {
	v := mult(s)
	for s.l == '*' || s.l == '/' {
		if s.l == '*' {
			scan(s)
			v *= mult(s)
		} else  if s.l == '/' {
			scan(s)
			v /= mult(s)
		}
	}
	return v
}

func expr(s *scanner) int {
	v := adder(s)
	for s.l == '+' || s.l == '-' {
		if s.l == '+' { 
			scan(s)
			v += adder(s)
		} else {
			scan(s)
			v -= adder(s)
		}
	}
	return v
}

func main() {
	var s scanner
	s.r = os.Stdin
	scan(&s)
	v := expr(&s)
	if s.l > -1 {
		fmt.Printf("%v\n", v)
	} else {
		fmt.Printf("error\n")
	}
}
