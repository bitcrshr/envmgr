// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// EnvironmentsColumns holds the columns for the "environments" table.
	EnvironmentsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "name", Type: field.TypeString},
		{Name: "kind", Type: field.TypeEnum, Enums: []string{"UNSPECIFIED", "DEVELOPMENT", "STAGING", "PRODUCTION"}},
		{Name: "project_environments", Type: field.TypeUUID, Nullable: true},
	}
	// EnvironmentsTable holds the schema information for the "environments" table.
	EnvironmentsTable = &schema.Table{
		Name:       "environments",
		Columns:    EnvironmentsColumns,
		PrimaryKey: []*schema.Column{EnvironmentsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "environments_projects_environments",
				Columns:    []*schema.Column{EnvironmentsColumns[3]},
				RefColumns: []*schema.Column{ProjectsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// ProjectsColumns holds the columns for the "projects" table.
	ProjectsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "display_name", Type: field.TypeString},
	}
	// ProjectsTable holds the schema information for the "projects" table.
	ProjectsTable = &schema.Table{
		Name:       "projects",
		Columns:    ProjectsColumns,
		PrimaryKey: []*schema.Column{ProjectsColumns[0]},
	}
	// VariablesColumns holds the columns for the "variables" table.
	VariablesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "key", Type: field.TypeString},
		{Name: "value", Type: field.TypeString, Nullable: true, Default: ""},
		{Name: "environment_vars", Type: field.TypeUUID, Nullable: true},
	}
	// VariablesTable holds the schema information for the "variables" table.
	VariablesTable = &schema.Table{
		Name:       "variables",
		Columns:    VariablesColumns,
		PrimaryKey: []*schema.Column{VariablesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "variables_environments_vars",
				Columns:    []*schema.Column{VariablesColumns[3]},
				RefColumns: []*schema.Column{EnvironmentsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		EnvironmentsTable,
		ProjectsTable,
		VariablesTable,
	}
)

func init() {
	EnvironmentsTable.ForeignKeys[0].RefTable = ProjectsTable
	VariablesTable.ForeignKeys[0].RefTable = EnvironmentsTable
}
