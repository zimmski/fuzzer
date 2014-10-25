package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/zimmski/tavor/fuzz/strategy"
	"github.com/zimmski/tavor/parser"
)

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(fmt.Sprintf("cannot open tavor file %s: %v", os.Args[1], err))
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	doc, err := parser.ParseTavor(file)
	if err != nil {
		panic(fmt.Sprintf("cannot parse tavor file: %v", err))
	}

	r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

	strat := strategy.NewAlmostAllPermutationsStrategy(doc)

	ch, err := strat.Fuzz(r)
	if err != nil {
		panic(err)
	}

	var errUnmarshal = 0
	var errMarshal = 0
	var unequal = 0
	var unequalAfterOKMinimizations = 0

	for i := range ch {
		out := doc.String()

		var document interface{}
		err = json.Unmarshal([]byte(out), &document)
		if err != nil {
			fmt.Sprintf("Error Unmarshal %q: %s\n", err)

			errUnmarshal++
		}

		res, err := json.Marshal(document)
		if err != nil {
			fmt.Sprintf("Error Marshal %q: %s\n", err)

			errMarshal++
		}

		if !reflect.DeepEqual(out, string(res)) {
			unequal++
		}

		if !reflect.DeepEqual(out, string(res)) {
			unequal++
		}

		orig := out

		for _, v := range [][2]string{
			[2]string{"\\/", "/"},
			[2]string{"\\b", "\\u0008"},
			[2]string{"\\f", "\\u000c"},
			[2]string{"\\t", "\\u0009"},
			[2]string{"\\uFFFF", "\uffff"},
			[2]string{"100E+99", "1e+101"},
			[2]string{"140E+99", "1.4e+101"},
			[2]string{"-999.99E+99", "-9.9999e+101"},
			[2]string{"0.40", "0.4"},
			[2]string{"0.00", "0"},
			[2]string{"0.0,", "0,"},
			[2]string{"0.0}", "0}"},
			[2]string{"0.0]", "0]"},
			[2]string{"0e00", "0"},
			[2]string{"0e0", "0"},
			[2]string{"0e4", "0"},
			[2]string{"0E9", "0"},
			[2]string{"0e-9", "0"},
			[2]string{"0e+0", "0"},
			[2]string{"0E+0", "0"},
			[2]string{"0E+4", "0"},
			[2]string{"0E-9", "0"},
			[2]string{"0e+9", "0"},
		} {
			out = strings.Replace(out, v[0], v[1], -1)
		}

		if !reflect.DeepEqual(out, string(res)) {
			fmt.Printf("%q vs %q (original was %q)\n", out, string(res), orig)

			unequalAfterOKMinimizations++
		}

		ch <- i
	}

	fmt.Println("Report:")
	fmt.Printf("\t%d unmarshal errors\n", errUnmarshal)
	fmt.Printf("\t%d marshal errors\n", errMarshal)
	fmt.Printf("\t%d differ to the original value\n", unequal)
	fmt.Printf("\t%d differ to the original value after OK minimizations have been ignored\n", unequalAfterOKMinimizations)
}
