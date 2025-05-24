package spacetimedb

import "math/big"

// InitialSubscriptionMessage matches the TS type for initial subscription messages.
type InitialSubscriptionMessage struct {
	Tag          string
	TableUpdates []TableUpdate
}

// TableUpdate matches the TS type imported from table_cache.ts.
type TableUpdate struct {
	TableName  string
	Operations []Operation
}

type Operation struct {
	Type  string
	RowId string
	Row   interface{}
}

// TransactionUpdateMessage matches the TS type for transaction update messages.
type TransactionUpdateMessage struct {
	Tag            string // 'TransactionUpdate'
	TableUpdates   []TableUpdate
	Identity       Identity
	ConnectionId   *ConnectionId
	ReducerInfo    *ReducerInfo
	Status         int // TODO: Make UpdateStatus
	Message        string
	Timestamp      int // TODO: make: Timestamp
	EnergyConsumed big.Int
}

type ReducerInfo struct {
	ReducerName string
	Args        []byte
}

// TransactionUpdateLightMessage matches the TS type for light transaction update messages.
type TransactionUpdateLightMessage struct {
	Tag          string // 'TransactionUpdateLight'
	TableUpdates []TableUpdate
}

// IdentityTokenMessage matches the TS type for identity token messages.
type IdentityTokenMessage struct {
	Tag          string // 'IdentityToken'
	Identity     Identity
	Token        string
	ConnectionId ConnectionId
}

// SubscribeAppliedMessage matches the TS type for subscribe applied messages.
type SubscribeAppliedMessage struct {
	Tag          string // 'SubscribeApplied'
	QueryId      int
	TableUpdates []TableUpdate
}

// UnsubscribeAppliedMessage matches the TS type for unsubscribe applied messages.
type UnsubscribeAppliedMessage struct {
	Tag          string // 'UnsubscribeApplied'
	QueryId      int
	TableUpdates []TableUpdate
}

// SubscriptionError matches the TS type for subscription errors.
type SubscriptionError struct {
	Tag     string // 'SubscriptionError'
	QueryId *int   // optional
	Error   string
}

// Message is a sum type for all message variants.
type Message interface{}
