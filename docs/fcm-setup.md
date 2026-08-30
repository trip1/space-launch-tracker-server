# Android FCM deployment

## One-time Firebase setup

1. Create or select the Firebase project for **Kotlin Space Launches**.
2. Register Android package `com.ds9.kotlin_space_launches`.
3. Download `google-services.json` and place it at `kotlin_space_launches/composeApp/google-services.json`. This client configuration is required for Android builds but is not committed.
4. Create a dedicated service account with only Firebase Cloud Messaging sending capability. Download its JSON key into protected deployment storage, outside the repository.

## Production server configuration

Set these protected production environment values:

```dotenv
FCM_ENABLED=true
FCM_PROJECT_ID=<Firebase-project-id>
FCM_CREDENTIALS_HOST_PATH=/protected/path/firebase-service-account.json
FCM_CREDENTIALS_FILE=/run/secrets/firebase-service-account.json
NOTIFICATION_POLL_INTERVAL=5m
```

`compose.production.yml` mounts the service-account file read-only. The Go server exits at startup if FCM is enabled but the sender cannot initialize.

## Runtime contract

- Android registers its FCM token at launch and when Firebase rotates it.
- `POST /v1/devices/fcm` returns a random per-installation device credential. Future rotations must supply that credential in `X-Device-Secret`; it is stored only on the device and as a SHA-256 hash on the server.
- The backend polls its canonical launch feed. First poll creates a baseline; subsequent changes to a launch NET or status produce a high-priority data-only FCM message containing only `event`, `launch_id`, and `kind`.
- The app fetches canonical launch data in `LaunchPushWorker`, then posts a deep-linked local notification. It never displays title/body text supplied by FCM.

## Release proof

On a physical Android device: install, grant notification permission, verify device registration returns 201 without printing the token, change a test launch revision, observe the FCM send accepted by Firebase, then verify one native notification opens the right launch detail. Confirm push behavior after swipe-away; do not promise delivery after Android Settings force-stop.
