# API Gateway/Lambda for query embedding database

## Topics

- Password stored in ssm, has to be migrated to secrets manager
- API Gateway/Lambda for query embedding database, 30 seconds timeout

## Call

```json
{
    "question": "This is my question?"
}
```

## Response


```json
{
    "answer": "lorem ipsum",
    "documents": [
        {
            "content": "short text",
            "context": "long text"
        }
      ]
}
```
