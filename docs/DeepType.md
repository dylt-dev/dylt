## DeepType

### How to create, populate, and test objects 100+ layers deep

Building a good Decoder + Encoder is all about testing. The better the testing the better the code. So you test, and test, and test.

Eventually you'll hit a wall. To keep going farther with manual test writing would be exhausting and error prone. You'll feel like you need to write tests to test your other tests. But if you quit, you won't know if you truly built something special. Did you even push your abilities to code? Or did you only find your limits for writing tests.

You have to keep going. To keep going would take extreme measures. So, take those measures.

Here are the measures that were taken to test the Eco Decoder.

- Gen types

- Gen Value tree

- Gen Value Reference

- Putting it together

By the time testing was complete, it was no longer possible to use Go's json decoder to verify results: our objects were so deep they broke the decoder.

### Generate Type Declaration

Here's a type that's '1 deep' - a basic Go type made up of scalars

    map[string]int

Here's a type that's '2 deep' - the map value is not a scalar but a 1-deep struct

    map[string]struct{Data int}

Here's a type that's 10(ish) deep

	map[string]struct{Data []map[string]struct{Slice []map[string]struct{Val []struct{N int}}}}

Creating a type like this by hand isn't too bad. It's just a matter of taking whatever the rightmost scalar type is -- eg `int` above -- and replacing it with a map, slice, or struct. Going any farther by hand would start to get painful. You could do it, but once you got past 25 or so you're never want to do it again. You'll end up with a maximum of 1 test, and likely 0.

Luckily the code to generate long + deep type declarations is pretty 
straightforward. For a given depth, it randomly selects map, slice, or struct,
and generates a type declaration where the element type* is gen'd by a 
recursive call with depth-1. When depth == 1, a scalar element type is 
generated instead, and the recursion terminates.

* element type refers to a struct's field type, a map's value type, or a
slice's element type. `reflect.Type.Elem()` returns the element type for
each.


```
func genDeclaration(ctx *EcoContext, n int, r rand.Source, w io.Writer) {
	if n < 1 {
		return
	}

	flavor := getRandFlavor(ctx)
	switch flavor {
	case Map:
		genMapDeclaration(ctx, n, r, w)
	case Slice:
		genSliceDeclaration(ctx, n, r, w)
	case Struct:
		genStructDeclaration(ctx, n, r, w)
	}
}


func genMapDeclaration(ctx *EcoContext, n int, r rand.Source, w io.Writer) {
	if n < 1 {
		return
	}

	w.Write([]byte("map["))
	w.Write([]byte(getRandScalar(ctx).String()))
	w.Write([]byte("]"))
	writeScalarOrRecurse(ctx, n, r, w)
}


func writeScalarOrRecurse(ctx *EcoContext, n int, r rand.Source, w io.Writer) {
	if n == 1 {
		w.Write([]byte(getRandScalar(ctx, ).String()))
	} else {
		genDeclaration(ctx, n-1, r, w)
	}
}
```

### Generating Scalar Values

Once we have our type, we need to populate our type. There are well-known
packages for populating deep objects. In our case however, we don't want to
populate our deep type. We want to populate a ValueTree which we can then
decode to test our Decoder. For this, we are on our own. We can however use a
Faker package to generate individual values.

We will also need to generate at least one unit test `require` statement, to
test the scalar value that ultimately gets written to our deep type's leaf node.

The ValueTree and the `require` value reference need to be coordinated. The leaf
value needs to match, and all intermediate slice keys, map keys, and field name
references need to match too. Otherwise there's no way to navigate to the 
leaf in our decoded object.

Since the ValueTree and the `require` value reference need exactly the same
values, it makes sense to generate a series of scalar values first, and then 
generate the ValueTree and the value ref from the same values.

Generating the values first is fairly simple. First, generate the value for the
key: a slice index, a map key, or a struct field name. Then, determine if the
element type is a scalar or not. If it's a scalar, we are at a leaf node, so
we generate one last value and return. If it's not a scalar we recurse.

All values are accumulated in a *[]any slice. A pointer is required so the
slice can grow as new values are appended.


```
func genScalarValues(ctx *EcoContext, typ reflect.Type, r rand.Source, values *[]any) {
	flavor := NewDeepType(typ).Flavor()
	switch flavor {
	case Map:
		genMapValues(ctx, typ, r, values, 1)
	case Slice:
		genSliceValues(ctx, typ, r, values, 1)
	case Struct:
		genStructValues(ctx, typ, r, values)
	}		
}

func genMapValues(ctx *EcoContext, typ reflect.Type, r rand.Source, values *[]any, n int) {
	ctx.Signature("genMapValues", typ, len(*values))
	ctx.Inc()
	defer ctx.Dec()

	for range n {
		// emit random key value
		ctx.Commentf("generating %d value(s) ...", n)
		a := genMapKeyValue(ctx, typ, r)
		ctx.Infof("value=%v", a)
		*values = append(*values, a)

		// emit value if scalar, else recurse
		if isElemScalar(typ) {
			a := genRandScalarValue(ctx, typ.Elem(), r)
			*values = append(*values, a)
		} else {
			genScalarValues(ctx, typ.Elem(), r, values)
		}
	}
}

```

### Generate Value Tree

Once we have a series of values for our deep type, generating the code to create
a ValueTree is fairly straightforward, though it does require a few choices.
First, we need to figure out how to write the statements. Here's code that will
create a ValueTree for the type `map[string]struct{Values []int}`

```
// ValueTree for map[string]struct{Values []int}

tree0 := NewValueTree(ctx, 0, 260)
tree1 := NewValueTree(ctx, "Values", tree0)
tree2 := NewValueTree(ctx, "odit", tree1)
```

