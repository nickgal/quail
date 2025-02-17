package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

type blitSpriteDef struct {
}

func (e *WLD) blitSpriteDefRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &blitSpriteDef{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	if dec.Error() != nil {
		return fmt.Errorf("blitSpriteDefRead: %v", dec.Error())
	}

	log.Debugf("%+v", def)
	e.Fragments[fragmentOffset] = def
	return nil
}

func (v *blitSpriteDef) build(e *WLD) error {
	return nil
}

func (e *WLD) blitSpriteDefWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
