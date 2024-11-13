package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/periaate/blume/fsio"
	. "github.com/periaate/blume/gen"
	. "github.com/periaate/blume/str"
)

func main() { fmt.Println(ParseTD(strings.Join(fsio.Args(), " ")).Unix()) }

func ParseTD(exp string) (res time.Time) {
	sar := SplitWithAll(exp, false, " ")
	abs := Any(Contains("abs", "absolute"))(sar)

	switch {
	case abs:
		res = time.Now()
	}

	for _, v := range sar {
		// fuck you @daniel
		switch {
		case Contains(":")(v):
			t, err := time.Parse(time.TimeOnly, exp)
			if err != nil {
				continue
			}
			res = t
			continue
		case strings.Count(v, "-") == 2:
			t, err := time.Parse(time.DateOnly, exp)
			if err != nil {
				continue
			}
			res = t
			continue
		}

		res = time.Now()
		if t, err := Parse(v); err == nil {
			res = res.Add(t)
		}
	}

	return
}

var (
	s = time.Second
	m = time.Minute
	h = time.Hour
	d = 24 * h
	M = 30 * d
	w = 7 * d
	y = 365 * d
)

func Parse(exp string) (t time.Duration, err error) {
	try, err := time.ParseDuration(exp)
	if err != nil {
		return try, nil
	}

	var neg time.Duration = 1
	if HasPrefix("-")(exp) {
		neg = -1
		exp = Shift(1)(exp)
	}

	mul := s
	switch {
	case HasSuffix("s")(exp):
		mul = s
	case HasSuffix("m")(exp):
		mul = m
	case HasSuffix("h")(exp):
		mul = h
	case HasSuffix("d")(exp):
		mul = d
	case HasSuffix("M")(exp):
		mul = M
	case HasSuffix("y")(exp):
		mul = y
	}

	exp = Pop(1)(exp)

	n, err := strconv.ParseInt(exp, 10, 64)
	t = mul * time.Duration(n) * neg
	return
}
