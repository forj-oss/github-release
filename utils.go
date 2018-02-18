package main

import "unicode"

func ReleaseStatus(draft, pre_release bool) (release_status string) {
	release_status = "release"

	if draft {
		release_status = "draft release"
		return
	}
	if pre_release {
		release_status = "pre-release"
	}
	return
}

func Capitalize(input string) string {
	output := []rune(input)
	output[0] = unicode.ToTitle(output[0])

	return string(output)
}
