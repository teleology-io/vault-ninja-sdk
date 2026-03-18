
# SecretFile

Metadata for a file attached to a secret. Use the `url` field to download the file contents — it points to `GET /secrets/{secretId}/files/{fileId}`. 

## Properties

Name | Type
------------ | -------------
`id` | string
`secretId` | string
`originalName` | string
`contentType` | string
`sizeBytes` | number
`position` | number
`createdAt` | Date
`url` | string

## Example

```typescript
import type { SecretFile } from ''

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "secretId": null,
  "originalName": null,
  "contentType": null,
  "sizeBytes": null,
  "position": null,
  "createdAt": null,
  "url": null,
} satisfies SecretFile

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as SecretFile
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


