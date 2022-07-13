package ter

import (
	"fmt"

	"github.com/g3n/engine/math32"
	"github.com/xackery/quail/common"
)

func (e *TER) MaterialAdd(name string, shaderName string) error {
	if shaderName == "" {
		shaderName = "Opaque_MaxCB1.fx"
	}
	e.materials = append(e.materials, &common.Material{
		Name:       name,
		ShaderName: shaderName,
		Properties: []*common.Property{},
	})
	return nil
}

func (e *TER) MaterialPropertyAdd(materialName string, propertyName string, category uint32, value string) error {
	for _, o := range e.materials {
		if o.Name != materialName {
			continue
		}
		o.Properties = append(o.Properties, &common.Property{
			Name:     propertyName,
			Category: category,
			Value:    value,
		})
		return nil
	}
	return fmt.Errorf("materialName not found: %s", materialName)
}

func (e *TER) VertexAdd(position *math32.Vector3, normal *math32.Vector3, uv *math32.Vector2) error {
	e.vertices = append(e.vertices, &common.Vertex{
		Position: position,
		Normal:   normal,
		Uv:       uv,
	})
	return nil
}

func (e *TER) FaceAdd(index [3]uint32, materialName string, flag uint32) error {
	if materialName == "" && len(e.materials) == 0 {
		e.faces = append(e.faces, &common.Face{
			Index:        index,
			MaterialName: materialName,
			Flag:         flag,
		})
		return nil
	}

	for _, o := range e.materials {
		if o.Name != materialName {
			continue
		}

		e.faces = append(e.faces, &common.Face{
			Index:        index,
			MaterialName: materialName,
			Flag:         flag,
		})
		return nil
	}

	return fmt.Errorf("materialName not found: %s", materialName)
}
