"use client"

import Head from "next/head"
import { Menu } from "lucide-react"

interface HeaderProps {
  onToggleSidebar: () => void;
  refresh: () => void;
}

function Header({ onToggleSidebar, refresh }: HeaderProps) {
    
  const createBlankNote = async () => {
      try {
        const res = await fetch("http://localhost:8000/notes", {
          method: "POST"
        })
        if (res.ok) {
          refresh()
        }
      } catch (err) {
        console.error("Delete failed", err)
      }
    }

  return (
    <>
      <Head>
        <title>Shared Notes</title>
      </Head>
      <div className="flex items-center justify-between bg-gray-50 text-black px-6 py-4 shadow-md">
        <div className="flex items-center gap-4">
          <button
            className="bg-gray-200 hover:bg-gray-300 text-black px-3 py-2 rounded"
            onClick={onToggleSidebar}
          >
            <Menu size={20} />
          </button>
          <h1 className="text-2xl font-bold">Shared Notes</h1>
        </div>

        <button className="bg-gray-500 hover:bg-gray-900 text-white text-xl font-bold px-4 py-2 rounded" onClick={createBlankNote}>
          + Add Note
        </button>
      </div>
    </>
  )
}

export default Header
