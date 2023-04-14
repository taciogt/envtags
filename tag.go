package envtags

import "strings"

func parseTagValue(tagValue string) tagDetails {
	options := strings.Split(tagValue, ",")
	result := tagDetails{EnvVarName: options[0]}

	for _, opt := range options[1:] {
		optSlice := strings.Split(opt, ":")
		switch k := optSlice[0]; k {
		case "prefix":
			v := optSlice[1]
			result.Prefix = v
		case "rune":
			result.IsRune = true
		}

	}

	return result
}

type tagDetails struct {
	Prefix     string
	EnvVarName string
	IsRune     bool
}

func (td tagDetails) Update(oldDetails tagDetails) tagDetails {
	td.Prefix = oldDetails.Prefix + td.Prefix
	return td
}

func (td tagDetails) GetEnvVar() string {
	return td.Prefix + td.EnvVarName
}
