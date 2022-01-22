package zon

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/g3n/engine/math32"
)

func TestSave(t *testing.T) {
	var err error
	z := &ZON{}
	err = z.AddModel("test")
	if err != nil {
		t.Fatalf("addModel: %s", err)
	}
	err = z.AddObject("test", "test01", math32.Vector3{X: 1, Y: 2, Z: 3}, math32.Vector3{}, 0)
	if err != nil {
		t.Fatalf("addObject: %s", err)
	}
	buf := bytes.NewBuffer(nil)

	err = z.Save(buf)
	if err != nil {
		t.Fatalf("save: %s", err.Error())
	}
	fmt.Println(hex.Dump(buf.Bytes()))
}

func TestCompare(t *testing.T) {
	compareFile(t, "test/out.zon", "test/soldungb.zon")
}

func compareFile(t *testing.T, path1 string, path2 string) {

	f1, err := os.Open(path1)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer f1.Close()
	f2, err := os.Open(path2)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer f2.Close()
	offset := 0

	fails := 0
	f1Data := []byte{}
	f2Data := []byte{}
	for {

		buf1 := make([]byte, 1)
		buf2 := make([]byte, 1)
		_, err1 := f1.Read(buf1)

		if err1 != nil {
			if err1 == io.EOF {
				break
			}
			buf1[0] = 0
		}

		_, err2 := f2.Read(buf2)
		if err2 != nil {
			if err2 == io.EOF {
				break
			}
			buf2[0] = 0
		}
		f1Data = append(f1Data, buf1[0])
		f2Data = append(f2Data, buf2[0])

		if offset == 0 {
			offset++
			continue
		}
		if buf1[0] != buf2[0] {
			fmt.Println(path1, "\n", hex.Dump([]byte(f1Data)))
			fmt.Println(path2, "\n", hex.Dump([]byte(f2Data)))
			if fails > 0 {
				t.Fatalf("mismatched at position %d (0x%02x) %s has value %d, wanted %d", offset, offset, path1, buf1[0], buf2[0])
			}
			fails++
		}

		offset++
	}
}
