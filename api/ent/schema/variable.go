package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Variable holds the schema definition for the Variable entity.
type Variable struct {
	ent.Schema
}

// Fields of the Variable.
func (Variable) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Immutable().Unique(),
		field.String("key").NotEmpty(),
		field.String("value").Optional().Default("").Sensitive(),
	}
}

// Edges of the Variable.
func (Variable) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("environment", Environment.Type).Ref("vars").Unique(),
	}
}
