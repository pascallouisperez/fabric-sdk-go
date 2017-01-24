// Code generated by protoc-gen-go.
// source: peer/transaction.proto
// DO NOT EDIT!

package peer

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf1 "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type InvalidTransaction_Cause int32

const (
	InvalidTransaction_TxIdAlreadyExists      InvalidTransaction_Cause = 0
	InvalidTransaction_RWConflictDuringCommit InvalidTransaction_Cause = 1
)

var InvalidTransaction_Cause_name = map[int32]string{
	0: "TxIdAlreadyExists",
	1: "RWConflictDuringCommit",
}
var InvalidTransaction_Cause_value = map[string]int32{
	"TxIdAlreadyExists":      0,
	"RWConflictDuringCommit": 1,
}

func (x InvalidTransaction_Cause) String() string {
	return proto.EnumName(InvalidTransaction_Cause_name, int32(x))
}
func (InvalidTransaction_Cause) EnumDescriptor() ([]byte, []int) { return fileDescriptor8, []int{1, 0} }

// This message is necessary to facilitate the verification of the signature
// (in the signature field) over the bytes of the transaction (in the
// transactionBytes field).
type SignedTransaction struct {
	// The bytes of the Transaction. NDD
	TransactionBytes []byte `protobuf:"bytes,1,opt,name=transactionBytes,proto3" json:"transactionBytes,omitempty"`
	// Signature of the transactionBytes The public key of the signature is in
	// the header field of TransactionAction There might be multiple
	// TransactionAction, so multiple headers, but there should be same
	// transactor identity (cert) in all headers
	Signature []byte `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (m *SignedTransaction) Reset()                    { *m = SignedTransaction{} }
func (m *SignedTransaction) String() string            { return proto.CompactTextString(m) }
func (*SignedTransaction) ProtoMessage()               {}
func (*SignedTransaction) Descriptor() ([]byte, []int) { return fileDescriptor8, []int{0} }

// This is used to wrap an invalid Transaction with the cause
type InvalidTransaction struct {
	Transaction *Transaction             `protobuf:"bytes,1,opt,name=transaction" json:"transaction,omitempty"`
	Cause       InvalidTransaction_Cause `protobuf:"varint,2,opt,name=cause,enum=protos.InvalidTransaction_Cause" json:"cause,omitempty"`
}

func (m *InvalidTransaction) Reset()                    { *m = InvalidTransaction{} }
func (m *InvalidTransaction) String() string            { return proto.CompactTextString(m) }
func (*InvalidTransaction) ProtoMessage()               {}
func (*InvalidTransaction) Descriptor() ([]byte, []int) { return fileDescriptor8, []int{1} }

func (m *InvalidTransaction) GetTransaction() *Transaction {
	if m != nil {
		return m.Transaction
	}
	return nil
}

// The transaction to be sent to the ordering service. A transaction contains
// one or more TransactionAction. Each TransactionAction binds a proposal to
// potentially multiple actions. The transaction is atomic meaning that either
// all actions in the transaction will be committed or none will.  Note that
// while a Transaction might include more than one Header, the Header.creator
// field must be the same in each.
// A single client is free to issue a number of independent Proposal, each with
// their header (Header) and request payload (ChaincodeProposalPayload).  Each
// proposal is independently endorsed generating an action
// (ProposalResponsePayload) with one signature per Endorser. Any number of
// independent proposals (and their action) might be included in a transaction
// to ensure that they are treated atomically.
type Transaction struct {
	// Version indicates message protocol version.
	Version int32 `protobuf:"varint,1,opt,name=version" json:"version,omitempty"`
	// Timestamp is the local time that the
	// message was created by the sender
	Timestamp *google_protobuf1.Timestamp `protobuf:"bytes,2,opt,name=timestamp" json:"timestamp,omitempty"`
	// The payload is an array of TransactionAction. An array is necessary to
	// accommodate multiple actions per transaction
	Actions []*TransactionAction `protobuf:"bytes,3,rep,name=actions" json:"actions,omitempty"`
}

func (m *Transaction) Reset()                    { *m = Transaction{} }
func (m *Transaction) String() string            { return proto.CompactTextString(m) }
func (*Transaction) ProtoMessage()               {}
func (*Transaction) Descriptor() ([]byte, []int) { return fileDescriptor8, []int{2} }

func (m *Transaction) GetTimestamp() *google_protobuf1.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *Transaction) GetActions() []*TransactionAction {
	if m != nil {
		return m.Actions
	}
	return nil
}

// TransactionAction binds a proposal to its action.  The type field in the
// header dictates the type of action to be applied to the ledger.
type TransactionAction struct {
	// The header of the proposal action, which is the proposal header
	Header []byte `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	// The payload of the action as defined by the type in the header For
	// chaincode, it's the bytes of ChaincodeActionPayload
	Payload []byte `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (m *TransactionAction) Reset()                    { *m = TransactionAction{} }
func (m *TransactionAction) String() string            { return proto.CompactTextString(m) }
func (*TransactionAction) ProtoMessage()               {}
func (*TransactionAction) Descriptor() ([]byte, []int) { return fileDescriptor8, []int{3} }

// ChaincodeActionPayload is the message to be used for the TransactionAction's
// payload when the Header's type is set to CHAINCODE.  It carries the
// chaincodeProposalPayload and an endorsed action to apply to the ledger.
type ChaincodeActionPayload struct {
	// This field contains the bytes of the ChaincodeProposalPayload message from
	// the original invocation (essentially the arguments) after the application
	// of the visibility function. The main visibility modes are "full" (the
	// entire ChaincodeProposalPayload message is included here), "hash" (only
	// the hash of the ChaincodeProposalPayload message is included) or
	// "nothing".  This field will be used to check the consistency of
	// ProposalResponsePayload.proposalHash.  For the CHAINCODE type,
	// ProposalResponsePayload.proposalHash is supposed to be H(ProposalHeader ||
	// f(ChaincodeProposalPayload)) where f is the visibility function.
	ChaincodeProposalPayload []byte `protobuf:"bytes,1,opt,name=chaincodeProposalPayload,proto3" json:"chaincodeProposalPayload,omitempty"`
	// The list of actions to apply to the ledger
	Action *ChaincodeEndorsedAction `protobuf:"bytes,2,opt,name=action" json:"action,omitempty"`
}

func (m *ChaincodeActionPayload) Reset()                    { *m = ChaincodeActionPayload{} }
func (m *ChaincodeActionPayload) String() string            { return proto.CompactTextString(m) }
func (*ChaincodeActionPayload) ProtoMessage()               {}
func (*ChaincodeActionPayload) Descriptor() ([]byte, []int) { return fileDescriptor8, []int{4} }

func (m *ChaincodeActionPayload) GetAction() *ChaincodeEndorsedAction {
	if m != nil {
		return m.Action
	}
	return nil
}

// ChaincodeEndorsedAction carries information about the endorsement of a
// specific proposal
type ChaincodeEndorsedAction struct {
	// This is the bytes of the ProposalResponsePayload message signed by the
	// endorsers.  Recall that for the CHAINCODE type, the
	// ProposalResponsePayload's extenstion field carries a ChaincodeAction
	ProposalResponsePayload []byte `protobuf:"bytes,1,opt,name=proposalResponsePayload,proto3" json:"proposalResponsePayload,omitempty"`
	// The endorsement of the proposal, basically the endorser's signature over
	// proposalResponsePayload
	Endorsements []*Endorsement `protobuf:"bytes,2,rep,name=endorsements" json:"endorsements,omitempty"`
}

func (m *ChaincodeEndorsedAction) Reset()                    { *m = ChaincodeEndorsedAction{} }
func (m *ChaincodeEndorsedAction) String() string            { return proto.CompactTextString(m) }
func (*ChaincodeEndorsedAction) ProtoMessage()               {}
func (*ChaincodeEndorsedAction) Descriptor() ([]byte, []int) { return fileDescriptor8, []int{5} }

func (m *ChaincodeEndorsedAction) GetEndorsements() []*Endorsement {
	if m != nil {
		return m.Endorsements
	}
	return nil
}

func init() {
	proto.RegisterType((*SignedTransaction)(nil), "protos.SignedTransaction")
	proto.RegisterType((*InvalidTransaction)(nil), "protos.InvalidTransaction")
	proto.RegisterType((*Transaction)(nil), "protos.Transaction")
	proto.RegisterType((*TransactionAction)(nil), "protos.TransactionAction")
	proto.RegisterType((*ChaincodeActionPayload)(nil), "protos.ChaincodeActionPayload")
	proto.RegisterType((*ChaincodeEndorsedAction)(nil), "protos.ChaincodeEndorsedAction")
	proto.RegisterEnum("protos.InvalidTransaction_Cause", InvalidTransaction_Cause_name, InvalidTransaction_Cause_value)
}

func init() { proto.RegisterFile("peer/transaction.proto", fileDescriptor8) }

var fileDescriptor8 = []byte{
	// 481 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x74, 0x93, 0x51, 0x6b, 0xdb, 0x30,
	0x10, 0xc7, 0xe7, 0x96, 0xa4, 0xf4, 0x5c, 0x46, 0xa2, 0xb1, 0xd4, 0x0b, 0x85, 0x06, 0x3f, 0x75,
	0x1b, 0xd8, 0x90, 0xb2, 0xb5, 0xf4, 0xad, 0xcd, 0xf2, 0xd0, 0xb7, 0xe2, 0x05, 0x06, 0x83, 0x31,
	0x14, 0xfb, 0xe2, 0x08, 0x6c, 0xc9, 0x48, 0x72, 0x69, 0xbe, 0xc3, 0xf6, 0xba, 0xaf, 0xb3, 0xaf,
	0x36, 0x62, 0x49, 0x89, 0x43, 0x96, 0x27, 0x73, 0xfa, 0xff, 0xfc, 0xbf, 0xd3, 0xdd, 0x09, 0x06,
	0x15, 0xa2, 0x8c, 0xb5, 0xa4, 0x5c, 0xd1, 0x54, 0x33, 0xc1, 0xa3, 0x4a, 0x0a, 0x2d, 0x48, 0xb7,
	0xf9, 0xa8, 0xe1, 0x65, 0x2e, 0x44, 0x5e, 0x60, 0xdc, 0x84, 0xf3, 0x7a, 0x11, 0x6b, 0x56, 0xa2,
	0xd2, 0xb4, 0xac, 0x0c, 0x38, 0xbc, 0x68, 0x0c, 0x2a, 0x29, 0x2a, 0xa1, 0x68, 0xf1, 0x53, 0xa2,
	0xaa, 0x04, 0x57, 0x68, 0xd4, 0xf0, 0x07, 0xf4, 0xbf, 0xb2, 0x9c, 0x63, 0x36, 0xdb, 0x66, 0x20,
	0x1f, 0xa0, 0xd7, 0x4a, 0xf8, 0xb0, 0xd2, 0xa8, 0x02, 0x6f, 0xe4, 0x5d, 0x9d, 0x25, 0x7b, 0xe7,
	0xe4, 0x02, 0x4e, 0x15, 0xcb, 0x39, 0xd5, 0xb5, 0xc4, 0xe0, 0xa8, 0x81, 0xb6, 0x07, 0xe1, 0x5f,
	0x0f, 0xc8, 0x23, 0x7f, 0xa6, 0x05, 0xdb, 0x49, 0xf0, 0x09, 0xfc, 0x96, 0x51, 0xe3, 0xed, 0x8f,
	0xdf, 0x98, 0x92, 0x54, 0xd4, 0x22, 0x93, 0x36, 0x47, 0x3e, 0x43, 0x27, 0xa5, 0xb5, 0x32, 0x79,
	0x5e, 0x8f, 0x47, 0xee, 0x87, 0xfd, 0x0c, 0xd1, 0x64, 0xcd, 0x25, 0x06, 0x0f, 0xef, 0xa0, 0xd3,
	0xc4, 0xe4, 0x2d, 0xf4, 0x67, 0x2f, 0x8f, 0xd9, 0x7d, 0x21, 0x91, 0x66, 0xab, 0xe9, 0x0b, 0x53,
	0x5a, 0xf5, 0x5e, 0x91, 0x21, 0x0c, 0x92, 0x6f, 0x13, 0xc1, 0x17, 0x05, 0x4b, 0xf5, 0x97, 0x5a,
	0x32, 0x9e, 0x4f, 0x44, 0x59, 0x32, 0xdd, 0xf3, 0xc2, 0x3f, 0x1e, 0xf8, 0xed, 0xd2, 0x03, 0x38,
	0x79, 0x46, 0xa9, 0x5c, 0xd9, 0x9d, 0xc4, 0x85, 0xe4, 0x16, 0x4e, 0x37, 0xbd, 0x6f, 0x2a, 0xf4,
	0xc7, 0xc3, 0xc8, 0x4c, 0x27, 0x72, 0xd3, 0x89, 0x66, 0x8e, 0x48, 0xb6, 0x30, 0xb9, 0x86, 0x13,
	0xe3, 0xae, 0x82, 0xe3, 0xd1, 0xf1, 0x95, 0x3f, 0x7e, 0xf7, 0x9f, 0x56, 0xdc, 0x9b, 0x86, 0x38,
	0x32, 0x9c, 0x42, 0x7f, 0x4f, 0x25, 0x03, 0xe8, 0x2e, 0x91, 0x66, 0x28, 0xed, 0xbc, 0x6c, 0xb4,
	0xae, 0xba, 0xa2, 0xab, 0x42, 0xd0, 0xcc, 0xce, 0xc8, 0x85, 0xe1, 0x6f, 0x0f, 0x06, 0x93, 0x25,
	0x65, 0x3c, 0x15, 0x19, 0x1a, 0x97, 0x27, 0x23, 0x91, 0x3b, 0x08, 0x52, 0xa7, 0x3c, 0xd9, 0xfd,
	0xb1, 0x9a, 0xb5, 0x3f, 0xa8, 0x93, 0x1b, 0xe8, 0xda, 0xe1, 0x9a, 0x4e, 0x5c, 0xba, 0x1b, 0x6d,
	0x72, 0x4d, 0x79, 0x26, 0xa4, 0xc2, 0xcc, 0xde, 0xcb, 0xe2, 0xe1, 0x2f, 0x0f, 0xce, 0x0f, 0x30,
	0xe4, 0x16, 0xce, 0xdd, 0x1e, 0x27, 0x76, 0x8d, 0x77, 0xeb, 0x39, 0x24, 0x93, 0x1b, 0x38, 0x43,
	0xe3, 0x55, 0x22, 0xd7, 0x2a, 0x38, 0x6a, 0xda, 0xbc, 0xd9, 0xb8, 0xe9, 0x56, 0x4b, 0x76, 0xc0,
	0x87, 0x8f, 0xdf, 0xdf, 0xe7, 0x4c, 0x2f, 0xeb, 0x79, 0x94, 0x8a, 0x32, 0x5e, 0xae, 0x2a, 0x94,
	0x05, 0x66, 0x39, 0xca, 0x78, 0x41, 0xe7, 0x92, 0xa5, 0xe6, 0xdd, 0xa9, 0x78, 0xfd, 0xc8, 0xe6,
	0xe6, 0x4d, 0x5e, 0xff, 0x0b, 0x00, 0x00, 0xff, 0xff, 0x52, 0xfc, 0xf5, 0x8d, 0xb4, 0x03, 0x00,
	0x00,
}
