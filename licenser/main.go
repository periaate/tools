package main

import (
	_ "embed"
	"fmt"

	"github.com/periaate/blume/fsio"
	"github.com/periaate/blume/gen"
	. "github.com/periaate/blume/gen/T"
	. "github.com/periaate/blume/typ"
	"github.com/periaate/blume/yap"
)

func AtLeast(n int) Condition[int] {
	return func(i int) (er Error[string]) {
		if i < n {
			er = Errs("at least "+string(n), "input is less than "+string(n), "")
		}
		return nil
	}
}

func AtMost(n int) Condition[int] {
	return func(i int) (er Error[string]) {
		if i > n {
			er = Errs("at most "+string(n), "input is more than "+string(n), "")
		}
		return nil
	}
}

func Between(min, max int) Condition[int] {
	return func(i int) (er Error[string]) {
		if i < min {
			er = Errs("between "+string(min)+" and "+string(max), "input is less than "+string(min), "")
		} else if i > max {
			er = Errs("between "+string(min)+" and "+string(max), "input is more than "+string(max), "")
		}
		return nil
	}
}

func Exactly(n int) Condition[int] {
	return func(i int) (er Error[string]) {
		if i != n {
			er = Errs("exactly "+string(n), "input is not "+string(n), "")
		}
		return nil
	}
}

func Len[A any](cond Condition[int]) Condition[[]A] {
	return func(ar []A) Error[string] {
		return cond(len(ar))
	}
}

func main() {
	if gen.Any(gen.Contains("LICENSE", "License", "license"))(gen.Must(fsio.ReadDir("./"))) {
		yap.Fatal("license file already exists")
	}

	fsio.QArgs(Len[string](Exactly(1))).Match(
		func(s []String) {
			yap.Info("licenser called with", s[0])
			if lic, ok := licenses[s[0].String()]; ok {
				err := fsio.WriteNew("LICENSE", fsio.B(lic))
				if err != nil {
					yap.Fatal("error writing license", "err", err)
				}
				yap.Info("license written", "license", s[0])
			} else {
				yap.Error("license not found", "tried to find", s[0])
				fmt.Println("")
				fmt.Println("Available licenses:")
				fmt.Println(licensesStr)
			}
		},
		func(e Error[any]) {
			yap.Fatal(e.Error(), "reason", e.Reason(), "data", e.Data())
		},
	)
}

var licensesStr = `AGPL-3.0
GPL-3.0
GPL-2.0
LGPL-3.0
LGPL-2.1
NON-AI-MPL-2.0
MPL-2.0
NON-AI-APACHE-2.0
Apache-2.0
NON-AI-UNLICENSE
UNLICENSE
NON-AI-MIT
MIT`

var licenses = map[string][]byte{
	"GPL-2.0":           gpl20,
	"GPL-3.0":           gpl30,
	"AGPL-3.0":          agpl30,
	"MIT":               mit,
	"MPL-2.0":           mpl20,
	"Apache-2.0":        apache20,
	"LGPL-2.1":          lgpl21,
	"LGPL-3.0":          lgpl30,
	"UNLICENSE":         unlicense,
	"NON-AI-MIT":        non_ai_mit,
	"NON-AI-MPL-2.0":    non_ai_mpl20,
	"NON-AI-UNLICENSE":  non_ai_unlicense,
	"NON-AI-APACHE-2.0": non_ai_apache20,
}

//go:embed GPL-2.0
var gpl20 []byte

//go:embed GPL-3.0
var gpl30 []byte

//go:embed AGPL-3.0
var agpl30 []byte

//go:embed MIT
var mit []byte

//go:embed MPL-2.0
var mpl20 []byte

//go:embed Apache-2.0
var apache20 []byte

//go:embed LGPL-2.1
var lgpl21 []byte

//go:embed LGPL-3.0
var lgpl30 []byte

//go:embed UNILICENSE
var unlicense []byte

//go:embed NON-AI-MIT
var non_ai_mit []byte

//go:embed NON-AI-MPL-2.0
var non_ai_mpl20 []byte

//go:embed NON-AI-Unlicense
var non_ai_unlicense []byte

//go:embed NON-AI-Apache-2.0
var non_ai_apache20 []byte