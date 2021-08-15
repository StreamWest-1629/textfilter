// restricted_cloning.go
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

import "errors"

var (
	ErrLimitedCloning = errors.New("restricted cloning or duplicate assign")
)

func Identifier() *LimitedCloning { return LimitClone(0) }
func LimitClone(limit int) *LimitedCloning {
	return &LimitedCloning{
		cloneLimitter: limit,
		mapping:       make(map[string]int),
	}
}

func (limited *LimitedCloning) check(chk string) error {
	if times, exist := limited.mapping[chk]; !exist {
		return nil
	} else if times < limited.cloneLimitter {
		return nil
	} else {
		return ErrLimitedCloning
	}
}

func (limited *LimitedCloning) decl(dec string) {
	if times, exist := limited.mapping[dec]; exist {
		if times > 0 {
			limited.mapping[dec] = times - 1
		} else {
			delete(limited.mapping, dec)
		}
	}
}

func (limited *LimitedCloning) incl(inc string) {
	if times, exist := limited.mapping[inc]; !exist {
		limited.mapping[inc] = 0
	} else {
		limited.mapping[inc] = times + 1
	}
}

func (limited *LimitedCloning) Add(chk string, internal Guaranteed) (done bool, err error) {

	if err := limited.check(chk); err != nil {
		return false, err
	} else {
		done, err := internal(chk)
		if done {
			limited.incl(chk)
		}
		return done, err
	}
}

func (limited *LimitedCloning) Assign(previous, chk string, internal GuaranteedAssigning) (done bool, err error) {

	if previous == chk {
		return internal(previous, chk)
	} else if err := limited.check(chk); err != nil {
		return false, err
	} else {
		done, err := internal(previous, chk)
		if done {
			limited.incl(chk)
			limited.decl(previous)
		}
		return done, err
	}
}

func (limited *LimitedCloning) Clone(original string, internal Guaranteed) (done bool, err error) {

	if err := limited.check(original); err != nil {
		return false, err
	} else {
		done, err := internal(original)
		if done {
			limited.incl(original)
		}
		return done, err
	}
}

func (limited *LimitedCloning) Remove(deleted string, internal Guaranteed) (done bool, err error) {
	done, err = internal(deleted)
	if done {
		limited.decl(deleted)
	}
	return
}
