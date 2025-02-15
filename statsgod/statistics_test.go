/**
 * Copyright 2015 Acquia, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package statsgod_test

import (
	"math"
	"sort"

	. "github.com/acquia/statsgod/statsgod"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var statsBenchmarkTimeLimit = 0.25

var _ = Describe("Statistics", func() {

	Describe("Testing the ValueSlice struct", func() {
		values := ValueSlice{1, 5, 2, 4, 3}

		It("should contain values", func() {
			// Ensure we can get a correct length.
			Expect(values.Len()).Should(Equal(5))
			// Ensure we can get a value at an index.
			Expect(values.Get(3)).Should(Equal(float64(4)))
			// Test the sort interface.
			sort.Sort(values)
			Expect(values.Get(0)).Should(Equal(float64(1)))
			Expect(values.Get(1)).Should(Equal(float64(2)))
			Expect(values.Get(2)).Should(Equal(float64(3)))
			Expect(values.Get(3)).Should(Equal(float64(4)))
			Expect(values.Get(4)).Should(Equal(float64(5)))

		})
	})

	Describe("Testing the statistics calculations", func() {
		Context("when the UniqueCount is applied", func() {
			values := ValueSlice{4, 5, 2, 3, 2, 4, 4, 5, 6, 7, 8, 9, 0, 4}
			count := values.UniqueCount()
			It("should find the correct number of unique values", func() {
				Expect(count).Should(Equal(9))
			})
		})

		Context("when the Minmax is applied", func() {
			values := ValueSlice{5, math.NaN(), 2, 3, 4, 1}
			min, max, err := values.Minmax()
			It("should find the min/max values", func() {
				Expect(min).Should(Equal(float64(1)))
				Expect(max).Should(Equal(float64(5)))
				Expect(err).ShouldNot(BeNil())
			})

			Measure("it should find min/max quickly.", func(b Benchmarker) {
				runtime := b.Time("runtime", func() {
					min, max, _ = values.Minmax()
				})
				Expect(runtime.Seconds()).Should(BeNumerically("<", statsBenchmarkTimeLimit), "it should find min/max quickly.")
			}, 100000)

		})

		Context("when the Median is applied", func() {
			values := ValueSlice{123, 234, 345, 456, 567, 678, 789}
			median := values.Median()

			It("should find the median value", func() {
				Expect(median).Should(Equal(float64(456)))

				values = append(values, 890)
				median = values.Median()
				Expect(median).Should(Equal(float64(511.5)))

				valuesOriginal := values
				values = ValueSlice{}
				median = values.Median()
				Expect(median).Should(Equal(float64(0.0)))
				values = valuesOriginal
			})

			Measure("it should find median quickly.", func(b Benchmarker) {
				runtime := b.Time("runtime", func() {
					median = values.Median()
				})
				Expect(runtime.Seconds()).Should(BeNumerically("<", statsBenchmarkTimeLimit), "it should find median quickly.")
			}, 100000)

		})

		Context("when the Mean is applied", func() {
			values := ValueSlice{123, 234, 345, 456, 567, 678, 789, 890}
			mean := values.Mean()
			It("should find the mean value", func() {
				Expect(mean).Should(Equal(float64(510.25)))
			})

			Measure("it should find mean quickly.", func(b Benchmarker) {
				runtime := b.Time("runtime", func() {
					mean = values.Mean()
				})
				Expect(runtime.Seconds()).Should(BeNumerically("<", statsBenchmarkTimeLimit), "it should find mean quickly.")
			}, 100000)

		})

		Context("when the Quantile is applied", func() {
			values := ValueSlice{123, 234, 345, 456, 567, 678, 789, 890, 910, 1011}
			It("should find the requested quantile value", func() {
				q100 := values.Quantile(1)
				Expect(q100).Should(Equal(float64(1011)))
				q90 := values.Quantile(0.9)
				Expect(q90).Should(Equal(float64(920.1)))
				q80 := values.Quantile(0.8)
				Expect(q80).Should(Equal(float64(894)))
				q75 := values.Quantile(0.75)
				Expect(q75).Should(Equal(float64(864.75)))
				q50 := values.Quantile(0.5)
				Expect(q50).Should(Equal(float64(622.5)))
				q25 := values.Quantile(0.25)
				Expect(q25).Should(Equal(float64(372.75)))

				valuesOriginal := values
				values = ValueSlice{}
				q0 := values.Quantile(1)
				Expect(q0).Should(Equal(float64(0.0)))
				values = valuesOriginal
			})

			testQuantile := float64(0.0)
			Measure("it should find quantile quickly.", func(b Benchmarker) {
				runtime := b.Time("runtime", func() {
					testQuantile = values.Quantile(0.9)
					_ = testQuantile
				})
				Expect(runtime.Seconds()).Should(BeNumerically("<", statsBenchmarkTimeLimit), "it should find quantile quickly.")
			}, 100000)

		})

		Context("when the sum is applied", func() {
			values := ValueSlice{123, 234, 345, 456, 567, 678, 789, 890, 910, 1011}
			It("should find the correct sum value", func() {
				sum := values.Sum()
				Expect(sum).Should(Equal(float64(6003.0)))
			})
		})
	})
})
