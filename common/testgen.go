package common

/*
	type typ [][]struct {
		Tempora struct {
			Eum map[string]struct{ Dolorem map[int][]map[bool][]string }
		}
	}

	x0 := "meat"
	x1 := []string{x0}
	x2 := map[bool][]string{true: x1}
	x3 := []map[bool][]string{x2}
	x4 := map[int][]map[bool][]string{13: x3}
	x5 := struct{Dolorem map[int][]map[bool][]string}{Dolorem: x4}
	x6 := map[string]struct{Dolorem map[int][]map[bool][]string}{"foo": x5}
	x7 := struct{Eum map[string]struct{Dolorem map[int][]map[bool][]string}}{Eum: x6}
	x8 := struct{Tempora struct{Eum map[string]struct{Dolorem map[int][]map[bool][]string}}}{Tempora: x7}
	x9 := []struct{Tempora struct{Eum map[string]struct{Dolorem map[int][]map[bool][]string}}}{x8}
	x10 := [][]struct{Tempora struct{Eum map[string]struct{Dolorem map[int][]map[bool][]string}}}{x9}
	var x typ = x10

	expected, err := json.Marshal(x0)
	require.NoError(t, err)
	kvs := encode(ctx, x)
	require.NotNil(t, kvs)
	require.Equal(t, 1, len(kvs))
	require.Equal(t, KeyString("/0/0/Tempora/Eum/foo/Dolorem/13/0/true/0"), kvs[0].Key)
	require.Equal(t, expected, kvs[0].Value)
	fmt.Fprint(t.Output(), kvs)
*/

