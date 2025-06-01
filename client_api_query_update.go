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
	result := ""
	if it.Deletes != nil {
		result += "    Deletes:\n" + it.Deletes.String()
	}
	if it.Inserts != nil {
		result += "    Inserts:\n" + it.Inserts.String()
	}
	return result
}
