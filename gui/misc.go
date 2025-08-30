package gui

import "strings"

/* Remove string r from slice s (assuming r occurs at most once) */
func remove(s []string, r string) []string {
	for i, v := range s {
		if strings.TrimSpace(v) == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

/* Remove all strings in slice r from slice s (assuming at most 1 present of every string found in r) */
func removeMany(s []string, r []string) []string {
	for _, u := range r {
		for i, v := range s {
			if strings.TrimSpace(v) == u {
				s = append(s[:i], s[i+1:]...)
				continue
			}
		}
	}
	return s
}
