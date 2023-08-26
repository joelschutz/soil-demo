package util

import (
	"fmt"
	"math/rand"

	"github.com/go-gl/mathgl/mgl32"
)

func Merge2Matrixes(matrixes [2][][]float32) ([][]mgl32.Vec2, error) {
	m0 := matrixes[0]
	m1 := matrixes[1]
	lx := len(m0)
	ly := len(m0[0])
	if lx != len(m1) || ly != len(m1[0]) {
		return nil, fmt.Errorf("Must have matrixes of the same size")
	}

	r := [][]mgl32.Vec2{}

	for x := 0; x < lx; x++ {
		row := []mgl32.Vec2{}
		for y := 0; y < ly; y++ {
			row = append(row, mgl32.Vec2{m0[x][y], m1[x][y]})
		}
		r = append(r, row)
	}

	return r, nil
}

func ApplyMaskOnMatrix[V any](dst [][]V, mask [][]bool, value V) ([][]V, error) {
	lx := len(dst)
	ly := len(dst[0])

	if lx != len(mask) || ly != len(mask[0]) {
		return dst, fmt.Errorf("dst and mask must have the same size")
	}

	for x := 0; x < lx; x++ {
		for y := 0; y < ly; y++ {
			if mask[x][y] {
				dst[x][y] = value
			}
		}
	}

	return dst, nil
}

func MakeMatrixWH[V any](width, height int, value V) [][]V {
	m := [][]V{}

	for x := 0; x < width; x++ {
		row := []V{}
		for y := 0; y < height; y++ {
			row = append(row, value)
		}
		m = append(m, row)
	}

	return m
}

func MakeMatrix[V any](size int, value V) [][]V {
	m := [][]V{}

	for x := 0; x < size; x++ {
		row := []V{}
		for y := 0; y < size; y++ {
			row = append(row, value)
		}
		m = append(m, row)
	}

	return m
}

func MakeRandMatrixUint8(size, maxValue int) [][]uint8 {
	m := [][]uint8{}

	for x := 0; x < size; x++ {
		row := []uint8{}
		for y := 0; y < size; y++ {
			row = append(row, uint8(rand.Intn(maxValue)+1))
		}
		m = append(m, row)
	}

	return m
}

func MakeMatrixBool(size int) [][]bool {
	m := [][]bool{}

	for x := 0; x < size; x++ {
		row := []bool{}
		for y := 0; y < size; y++ {
			row = append(row, false)
		}
		m = append(m, row)
	}

	return m
}

func MakeRandMatrixBool(size, threshold int) [][]bool {
	m := [][]bool{}

	for x := 0; x < size; x++ {
		row := []bool{}
		for y := 0; y < size; y++ {
			row = append(row, rand.Intn(101) < threshold)
		}
		m = append(m, row)
	}

	return m
}
