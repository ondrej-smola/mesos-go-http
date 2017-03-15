// Code generated by protoc-gen-gogo.
// source: operator/quota/quota.proto
// DO NOT EDIT!

/*
	Package quota is a generated protocol buffer package.

	It is generated from these files:
		operator/quota/quota.proto

	It has these top-level messages:
		QuotaInfo
		QuotaRequest
		QuotaStatus
*/
package quota

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import mesos_v1 "github.com/ondrej-smola/mesos-go-http"
import _ "github.com/gogo/protobuf/gogoproto"

import strings "strings"
import reflect "reflect"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// TODO(joerg84): Add limits, i.e. upper bounds of resources that a
// role is allowed to use.
type QuotaInfo struct {
	// Quota is granted per role and not per framework, similar to
	// dynamic reservations.
	Role *string `protobuf:"bytes,1,opt,name=role" json:"role,omitempty"`
	// Principal which set the quota. Currently only operators can set quotas.
	Principal *string `protobuf:"bytes,2,opt,name=principal" json:"principal,omitempty"`
	// The guarantee that these resources are allocatable for the above role.
	// NOTE: `guarantee.role` should not specify any role except '*',
	// because quota does not reserve specific resources.
	Guarantee []mesos_v1.Resource `protobuf:"bytes,3,rep,name=guarantee" json:"guarantee"`
}

func (m *QuotaInfo) Reset()                    { *m = QuotaInfo{} }
func (*QuotaInfo) ProtoMessage()               {}
func (*QuotaInfo) Descriptor() ([]byte, []int) { return fileDescriptorQuota, []int{0} }

func (m *QuotaInfo) GetRole() string {
	if m != nil && m.Role != nil {
		return *m.Role
	}
	return ""
}

func (m *QuotaInfo) GetPrincipal() string {
	if m != nil && m.Principal != nil {
		return *m.Principal
	}
	return ""
}

func (m *QuotaInfo) GetGuarantee() []mesos_v1.Resource {
	if m != nil {
		return m.Guarantee
	}
	return nil
}

// *
// `QuotaRequest` provides a schema for set quota JSON requests.
type QuotaRequest struct {
	// Disables the capacity heuristic check if set to `true`.
	Force *bool `protobuf:"varint,1,opt,name=force,def=0" json:"force,omitempty"`
	// The role for which to set quota.
	Role *string `protobuf:"bytes,2,opt,name=role" json:"role,omitempty"`
	// The requested guarantee that these resources will be allocatable for
	// the above role.
	Guarantee []mesos_v1.Resource `protobuf:"bytes,3,rep,name=guarantee" json:"guarantee"`
}

func (m *QuotaRequest) Reset()                    { *m = QuotaRequest{} }
func (*QuotaRequest) ProtoMessage()               {}
func (*QuotaRequest) Descriptor() ([]byte, []int) { return fileDescriptorQuota, []int{1} }

const Default_QuotaRequest_Force bool = false

func (m *QuotaRequest) GetForce() bool {
	if m != nil && m.Force != nil {
		return *m.Force
	}
	return Default_QuotaRequest_Force
}

func (m *QuotaRequest) GetRole() string {
	if m != nil && m.Role != nil {
		return *m.Role
	}
	return ""
}

func (m *QuotaRequest) GetGuarantee() []mesos_v1.Resource {
	if m != nil {
		return m.Guarantee
	}
	return nil
}

// *
// `QuotaStatus` describes the internal representation for the /quota/status
// response.
type QuotaStatus struct {
	// Quotas which are currently set, i.e. known to the master.
	Infos []QuotaInfo `protobuf:"bytes,1,rep,name=infos" json:"infos"`
}

func (m *QuotaStatus) Reset()                    { *m = QuotaStatus{} }
func (*QuotaStatus) ProtoMessage()               {}
func (*QuotaStatus) Descriptor() ([]byte, []int) { return fileDescriptorQuota, []int{2} }

func (m *QuotaStatus) GetInfos() []QuotaInfo {
	if m != nil {
		return m.Infos
	}
	return nil
}

