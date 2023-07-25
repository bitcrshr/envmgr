package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Immutable().Unique().Default(uuid.New),
		field.String("auth0_id").NotEmpty(),
		field.String("email").Optional().Nillable(),
		field.String("family_name").Optional().Nillable(),
		field.String("given_name").Optional().Nillable(),
		field.String("name").Optional().Nillable(),
		field.String("nickname").Optional().Nillable(),
		field.String("phone_number").Optional().Nillable().Sensitive(),
		field.String("picture").Optional().Nillable(),
		field.String("username").Optional().Nillable(),
		field.JSON("app_metadata", map[string]string{}).Optional().Sensitive().Default(map[string]string{}),
		field.JSON("user_metadata", map[string]string{}).Optional().Sensitive().Default(map[string]string{}),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("projects", Project.Type),
	}
}
