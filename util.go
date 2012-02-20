// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

func maxi(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func mini(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func boolToInt(value bool) int {
	if value {
		return 1
	}

	return 0
}
