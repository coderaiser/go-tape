//go:build !no_external

package diff

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

const (
	red   = "\033[31m"
	green = "\033[32m"
	reset = "\033[0m"
)

// Diff returns a colorized line-diff between expected and result,
// matching jest-diff output style: red "- expected", green "+ received".
// Returns "" when values are deeply equal.
func Diff(expected, result any) string {
	if reflect.DeepEqual(expected, result) {
		return ""
	}
	aLines := strings.Split(prettyValue(reflect.ValueOf(expected), 0), "\n")
	bLines := strings.Split(prettyValue(reflect.ValueOf(result), 0), "\n")
	edits := buildEdits(aLines, bLines)

	var sb strings.Builder
	sb.WriteString(red + "- expected" + reset + "\n")
	sb.WriteString(green + "+ received" + reset + "\n\n")
	for _, e := range edits {
		switch e.kind {
		case opEqual:
			sb.WriteString("  " + e.line + "\n")
		case opDelete:
			sb.WriteString(red + "- " + e.line + reset + "\n")
		case opInsert:
			sb.WriteString(green + "+ " + e.line + reset + "\n")
		}
	}
	return sb.String()
}

// prettyValue renders v as indented, line-per-field text for diffing.
func prettyValue(v reflect.Value, depth int) string {
	if !v.IsValid() {
		return "nil"
	}
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return fmt.Sprintf("(%s)(nil)", v.Type())
		}
		return "&" + prettyValue(v.Elem(), depth)
	}
	if v.Kind() == reflect.Interface {
		if v.IsNil() {
			return "nil"
		}
		return prettyValue(v.Elem(), depth)
	}

	ind := strings.Repeat("  ", depth)
	ind1 := strings.Repeat("  ", depth+1)

	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		if v.NumField() == 0 {
			return t.String() + "{}"
		}
		var sb strings.Builder
		sb.WriteString(t.String() + "{\n")
		for i := 0; i < v.NumField(); i++ {
			f := t.Field(i)
			sb.WriteString(ind1 + f.Name + ": " + prettyValue(v.Field(i), depth+1) + ",\n")
		}
		sb.WriteString(ind + "}")
		return sb.String()

	case reflect.Map:
		if v.IsNil() || v.Len() == 0 {
			return v.Type().String() + "{}"
		}
		keys := v.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			return fmt.Sprintf("%v", keys[i]) < fmt.Sprintf("%v", keys[j])
		})
		var sb strings.Builder
		sb.WriteString(v.Type().String() + "{\n")
		for _, k := range keys {
			sb.WriteString(ind1 + prettyValue(k, depth+1) + ": " + prettyValue(v.MapIndex(k), depth+1) + ",\n")
		}
		sb.WriteString(ind + "}")
		return sb.String()

	case reflect.Slice, reflect.Array:
		if v.Kind() == reflect.Slice && v.IsNil() {
			return v.Type().String() + "(nil)"
		}
		if v.Len() == 0 {
			return v.Type().String() + "{}"
		}
		var sb strings.Builder
		sb.WriteString(v.Type().String() + "{\n")
		for i := 0; i < v.Len(); i++ {
			sb.WriteString(ind1 + prettyValue(v.Index(i), depth+1) + ",\n")
		}
		sb.WriteString(ind + "}")
		return sb.String()

	case reflect.String:
		return fmt.Sprintf("%q", v.String())

	case reflect.Invalid, reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr, reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128:
		return fmt.Sprintf("%v", v.Interface())

	case reflect.Chan, reflect.Func,
		reflect.Interface, reflect.Pointer, reflect.UnsafePointer:
		panic(fmt.Sprintf("unhandled reflect.Kind in prettyValue: %v", v.Kind()))
	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}

// Myers diff via LCS.
type opKind int

const (
	opEqual  opKind = iota
	opDelete        // in expected, not in result  → red -
	opInsert        // in result, not in expected  → green +
)

type edit struct {
	kind opKind
	line string
}

func lcs(a, b []string) [][]int {
	m, n := len(a), len(b)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}
	return dp
}

func buildEdits(a, b []string) []edit {
	dp := lcs(a, b)
	var edits []edit
	i, j := len(a), len(b)
	for i > 0 || j > 0 {
		switch {
		case i > 0 && j > 0 && a[i-1] == b[j-1]:
			edits = append(edits, edit{opEqual, a[i-1]})
			i--
			j--

		case j > 0 && (i == 0 || dp[i][j-1] >= dp[i-1][j]):
			edits = append(edits, edit{opInsert, b[j-1]})
			j--

		default:
			edits = append(edits, edit{opDelete, a[i-1]})
			i--
		}
	}
	for l, r := 0, len(edits)-1; l < r; l, r = l+1, r-1 {
		edits[l], edits[r] = edits[r], edits[l]
	}
	return edits
}
