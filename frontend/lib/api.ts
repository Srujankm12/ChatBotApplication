const API_URL = process.env.NEXT_PUBLIC_API_URL || "https://chatbotbackendd.onrender.com";

export async function sendMessage(
  message: string,
  sessionId: string,
  signal?: AbortSignal
): Promise<string> {
  const res = await fetch(`${API_URL}/chat`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ message, session_id: sessionId }),
    signal,
  });

  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: "Unknown error" }));
    throw new Error(err.error || `HTTP ${res.status}`);
  }

  const data = await res.json();
  return data.reply as string;
}

export async function fetchMessages(
  sessionId: string
): Promise<{ role: "user" | "assistant"; content: string }[]> {
  try {
    const res = await fetch(`${API_URL}/chat/${sessionId}`);
    if (!res.ok) return [];
    const data = await res.json();
    return Array.isArray(data.messages) ? data.messages : [];
  } catch {
    return [];
  }
}

export async function fetchSessions(): Promise<string[]> {
  try {
    const res = await fetch(`${API_URL}/sessions`);
    if (!res.ok) return [];
    const data = await res.json();
    return Array.isArray(data.sessions) ? data.sessions : [];
  } catch {
    return [];
  }
}
