
# SecretField

A single field belonging to a secret. The `value` is returned decrypted.

## Properties

Name | Type
------------ | -------------
`id` | string
`secretId` | string
`position` | number
`fieldType` | [FieldType](FieldType.md)
`label` | string
`value` | string
`createdAt` | Date

## Example

```typescript
import type { SecretField } from ''

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "secretId": null,
  "position": null,
  "fieldType": null,
  "label": null,
  "value": null,
  "createdAt": null,
} satisfies SecretField

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as SecretField
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


