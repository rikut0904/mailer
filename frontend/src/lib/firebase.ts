import { type FirebaseApp } from "firebase/app";
import { type Auth, type User } from "firebase/auth";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

interface FirebaseClientConfig {
  firebase_api_key: string;
  firebase_auth_domain: string;
  firebase_project_id: string;
}

let configCache: FirebaseClientConfig | undefined;
let app: FirebaseApp | undefined;
let auth: Auth | undefined;

async function fetchConfig(): Promise<FirebaseClientConfig> {
  if (!configCache) {
    const res = await fetch(`${API_URL}/api/config`);
    if (!res.ok) throw new Error("Failed to fetch config");
    configCache = await res.json();
  }
  return configCache!;
}

async function getFirebaseApp(): Promise<FirebaseApp> {
  if (!app) {
    const [{ initializeApp, getApps }, cfg] = await Promise.all([
      import("firebase/app"),
      fetchConfig(),
    ]);
    const firebaseConfig = {
      apiKey: cfg.firebase_api_key,
      authDomain: cfg.firebase_auth_domain,
      projectId: cfg.firebase_project_id,
    };
    app =
      getApps().length === 0
        ? initializeApp(firebaseConfig)
        : getApps()[0];
  }
  return app;
}

async function getFirebaseAuth(): Promise<Auth> {
  if (!auth) {
    const { getAuth } = await import("firebase/auth");
    const fbApp = await getFirebaseApp();
    auth = getAuth(fbApp);
  }
  return auth;
}

export { getFirebaseAuth as getAuth };
export type { User };
