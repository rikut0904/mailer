"use client";

import { useState, useEffect } from "react";
import { getAuth, type User } from "@/lib/firebase";

export function useAuth() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let unsubscribe: (() => void) | undefined;

    (async () => {
      const { onAuthStateChanged } = await import("firebase/auth");
      const auth = await getAuth();
      unsubscribe = onAuthStateChanged(auth, (u) => {
        setUser(u);
        setLoading(false);
      });
    })();

    return () => {
      unsubscribe?.();
    };
  }, []);

  const signIn = async (email: string, password: string) => {
    const { signInWithEmailAndPassword } = await import("firebase/auth");
    const auth = await getAuth();
    return signInWithEmailAndPassword(auth, email, password);
  };

  const signUp = async (email: string, password: string) => {
    const { createUserWithEmailAndPassword } = await import("firebase/auth");
    const auth = await getAuth();
    return createUserWithEmailAndPassword(auth, email, password);
  };

  const signOut = async () => {
    const { signOut: fbSignOut } = await import("firebase/auth");
    const auth = await getAuth();
    return fbSignOut(auth);
  };

  return { user, loading, signIn, signUp, signOut };
}
