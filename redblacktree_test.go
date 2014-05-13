/*
Copyright 2014 Gavin Bong.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
either express or implied. See the License for the specific
language governing permissions and limitations under the
License.
*/

package redblacktree

import (
    _ "fmt"
    "reflect"
    "testing"
)

var funcs map[string]reflect.Method

func init() {
    var found bool
    var put, rotateLeft, rotateRight reflect.Method

    t := reflect.TypeOf(NewTree())
    put, found = t.MethodByName("Put")
    if !found {
        panic("No method `Put` in Tree")
    }
    rotateLeft, found = t.MethodByName("RotateLeft")
    if !found {
        panic("No method `RotateLeft` in Tree")
    }
    rotateRight, found = t.MethodByName("RotateRight")
    if !found {
        panic("No method `RotateRight` in Tree")
    }

    funcs = map[string]reflect.Method{
        "rotateRight": rotateRight,
        "rotateLeft": rotateLeft,
        "put": put,
    }

    TraceOff()
    //TraceOn()
}

func True(b bool, t *testing.T) {
	if !b {
		t.Errorf("Expected [ %t ] got [ %t ]", true, b)
	}
}

func False(b bool, t *testing.T) {
	if b {
		t.Errorf("Expected [ %t ] got [ %t ]", false, b)
	}
}

func assertDirection(expected Direction, actual Direction, t *testing.T) {
    if actual != expected {
		t.Errorf("Expected (%s) got (%s)", expected, actual)
	}
}

func assertNodeKey(n *Node, expected int, t *testing.T) {
    if n.value != expected {
        t.Errorf("Expected (%#v) got (%#v)", expected, n.value)
    }
}

// using the inorder walk of the tree for equality
func assertEqualTree(tr *Tree, t *testing.T, expected string) {
    visitor := &InorderVisitor{}
    tr.Walk(visitor)
    if visitor.String() != expected {
        t.Errorf("Expected [ %s ] got [ %s ]", expected, visitor)
    }
}

// Value.IsNil returns true if v is a nil value. It panics if 
// v's Kind is not Chan, Func, Interface, Map, Ptr, or Slice.
func nillable(k reflect.Kind) bool {
    switch k {
    case reflect.Chan: fallthrough
    case reflect.Func: fallthrough
    case reflect.Interface: fallthrough
    case reflect.Map: fallthrough
    case reflect.Ptr: fallthrough
    case reflect.Slice: return true
    default: return false
    }
}

// asserts that @param `a` is nil
func Nil(a interface{}, t *testing.T) {
    if a == nil {
        return
    }
    value := reflect.ValueOf(a)
    if nillable(value.Kind()) {
        if !value.IsNil() {
			t.Errorf("%#v is not nil", a)
		}
    } else {
        t.Errorf("%#v is not nil", a)
	}
}

// asserts that @param `a` is NOT nil
func NotNil(a interface{}, t *testing.T) {
    if a == nil {
        t.Errorf("Expected NOT nil")
        return
    }
    value := reflect.ValueOf(a)
    if nillable(value.Kind()) {
        if value.IsNil() {
			t.Errorf("%#v is nil but we expected it to be NOT nil", a)
		}
    }
}

type KV struct {
    key int
    arg string
}

func ToArgs(a ...interface{}) []reflect.Value {
    in := make([]reflect.Value, len(a))
    it := make([]struct{}, len(a))
    for i, _ := range it {
        in[i] = reflect.ValueOf(a[i])
    }
    return in
}

var fixtureSmall = []struct {
    ops string
    kv  KV
    expected string
}{
    {"put", KV{7, "payload7"}, "(.7.)"},
    {"put", KV{3, "payload3"}, "((.3.)7.)"},
    {"put", KV{8, "payload8"}, "((.3.)7(.8.))"},
}

func TestRedBlackSmall(t *testing.T) {
    t1 := NewTree()
    if t1.root != nil {
        t.Errorf("root starts out as nil but got (%t)", t1.root != nil)
    }
    assertEqualTree(t1, t, ".")

    t2 := NewTree()
    for _, tt := range fixtureSmall {
        method := funcs[tt.ops]
        in := ToArgs(t2, tt.kv.key, tt.kv.arg)
        method.Func.Call(in) // @TODO returns []reflect.Value ?
        assertEqualTree(t2, t, tt.expected)
    }
}

