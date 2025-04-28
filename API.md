# API Documentation

The API documentation covers the APIs exposed by the service. It gives an overview about the request method, payloads, and responses that each API returns.

Each API can be invoked and configured through the ReST client of your choice.

#### POST - /file/upload

The file-upload API is the entry point for the file upload, storage, and processing flows. Currently, the upload API limits the file size to 10MB and is stored in the local file system, underneath the uploads directory.

The service utilizes a background job processor to process uploaded content based on the content type.
- For video uploads, metadata generation is automatically done for you assuming that you were able to successfully set up application with redis.

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

#### PUT - /file/{id}/resize

The resize endpoint allows us to resize a file. Currently, only images can be resized and the task
is queued in Redis and carried out through a background job.

+ Request 

The resize API takes the following payload:

```
{
    "width: 123 // integer
    "height": 123 // integer
}
```

+ Response (202)

```
{
    message: "Image resize task enqueued"
}
```

+ Response (400) - 
```
{
    error: "Failed to parse request"
}
```

+ Response (404) - File is not found
```
{
    error: "File not found"
}
```

+ Response (422) - This response is returned for invalid requests and there are many cases where requests may be invalid. The respone structure coveys a message such that it is clear as to what is incorrect about the request.

```
{
    error: "Failed to enqueue resize task"
}
```

```
{
    error: "File is not an image"
}
```

```
{
    error: "File id is a required path parameter"
}
```