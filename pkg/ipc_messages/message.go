package messages

//
//import (
//  "github.com/golang/protobuf/ptypes/any"
//  "google.golang.org/protobuf/runtime/protoimpl"
//)
//
//// BaseMessage The request message containing the user's name.
//type BaseMessage struct {
//	MessageType message.MessageType `protobuf:"varint,1,opt,name=messageType,proto3,enum=MessageType" json:"messageType,omitempty"`
//	Data        []byte              `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
//}
//
//type NotifierAddedNotificationMessage struct {
//	MessageType message.MessageType                    `protobuf:"varint,1,opt,name=messageType,proto3,enum=MessageType" json:"messageType,omitempty"`
//	Data        *NotifierAddedNotificationMessage_Data `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
//}
//
//type PluginRegisterRequestMessage struct {
//
//
//	MessageType message.MessageType                `protobuf:"varint,1,opt,name=messageType,proto3,enum=MessageType" json:"messageType,omitempty"`
//	Data        *PluginRegisterRequestMessage_Data `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
//}
//
//type PluginRegisterResponseMessage struct {
//
//
//	MessageType message.MessageType                 `protobuf:"varint,1,opt,name=messageType,proto3,enum=MessageType" json:"messageType,omitempty"`
//	Data        *PluginRegisterResponseMessage_Data `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
//}
//
//type Preferences struct {
//	Language string             `protobuf:"bytes,1,opt,name=language,proto3" json:"language,omitempty"`
//	Units    *Preferences_Units `protobuf:"bytes,2,opt,name=units,proto3" json:"units,omitempty"`
//}
//
//type UsrProfile struct {
//	AddonsDir  string `protobuf:"bytes,1,opt,name=addonsDir,proto3" json:"addonsDir,omitempty"`
//	BaseDir    string `protobuf:"bytes,2,opt,name=baseDir,proto3" json:"baseDir,omitempty"`
//	ConfigDir  string `protobuf:"bytes,3,opt,name=configDir,proto3" json:"configDir,omitempty"`
//	DataDir    string `protobuf:"bytes,4,opt,name=dataDir,proto3" json:"dataDir,omitempty"`
//	MediaDir   string `protobuf:"bytes,5,opt,name=mediaDir,proto3" json:"mediaDir,omitempty"`
//	LogDir     string `protobuf:"bytes,6,opt,name=logDir,proto3" json:"logDir,omitempty"`
//	GatewayDir string `protobuf:"bytes,7,opt,name=gatewayDir,proto3" json:"gatewayDir,omitempty"`
//}
//
//type AdapterAddedNotificationMessage struct {
//	MessageType message.MessageType                   `protobuf:"varint,1,opt,name=messageType,proto3,enum=MessageType" json:"messageType,omitempty"`
//	Data        *AdapterAddedNotificationMessage_Data `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
//}
//
//type DeviceAddedNotificationMessage struct {
//	MessageType message.MessageType                  `protobuf:"varint,1,opt,name=messageType,proto3,enum=MessageType" json:"messageType,omitempty"`
//	Data        *DeviceAddedNotificationMessage_Data `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
//}
//
//type DevicePropertyChangedNotificationMessage struct {
//	MessageType message.MessageType                            `protobuf:"varint,1,opt,name=messageType,proto3,enum=MessageType" json:"messageType,omitempty"`
//	Data        *DevicePropertyChangedNotificationMessage_Data `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
//}
//
//type DeviceActionStatusNotificationMessage struct {
//	MessageType message.MessageType                         `protobuf:"varint,1,opt,name=messageType,proto3,enum=MessageType" json:"messageType,omitempty"`
//	Data        *DeviceActionStatusNotificationMessage_Data `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
//}
//
//type DeviceConnectedStateNotificationMessage struct {
//	MessageType message.MessageType                           `protobuf:"varint,1,opt,name=messageType,proto3,enum=MessageType" json:"messageType,omitempty"`
//	Data        *DeviceConnectedStateNotificationMessage_Data `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
//}
//
//// Addon_Device
//type Addon_Device struct {
//	AtContext           string               `protobuf:"bytes,1,opt,name=atContext,json=@context,proto3" json:"atContext,omitempty"`
//	AtType              string               `protobuf:"bytes,2,opt,name=atType,json=@type,proto3" json:"atType,omitempty"`
//	id                  string               `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
//	Title               string               `protobuf:"bytes,4,opt,name=title,proto3" json:"title,omitempty"`
//	Description         string               `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty"`
//	Links               []*Link              `protobuf:"bytes,6,rep,name=links,proto3" json:"links,omitempty"`
//	BaseHref            string               `protobuf:"bytes,7,opt,name=baseHref,proto3" json:"baseHref,omitempty"`
//	Pin                 *DevicePin           `protobuf:"bytes,8,opt,name=pin,proto3" json:"pin,omitempty"`
//	Properties          map[string]*Property `protobuf:"bytes,9,rep,name=properties,proto3" json:"properties,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
//	actions             map[string]*Action   `protobuf:"bytes,10,rep,name=actions,proto3" json:"actions,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
//	Events              map[string]*Event    `protobuf:"bytes,11,rep,name=events,proto3" json:"events,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
//	CredentialsRequired bool                 `protobuf:"varint,12,opt,name=credentialsRequired,proto3" json:"credentialsRequired,omitempty"`
//}
//
//type Property struct {
//	Instance        string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
//	AtType      string   `protobuf:"bytes,2,opt,name=atType,json=@type,proto3" json:"atType,omitempty"`
//	Title       string   `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
//	Type        string   `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
//	Unit        string   `protobuf:"bytes,5,opt,name=unit,proto3" json:"unit,omitempty"`
//	Description string   `protobuf:"bytes,6,opt,name=description,proto3" json:"description,omitempty"`
//	Minimum     float32  `protobuf:"fixed32,7,opt,name=minimum,proto3" json:"minimum,omitempty"`
//	Maximum     float32  `protobuf:"fixed32,8,opt,name=maximum,proto3" json:"maximum,omitempty"`
//	Enum        [][]byte `protobuf:"bytes,9,rep,name=enum,proto3" json:"enum,omitempty"`
//	ReadOnly    bool     `protobuf:"varint,10,opt,name=readOnly,proto3" json:"readOnly,omitempty"`
//	MultipleOf  float32  `protobuf:"fixed32,11,opt,name=multipleOf,proto3" json:"multipleOf,omitempty"`
//	Links       []*Link  `protobuf:"bytes,103,rep,name=links,proto3" json:"links,omitempty"`
//	Value       []byte   `protobuf:"bytes,12,opt,name=value,proto3" json:"value,omitempty"`
//}
//
//type Action struct {
//	AtType      string  `protobuf:"bytes,1,opt,name=atType,json=@type,proto3" json:"atType,omitempty"`
//	Title       string  `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
//	Description string  `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
//	Links       []*Link `protobuf:"bytes,4,rep,name=links,proto3" json:"links,omitempty"`
//	Input       []byte  `protobuf:"bytes,5,opt,name=input,proto3,oneof" json:"input,omitempty"`
//}
//
//type Event struct {
//	AtType      string     `protobuf:"bytes,1,opt,name=atType,json=@type,proto3" json:"atType,omitempty"`
//	Instance        string     `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
//	Title       string     `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
//	Description string     `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
//	Links       []*Link    `protobuf:"bytes,5,rep,name=links,proto3" json:"links,omitempty"`
//	Type        string     `protobuf:"bytes,6,opt,name=type,proto3" json:"type,omitempty"`
//	Unit        string     `protobuf:"bytes,7,opt,name=unit,proto3" json:"unit,omitempty"`
//	Minimum     float32    `protobuf:"fixed32,8,opt,name=minimum,proto3" json:"minimum,omitempty"`
//	Maximum     float32    `protobuf:"fixed32,9,opt,name=maximum,proto3" json:"maximum,omitempty"`
//	MultipleOf  float32    `protobuf:"fixed32,10,opt,name=multipleOf,proto3" json:"multipleOf,omitempty"`
//	Enum        []*any.Any `protobuf:"bytes,11,rep,name=enum,proto3" json:"enum,omitempty"`
//}
//
//type Link struct {
//	Href      string `protobuf:"bytes,1,opt,name=href,proto3" json:"href,omitempty"`
//	Rel       string `protobuf:"bytes,2,opt,name=rel,proto3" json:"rel,omitempty"`
//	MediaType string `protobuf:"bytes,3,opt,name=mediaType,proto3" json:"mediaType,omitempty"`
//}
//
//type DevicePin struct {
//	Required bool   `protobuf:"varint,1,opt,name=required,proto3" json:"required,omitempty"`
//	Pattern  string `protobuf:"bytes,2,opt,name=pattern,proto3" json:"pattern,omitempty"`
//}
//
//type ObjectActionInput struct {
//	Type       string                          `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
//	Properties map[string]*ActionInputProperty `protobuf:"bytes,2,rep,name=properties,proto3" json:"properties,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
//}
//
//type ActionInputProperty struct {
//	AtType     string     `protobuf:"bytes,1,opt,name=atType,json=@type,proto3" json:"atType,omitempty"`
//	Type       string     `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
//	Unit       string     `protobuf:"bytes,3,opt,name=unit,proto3" json:"unit,omitempty"`
//	Minimum    float32    `protobuf:"fixed32,4,opt,name=minimum,proto3" json:"minimum,omitempty"`
//	Maximum    float32    `protobuf:"fixed32,5,opt,name=maximum,proto3" json:"maximum,omitempty"`
//	MultipleOf float32    `protobuf:"fixed32,6,opt,name=multipleOf,proto3" json:"multipleOf,omitempty"`
//	Enum       []*any.Any `protobuf:"bytes,7,rep,name=enum,proto3" json:"enum,omitempty"`
//}
//
//type ActionDescription struct {
//	id            string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
//	Instance          string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
//	Input         []byte `protobuf:"bytes,3,opt,name=input,proto3,oneof" json:"input,omitempty"`
//	Status        string `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty"`
//	TimeRequested string `protobuf:"bytes,5,opt,name=timeRequested,proto3" json:"timeRequested,omitempty"`
//	TimeCompleted string `protobuf:"bytes,6,opt,name=timeCompleted,proto3" json:"timeCompleted,omitempty"`
//}
//
//type EventDescription struct {
//	state         protoimpl.MessageState
//	sizeCache     protoimpl.SizeCache
//	unknownFields protoimpl.UnknownFields
//
//	Instance      string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
//	Data      *any.Any `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
//	Timestamp string   `protobuf:"bytes,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
//}
//
//// Things
//type ThingDescription struct {
//	id                  string               `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
//	Title               string               `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
//	AtContext           string               `protobuf:"bytes,3,opt,name=atContext,json=@context,proto3" json:"atContext,omitempty"`
//	AtType              string               `protobuf:"bytes,4,opt,name=atType,json=@type,proto3" json:"atType,omitempty"`
//	Description         string               `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty"`
//	Base                string               `protobuf:"bytes,6,opt,name=base,proto3" json:"base,omitempty"`
//	BaseHref            string               `protobuf:"bytes,7,opt,name=baseHref,proto3" json:"baseHref,omitempty"`
//	Links               []*Link              `protobuf:"bytes,8,rep,name=links,proto3" json:"links,omitempty"`
//	Pin                 *DevicePin           `protobuf:"bytes,9,opt,name=pin,proto3" json:"pin,omitempty"`
//	Properties          map[string]*Property `protobuf:"bytes,10,rep,name=properties,proto3" json:"properties,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
//	actions             map[string]*Action   `protobuf:"bytes,11,rep,name=actions,proto3" json:"actions,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
//	Events              map[string]*Event    `protobuf:"bytes,12,rep,name=events,proto3" json:"events,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
//	CredentialsRequired bool                 `protobuf:"varint,13,opt,name=credentialsRequired,proto3" json:"credentialsRequired,omitempty"`
//	FloorplanVisibility bool                 `protobuf:"varint,14,opt,name=floorplanVisibility,proto3" json:"floorplanVisibility,omitempty"`
//	FloorplanX          uint32               `protobuf:"varint,15,opt,name=floorplanX,proto3" json:"floorplanX,omitempty"`
//	FloorplanY          uint32               `protobuf:"varint,16,opt,name=floorplanY,proto3" json:"floorplanY,omitempty"`
//	LayoutIndex         uint32               `protobuf:"varint,17,opt,name=layoutIndex,proto3" json:"layoutIndex,omitempty"`
//	SelectedCapability  string               `protobuf:"bytes,18,opt,name=selectedCapability,proto3" json:"selectedCapability,omitempty"`
//	IconHref            string               `protobuf:"bytes,19,opt,name=iconHref,proto3" json:"iconHref,omitempty"`
//	IconData            *IconData            `protobuf:"bytes,20,opt,name=iconData,proto3" json:"iconData,omitempty"`
//	Security            string               `protobuf:"bytes,21,opt,name=security,proto3" json:"security,omitempty"`
//	SecurityDefinitions *SecurityDefinition  `protobuf:"bytes,22,opt,name=securityDefinitions,proto3" json:"securityDefinitions,omitempty"`
//	GroupId             string               `protobuf:"bytes,23,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
//}
//
//type IconData struct {
//	state         protoimpl.MessageState
//	sizeCache     protoimpl.SizeCache
//	unknownFields protoimpl.UnknownFields
//
//	Data string `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
//	Mime string `protobuf:"bytes,2,opt,name=mime,proto3" json:"mime,omitempty"`
//}
//
//type SecurityDefinition struct {
//	state         protoimpl.MessageState
//	sizeCache     protoimpl.SizeCache
//	unknownFields protoimpl.UnknownFields
//
//	Oauth2Sc *OAuth2 `protobuf:"bytes,1,opt,name=oauth2_sc,json=oauth2Sc,proto3" json:"oauth2_sc,omitempty"`
//}
//
//type OAuth2 struct {
//	Scheme        string   `protobuf:"bytes,1,opt,name=scheme,proto3" json:"scheme,omitempty"`
//	Flow          string   `protobuf:"bytes,2,opt,name=flow,proto3" json:"flow,omitempty"`
//	Authorization string   `protobuf:"bytes,3,opt,name=authorization,proto3" json:"authorization,omitempty"`
//	Token         string   `protobuf:"bytes,4,opt,name=token,proto3" json:"token,omitempty"`
//	Scopes        []string `protobuf:"bytes,5,rep,name=scopes,proto3" json:"scopes,omitempty"`
//}
//
//type NotifierAddedNotificationMessage_Data struct {
//	PluginId    string `protobuf:"bytes,2,opt,name=pluginId,proto3" json:"pluginId,omitempty"`
//	NotifierId  string `protobuf:"bytes,3,opt,name=notifierId,proto3" json:"notifierId,omitempty"`
//	Instance        string `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
//	PackageName string `protobuf:"bytes,5,opt,name=packageName,proto3" json:"packageName,omitempty"`
//}
//
//type PluginRegisterRequestMessage_Data struct {
//	state         protoimpl.MessageState
//	sizeCache     protoimpl.SizeCache
//	unknownFields protoimpl.UnknownFields
//
//	PluginId string `protobuf:"bytes,11,opt,name=pluginId,proto3" json:"pluginId,omitempty"`
//}
//
//type PluginRegisterResponseMessage_Data struct {
//	state         protoimpl.MessageState
//	sizeCache     protoimpl.SizeCache
//	unknownFields protoimpl.UnknownFields
//
//	PluginId       string       `protobuf:"bytes,1,opt,name=pluginId,proto3" json:"pluginId,omitempty"`
//	GatewayVersion string       `protobuf:"bytes,2,opt,name=gatewayVersion,proto3" json:"gatewayVersion,omitempty"`
//	UserProfile    *UsrProfile  `protobuf:"bytes,3,opt,name=userProfile,proto3" json:"userProfile,omitempty"`
//	Preferences    *Preferences `protobuf:"bytes,4,opt,name=preferences,proto3" json:"preferences,omitempty"`
//}
//
//type Preferences_Units struct {
//	state         protoimpl.MessageState
//	sizeCache     protoimpl.SizeCache
//	unknownFields protoimpl.UnknownFields
//
//	Temperature string `protobuf:"bytes,1,opt,name=temperature,proto3" json:"temperature,omitempty"`
//}
//
//type AdapterAddedNotificationMessage_Data struct {
//	state         protoimpl.MessageState
//	sizeCache     protoimpl.SizeCache
//	unknownFields protoimpl.UnknownFields
//
//	PluginId    string `protobuf:"bytes,1,opt,name=pluginId,proto3" json:"pluginId,omitempty"`
//	AdapterId   string `protobuf:"bytes,2,opt,name=adapterId,proto3" json:"adapterId,omitempty"`
//	Instance        string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
//	PackageName string `protobuf:"bytes,4,opt,name=packageName,proto3" json:"packageName,omitempty"`
//}
//
//type DeviceAddedNotificationMessage_Data struct {
//	state         protoimpl.MessageState
//	sizeCache     protoimpl.SizeCache
//	unknownFields protoimpl.UnknownFields
//
//	PluginId  string  `protobuf:"bytes,1,opt,name=pluginId,proto3" json:"pluginId,omitempty"`
//	AdapterId string  `protobuf:"bytes,2,opt,name=adapterId,proto3" json:"adapterId,omitempty"`
//	Addon_Device    *Addon_Device `protobuf:"bytes,3,opt,name=addon,proto3" json:"addon,omitempty"`
//}
//
//type DevicePropertyChangedNotificationMessage_Data struct {
//	PluginId  string    `protobuf:"bytes,1,opt,name=pluginId,proto3" json:"pluginId,omitempty"`
//	AdapterId string    `protobuf:"bytes,2,opt,name=adapterId,proto3" json:"adapterId,omitempty"`
//	DeviceId  string    `protobuf:"bytes,3,opt,name=deviceId,proto3" json:"deviceId,omitempty"`
//	Property  *Property `protobuf:"bytes,4,opt,name=property,proto3" json:"property,omitempty"`
//}
//
//type DeviceActionStatusNotificationMessage_Data struct {
//	PluginId  string             `protobuf:"bytes,3,opt,name=pluginId,proto3" json:"pluginId,omitempty"`
//	AdapterId string             `protobuf:"bytes,4,opt,name=adapterId,proto3" json:"adapterId,omitempty"`
//	DeviceId  string             `protobuf:"bytes,5,opt,name=deviceId,proto3" json:"deviceId,omitempty"`
//	Action    *ActionDescription `protobuf:"bytes,6,opt,name=actions,proto3" json:"actions,omitempty"`
//}
//
//type DeviceConnectedStateNotificationMessage_Data struct {
//	PluginId  string `protobuf:"bytes,11,opt,name=pluginId,proto3" json:"pluginId,omitempty"`
//	AdapterId string `protobuf:"bytes,21,opt,name=adapterId,proto3" json:"adapterId,omitempty"`
//	DeviceId  string `protobuf:"bytes,22,opt,name=deviceId,proto3" json:"deviceId,omitempty"`
//	Connected bool   `protobuf:"varint,23,opt,name=connected,proto3" json:"connected,omitempty"`
//}
