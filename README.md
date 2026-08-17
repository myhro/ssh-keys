SSH Keys
========

Sync SSH `authorized_keys` from a remote URL. Defaults to GitHub's `https://github.com/{USER}.keys`, but it can be any HTTPS URL.

Configurable via environment variables:

| Config | Default |
| :-- | :-- |
| `SSH_KEYS_URL` | `https://github.com/{USER}.keys` |
| `SSH_KEYS_USER` | Current `$USER` |
| `SSH_KEYS_FILE` | `~/.ssh/authorized_keys` |

Tries to be as safe as possible:

- Checks ETag headers to avoid downloading keys every time.
- Ensures fetched keys are valid.
- Detects local modifications and restores the keys from the remote.
- Only overwrites the target file when at least one valid key was retrieved.
- Writes are atomic, ensuring no partial/corrupted `authorized_keys` file.
