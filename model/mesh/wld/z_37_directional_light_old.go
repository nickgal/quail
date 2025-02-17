package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

type directionalLightOld struct {
}

func (e *WLD) directionalLightOldRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &directionalLightOld{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	if dec.Error() != nil {
		return fmt.Errorf("directionalLightOldRead: %v", dec.Error())
	}

	log.Debugf("%+v", def)
	e.Fragments[fragmentOffset] = def
	return nil
}

func (v *directionalLightOld) build(e *WLD) error {
	return nil
}

func (e *WLD) directionalLightOldWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
