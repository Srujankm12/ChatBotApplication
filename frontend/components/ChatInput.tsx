"use client";

import { useState, useRef, KeyboardEvent, useEffect } from "react";

interface Props {
  onSend: (message: string) => void;
  onStop: () => void;
  disabled: boolean;
  isLoading: boolean;
}

export default function ChatInput({ onSend, onStop, disabled, isLoading }: Props) {
  const [input, setInput] = useState("");
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  useEffect(() => {
    if (!disabled) textareaRef.current?.focus();
  }, [disabled]);

  const autoResize = () => {
    const el = textareaRef.current;
    if (!el) return;
    el.style.height = "auto";
    el.style.height = Math.min(el.scrollHeight, 96) + "px"; // 4 lines max
  };

  const handleSend = () => {
    const trimmed = input.trim();
    if (!trimmed || disabled) return;
    onSend(trimmed);
    setInput("");
    if (textareaRef.current) textareaRef.current.style.height = "auto";
  };

  const handleKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  const isEmpty = !input.trim();

  return (
    <div className="flex items-end gap-2 bg-white border border-gray-200 rounded-full px-4 py-2 ring-1 ring-transparent focus-within:ring-blue-200 transition-shadow duration-150">
      <textarea
        ref={textareaRef}
        rows={1}
        className="flex-1 resize-none bg-transparent text-sm text-gray-800 placeholder-gray-400 outline-none leading-5 py-0.5"
        placeholder="Message..."
        value={input}
        onChange={(e) => {
          setInput(e.target.value);
          autoResize();
        }}
        onKeyDown={handleKeyDown}
        disabled={disabled}
      />
      {isLoading ? (
        <button
          onClick={onStop}
          aria-label="Stop"
          className="flex-none mb-0.5 w-7 h-7 flex items-center justify-center rounded-full bg-gray-900 text-white hover:bg-gray-700 active:scale-95 transition-all duration-150"
        >
          <svg className="w-3 h-3" viewBox="0 0 24 24" fill="currentColor">
            <rect x="4" y="4" width="16" height="16" rx="2" />
          </svg>
        </button>
      ) : (
        <button
          onClick={handleSend}
          disabled={disabled || isEmpty}
          aria-label="Send"
          className="flex-none mb-0.5 w-7 h-7 flex items-center justify-center rounded-full bg-gray-900 text-white disabled:opacity-25 hover:bg-gray-700 active:scale-95 transition-all duration-150"
        >
          <svg
            className="w-3.5 h-3.5"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            strokeWidth={2.5}
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M6 12L3.269 3.126A59.768 59.768 0 0121.485 12 59.77 59.77 0 013.27 20.876L5.999 12zm0 0h7.5"
            />
          </svg>
        </button>
      )}
    </div>
  );
}
