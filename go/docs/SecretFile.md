# SecretFile

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**SecretId** | **string** |  | 
**OriginalName** | **string** |  | 
**ContentType** | **string** |  | 
**SizeBytes** | **int64** |  | 
**Position** | **int32** |  | 
**CreatedAt** | **time.Time** |  | 
**Url** | **string** | SDK download URL for this file | 

## Methods

### NewSecretFile

`func NewSecretFile(id string, secretId string, originalName string, contentType string, sizeBytes int64, position int32, createdAt time.Time, url string, ) *SecretFile`

NewSecretFile instantiates a new SecretFile object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSecretFileWithDefaults

`func NewSecretFileWithDefaults() *SecretFile`

NewSecretFileWithDefaults instantiates a new SecretFile object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *SecretFile) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *SecretFile) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *SecretFile) SetId(v string)`

SetId sets Id field to given value.


### GetSecretId

`func (o *SecretFile) GetSecretId() string`

GetSecretId returns the SecretId field if non-nil, zero value otherwise.

### GetSecretIdOk

`func (o *SecretFile) GetSecretIdOk() (*string, bool)`

GetSecretIdOk returns a tuple with the SecretId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSecretId

`func (o *SecretFile) SetSecretId(v string)`

SetSecretId sets SecretId field to given value.


### GetOriginalName

`func (o *SecretFile) GetOriginalName() string`

GetOriginalName returns the OriginalName field if non-nil, zero value otherwise.

### GetOriginalNameOk

`func (o *SecretFile) GetOriginalNameOk() (*string, bool)`

GetOriginalNameOk returns a tuple with the OriginalName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOriginalName

`func (o *SecretFile) SetOriginalName(v string)`

SetOriginalName sets OriginalName field to given value.


### GetContentType

`func (o *SecretFile) GetContentType() string`

GetContentType returns the ContentType field if non-nil, zero value otherwise.

### GetContentTypeOk

`func (o *SecretFile) GetContentTypeOk() (*string, bool)`

GetContentTypeOk returns a tuple with the ContentType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContentType

`func (o *SecretFile) SetContentType(v string)`

SetContentType sets ContentType field to given value.


### GetSizeBytes

`func (o *SecretFile) GetSizeBytes() int64`

GetSizeBytes returns the SizeBytes field if non-nil, zero value otherwise.

### GetSizeBytesOk

`func (o *SecretFile) GetSizeBytesOk() (*int64, bool)`

GetSizeBytesOk returns a tuple with the SizeBytes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSizeBytes

`func (o *SecretFile) SetSizeBytes(v int64)`

SetSizeBytes sets SizeBytes field to given value.


### GetPosition

`func (o *SecretFile) GetPosition() int32`

GetPosition returns the Position field if non-nil, zero value otherwise.

### GetPositionOk

`func (o *SecretFile) GetPositionOk() (*int32, bool)`

GetPositionOk returns a tuple with the Position field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPosition

`func (o *SecretFile) SetPosition(v int32)`

SetPosition sets Position field to given value.


### GetCreatedAt

`func (o *SecretFile) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *SecretFile) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *SecretFile) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetUrl

`func (o *SecretFile) GetUrl() string`

GetUrl returns the Url field if non-nil, zero value otherwise.

### GetUrlOk

`func (o *SecretFile) GetUrlOk() (*string, bool)`

GetUrlOk returns a tuple with the Url field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUrl

`func (o *SecretFile) SetUrl(v string)`

SetUrl sets Url field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


