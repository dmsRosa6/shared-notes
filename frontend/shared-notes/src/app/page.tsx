"use client"

import { useState } from "react"
import TxtEditor from "./txtEditor"
import Header from "./header"
import NameList from "./fileTree"
import AuthGate from "./auth"

export default function Home() {
  const [showSidebar, setShowSidebar] = useState(true)
  const pass = "secret"
  return (
    <AuthGate password={pass}>
      <div className="min-h-screen flex flex-col bg-gray-200">
        <Header onToggleSidebar={() => setShowSidebar(!showSidebar)} />

        <div className="flex flex-grow justify-start p-4 gap-6 overflow-hidden">
          {showSidebar && (
            <div className="bg-gray-100 border border-gray-300 rounded-xl p-4 w-72 shadow-sm flex-shrink-0">
              <div className="bg-white rounded-lg border border-gray-200 p-3 max-h-64 overflow-y-auto">
                <NameList />
              </div>
            </div>
          )}

          <div className="flex-grow bg-gray-700 text-white p-4 rounded-xl shadow-inner overflow-auto">
            <TxtEditor />
          </div>
        </div>
      </div>
    </AuthGate>
  )
}
