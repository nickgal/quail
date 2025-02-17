package quail

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/log"
	"github.com/xackery/quail/model"
	"github.com/xackery/quail/model/metadata/ani"
	"github.com/xackery/quail/model/metadata/prt"
	"github.com/xackery/quail/model/metadata/pts"
	"github.com/xackery/quail/model/metadata/zon"
	"github.com/xackery/quail/pfs/eqg"
	"github.com/xackery/quail/tag"
)

// Export exports the quail target
func (e *Quail) PFSExport(fileVersion uint32, pfsVersion int, path string) error {
	if len(path) == 0 {
		return fmt.Errorf("path is empty")
	}
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".eqg":
		return e.EQGExport(fileVersion, pfsVersion, path)
	case ".s3d":
		return e.S3DExport(fileVersion, pfsVersion, path)
	default:
		if len(ext) < 2 {
			return fmt.Errorf("unknown pfs type %s, valid options are eqg and pfs", path)
		}

		return fmt.Errorf("unknown pfs type %s, valid options are eqg and pfs", ext[1:])
	}
}

// EQGExport exports the quail target to an EQG file
func (e *Quail) EQGExport(fileVersion uint32, pfsVersion int, path string) error {
	pfs, err := eqg.New(path)
	if err != nil {
		return fmt.Errorf("eqg new: %w", err)
	}
	defer pfs.Close()

	if e.Zone != nil {
		buf := &bytes.Buffer{}
		err = zon.Encode(e.Zone, fileVersion, buf)
		if err != nil {
			return fmt.Errorf("encodeZon %s: %w", e.Zone.Name, err)
		}
		os.WriteFile(fmt.Sprintf("%s/%s-raw-out.zon", "testdata", e.Zone.Name), buf.Bytes(), 0644)
		tag.Write(fmt.Sprintf("%s/%s-raw-out.zon.tags", "testdata", e.Zone.Name))

		err = pfs.Add(fmt.Sprintf("%s.zon", e.Zone.Name), buf.Bytes())
		if err != nil {
			return fmt.Errorf("addZon %s: %w", e.Zone.Name, err)
		}
	}

	for _, entry := range e.Models {
		buf := &bytes.Buffer{}
		err = model.Encode(entry, fileVersion, buf)
		if err != nil {
			return fmt.Errorf("encodeMod %s: %w", entry.Name, err)
		}

		os.WriteFile(fmt.Sprintf("%s/%s-raw-out.%s", "testdata", entry.Name, entry.FileType), buf.Bytes(), 0644)
		tag.Write(fmt.Sprintf("%s/%s-raw-out.%s.tags", "testdata", entry.Name, entry.FileType))

		err = pfs.Add(fmt.Sprintf("%s.%s", entry.Name, entry.FileType), buf.Bytes())
		if err != nil {
			return fmt.Errorf("addMod %s: %w", entry.Name, err)
		}
		for _, material := range entry.Materials {
			for _, property := range material.Properties {
				if len(property.Data) == 0 {
					continue
				}
				err = pfs.Add(property.Value, property.Data)
				if err != nil {
					return fmt.Errorf("model %s addMaterial %s texture %s: %w", entry.Name, material.Name, property.Value, err)
				}
			}
		}
	}

	for _, anim := range e.Animations {
		buf := &bytes.Buffer{}
		err = ani.Encode(anim, fileVersion, buf)
		if err != nil {
			return fmt.Errorf("encodeAni %s: %w", anim.Name, err)
		}
		err = pfs.Add(fmt.Sprintf("%s.ani", anim.Name), buf.Bytes())
		if err != nil {
			return fmt.Errorf("addMds %s: %w", anim.Name, err)
		}
	}

	for _, model := range e.Models {
		for _, render := range model.ParticleRenders {
			buf := &bytes.Buffer{}
			err = prt.Encode(render, 4, buf) // TODO: add support for other versions
			if err != nil {
				return fmt.Errorf("encodePtr %s: %w", render.Name, err)
			}
			err = pfs.Add(fmt.Sprintf("%s.prt", render.Name), buf.Bytes())
			if err != nil {
				return fmt.Errorf("addPtr %s: %w", render.Name, err)
			}
		}

		for _, point := range model.ParticlePoints {
			buf := &bytes.Buffer{}
			err = pts.Encode(point, fileVersion, buf)
			if err != nil {
				return fmt.Errorf("encodePts %s: %w", point.Name, err)
			}
			err = pfs.Add(fmt.Sprintf("%s.pts", point.Name), buf.Bytes())
			if err != nil {
				return fmt.Errorf("addPts %s: %w", point.Name, err)
			}
		}
	}

	w, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	defer w.Close()

	err = pfs.Encode(w)
	if err != nil {
		return fmt.Errorf("encode %s: %w", path, err)
	}
	log.Debugf("wrote %s with %d entries", path, pfs.Len())
	return nil
}

// S3DExport exports the quail target to an S3D file
func (e *Quail) S3DExport(fileVersion uint32, pfsVersion int, path string) error {
	return nil
}
