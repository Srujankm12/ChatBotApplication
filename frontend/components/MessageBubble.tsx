"use client";

interface Props {
  role: "user" | "assistant";
  content: string;
}

export default function MessageBubble({ role, content }: Props) {
  const isUser = role === "user";

  return (
    <div className={`flex ${isUser ? "justify-end" : "justify-start"}`}>
      <div
        className={`
          max-w-[75%] px-4 py-2 rounded-xl
          text-sm leading-relaxed whitespace-pre-wrap break-words
          ${isUser
            ? "bg-blue-500/10 text-gray-900 border border-blue-100"
            : "bg-gray-100 text-gray-800"
          }
        `}
      >
        {content}
      </div>
    </div>
  );
}
