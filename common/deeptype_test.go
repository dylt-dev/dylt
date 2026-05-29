package common

import (
	"bytes"
	"embed"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"testing"
	"text/template"
	"time"
	"unicode"

	"github.com/stretchr/testify/require"
)

//go:embed content/*
var content embed.FS

func TestGenGenGen(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)

	depth := 10
	r := rand.NewSource(time.Now().UTC().UnixNano())

	// generate type declaration
	bbTypeDecl := bytes.NewBuffer([]byte{})
	GenDeclaration(ctx, depth, r, bbTypeDecl)
	typeDecl := bbTypeDecl.String()

	// load template
	buf, err := content.ReadFile("content/TestGenTest.tmpl")
	require.NoError(t, err)
	require.NotNil(t, buf)
	tmpl, err := template.New("gengengen").Parse(string(buf))
	require.NoError(t, err)

	data := map[string]any{
		"depth":           depth,
		"typeDeclaration": typeDecl,
	}

	tmpl.Execute(t.Output(), data)
}

func TestGenDecodeDeepTest(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)

	// type declaration
	type typ struct {
		Data []map[string]struct {
			Slice []map[string]struct{ Val []struct{ N int } }
		}
	}
	rt := reflect.TypeFor[typ]()
	dt := NewDeepType(rt)
	var n int = 0

	// generate scalar values
	values := []any{}
	r := rand.NewSource(time.Now().UTC().UnixNano())
	GenScalarValues(ctx, rt, r, &values)

	// generate value tree
	bbValueTree := bytes.NewBuffer([]byte{})
	dt.EmitTreeDecl(&n, values, bbValueTree)
	sValueTree := bbValueTree.String()

	// generate value ref
	dt.EmitValueRef(values, t.Output())
	// dt.EmitValueRef(values, t.Output())
	// dt.EmitValueRef(values, t.Output())
	bbValueRef := bytes.NewBuffer([]byte{})
	dt.EmitValueRef(values, bbValueRef)
	sValueRef := bbValueRef.String()
	t.Logf("sValueRef=%s", sValueRef)

	// load template
	buf, err := content.ReadFile("content/deeptest.tmpl")
	require.NoError(t, err)
	require.NotNil(t, buf)
	tmpl, err := template.New("deeptest").Parse(string(buf))
	require.NoError(t, err)

	depth := 10
	data := map[string]any{
		"depth":           depth,
		"expectedVal":     values[len(values)-1],
		"lastIndex":       depth - 1,
		"typeDeclaration": "{ Data []map[string]struct { Slice []map[string]struct{ Val []struct{ N int } } } }",
		"valueRef":        sValueRef,
		"valueTree":       sValueTree,
	}

	tmpl.Execute(t.Output(), data)
}

func TestGenDecodeDeepTestGenned(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)

	// type declaration
	type typ struct {
		Aut map[int]struct {
			Perspiciatis struct {
				Repudiandae map[int]struct {
					Omnis struct {
						Maiores map[int]map[bool]struct{ Animi int }
					}
				}
			}
		}
	}
	rt := reflect.TypeFor[typ]()
	dt := NewDeepType(rt)
	var n int = 0
	depth := 10

	// generate scalar values
	values := []any{}
	r := rand.NewSource(time.Now().UTC().UnixNano())
	GenScalarValues(ctx, rt, r, &values)

	// generate value tree
	bbValueTree := bytes.NewBuffer([]byte{})
	dt.EmitTreeDecl(&n, values, bbValueTree)
	sValueTree := bbValueTree.String()

	// generate value ref
	bbValueRef := bytes.NewBuffer([]byte{})
	dt.EmitValueRef(values, bbValueRef)
	sValueRef := bbValueRef.String()

	// load template
	buf, err := content.ReadFile("content/deeptest.tmpl")
	require.NoError(t, err)
	require.NotNil(t, buf)
	tmpl, err := template.New("deeptest").Parse(string(buf))
	require.NoError(t, err)

	data := map[string]any{
		"depth":           depth,
		"expectedVal":     values[len(values)-1],
		"lastIndex":       depth - 1,
		"typeDeclaration": "struct{Aut map[int]struct{Perspiciatis struct{Repudiandae map[int]struct{Omnis struct{Maiores map[int]map[bool]struct{Animi int}}}}}}",
		"valueRef":        sValueRef,
		"valueTree":       sValueTree,
	}

	tmpl.Execute(t.Output(), data)
}

