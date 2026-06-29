// Package config wires up application configuration and external service
// clients (Firebase/Firestore).
//
// Service-account credentials are intentionally not embedded into the binary.
// Provide them at runtime via the FIREBASE_CREDENTIALS environment variable
// (path to a service-account JSON file) or via Application Default Credentials.
package config
