package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/bitcrshr/envmgr/api/ent/schema/gotype"
	"github.com/google/uuid"
)

// Environment holds the schema definition for the Environment entity.
type Environment struct {
	ent.Schema
}

// Fields of the Environment.
func (Environment) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Unique().Immutable(),
		field.String("name").NotEmpty(),
		field.Enum("kind").GoType(gotype.EnvironmentKind("")),
	}
}

// Edges of the Environment.
func (Environment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("project", Project.Type).Ref("environments").Unique(),
		edge.To("vars", Variable.Type),
	}
}