var fixtureSimpleRightRotation = []struct {
    ops string
    kv  KV
    expected string
}{
    {"put", KV{7, "payload7"}, "(.7.)"},
    {"put", KV{3, "payload3"}, "((.3.)7.)"},
    {"put", KV{1, "payload1"}, "((.1.)3(.7.))"},
}

func TestRedBlackSimpleRightRotation(t *testing.T) {
    tr := NewTree()
    for _, tt := range fixtureSimpleRightRotation {
        method := funcs[tt.ops]
        method.Func.Call(ToArgs(tr, tt.kv.key, tt.kv.arg))
        assertEqualTree(tr, t, tt.expected)
    }
}

var fixtureCase1 = []struct {
    ops string
    kv  KV
    expected string
}{
    {"put", KV{7, "payload7"}, "(.7.)"},
    {"put", KV{8, "payload8"}, "(.7(.8.))"},
    {"put", KV{9, "payload9"}, "((.7.)8(.9.))"},
    {"put", KV{11, "payload11"}, "((.7.)8(.9(.11.)))"},
    {"put", KV{10, "payload10"}, "((.7.)8((.9.)10(.11.)))"},
}

func TestRedBlackCase1(t *testing.T) {
    tr := NewTree()
    for _, tt := range fixtureCase1 {
        method := funcs[tt.ops]
        method.Func.Call(ToArgs(tr, tt.kv.key, tt.kv.arg))
        assertEqualTree(tr, t, tt.expected)
    }
}

var fixtureRotationLeft = []struct {
    ops string
    kv  KV
    expected string
}{
    {"put", KV{7, "payload7"}, "(.7.)"},
    {"rotateLeft", KV{}, "(.7.)"},
    {"put", KV{8, "payload8"}, "(.7(.8.))"},
    {"put", KV{9, "payload9"}, "((.7.)8(.9.))"},
    {"put", KV{11, "payload11"}, "((.7.)8(.9(.11.)))"},
    {"rotateLeft", KV{}, "(((.7.)8.)9(.11.))"},
    {"rotateLeft", KV{}, "((((.7.)8.)9.)11.)"},
    {"rotateLeft", KV{}, "((((.7.)8.)9.)11.)"},
}

func TestRedBlackLeft1(t *testing.T) {
    t1 := NewTree()
    assertEqualTree(t1, t, ".")
    t1.RotateLeft(t1.root)
    assertEqualTree(t1, t, ".")

    for _, tt := range fixtureRotationLeft {
        method := funcs[tt.ops]
        switch {
        case tt.ops == "put":
            method.Func.Call(ToArgs(t1, tt.kv.key, tt.kv.arg))

        case tt.ops == "rotateLeft":
            method.Func.Call(ToArgs(t1, t1.root))
        }
        assertEqualTree(t1, t, tt.expected)
    }
}

var fixtureRotationRight = []struct {
    ops string
    kv  KV
    expected string
}{
    {"put", KV{7, "payload7"}, "(.7.)"},
    {"rotateRight", KV{}, "(.7.)"},
    {"put", KV{3, "payload3"}, "((.3.)7.)"},
    {"put", KV{2, "payload2"}, "((.2.)3(.7.))"},
    {"put", KV{1, "payload1"}, "(((.1.)2.)3(.7.))"},
    {"rotateRight", KV{}, "((.1.)2(.3(.7.)))"},
    {"rotateRight", KV{}, "(.1(.2(.3(.7.))))"},
    {"rotateRight", KV{}, "(.1(.2(.3(.7.))))"},
}

func TestRedBlackRight(t *testing.T) {
    t1 := NewTree()
    assertEqualTree(t1, t, ".")
    t1.RotateRight(t1.root)
    assertEqualTree(t1, t, ".")

    for _, tt := range fixtureRotationRight {
        method := funcs[tt.ops]
        switch {
        case tt.ops == "put":
            method.Func.Call(ToArgs(t1, tt.kv.key, tt.kv.arg))

        case tt.ops == "rotateRight":
            method.Func.Call(ToArgs(t1, t1.root))
        }
        assertEqualTree(t1, t, tt.expected)
    }
}

