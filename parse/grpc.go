package parse

import (
	"encoding/json"
	"fmt"
	"github.com/fullstorydev/grpcurl"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"net/http"
	"sort"
	"strings"
)


const (
	typeString   FieldType = "string"
	typeBytes    FieldType = "bytes"
	typeInt32    FieldType = "int32"
	typeInt64    FieldType = "int64"
	typeSint32   FieldType = "sint32"
	typeSint64   FieldType = "sint64"
	typeUint32   FieldType = "uint32"
	typeUint64   FieldType = "uint64"
	typeFixed32  FieldType = "fixed32"
	typeFixed64  FieldType = "fixed64"
	typeSfixed32 FieldType = "sfixed32"
	typeSfixed64 FieldType = "sfixed64"
	typeFloat    FieldType = "float"
	typeDouble   FieldType = "double"
	typeBool     FieldType = "bool"
	typeOneOf    FieldType = "oneof"
)

var typeMap = map[descriptor.FieldDescriptorProto_Type]FieldType{
	descriptor.FieldDescriptorProto_TYPE_STRING:   typeString,
	descriptor.FieldDescriptorProto_TYPE_BYTES:    typeBytes,
	descriptor.FieldDescriptorProto_TYPE_INT32:    typeInt32,
	descriptor.FieldDescriptorProto_TYPE_INT64:    typeInt64,
	descriptor.FieldDescriptorProto_TYPE_SINT32:   typeSint32,
	descriptor.FieldDescriptorProto_TYPE_SINT64:   typeSint64,
	descriptor.FieldDescriptorProto_TYPE_UINT32:   typeUint32,
	descriptor.FieldDescriptorProto_TYPE_UINT64:   typeUint64,
	descriptor.FieldDescriptorProto_TYPE_FIXED32:  typeFixed32,
	descriptor.FieldDescriptorProto_TYPE_FIXED64:  typeFixed64,
	descriptor.FieldDescriptorProto_TYPE_SFIXED32: typeSfixed32,
	descriptor.FieldDescriptorProto_TYPE_SFIXED64: typeSfixed64,
	descriptor.FieldDescriptorProto_TYPE_FLOAT:    typeFloat,
	descriptor.FieldDescriptorProto_TYPE_DOUBLE:   typeDouble,
	descriptor.FieldDescriptorProto_TYPE_BOOL:     typeBool,
}

type SvcConfig struct {
	includeService bool
	includeMethods map[string]struct{}
}

type EnumValDef struct {
	Num  int32  `json:"num"`
	Name string `json:"name"`
}

type Schema struct {
	RequestType   string                  `json:"requestType"`
	RequestStream bool                    `json:"requestStream"`
	MessageTypes  map[string][]FieldDef   `json:"messageTypes"`
	EnumTypes     map[string][]EnumValDef `json:"enumTypes"`
}

type FieldType string

type FieldDef struct {
	Name        string      `json:"name"`
	ProtoName   string      `json:"protoName"`
	Type        FieldType   `json:"type"`
	OneOfFields []FieldDef  `json:"oneOfFields"`
	IsMessage   bool        `json:"isMessage"`
	IsEnum      bool        `json:"isEnum"`
	IsArray     bool        `json:"isArray"`
	IsMap       bool        `json:"isMap"`
	IsRequired  bool        `json:"isRequired"`
	DefaultVal  interface{} `json:"defaultVal"`
}

func GatherMetadataForMethod(md *desc.MethodDescriptor) (*Schema, error) {
	msg := md.GetInputType()
	result := &Schema{
		RequestType:   msg.GetFullyQualifiedName(),
		RequestStream: md.IsClientStreaming(),
		MessageTypes:  map[string][]FieldDef{},
		EnumTypes:     map[string][]EnumValDef{},
	}

	result.VisitMessage(msg)

	return result, nil
}

func GetMethods(source grpcurl.DescriptorSource, configs map[string]*SvcConfig) ([]*desc.MethodDescriptor, error) {
	allServices, err := source.ListServices()
	if err != nil {
		return nil, err
	}

	var descs []*desc.MethodDescriptor
	for _, svc := range allServices {
		if svc == "grpc.reflection.v1alpha.ServerReflection" {
			continue
		}
		d, err := source.FindSymbol(svc)
		if err != nil {
			return nil, err
		}
		sd, ok := d.(*desc.ServiceDescriptor)
		if !ok {
			return nil, fmt.Errorf("%s should be a service descriptor but instead is a %T", d.GetFullyQualifiedName(), d)
		}
		cfg := configs[svc]
		if cfg == nil && len(configs) != 0 {
			// not configured to expose this service
			continue
		}
		delete(configs, svc)
		for _, md := range sd.GetMethods() {
			if cfg == nil {
				descs = append(descs, md)
				continue
			}
			_, found := cfg.includeMethods[md.GetName()]
			delete(cfg.includeMethods, md.GetName())
			if found && cfg.includeService {
				fmt.Print("Service %s already configured, so -method %s is unnecessary", svc, md.GetName())
			}
			if found || cfg.includeService {
				descs = append(descs, md)
			}
		}
		if cfg != nil && len(cfg.includeMethods) > 0 {
			// configured methods not found
			methodNames := make([]string, 0, len(cfg.includeMethods))
			for m := range cfg.includeMethods {
				methodNames = append(methodNames, fmt.Sprintf("%s/%s", svc, m))
			}
			sort.Strings(methodNames)
			return nil, fmt.Errorf("configured methods not found: %s", strings.Join(methodNames, ", "))
		}
	}

	if len(configs) > 0 {
		// configured services not found
		svcNames := make([]string, 0, len(configs))
		for s := range configs {
			svcNames = append(svcNames, s)
		}
		sort.Strings(svcNames)
		return nil, fmt.Errorf("configured services not found: %s", strings.Join(svcNames, ", "))
	}

	return descs, nil
}


