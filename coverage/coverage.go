package coverage

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"

	fuzz "github.com/AdaLogics/go-fuzz-headers"
)

type Walker struct {
	args []string
	fuzzerName  string
	fset  *token.FileSet
	src  []byte // file contents
}

// Main walker func to traverse a fuzz harness when obtaining
// the fuzzers args. Does not add the first add (t *testing.T)
func (walker *Walker) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return walker
	}
	switch n := node.(type) {
	case *ast.FuncDecl:
		if n.Name.Name == walker.fuzzerName {
			bw := &BodyWalker{
				args: make([]string, 0),
				fuzzerName: walker.fuzzerName,
				fset: walker.fset,
				src: walker.src,
			}
			ast.Walk(bw, n.Body)
			walker.args = bw.args
		}
	}
	return walker
}

type BodyWalker struct {
	args []string
	fuzzerName  string
	fset  *token.FileSet
	src  []byte // file contents
}

func (walker *BodyWalker) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return walker
	}
	switch n := node.(type) {
	case *ast.CallExpr:
		if aa, ok := n.Fun.(*ast.SelectorExpr); ok {
			if _, ok := aa.X.(*ast.Ident); ok {
				if aa.X.(*ast.Ident).Name == "f" && aa.Sel.Name == "Fuzz" {

					// Get the func() arg to f.Fuzz:
					funcArg := n.Args[0].(*ast.FuncLit)

					walker.addArgs(funcArg.Type.Params.List[1:])
				}
			}
		}
	}
	return walker
}

// Receives a list of *ast.Field and adds them to the walker
func (walker *BodyWalker) addArgs(n []*ast.Field) {
	for _, names := range n {
		for _, _ = range names.Names {
			if a, ok := names.Type.(*ast.ArrayType); ok {
				walker.addArg(getArrayType(a))
			} else {
				walker.addArg(names.Type.(*ast.Ident).Name)
			}
		}
	}
}

func (walker *BodyWalker) addArg(arg string) {
	walker.args = append(walker.args, arg)
}

func getArrayType(n *ast.ArrayType) string {
	typeName := n.Elt.(*ast.Ident).Name
	return fmt.Sprintf("[]%s", typeName)
}

func getFuzzArgs(fuzzerFileContents, fuzzerName string) ([]string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "fuzz_test.go", fuzzerFileContents, 0)
	if err != nil {
		panic(err)
	}
	w := &Walker{
		args: []string{},
		fuzzerName: fuzzerName,
		fset: fset,
		src: []byte(fuzzerFileContents),
	}
	ast.Walk(w, f)
	return w.args, nil
}

// This is the API that should be called externally.
// Params:
// fuzzerFileContents: the contents of the fuzzerfile. This should be
// obtained with os.ReadFile().
// testCase: The libFuzzer testcase. This should also be obtained
// with os.ReadFile().
func ConvertLibfuzzerSeedToGoSeed(fuzzerFileContents, testCase []byte, fuzzerName string) string {
	args, err := getFuzzArgs(string(fuzzerFileContents), fuzzerName)
	if err != nil {
		panic(err)
	}
	newSeed := libFuzzerSeedToGoSeed(testCase, args)
	return newSeed
}

