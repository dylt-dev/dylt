### Assumptions

Decoders need to work with `reflect.Value`s, not simple pointers
Parameters should be Go types where possible, not reflect.Values. Reflection is an implementation choice.

### Data Structures

Name+Value
KV
KVSlice
KVChildMap
ValueTree


### Algorithms

Decode(ctx, cli, rootKey, a any) error
------------------------------------
Confirm `a` is a non-nil pointer
Get the underlying Type for a
Lookup KVs
Convert to KvSlice
MainDecoder.Decode(ctx, kvSlice, a)

MainDecoder.decode(ctx, kvSlice, a)
-----------------------------------
Walk the pointer
CreateOrGet(ptr, kvSlice)
Get the underlying Type for a
Get the appropriate decoder for the kind
decoder.Decode(ctx, kvSlice, a)

MapDecoder.decode()
-------------------
dim := len(valueNode.Children)
ptr := CreateOrGetMapPointer(rv, dim)
for name, childNode := range valueNode.Children) {
    i := strconv.Atoi(name)
    rvKeyPtr := GetMapKeyPointer()
    rvValPtr := GetMapValuePointer()
    MainDecoder(childNode, rvValPtr)
    ptr.SetMapIndex(rvKeyPtr, rvRvValPtr)
}

ScalarDecoder.decode()
----------------------
ptr := CreateOrGetPointer(rv)
json.Unmarshal(valueNode.value, ptr)

SliceDecoder.decode()
---------------------
dim := MaxIndex(valueNode.Children) // @note this might be a good operation on the ValueNodeChildMap
ptr := CreateOrGetSlicePointer(rv, dim)
for name, childNode := range valueNode.Children) {
    i := strconv.Atoi(name)
    rvPtr := GetSliceElPointer(ptr, i) // @note this might be a good operation on an reflect.Value-based type
    MainDecoder(childNode)
}

StructDecoder.decode(ctx, kvSlice, a any)
-----------------------------------------
// a is either a pointer to an allocated struct, or a pointer to a nil pointer
// @note - it might be nice if the pointer games were resolved before we got here
//         like, in MainDecoder


pre := CreateOrGetStructPointer(rv)
for name, childNode := range valueNode.Children) {
for structField := range structType.Fields()
    i := strconv.Atoi(name)
    rvPtr := GetStructFieldPointer(ptr, structField) // @note this might be a good operation on an reflect.Value-based type
    childNode := valueNode.Children[structField.Name]
    MainDecoder(childNode, rvPtr)
}

 - Walk the incoming non-nil pointer, and either get the address of the underlying object or allocate a new object that will accommodate the number of children, whatever it takes
  - (ScalarDecoder only) Unmarshal into the pointer the value
  - (Every other decoder) Iterate over the Children collection, looking up the Decoder for each child and invoking it recursively

### Operations

CreateOrGetScalarPointer()
--------------------------
p := WalkPointer()
if p is a pointer-to-a-pointer
    get the underlying type
    allocate the underlying type
    set *p to the new address
else
   p is the pointer we need so we're all set
@note this might be unnecessary because json.Unmarshal() might handle it.
 but maybe I can just handle it myself too


CreateOrGetMapPointer()
-----------------------
p := WalkPointer()
if p.IsNil()
    allocate a new map
    *p = new map
return p


CreateOrGetSlicePointer(n)
-------------------------
p := WalkPointer()
if p.IsNil() or p is a slice that is too small
    allocate a new slice
    *p = new slice
return p


CreateOrGetStructPointer()
--------------------------
p := WalkPointer()
if p is a pointer-to-a-pointer
    get the underlying type
    allocate the underlying type
    set *p to the new address
else
   p is the pointer we need so we're all set



MapDecoder.decode()
-------------------
dim := len(valueNode.Children)
ptr := CreateOrGetMapPointer(rv, dim)
for name, childNode := range valueNode.Children) {
    i := strconv.Atoi(name)
    rvKeyPtr := GetMapKeyPointer()
    rvValPtr := GetMapValuePointer()
    MainDecoder(childNode, rvValPtr)
    ptr.SetMapIndex(rvKeyPtr, rvRvValPtr)
}

StructDecoder.decode()
----------------------
pre := CreateOrGetStructPointer(rv)
for name, childNode := range valueNode.Children) {
for structField := range structType.Fields()
    i := strconv.Atoi(name)
    rvPtr := GetStructFieldPointer(ptr, structField) // @note this might be a good operation on an reflect.Value-based type
    childNode := valueNode.Children[structField.Name]
    MainDecoder(childNode, rvPtr)
}

 - Walk the incoming non-nil pointer, and either get the address of the underlying object or allocate a new object that will accommodate the number of children, whatever it takes
  - (ScalarDecoder only) Unmarshal into the pointer the value
  - (Every other decoder) Iterate over the Children collection, looking up the Decoder for each child and invoking it recursively


### CreateValueTree(key, []*mvccpb)
kv = popKv(kv, kvs)

// Terminal case 1: the specified key does not exist in the collection at all
if kv == nil { return nil }

tree := new(ValueTree)
this.Value  kv.Value
children := kvs.ChildrenOf(key)

// Terminal case 2: the specified key exists but has no children
if len(children) == 0 {
    this.Children = nil
    return tree
}

