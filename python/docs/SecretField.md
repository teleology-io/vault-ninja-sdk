# SecretField

A single field belonging to a secret. The `value` is returned decrypted.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **UUID** |  | 
**secret_id** | **UUID** |  | 
**position** | **int** |  | 
**field_type** | [**FieldType**](FieldType.md) |  | 
**label** | **str** |  | 
**value** | **str** | Decrypted plaintext value | 
**created_at** | **datetime** |  | 

## Example

```python
from vault_ninja.models.secret_field import SecretField

# TODO update the JSON string below
json = "{}"
# create an instance of SecretField from a JSON string
secret_field_instance = SecretField.from_json(json)
# print the JSON string representation of the object
print(SecretField.to_json())

# convert the object into a dict
secret_field_dict = secret_field_instance.to_dict()
# create an instance of SecretField from a dict
secret_field_from_dict = SecretField.from_dict(secret_field_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


