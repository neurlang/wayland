package main

import "github.com/neurlang/gm/mat4"

func Frustum(result *mat4.Mat4, left, right, bottom, top, nearZ, farZ float32) {
	var deltaX = right - left
	var deltaY = top - bottom
	var deltaZ = farZ - nearZ
	var frust [16]float32

	if (nearZ <= 0.0) || (farZ <= 0.0) ||
		(deltaX <= 0.0) || (deltaY <= 0.0) || (deltaZ <= 0.0) {
		return
	}

	frust[0] = 2.0 * nearZ / deltaX
	frust[1] = 0.0
	frust[2] = 0.0
	frust[3] = 0.0

	frust[4+1] = 2.0 * nearZ / deltaY
	frust[4+0] = 0.0
	frust[4+2] = 0.0
	frust[4+3] = 0.0

	frust[8+0] = (right + left) / deltaX
	frust[8+1] = (top + bottom) / deltaY
	frust[8+2] = -farZ / deltaZ
	frust[8+3] = -1.0

	frust[12+2] = -(nearZ * farZ) / deltaZ
	frust[12+0] = 0.0
	frust[12+1] = 0.0
	frust[12+3] = 0.0

	*result = frust
}
