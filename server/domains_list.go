package main

import "strings"

var longTLDs = [...]string{
	"co.jp",
	"co.uk",
	"or.us",
	"k12.or.us",
}

// hasLongTLD indicates whether the domain name has a long TLD, such as ".co.uk".
func hasLongTLD(domain string) bool {
	for _, tld := range longTLDs {
		if strings.HasSuffix(domain, "."+tld) {
			return true
		}
	}
	return false
}

// prepareDomainNames adds into the domains slice the wildcard sub-domain variant for each domain name.
// If the output should include "example.com" as well as "*.example.com" then the given domains slice
// must include both "example.com" and "*.example.com".
func prepareDomainNames(domains []string) []string {
	made := make([]string, 0, len(domains)*2)
	// remove will contain any domain names that will need to be removed after this initial processing.
	var remove []string

	for _, d := range domains {
		d = strings.TrimSpace(d)
		if d == "" {
			continue
		}

		if stringSliceContains(made, d) {
			continue
		}

		made = append(made, d)

		// If this domain has two dots in the TLD extension and it's a sub-domain, or if this domain
		// has a single part in the TLD and it's a sub domain, then don't include any other variant.
		// The sub-domain could just be "www", in which case no other variant is retrieved.
		longTLD := hasLongTLD(d)
		dots := strings.Count(d, ".")
		if (longTLD && dots > 2) ||
			(!longTLD && dots > 1) {
			// If there is a wildcard sub-domain, simply include this and make sure the "www" sub-domain
			// variant is not also included in the output slice.
			// We cannot include any particular sub-domains if at that sub-level we want a wildcard.
			// So you cannot request a certificate for www.example.com and *.example.com.
			if d[0] == '*' {
				remove = append(remove, d[1:])
			}
			continue
		}

		// Include the "www" sub-domain variant because d is not a sub-domain.
		www := "www." + d
		if !stringSliceContains(made, www) {
			made = append(made, www)
		}
	}

	// Pruning must be done here because we can't know the order in which the domain names in the initial
	// domains slice will be given.
	for _, d := range remove {
		made = removeStringsSuffixed(made, d)
	}
	return made
}

// stringSliceContains indicates whether the slice contains the string str.
func stringSliceContains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// removeStringsSuffixed returns a slice with the given slice's strings but without any strings that
// have the suffix but are not the string "*"+suffix.
func removeStringsSuffixed(slice []string, suffix string) []string {
	for i := 0; i < len(slice); i++ {
		if strings.HasSuffix(slice[i], suffix) && slice[i] != "*"+suffix {
			slice = append(slice[:i], slice[i+1:]...)
			i--
		}
	}
	return slice
}
