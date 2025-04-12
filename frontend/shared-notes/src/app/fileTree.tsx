import { Trash2 } from "lucide-react"

export default function NameList() {
  const names = ["Note 1", "Note 2", "Note 3", "Note 4", "Note 5", "Note 6", "Note 7", "Note 8", "Note 9", "Note 10"]

  return (
    <div className="text-black font-bold overflow-y-auto rounded-lg shadow-sm">
      {names.map((name, index) => (
        <div
          key={index}
          className={`flex items-center justify-between p-4 transition-colors ${
            index % 2 === 0 ? "bg-gray-50" : "bg-gray-100"
          } hover:bg-gray-200 border-b border-gray-200 text-lg`}
        >
          <span>{name}</span>
          <button className="text-gray-500 hover:text-red-600 transition">
            <Trash2 size={20} />
          </button>
        </div>
      ))}
    </div>
  )
}
