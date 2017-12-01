// Code generated by protoc-gen-go. DO NOT EDIT.
// source: meta_data.proto

package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type MetaData struct {
	Term            uint64              `protobuf:"varint,1,opt,name=term" json:"term,omitempty"`
	Index           uint64              `protobuf:"varint,2,opt,name=index" json:"index,omitempty"`
	MaxShardGroupID uint64              `protobuf:"varint,5,opt,name=maxShardGroupID" json:"maxShardGroupID,omitempty"`
	MaxShardID      uint64              `protobuf:"varint,6,opt,name=maxShardID" json:"maxShardID,omitempty"`
	MetaNodes       []*MetaNodeInfo     `protobuf:"bytes,7,rep,name=metaNodes" json:"metaNodes,omitempty"`
	DataNodes       []*MetaNodeInfo     `protobuf:"bytes,8,rep,name=dataNodes" json:"dataNodes,omitempty"`
	Databases       []*MetaDatabaseInfo `protobuf:"bytes,9,rep,name=databases" json:"databases,omitempty"`
	Users           []*MetaUserInfo     `protobuf:"bytes,10,rep,name=users" json:"users,omitempty"`
}

func (m *MetaData) Reset()                    { *m = MetaData{} }
func (m *MetaData) String() string            { return proto.CompactTextString(m) }
func (*MetaData) ProtoMessage()               {}
func (*MetaData) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

func (m *MetaData) GetTerm() uint64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *MetaData) GetIndex() uint64 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *MetaData) GetMaxShardGroupID() uint64 {
	if m != nil {
		return m.MaxShardGroupID
	}
	return 0
}

func (m *MetaData) GetMaxShardID() uint64 {
	if m != nil {
		return m.MaxShardID
	}
	return 0
}

func (m *MetaData) GetMetaNodes() []*MetaNodeInfo {
	if m != nil {
		return m.MetaNodes
	}
	return nil
}

func (m *MetaData) GetDataNodes() []*MetaNodeInfo {
	if m != nil {
		return m.DataNodes
	}
	return nil
}

func (m *MetaData) GetDatabases() []*MetaDatabaseInfo {
	if m != nil {
		return m.Databases
	}
	return nil
}

func (m *MetaData) GetUsers() []*MetaUserInfo {
	if m != nil {
		return m.Users
	}
	return nil
}

type MetaLease struct {
	Name       string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Expiration int64  `protobuf:"varint,2,opt,name=expiration" json:"expiration,omitempty"`
	Owner      uint64 `protobuf:"varint,3,opt,name=owner" json:"owner,omitempty"`
}

func (m *MetaLease) Reset()                    { *m = MetaLease{} }
func (m *MetaLease) String() string            { return proto.CompactTextString(m) }
func (*MetaLease) ProtoMessage()               {}
func (*MetaLease) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{1} }

func (m *MetaLease) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *MetaLease) GetExpiration() int64 {
	if m != nil {
		return m.Expiration
	}
	return 0
}

func (m *MetaLease) GetOwner() uint64 {
	if m != nil {
		return m.Owner
	}
	return 0
}

type MetaNodeInfo struct {
	ID      uint64 `protobuf:"varint,1,opt,name=ID" json:"ID,omitempty"`
	Address string `protobuf:"bytes,2,opt,name=address" json:"address,omitempty"`
}

func (m *MetaNodeInfo) Reset()                    { *m = MetaNodeInfo{} }
func (m *MetaNodeInfo) String() string            { return proto.CompactTextString(m) }
func (*MetaNodeInfo) ProtoMessage()               {}
func (*MetaNodeInfo) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{2} }

func (m *MetaNodeInfo) GetID() uint64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *MetaNodeInfo) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

