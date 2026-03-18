# SecretsApi

All URIs are relative to *http://localhost:8080/api/sdk/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**getField**](SecretsApi.md#getfield) | **GET** /secrets/{id}/fields/{fid} | Get a single decrypted field by ID |
| [**getFile**](SecretsApi.md#getfile) | **GET** /secrets/{id}/files/{fid} | Download a file by ID |
| [**getSecret**](SecretsApi.md#getsecret) | **GET** /secrets/{id} | Get secret with all decrypted fields and file metadata |
| [**listSecrets**](SecretsApi.md#listsecrets) | **GET** /secrets | List org secrets |



## getField

> SecretField getField(id, fid)

Get a single decrypted field by ID

Returns one field with its decrypted value. Useful when you only need a specific credential without fetching the entire secret (e.g. injecting a password into a CI environment). 

### Example

```ts
import {
  Configuration,
  SecretsApi,
} from '';
import type { GetFieldRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: apiKeyAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new SecretsApi(config);

  const body = {
    // string | Secret ID
    id: 38400000-8cf0-11bd-b23e-10b96e4ef00d,
    // string | Field ID
    fid: 38400000-8cf0-11bd-b23e-10b96e4ef00d,
  } satisfies GetFieldRequest;

  try {
    const data = await api.getField(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | `string` | Secret ID | [Defaults to `undefined`] |
| **fid** | `string` | Field ID | [Defaults to `undefined`] |

### Return type

[**SecretField**](SecretField.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **401** | Missing or invalid API key |  -  |
| **404** | Resource not found or does not belong to this org |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## getFile

> Blob getFile(id, fid)

Download a file by ID

Streams the raw file bytes. Use the &#x60;Content-Disposition&#x60; header to determine the suggested filename. The &#x60;Content-Type&#x60; header reflects the original file\&#39;s MIME type. 

### Example

```ts
import {
  Configuration,
  SecretsApi,
} from '';
import type { GetFileRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: apiKeyAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new SecretsApi(config);

  const body = {
    // string | Secret ID
    id: 38400000-8cf0-11bd-b23e-10b96e4ef00d,
    // string | File ID
    fid: 38400000-8cf0-11bd-b23e-10b96e4ef00d,
  } satisfies GetFileRequest;

  try {
    const data = await api.getFile(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | `string` | Secret ID | [Defaults to `undefined`] |
| **fid** | `string` | File ID | [Defaults to `undefined`] |

### Return type

**Blob**

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/octet-stream`, `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | File bytes |  * Content-Disposition -  <br>  * Content-Type -  <br>  |
| **401** | Missing or invalid API key |  -  |
| **404** | Resource not found or does not belong to this org |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## getSecret

> Secret getSecret(id)

Get secret with all decrypted fields and file metadata

Returns the full secret object including: - All fields with their **decrypted** values - All file metadata with SDK download URLs  The secret must belong to the org that owns the API key. 

### Example

```ts
import {
  Configuration,
  SecretsApi,
} from '';
import type { GetSecretRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: apiKeyAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new SecretsApi(config);

  const body = {
    // string | Secret ID
    id: 38400000-8cf0-11bd-b23e-10b96e4ef00d,
  } satisfies GetSecretRequest;

  try {
    const data = await api.getSecret(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | `string` | Secret ID | [Defaults to `undefined`] |

### Return type

[**Secret**](Secret.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **401** | Missing or invalid API key |  -  |
| **404** | Resource not found or does not belong to this org |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## listSecrets

> Array&lt;SecretSummary&gt; listSecrets()

List org secrets

Returns all secrets belonging to the org associated with the API key. Fields and files are **not** included in the list — fetch them via &#x60;GET /secrets/{id}&#x60; or the individual field/file endpoints. 

### Example

```ts
import {
  Configuration,
  SecretsApi,
} from '';
import type { ListSecretsRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: apiKeyAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new SecretsApi(config);

  try {
    const data = await api.listSecrets();
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**Array&lt;SecretSummary&gt;**](SecretSummary.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **401** | Missing or invalid API key |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

