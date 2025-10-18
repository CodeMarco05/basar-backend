# BASAR Backend API Documentation

## Table of Contents
- [Health Check](#health-check)
- [Posts](#posts)
  - [Get All Posts](#get-all-posts)
  - [Get Single Post](#get-single-post)
  - [Create Post](#create-post)
- [Comments](#comments)
  - [Create Comment](#create-comment)
  - [Delete Comment](#delete-comment)
- [User Posts](#user-posts)
  - [Get Posts by Creator](#get-posts-by-creator)
  - [Update Post](#update-post)
  - [Delete Post](#delete-post)

---

## Health Check

### Health Check Endpoint

**Endpoint:** `GET /health`

**Description:** Verifies that the application and database are running.

**Authentication:** Not required

**Request:**
```http
GET /health
```

**Response:**

Success (200 OK):
```
The application appears to be up and running
```

Error (401 Unauthorized):
```
The database isn't available.
```

---

## Posts

### Get All Posts

**Endpoint:** `GET /api/v1/posts`

**Description:** Retrieves all posts from the database.

**Authentication:** Not required

**Request:**
```http
GET /api/v1/posts
```

**Response:**

Success (200 OK):
```json
[
  {
    "id": "507f1f77bcf86cd799439011",
    "creatorId": "firebaseUUID",
    "title": "Vintage Camera for Sale",
    "description": "Beautiful vintage camera in excellent condition",
    "tags": ["electronics", "vintage", "camera"],
    "text": "Detailed description of the camera...",
    "payPalMail": "seller@example.com",
    "images": [
      "https://example.com/image1.jpg",
      "https://example.com/image2.jpg"
    ],
    "created_at": "2024-01-15T10:30:00Z"
  }
]
```

Error (500 Internal Server Error):
```json
{
  "error": "Internal Server Error"
}
```

---

### Get Single Post

**Endpoint:** `GET /api/v1/posts/{postId}`

**Description:** Retrieves a specific post by its ID.

**Authentication:** Not required

**Path Parameters:**
- `postId` (string, required): MongoDB ObjectId of the post

**Request:**
```http
GET /api/v1/posts/507f1f77bcf86cd799439011
```

**Response:**

Success (200 OK):
```json
{
  "id": "507f1f77bcf86cd799439011",
  "creatorId": "firebaseUUID",
  "title": "Vintage Camera for Sale",
  "description": "Beautiful vintage camera in excellent condition",
  "tags": ["electronics", "vintage", "camera"],
  "text": "Detailed description of the camera...",
  "payPalMail": "seller@example.com",
  "images": [
    "https://example.com/image1.jpg",
    "https://example.com/image2.jpg"
  ],
  "created_at": "2024-01-15T10:30:00Z"
}
```

Error (400 Bad Request):
```json
{
  "error": "PostId was missing or given under a false key"
}
```
or
```json
{
  "error": "PostId was not found or it was an internal server error"
}
```

---

### Create Post

**Endpoint:** `POST /api/v1/posts`

**Description:** Creates a new post in the database.

**Authentication:** Not required

**Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
  "creatorId": "firebaseUUID",
  "title": "Vintage Camera for Sale",
  "description": "Beautiful vintage camera in excellent condition",
  "tags": ["electronics", "vintage", "camera"],
  "text": "This is a rare vintage camera from the 1970s. It has been well maintained and comes with the original leather case. Perfect for collectors or photography enthusiasts.",
  "payPalMail": "seller@example.com",
  "images": [
    "https://example.com/image1.jpg",
    "https://example.com/image2.jpg"
  ]
}
```

**Field Requirements:**
- `creatorId` (string, required): Firebase UUID of the post creator
- `title` (string, required): Post title
- `description` (string, required): Short description of the item
- `tags` (array of strings, required): Non-empty array of tags/categories
- `text` (string, required): Detailed description/content
- `payPalMail` (string, required): Valid email address for PayPal payments
- `images` (array of strings, required): Non-empty array of image URLs

**Response:**

Success (200 OK):
```json
{
  "id": "507f1f77bcf86cd799439011"
}
```

Error (400 Bad Request):
```json
{
  "error": "Invalid request body when transforming to the required object."
}
```

Error (500 Internal Server Error):
```json
{
  "error": "Internal Server Error"
}
```

---

## Comments

### Create Comment

**Endpoint:** `POST /api/v1/posts/{postId}/comments`

**Description:** Creates a new comment on a specific post.

**Authentication:** Not required

**Path Parameters:**
- `postId` (string, required): MongoDB ObjectId of the post

**Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
  "creatorId": "firebaseUUID",
  "message": "Great item! Is it still available?",
  "commenterMail": "commenter@example.com",
  "commenterName": "John Doe"
}
```

**Field Requirements:**
- `creatorId` (string, required): Firebase UUID of the comment creator
- `message` (string, required): The comment text/message
- `commenterMail` (string, required): Email address of the commenter
- `commenterName` (string, required): Display name of the commenter

**Response:**

Success (200 OK):
```
(empty response body)
```

Error (400 Bad Request):
```
postId is required
```
or
```
comment is required
```
or
```json
{
  "error": "Validation failed",
  "fields": {
    "Message": "Field 'Message' failed validation: required"
  }
}
```

Error (500 Internal Server Error):
```
Failed to create comment
```

---

### Delete Comment

**Endpoint:** `DELETE /api/v1/posts/{postId}/comments/{commentId}/{commentCreatorId}`

**Description:** Deletes a specific comment from a post. Only the comment creator can delete their own comment.

**Authentication:** Not required

**Path Parameters:**
- `postId` (string, required): MongoDB ObjectId of the post
- `commentId` (string, required): UUID of the comment to delete
- `commentCreatorId` (string, required): Firebase UUID of the comment creator

**Request:**
```http
DELETE /api/v1/posts/507f1f77bcf86cd799439011/comments/comment-uuid-123/firebaseUUID
```

**Response:**

Success (200 OK):
```
(empty response body)
```

Error (400 Bad Request):
```
postId is required
```
or
```
commentCreatorId is required
```
or
```
commentId is required
```

Error (500 Internal Server Error):
```
Failed to delete comment
```

---

## User Posts

### Get Posts by Creator

**Endpoint:** `GET /api/v1/users/{creatorId}/posts`

**Description:** Retrieves all posts created by a specific user.

**Authentication:** Not required

**Path Parameters:**
- `creatorId` (string, required): Firebase UUID of the post creator

**Request:**
```http
GET /api/v1/users/firebaseUUID/posts
```

**Response:**

Success (200 OK):
```json
[
  {
    "id": "507f1f77bcf86cd799439011",
    "creatorId": "firebaseUUID",
    "title": "Vintage Camera for Sale",
    "description": "Beautiful vintage camera in excellent condition",
    "tags": ["electronics", "vintage", "camera"],
    "text": "Detailed description of the camera...",
    "payPalMail": "seller@example.com",
    "images": [
      "https://example.com/image1.jpg",
      "https://example.com/image2.jpg"
    ],
    "created_at": "2024-01-15T10:30:00Z"
  }
]
```

Error (400 Bad Request):
```json
{
  "error": "creatorId was missing or given under a false key"
}
```

Error (500 Internal Server Error):
```json
{
  "error": "Internal Server Error"
}
```

---

### Update Post

**Endpoint:** `PATCH /api/v1/users/{creatorId}/{postId}`

**Description:** Updates an existing post. Only the creator can update their own post.

**Authentication:** Not required

**Path Parameters:**
- `creatorId` (string, required): Firebase UUID of the post creator
- `postId` (string, required): MongoDB ObjectId of the post

**Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
  "creatorId": "firebaseUUID",
  "title": "Vintage Camera for Sale - Price Reduced",
  "description": "Beautiful vintage camera in excellent condition - Now with reduced price!",
  "tags": ["electronics", "vintage", "camera", "sale"],
  "text": "Updated description with new pricing information...",
  "payPalMail": "seller@example.com",
  "images": [
    "https://example.com/image1.jpg",
    "https://example.com/image2.jpg",
    "https://example.com/image3.jpg"
  ]
}
```

**Field Requirements:** Same as Create Post (all fields required)

**Response:**

Success (200 OK):
```
(empty response body)
```

Error (400 Bad Request):
```json
{
  "error": "postId was missing or given under a false key"
}
```
or
```json
{
  "error": "creatorId was missing or given under a false key"
}
```
or
```json
{
  "error": "Validation failed",
  "fields": {
    "PayPalMail": "Field 'PayPalMail' failed validation: email"
  }
}
```

Error (500 Internal Server Error):
```json
{
  "error": "Error during post update: [error details]"
}
```

---

### Delete Post

**Endpoint:** `DELETE /api/v1/users/{creatorId}/{postId}`

**Description:** Deletes a post. Only the creator can delete their own post.

**Authentication:** Not required

**Path Parameters:**
- `creatorId` (string, required): Firebase UUID of the post creator
- `postId` (string, required): MongoDB ObjectId of the post

**Request:**
```http
DELETE /api/v1/users/firebaseUUID/posts/507f1f77bcf86cd799439011
```

**Response:**

Success (200 OK):
```
(empty response body)
```

Error (400 Bad Request):
```json
{
  "error": "postId was missing or given under a false key"
}
```
or
```json
{
  "error": "creatorId was missing or given under a false key"
}
```

Error (500 Internal Server Error):
```json
{
  "error": "Error during post delete: [error details]"
}
```

---

## Data Models

### Post Object
```json
{
  "id": "string (MongoDB ObjectId)",
  "creatorId": "string",
  "title": "string",
  "description": "string",
  "tags": ["string"],
  "text": "string",
  "payPalMail": "string (email format)",
  "images": ["string (URLs)"],
  "created_at": "timestamp (ISO 8601)"
}
```

### InsertPost Object (Create/Update Request)
```json
{
  "creatorId": "string (required, Firebase UUID)",
  "title": "string (required)",
  "description": "string (required)",
  "tags": ["string"] (required, non-empty array),
  "text": "string (required)",
  "payPalMail": "string (required, valid email)",
  "images": ["string"] (required, non-empty array)
}
```

### Comment Object
```json
{
  "id": "string (UUID)",
  "creatorId": "string (Firebase UUID)",
  "message": "string",
  "commenterMail": "string (email format)",
  "commenterName": "string",
  "createdAt": "timestamp (ISO 8601)"
}
```

---

## Validation Rules

### Email Validation
- `payPalMail` must be a valid email format
- Example: `user@example.com`

### Array Validation
- `tags` must be a non-empty array
- `images` must be a non-empty array

### Required Fields
All fields in InsertPost object are required:
- creatorId
- title
- description
- tags
- text
- payPalMail
- images

---

## Error Handling

### HTTP Status Codes
- `200 OK` - Request successful
- `400 Bad Request` - Invalid request parameters or body
- `401 Unauthorized` - Database connection issue (health check only)
- `500 Internal Server Error` - Server-side error

### Error Response Format
```json
{
  "error": "Error message description"
}
```

For validation errors:
```json
{
  "error": "Validation failed",
  "fields": {
    "FieldName": "Field 'FieldName' failed validation: validationType"
  }
}
```

---

## Notes

- **Authentication:** While authentication middleware is implemented (Firebase-based), it is currently disabled for all endpoints.
- **CORS:** Check server configuration for CORS settings if accessing from web browsers.
- **MongoDB ObjectIds:** All `id` and `postId` fields use MongoDB's ObjectId format (24 character hex string).
- **Timestamps:** The `created_at` field follows ISO 8601 format.

---

## Quick Reference

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/api/v1/posts` | Get all posts |
| GET | `/api/v1/posts/{postId}` | Get single post |
| POST | `/api/v1/posts` | Create new post |
| POST | `/api/v1/posts/{postId}/comments` | Create comment on post |
| DELETE | `/api/v1/posts/{postId}/comments/{commentId}/{commentCreatorId}` | Delete comment |
| GET | `/api/v1/users/{creatorId}/posts` | Get posts by creator |
| PATCH | `/api/v1/users/{creatorId}/{postId}` | Update post |
| DELETE | `/api/v1/users/{creatorId}/{postId}` | Delete post |
