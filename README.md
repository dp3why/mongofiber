# Mongofiber

Go Fiber + MongoDB + Google Cloud Storage

## Introduction

- Leveraged MongoDB for data management, Google Cloud Storage for storing images

## Endpoints

#### GET /books

Returns all books

#### GET /books/:id

Returns a single book

#### POST /books

Creates a new book

input:

```
{
    "title": "test book",
    "author": "me",
    "year": "2022"
}
```

#### PUT /books/:id

Updates a book

input:

```
{
    "title": "test book",
    "author": "me",
    "year": "2022"
}
```

#### DELETE /books/:id

Deletes a book
