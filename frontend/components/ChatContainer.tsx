"use client";

import { useState, useEffect, useRef } from "react";
import { v4 as uuidv4 } from "uuid";
import MessageBubble from "./MessageBubble";
import ChatInput from "./ChatInput";
import { sendMessage, fetchMessages, fetchSessions } from "@/lib/api";

interface Message {
  role: "user" | "assistant";
  content: string;
  id: string;
}

const SESSION_KEY = "chatbot_session_id";

export default function ChatContainer() {
  const [messages, setMessages] = useState<Message[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [sessionId, setSessionId] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [sessions, setSessions] = useState<string[]>([]);
  const [showSessions, setShowSessions] = useState(false);
  const bottomRef = useRef<HTMLDivElement>(null);
  const sessionsRef = useRef<HTMLDivElement>(null);
  const abortRef = useRef<AbortController | null>(null);

  useEffect(() => {
    let sid = localStorage.getItem(SESSION_KEY);
    if (!sid) {
      sid = uuidv4();
      localStorage.setItem(SESSION_KEY, sid);
    }
    setSessionId(sid);

    fetchMessages(sid).then((history) => {
      if (history.length > 0) {
        setMessages(history.map((m) => ({ ...m, id: uuidv4() })));
      }
    });
  }, []);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages, isLoading]);

  const handleSend = async (text: string) => {
    setError(null);
    const userMsg = { role: "user" as const, content: text, id: uuidv4() };
    setMessages((prev) => [...prev, userMsg]);
    setIsLoading(true);

    const controller = new AbortController();
    abortRef.current = controller;

    try {
      const reply = await sendMessage(text, sessionId, controller.signal);
      setMessages((prev) => [
        ...prev,
        { role: "assistant", content: reply, id: uuidv4() },
      ]);
    } catch (err) {
      if (err instanceof Error && err.name === "AbortError") {
        setMessages((prev) => prev.filter((m) => m.id !== userMsg.id));
      } else {
        setError(err instanceof Error ? err.message : "Something went wrong");
      }
    } finally {
      setIsLoading(false);
      abortRef.current = null;
    }
  };

  const handleStop = () => {
    abortRef.current?.abort();
  };

  const handleNewChat = () => {
    const newSid = uuidv4();
    localStorage.setItem(SESSION_KEY, newSid);
    setSessionId(newSid);
    setMessages([]);
    setError(null);
    setShowSessions(false);
  };

  const handleOpenSessions = async () => {
    const next = !showSessions;
    setShowSessions(next);
    if (next) {
      const list = await fetchSessions();
      setSessions(list);
    }
  };

  const handleSelectSession = async (sid: string) => {
    localStorage.setItem(SESSION_KEY, sid);
    setSessionId(sid);
    setMessages([]);
    setError(null);
    setShowSessions(false);
    const history = await fetchMessages(sid);
    if (history.length > 0) {
      setMessages(history.map((m) => ({ ...m, id: uuidv4() })));
    }
  };

  return (
    <div className="h-screen flex flex-col bg-gray-50">

      {/* Header — fixed, never scrolls */}
      <header className="flex-none bg-white border-b border-gray-200 px-4 py-3">
        <div className="max-w-[720px] mx-auto flex items-center justify-between">
          <div>
            <h1 className="text-sm font-semibold text-gray-800 tracking-tight">
              AI Chat
            </h1>
            <p className="text-xs text-gray-400 mt-0.5">Gemini 1.5 Flash</p>
          </div>
          <div className="flex items-center gap-2">
            <div className="relative" ref={sessionsRef}>
              <button
                onClick={handleOpenSessions}
                className="text-xs text-gray-500 border border-gray-200 rounded-md px-2.5 py-1.5 hover:bg-gray-50 hover:text-gray-700 transition-colors duration-150"
              >
                History
              </button>
              {showSessions && (
                <div className="absolute right-0 mt-1 w-56 bg-white border border-gray-200 rounded-md shadow-md z-10 max-h-64 overflow-y-auto">
                  {sessions.length === 0 ? (
                    <p className="text-xs text-gray-400 px-3 py-2">No sessions found</p>
                  ) : (
                    sessions.map((sid) => (
                      <button
                        key={sid}
                        onClick={() => handleSelectSession(sid)}
                        className={`block w-full text-left px-3 py-2 text-xs hover:bg-gray-50 transition-colors ${
                          sid === sessionId
                            ? "font-medium text-gray-800 bg-gray-50"
                            : "text-gray-500"
                        }`}
                      >
                        {sid.slice(0, 8)}…
                      </button>
                    ))
                  )}
                </div>
              )}
            </div>
            <button
              onClick={handleNewChat}
              className="text-xs text-gray-500 border border-gray-200 rounded-md px-2.5 py-1.5 hover:bg-gray-50 hover:text-gray-700 transition-colors duration-150"
            >
              New chat
            </button>
          </div>
        </div>
      </header>

      {/* Chat area — grows to fill, scrolls independently */}
      <main className="flex-1 overflow-y-auto scrollbar-thin">
        <div className="max-w-[720px] mx-auto px-4 py-6">

          {messages.length === 0 && !isLoading && (
            <p className="text-gray-400 text-sm text-center mt-24 select-none">
              Send a message to start.
            </p>
          )}

          <div className="space-y-3">
            {messages.map((msg) => (
              <MessageBubble
                key={msg.id}
                role={msg.role}
                content={msg.content}
              />
            ))}

            {isLoading && (
              <div className="flex justify-start">
                <span className="bg-gray-100 text-gray-400 rounded-xl px-4 py-2 text-sm select-none">
                  ...
                </span>
              </div>
            )}

            {error && (
              <p className="text-red-400 text-xs text-center py-1">{error}</p>
            )}
          </div>

          <div ref={bottomRef} />
        </div>
      </main>

      {/* Input bar — fixed, never scrolls */}
      <footer className="flex-none bg-white border-t border-gray-100 px-4 py-3">
        <div className="max-w-[720px] mx-auto">
          <ChatInput
            onSend={handleSend}
            onStop={handleStop}
            disabled={isLoading || !sessionId}
            isLoading={isLoading}
          />
          <p className="text-xs text-gray-300 text-center mt-2 select-none">
            Enter to send · Shift+Enter for new line
          </p>
        </div>
      </footer>

    </div>
  );
}
