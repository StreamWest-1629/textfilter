// multiple.go
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

func (mlt Multiple) Add(chk string, internal Guaranteed) (done bool, err error) {
	return mlt.add(0, internal)(chk)
}

func (mlt Multiple) add(next int, finally Guaranteed) Guaranteed {
	if next < len(mlt) {
		return func(passed string) (done bool, err error) {
			return mlt[next].Add(passed, mlt.add(next+1, finally))
		}
	} else {
		return finally
	}
}

func (mlt Multiple) Assign(previous, chk string, internal GuaranteedAssigning) (done bool, err error) {
	return mlt.assign(0, internal)(previous, chk)
}

func (mlt Multiple) assign(next int, finally GuaranteedAssigning) GuaranteedAssigning {
	if next < len(mlt) {
		return func(previous, passed string) (assigned bool, err error) {
			return mlt[next].Assign(previous, passed, mlt.assign(next+1, finally))
		}
	} else {
		return finally
	}
}

func (mlt Multiple) Clone(original string, internal Guaranteed) (done bool, err error) {
	return mlt.clone(0, internal)(original)
}

func (mlt Multiple) clone(next int, finally Guaranteed) Guaranteed {
	if next < len(mlt) {
		return func(passed string) (done bool, err error) {
			return mlt[next].Clone(passed, mlt.clone(next+1, finally))
		}
	} else {
		return finally
	}
}

func (mlt Multiple) Remove(deleted string, internal Guaranteed) (done bool, err error) {
	return mlt.remove(0, internal)(deleted)
}

func (mlt Multiple) remove(next int, finally Guaranteed) Guaranteed {
	if next < len(mlt) {
		return func(passed string) (done bool, err error) {
			return mlt[next].Remove(passed, mlt.remove(next+1, finally))
		}
	} else {
		return finally
	}
}
