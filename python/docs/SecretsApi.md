# vault_ninja.SecretsApi

All URIs are relative to *http://localhost:8080/api/sdk/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**get_field**](SecretsApi.md#get_field) | **GET** /secrets/{id}/fields/{fid} | Get a single decrypted field by ID
[**get_file**](SecretsApi.md#get_file) | **GET** /secrets/{id}/files/{fid} | Download a file by ID
[**get_secret**](SecretsApi.md#get_secret) | **GET** /secrets/{id} | Get secret with all decrypted fields and file metadata
[**list_secrets**](SecretsApi.md#list_secrets) | **GET** /secrets | List org secrets


# **get_field**
> SecretField get_field(id, fid)

Get a single decrypted field by ID

Returns one field with its decrypted value.
Useful when you only need a specific credential without fetching the
entire secret (e.g. injecting a password into a CI environment).


### Example

* Bearer (vn_org_<hex>) Authentication (apiKeyAuth):

```python
import vault_ninja
from vault_ninja.models.secret_field import SecretField
from vault_ninja.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost:8080/api/sdk/v1
# See configuration.py for a list of all supported configuration parameters.
configuration = vault_ninja.Configuration(
    host = "http://localhost:8080/api/sdk/v1"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (vn_org_<hex>): apiKeyAuth
configuration = vault_ninja.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
with vault_ninja.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = vault_ninja.SecretsApi(api_client)
    id = UUID('38400000-8cf0-11bd-b23e-10b96e4ef00d') # UUID | Secret ID
    fid = UUID('38400000-8cf0-11bd-b23e-10b96e4ef00d') # UUID | Field ID

    try:
        # Get a single decrypted field by ID
        api_response = api_instance.get_field(id, fid)
        print("The response of SecretsApi->get_field:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling SecretsApi->get_field: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **UUID**| Secret ID | 
 **fid** | **UUID**| Field ID | 

### Return type

[**SecretField**](SecretField.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**401** | Missing or invalid API key |  -  |
**404** | Resource not found or does not belong to this org |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **get_file**
> bytearray get_file(id, fid)

Download a file by ID

Streams the raw file bytes. Use the `Content-Disposition` header
to determine the suggested filename. The `Content-Type` header
reflects the original file's MIME type.


### Example

* Bearer (vn_org_<hex>) Authentication (apiKeyAuth):

```python
import vault_ninja
from vault_ninja.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost:8080/api/sdk/v1
# See configuration.py for a list of all supported configuration parameters.
configuration = vault_ninja.Configuration(
    host = "http://localhost:8080/api/sdk/v1"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (vn_org_<hex>): apiKeyAuth
configuration = vault_ninja.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
with vault_ninja.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = vault_ninja.SecretsApi(api_client)
    id = UUID('38400000-8cf0-11bd-b23e-10b96e4ef00d') # UUID | Secret ID
    fid = UUID('38400000-8cf0-11bd-b23e-10b96e4ef00d') # UUID | File ID

    try:
        # Download a file by ID
        api_response = api_instance.get_file(id, fid)
        print("The response of SecretsApi->get_file:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling SecretsApi->get_file: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **UUID**| Secret ID | 
 **fid** | **UUID**| File ID | 

### Return type

**bytearray**

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/octet-stream, application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | File bytes |  * Content-Disposition -  <br>  * Content-Type -  <br>  |
**401** | Missing or invalid API key |  -  |
**404** | Resource not found or does not belong to this org |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **get_secret**
> Secret get_secret(id)

Get secret with all decrypted fields and file metadata

Returns the full secret object including:
- All fields with their **decrypted** values
- All file metadata with SDK download URLs

The secret must belong to the org that owns the API key.


### Example

* Bearer (vn_org_<hex>) Authentication (apiKeyAuth):

```python
import vault_ninja
from vault_ninja.models.secret import Secret
from vault_ninja.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost:8080/api/sdk/v1
# See configuration.py for a list of all supported configuration parameters.
configuration = vault_ninja.Configuration(
    host = "http://localhost:8080/api/sdk/v1"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (vn_org_<hex>): apiKeyAuth
configuration = vault_ninja.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
with vault_ninja.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = vault_ninja.SecretsApi(api_client)
    id = UUID('38400000-8cf0-11bd-b23e-10b96e4ef00d') # UUID | Secret ID

    try:
        # Get secret with all decrypted fields and file metadata
        api_response = api_instance.get_secret(id)
        print("The response of SecretsApi->get_secret:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling SecretsApi->get_secret: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **UUID**| Secret ID | 

### Return type

[**Secret**](Secret.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**401** | Missing or invalid API key |  -  |
**404** | Resource not found or does not belong to this org |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **list_secrets**
> List[SecretSummary] list_secrets()

List org secrets

Returns all secrets belonging to the org associated with the API key.
Fields and files are **not** included in the list — fetch them via
`GET /secrets/{id}` or the individual field/file endpoints.


### Example

* Bearer (vn_org_<hex>) Authentication (apiKeyAuth):

```python
import vault_ninja
from vault_ninja.models.secret_summary import SecretSummary
from vault_ninja.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost:8080/api/sdk/v1
# See configuration.py for a list of all supported configuration parameters.
configuration = vault_ninja.Configuration(
    host = "http://localhost:8080/api/sdk/v1"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (vn_org_<hex>): apiKeyAuth
configuration = vault_ninja.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
with vault_ninja.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = vault_ninja.SecretsApi(api_client)

    try:
        # List org secrets
        api_response = api_instance.list_secrets()
        print("The response of SecretsApi->list_secrets:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling SecretsApi->list_secrets: %s\n" % e)
```



### Parameters

This endpoint does not need any parameter.

### Return type

[**List[SecretSummary]**](SecretSummary.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**401** | Missing or invalid API key |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

