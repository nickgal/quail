package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/log"
)

type sphereListDef struct {
	nameRef     int32
	flags       uint32
	sphereCount uint32
	radius      float32
	scale       float32
	spheres     []common.Quad4
}

func (e *WLD) sphereListDefRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &sphereListDef{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.nameRef = dec.Int32()
	def.flags = dec.Uint32()
	def.sphereCount = dec.Uint32()
	def.radius = dec.Float32()
	def.scale = dec.Float32()
	for i := uint32(0); i < def.sphereCount; i++ {
		var sphere common.Quad4
		sphere.X = dec.Float32()
		sphere.Y = dec.Float32()
		sphere.Z = dec.Float32()
		sphere.W = dec.Float32()
		def.spheres = append(def.spheres, sphere)
	}

	if dec.Error() != nil {
		return dec.Error()
	}

	log.Debugf("%+v", def)
	e.Fragments[fragmentOffset] = def
	return nil
}

func (v *sphereListDef) build(e *WLD) error {
	return nil
}

func (e *WLD) sphereListDefWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
