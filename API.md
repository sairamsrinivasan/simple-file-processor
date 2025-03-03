# API Documentation

The API documentation covers the APIs exposed by the service. It gives an overview about the request method, payloads, and responses that each API returns.

Each API can be invoked and configured through the ReST client of your choice.

#### POST - /upload

The upload API is the entry point for the file upload, storage, and processing flows. Currently, the upload API contains a limit of 10MB.

+ Request 

The upload API takes form data as input, with "file" as the key and the value being the file selected for upload.

+ Response (200)

```
{
    "ID": "a0de50ee-d9f6-4fc3-8b26-16242724f0e9",
    "generated_name": "a0de50ee-d9f6-4fc3-8b26-16242724f0e9/dj.jpeg",
    "mime_type": "image/jpeg",
    "original_name": "dj.jpeg",
    "size": 4417534,
    "status": "pending",
    "storage_path": "uploads/a0de50ee-d9f6-4fc3-8b26-16242724f0e9",
    "type": "image",
    "uploaded_extension": "jpeg",
    "created_at": "2025-02-17T19:40:06.612592-08:00",
    "updated_at": "2025-02-17T19:40:06.612592-08:00"
}
```


+ Response (413) - File is too large
+ Response (400) - payload does not contain the "file" key within form-data
+ Response (500) - failure creating the file on the server or inserting file metadata information into database