// validator.go
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

import (
	"errors"
	"regexp"
)

var (
	ErrUnmatchRegexp = errors.New("unmatches with regular expression")
	ErrUnmatchList   = errors.New("unmatches with list")
)

func (v Validator) Add(chk string, internal Guaranteed) (done bool, err error) {
	if err := v(chk); err != nil {
		return false, err
	} else {
		return internal(chk)
	}
}

func (v Validator) Assign(previous, chk string, internal GuaranteedAssigning) (done bool, err error) {
	if err := v(chk); err != nil {
		return false, err
	} else {
		return internal(previous, chk)
	}
}

func (v Validator) Clone(original string, internal Guaranteed) (done bool, err error) {
	return internal(original)
}

func (v Validator) Remove(deleted string, internal Guaranteed) (done bool, err error) {
	return internal(deleted)
}

func RegexpMatches(expression string) Validator {
	regex := regexp.MustCompile(expression)
	regex.Longest()

	return func(chk string) error {
		if regex.MatchString(expression) {
			return nil
		}
		return ErrUnmatchRegexp
	}
}

func ListMatches(allowList ...string) Validator {

	mapping := make(map[string]interface{})

	for _, val := range allowList {
		mapping[val] = nil
	}

	return func(chk string) error {
		if _, exist := mapping[chk]; exist {
			return nil
		} else {
			return ErrUnmatchList
		}
	}

}