func (s *Schema) VisitMessage(md *desc.MessageDescriptor) {
	if _, ok := s.MessageTypes[md.GetFullyQualifiedName()]; ok {
		// already visited
		return
	}

	fields := make([]FieldDef, 0, len(md.GetFields()))
	s.MessageTypes[md.GetFullyQualifiedName()] = fields

	oneOfsSeen := map[*desc.OneOfDescriptor]struct{}{}
	for _, fd := range md.GetFields() {
		ood := fd.GetOneOf()
		if ood != nil {
			if _, ok := oneOfsSeen[ood]; ok {
				// already processed this one
				continue
			}
			oneOfsSeen[ood] = struct{}{}
			fields = append(fields, s.ProcessOneOf(ood))
		} else {
			fields = append(fields, s.ProcessField(fd))
		}
	}

	s.MessageTypes[md.GetFullyQualifiedName()] = fields
}

func (s *Schema) ProcessOneOf(ood *desc.OneOfDescriptor) FieldDef {
	choices := make([]FieldDef, len(ood.GetChoices()))
	for i, fd := range ood.GetChoices() {
		choices[i] = s.ProcessField(fd)
	}
	return FieldDef{
		Name:        ood.GetName(),
		Type:        typeOneOf,
		OneOfFields: choices,
	}
}

func (s *Schema) ProcessField(fd *desc.FieldDescriptor) FieldDef {
	def := FieldDef{
		Name:       fd.GetJSONName(),
		ProtoName:  fd.GetName(),
		IsEnum:     fd.GetEnumType() != nil,
		IsMessage:  fd.GetMessageType() != nil,
		IsArray:    fd.IsRepeated() && !fd.IsMap(),
		IsMap:      fd.IsMap(),
		IsRequired: fd.IsRequired(),
		DefaultVal: fd.GetDefaultValue(),
	}

	if def.IsMap {
		// fd.GetDefaultValue returns empty map[interface{}]interface{}
		// as the default for map fields, but "encoding/json" refuses
		// to encode a map with interface{} keys (even if it's empty).
		// So we fix up the key type here.
		def.DefaultVal = map[string]interface{}{}
	}

	// 64-bit int values are represented as strings in JSON
	if i, ok := def.DefaultVal.(int64); ok {
		def.DefaultVal = fmt.Sprintf("%d", i)
	} else if u, ok := def.DefaultVal.(uint64); ok {
		def.DefaultVal = fmt.Sprintf("%d", u)
	} else if b, ok := def.DefaultVal.([]byte); ok && b == nil {
		// bytes fields may have []byte(nil) as default value, but
		// that gets rendered as JSON null, not empty array
		def.DefaultVal = []byte{}
	}

	switch fd.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		def.Type = FieldType(fd.GetEnumType().GetFullyQualifiedName())
		s.VisitEnum(fd.GetEnumType())
		// DefaultVal will be int32 for enums, but we want to instead
		// send enum name as string
		if val, ok := def.DefaultVal.(int32); ok {
			valDesc := fd.GetEnumType().FindValueByNumber(val)
			if valDesc != nil {
				def.DefaultVal = valDesc.GetName()
			}
		}

	case descriptor.FieldDescriptorProto_TYPE_GROUP, descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		def.Type = FieldType(fd.GetMessageType().GetFullyQualifiedName())
		s.VisitMessage(fd.GetMessageType())

	default:
		def.Type = typeMap[fd.GetType()]
	}

	return def
}

func GatherAllMessageMetadata(files []*desc.FileDescriptor) *Schema {
	result := &Schema{
		MessageTypes: map[string][]FieldDef{},
		EnumTypes:    map[string][]EnumValDef{},
	}
	for _, fd := range files {
		gatherAllMessages(fd.GetMessageTypes(), result)
	}
	return result
}


func gatherAllMessages(msgs []*desc.MessageDescriptor, result *Schema) {
	for _, md := range msgs {
		result.VisitMessage(md)
		gatherAllMessages(md.GetNestedMessageTypes(), result)
	}
}


func (s *Schema) VisitEnum(ed *desc.EnumDescriptor) {
	if _, ok := s.EnumTypes[ed.GetFullyQualifiedName()]; ok {
		// already visited
		return
	}

	enumValDefs := make([]EnumValDef, len(ed.GetValues()))
	for i, evd := range ed.GetValues() {
		enumValDefs[i] = EnumValDef{
			Num:  evd.GetNumber(),
			Name: evd.GetName(),
		}
	}

	s.EnumTypes[ed.GetFullyQualifiedName()] = enumValDefs
}


func RPCMetadataHandler(methods []*desc.MethodDescriptor, files []*desc.FileDescriptor) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.Header().Set("Allow", "GET")
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		method := r.URL.Query().Get("method")
		var results *Schema
		if method == "*" {
			// This means gather *all* message types. This is used to
			// provide a drop-down for Any messages.
			results = GatherAllMessageMetadata(files)
		} else {
			for _, md := range methods {
				if md.GetFullyQualifiedName() == method {
					r, err := GatherMetadataForMethod(md)
					if err != nil {
						http.Error(w, "Failed to gather metadata for RPC Method", http.StatusUnprocessableEntity)
						return
					}
					fmt.Println(r)
					results = r
					break
				}
			}
		}

		if results == nil {
			http.Error(w, "Unknown RPC Method", http.StatusUnprocessableEntity)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)

		enc.SetIndent("", "  ")
		// TODO: what if enc.Encode returns a non-I/O error?
		enc.Encode(results)
	})
}