type MetaDatabaseInfo struct {
	Name                   string                     `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	DefaultRetentionPolicy string                     `protobuf:"bytes,2,opt,name=defaultRetentionPolicy" json:"defaultRetentionPolicy,omitempty"`
	RetentionPolicies      []*MetaRetentionPolicyInfo `protobuf:"bytes,3,rep,name=retentionPolicies" json:"retentionPolicies,omitempty"`
	ContinuousQueries      []*MetaContinuousQueryInfo `protobuf:"bytes,4,rep,name=continuousQueries" json:"continuousQueries,omitempty"`
}

func (m *MetaDatabaseInfo) Reset()                    { *m = MetaDatabaseInfo{} }
func (m *MetaDatabaseInfo) String() string            { return proto.CompactTextString(m) }
func (*MetaDatabaseInfo) ProtoMessage()               {}
func (*MetaDatabaseInfo) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{3} }

func (m *MetaDatabaseInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *MetaDatabaseInfo) GetDefaultRetentionPolicy() string {
	if m != nil {
		return m.DefaultRetentionPolicy
	}
	return ""
}

func (m *MetaDatabaseInfo) GetRetentionPolicies() []*MetaRetentionPolicyInfo {
	if m != nil {
		return m.RetentionPolicies
	}
	return nil
}

func (m *MetaDatabaseInfo) GetContinuousQueries() []*MetaContinuousQueryInfo {
	if m != nil {
		return m.ContinuousQueries
	}
	return nil
}

type MetaRetentionPolicySpec struct {
	Name               *Str `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Duration           *I64 `protobuf:"bytes,2,opt,name=duration" json:"duration,omitempty"`
	ShardGroupDuration *I64 `protobuf:"bytes,3,opt,name=shardGroupDuration" json:"shardGroupDuration,omitempty"`
	ReplicaN           *I64 `protobuf:"bytes,4,opt,name=replicaN" json:"replicaN,omitempty"`
}

func (m *MetaRetentionPolicySpec) Reset()                    { *m = MetaRetentionPolicySpec{} }
func (m *MetaRetentionPolicySpec) String() string            { return proto.CompactTextString(m) }
func (*MetaRetentionPolicySpec) ProtoMessage()               {}
func (*MetaRetentionPolicySpec) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{4} }

func (m *MetaRetentionPolicySpec) GetName() *Str {
	if m != nil {
		return m.Name
	}
	return nil
}

func (m *MetaRetentionPolicySpec) GetDuration() *I64 {
	if m != nil {
		return m.Duration
	}
	return nil
}

func (m *MetaRetentionPolicySpec) GetShardGroupDuration() *I64 {
	if m != nil {
		return m.ShardGroupDuration
	}
	return nil
}

func (m *MetaRetentionPolicySpec) GetReplicaN() *I64 {
	if m != nil {
		return m.ReplicaN
	}
	return nil
}

type MetaRetentionPolicyInfo struct {
	Name               string                  `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Duration           int64                   `protobuf:"varint,2,opt,name=duration" json:"duration,omitempty"`
	ShardGroupDuration int64                   `protobuf:"varint,3,opt,name=shardGroupDuration" json:"shardGroupDuration,omitempty"`
	ReplicaN           int64                   `protobuf:"varint,4,opt,name=replicaN" json:"replicaN,omitempty"`
	ShardGroups        []*MetaShardGroupInfo   `protobuf:"bytes,5,rep,name=shardGroups" json:"shardGroups,omitempty"`
	Subscriptions      []*MetaSubscriptionInfo `protobuf:"bytes,6,rep,name=subscriptions" json:"subscriptions,omitempty"`
}

func (m *MetaRetentionPolicyInfo) Reset()                    { *m = MetaRetentionPolicyInfo{} }
func (m *MetaRetentionPolicyInfo) String() string            { return proto.CompactTextString(m) }
func (*MetaRetentionPolicyInfo) ProtoMessage()               {}
func (*MetaRetentionPolicyInfo) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{5} }

func (m *MetaRetentionPolicyInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *MetaRetentionPolicyInfo) GetDuration() int64 {
	if m != nil {
		return m.Duration
	}
	return 0
}

func (m *MetaRetentionPolicyInfo) GetShardGroupDuration() int64 {
	if m != nil {
		return m.ShardGroupDuration
	}
	return 0
}

func (m *MetaRetentionPolicyInfo) GetReplicaN() int64 {
	if m != nil {
		return m.ReplicaN
	}
	return 0
}

func (m *MetaRetentionPolicyInfo) GetShardGroups() []*MetaShardGroupInfo {
	if m != nil {
		return m.ShardGroups
	}
	return nil
}

func (m *MetaRetentionPolicyInfo) GetSubscriptions() []*MetaSubscriptionInfo {
	if m != nil {
		return m.Subscriptions
	}
	return nil
}

type MetaShardGroupInfo struct {
	ID          uint64           `protobuf:"varint,1,opt,name=ID" json:"ID,omitempty"`
	StartTime   int64            `protobuf:"varint,2,opt,name=startTime" json:"startTime,omitempty"`
	EndTime     int64            `protobuf:"varint,3,opt,name=endTime" json:"endTime,omitempty"`
	DeletedAt   int64            `protobuf:"varint,4,opt,name=deletedAt" json:"deletedAt,omitempty"`
	Shards      []*MetaShardInfo `protobuf:"bytes,5,rep,name=shards" json:"shards,omitempty"`
	TruncatedAt int64            `protobuf:"varint,6,opt,name=truncatedAt" json:"truncatedAt,omitempty"`
}

func (m *MetaShardGroupInfo) Reset()                    { *m = MetaShardGroupInfo{} }
func (m *MetaShardGroupInfo) String() string            { return proto.CompactTextString(m) }
func (*MetaShardGroupInfo) ProtoMessage()               {}
func (*MetaShardGroupInfo) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{6} }

func (m *MetaShardGroupInfo) GetID() uint64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *MetaShardGroupInfo) GetStartTime() int64 {
	if m != nil {
		return m.StartTime
	}
	return 0
}

func (m *MetaShardGroupInfo) GetEndTime() int64 {
	if m != nil {
		return m.EndTime
	}
	return 0
}

func (m *MetaShardGroupInfo) GetDeletedAt() int64 {
	if m != nil {
		return m.DeletedAt
	}
	return 0
}

func (m *MetaShardGroupInfo) GetShards() []*MetaShardInfo {
	if m != nil {
		return m.Shards
	}
	return nil
}

func (m *MetaShardGroupInfo) GetTruncatedAt() int64 {
	if m != nil {
		return m.TruncatedAt
	}
	return 0
}

type MetaShardInfo struct {
	ID     uint64            `protobuf:"varint,1,opt,name=ID" json:"ID,omitempty"`
	Owners []*MetaShardOwner `protobuf:"bytes,2,rep,name=owners" json:"owners,omitempty"`
}

func (m *MetaShardInfo) Reset()                    { *m = MetaShardInfo{} }
func (m *MetaShardInfo) String() string            { return proto.CompactTextString(m) }
func (*MetaShardInfo) ProtoMessage()               {}
func (*MetaShardInfo) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{7} }

func (m *MetaShardInfo) GetID() uint64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *MetaShardInfo) GetOwners() []*MetaShardOwner {
	if m != nil {
		return m.Owners
	}
	return nil
}

type MetaSubscriptionInfo struct {
	Name         string   `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Mode         string   `protobuf:"bytes,2,opt,name=mode" json:"mode,omitempty"`
	Destinations []string `protobuf:"bytes,3,rep,name=Destinations" json:"Destinations,omitempty"`
}