func TestRedBlackParentLookup(t *testing.T) {
    tr := NewTree()

    // Lookup for something in an empty tree
    found, parent, dir := tr.GetParent(5)
    False(found, t)
    Nil(parent, t)
    assertDirection(NODIR, dir, t)

    key7 := 7
    tr.Put(key7, "payload7")
    // Lookup the root node
    found, parent, dir = tr.GetParent(key7)
    True(found, t)
    Nil(parent, t)
    assertDirection(NODIR, dir, t)

    key3, key11 := 3, 11
    tr.Put(key3, "payload3")
    tr.Put(key11, "payload11")

    found, parent, dir = tr.GetParent(key3)
    True(found, t)
    NotNil(parent, t)
    if parent.value != key7 {
		t.Errorf("Expected root node 7")
    }
    assertDirection(LEFT, dir, t)

    found, parent, dir = tr.GetParent(key11)
    True(found, t)
    NotNil(parent, t)
    if parent.value != key7 {
		t.Errorf("Expected root node 7")
    }
    assertDirection(RIGHT, dir, t)
}

var treeData = []struct {
    ops string
    kv  KV
}{
    {"put", KV{7, "payload7"}},
    {"put", KV{3, "payload3"}},
    {"put", KV{18, "payload18"}},
    {"put", KV{10, "payload10"}},
    {"put", KV{8, "payload8"}},
    {"put", KV{11, "payload11"}},
    {"put", KV{22, "payload22"}},
    {"put", KV{26, "payload26"}},
    {"put", KV{30, "payload30"}},
    {"put", KV{45, "payload45"}},
    {"put", KV{35, "payload35"}},
    {"put", KV{90, "payload90"}},
    {"put", KV{85, "payload85"}},
    {"put", KV{83, "payload83"}},
    {"put", KV{100, "payload100"}},
}

func TestRedBlackNodeLookup(t *testing.T) {
    t1 := NewTree()
    for _, tt := range treeData {
        method := funcs[tt.ops]
        switch {
        case tt.ops == "put":
            method.Func.Call(ToArgs(t1, tt.kv.key, tt.kv.arg))
        }
    }

    key10 := 10
    assertEqualTree(t1, t, "(((.3.)7(.8.))10(((.11.)18(.22.))26((.30.)35((.45(.83.))85(.90(.100.))))))")
    assertNodeKey(t1.root, key10, t)

    // search for the root
    {
        ok, node10 := t1.Get(key10)
        True(ok, t)
        NotNil(node10, t)
        assertNodeKey(node10, key10, t)
    }

    // search for non-existent node
    {
        ok, node6 := t1.Get(6)
        False(ok, t)
        Nil(node6, t)
    }

    // search for nodes that exist
    {
        key3, key100 := 3, 100

        ok, node3 := t1.Get(key3)
        True(ok, t)
        NotNil(node3, t)
        assertNodeKey(node3, key3, t)

        ok, node100 := t1.Get(key100)
        True(ok, t)
        NotNil(node100, t)
        assertNodeKey(node100, key100, t)
    }
}

// Removed prefix `Ignore` to test it.
func IgnoreTestLeftRotateProperly(t *testing.T) {
    t1 := NewTree()
    for i, tt := range treeData {
        if i == 9 {
            break
        }
        method := funcs[tt.ops]
        switch {
        case tt.ops == "put":
            method.Func.Call(ToArgs(t1, tt.kv.key, tt.kv.arg))
        }
    }

    key18 := 18
    ok, node18 := t1.Get(key18)
    True(ok, t)
    NotNil(node18, t)
    assertNodeKey(node18, key18, t)

    /*
      (n) = black
             (10)
            /    \
           7     18
          / \   /  \
        (3) (8)(11)(26)
                    / \
                   22  30
    */
    assertEqualTree(t1, t, "(((.3.)7(.8.))10((.11.)18((.22.)26(.30.))))")

    t1.RotateLeft(node18)
    assertEqualTree(t1, t, "(((.3.)7(.8.))10(((.11.)18(.22.))26(.30.)))")
}
