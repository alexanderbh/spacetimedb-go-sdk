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

func (u *User) PrimaryKey() string {
	return u.Identity.ToHexString()
}

func (u *User) Deserialize(reader *spacetimedb.BinaryReader) error {
	u.Identity = &spacetimedb.Identity{}
	if err := u.Identity.Deserialize(reader); err != nil {
		return fmt.Errorf("failed to deserialize User.Identity: %w", err)
	}
	if reader.ReadU8() == 0 {
		name := reader.ReadString()
		u.Name = &name
	} else {
		u.Name = nil
	}
	u.Online = reader.ReadBool()

	return nil
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

func (t *UserTable) Insert(reader *spacetimedb.BinaryReader) error {
	row := &User{}
	if err := row.Deserialize(reader); err != nil {
		return fmt.Errorf("failed to deserialize User row for insert: %w", err)
	}

	t.TableCache.Rows[row.PrimaryKey()] = row
	return nil
}

func (t *UserTable) Delete(reader *spacetimedb.BinaryReader) error {
	row := &User{}
	if err := row.Deserialize(reader); err != nil {
		return fmt.Errorf("failed to deserialize User row for delete: %w", err)
	}

	delete(t.TableCache.Rows, row.PrimaryKey())

	return nil
}