func TestEmitTree3(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)

	type deepType map[string]struct{ Values []int }

	typ := reflect.TypeFor[deepType]()

	values := []any{}
	r := rand.NewSource(time.Now().UTC().UnixNano())
	GenScalarValues(ctx, reflect.TypeFor[deepType](), r, &values)
	t.Log(values)
	level := 0
	DeepType{typ}.EmitTreeDecl(&level, values, t.Output())
}


func TestEmitTree4(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)

	type deepType map[int][]string
	
	typ := reflect.TypeFor[deepType]()

	values := []any{}
	r := rand.NewSource(time.Now().UTC().UnixNano())
	GenScalarValues(ctx, reflect.TypeFor[deepType](), r, &values)
	t.Log(values)
	level := 0
	DeepType{typ}.EmitTreeDecl(&level, values, t.Output())
}

func TestEmitTree10(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)

	type deepType struct {
		Data []map[string]struct {
			Slice []map[string]struct{ Val []struct{ N int } }
		}
	}

	typ := reflect.TypeFor[deepType]()

	values := []any{}
	r := rand.NewSource(time.Now().UTC().UnixNano())
	GenScalarValues(ctx, reflect.TypeFor[deepType](), r, &values)
	t.Log(values)
	level := 0
	DeepType{typ}.EmitTreeDecl(&level, values, t.Output())
}

// x.Data[2]["bar"].Slice[0]["foo"].Val[3].N
func TestEmitValueRef10(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)

	type typ struct {
		Data []map[string]struct {
			Slice []map[string]struct{ Val []struct{ N int } }
		}
	}
	r := rand.NewSource(time.Now().UTC().UnixNano())
	values := []any{}
	GenScalarValues(ctx, reflect.TypeFor[typ](), r, &values)
	t.Log(values)

	fmt.Print("x")
	DeepType{reflect.TypeFor[typ]()}.EmitValueRef(values, t.Output())
	fmt.Println()

}

func TestGetDeclFromType1(t *testing.T) {
	type typ []int
	s := GetDeclFromType(reflect.TypeFor[typ]())
	t.Log(s)
}

func TestGetDeclFromType2(t *testing.T) {
	type typ [][][][][][]int
	s := GetDeclFromType(reflect.TypeFor[typ]())
	t.Log(s)
}

func TestGetDeclFromType3(t *testing.T) {
	type typ map[string]int
	s := GetDeclFromType(reflect.TypeFor[typ]())
	t.Log(s)
}

func TestGetDeclFromType4(t *testing.T) {
	type typ map[string]map[int]map[bool]string
	s := GetDeclFromType(reflect.TypeFor[typ]())
	t.Log(s)
}

func TestGetDeclFromType5(t *testing.T) {
	type typ map[string][]int
	s := GetDeclFromType(reflect.TypeFor[typ]())
	t.Log(s)
}

func TestGetDeclFromType6(t *testing.T) {
	type typ struct {
		Val      int
		Name     string
		RedFlags []bool
	}
	s := GetDeclFromType(reflect.TypeFor[typ]())
	t.Log(s)
	require.Equal(t, "struct{Val int;Name string;RedFlags []bool;}", s)
}

func TestGetDeclFromType10(t *testing.T) {
	type typ struct {
		Data []map[string]struct {
			Slice []map[string]struct{ Val []struct{ N int } }
		}
	}
	s := GetDeclFromType(reflect.TypeFor[typ]())
	t.Log(s)
	require.Equal(t, "struct{Data []map[string]struct{Slice []map[string]struct{Val []struct{N int;};};};}", s)
}

func TestGenDeclaration1(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)
	r := rand.NewSource(time.Now().UTC().UnixNano())

	for range 10 {
		GenDeclaration(ctx, 1, r, t.Output())
		t.Output().Write([]byte("\n"))
	}
}

func TestGenDeclaration2(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)
	r := rand.NewSource(time.Now().UTC().UnixNano())

	for range 10 {
		GenDeclaration(ctx, 2, r, t.Output())
		t.Output().Write([]byte("\n"))
	}
}

func TestGenDeclaration100(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)
	r := rand.NewSource(time.Now().UTC().UnixNano())

	GenDeclaration(ctx, 100, r, t.Output())
	t.Output().Write([]byte("\n"))
}

