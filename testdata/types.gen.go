// Code generated by sql_update_by_struct DO NOT EDIT.

package testdata

import "strconv"

func (p *PersonUpdate) Update(args []any) ([]string, []any) {
	cols := make([]string, 0, 4)
	if p.Name != nil {
		args = append(args, *p.Name)
		cols = append(cols, "uname = $"+strconv.Itoa(len(args)))
	}
	if p.Age != nil {
		args = append(args, *p.Age)
		cols = append(cols, "age = $"+strconv.Itoa(len(args)))
	}
	if p.Email != nil {
		args = append(args, *p.Email)
		cols = append(cols, "email = $"+strconv.Itoa(len(args)))
	}
	if p.Times1 != nil {
		_cols, _args := p.Times1.Update(args)
		args = append(args, _args...)
		cols = append(cols, _cols...)
	}
	return cols, args
}

func (t *Times) Update(args []any) ([]string, []any) {
	cols := make([]string, 0, 3)
	if t.CreatedAt != nil {
		args = append(args, *t.CreatedAt)
		cols = append(cols, "createdat = $"+strconv.Itoa(len(args)))
	}
	if t.UpdatedAt != nil {
		args = append(args, *t.UpdatedAt)
		cols = append(cols, "updatedat = $"+strconv.Itoa(len(args)))
	}
	if t.DeletedAt != nil {
		args = append(args, *t.DeletedAt)
		cols = append(cols, "deletedat = $"+strconv.Itoa(len(args)))
	}
	return cols, args
}