The ValueTree is created from the bottom up. First a ValueTree is created for
an []int with a value of 260 at [0]. Then a VT is created with a field name of
`Values` with a field value of  `tree0`. Then a map entry with a key of `odit`
is created, with a value of `tree1`.

To generate this code, the `ValueTree`s will need to be generated from the
bottom up, and we will need to number the variables appropriately, so the first
variable is `tree0` and the last is `tree[n-1]`. To do the latter we need to do
a little trick with an `*int`, but it's not so bad.

```
func (dt DeepType) EmitTreeDecl(plevel *int, values []any) {
	// If scalar, emit the current tree
	// Else, recurse on the element type, skipping the first value,
	//       then bump the level and emit the current tree
	if dt.isScalar() {
		key := dt.createValueTreeKey(values[0])
		val := values[1]
		fmt.Printf("tree%d := NewValueTree(ctx, %v, %v)\n", *plevel, key, val)
	} else {
		key := dt.createValueTreeKey(values[0])
		dt.nextType().EmitTreeDecl(plevel, values[1:])
		*plevel++
		fmt.Printf("tree%d := NewValueTree(ctx, %v, tree%d)\n", *plevel, key, *plevel-1)
	}
}
```


### Generate Value Reference

A Value reference is an expression that refers to the scalar value at the leaf
node of an instance of a deep type. Its purpose is to be used in a 
`require.Equal()` call to see if the leaf node of a decoded object matches its
expected value.

A ValueReference for a 10-deep type looks like this

`x.Data[0][""].Slice[0][""].Val[0].N`

A `require.Equal()` call using this value would look like this

`require.Equal(t, 13, x.Data[0][""].Slice[0][""].Val[0].N)`

Generating value references is simple: generate an expression for an object's
top level key, then recurse and keep appending keys till the leaf node is 
reached. The leaf node value will not be included in the Value Reference, though
it will be used as the expected value in the `require.Equal()`.

```
func (dt DeepType) EmitValueRef(values []any) {
	dt.emitKeyRef(values)
	// if not scalar, recurse
	if !dt.isScalar() {
		dt.nextType().EmitValueRef(values[1:])
	}
}

func (t MapType) emitKeyRef(values []any) {
	// strings get quoted; all other types as-is
	if t.typ.Key().Kind() == reflect.String {
		fmt.Printf("[%q]", values[0])
	} else {
		fmt.Printf("[%v]", values[0])
	}
}
```

### Putting it all together to create a test

None of the above actually performs tests, or generates tests. It just emits
snippets of Go code that can be used to build a test. To generate all the 
snippets required for a complete test, you'll need to write a test to generate a
deep type of the desired depth. Then copy paste that snippet into a new test,
and add lines to ...

- Generate scalar values
- Generate a ValueTree instantiation from the scalar values
- Generate a Value Reference from the scalar values

Then the Value Tree and Value Reference pieces can be copy/pasted into a new
test, along with the deep type declaration.

Here's a 10-Deep test assembled from generated snuppets. The 100-Deep test looks
the same. Just 10x longer.

```
func TestDecodeDeepType10Genned(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	decoder := MainDecoder{}
	type deepType struct{Data []map[string]struct{Slice []map[string]struct{ Val []struct{N int}}}}

	tree0 := NewValueTree(ctx, "N", 168)
	tree1 := NewValueTree(ctx, 0, tree0)
	tree2 := NewValueTree(ctx, "Val", tree1)
	tree3 := NewValueTree(ctx, "mollitia", tree2)
	tree4 := NewValueTree(ctx, 0, tree3)
	tree5 := NewValueTree(ctx, "Slice", tree4)
	tree6 := NewValueTree(ctx, "consequatur", tree5)
	tree7 := NewValueTree(ctx, 0, tree6)
	tree8 := NewValueTree(ctx, "Data", tree7)
	tree := tree8

	var x deepType
	p := &x
	err := decoder.Decode(ctx, tree, p)
	require.NoError(t, err)
	require.Equal(t, 168, x.Data[0]["consequatur"].Slice[0]["mollitia"].Val[0].N)
	t.Log(x)
	buf, err := json.MarshalIndent(x, "", "\t")
	require.NoError(t, err)
	t.Log(string(buf))
	
}

```

### Generating Tests

In addition to generating snippets, generating full tests and test files is 
supported. Generating tests is tricky. Often it is not possible to generate
tests in a single pass, and it is necessary to generate intermediate tests
that can eventually generate the final test.

The biggest snag is regarding type declaration. It is straightforward in go
to generate a string that compiles to a type declaration. But there's no way
to actually generate a reflect.Type from a string. That requires generating
a file of Go code that includes the type declaration as a statement, along
with other statements that generate test pieces like the Value Tree and
Value Reference. Then, this newly generated Go code can be executed to generate
an actual file of tests.

Stage 1: Generate one or more deep type declarations
Stage 2: For each deep type declaration,
           Gen decl
           Gen values
		   Gen tree
		   Gen ref
		   Execute a template to geneate a test that contains type decl, the
		      value tree, and the 
Stage 3: Run the actual test(s)

It's very clumsy to run multiple test generation stages and copy/paste from a 
console window into an editor, plus hand-editing to change test generation
parameters like type depth and # of tests. There are techniques that can make
this less clumsy
	
	- Parameterize desired behavior through environment variables
	- Write output to new files through output redirection
	- Script tests to run in sequence from the command line

All tests will be generated except a single bootstrap test. This test will have
no domain-specific logic in it whatsoever. All it will do is retrieve parameters
from environment variables, and execute templates to generate code.

@note perhaps an uber-bootstrap test would be desirable, that takes a template
name as a parameter. Something to think about.