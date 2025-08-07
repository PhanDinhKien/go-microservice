package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			MinLen(2).
			MaxLen(100).
			Comment("User's full name"),
		field.String("email").
			Unique().
			NotEmpty().
			MaxLen(255).
			Comment("User's email address"),
		field.String("phone").
			Optional().
			MaxLen(20).
			Comment("User's phone number"),
		field.Enum("status").
			Values("active", "inactive", "suspended").
			Default("active").
			Comment("User's account status"),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("User creation timestamp"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("User last update timestamp"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("students", Student.Type).Comment("danh sach hoc sinh cua user"),
	}
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("email").Unique(),
		index.Fields("status"),
		index.Fields("created_at"),
	}
}
