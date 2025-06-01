package module_bindings

import (
	"fmt"

	"github.com/alexanderbh/spacetimedb-go-sdk"
)

type User struct {
	Identity *spacetimedb.Identity
	Name     *string
	Online   bool
}

type UserTable struct {
	TableCache *spacetimedb.TableCache[*User]
}

func NewUserTable() *UserTable {
	tableCache := &spacetimedb.TableCache[*User]{
		Rows: make(map[string]*User),
	}
	return &UserTable{
		TableCache: tableCache,
	}
}

func (*UserTable) DeserializeRow(reader *spacetimedb.BinaryReader) (any, error) {
	u := &User{}
	u.Identity = &spacetimedb.Identity{}
	if err := u.Identity.Deserialize(reader); err != nil {
		return nil, fmt.Errorf("failed to deserialize User.Identity: %w", err)
	}
	if reader.ReadU8() == 0 {
		name := reader.ReadString()
		u.Name = &name
	} else {
		u.Name = nil
	}
	u.Online = reader.ReadBool()

	return u, nil
}

func (t *UserTable) Insert(rowAny any) error {
	row := rowAny.(*User)
	if row.Identity == nil {
		return fmt.Errorf("User.Identity cannot be nil")
	}
	if row.Name != nil && *row.Name == "" {
		return fmt.Errorf("User.Name cannot be an empty string")
	}

	t.TableCache.Rows[row.Identity.ToHexString()] = row
	return nil
}
