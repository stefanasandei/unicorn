# Secret Key Rotation

This document describes the key rotation functionality implemented in the secret manager service.

## Overview

The secret manager now supports automatic key rotation, which allows users to rotate their encryption keys while maintaining access to their existing secrets. This feature enhances security by regularly changing encryption keys.

## Key Features

### 1. Key Versioning

- Each user has multiple key versions
- Keys are derived from user ID and version number
- Only one key version is active at a time
- Old keys are retained for decryption of existing secrets

### 2. Automatic Re-encryption

- When keys are rotated, all existing secrets are automatically re-encrypted with the new key
- The process is transparent to users
- No data loss occurs during rotation

### 3. Key Management

- Keys are cached in memory for performance
- Key versions are stored in the database
- Key hashes are stored for audit purposes

## API Endpoints

### Rotate Keys

```
POST /api/v1/secrets/rotate-keys
```

Rotates all encryption keys for the authenticated user. This will:

1. Create a new key version
2. Re-encrypt all existing secrets with the new key
3. Deactivate the old key version

**Response:**

```json
{
  "message": "Keys rotated successfully",
  "user_id": "user-uuid"
}
```

### Get Key Versions

```
GET /api/v1/secrets/key-versions
```

Returns all key versions for the authenticated user.

**Response:**

```json
[
  {
    "id": "key-version-uuid",
    "user_id": "user-uuid",
    "version": 2,
    "key_hash": "hash-of-key",
    "created_at": "2024-01-01T00:00:00Z",
    "expires_at": null,
    "is_active": true
  },
  {
    "id": "key-version-uuid",
    "user_id": "user-uuid",
    "version": 1,
    "key_hash": "hash-of-key",
    "created_at": "2024-01-01T00:00:00Z",
    "expires_at": null,
    "is_active": false
  }
]
```

## Database Schema

### KeyVersion Table

```sql
CREATE TABLE key_versions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    version INTEGER NOT NULL,
    key_hash TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    expires_at DATETIME,
    is_active BOOLEAN NOT NULL DEFAULT true
);
```

### Updated Secret Table

The `secrets` table now includes a `key_version` field:

```sql
ALTER TABLE secrets ADD COLUMN key_version INTEGER NOT NULL DEFAULT 1;
```

## Security Considerations

### Key Derivation

- Keys are derived from user ID and version number using SHA-256
- In production, use a proper Key Management Service (KMS)
- Keys are never stored in plain text

### Access Control

- Key rotation requires write permissions
- Key version viewing requires read permissions
- All operations are user-scoped

### Audit Trail

- Key hashes are stored for audit purposes
- Key creation timestamps are recorded
- Active/inactive status is tracked

## Implementation Details

### KeyManager

The `KeyManager` class handles:

- Key version creation and management
- Key caching for performance
- Automatic key rotation
- Re-encryption of existing secrets

### SecretStore Integration

The `SecretStore` has been updated to:

- Use the KeyManager for all encryption/decryption
- Track key versions for each secret
- Support key rotation operations

### Backward Compatibility

- Existing secrets without key versions are treated as version 1
- Automatic migration creates initial key versions
- No data migration required for existing secrets

## Usage Examples

### Rotating Keys via API

```bash
curl -X POST http://localhost:8080/api/v1/secrets/rotate-keys \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Checking Key Versions

```bash
curl -X GET http://localhost:8080/api/v1/secrets/key-versions \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Testing

Run the key manager tests:

```bash
go test ./internal/stores -v -run TestKeyManager
```

## Future Enhancements

1. **Scheduled Rotation**: Automatic key rotation based on time intervals
2. **Key Expiration**: Automatic deactivation of old keys after a period
3. **Audit Logging**: Detailed logs of key rotation events
4. **KMS Integration**: Integration with external Key Management Services
5. **Bulk Operations**: Support for rotating keys across multiple users