/*
func WriteDeclaration(ctx *EcoContext, n int, w io.Writer) {
	ctx.Signature("genDeclaration", n)
	ctx.Inc()
	defer ctx.Dec()

	flavor := getRandFlavor(ctx)
	if n < 1 {
		return
	}

	switch flavor {
	case Map:
		writeMapDeclaration(ctx, n, w)
	case Slice:
		writeSliceDeclaration(ctx, n, w)
	case Struct:
		writeStructDeclaration(ctx, n, w)
	default:
		panic(fmt.Errorf("How'd I get a flavor of %s???", flavor))
	}
}

func genMapKeyString(ctx *EcoContext) string {
	mapKey := genRandScalarValue(ctx, reflect.TypeFor[string]()).(string)
	return strings.ToLower(mapKey)
}

func genMapKeyValue(ctx *EcoContext, typ reflect.Type) any {
	ctx.Signature("genMapKeyValue", typ)
	ctx.Inc()
	defer ctx.Dec()

	var a any
	keyType := typ.Key()
	keyDeep := DeepType{keyType}
	keyFlavor := keyDeep.Flavor()
	if keyFlavor != Scalar {
		panic("inconthievalble!")
	}
	if keyType.Kind() == reflect.String {
		a = genMapKeyString(ctx)
	} else {
		a = genRandScalarValue(ctx, keyType)
	}

	return a
}

func genMapValues(ctx *EcoContext, typ reflect.Type, values *[]any, n int) {
	ctx.Signature("genMapValues", typ, len(*values))
	ctx.Inc()
	defer ctx.Dec()

	for range n {
		// emit random key value
		ctx.Commentf("generating %d value(s) ...", n)
		a := genMapKeyValue(ctx, typ)
		ctx.Infof("value=%v", a)
		*values = append(*values, a)

		// emit value if scalar, else recurse
		if isElemScalar(typ) {
			a := genRandScalarValue(ctx, typ.Elem())
			*values = append(*values, a)
		} else {
			GenScalarValues(ctx, typ.Elem(), values)
		}
	}
}

func genMapDeclaration(ctx *EcoContext, n int) string {
	ctx.Signature("genMapDeclaration", n)
	ctx.Inc()
	defer ctx.Dec()

	if n < 1 {
		return ""
	}

	sb := strings.Builder{}
	sb.WriteString("map[")
	sb.WriteString(getRandScalar(ctx).String())
	sb.WriteString("]")
	sb.WriteString(genScalarOrRecurse(ctx, n))

	return sb.String()
}

func getRandFlavor(ctx *EcoContext) Flavor {
	ctx.Signature("getRandFlavor")
	ctx.Inc()
	defer ctx.Dec()

	nMax := 3
	switch rand.IntN(nMax) {
	case 0:
		return Map
	case 1:
		return Slice
	case 2:
		return Struct
	default:
		panic("inconthievable!")
	}
}

func getRandScalar(ctx *EcoContext) reflect.Kind {
	ctx.Signature("getRandScalar")
	ctx.Inc()
	defer ctx.Dec()

	nMax := int(reflect.UnsafePointer)
	for {
		n := rand.IntN(nMax)
		knd := reflect.Kind(n)
		switch knd {
		case reflect.Bool,
			reflect.Int,
			reflect.String:
			return knd
		default:
			continue
		}
	}
}

func genRandScalarValue(ctx *EcoContext, typ reflect.Type) any {
	ctx.Signature("genRandScalarValue", typ)
	ctx.Inc()
	defer ctx.Dec()

	switch typ.Kind() {
	case reflect.Bool:
		return faker.Bool()
	case reflect.Int:
		return faker.Int1000()
	case reflect.String:
		return faker.LoremWord()
	default:
		panic("inconthievalble!")
	}
}

func genScalarOrRecurse(ctx *EcoContext, n int) string {
	ctx.Signature("genScalarOrRecurse", n)
	ctx.Inc()
	defer ctx.Dec()

	if n == 1 {
		return getRandScalar(ctx).String()
	}

	return GenDeclaration(ctx, n-1)
}

func genSliceDeclaration(ctx *EcoContext, n int) string {
	ctx.Signature("genSliceDeclaration", n)
	ctx.Inc()
	defer ctx.Dec()

	if n < 1 {
		return ""
	}

	sb := strings.Builder{}
	sb.WriteString("[]")
	return genScalarOrRecurse(ctx, n)
}

func genSliceValues(ctx *EcoContext, typ reflect.Type, values *[]any, n int) {
	ctx.Signature("genSliceValues", typ, len(*values))
	ctx.Inc()
	defer ctx.Dec()

	// emit random slice index from [0, n)
	*values = append(*values, rand.IntN(n))

	// if slice type is scalar, emit scalar, else recurse
	tyElem := typ.Elem()
	flavorElem := NewFlavor(tyElem.Kind())
	for range n {
		if flavorElem == Scalar {
			a := genRandScalarValue(ctx, tyElem)
			*values = append(*values, a)
		} else {
			GenScalarValues(ctx, tyElem, values)
		}
	}
}

func genStructFieldName(ctx *EcoContext) string {
	fieldName := genRandScalarValue(ctx, reflect.TypeFor[string]()).(string)
	bufSrc := []byte(fieldName)
	bufDst := make([]byte, len(bufSrc))
	caser := cases.Title(language.English)
	nDst, nSrc, err := caser.Transform(bufDst, bufSrc, true)
	if err != nil {
		panic(err)
	}
	if nDst < len(bufDst) {
		panic("nDst too small")
	}
	if nSrc < len(bufSrc) {
		panic("nSrc too small")
	}
	fieldName = string(bufDst)

	return fieldName
}

func genStructDeclaration(ctx *EcoContext, n int) string {
	ctx.Signature("genStructDeclaration", n)
	ctx.Inc()
	defer ctx.Dec()

	if n < 1 {
		return ""
	}

	fieldName := genStructFieldName(ctx)
	sb := strings.Builder{}
	fmt.Fprintf(&sb, "struct{%s ", fieldName)
	sb.WriteString(genScalarOrRecurse(ctx, n))
	sb.WriteString("}")

	return sb.String()
}

func genStructValues(ctx *EcoContext, typ reflect.Type, values *[]any) {
	ctx.Signature("genStructValues", typ, len(*values))
	ctx.Inc()
	defer ctx.Dec()

	if typ.NumField() == 0 {
		return
	}

	field := typ.Field(0)
	fieldType := field.Type
	fieldDeep := DeepType{fieldType}
	fieldFlavor := fieldDeep.Flavor()

	// emit first struct field name
	*values = append(*values, field.Name)

	// emit field value if scalar, else recurse
	if fieldFlavor == Scalar {
		a := genRandScalarValue(ctx, fieldType)
		*values = append(*values, a)
	} else {
		GenScalarValues(ctx, fieldType, values)
	}
}
*/