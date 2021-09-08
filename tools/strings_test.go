package tools

import (
	"reflect"
	"sort"
	"testing"
)

func TestJoinStrings(t *testing.T) {
	type joinStringTest struct {
		elem   []string
		output string
	}

	var joinStringTests = []joinStringTest{
		{[]string{"Hello", "World"}, "HelloWorld"},
		{[]string{"Hello"}, "Hello"},
		{[]string{"Hello", " ", "World", "!!!"}, "Hello World!!!"},
	}

	for _, jst := range joinStringTests {
		actual := JoinStrings(jst.elem...)
		if actual != jst.output {
			t.Errorf("JoinStrings(%v): expected %s, actual %s", jst.elem, jst.output, actual)
		}
	}
}

func TestJoinStringsSep(t *testing.T) {
	type joinStringTest struct {
		separator string
		elem      []string
		output    string
	}

	var joinStringTests = []joinStringTest{
		{"-", []string{"Hello", "World"}, "Hello-World"},
		{"|", []string{"Hello"}, "Hello"},
		{" ", []string{"Hello", "World", "!!!"}, "Hello World !!!"},
	}

	for _, jst := range joinStringTests {
		actual := JoinStringsSep(jst.separator, jst.elem...)
		if actual != jst.output {
			t.Errorf("JoinStringsSep(%v): expected %s, actual %s", jst.elem, jst.output, actual)
		}
	}
}

func TestGetMD5Hash(t *testing.T) {
	type md5Test struct {
		text   string
		output string
	}

	var md5Tests = []md5Test{
		{"Hello World", "b10a8db164e0754105b7a99be72e3fe5"},
		{"I am Groot!", "a3a25a9b537682816c35feac19dc5edc"},
		{"Luke, I am your Father!", "42d3611e3bb9ae93a647399f54ac1766"},
		{"Avengers, assemble.", "56518ca336d16939f703fe870f49be53"},
	}

	for _, mdt := range md5Tests {
		actual := GetMD5Hash(mdt.text)
		if actual != mdt.output {
			t.Errorf("GetMD5Hash(%s): expected %s, actual %s", mdt.text, mdt.output, actual)
		}
	}
}

func TestSliceContainsString(t *testing.T) {
	type sliceContainingTest struct {
		slice []string
		text  string
		found bool
		pos   int
	}

	var sliceContainingTests = []sliceContainingTest{
		{[]string{"Hello", "World"}, "Hello", true, 0},
		{[]string{"Hello", "World"}, "Bye", false, -1},
	}

	for _, sct := range sliceContainingTests {
		actPos, actFound := SliceContainsString(sct.slice, sct.text)
		if actFound != sct.found {
			t.Errorf("SliceContainsString(%v, %s): expected %v, actual %v", sct.slice, sct.text, sct.found, actFound)
		}
		if actPos != sct.pos {
			t.Errorf("SliceContainsString(%v, %s): expected %v, actual %v", sct.slice, sct.text, sct.found, actPos)
		}
	}
}

func TestStrBeforeSubstr(t *testing.T) {
	type strSubstrTest struct {
		str    string
		substr string
		output string
	}

	var strBeforeSubstrTests = []strSubstrTest{
		{"Hello World", "World", "Hello "},
		{"Hello", "Hell", ""},
		{"Hello", "world", ""},
	}

	for _, sbs := range strBeforeSubstrTests {
		actual := StrBeforeSubstr(sbs.str, sbs.substr)
		if actual != sbs.output {
			t.Errorf("StrBeforeSubstr(%s, %s): expected %s, actual %s", sbs.str, sbs.substr, sbs.output, actual)
		}
	}
}

func TestStrAfterSubstr(t *testing.T) {
	type strSubstrTest struct {
		str    string
		substr string
		output string
	}

	var strAfterSubstrTests = []strSubstrTest{
		{"Hello World", "Hello", " World"},
		{"Hello", "World", ""},
		{"Hello", "lo", ""},
	}

	for _, sas := range strAfterSubstrTests {
		actual := StrAfterSubstr(sas.str, sas.substr)
		if actual != sas.output {
			t.Errorf("StrAfterSubstr(%s, %s): expected %s, actual %s", sas.str, sas.substr, sas.output, actual)
		}
	}
}