func TestGenMapScalars(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)

	types := []any{
		new(map[string]bool),
		new(map[string]int),
		new(map[string]string),
		new(map[string]struct{}),
	}

	r := rand.NewSource(time.Now().UTC().UnixNano())
	for _, p := range types {
		values := []any{}
		typ := reflect.TypeOf(p).Elem()
		GenScalarValues(ctx, typ, r, &values)
		t.Logf("%s => %v", typ, values)
	}
}

func TestGetRandomFlavor(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)

	flavorCount := map[Flavor]int{
		Map:    0,
		Slice:  0,
		Struct: 0,
	}
	for range 10000 {
		flavorCount[getRandFlavor(ctx)]++
	}

	t.Log(flavorCount)
}

func TestGetRandomScalarKind(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)

	kindCount := map[reflect.Kind]int{
		reflect.Bool:   0,
		reflect.Int:    0,
		reflect.String: 0,
	}
	for range 10000 {
		kindCount[getRandScalar(ctx)]++
	}

	t.Log(kindCount)
}

func TestGenScalars1(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)

	types := []any{
		new(map[string]bool),
		new(struct{ Field bool }),
		new(struct{ Field string }),
		new([]int),
		new(struct{ Field string }),
		new([]string),
		new([]bool),
		new(struct{ Field int }),
		new(struct{ Field bool }),
		new(map[bool]string),
	}

	r := rand.NewSource(time.Now().UTC().UnixNano())
	for _, p := range types {
		values := []any{}
		typ := reflect.TypeOf(p).Elem()
		GenScalarValues(ctx, typ, r, &values)
		t.Logf("%s => %v", typ, values)
	}
}

func TestGenScalars2(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)

	types := []any{
		new(map[int]map[string]int),
		new(struct{ Field struct{ Field int } }),
		new([]map[string]bool),
		new(map[string][]string),
		new(struct{ Field []string }),
		new(map[string][]bool),
		new(struct{ Field map[bool]string }),
		new(struct{ Field struct{ Field bool } }),
		new(struct{ Field []bool }),
		new(struct{ Field map[int]bool }),
	}

	r := rand.NewSource(time.Now().UTC().UnixNano())
	for _, p := range types {
		values := []any{}
		typ := reflect.TypeOf(p).Elem()
		GenScalarValues(ctx, typ, r, &values)
		t.Logf("%s => %v", typ, values)
	}
}

func TestGenScalars10(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)

	// x.Data[0][""].Slice[0][""].Val[0].N
	type deepType struct {
		Data []map[string]struct {
			Slice []map[string]struct{ Val []struct{ N int } }
		}
	}

	values := []any{}
	r := rand.NewSource(time.Now().UTC().UnixNano())
	typ := reflect.TypeFor[deepType]()
	GenScalarValues(ctx, typ, r, &values)
	t.Logf("%s => %v", typ, values)
}