func init() {
	proto.RegisterType((*QuotaInfo)(nil), "mesos.v1.quota.QuotaInfo")
	proto.RegisterType((*QuotaRequest)(nil), "mesos.v1.quota.QuotaRequest")
	proto.RegisterType((*QuotaStatus)(nil), "mesos.v1.quota.QuotaStatus")
}
func (this *QuotaInfo) VerboseEqual(that interface{}) error {
	if that == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that == nil && this != nil")
	}

	that1, ok := that.(*QuotaInfo)
	if !ok {
		that2, ok := that.(QuotaInfo)
		if ok {
			that1 = &that2
		} else {
			return fmt.Errorf("that is not of type *QuotaInfo")
		}
	}
	if that1 == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that is type *QuotaInfo but is nil && this != nil")
	} else if this == nil {
		return fmt.Errorf("that is type *QuotaInfo but is not nil && this == nil")
	}
	if this.Role != nil && that1.Role != nil {
		if *this.Role != *that1.Role {
			return fmt.Errorf("Role this(%v) Not Equal that(%v)", *this.Role, *that1.Role)
		}
	} else if this.Role != nil {
		return fmt.Errorf("this.Role == nil && that.Role != nil")
	} else if that1.Role != nil {
		return fmt.Errorf("Role this(%v) Not Equal that(%v)", this.Role, that1.Role)
	}
	if this.Principal != nil && that1.Principal != nil {
		if *this.Principal != *that1.Principal {
			return fmt.Errorf("Principal this(%v) Not Equal that(%v)", *this.Principal, *that1.Principal)
		}
	} else if this.Principal != nil {
		return fmt.Errorf("this.Principal == nil && that.Principal != nil")
	} else if that1.Principal != nil {
		return fmt.Errorf("Principal this(%v) Not Equal that(%v)", this.Principal, that1.Principal)
	}
	if len(this.Guarantee) != len(that1.Guarantee) {
		return fmt.Errorf("Guarantee this(%v) Not Equal that(%v)", len(this.Guarantee), len(that1.Guarantee))
	}
	for i := range this.Guarantee {
		if !this.Guarantee[i].Equal(&that1.Guarantee[i]) {
			return fmt.Errorf("Guarantee this[%v](%v) Not Equal that[%v](%v)", i, this.Guarantee[i], i, that1.Guarantee[i])
		}
	}
	return nil
}
func (this *QuotaInfo) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*QuotaInfo)
	if !ok {
		that2, ok := that.(QuotaInfo)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if this.Role != nil && that1.Role != nil {
		if *this.Role != *that1.Role {
			return false
		}
	} else if this.Role != nil {
		return false
	} else if that1.Role != nil {
		return false
	}
	if this.Principal != nil && that1.Principal != nil {
		if *this.Principal != *that1.Principal {
			return false
		}
	} else if this.Principal != nil {
		return false
	} else if that1.Principal != nil {
		return false
	}
	if len(this.Guarantee) != len(that1.Guarantee) {
		return false
	}
	for i := range this.Guarantee {
		if !this.Guarantee[i].Equal(&that1.Guarantee[i]) {
			return false
		}
	}
	return true
}
func (this *QuotaRequest) VerboseEqual(that interface{}) error {
	if that == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that == nil && this != nil")
	}

	that1, ok := that.(*QuotaRequest)
	if !ok {
		that2, ok := that.(QuotaRequest)
		if ok {
			that1 = &that2
		} else {
			return fmt.Errorf("that is not of type *QuotaRequest")
		}
	}
	if that1 == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that is type *QuotaRequest but is nil && this != nil")
	} else if this == nil {
		return fmt.Errorf("that is type *QuotaRequest but is not nil && this == nil")
	}
	if this.Force != nil && that1.Force != nil {
		if *this.Force != *that1.Force {
			return fmt.Errorf("Force this(%v) Not Equal that(%v)", *this.Force, *that1.Force)
		}
	} else if this.Force != nil {
		return fmt.Errorf("this.Force == nil && that.Force != nil")
	} else if that1.Force != nil {
		return fmt.Errorf("Force this(%v) Not Equal that(%v)", this.Force, that1.Force)
	}
	if this.Role != nil && that1.Role != nil {
		if *this.Role != *that1.Role {
			return fmt.Errorf("Role this(%v) Not Equal that(%v)", *this.Role, *that1.Role)
		}
	} else if this.Role != nil {
		return fmt.Errorf("this.Role == nil && that.Role != nil")
	} else if that1.Role != nil {
		return fmt.Errorf("Role this(%v) Not Equal that(%v)", this.Role, that1.Role)
	}
	if len(this.Guarantee) != len(that1.Guarantee) {
		return fmt.Errorf("Guarantee this(%v) Not Equal that(%v)", len(this.Guarantee), len(that1.Guarantee))
	}
	for i := range this.Guarantee {
		if !this.Guarantee[i].Equal(&that1.Guarantee[i]) {
			return fmt.Errorf("Guarantee this[%v](%v) Not Equal that[%v](%v)", i, this.Guarantee[i], i, that1.Guarantee[i])
		}
	}
	return nil
}
func (this *QuotaRequest) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*QuotaRequest)
	if !ok {
		that2, ok := that.(QuotaRequest)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if this.Force != nil && that1.Force != nil {
		if *this.Force != *that1.Force {
			return false
		}
	} else if this.Force != nil {
		return false
	} else if that1.Force != nil {
		return false
	}
	if this.Role != nil && that1.Role != nil {
		if *this.Role != *that1.Role {
			return false
		}
	} else if this.Role != nil {
		return false
	} else if that1.Role != nil {
		return false
	}
	if len(this.Guarantee) != len(that1.Guarantee) {
		return false
	}
	for i := range this.Guarantee {
		if !this.Guarantee[i].Equal(&that1.Guarantee[i]) {
			return false
		}
	}
	return true
}
func (this *QuotaStatus) VerboseEqual(that interface{}) error {
	if that == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that == nil && this != nil")
	}

	that1, ok := that.(*QuotaStatus)
	if !ok {
		that2, ok := that.(QuotaStatus)
		if ok {
			that1 = &that2
		} else {
			return fmt.Errorf("that is not of type *QuotaStatus")
		}
	}
	if that1 == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that is type *QuotaStatus but is nil && this != nil")
	} else if this == nil {
		return fmt.Errorf("that is type *QuotaStatus but is not nil && this == nil")
	}
	if len(this.Infos) != len(that1.Infos) {
		return fmt.Errorf("Infos this(%v) Not Equal that(%v)", len(this.Infos), len(that1.Infos))
	}
	for i := range this.Infos {
		if !this.Infos[i].Equal(&that1.Infos[i]) {
			return fmt.Errorf("Infos this[%v](%v) Not Equal that[%v](%v)", i, this.Infos[i], i, that1.Infos[i])
		}
	}
	return nil
}
func (this *QuotaStatus) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*QuotaStatus)
	if !ok {
		that2, ok := that.(QuotaStatus)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if len(this.Infos) != len(that1.Infos) {
		return false
	}
	for i := range this.Infos {
		if !this.Infos[i].Equal(&that1.Infos[i]) {
			return false
		}
	}
	return true
}
func (this *QuotaInfo) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&quota.QuotaInfo{")
	if this.Role != nil {
		s = append(s, "Role: "+valueToGoStringQuota(this.Role, "string")+",\n")
	}
	if this.Principal != nil {
		s = append(s, "Principal: "+valueToGoStringQuota(this.Principal, "string")+",\n")
	}
	if this.Guarantee != nil {
		s = append(s, "Guarantee: "+fmt.Sprintf("%#v", this.Guarantee)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *QuotaRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&quota.QuotaRequest{")
	if this.Force != nil {
		s = append(s, "Force: "+valueToGoStringQuota(this.Force, "bool")+",\n")
	}
	if this.Role != nil {
		s = append(s, "Role: "+valueToGoStringQuota(this.Role, "string")+",\n")
	}
	if this.Guarantee != nil {
		s = append(s, "Guarantee: "+fmt.Sprintf("%#v", this.Guarantee)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *QuotaStatus) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 5)
	s = append(s, "&quota.QuotaStatus{")
	if this.Infos != nil {
		s = append(s, "Infos: "+fmt.Sprintf("%#v", this.Infos)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringQuota(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *QuotaInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QuotaInfo) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Role != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintQuota(dAtA, i, uint64(len(*m.Role)))
		i += copy(dAtA[i:], *m.Role)
	}
	if m.Principal != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintQuota(dAtA, i, uint64(len(*m.Principal)))
		i += copy(dAtA[i:], *m.Principal)
	}
	if len(m.Guarantee) > 0 {
		for _, msg := range m.Guarantee {
			dAtA[i] = 0x1a
			i++
			i = encodeVarintQuota(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *QuotaRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QuotaRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Force != nil {
		dAtA[i] = 0x8
		i++
		if *m.Force {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if m.Role != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintQuota(dAtA, i, uint64(len(*m.Role)))
		i += copy(dAtA[i:], *m.Role)
	}
	if len(m.Guarantee) > 0 {
		for _, msg := range m.Guarantee {
			dAtA[i] = 0x1a
			i++
			i = encodeVarintQuota(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *QuotaStatus) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QuotaStatus) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Infos) > 0 {
		for _, msg := range m.Infos {
			dAtA[i] = 0xa
			i++
			i = encodeVarintQuota(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func encodeFixed64Quota(dAtA []byte, offset int, v uint64) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	dAtA[offset+4] = uint8(v >> 32)
	dAtA[offset+5] = uint8(v >> 40)
	dAtA[offset+6] = uint8(v >> 48)
	dAtA[offset+7] = uint8(v >> 56)
	return offset + 8
}
func encodeFixed32Quota(dAtA []byte, offset int, v uint32) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintQuota(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *QuotaInfo) Size() (n int) {
	var l int
	_ = l
	if m.Role != nil {
		l = len(*m.Role)
		n += 1 + l + sovQuota(uint64(l))
	}
	if m.Principal != nil {
		l = len(*m.Principal)
		n += 1 + l + sovQuota(uint64(l))
	}
	if len(m.Guarantee) > 0 {
		for _, e := range m.Guarantee {
			l = e.Size()
			n += 1 + l + sovQuota(uint64(l))
		}
	}
	return n
}

func (m *QuotaRequest) Size() (n int) {
	var l int
	_ = l
	if m.Force != nil {
		n += 2
	}
	if m.Role != nil {
		l = len(*m.Role)
		n += 1 + l + sovQuota(uint64(l))
	}
	if len(m.Guarantee) > 0 {
		for _, e := range m.Guarantee {
			l = e.Size()
			n += 1 + l + sovQuota(uint64(l))
		}
	}
	return n
}

func (m *QuotaStatus) Size() (n int) {
	var l int
	_ = l
	if len(m.Infos) > 0 {
		for _, e := range m.Infos {
			l = e.Size()
			n += 1 + l + sovQuota(uint64(l))
		}
	}
	return n
}

func sovQuota(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozQuota(x uint64) (n int) {
	return sovQuota(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *QuotaInfo) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&QuotaInfo{`,
		`Role:` + valueToStringQuota(this.Role) + `,`,
		`Principal:` + valueToStringQuota(this.Principal) + `,`,
		`Guarantee:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.Guarantee), "Resource", "mesos_v1.Resource", 1), `&`, ``, 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *QuotaRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&QuotaRequest{`,
		`Force:` + valueToStringQuota(this.Force) + `,`,
		`Role:` + valueToStringQuota(this.Role) + `,`,
		`Guarantee:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.Guarantee), "Resource", "mesos_v1.Resource", 1), `&`, ``, 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *QuotaStatus) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&QuotaStatus{`,
		`Infos:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.Infos), "QuotaInfo", "QuotaInfo", 1), `&`, ``, 1) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringQuota(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *QuotaInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuota
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QuotaInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QuotaInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Role", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuota
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuota
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(dAtA[iNdEx:postIndex])
			m.Role = &s
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Principal", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuota
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuota
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(dAtA[iNdEx:postIndex])
			m.Principal = &s
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Guarantee", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuota
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuota
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Guarantee = append(m.Guarantee, mesos_v1.Resource{})
			if err := m.Guarantee[len(m.Guarantee)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuota(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthQuota
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QuotaRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuota
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QuotaRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QuotaRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Force", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuota
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			b := bool(v != 0)
			m.Force = &b
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Role", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuota
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuota
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(dAtA[iNdEx:postIndex])
			m.Role = &s
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Guarantee", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuota
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuota
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Guarantee = append(m.Guarantee, mesos_v1.Resource{})
			if err := m.Guarantee[len(m.Guarantee)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuota(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthQuota
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QuotaStatus) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuota
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QuotaStatus: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QuotaStatus: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Infos", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuota
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuota
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Infos = append(m.Infos, QuotaInfo{})
			if err := m.Infos[len(m.Infos)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuota(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthQuota
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipQuota(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuota
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuota
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuota
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthQuota
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowQuota
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipQuota(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthQuota = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuota   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("operator/quota/quota.proto", fileDescriptorQuota) }

var fileDescriptorQuota = []byte{
	// 312 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x90, 0xbd, 0x4e, 0xc3, 0x30,
	0x14, 0x85, 0x63, 0xda, 0x48, 0xc4, 0xad, 0x90, 0xb0, 0x18, 0x42, 0x07, 0x53, 0x75, 0xaa, 0x84,
	0xe2, 0xa8, 0xb0, 0x21, 0xb1, 0x74, 0x63, 0xa4, 0x4c, 0xb0, 0x20, 0x37, 0x38, 0x69, 0x50, 0x9a,
	0x9b, 0xfa, 0x87, 0x99, 0x47, 0xe0, 0x31, 0x78, 0x94, 0x8e, 0x8c, 0x4c, 0xa8, 0x35, 0x0b, 0x63,
	0x1f, 0x01, 0x61, 0x4b, 0x11, 0xac, 0x2c, 0x96, 0xcf, 0xd1, 0x3d, 0xe7, 0xb3, 0x2f, 0x1e, 0x40,
	0x23, 0x24, 0xd7, 0x20, 0xd3, 0x95, 0x01, 0xcd, 0xfd, 0xc9, 0x1a, 0x09, 0x1a, 0xc8, 0xc1, 0x52,
	0x28, 0x50, 0xec, 0x69, 0xc2, 0x9c, 0x3b, 0x98, 0x14, 0xa5, 0x5e, 0x98, 0x39, 0xcb, 0x60, 0x99,
	0x42, 0xfd, 0x20, 0xc5, 0x63, 0xa2, 0x96, 0x50, 0xf1, 0xd4, 0xcd, 0x25, 0x05, 0x24, 0x0b, 0xad,
	0x1b, 0xaf, 0x7c, 0xc5, 0x20, 0xf9, 0x15, 0x29, 0xa0, 0x80, 0xd4, 0xd9, 0x73, 0x93, 0x3b, 0xe5,
	0x84, 0xbb, 0xf9, 0xf1, 0xd1, 0x2d, 0x8e, 0xae, 0x7f, 0x50, 0x57, 0x75, 0x0e, 0xa4, 0x8f, 0xbb,
	0x12, 0x2a, 0x11, 0xa3, 0x21, 0x1a, 0x47, 0xe4, 0x10, 0x47, 0x8d, 0x2c, 0xeb, 0xac, 0x6c, 0x78,
	0x15, 0xef, 0x39, 0xeb, 0x14, 0x47, 0x85, 0xe1, 0x92, 0xd7, 0x5a, 0x88, 0xb8, 0x33, 0xec, 0x8c,
	0x7b, 0x67, 0x84, 0xb5, 0x6f, 0x9e, 0x09, 0x05, 0x46, 0x66, 0x62, 0xda, 0x5d, 0x7f, 0x9c, 0x04,
	0xa3, 0x7b, 0xdc, 0x77, 0xd5, 0x33, 0xb1, 0x32, 0x42, 0x69, 0x72, 0x84, 0xc3, 0x1c, 0x64, 0xe6,
	0xeb, 0xf7, 0x2f, 0xc2, 0x9c, 0x57, 0x4a, 0xb4, 0xcc, 0x7f, 0x00, 0x2e, 0x71, 0xcf, 0x01, 0x6e,
	0x34, 0xd7, 0x46, 0x11, 0x86, 0xc3, 0xb2, 0xce, 0x41, 0xc5, 0xc8, 0xe5, 0x8e, 0xd9, 0xdf, 0x65,
	0xb2, 0xf6, 0x9f, 0x3e, 0x3e, 0x3d, 0x7f, 0xdf, 0xd2, 0x60, 0xb3, 0xa5, 0x68, 0xb7, 0xa5, 0xe8,
	0xd9, 0x52, 0xf4, 0x6a, 0x29, 0x5a, 0x5b, 0x8a, 0xde, 0x2c, 0x45, 0x1b, 0x4b, 0xd1, 0x97, 0xa5,
	0xc1, 0xce, 0x52, 0xf4, 0xf2, 0x49, 0x83, 0xbb, 0xd0, 0x95, 0x7c, 0x07, 0x00, 0x00, 0xff, 0xff,
	0xbc, 0x77, 0xa5, 0x5c, 0xbe, 0x01, 0x00, 0x00,
}
