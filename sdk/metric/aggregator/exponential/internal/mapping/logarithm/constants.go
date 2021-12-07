// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logarithm

// maxBoundaries is the exact lower boundary value corresponding with
// the largest bucket mapped by finite float64 values up to
// math.MaxFloat64.  This is offset by one, since the logarithm
// implementation is used with scale >= 1.
var prebuiltMappings = [MaxScale]logarithmMapping{
	{
		scale:	       1,
		minIndex:      -0x7fc,
		maxIndex:      0x7ff,
		maxBoundary:   0x1.6a09e667f3bcdp+1023,
		scaleFactor:   0x1.71547652b82fep+01,
		inverseFactor: 0x1.62e42fefa39efp-02,
	},
	{
		scale:	       2,
		minIndex:      -0xff8,
		maxIndex:      0xfff,
		maxBoundary:   0x1.ae89f995ad3adp+1023,
		scaleFactor:   0x1.71547652b82fep+02,
		inverseFactor: 0x1.62e42fefa39efp-03,
	},
	{
		scale:	       3,
		minIndex:      -0x1ff0,
		maxIndex:      0x1fff,
		maxBoundary:   0x1.d5818dcfba487p+1023,
		scaleFactor:   0x1.71547652b82fep+03,
		inverseFactor: 0x1.62e42fefa39efp-04,
	},
	{
		scale:	       4,
		minIndex:      -0x3fe0,
		maxIndex:      0x3fff,
		maxBoundary:   0x1.ea4afa2a490dap+1023,
		scaleFactor:   0x1.71547652b82fep+04,
		inverseFactor: 0x1.62e42fefa39efp-05,
	},
	{
		scale:	       5,
		minIndex:      -0x7fc0,
		maxIndex:      0x7fff,
		maxBoundary:   0x1.f50765b6e454p+1023,
		scaleFactor:   0x1.71547652b82fep+05,
		inverseFactor: 0x1.62e42fefa39efp-06,
	},
	{
		scale:	       6,
		minIndex:      -0xff80,
		maxIndex:      0xffff,
		maxBoundary:   0x1.fa7c1819e90d8p+1023,
		scaleFactor:   0x1.71547652b82fep+06,
		inverseFactor: 0x1.62e42fefa39efp-07,
	},
	{
		scale:	       7,
		minIndex:      -0x1ff00,
		maxIndex:      0x1ffff,
		maxBoundary:   0x1.fd3c22b8f71f1p+1023,
		scaleFactor:   0x1.71547652b82fep+07,
		inverseFactor: 0x1.62e42fefa39efp-08,
	},
	{
		scale:	       8,
		minIndex:      -0x3fe00,
		maxIndex:      0x3ffff,
		maxBoundary:   0x1.fe9d96b2a23d9p+1023,
		scaleFactor:   0x1.71547652b82fep+08,
		inverseFactor: 0x1.62e42fefa39efp-09,
	},
	{
		scale:	       9,
		minIndex:      -0x7fc00,
		maxIndex:      0x7ffff,
		maxBoundary:   0x1.ff4eaca4391b6p+1023,
		scaleFactor:   0x1.71547652b82fep+09,
		inverseFactor: 0x1.62e42fefa39efp-10,
	},
	{
		scale:	       10,
		minIndex:      -0xff800,
		maxIndex:      0xfffff,
		maxBoundary:   0x1.ffa74ea381efcp+1023,
		scaleFactor:   0x1.71547652b82fep+10,
		inverseFactor: 0x1.62e42fefa39efp-11,
	},
	{
		scale:	       11,
		minIndex:      -0x1ff000,
		maxIndex:      0x1fffff,
		maxBoundary:   0x1.ffd3a565efb65p+1023,
		scaleFactor:   0x1.71547652b82fep+11,
		inverseFactor: 0x1.62e42fefa39efp-12,
	},
	{
		scale:	       12,
		minIndex:      -0x3fe000,
		maxIndex:      0x3fffff,
		maxBoundary:   0x1.ffe9d237fe372p+1023,
		scaleFactor:   0x1.71547652b82fep+12,
		inverseFactor: 0x1.62e42fefa39efp-13,
	},
	{
		scale:	       13,
		minIndex:      -0x7fc000,
		maxIndex:      0x7fffff,
		maxBoundary:   0x1.fff4e8fd40081p+1023,
		scaleFactor:   0x1.71547652b82fep+13,
		inverseFactor: 0x1.62e42fefa39efp-14,
	},
	{
		scale:	       14,
		minIndex:      -0xff8000,
		maxIndex:      0xffffff,
		maxBoundary:   0x1.fffa7476f029dp+1023,
		scaleFactor:   0x1.71547652b82fep+14,
		inverseFactor: 0x1.62e42fefa39efp-15,
	},
	{
		scale:	       15,
		minIndex:      -0x1ff0000,
		maxIndex:      0x1ffffff,
		maxBoundary:   0x1.fffd3a398c1bbp+1023,
		scaleFactor:   0x1.71547652b82fep+15,
		inverseFactor: 0x1.62e42fefa39efp-16,
	},
	{
		scale:	       16,
		minIndex:      -0x3fe0000,
		maxIndex:      0x3ffffff,
		maxBoundary:   0x1.fffe9d1c4b0f3p+1023,
		scaleFactor:   0x1.71547652b82fep+16,
		inverseFactor: 0x1.62e42fefa39efp-17,
	},
	{
		scale:	       17,
		minIndex:      -0x7fc0000,
		maxIndex:      0x7ffffff,
		maxBoundary:   0x1.ffff4e8e06c7fp+1023,
		scaleFactor:   0x1.71547652b82fep+17,
		inverseFactor: 0x1.62e42fefa39efp-18,
	},
	{
		scale:	       18,
		minIndex:      -0xff80000,
		maxIndex:      0xfffffff,
		maxBoundary:   0x1.ffffa746fbb4p+1023,
		scaleFactor:   0x1.71547652b82fep+18,
		inverseFactor: 0x1.62e42fefa39efp-19,
	},
	{
		scale:	       19,
		minIndex:      -0x1ff00000,
		maxIndex:      0x1fffffff,
		maxBoundary:   0x1.ffffd3a37beep+1023,
		scaleFactor:   0x1.71547652b82fep+19,
		inverseFactor: 0x1.62e42fefa39efp-20,
	},
	{
		scale:	       20,
		minIndex:      -0x3fe00000,
		maxIndex:      0x3fffffff,
		maxBoundary:   0x1.ffffe9d1bd7cp+1023,
		scaleFactor:   0x1.71547652b82fep+20,
		inverseFactor: 0x1.62e42fefa39efp-21,
	},
}