func TestGenScalars100(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)
	r := rand.NewSource(time.Now().UTC().UnixNano())

	type deepType struct {
		Nemo struct {
			Eligendi struct {
				Illum map[bool]map[bool]map[bool]struct {
					Est struct {
						Tempore struct {
							Magnam struct {
								Fugiat struct {
									Sed struct {
										Et [][][]map[int]map[string]struct {
											Quis []map[bool]map[bool][]struct {
												Ullam []map[int]struct {
													Assumenda struct {
														Rerum struct {
															Possimus [][][]map[int]struct {
																Harum map[bool]map[int]struct {
																	Aperiam struct {
																		Incidunt struct {
																			Beatae []struct {
																				Nam struct {
																					Sunt []map[string]struct {
																						Dolore [][]map[bool][]struct {
																							Assumenda []struct {
																								Consequatur struct {
																									Iste []map[bool]struct {
																										Et [][]map[bool][][][]struct {
																											Est map[bool]struct {
																												Ut [][]map[string]struct {
																													Facere struct {
																														Velit struct {
																															Ut map[int]map[string]struct {
																																Quidem struct {
																																	Ipsa map[int]struct {
																																		Molestiae map[bool][]map[int]struct {
																																			Repellat [][][]struct {
																																				Reprehenderit struct {
																																					Eaque struct {
																																						Beatae map[string]struct {
																																							Unde struct {
																																								Perspiciatis [][]struct {
																																									Possimus map[bool]map[bool]struct {
																																										Eligendi struct {
																																											Modi struct{ Vel []struct{ Possimus []int } }
																																										}
																																									}
																																								}
																																							}
																																						}
																																					}
																																				}
																																			}
																																		}
																																	}
																																}
																															}
																														}
																													}
																												}
																											}
																										}
																									}
																								}
																							}
																						}
																					}
																				}
																			}
																		}
																	}
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	values := []any{}
	typ := reflect.TypeFor[deepType]()
	GenScalarValues(ctx, typ, r, &values)
	t.Logf("%s => %v", typ, values)
}

func TestGenMapKeyString(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)

	r := rand.NewSource(time.Now().UTC().UnixNano())
	for range 1000 {
		mapKey := genMapKeyString(ctx, r)
		require.NotNil(t, mapKey)
		require.True(t, unicode.IsLower(rune(mapKey[0])))
	}
}

func TestGenStructFieldName(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)

	r := rand.NewSource(time.Now().UTC().UnixNano())
	for range 1000 {
		fieldName := genStructFieldName(ctx, r)
		require.NotNil(t, fieldName)
		require.True(t, unicode.IsUpper(rune(fieldName[0])))
	}
}

/*
Object creation     structTree1 := NewValueTree(ctx, "N", 13)

	sliceTree1  := NewValueTree(ctx, 3, structTree1)
	valTree     := NewValueTree(ctx, "Val", sliceTree1)
	mapTree     := NewValueTree(ctx, "foo", valTree)
	sliceTree2  := NewValueTree(ctx, 0, mapTree)
	structTree2 := NewValueTree(ctx, "Slice", sliceTree2)
	mapTree2    := NewValueTree(ctx, "bar", structTree2)
	sliceTree3  := NewValueTree(ctx, 2, mapTree2)
	structTree3 := NewValueTree(ctx, "Data", sliceTree3)

Field access        x.Data[2]["bar"].Slice[0]["foo"].Val[3].N

	x.Data[0][""].Slice[0][""].Val[0].N
*/
func TestGenTest10(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)

	type typ struct {
		Data []map[string]struct {
			Slice []map[string]struct{ Val []struct{ N int } }
		}
	}
	r := rand.NewSource(time.Now().UTC().UnixNano())
	values := []any{}
	GenScalarValues(ctx, reflect.TypeFor[typ](), r, &values)
	t.Log(values)

	// emit tree
	n := 0
	DeepType{reflect.TypeFor[typ]()}.EmitTreeDecl(&n, values, t.Output())

	// emit value ref
	fmt.Print("x")
	DeepType{reflect.TypeFor[typ]()}.EmitValueRef(values, t.Output())
	fmt.Println()
}

func TestGenTest100(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)
	r := rand.NewSource(time.Now().UTC().UnixNano())

	type typ struct {
		Nemo struct {
			Eligendi struct {
				Illum map[bool]map[bool]map[bool]struct {
					Est struct {
						Tempore struct {
							Magnam struct {
								Fugiat struct {
									Sed struct {
										Et [][][]map[int]map[string]struct {
											Quis []map[bool]map[bool][]struct {
												Ullam []map[int]struct {
													Assumenda struct {
														Rerum struct {
															Possimus [][][]map[int]struct {
																Harum map[bool]map[int]struct {
																	Aperiam struct {
																		Incidunt struct {
																			Beatae []struct {
																				Nam struct {
																					Sunt []map[string]struct {
																						Dolore [][]map[bool][]struct {
																							Assumenda []struct {
																								Consequatur struct {
																									Iste []map[bool]struct {
																										Et [][]map[bool][][][]struct {
																											Est map[bool]struct {
																												Ut [][]map[string]struct {
																													Facere struct {
																														Velit struct {
																															Ut map[int]map[string]struct {
																																Quidem struct {
																																	Ipsa map[int]struct {
																																		Molestiae map[bool][]map[int]struct {
																																			Repellat [][][]struct {
																																				Reprehenderit struct {
																																					Eaque struct {
																																						Beatae map[string]struct {
																																							Unde struct {
																																								Perspiciatis [][]struct {
																																									Possimus map[bool]map[bool]struct {
																																										Eligendi struct {
																																											Modi struct{ Vel []struct{ Possimus []int } }
																																										}
																																									}
																																								}
																																							}
																																						}
																																					}
																																				}
																																			}
																																		}
																																	}
																																}
																															}
																														}
																													}
																												}
																											}
																										}
																									}
																								}
																							}
																						}
																					}
																				}
																			}
																		}
																	}
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	values := []any{}
	GenScalarValues(ctx, reflect.TypeFor[typ](), r, &values)
	t.Log(values)

	// emit tree
	n := 0
	DeepType{reflect.TypeFor[typ]()}.EmitTreeDecl(&n, values, t.Output())

	// emit value ref
	fmt.Print("x")
	DeepType{reflect.TypeFor[typ]()}.EmitValueRef(values, t.Output())
	fmt.Println()

}

func TestMapTypeEmitValueRef1a(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)

	type typ map[bool]int
	// mapType := NewDeepType(reflect.TypeFor[m]())
	r := rand.NewSource(time.Now().UTC().UnixNano())
	values := []any{}
	GenScalarValues(ctx, reflect.TypeFor[typ](), r, &values)
	t.Log(values)

	mapType := NewDeepType(reflect.TypeFor[typ]())
	mapType.EmitValueRef(values, t.Output())
}

func TestSliceTypeIsScalar1(t *testing.T) {
	type sl []int
	sliceType := NewDeepType(reflect.TypeFor[sl]())
	require.True(t, sliceType.isScalar())
}

func TestSliceTypeIsScalar2(t *testing.T) {
	type sl []struct{}
	sliceType := NewDeepType(reflect.TypeFor[sl]())
	require.False(t, sliceType.isScalar())
}

func TestSliceTypeKeyName1(t *testing.T) {
	type sl []int
	sliceType := NewDeepType(reflect.TypeFor[sl]())
	require.Equal(t, "0", sliceType.keyName())
}

func TestSliceTypeKeyName2(t *testing.T) {
	type sl []string
	sliceType := NewDeepType(reflect.TypeFor[sl]())
	require.Equal(t, "0", sliceType.keyName())
}

func TestSliceTypeKeyName3(t *testing.T) {
	type sl []bool
	sliceType := NewDeepType(reflect.TypeFor[sl]())
	require.Equal(t, "0", sliceType.keyName())
}

func TestSliceTypeKeyName4(t *testing.T) {
	type sl []struct{}
	sliceType := NewDeepType(reflect.TypeFor[sl]())
	require.Equal(t, "0", sliceType.keyName())
}

func TestSliceTypeZeroValue1(t *testing.T) {
	type sl []string
	sliceType := NewDeepType(reflect.TypeFor[sl]())
	require.Equal(t, "", sliceType.zeroValue())
}

func TestSliceTypeZeroValue2(t *testing.T) {
	type sl []bool
	sliceType := NewDeepType(reflect.TypeFor[sl]())
	require.Equal(t, "false", sliceType.zeroValue())
}

func TestSliceTypeZeroValue3(t *testing.T) {
	type sl []int
	sliceType := NewDeepType(reflect.TypeFor[sl]())
	require.Equal(t, "0", sliceType.zeroValue())
}

func TestStructTypeIsScalar1(t *testing.T) {
	type st struct{ Value int }
	structType := NewDeepType(reflect.TypeFor[st]())
	require.True(t, structType.isScalar())
}

func TestStructTypeIsScalar2(t *testing.T) {
	type st struct{ Value []int }
	structType := NewDeepType(reflect.TypeFor[st]())
	require.False(t, structType.isScalar())
}

func TestStructTypeKeyName1(t *testing.T) {
	type st struct{ Value string }
	structType := NewDeepType(reflect.TypeFor[st]())
	require.Equal(t, "Value", structType.keyName())
}

func TestStructTypeKeyName2(t *testing.T) {
	type st struct{ Value int }
	structType := NewDeepType(reflect.TypeFor[st]())
	require.Equal(t, "Value", structType.keyName())
}

func TestStructTypeKeyName3(t *testing.T) {
	type st struct{ Value struct{} }
	structType := NewDeepType(reflect.TypeFor[st]())
	require.Equal(t, "Value", structType.keyName())
}

func TestStructTypeZeroValue1(t *testing.T) {
	type st struct{ Value bool }
	structType := NewDeepType(reflect.TypeFor[st]())
	require.Equal(t, "false", structType.zeroValue())
}

func TestStructTypeZeroValue2(t *testing.T) {
	type st struct{ Value int }
	structType := NewDeepType(reflect.TypeFor[st]())
	require.Equal(t, "0", structType.zeroValue())
}

func TestStructTypeZeroValue3(t *testing.T) {
	type st struct{ Value string }
	structType := NewDeepType(reflect.TypeFor[st]())
	require.Equal(t, "", structType.zeroValue())
}
