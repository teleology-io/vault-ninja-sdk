# SecretField

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**SecretId** | **string** |  | 
**Position** | **int32** |  | 
**FieldType** | [**FieldType**](FieldType.md) |  | 
**Label** | **string** |  | 
**Value** | **string** | Decrypted plaintext value | 
**CreatedAt** | **time.Time** |  | 

## Methods

### NewSecretField

`func NewSecretField(id string, secretId string, position int32, fieldType FieldType, label string, value string, createdAt time.Time, ) *SecretField`

NewSecretField instantiates a new SecretField object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSecretFieldWithDefaults

`func NewSecretFieldWithDefaults() *SecretField`

NewSecretFieldWithDefaults instantiates a new SecretField object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *SecretField) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *SecretField) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *SecretField) SetId(v string)`

SetId sets Id field to given value.


### GetSecretId

`func (o *SecretField) GetSecretId() string`

GetSecretId returns the SecretId field if non-nil, zero value otherwise.

### GetSecretIdOk

`func (o *SecretField) GetSecretIdOk() (*string, bool)`

GetSecretIdOk returns a tuple with the SecretId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSecretId

`func (o *SecretField) SetSecretId(v string)`

SetSecretId sets SecretId field to given value.


### GetPosition

`func (o *SecretField) GetPosition() int32`

GetPosition returns the Position field if non-nil, zero value otherwise.

### GetPositionOk

`func (o *SecretField) GetPositionOk() (*int32, bool)`

GetPositionOk returns a tuple with the Position field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPosition

`func (o *SecretField) SetPosition(v int32)`

SetPosition sets Position field to given value.


### GetFieldType

`func (o *SecretField) GetFieldType() FieldType`

GetFieldType returns the FieldType field if non-nil, zero value otherwise.

### GetFieldTypeOk

`func (o *SecretField) GetFieldTypeOk() (*FieldType, bool)`

GetFieldTypeOk returns a tuple with the FieldType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFieldType

`func (o *SecretField) SetFieldType(v FieldType)`

SetFieldType sets FieldType field to given value.


### GetLabel

`func (o *SecretField) GetLabel() string`

GetLabel returns the Label field if non-nil, zero value otherwise.

### GetLabelOk

`func (o *SecretField) GetLabelOk() (*string, bool)`

GetLabelOk returns a tuple with the Label field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabel

`func (o *SecretField) SetLabel(v string)`

SetLabel sets Label field to given value.


### GetValue

`func (o *SecretField) GetValue() string`

GetValue returns the Value field if non-nil, zero value otherwise.

### GetValueOk

`func (o *SecretField) GetValueOk() (*string, bool)`

GetValueOk returns a tuple with the Value field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetValue

`func (o *SecretField) SetValue(v string)`

SetValue sets Value field to given value.


### GetCreatedAt

`func (o *SecretField) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *SecretField) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *SecretField) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


