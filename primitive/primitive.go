package primitive

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// Mode defines the shapes used when transforming images
type Mode int

// Modes supported by the primitive package
const (
	ModeCombo Mode = iota
	ModeTriangle
	ModeRect
	ModeEllipse
	ModeCircle
	ModeRotatedRect
	ModeBeziers
	ModeRotatedESllipse
	ModePolygon
)

// WithMode is a function with the Transform function that will define what
// mode will be used. By default, ModeTriangle will be used.
func WithMode(m Mode) func() []string {
	return func() []string {
		return []string{"-n", fmt.Sprintf("%d", m)}
	}
}

// Transform will take the provided image and provide a primitive transformation.
// It will return a reader to the finish image.
func Transform(image io.Reader, numShapes int, opts ...func() []string) (io.Reader, error) {
	in, err := ioutil.TempFile("", "in_")
	if err != nil {
		return nil, err
	}
	defer os.Remove(in.Name())

	out, err := ioutil.TempFile("", "out_")
	if err != nil {
		return nil, err
	}
	defer os.Remove(out.Name())

	// Read image into in Failed
	_, err = io.Copy(in, image)
	if err != nil {
		return nil, err
	}

	// Run primitive w/ -i in.Name() -o out.Name()
	stdCombo, err := primitive(in.Name(), out.Name(), numShapes, ModeCombo)
	if err != nil {
		return nil, err
	}
	fmt.Println(stdCombo)

	// Read out into a reader, return reader, delete out
	b := bytes.NewBuffer(nil)
	_, err = io.Copy(b, out)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func primitive(inputFile string, outputFile string, numShapes int, mode Mode) (string, error) {
	argStr := fmt.Sprintf("-i %s -o %s -n %d -m %d", inputFile, outputFile, numShapes, mode)
	cmd := exec.Command("primitive", strings.Fields(argStr)...)
	d, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(d[:]), nil
}
