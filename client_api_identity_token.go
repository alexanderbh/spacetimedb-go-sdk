package spacetimedb

type IdentityToken struct {
	Identity     *Identity     `json:"identity"`
	Token        string        `json:"token"`
	ConnectionId *ConnectionId `json:"connectionId"`
}

func (it *IdentityToken) Deserialize(reader *BinaryReader) error {

	it.Identity = &Identity{}
	it.Identity.Deserialize(reader)

	it.Token = reader.ReadString()

	it.ConnectionId = &ConnectionId{}
	it.ConnectionId.Deserialize(reader)

	return nil
}
