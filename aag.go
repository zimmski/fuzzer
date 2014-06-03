package main

import (
	"fmt"
	randMath "math/rand"
	"strconv"
	"time"

	"github.com/zimmski/tavor/token/aggregates"
	"github.com/zimmski/tavor/token/expressions"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
	"github.com/zimmski/tavor/token/sequences"
)

/*

	This is a fuzzer made using Tavor[https://github.com/zimmski/tavor].
	It fuzzes the AAG ASCII format [http://fmv.jku.at/aiger/FORMAT].

*/

func main() {
	// constants
	var maxInputCount int64 = 5
	var maxLatchCount int64 = 5
	var maxOutputCount int64 = 5
	var maxAndCount int64 = 5

	// special tokens
	ws := primitives.NewConstantString(" ")
	nl := primitives.NewConstantString("\n")

	// construct body parts
	idSequence := sequences.NewSequence(2, 2)

	inputID := lists.NewOne(
		primitives.NewConstantInt(0),
		primitives.NewConstantInt(1),
		lists.NewOne(
			idSequence.ExistingItem(),
			expressions.NewAddArithmetic(idSequence.ExistingItem(), primitives.NewConstantInt(1)),
		),
	)

	input := lists.NewAll(
		idSequence.Item(),
		nl,
	)
	inputList := lists.NewRepeat(input, 0, maxInputCount)

	latch := lists.NewAll(
		idSequence.Item(),
		ws,
		inputID.Clone(),
		nl,
	)
	latchList := lists.NewRepeat(latch, 0, maxLatchCount)

	output := lists.NewAll(
		inputID.Clone(),
		nl,
	)
	outputList := lists.NewRepeat(output, 0, maxOutputCount)

	and := lists.NewAll(
		idSequence.Item(),
		ws,
		inputID.Clone(),
		ws,
		inputID.Clone(),
		nl,
	)
	andList := lists.NewRepeat(and, 0, maxAndCount)

	// head
	docType := primitives.NewConstantString("aag")

	maxVariableIndex := expressions.NewFuncExpression(func() string {
		return strconv.Itoa(inputList.Len() + latchList.Len() + andList.Len())
	})
	numberOfInputs := aggregates.NewLen(inputList)
	numberOfLatches := aggregates.NewLen(latchList)
	numberOfOutputs := aggregates.NewLen(outputList)
	numberOfAnds := aggregates.NewLen(andList)

	header := lists.NewAll(
		docType, ws,
		maxVariableIndex, ws,
		numberOfInputs, ws,
		numberOfLatches, ws,
		numberOfOutputs, ws,
		numberOfAnds, nl,
	)

	// body
	body := lists.NewAll(
		inputList,
		latchList,
		outputList,
		andList,
	)

	// doc
	doc := lists.NewAll(header, body)

	// fuzz the document
	r := randMath.New(randMath.NewSource(time.Now().UTC().UnixNano()))

	// TODO make this somehow easier
	idSequence.Reset()

	doc.Fuzz(r)

	// output
	fmt.Print(doc.String())
}