// Takes a libFuzzer testcase and returns a Native Go testcase
func libFuzzerSeedToGoSeed(testcase []byte, args []string) string {
	var b strings.Builder
	b.WriteString("go test fuzz v1\n")

	fuzzConsumer := fuzz.NewConsumer(testcase)
	for argNumber, arg := range args {
		fmt.Println(argNumber)
		switch arg {
		case "[]uint8", "[]byte":
			randBytes, err := fuzzConsumer.GetBytes()
			if err != nil {
				panic(err)
			}
			b.WriteString(fmt.Sprintf("[]byte(\"%s\")", string(randBytes)))
			if argNumber != len(args)-1 {
				b.WriteString("\n")
			}
		case "string":
			s, err := fuzzConsumer.GetString()
			if err != nil {
				panic(err)
			}
			b.WriteString(fmt.Sprintf("string(\"%s\")", s))
			if argNumber != len(args)-1 {
				b.WriteString("\n")
			}
		case "int":
			randInt, err := fuzzConsumer.GetInt()
			if err != nil {
				panic(err)
			}
			b.WriteString(fmt.Sprintf("int(%s)", strconv.Itoa(randInt)))
			if argNumber != len(args)-1 {
				b.WriteString("\n")
			}
		case "int8":
			randInt, err := fuzzConsumer.GetInt()
			if err != nil {
				panic(err)
			}
			b.WriteString(fmt.Sprintf("int8(%s)", strconv.Itoa(randInt)))
			if argNumber != len(args)-1 {
				b.WriteString("\n")
			}
		case "int16":
			randInt, err := fuzzConsumer.GetInt()
			if err != nil {
				panic(err)
			}
			b.WriteString(fmt.Sprintf("int16(%s)", strconv.Itoa(randInt)))
			if argNumber != len(args)-1 {
				b.WriteString("\n")
			}
		case "int32":
			randInt, err := fuzzConsumer.GetInt()
			if err != nil {
				panic(err)
			}
			b.WriteString(fmt.Sprintf("int32(%s)", strconv.Itoa(randInt)))
			if argNumber != len(args)-1 {
				b.WriteString("\n")
			}
		case "int64":
			randInt, err := fuzzConsumer.GetInt()
			if err != nil {
				panic(err)
			}
			b.WriteString(fmt.Sprintf("int64(%s)", strconv.Itoa(randInt)))
			if argNumber != len(args)-1 {
				b.WriteString("\n")
			}
		case "uint":
			randInt, err := fuzzConsumer.GetInt()
			if err != nil {
				panic(err)
			}
			b.WriteString(fmt.Sprintf("uint(%s)", strconv.Itoa(randInt)))
			if argNumber != len(args)-1 {
				b.WriteString("\n")
			}
		case "uint8":
			randInt, err := fuzzConsumer.GetInt()
			if err != nil {
				panic(err)
			}
			b.WriteString(fmt.Sprintf("uint8(%s)", strconv.Itoa(randInt)))
			if argNumber != len(args)-1 {
				b.WriteString("\n")
			}
		case "uint16":
			randInt, err := fuzzConsumer.GetUint16()
			if err != nil {
				panic(err)
			}
			b.WriteString(fmt.Sprintf("uint16(%d)", randInt))
			if argNumber != len(args)-1 {
				b.WriteString("\n")
			}
		case "uint32":
			randInt, err := fuzzConsumer.GetUint32()
			if err != nil {
				panic(err)
			}
			b.WriteString(fmt.Sprintf("uint32(%d)", randInt))
			if argNumber != len(args)-1 {
				b.WriteString("\n")
			}
		case "uint64":
			randInt, err := fuzzConsumer.GetUint64()
			if err != nil {
				panic(err)
			}
			b.WriteString(fmt.Sprintf("uint64(%d)", randInt))
			if argNumber != len(args)-1 {
				b.WriteString("\n")
			}
		case "rune":
			randRune, err := fuzzConsumer.GetRune()
			if err != nil {
				panic(err)
			}
			b.WriteString(fmt.Sprintf("rune(%s)", string(randRune)))
			if argNumber != len(args)-1 {
				b.WriteString("\n")
			}
		case "float32":
			randFloat, err := fuzzConsumer.GetFloat32()
			if err != nil {
				panic(err)
			}
			b.WriteString("float32(")
			b.WriteString(fmt.Sprintf("%f", randFloat))
			b.WriteString(")")
			if argNumber != len(args)-1 {
				b.WriteString("\n")
			}
		case "float64":
			randFloat, err := fuzzConsumer.GetFloat64()
			if err != nil {
				panic(err)
			}
			b.WriteString("float64(")
			b.WriteString(fmt.Sprintf("%f", randFloat))
			b.WriteString(")")
			if argNumber != len(args)-1 {
				b.WriteString("\n")
			}
		case "bool":
			randBool, err := fuzzConsumer.GetBool()
			if err != nil {
				panic(err)
			}
			b.WriteString(fmt.Sprintf("bool(%t)", randBool))
			if argNumber != len(args)-1 {
				b.WriteString("\n")
			}
		default:
			panic("Fuzzer uses unsupported type")
		}
	}
	return b.String()
}