func TestStrWithSubstrMaxLen(t *testing.T) {
	type strSubstrMaxLenTest struct {
		str    string
		substr string
		maxLen int
		addLen int
		output string
	}

	var strSubstrMaxLenTests = []strSubstrMaxLenTest{
		{"The quick brown fox jumps over a lazy dog.", "fox", 20, 6, "brown fox jumps"},
		{"The quick brown fox jumps over a lazy dog.", "fox", 50, 6, "The quick brown fox jumps over a lazy dog."},
		{"The quick brown fox jumps over a lazy dog.", "The", 15, 6, "The quick brown"},
		{"The quick brown fox jumps over a lazy dog.", "dog.", 16, 6, "over a lazy dog."},
	}

	for _, mxt := range strSubstrMaxLenTests {
		actual := StrWithSubstrMaxLen(mxt.str, mxt.substr, mxt.maxLen, mxt.addLen)
		if actual != mxt.output {
			t.Errorf("StrWithSubstrMaxLen(%s, %s, %d, %d): expected %s, actual %s",
				mxt.str, mxt.substr, mxt.maxLen, mxt.addLen, mxt.output, actual)
		}
	}
}

func TestStrWithSubstrMaxLenIdx(t *testing.T) {
	type strSubstrMaxIdxLenTest struct {
		str    string
		substr string
		idx    int
		maxLen int
		addLen int
		output string
	}

	var strSubstrMaxIdxLenTests = []strSubstrMaxIdxLenTest{
		{"The quick brown fox jumps over a lazy dog.", "fox", 10, 50, 6, "The quick brown fox jumps over a lazy dog."},
		{"The quick brown fox jumps over a lazy dog.", "The", 10, 15, 6, "The quick brown"},
		{"The quick brown fox jumps over a lazy dog.", "dog.", 10, 16, 6, "over a lazy dog."},
		{"The quick brown fox jumps over a lazy dog.", "fox", 12, 20, 6, "ick brfox fox j"},
	}

	for _, mxt := range strSubstrMaxIdxLenTests {
		actual := StrWithSubstrMaxLenIdx(mxt.str, mxt.substr, mxt.idx, mxt.maxLen, mxt.addLen)
		if actual != mxt.output {
			t.Errorf("StrWithSubstrMaxLenIdx(%s, %s, %d, %d, %d): expected %s, actual %s",
				mxt.str, mxt.substr, mxt.idx, mxt.maxLen, mxt.addLen, mxt.output, actual)
		}
	}
}

func TestSliceStringUnique(t *testing.T) {
	type uniqueTest struct {
		input         []string
		output        []string
		caseSensitive bool
	}

	uniqueTests := []uniqueTest{
		{[]string{}, []string{}, false},
		{nil, nil, true},
		{[]string{"Hello", "World"}, []string{"Hello", "World"}, true},
		{[]string{"Hello", "World", "My", "World"}, []string{"Hello", "World", "My"}, true},
		{[]string{"Hello", "World", "My", "world", "Your", "hello", "World"}, []string{"Hello", "World", "My", "world", "Your", "hello"}, true},
		{[]string{"Hello", "World", "My", "World", "Your", "Hello", "World"}, []string{"Hello", "World", "My", "Your"}, false},
	}

	for _, ut := range uniqueTests {
		actual := SliceStringUnique(ut.input, ut.caseSensitive)
		sort.Strings(actual)
		sort.Strings(ut.output)
		if !reflect.DeepEqual(ut.output, actual) {
			t.Errorf("expected %v, got %v", ut.output, actual)
		}
	}
}

func TestSliceStringEqual(t *testing.T) {
	type equalsTest struct {
		left   []string
		right  []string
		result bool
	}

	equalTests := []equalsTest{
		{[]string{}, []string{}, true},
		{nil, nil, true},
		{[]string{"Go", "is", "cool"}, []string{"Go", "is", "cool"}, true},
		{nil, []string{"Hello", "World"}, false},
		{[]string{"Hello", "World"}, []string{"Hello"}, false},
		{[]string{"Hello", "World"}, []string{"Hello", "Shmorld"}, false},
	}

	for _, ut := range equalTests {
		actual := SliceStringEqual(ut.left, ut.right)
		if !reflect.DeepEqual(ut.result, actual) {
			t.Errorf("expected %v, got %v", ut.result, actual)
		}
	}
}

func TestGetFullGenderText(t *testing.T) {
	type genderTest struct {
		str    string
		output string
	}

	var genderTests = []genderTest{
		{"m", "male"},
		{"f", "female"},
		{"t", "female"},
	}

	for _, sas := range genderTests {
		actual := GetFullGenderText(sas.str)
		if actual != sas.output {
			t.Errorf("GetFullGenderText(%s): expected %s, actual %s", sas.str, sas.output, actual)
		}
	}
}
