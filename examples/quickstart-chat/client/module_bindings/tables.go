package module_bindings

import "github.com/alexanderbh/spacetimedb-go-sdk"

var Tables = map[string]spacetimedb.Table{
	"user": NewUserTable(),
}
