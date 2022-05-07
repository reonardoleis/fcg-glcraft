package lib

// Utilizamos este arquivo para carregar OBJ pois vimos
// que em C poderia utilizar a biblioteca tinyobjloader
// FONTE:  https://gist.github.com/davemackintosh/67959fa9dfd9018d79a4

import (
	"fmt"
	"io"
	"os"

	"github.com/go-gl/mathgl/mgl32"
)

// Model is a renderable collection of vecs.
type Model struct {
	// For the v, vt and vn in the obj file.
	Normals, Vecs []mgl32.Vec3
	Uvs           []mgl32.Vec2

	// For the fun "f" in the obj file.
	VecIndices, NormalIndices, UvIndices []float32
}

// NewModel will read an OBJ model file and create a Model from its contents
func NewModel(file string) Model {
	// Open the file for reading and check for errors.
	objFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	// Don't forget to close the file reader.
	defer objFile.Close()

	// Create a model to store stuff.
	model := Model{}

	// Read the file and get it's contents.
	for {
		var lineType string

		// Scan the type field.
		_, err := fmt.Fscanf(objFile, "%s", &lineType)

		// Check if it's the end of the file
		// and break out of the loop.
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		// Check the type.
		switch lineType {
		// VERTICES.
		case "v":
			// Create a vec to assign digits to.
			vec := mgl32.Vec3{}
			x := float32(0)

			// Get the digits from the file.
			fmt.Fscanf(objFile, "%f %f %f $f\n", &vec[0], &vec[1], &vec[2], &x)

			// Add the vector to the model.
			model.Vecs = append(model.Vecs, vec)

		// INDICES.
		case "f":
			// Create a vec to assign digits to.

			vec := make([]float32, 3)
			x := float32(0)
			// Get the digits from the file.
			matches, _ := fmt.Fscanf(objFile, "%f %f %f %f\n", &vec[0], &vec[1], &vec[2], &x)

			if matches != 3 {
				panic("Cannot read your file")
			}

			// Add the numbers to the model.

			model.VecIndices = append(model.VecIndices, vec[0])
			model.VecIndices = append(model.VecIndices, vec[1])
			model.VecIndices = append(model.VecIndices, vec[2])

		}
	}

	// Return the newly created Model.
	return model
}

// GetRenderableVertices returns a slice of float32s
// formatted in X, Y, Z, U, V. That is, XYZ of the
// vertex and the texture position.
func (model Model) GetRenderableVertices() []float32 {
	// Create a slice for the outward float32s.
	var out []float32

	// Loop over each vec3 in the indices property.
	for _, position := range model.VecIndices {
		index := int(position) - 1
		vec := model.Vecs[index]

		out = append(out, vec.X(), vec.Y(), vec.Z(), 1.0)
	}

	// Return the array.
	return out
}
