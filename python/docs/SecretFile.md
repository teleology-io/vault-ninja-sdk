# SecretFile

Metadata for a file attached to a secret. Use the `url` field to download the file contents — it points to `GET /secrets/{secretId}/files/{fileId}`. 

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **UUID** |  | 
**secret_id** | **UUID** |  | 
**original_name** | **str** |  | 
**content_type** | **str** |  | 
**size_bytes** | **int** |  | 
**position** | **int** |  | 
**created_at** | **datetime** |  | 
**url** | **str** | SDK download URL for this file | 

## Example

```python
from vault_ninja.models.secret_file import SecretFile

# TODO update the JSON string below
json = "{}"
# create an instance of SecretFile from a JSON string
secret_file_instance = SecretFile.from_json(json)
# print the JSON string representation of the object
print(SecretFile.to_json())

# convert the object into a dict
secret_file_dict = secret_file_instance.to_dict()
# create an instance of SecretFile from a dict
secret_file_from_dict = SecretFile.from_dict(secret_file_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


