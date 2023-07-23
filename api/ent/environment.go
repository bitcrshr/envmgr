// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/bitcrshr/envmgr/api/ent/environment"
	"github.com/bitcrshr/envmgr/api/ent/project"
	"github.com/bitcrshr/envmgr/api/ent/schema/gotype"
	"github.com/google/uuid"
)

// Environment is the model entity for the Environment schema.
type Environment struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Kind holds the value of the "kind" field.
	Kind gotype.EnvironmentKind `json:"kind,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the EnvironmentQuery when eager-loading is set.
	Edges                EnvironmentEdges `json:"edges"`
	project_environments *uuid.UUID
	selectValues         sql.SelectValues
}

// EnvironmentEdges holds the relations/edges for other nodes in the graph.
type EnvironmentEdges struct {
	// Project holds the value of the project edge.
	Project *Project `json:"project,omitempty"`
	// Vars holds the value of the vars edge.
	Vars []*Variable `json:"vars,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// ProjectOrErr returns the Project value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e EnvironmentEdges) ProjectOrErr() (*Project, error) {
	if e.loadedTypes[0] {
		if e.Project == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: project.Label}
		}
		return e.Project, nil
	}
	return nil, &NotLoadedError{edge: "project"}
}

// VarsOrErr returns the Vars value or an error if the edge
// was not loaded in eager-loading.
func (e EnvironmentEdges) VarsOrErr() ([]*Variable, error) {
	if e.loadedTypes[1] {
		return e.Vars, nil
	}
	return nil, &NotLoadedError{edge: "vars"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Environment) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case environment.FieldName, environment.FieldKind:
			values[i] = new(sql.NullString)
		case environment.FieldID:
			values[i] = new(uuid.UUID)
		case environment.ForeignKeys[0]: // project_environments
			values[i] = &sql.NullScanner{S: new(uuid.UUID)}
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Environment fields.
func (e *Environment) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case environment.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				e.ID = *value
			}
		case environment.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				e.Name = value.String
			}
		case environment.FieldKind:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field kind", values[i])
			} else if value.Valid {
				e.Kind = gotype.EnvironmentKind(value.String)
			}
		case environment.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullScanner); !ok {
				return fmt.Errorf("unexpected type %T for field project_environments", values[i])
			} else if value.Valid {
				e.project_environments = new(uuid.UUID)
				*e.project_environments = *value.S.(*uuid.UUID)
			}
		default:
			e.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Environment.
// This includes values selected through modifiers, order, etc.
func (e *Environment) Value(name string) (ent.Value, error) {
	return e.selectValues.Get(name)
}

// QueryProject queries the "project" edge of the Environment entity.
func (e *Environment) QueryProject() *ProjectQuery {
	return NewEnvironmentClient(e.config).QueryProject(e)
}

// QueryVars queries the "vars" edge of the Environment entity.
func (e *Environment) QueryVars() *VariableQuery {
	return NewEnvironmentClient(e.config).QueryVars(e)
}

// Update returns a builder for updating this Environment.
// Note that you need to call Environment.Unwrap() before calling this method if this Environment
// was returned from a transaction, and the transaction was committed or rolled back.
func (e *Environment) Update() *EnvironmentUpdateOne {
	return NewEnvironmentClient(e.config).UpdateOne(e)
}

// Unwrap unwraps the Environment entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (e *Environment) Unwrap() *Environment {
	_tx, ok := e.config.driver.(*txDriver)
	if !ok {
		panic("ent: Environment is not a transactional entity")
	}
	e.config.driver = _tx.drv
	return e
}

// String implements the fmt.Stringer.
func (e *Environment) String() string {
	var builder strings.Builder
	builder.WriteString("Environment(")
	builder.WriteString(fmt.Sprintf("id=%v, ", e.ID))
	builder.WriteString("name=")
	builder.WriteString(e.Name)
	builder.WriteString(", ")
	builder.WriteString("kind=")
	builder.WriteString(fmt.Sprintf("%v", e.Kind))
	builder.WriteByte(')')
	return builder.String()
}

// Environments is a parsable slice of Environment.
type Environments []*Environment