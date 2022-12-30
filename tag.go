package envtags

import "strings"

func parseTagValue(tagValue string) tagDetails {
	options := strings.Split(tagValue, ",")
	result := tagDetails{EnvVarName: options[0]}

	for _, opt := range options[1:] {
		optSlice := strings.Split(opt, ":")
		switch k, v := optSlice[0], optSlice[1]; k {
		case "prefix":
			result.Prefix = v
		}

	}

	return result
}

type tagDetails struct {
	Prefix     string
	EnvVarName string
}

func (td tagDetails) Update(oldDetails tagDetails) tagDetails {
	td.Prefix = oldDetails.Prefix + td.Prefix
	return td
}

func (td tagDetails) GetEnvVar() string {
	return td.Prefix + td.EnvVarName
}
