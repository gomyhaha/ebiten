// Copyright 2019 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package restorable

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/internal/graphicscommand"
)

type rectToPixels struct {
	m map[image.Rectangle][]byte

	last image.Rectangle
}

func (rtp *rectToPixels) addOrReplace(pixels []byte, x, y, width, height int) {
	if len(pixels) != 4*width*height {
		panic(fmt.Sprintf("restorable: len(pixels) must be %d but %d", 4*width*height, len(pixels)))
	}

	if rtp.m == nil {
		rtp.m = map[image.Rectangle][]byte{}
	}

	copied := make([]byte, len(pixels))
	copy(copied, pixels)

	newr := image.Rect(x, y, x+width, y+height)
	for r := range rtp.m {
		if r == newr {
			// Replace the region.
			rtp.m[r] = copied
			return
		}
		if r.Overlaps(newr) {
			panic(fmt.Sprintf("restorable: region (%#v) conflicted with the other region (%#v)", newr, r))
		}
	}

	// Add the region.
	rtp.m[newr] = copied
}

func (rtp *rectToPixels) remove(x, y, width, height int) {
	if rtp.m == nil {
		return
	}

	newr := image.Rect(x, y, x+width, y+height)
	for r := range rtp.m {
		if r == newr {
			delete(rtp.m, r)
			return
		}
	}
}

func (rtp *rectToPixels) at(i, j int) (byte, byte, byte, byte, bool) {
	if rtp.m == nil {
		return 0, 0, 0, 0, false
	}

	var r *image.Rectangle
	pt := image.Pt(i, j)
	if pt.In(rtp.last) {
		r = &rtp.last
	} else {
		for rr := range rtp.m {
			if pt.In(rr) {
				r = &rr
				rtp.last = rr
				break
			}
		}
	}

	if r == nil {
		return 0, 0, 0, 0, false
	}

	pix := rtp.m[*r]
	idx := 4 * ((j-r.Min.Y)*r.Dx() + (i - r.Min.X))
	return pix[idx], pix[idx+1], pix[idx+2], pix[idx+3], true
}

func (rtp *rectToPixels) apply(img *graphicscommand.Image) {
	// TODO: Isn't this too heavy? Can we merge the operations?
	for r, pix := range rtp.m {
		img.ReplacePixels(pix, r.Min.X, r.Min.Y, r.Dx(), r.Dy())
	}
}
