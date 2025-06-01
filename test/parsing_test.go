package test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/alexanderbh/spacetimedb-go-sdk"
)

func TestParsingIdentityTokenMessage(t *testing.T) {
	// The first byte is removed from this as it indicates the CompressionType. Which is assumed is handled here.

	// This is actual response from the server (but without the first compression byte).
	hexStr := "030cf431a3f5d27800aaec59ec8ec09883b2f9a9a9b41480056a0487210af600c28201000065794a30655841694f694a4b563151694c434a68624763694f694a46557a49314e694a392e65794a6f5a5868666157526c626e527064486b694f694a6a4d6a41775a6a5977595449784f4463774e445a684d4455344d444530596a52684f5745355a6a6c694d6a677a4f54686a4d44686c5a574d314f57566a595745774d4463345a444a6d4e57457a4d7a466d4e44426a4969776963335669496a6f694d6d56684d5449305a4441745a444935597930305a4759774c546b774f4451745a4745324d6a4e6a4e3259304e5756694969776961584e7a496a6f696247396a5957786f62334e30496977695958566b496a7062496e4e7759574e6c64476c745a575269496c3073496d6c68644349364d5463304f44417a4e4441324f4377695a586877496a70756457787366512e5a4a4b6e734f633854736e547033584a576473535a64755379315876366f2d55585f494672334a4f7a6e7072723970424f68497549426e6b4c5a5f6c6f34777351306863523667624476446672444857686a63664167daa835213b727f77283d0ff5d51884c1"

	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		t.Fatalf("failed to decode hex string: %v", err)
	}

	reader := spacetimedb.NewBinaryReader(bytes)

	got := &spacetimedb.ServerMessage{}
	if err := got.Deserialize(reader); err != nil {
		t.Fatalf("failed to deserialize server message: %v", err)
	}
	fmt.Println("Deserialized ServerMessage:", got)

	identityToken, ok := got.Message.(*spacetimedb.IdentityToken)

	if !ok {
		t.Errorf("failed to cast to spacetimedb.IdentityToken")
	}

	idBigInt, err := spacetimedb.HexStringToU256("c200f60a2187046a058014b4a9a9f9b28398c08eec59ecaa0078d2f5a331f40c")
	wantIdentity, err := spacetimedb.NewIdentity(idBigInt)
	if identityToken.Identity.IsEqual(wantIdentity) == false {
		t.Errorf("identity mismatch: got %s, want %s", identityToken.Identity.ToHexString(), wantIdentity.ToHexString())
	}

	// CONNECTION ID
	wantConnectionId := spacetimedb.NewConnectionId(big.NewInt(456))
	if identityToken.ConnectionId.IsEqual(wantConnectionId) {
		t.Errorf("connectionId mismatch: got %v, want %v", identityToken.ConnectionId, wantConnectionId)
	}

	// TOKEN
	wantToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJFUzI1NiJ9.eyJoZXhfaWRlbnRpdHkiOiJjMjAwZjYwYTIxODcwNDZhMDU4MDE0YjRhOWE5ZjliMjgzOThjMDhlZWM1OWVjYWEwMDc4ZDJmNWEzMzFmNDBjIiwic3ViIjoiMmVhMTI0ZDAtZDI5Yy00ZGYwLTkwODQtZGE2MjNjN2Y0NWViIiwiaXNzIjoibG9jYWxob3N0IiwiYXVkIjpbInNwYWNldGltZWRiIl0sImlhdCI6MTc0ODAzNDA2OCwiZXhwIjpudWxsfQ.ZJKnsOc8TsnTp3XJWdsSZduSy1Xv6o-UX_IFr3JOznprr9pBOhIuIBnkLZ_lo4wsQ0hcR6gbDvDfrDHWhjcfAg"
	if identityToken.Token != wantToken {
		t.Errorf("token mismatch: got %s, want %s", identityToken.Token, wantToken)
	}
}
