// filter.go
// Copyright (C) 2021 Kasai Koji

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package textfilter

type (
	Guaranteed          func(passed string) (done bool, err error)
	GuaranteedAssigning func(previous, passed string) (assigned bool, err error)

	// Interface to check string expression is valid or not.
	// When passed filter, call internalFunc.
	Filter interface {
		Add(chk string, internalFunc Guaranteed) (done bool, err error)
		Assign(previous, chk string, internalFunc GuaranteedAssigning) (done bool, err error)
		Clone(original string, internalFunc Guaranteed) (done bool, err error)
		Remove(deleted string, internalFunc Guaranteed) (done bool, err error)
	}

	// Defines stacked filter instance.
	Multiple []Filter

	// Defines validation filter function.
	Validator func(chk string) error

	// Defines limited cloning expression filter instance.
	LimitedCloning struct {
		cloneLimitter int
		mapping       map[string]int
	}
)

// Deep clone string expression.
func CloneString(src string) (dest string) {
	buffer := make([]byte, len(src)+1)
	copy(buffer, src)
	return string(buffer)
}

func RegisterFiltering(filter Filter, chkStr string, whenSuccess func() error) error {

	_, err := filter.Add(chkStr, func(passed string) (done bool, err error) {
		return true, whenSuccess()
	})

	return err
}

func RegisterFilteringPairs(filters []Filter, chkStrs []string, whenSuccess func() error) error {
	if len(filters) == 0 {
		panic("no filters in arguments")
	} else if len(filters) != len(chkStrs) {
		panic("length of filters and chkStrs arguments are not equal")
	}
	_, err := regPairs(filters, chkStrs, whenSuccess, 0)
	return err
}

func regPairs(filters []Filter, chkStrs []string, whenSuccess func() error, idx int) (bool, error) {
	if idx < len(filters) {
		return filters[idx].Add(chkStrs[idx], func(passed string) (done bool, err error) {
			return regPairs(filters, chkStrs, whenSuccess, idx+1)
		})
	} else {
		return filters[idx].Add(chkStrs[idx], func(passed string) (done bool, err error) {
			return true, whenSuccess()
		})
	}
}
