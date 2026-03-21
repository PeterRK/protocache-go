package test

import (
	"bytes"
	"os"
	"testing"

	"github.com/peterrk/protocache-go"
	"github.com/peterrk/protocache-go/test/pb"
	"github.com/peterrk/protocache-go/test/pc"
	"google.golang.org/protobuf/encoding/protojson"
)

func loadMain(t *testing.T) []byte {
	raw, err := os.ReadFile("test.json")
	assert(t, err == nil)

	message := &pb.Main{}
	err = protojson.Unmarshal(raw, message)
	assert(t, err == nil)

	raw, err = protocache.Serialize(message)
	assert(t, err == nil)
	return raw
}

func TestEXReplayWithoutAccess(t *testing.T) {
	raw := loadMain(t)

	root := pc.TO_MainEX(raw)
	out, err := root.Serialize()
	assert(t, err == nil)
	assert(t, bytes.Equal(raw, out))
}

func TestEXModifyScalarsAndNestedFields(t *testing.T) {
	raw := loadMain(t)

	root := pc.TO_MainEX(raw)
	root.SetStr("patched")
	root.SetI32V([]int32{9, 8, 7})
	root.SetStrv([]string{"left", "right"})
	root.SetModev([]pc.Mode{pc.Mode_MODE_A, pc.Mode_MODE_B})

	obj := root.GetObject()
	obj.SetStr("nested")
	obj.SetI32(777)

	out, err := root.Serialize()
	assert(t, err == nil)

	view := pc.AS_Main(out)
	assert(t, view.GetStr() == "patched")
	assert(t, len(view.GetI32V()) == 3)
	assert(t, view.GetI32V()[0] == 9 && view.GetI32V()[2] == 7)

	strv := view.GetStrv()
	assert(t, strv.Size() == 2)
	assert(t, strv.Get(0) == "left")
	assert(t, strv.Get(1) == "right")

	modev := view.GetModev()
	assert(t, len(modev) == 2)
	assert(t, modev[0] == pc.Mode_MODE_A)
	assert(t, modev[1] == pc.Mode_MODE_B)

	objView := view.GetObject()
	assert(t, objView.GetStr() == "nested")
	assert(t, objView.GetI32() == 777)
}

func TestEXModifyCollections(t *testing.T) {
	raw := loadMain(t)

	root := pc.TO_MainEX(raw)

	items := root.GetObjectv()
	assert(t, len(items) == 3)
	items[0].SetStr("first")
	root.SetObjectv(items[:2])

	index := root.GetIndex()
	index["abc-1"] = 101
	index["new-key"] = 202
	root.SetIndex(index)

	objects := root.GetObjects()
	objects[2].SetStr("updated")
	objects[9] = pc.TO_SmallEX(nil)
	objects[9].SetI32(909)
	root.SetObjects(objects)

	out, err := root.Serialize()
	assert(t, err == nil)

	view := pc.AS_Main(out)

	objs := view.GetObjectv()
	assert(t, objs.Size() == 2)
	obj0 := objs.Get(0)
	obj1 := objs.Get(1)
	assert(t, obj0.GetStr() == "first")
	assert(t, obj1.GetFlag())

	gotIndex := view.GetIndex()
	val, found := gotIndex.Find("new-key")
	assert(t, found && val == 202)
	val, found = gotIndex.Find("abc-1")
	assert(t, found && val == 101)

	gotObjects := view.GetObjects()
	obj2, found := gotObjects.Find(2)
	assert(t, found && obj2.GetStr() == "updated")
	obj9, found := gotObjects.Find(9)
	assert(t, found && obj9.GetI32() == 909)
}

func TestEXModifyAliasFields(t *testing.T) {
	raw := loadMain(t)

	root := pc.TO_MainEX(raw)

	matrix := root.GetMatrix()
	matrix[0][0] = 101
	matrix = append(matrix, pc.Vec2D_Vec1DEX{10, 11})
	root.SetMatrix(matrix)

	vector := root.GetVector()
	vector[0]["lv1"] = pc.ArrMap_ArrayEX{111, 112}
	vector = append(vector, pc.ArrMapEX{
		"lv9": pc.ArrMap_ArrayEX{901, 902},
	})
	root.SetVector(vector)

	arrays := root.GetArrays()
	arrays["lv4"] = pc.ArrMap_ArrayEX{401, 402, 403}
	arrays["lv6"] = pc.ArrMap_ArrayEX{601, 602}
	root.SetArrays(arrays)

	out, err := root.Serialize()
	assert(t, err == nil)

	view := pc.AS_Main(out)

	gotMatrix := view.GetMatrix()
	assert(t, gotMatrix.Size() == 4)
	row0 := gotMatrix.Get(0)
	row3 := gotMatrix.Get(3)
	assert(t, row0.Get(0) == 101)
	assert(t, row3.Size() == 2)
	assert(t, row3.Get(1) == 11)

	gotVector := view.GetVector()
	assert(t, gotVector.Size() == 3)
	vec0 := gotVector.Get(0)
	lv1, found := vec0.Find("lv1")
	assert(t, found)
	assert(t, lv1.Size() == 2)
	assert(t, lv1.Get(0) == 111)
	vec2 := gotVector.Get(2)
	lv9, found := vec2.Find("lv9")
	assert(t, found)
	assert(t, lv9.Size() == 2)
	assert(t, lv9.Get(1) == 902)

	gotArrays := view.GetArrays()
	lv4, found := gotArrays.Find("lv4")
	assert(t, found)
	assert(t, lv4.Size() == 3)
	assert(t, lv4.Get(2) == 403)
	lv6, found := gotArrays.Find("lv6")
	assert(t, found)
	assert(t, lv6.Size() == 2)
	assert(t, lv6.Get(0) == 601)
}

func TestEXNewObjectWithoutBaseCanSerialize(t *testing.T) {
	root := pc.TO_MainEX(nil)
	assert(t, !root.HasSource())

	root.SetI32(321)
	root.SetStr("fresh")
	root.SetModev([]pc.Mode{pc.Mode_MODE_A})
	root.SetObjects(map[int32]*pc.SmallEX{
		1: pc.TO_SmallEX(nil),
	})
	root.GetObjects()[1].SetStr("x")

	out, err := root.Serialize()
	assert(t, err == nil)

	view := pc.AS_Main(out)
	assert(t, view.GetI32() == 321)
	assert(t, view.GetStr() == "fresh")

	modev := view.GetModev()
	assert(t, len(modev) == 1)
	assert(t, modev[0] == pc.Mode_MODE_A)

	gotObjects := view.GetObjects()
	obj, found := gotObjects.Find(1)
	assert(t, found)
	assert(t, obj.GetStr() == "x")
}
