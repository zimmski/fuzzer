package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jessevdk/go-flags"

	"github.com/zimmski/tavor"
	"github.com/zimmski/tavor/fuzz/strategy"
	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/aggregates"
	"github.com/zimmski/tavor/token/constraints"
	"github.com/zimmski/tavor/token/expressions"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
	"github.com/zimmski/tavor/token/sequences"
	"github.com/zimmski/tavor/token/variables"
)

/*

	This is a fuzzer made using Tavor[https://github.com/zimmski/tavor].
	It fuzzes the AAG ASCII format [http://fmv.jku.at/aiger/FORMAT].

	See aag.tavor for the corresponding Tavor format file.

*/

func aagToken() token.Token {
	// constants
	maxRepeat := int64(tavor.MaxRepeat)

	// special tokens
	ws := primitives.NewConstantString(" ")
	nl := primitives.NewConstantString("\n")

	// construct body parts
	literalSequence := sequences.NewSequence(2, 2)

	existingLiteral := lists.NewOne(
		primitives.NewConstantInt(0),
		primitives.NewConstantInt(1),
		lists.NewOne(
			literalSequence.ExistingItem(nil),
			expressions.NewAddArithmetic(literalSequence.ExistingItem(nil), primitives.NewConstantInt(1)),
		),
	)

	input := lists.NewAll(
		literalSequence.Item(),
		nl,
	)
	inputList := lists.NewRepeat(input, 0, maxRepeat)

	latch := lists.NewAll(
		literalSequence.Item(),
		ws,
		existingLiteral.Clone(),
		nl,
	)
	latchList := lists.NewRepeat(latch, 0, maxRepeat)

	output := lists.NewAll(
		existingLiteral.Clone(),
		nl,
	)
	outputList := lists.NewRepeat(output, 0, maxRepeat)

	andListVar := variables.NewVariableReference(variables.NewVariable("andList", nil))
	andListVarEntry := variables.NewVariable("e", nil)
	andLiteral := variables.NewVariable("andLiteral", literalSequence.Item())

	andCycle, err := expressions.NewPath(
		andListVar,
		variables.NewVariableValue(andLiteral),
		variables.NewVariableItem(primitives.NewConstantInt(0), andListVarEntry),
		[]token.Token{
			expressions.NewMulArithmetic(expressions.NewDivArithmetic(variables.NewVariableItem(primitives.NewConstantInt(2), andListVarEntry), primitives.NewConstantInt(2)), primitives.NewConstantInt(2)),
			expressions.NewMulArithmetic(expressions.NewDivArithmetic(variables.NewVariableItem(primitives.NewConstantInt(4), andListVarEntry), primitives.NewConstantInt(2)), primitives.NewConstantInt(2)),
		},
		[]token.Token{
			primitives.NewConstantInt(0),
			primitives.NewConstantInt(1),
		},
	)
	if err != nil {
		panic(err)
	}

	existingLiteralAnd := lists.NewOne(
		primitives.NewConstantInt(0),
		primitives.NewConstantInt(1),
		lists.NewOne(
			literalSequence.ExistingItem([]token.Token{andCycle.Clone()}),
			expressions.NewAddArithmetic(literalSequence.ExistingItem([]token.Token{andCycle.Clone()}), primitives.NewConstantInt(1)),
		),
	)

	and := lists.NewAll(
		andLiteral,
		ws,
		existingLiteralAnd.Clone(),
		ws,
		existingLiteralAnd.Clone(),
		nl,
	)
	andList := lists.NewRepeat(and, 0, maxRepeat)

	// head
	docType := primitives.NewConstantString("aag")

	numberOfInputs := aggregates.NewLen(inputList)
	numberOfLatches := aggregates.NewLen(latchList)
	numberOfOutputs := aggregates.NewLen(outputList)
	numberOfAnds := aggregates.NewLen(andList)
	maxVariableIndex := lists.NewOne(
		expressions.NewAddArithmetic(numberOfInputs.Clone(), expressions.NewAddArithmetic(numberOfLatches.Clone(), numberOfAnds.Clone())),
		expressions.NewAddArithmetic(numberOfInputs.Clone(), expressions.NewAddArithmetic(numberOfLatches.Clone(), expressions.NewAddArithmetic(numberOfAnds.Clone(), primitives.NewConstantInt(1)))), // M does not have to be exactly I + L + A there can be unused Literals
	)

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
		variables.NewVariable("andList", primitives.NewScope(andList)),
	)

	// symbols
	vi := variables.NewVariableSave("e", lists.NewUniqueItem(inputList))
	symbolInput := lists.NewAll(
		primitives.NewConstantString("i"),
		vi,
		lists.NewIndexItem(variables.NewVariableValue(vi)),
		primitives.NewConstantString(" "),
		lists.NewRepeat(
			primitives.NewCharacterClass("\\w "),
			1,
			maxRepeat,
		),
		primitives.NewConstantString("\n"),
	)

	vl := variables.NewVariableSave("e", lists.NewUniqueItem(latchList))
	symbolLatch := lists.NewAll(
		primitives.NewConstantString("l"),
		vl,
		lists.NewIndexItem(variables.NewVariableValue(vl)),
		primitives.NewConstantString(" "),
		lists.NewRepeat(
			primitives.NewCharacterClass("\\w "),
			1,
			maxRepeat,
		),
		primitives.NewConstantString("\n"),
	)

	vo := variables.NewVariableSave("e", lists.NewUniqueItem(outputList))
	symbolOutput := lists.NewAll(
		primitives.NewConstantString("o"),
		vo,
		lists.NewIndexItem(variables.NewVariableValue(vo)),
		primitives.NewConstantString(" "),
		lists.NewRepeat(
			primitives.NewCharacterClass("\\w "),
			1,
			maxRepeat,
		),
		primitives.NewConstantString("\n"),
	)

	symbols := lists.NewAll(
		lists.NewRepeatWithTokens(
			symbolInput,
			primitives.NewConstantInt(0),
			aggregates.NewLen(inputList),
		),
		lists.NewRepeatWithTokens(
			symbolLatch,
			primitives.NewConstantInt(0),
			aggregates.NewLen(latchList),
		),
		lists.NewRepeatWithTokens(
			symbolOutput,
			primitives.NewConstantInt(0),
			aggregates.NewLen(outputList),
		),
	)

	// comments
	comment := lists.NewAll(
		lists.NewRepeat(
			primitives.NewCharacterClass("\\w "),
			1,
			maxRepeat,
		),
		primitives.NewConstantString("\n"),
	)

	comments := lists.NewAll(
		primitives.NewConstantString("c\n"),
		lists.NewRepeat(
			comment,
			0,
			maxRepeat,
		),
	)

	// doc
	doc := lists.NewAll(
		literalSequence.ResetItem(),
		header,
		body,
		constraints.NewOptional(symbols),
		constraints.NewOptional(comments),
	)

	return doc
}

func main() {
	var opts struct {
		Seed int64 `long:"seed" description:"Seed for all the randomness"`
	}

	p := flags.NewParser(&opts, flags.None)

	_, err := p.Parse()
	if err != nil {
		panic(err)
	}

	if opts.Seed == 0 {
		opts.Seed = time.Now().UTC().UnixNano()
	}

	log.Infof("using seed %d", opts.Seed)

	doc := aagToken()

	ch, err := strategy.NewRandom(doc, rand.New(rand.NewSource(opts.Seed)))
	if err != nil {
		panic(err)
	}

	for i := range ch {
		fmt.Print(doc.String())

		ch <- i
	}
}
