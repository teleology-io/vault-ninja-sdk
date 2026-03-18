# Secret

A secret with all its decrypted fields and file metadata. Secrets are always scoped to the org that owns the API key. 

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **UUID** |  | 
**title** | **str** |  | 
**description** | **str** |  | 
**tags** | **List[str]** |  | 
**org_id** | **UUID** |  | 
**created_by** | **UUID** |  | 
**created_at** | **datetime** |  | 
**updated_at** | **datetime** |  | 
**fields** | [**List[SecretField]**](SecretField.md) |  | 
**files** | [**List[SecretFile]**](SecretFile.md) |  | 

## Example

```python
from vault_ninja.models.secret import Secret

# TODO update the JSON string below
json = "{}"
# create an instance of Secret from a JSON string
secret_instance = Secret.from_json(json)
# print the JSON string representation of the object
print(Secret.to_json())

# convert the object into a dict
secret_dict = secret_instance.to_dict()
# create an instance of Secret from a dict
secret_from_dict = Secret.from_dict(secret_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


