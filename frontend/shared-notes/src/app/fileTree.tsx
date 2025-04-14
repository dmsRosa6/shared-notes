"use client"

import { Trash2 } from "lucide-react"

interface Note {
  _id: string;
  _t: string;
  _c?: string;
}

interface NameListProps {
  notes: Note[];
  refresh: () => void;
}

async function deleteNote(id: string, refresh: () => void) {
  try {
    const res = await fetch("http://localhost:8000/notes/" + id, {
      method: "DELETE"
    })
    if (res.ok) {
      refresh()
    }
  } catch (err) {
    console.error("Delete failed", err)
  }
}

export default function NameList({ notes, refresh }: NameListProps) {
  if (!notes || notes.length === 0) {
    return <div className="text-gray-600 italic">No notes available</div>
  }

  return (
    <div className="text-black font-bold overflow-y-auto rounded-lg shadow-sm">
      {notes.map((note, index) => (
        <div
          key={note._id}
          className={`flex items-center justify-between p-4 transition-colors ${
            index % 2 === 0 ? "bg-gray-50" : "bg-gray-100"
          } hover:bg-gray-200 border-b border-gray-200 text-lg`}
        >
          <span>{note._t}</span>
          <button
            className="text-gray-500 hover:text-red-600 transition"
            onClick={() => deleteNote(note._id, refresh)}
          >
            <Trash2 size={20} />
          </button>
        </div>
      ))}
    </div>
  )
}
