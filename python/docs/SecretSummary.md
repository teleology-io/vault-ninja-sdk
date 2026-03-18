# SecretSummary

Lightweight secret listing — no fields or files.

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

## Example

```python
from vault_ninja.models.secret_summary import SecretSummary

# TODO update the JSON string below
json = "{}"
# create an instance of SecretSummary from a JSON string
secret_summary_instance = SecretSummary.from_json(json)
# print the JSON string representation of the object
print(SecretSummary.to_json())

# convert the object into a dict
secret_summary_dict = secret_summary_instance.to_dict()
# create an instance of SecretSummary from a dict
secret_summary_from_dict = SecretSummary.from_dict(secret_summary_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