func (m *MetaSubscriptionInfo) Reset()                    { *m = MetaSubscriptionInfo{} }
func (m *MetaSubscriptionInfo) String() string            { return proto.CompactTextString(m) }
func (*MetaSubscriptionInfo) ProtoMessage()               {}
func (*MetaSubscriptionInfo) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{8} }

func (m *MetaSubscriptionInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *MetaSubscriptionInfo) GetMode() string {
	if m != nil {
		return m.Mode
	}
	return ""
}

func (m *MetaSubscriptionInfo) GetDestinations() []string {
	if m != nil {
		return m.Destinations
	}
	return nil
}

type MetaShardOwner struct {
	NodeID uint64 `protobuf:"varint,1,opt,name=nodeID" json:"nodeID,omitempty"`
}

func (m *MetaShardOwner) Reset()                    { *m = MetaShardOwner{} }
func (m *MetaShardOwner) String() string            { return proto.CompactTextString(m) }
func (*MetaShardOwner) ProtoMessage()               {}
func (*MetaShardOwner) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{9} }

func (m *MetaShardOwner) GetNodeID() uint64 {
	if m != nil {
		return m.NodeID
	}
	return 0
}

type MetaContinuousQueryInfo struct {
	Name  string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Query string `protobuf:"bytes,2,opt,name=query" json:"query,omitempty"`
}

func (m *MetaContinuousQueryInfo) Reset()                    { *m = MetaContinuousQueryInfo{} }
func (m *MetaContinuousQueryInfo) String() string            { return proto.CompactTextString(m) }
func (*MetaContinuousQueryInfo) ProtoMessage()               {}
func (*MetaContinuousQueryInfo) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{10} }

func (m *MetaContinuousQueryInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *MetaContinuousQueryInfo) GetQuery() string {
	if m != nil {
		return m.Query
	}
	return ""
}

