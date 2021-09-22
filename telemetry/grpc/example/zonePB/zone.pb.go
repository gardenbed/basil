// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: zone.proto

package zonePB

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Location struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Latitude  float64 `protobuf:"fixed64,1,opt,name=latitude,proto3" json:"latitude,omitempty"`
	Longitude float64 `protobuf:"fixed64,2,opt,name=longitude,proto3" json:"longitude,omitempty"`
}

func (x *Location) Reset() {
	*x = Location{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zone_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Location) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Location) ProtoMessage() {}

func (x *Location) ProtoReflect() protoreflect.Message {
	mi := &file_zone_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Location.ProtoReflect.Descriptor instead.
func (*Location) Descriptor() ([]byte, []int) {
	return file_zone_proto_rawDescGZIP(), []int{0}
}

func (x *Location) GetLatitude() float64 {
	if x != nil {
		return x.Latitude
	}
	return 0
}

func (x *Location) GetLongitude() float64 {
	if x != nil {
		return x.Longitude
	}
	return 0
}

type Zone struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Location *Location `protobuf:"bytes,1,opt,name=location,proto3" json:"location,omitempty"`
	Radius   float64   `protobuf:"fixed64,2,opt,name=radius,proto3" json:"radius,omitempty"`
}

func (x *Zone) Reset() {
	*x = Zone{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zone_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Zone) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Zone) ProtoMessage() {}

func (x *Zone) ProtoReflect() protoreflect.Message {
	mi := &file_zone_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Zone.ProtoReflect.Descriptor instead.
func (*Zone) Descriptor() ([]byte, []int) {
	return file_zone_proto_rawDescGZIP(), []int{1}
}

func (x *Zone) GetLocation() *Location {
	if x != nil {
		return x.Location
	}
	return nil
}

func (x *Zone) GetRadius() float64 {
	if x != nil {
		return x.Radius
	}
	return 0
}

type Place struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       string    `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name     string    `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Location *Location `protobuf:"bytes,3,opt,name=location,proto3" json:"location,omitempty"`
}

func (x *Place) Reset() {
	*x = Place{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zone_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Place) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Place) ProtoMessage() {}

func (x *Place) ProtoReflect() protoreflect.Message {
	mi := &file_zone_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Place.ProtoReflect.Descriptor instead.
func (*Place) Descriptor() ([]byte, []int) {
	return file_zone_proto_rawDescGZIP(), []int{2}
}

func (x *Place) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Place) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Place) GetLocation() *Location {
	if x != nil {
		return x.Location
	}
	return nil
}

type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zone_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_zone_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_zone_proto_rawDescGZIP(), []int{3}
}

