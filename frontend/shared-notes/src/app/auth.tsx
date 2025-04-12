"use client"

import { useState } from "react"

type AuthGateProps = {
  password: string
  children: React.ReactNode
}

export default function AuthGate({ password, children }: AuthGateProps) {
  const [input, setInput] = useState("")
  const [unlocked, setUnlocked] = useState(false)

  if (!unlocked) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-900 text-white">
        <div className="bg-gray-800 p-8 rounded-lg shadow-md w-80 flex flex-col gap-4">
          <h2 className="text-2xl font-bold text-center">Enter Password</h2>
          <input
            type="password"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            className="p-2 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Password"
          />
          <button
            onClick={() => {
              if (input === password) {
                setUnlocked(true)
              } else {
                alert("Incorrect password")
              }
            }}
            className="bg-blue-600 hover:bg-blue-700 text-white py-2 rounded font-semibold transition"
          >
            Unlock
          </button>
        </div>
      </div>
    )
  }

  return <>{children}</>
}
