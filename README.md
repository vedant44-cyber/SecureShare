# secure-share
| Domain       | Technology     | Specific Library/Tool       | Why We Are Using It                                                                   |
| ------------ | -------------- | --------------------------- | ------------------------------------------------------------------------------------- |
| **Frontend** | React (Vite)   | axios                       | Fast UI development and easier API communication.                                     |
|              | Web Crypto API | *(Native Browser API)*      | AES-256-GCM encryption with native performance—no need for heavy JS crypto libraries. |
| **Backend**  | Go (Golang)    | gin-gonic                   | High-performance HTTP web framework with minimal overhead.                            |
|              |                | aws-sdk-go-v2               | Efficient S3 streaming support for upload/download handling.                          |
| **Database** | Redis          | go-redis                    | Atomic counters (DECR) for download limits + TTL for automatic expiry.                |
| **Storage**  | Object Storage | AWS S3 (Prod) / MinIO (Dev) | Encrypted binary blob storage; MinIO replicates S3 behavior locally.                  |
| **DevOps**   | Docker         | docker-compose              | Orchestrates App + Redis + MinIO for a unified development environment.               |