func (x *User) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *User) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type UserInZone struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Location *Location `protobuf:"bytes,1,opt,name=location,proto3" json:"location,omitempty"`
	User     *User     `protobuf:"bytes,2,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *UserInZone) Reset() {
	*x = UserInZone{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zone_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserInZone) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserInZone) ProtoMessage() {}

func (x *UserInZone) ProtoReflect() protoreflect.Message {
	mi := &file_zone_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserInZone.ProtoReflect.Descriptor instead.
func (*UserInZone) Descriptor() ([]byte, []int) {
	return file_zone_proto_rawDescGZIP(), []int{4}
}

func (x *UserInZone) GetLocation() *Location {
	if x != nil {
		return x.Location
	}
	return nil
}

func (x *UserInZone) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

type GetPlacesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Zone   *Zone    `protobuf:"bytes,1,opt,name=zone,proto3" json:"zone,omitempty"`
	Places []*Place `protobuf:"bytes,2,rep,name=places,proto3" json:"places,omitempty"`
}

func (x *GetPlacesResponse) Reset() {
	*x = GetPlacesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zone_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPlacesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPlacesResponse) ProtoMessage() {}

func (x *GetPlacesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_zone_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPlacesResponse.ProtoReflect.Descriptor instead.
func (*GetPlacesResponse) Descriptor() ([]byte, []int) {
	return file_zone_proto_rawDescGZIP(), []int{5}
}

func (x *GetPlacesResponse) GetZone() *Zone {
	if x != nil {
		return x.Zone
	}
	return nil
}

func (x *GetPlacesResponse) GetPlaces() []*Place {
	if x != nil {
		return x.Places
	}
	return nil
}

var File_zone_proto protoreflect.FileDescriptor

var file_zone_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x7a, 0x6f, 0x6e, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x7a, 0x6f,
	0x6e, 0x65, 0x50, 0x42, 0x22, 0x44, 0x0a, 0x08, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x08, 0x6c, 0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x12, 0x1c, 0x0a, 0x09,
	0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x22, 0x4c, 0x0a, 0x04, 0x5a, 0x6f,
	0x6e, 0x65, 0x12, 0x2c, 0x0a, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x7a, 0x6f, 0x6e, 0x65, 0x50, 0x42, 0x2e, 0x4c, 0x6f,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x16, 0x0a, 0x06, 0x72, 0x61, 0x64, 0x69, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x06, 0x72, 0x61, 0x64, 0x69, 0x75, 0x73, 0x22, 0x59, 0x0a, 0x05, 0x50, 0x6c, 0x61, 0x63,
	0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x2c, 0x0a, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x7a, 0x6f, 0x6e, 0x65, 0x50, 0x42,
	0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x22, 0x2a, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22,
	0x5c, 0x0a, 0x0a, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x5a, 0x6f, 0x6e, 0x65, 0x12, 0x2c, 0x0a,
	0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x10, 0x2e, 0x7a, 0x6f, 0x6e, 0x65, 0x50, 0x42, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x20, 0x0a, 0x04, 0x75,
	0x73, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x7a, 0x6f, 0x6e, 0x65,
	0x50, 0x42, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x22, 0x5c, 0x0a,
	0x11, 0x47, 0x65, 0x74, 0x50, 0x6c, 0x61, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x20, 0x0a, 0x04, 0x7a, 0x6f, 0x6e, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0c, 0x2e, 0x7a, 0x6f, 0x6e, 0x65, 0x50, 0x42, 0x2e, 0x5a, 0x6f, 0x6e, 0x65, 0x52, 0x04,
	0x7a, 0x6f, 0x6e, 0x65, 0x12, 0x25, 0x0a, 0x06, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x7a, 0x6f, 0x6e, 0x65, 0x50, 0x42, 0x2e, 0x50, 0x6c,
	0x61, 0x63, 0x65, 0x52, 0x06, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x73, 0x32, 0xef, 0x01, 0x0a, 0x0b,
	0x5a, 0x6f, 0x6e, 0x65, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x12, 0x35, 0x0a, 0x11, 0x47,
	0x65, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x69, 0x6e, 0x67, 0x5a, 0x6f, 0x6e, 0x65,
	0x12, 0x10, 0x2e, 0x7a, 0x6f, 0x6e, 0x65, 0x50, 0x42, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x1a, 0x0c, 0x2e, 0x7a, 0x6f, 0x6e, 0x65, 0x50, 0x42, 0x2e, 0x5a, 0x6f, 0x6e, 0x65,
	0x28, 0x01, 0x12, 0x3a, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x50, 0x6c, 0x61, 0x63, 0x65, 0x73, 0x49,
	0x6e, 0x5a, 0x6f, 0x6e, 0x65, 0x12, 0x0c, 0x2e, 0x7a, 0x6f, 0x6e, 0x65, 0x50, 0x42, 0x2e, 0x5a,
	0x6f, 0x6e, 0x65, 0x1a, 0x19, 0x2e, 0x7a, 0x6f, 0x6e, 0x65, 0x50, 0x42, 0x2e, 0x47, 0x65, 0x74,
	0x50, 0x6c, 0x61, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x34,
	0x0a, 0x0e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x73, 0x49, 0x6e, 0x5a, 0x6f, 0x6e, 0x65,
	0x12, 0x0c, 0x2e, 0x7a, 0x6f, 0x6e, 0x65, 0x50, 0x42, 0x2e, 0x5a, 0x6f, 0x6e, 0x65, 0x1a, 0x12,
	0x2e, 0x7a, 0x6f, 0x6e, 0x65, 0x50, 0x42, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x5a, 0x6f,
	0x6e, 0x65, 0x30, 0x01, 0x12, 0x37, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x73,
	0x49, 0x6e, 0x5a, 0x6f, 0x6e, 0x65, 0x73, 0x12, 0x0c, 0x2e, 0x7a, 0x6f, 0x6e, 0x65, 0x50, 0x42,
	0x2e, 0x5a, 0x6f, 0x6e, 0x65, 0x1a, 0x12, 0x2e, 0x7a, 0x6f, 0x6e, 0x65, 0x50, 0x42, 0x2e, 0x55,
	0x73, 0x65, 0x72, 0x49, 0x6e, 0x5a, 0x6f, 0x6e, 0x65, 0x28, 0x01, 0x30, 0x01, 0x42, 0x3a, 0x5a,
	0x38, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x61, 0x72, 0x64,
	0x65, 0x6e, 0x62, 0x65, 0x64, 0x2f, 0x62, 0x61, 0x73, 0x69, 0x6c, 0x2f, 0x74, 0x65, 0x6c, 0x65,
	0x6d, 0x65, 0x74, 0x72, 0x79, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70,
	0x6c, 0x65, 0x2f, 0x7a, 0x6f, 0x6e, 0x65, 0x50, 0x42, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_zone_proto_rawDescOnce sync.Once
	file_zone_proto_rawDescData = file_zone_proto_rawDesc
)

func file_zone_proto_rawDescGZIP() []byte {
	file_zone_proto_rawDescOnce.Do(func() {
		file_zone_proto_rawDescData = protoimpl.X.CompressGZIP(file_zone_proto_rawDescData)
	})
	return file_zone_proto_rawDescData
}

var file_zone_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_zone_proto_goTypes = []interface{}{
	(*Location)(nil),          // 0: zonePB.Location
	(*Zone)(nil),              // 1: zonePB.Zone
	(*Place)(nil),             // 2: zonePB.Place
	(*User)(nil),              // 3: zonePB.User
	(*UserInZone)(nil),        // 4: zonePB.UserInZone
	(*GetPlacesResponse)(nil), // 5: zonePB.GetPlacesResponse
}
var file_zone_proto_depIdxs = []int32{
	0,  // 0: zonePB.Zone.location:type_name -> zonePB.Location
	0,  // 1: zonePB.Place.location:type_name -> zonePB.Location
	0,  // 2: zonePB.UserInZone.location:type_name -> zonePB.Location
	3,  // 3: zonePB.UserInZone.user:type_name -> zonePB.User
	1,  // 4: zonePB.GetPlacesResponse.zone:type_name -> zonePB.Zone
	2,  // 5: zonePB.GetPlacesResponse.places:type_name -> zonePB.Place
	0,  // 6: zonePB.ZoneManager.GetContainingZone:input_type -> zonePB.Location
	1,  // 7: zonePB.ZoneManager.GetPlacesInZone:input_type -> zonePB.Zone
	1,  // 8: zonePB.ZoneManager.GetUsersInZone:input_type -> zonePB.Zone
	1,  // 9: zonePB.ZoneManager.GetUsersInZones:input_type -> zonePB.Zone
	1,  // 10: zonePB.ZoneManager.GetContainingZone:output_type -> zonePB.Zone
	5,  // 11: zonePB.ZoneManager.GetPlacesInZone:output_type -> zonePB.GetPlacesResponse
	4,  // 12: zonePB.ZoneManager.GetUsersInZone:output_type -> zonePB.UserInZone
	4,  // 13: zonePB.ZoneManager.GetUsersInZones:output_type -> zonePB.UserInZone
	10, // [10:14] is the sub-list for method output_type
	6,  // [6:10] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_zone_proto_init() }
func file_zone_proto_init() {
	if File_zone_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_zone_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Location); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_zone_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Zone); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_zone_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Place); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_zone_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*User); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_zone_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserInZone); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_zone_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPlacesResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_zone_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_zone_proto_goTypes,
		DependencyIndexes: file_zone_proto_depIdxs,
		MessageInfos:      file_zone_proto_msgTypes,
	}.Build()
	File_zone_proto = out.File
	file_zone_proto_rawDesc = nil
	file_zone_proto_goTypes = nil
	file_zone_proto_depIdxs = nil
}
