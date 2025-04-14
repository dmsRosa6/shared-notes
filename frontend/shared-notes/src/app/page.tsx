"use client"

import { useEffect, useState } from "react"
import RichEditor from "./editor"
import Header from "./header"
import NameList from "./fileTree"
import AuthGate from "./auth"

interface Note {
  _id: string;
  _t: string;
  _c?: string;
}

export default function Home() {
  const [showSidebar, setShowSidebar] = useState(true)
  const [notes, setNotes] = useState<Note[]>([])

  const fetchNotes = async () => {
    try {
      const res = await fetch("http://localhost:8000/notes")
      const data = await res.json()
      setNotes(data)
    } catch (err) {
      console.error("Failed to load notes", err)
    }
  }

  useEffect(() => {
    fetchNotes()
  }, [])

  const pass = "secret"
  return (
//    <AuthGate password={pass}>
      <div className="min-h-screen flex flex-col bg-gray-200">
        <Header onToggleSidebar={() => setShowSidebar(!showSidebar)} refresh={fetchNotes} />

        <div className="flex flex-grow justify-start p-4 gap-6 overflow-hidden">
          {showSidebar && (
            <div className="bg-gray-100 border border-gray-300 rounded-xl p-4 w-72 shadow-sm flex-shrink-0">
              <div className="bg-white rounded-lg border border-gray-200 p-3 max-h-64 overflow-y-auto">
                <NameList notes={notes} refresh={fetchNotes} />
              </div>
            </div>
          )}
          <div className="flex-grow bg-gray-500 text-white p-4 rounded-xl shadow-inner overflow-auto flex">
            <RichEditor />
          </div>
        </div>
      </div>
//    </AuthGate>
  )
}
