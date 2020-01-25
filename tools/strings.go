package tools

import (
    "crypto/md5"
    "encoding/hex"
    "os"
    "strings"
)

// JoinStrings returns a string composed from a given set of args.
func JoinStrings(elem ...string) string {
    return strings.Join(elem, "")
}

// JoinStringsSep returns a string composed from the second and other args, separated by the separator as first arg.
func JoinStringsSep(sep string, elem ...string) string {
    return strings.Join(elem, sep)
}

// SliceContainsString returns idx and true if a found in s. Otherwise -1 and false.
func SliceContainsString(s []string, a string) (int, bool) {
    for i, b := range s {
        if b == a {
            return i, true
        }
    }

    return -1, false
}

// SliceContainsAnyString returns true if any string of the given list is found
func SliceContainsAnyString(s, matches []string) bool {
    for _, m := range matches {
        _, f := SliceContainsString(s, m)
        if f {
            return true
        }
    }

    return false
}

// StrBeforeSubstr returns a string preceding a given substr.
func StrBeforeSubstr(str, substr string) string {
    p := strings.Index(str, substr)
    if p == -1 {
        return ""
    }
    return str[0:p]
}

// StrAfterSubstr returns a string following after a given substr.
func StrAfterSubstr(str, substr string) string {
    p := strings.Index(str, substr)
    if p == -1 {
        return ""
    }
    end := p + len(substr)
    if end >= len(str) {
        return ""
    }
    return str[end:]
}

// StrWithSubstrMaxLen returns a string containing a given substr.
// The returned string is no longer than maxLen.
//
// If the original str is longer than maxLen, the returned string is:
// - if str starts with substr => result is str up to maxLen
// - if str ends with substr => result is str from len-maxLen up to the end
// - if str contains substr in the middle => substr gets surrounded by addLen such that maxLen = addLen + len(substr) addLen
func StrWithSubstrMaxLen(str, substr string, maxLen, addLen int) string {
    if len(str) <= maxLen {
        return str
    }

    if strings.HasPrefix(str, substr) {
        return str[:maxLen]
    }

    if strings.HasSuffix(str, substr) {
        return str[len(str)-maxLen:]
    }

    before := StrBeforeSubstr(str, substr)
    after := StrAfterSubstr(str, substr)
    if len(before) > 0 && len(before) > addLen {
        before = before[len(before)-addLen:]
    }
    if len(after) > 0 && len(after) > addLen {
        after = after[:addLen]
    }

    return JoinStrings(before, substr, after)
}

// StrWithSubstrMaxLenIdx returns a string containing a given substr.
func StrWithSubstrMaxLenIdx(str, substr string, idx, maxLen, addLen int) string {
    if len(str) == 0 || len(substr) == 0 {
        return ""
    }

    // It should not happen, but we need to be sure.
    if len(str) < idx {
        return ""
    }

    strLen := len(str)
    substrLen := len(substr)

    if strLen <= maxLen {
        return str
    }

    if strings.HasPrefix(str, substr) {
        return str[:maxLen]
    }

    if strings.HasSuffix(str, substr) {
        return str[strLen-maxLen:]
    }

    before := str[0:idx]
    after := ""
    end := idx + substrLen
    if end < strLen {
        after = str[end:]
    }

    if len(before) > 0 && len(before) > addLen {
        before = before[len(before)-addLen:]
    }
    if len(after) > 0 && len(after) > addLen {
        after = after[:addLen]
    }

    return JoinStrings(before, substr, after)
}

// GetMD5Hash returns md5 hash of a string
func GetMD5Hash(text string) string {
    hasher := md5.New()
    _, err := hasher.Write([]byte(text))
    if err != nil {
        return ""
    }
    return hex.EncodeToString(hasher.Sum(nil))
}

// SliceStringUnique returns a slice of unique strings by discarding duplicates from the original.
func SliceStringUnique(original []string, caseSensitive bool) []string {
    if original == nil {
        return nil
    }

    unique := make([]string, 0)
    keys := make(map[string]struct{})
    for _, val := range original {
        keyToCheck := val
        if !caseSensitive {
            keyToCheck = strings.ToLower(val)
        }

        if _, ok := keys[keyToCheck]; !ok {
            keys[keyToCheck] = struct{}{}
            unique = append(unique, val)
        }
    }

    return unique
}

// SliceStringEqual performs a simple check for equality for slices of strings.
func SliceStringEqual(a, b []string) bool {
    if a == nil && b == nil {
        return true
    }

    if a == nil || b == nil {
        return false
    }

    if len(a) != len(b) {
        return false
    }

    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }

    return true
}

// StringBetweenDelimiters will return the text between the given delimiters
func StringBetweenDelimiters(s, leftDelimiter, rightDelimiter string) string {
    i := strings.Index(s, leftDelimiter)
    if i >= 0 {
        j := strings.Index(s[i+1:], rightDelimiter)
        if j >= 0 && i+j+1 <= len(s) {
            return s[i+len(leftDelimiter) : i+j+1]
        }
    }

    return ""
}

// GetCurrentDirectory function
func GetCurrentDirectory() string {
    pwd, err := os.Getwd()
    if err != nil {
        return ""
    }

    return pwd
}

// GetFullGenderText returns full gender text
func GetFullGenderText(g string) string {
    gender := "female"

    if strings.ToLower(g) == "m" {
        gender = "male"
    }

    return gender
}