// Recursive case: for each child, create a subtree of the child's children
// @note it's a little unsatisfying to gather children and child descendants
//       in a single step. it feels like a clean recursive algorithm would
//       gather children once, and leave the gathering of descendants to 
//       subsequent recursive calls. And maybe this is exactly what I can do,
//       and gathering descendants is merely an optimization, that isn't really
//       much of an optimization, since subsequent recursive calls are just going
//       to gather childen anyway. Or maybe it's such an obvious optimization that
//       it makes sense in the algorithm and maybe just needs more clarity
for each child in children {
    name := child.Name()
    descendants := children.ChildrenOf(name)
    this.ChildName[name] = CreateValueTree(name, descendants)
}



### Operations
type RvSlice reflect.Value
func (this RvSlice) GetElPointer(i) any {
    // If we're going to unmarshal into a slice element we need its pointer
    rv = reflect.Value(this)
    return rv.Index(i).Interface()
}

GetMapKeyPointer()
------------------

GetMapValuePointer()
--------------------

type RvStruct reflect.Value
func (this RvStruct) GetFieldPointer(fieldName) {
    // If we're going to unmarshal into a slice element we need its pointer
    rv = reflect.Value(this)
    return rv.FieldByName(fieldName).Interface()
}

WalkPointer (i any) (any, error)
--------------------------------
// Confirm we are starting with a non-nil pointer
rv :=Reflect(i)
if !rv.IsValid() { return Error }
if rv.IsNil() { return Error }
if rv.Type(p) != reflect.Pointer { return Error }

// We definitely have a non-nil pointer
// Now, we care about 3 conditions
// - rv is a pointer to a Reference type (slice or map) - return
// - rv is a non-nil pointer to a Value type (scalar struct) - return
// - rv is a pointer to a non-pointer - return Error
// - rv is a pointer to a nil value pointer - return
// - rv := rv.Elem()


for {
    rv = rv.Elem()
    if rv.Kind() 
    if !rv.IsValid() { return Error }
    if rv.IsNil() { return Error }
    if rv.Type(p) != reflect.Pointer { return Error }
}

IsReference()
-------------
// A 'reference' is basically a variable whose value might be nil, and might
// require allocation. Some functions are responsible for allocating a 
// variable to assign to a nil reference. In these cases, the reference cannot
// be passed as an argument, because that would pass the reference by value and
// then the new value could not be assigned argument. Instead, a pointer to
// the reference must be passed as an argument. Then a new object can be allocated,
// and the address of the new object can be assigned to the pointer.
// This is a cricial part of the Decoding process. A Decoder might be asked to Decode
// into an allocated object, or to decode into a nil pointer.

func IsReference(a any) bool {
    rv := Reflect(a)
    knd := rv.Kind()
    flag := knd == reflect.Slice || reflect.Map || reflect.Pointer
    return flag
}


func (this RvPointer) CreateOrGet(kvs) (any, error) {
    /*
        // Walk the pointer
        ptr, err := this.Walk()
        if err ! nil {
            return nil, err
        }
        isAllocated, err := IsPointerAllocated(ptr, err)
        if err ! nil {
            return nil, err
        }
        rvPtr := reflect.ValueOf(ptr)

        // Handle slices differently
        if rvPtr.Type().Elem() == reflect.Slice {
            return CreateOrGetSlice(kvs)
        }
        
        // *** extract all the slice specific stuff into its own mmethod. Then
        //     pull the non-slice stuff together here, possible generalizing it
        // ***

        // If the pointer is non-nil our work is done - except for slices
        // For slices we need to check if the existing slice can hold the
        //    required # of elements
        if isAllocated {
            rv := reflect.ValueOf(ptr)
            rvElem := rv.Elem()
            typElem := rvElem.Type()
            If typElem.Kind() != reflect.Slice,
                return ptr, nil
            n = kvs.MaxSliceIndex()
            if n >= this.Cap()
                typ := rv.Type().Elem()
                ptr := MakeSlice(typ, n)
            return ptr, nil

        var ptr any
        Get underlying type
        switch kind {
            case Scalar:
                ptr := rv.Interface()
            case Map:
                map := reflect.MakeMap(typUnderlying)
                ptr := reflect.New(typUnderlying)
                ptr.Elem().Set(map) 
            case Struct:
                ptr := reflect.New(typUnderlying)
            case Slice:
                n = kvs.MaxSliceIndex()
                ptr := MakeSlice(typUnderlying, n))
            default: return nil, fmt.Errorf("unexpected type (%s)", knd.String())
        }
        rv.Elem().Set(ptr)
        return ptr, nil
    */
}


func (this RvPointer) CreateOrGetMap() (ptr any, error) {
    return nil, nil
}


func (this RvPointer) CreateOrGetScalar() (ptr any, error) {
    return nil, nil
}


func (this RvPointer) CreateOrGetSlice(kvs) (ptr any, error) {
    ptr := this.Interface()
    isAllocated, err := IsPointerAllocated(ptr, err)
    rv := reflect.ValueOf(ptr)
    n := kvs.MaxIndex()+1
    
    var ptr any
    if !sAllocated && rv.Elem().Cap() >= n  {
        ptr := 
        typ := rv.Type().Elem()
        ptr := MakeSlice(typ, n)
        this.Elem().Set(ptr)
    }
    
    return ptr, nil
}

func (this RvPointer) CreateOrGetStruct() (ptr any, error) {
    return nil, nil
}


I need to figure out the difference between how to handle a pointer to a nil ref,
and how to handle a nil pointer to a ref.
- How does Walk() handle each?
- How does CreateOrGet() handle each?