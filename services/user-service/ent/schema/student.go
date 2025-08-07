package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Student struct {
	ent.Schema
}

func (Student) Fields() []ent.Field {
	return []ent.Field{
		field.String("student_code").NotEmpty().Unique().Comment("ma hoc sinh"),
		field.String("name").NotEmpty().Comment("ten hoc sinh"),
		field.String("email").NotEmpty().Unique().Comment("email hoc sinh"),
		field.String("phone").NotEmpty().Unique().Comment("so dien thoai hoc sinh"),
		field.Time("date_of_birth").Comment("ngay sinh hoc sinh"),
		field.String("address").NotEmpty().Comment("dia chi hoc sinh"),
		field.Time("created_at").Default(time.Now).Comment("thoi gian tao"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("thoi gian cap nhat"),
		field.Time("deleted_at").Optional().Nillable().Comment("thoi gian xoa"),
		field.Int("user_id").Optional().Comment("id cua user quan ly hoc sinh"),
	}
}

func (Student) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("students").Unique().Field("user_id").Comment("moi quan he voi nguoi dung"),
	}
}

// Indexes of the Student.
func (Student) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("student_code").Unique(),
		index.Fields("email").Unique(),
		index.Fields("user_id"), // Index cho foreign key
		index.Fields("phone"),
		index.Fields("created_at"),
		index.Fields("user_id", "created_at"), // Compound index
	}
}
