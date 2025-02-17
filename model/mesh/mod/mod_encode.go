package mod

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/tag"
)

// Encode writes a mod file
func Encode(model *common.Model, version uint32, w io.Writer) error {
	var err error
	modelNames := []string{}

	if len(model.Bones) > 0 {
		modelNames = append(modelNames, model.Name)
	}

	names, nameData, err := model.NameBuild(modelNames)
	if err != nil {
		return fmt.Errorf("nameBuild: %w", err)
	}

	materialData, err := model.MaterialBuild(names)
	if err != nil {
		return fmt.Errorf("materialBuild: %w", err)
	}

	verticesData, err := model.VertexBuild(version, names)
	if err != nil {
		return fmt.Errorf("vertexBuild: %w", err)
	}

	triangleData, err := model.TriangleBuild(version, names)
	if err != nil {
		return fmt.Errorf("triangleBuild: %w", err)
	}

	boneData, err := model.BoneBuild(version, "mod", names)
	if err != nil {
		return fmt.Errorf("boneBuild: %w", err)
	}

	tag.New()
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.String("EQGM")
	enc.Uint32(version)
	enc.Uint32(uint32(len(nameData)))
	enc.Uint32(uint32(len(model.Materials)))
	enc.Uint32(uint32(len(model.Vertices)))
	enc.Uint32(uint32(len(model.Triangles)))
	enc.Uint32(uint32(len(model.Bones)))
	tag.Add(0, int(enc.Pos()-1), "red", "header")
	enc.Bytes(nameData)
	tag.Add(tag.LastPos(), int(enc.Pos()), "green", "names")
	enc.Bytes(materialData)
	tag.Add(tag.LastPos(), int(enc.Pos()), "blue", "materials")
	enc.Bytes(verticesData)
	tag.Add(tag.LastPos(), int(enc.Pos()), "yellow", "vertices")
	enc.Bytes(triangleData)
	tag.Add(tag.LastPos(), int(enc.Pos()), "purple", "triangles")
	enc.Bytes(boneData)
	tag.Add(tag.LastPos(), int(enc.Pos()), "orange", "bones")

	err = enc.Error()
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	log.Debugf("%s encoded %d verts, %d triangles, %d bones, %d materials", model.Name, len(model.Vertices), len(model.Triangles), len(model.Bones), len(model.Materials))
	return nil
}
