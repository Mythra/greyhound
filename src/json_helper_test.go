package main

import (
	"bytes"
	"reflect"
	"sort"
	"testing"
)

func TestByNameSort(t *testing.T) {
	fieldUno := field{"uno", []byte("uno"), bytes.EqualFold, true, []int{1}, reflect.TypeOf("uno"), true, true}
	fieldDuo := field{"duo", []byte("duo"), bytes.EqualFold, true, []int{2}, reflect.TypeOf("duo"), true, true}
	fieldTres := field{"tres", []byte("tres"), bytes.EqualFold, true, []int{3}, reflect.TypeOf("tres"), true, true}

	result := byName([]field{fieldUno, fieldDuo, fieldTres})
	sort.Sort(result)
	if result[0].name != "duo" || result[1].name != "tres" || result[2].name != "uno" {
		t.Fatalf("Sorted Values are not in correct sorted order: \n %+v", result)
	}
}

func TestByIndexSort(t *testing.T) {
	fieldUno := field{"uno", []byte("uno"), bytes.EqualFold, true, []int{1}, reflect.TypeOf("uno"), true, true}
	fieldDuo := field{"duo", []byte("duo"), bytes.EqualFold, true, []int{2}, reflect.TypeOf("duo"), true, true}
	fieldTres := field{"tres", []byte("tres"), bytes.EqualFold, true, []int{3}, reflect.TypeOf("tres"), true, true}

	result := byIndex([]field{fieldTres, fieldUno, fieldDuo})
	sort.Sort(result)
	if result[0].name != "uno" || result[1].name != "duo" || result[2].name != "tres" {
		t.Fatalf("Sorted Values are not in correct sorted order: \n %+v", result)
	}
}

func TestTagOptionContains(t *testing.T) {
	tagOptions := tagOptions("hello,world,my,name,is,kanye,i,am,your,savior")

	if !tagOptions.Contains("hello") || !tagOptions.Contains("savior") {
		t.Fatalf("Error contains: \n [ %+v ] \n %+v \n %+v ", tagOptions, tagOptions.Contains("hello"), tagOptions.Contains("savior"))
	}
}

func TestParseTag(t *testing.T) {
	parsedTag, options := parseTag("myName,omitempty,lel")

	if parsedTag != "myName" || !options.Contains("omitempty") || !options.Contains("lel") {
		t.Fatalf("Error in ParseTag: \n %s \n [ %+v ] \n", parsedTag, options)
	}
}

func TestValidateTag(t *testing.T) {
	shouldBeValidTag := isValidTag("myName")
	shouldBeInvalidTag := isValidTag("これは大丈夫です,,,,,,")

	if !shouldBeValidTag || shouldBeInvalidTag {
		t.Fatalf(
			"Valid Tags returned incorrect status: \n [ %s ] \n [ %s ] \n [ %+v ] \n [ %+v ]",
			"myName,omitempty",
			"これは大丈夫です,,,,,,",
			shouldBeValidTag,
			shouldBeInvalidTag,
		)
	}
}
