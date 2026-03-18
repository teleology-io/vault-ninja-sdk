
# SecretSummary

Lightweight secret listing — no fields or files.

## Properties

Name | Type
------------ | -------------
`id` | string
`title` | string
`description` | string
`tags` | Array&lt;string&gt;
`orgId` | string
`createdBy` | string
`createdAt` | Date
`updatedAt` | Date

## Example

```typescript
import type { SecretSummary } from ''

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "title": null,
  "description": null,
  "tags": null,
  "orgId": null,
  "createdBy": null,
  "createdAt": null,
  "updatedAt": null,
} satisfies SecretSummary

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as SecretSummary
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


