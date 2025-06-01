package spacetimedb

type QueryUpdate struct {
	Deletes *BsatnRowList
	Inserts *BsatnRowList
}

func (it *QueryUpdate) Deserialize(reader *BinaryReader) error {

	it.Deletes = &BsatnRowList{}
	it.Deletes.Deserialize(reader)

	it.Inserts = &BsatnRowList{}
	it.Inserts.Deserialize(reader)

	return nil
}

func (it *QueryUpdate) String() string {
	result := "QueryUpdate:\n"
	if it.Deletes != nil {
		result += "  Deletes: " + it.Deletes.String() + "\n"
	} else {
		result += "  Deletes: <nil>\n"
	}
	if it.Inserts != nil {
		result += "  Inserts: " + it.Inserts.String() + "\n"
	} else {
		result += "  Inserts: <nil>\n"
	}
	return result
}