type MetaUserInfo struct {
	Name       string               `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Hash       string               `protobuf:"bytes,2,opt,name=hash" json:"hash,omitempty"`
	Admin      bool                 `protobuf:"varint,3,opt,name=admin" json:"admin,omitempty"`
	Privileges []*MetaUserPrivilege `protobuf:"bytes,4,rep,name=privileges" json:"privileges,omitempty"`
}

func (m *MetaUserInfo) Reset()                    { *m = MetaUserInfo{} }
func (m *MetaUserInfo) String() string            { return proto.CompactTextString(m) }
func (*MetaUserInfo) ProtoMessage()               {}
func (*MetaUserInfo) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{11} }

func (m *MetaUserInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *MetaUserInfo) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

func (m *MetaUserInfo) GetAdmin() bool {
	if m != nil {
		return m.Admin
	}
	return false
}

func (m *MetaUserInfo) GetPrivileges() []*MetaUserPrivilege {
	if m != nil {
		return m.Privileges
	}
	return nil
}

type MetaUserPrivilege struct {
	Database  string `protobuf:"bytes,1,opt,name=database" json:"database,omitempty"`
	Privilege int32  `protobuf:"varint,2,opt,name=privilege" json:"privilege,omitempty"`
}

func (m *MetaUserPrivilege) Reset()                    { *m = MetaUserPrivilege{} }
func (m *MetaUserPrivilege) String() string            { return proto.CompactTextString(m) }
func (*MetaUserPrivilege) ProtoMessage()               {}
func (*MetaUserPrivilege) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{12} }

func (m *MetaUserPrivilege) GetDatabase() string {
	if m != nil {
		return m.Database
	}
	return ""
}

func (m *MetaUserPrivilege) GetPrivilege() int32 {
	if m != nil {
		return m.Privilege
	}
	return 0
}

func init() {
	proto.RegisterType((*MetaData)(nil), "pb.MetaData")
	proto.RegisterType((*MetaLease)(nil), "pb.MetaLease")
	proto.RegisterType((*MetaNodeInfo)(nil), "pb.MetaNodeInfo")
	proto.RegisterType((*MetaDatabaseInfo)(nil), "pb.MetaDatabaseInfo")
	proto.RegisterType((*MetaRetentionPolicySpec)(nil), "pb.MetaRetentionPolicySpec")
	proto.RegisterType((*MetaRetentionPolicyInfo)(nil), "pb.MetaRetentionPolicyInfo")
	proto.RegisterType((*MetaShardGroupInfo)(nil), "pb.MetaShardGroupInfo")
	proto.RegisterType((*MetaShardInfo)(nil), "pb.MetaShardInfo")
	proto.RegisterType((*MetaSubscriptionInfo)(nil), "pb.MetaSubscriptionInfo")
	proto.RegisterType((*MetaShardOwner)(nil), "pb.MetaShardOwner")
	proto.RegisterType((*MetaContinuousQueryInfo)(nil), "pb.MetaContinuousQueryInfo")
	proto.RegisterType((*MetaUserInfo)(nil), "pb.MetaUserInfo")
	proto.RegisterType((*MetaUserPrivilege)(nil), "pb.MetaUserPrivilege")
}

func init() { proto.RegisterFile("meta_data.proto", fileDescriptor3) }

var fileDescriptor3 = []byte{
	// 750 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x55, 0xdd, 0x6e, 0xd3, 0x4a,
	0x10, 0x56, 0xe2, 0x24, 0x4d, 0x26, 0xfd, 0x5d, 0xe5, 0xf4, 0xac, 0xda, 0xa3, 0xa3, 0x6a, 0x8f,
	0x74, 0x14, 0xb8, 0xc8, 0x45, 0x81, 0xd2, 0x2b, 0x24, 0x54, 0x4b, 0x28, 0x82, 0x96, 0xe2, 0xd2,
	0x6b, 0xb4, 0x89, 0xa7, 0xd4, 0x52, 0xfc, 0xc3, 0xee, 0x1a, 0xca, 0x35, 0x0f, 0xc0, 0xf3, 0xc0,
	0x0b, 0xf0, 0x42, 0x3c, 0x00, 0xda, 0x1f, 0xc7, 0x76, 0xe2, 0xdc, 0x79, 0xbe, 0xf9, 0x66, 0xf6,
	0x9b, 0x99, 0x9d, 0x35, 0xec, 0xc5, 0xa8, 0xf8, 0x87, 0x90, 0x2b, 0x3e, 0xc9, 0x44, 0xaa, 0x52,
	0xd2, 0xce, 0x66, 0x47, 0xdb, 0xf3, 0x34, 0x8e, 0xd3, 0xc4, 0x22, 0xec, 0x67, 0x1b, 0xfa, 0x97,
	0xa8, 0xb8, 0xcf, 0x15, 0x27, 0x04, 0x3a, 0x0a, 0x45, 0x4c, 0x5b, 0x27, 0xad, 0x71, 0x27, 0x30,
	0xdf, 0x64, 0x04, 0xdd, 0x28, 0x09, 0xf1, 0x81, 0xb6, 0x0d, 0x68, 0x0d, 0x32, 0x86, 0xbd, 0x98,
	0x3f, 0xdc, 0xdc, 0x73, 0x11, 0xbe, 0x12, 0x69, 0x9e, 0x4d, 0x7d, 0xda, 0x35, 0xfe, 0x55, 0x98,
	0xfc, 0x0b, 0x50, 0x40, 0x53, 0x9f, 0xf6, 0x0c, 0xa9, 0x82, 0x90, 0x09, 0x0c, 0xb4, 0xca, 0xab,
	0x34, 0x44, 0x49, 0xb7, 0x4e, 0xbc, 0xf1, 0xf0, 0x74, 0x7f, 0x92, 0xcd, 0x26, 0x97, 0x0e, 0x9c,
	0x26, 0x77, 0x69, 0x50, 0x52, 0x34, 0x5f, 0x17, 0x64, 0xf9, 0xfd, 0x4d, 0xfc, 0x25, 0x85, 0x9c,
	0x5a, 0xfe, 0x8c, 0x4b, 0x94, 0x74, 0x60, 0xf8, 0xa3, 0x82, 0xef, 0x3b, 0x47, 0x19, 0x63, 0x68,
	0xe4, 0x7f, 0xe8, 0xe6, 0x12, 0x85, 0xa4, 0x50, 0xcf, 0x7f, 0x2b, 0x51, 0x18, 0xae, 0x75, 0xb3,
	0x5b, 0x18, 0x68, 0xf8, 0x0d, 0x72, 0x89, 0xba, 0x79, 0x09, 0x8f, 0xd1, 0x34, 0x6f, 0x10, 0x98,
	0x6f, 0x5d, 0x3c, 0x3e, 0x64, 0x91, 0xe0, 0x2a, 0x4a, 0x13, 0xd3, 0x41, 0x2f, 0xa8, 0x20, 0xba,
	0xb9, 0xe9, 0x97, 0x04, 0x05, 0xf5, 0x6c, 0x73, 0x8d, 0xc1, 0xce, 0x61, 0xbb, 0x5a, 0x0d, 0xd9,
	0x85, 0xf6, 0xd4, 0x77, 0x43, 0x69, 0x4f, 0x7d, 0x42, 0x61, 0x8b, 0x87, 0xa1, 0x40, 0x29, 0x4d,
	0xca, 0x41, 0x50, 0x98, 0xec, 0x77, 0x0b, 0xf6, 0x57, 0x0b, 0x6b, 0x14, 0x76, 0x06, 0x87, 0x21,
	0xde, 0xf1, 0x7c, 0xa1, 0x02, 0x54, 0x98, 0x68, 0x31, 0xd7, 0xe9, 0x22, 0x9a, 0x7f, 0x75, 0x19,
	0x37, 0x78, 0xc9, 0x14, 0x0e, 0x44, 0x0d, 0x8a, 0x50, 0x52, 0xcf, 0x74, 0xe9, 0xb8, 0xe8, 0xd2,
	0x4a, 0x8c, 0x69, 0xd8, 0x7a, 0x94, 0x4e, 0x35, 0x4f, 0x13, 0x15, 0x25, 0x79, 0x9a, 0xcb, 0x77,
	0x39, 0x0a, 0x9d, 0xaa, 0x53, 0x4f, 0x75, 0x51, 0x23, 0xb8, 0x54, 0x6b, 0x51, 0xec, 0x47, 0x0b,
	0xfe, 0x6e, 0x38, 0xf9, 0x26, 0xc3, 0x39, 0x39, 0xae, 0x54, 0x3f, 0x3c, 0xdd, 0xd2, 0x99, 0x6f,
	0x94, 0x70, 0x6d, 0xf8, 0x0f, 0xfa, 0x61, 0x5e, 0x99, 0x8e, 0x23, 0x4c, 0xcf, 0x9e, 0x06, 0x4b,
	0x07, 0x79, 0x0e, 0x44, 0x2e, 0x6f, 0xb4, 0x5f, 0xd0, 0xbd, 0x3a, 0xbd, 0x81, 0xa2, 0xb3, 0x0b,
	0xcc, 0x16, 0xd1, 0x9c, 0x5f, 0xd1, 0xce, 0x4a, 0xf6, 0xc2, 0xc1, 0xbe, 0xb7, 0x1b, 0xb5, 0x6f,
	0x9c, 0xdc, 0xd1, 0x8a, 0x64, 0xaf, 0xa2, 0x74, 0xb2, 0x51, 0xa9, 0xd7, 0x28, 0xf0, 0x68, 0x45,
	0xa0, 0x57, 0xea, 0x22, 0xe7, 0x30, 0x2c, 0x23, 0x24, 0xed, 0x9a, 0xc1, 0x1c, 0x16, 0x83, 0xa9,
	0xac, 0xb8, 0x9e, 0x49, 0x95, 0x4a, 0x5e, 0xc0, 0x8e, 0xcc, 0x67, 0x72, 0x2e, 0xa2, 0x4c, 0x9f,
	0x22, 0x69, 0xcf, 0xc4, 0xd2, 0x65, 0x6c, 0xc5, 0x69, 0xa2, 0xeb, 0x74, 0xf6, 0xab, 0x05, 0x64,
	0xfd, 0x8c, 0xb5, 0x2d, 0xf8, 0x07, 0x06, 0x52, 0x71, 0xa1, 0xde, 0x47, 0x31, 0xba, 0x4e, 0x94,
	0x80, 0xde, 0x11, 0x4c, 0x42, 0xe3, 0xb3, 0xf5, 0x17, 0xa6, 0x8e, 0x0b, 0x71, 0x81, 0x0a, 0xc3,
	0x97, 0xca, 0x55, 0x5d, 0x02, 0xe4, 0x11, 0xf4, 0x4c, 0x2d, 0x45, 0xc5, 0x07, 0xb5, 0x8a, 0x8d,
	0x5c, 0x47, 0x20, 0x27, 0x30, 0x54, 0x22, 0x4f, 0xe6, 0xdc, 0xa6, 0xea, 0x99, 0x54, 0x55, 0x88,
	0xbd, 0x86, 0x9d, 0x5a, 0xe8, 0x5a, 0x0d, 0x8f, 0xa1, 0x67, 0x56, 0x5e, 0x2f, 0xb2, 0x3e, 0x8d,
	0xd4, 0x4e, 0x7b, 0xab, 0x5d, 0x81, 0x63, 0xb0, 0x19, 0x8c, 0x9a, 0xba, 0xd7, 0x78, 0x49, 0x08,
	0x74, 0xe2, 0x34, 0x44, 0xb7, 0xcc, 0xe6, 0x9b, 0x30, 0xd8, 0xf6, 0x51, 0xaa, 0x28, 0xe1, 0x76,
	0x2a, 0x7a, 0x6b, 0x07, 0x41, 0x0d, 0x63, 0x63, 0xd8, 0xad, 0x9f, 0x4e, 0x0e, 0xa1, 0x97, 0xe8,
	0x77, 0xa8, 0x50, 0xed, 0x2c, 0x76, 0x61, 0x6f, 0x6d, 0xc3, 0x82, 0x36, 0x0a, 0x1a, 0x41, 0xf7,
	0x93, 0x26, 0x38, 0x45, 0xd6, 0x60, 0xdf, 0x5a, 0xf6, 0xa5, 0x2b, 0xde, 0xd5, 0x4d, 0xb5, 0xdc,
	0x73, 0x79, 0x5f, 0xd4, 0xa2, 0xbf, 0x75, 0x3a, 0x1e, 0xc6, 0x91, 0xbd, 0xdb, 0xfd, 0xc0, 0x1a,
	0xe4, 0x19, 0x40, 0x26, 0xa2, 0xcf, 0xd1, 0x02, 0x3f, 0x2e, 0x9f, 0x92, 0xbf, 0xaa, 0x6f, 0xf7,
	0x75, 0xe1, 0x0d, 0x2a, 0x44, 0x76, 0x09, 0x07, 0x6b, 0x04, 0xb3, 0x66, 0xee, 0x11, 0x75, 0x6a,
	0x96, 0xb6, 0xbe, 0x41, 0xcb, 0x70, 0x23, 0xab, 0x1b, 0x94, 0xc0, 0xac, 0x67, 0x7e, 0xac, 0x4f,
	0xfe, 0x04, 0x00, 0x00, 0xff, 0xff, 0xd2, 0xb8, 0xdc, 0x9b, 0x7d, 0x07, 0x00, 0x00,
}
