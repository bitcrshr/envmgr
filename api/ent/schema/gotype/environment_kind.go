package gotype

type EnvironmentKind string

const (
	EnvironmentKind_Unspecified EnvironmentKind = "UNSPECIFIED"
	EnvironmentKind_Development EnvironmentKind = "DEVELOPMENT"
	EnvironmentKind_Staging     EnvironmentKind = "STAGING"
	EnvironmentKind_Production  EnvironmentKind = "PRODUCTION"
)

func (EnvironmentKind) Values() []string {
	return []string{
		string(EnvironmentKind_Unspecified),
		string(EnvironmentKind_Development),
		string(EnvironmentKind_Staging),
		string(EnvironmentKind_Production),
	}
}
