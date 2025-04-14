'use client';

import { useEffect, useRef, useState } from 'react';
import {Bold, Italic} from 'lucide-react'

export default function RichEditor() {
    const editorRef = useRef<HTMLDivElement | null>(null);
    const socketRef = useRef<WebSocket | null>(null);
    const [isBold, setIsBold] = useState(false);
    const [isItalic, setIsItalic] = useState(false);
  
    useEffect(() => {
      const socket = new WebSocket('ws://localhost:8080/ws');
      socketRef.current = socket;
  
      socket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        if (data.type === 'update' && editorRef.current) {
          editorRef.current.innerHTML = data.html;
        }
      };
  
      return () => {
        socket.close();
      };
    }, []);
  
    const handleInput = () => {
      const html = editorRef.current?.innerHTML || '';
      const message = JSON.stringify({ type: 'update', html });
      socketRef.current?.send(message);
    };
  
    const makeBold = () => {
      const selection = window.getSelection();
      if (!selection || selection.rangeCount === 0) return;
  
      const range = selection.getRangeAt(0);
      const strong = document.createElement('strong');
      try {
        range.surroundContents(strong);
      } catch (err) {
        console.warn('Could not wrap selection:', err);
      }
  
      selection.removeAllRanges();
      setIsBold((prev) => !prev);
    };

    const makeItalic = () => {
        setIsItalic((prev) => !prev);
      };
  
    return (
      <div className="w-full h-full flex flex-col space-y-4">
        <div className="space-x-2">
          <button
            onClick={makeBold}
            className={`px-3 py-1 ${isBold ? 'bg-gray-900' : 'bg-gray-500'} text-white rounded`}
          >
            <Bold size={20} />
          </button>
          <button
            onClick={makeItalic}
            className={`px-3 py-1 ${isItalic ? 'bg-gray-900' : 'bg-gray-500'} text-white rounded`}
          >
            <Italic size={20} />
          </button>
        </div>
  
        <div
          ref={editorRef}
          onInput={handleInput}
          contentEditable
          suppressContentEditableWarning
          className="flex-grow overflow-y-auto border p-4 rounded bg-white text-black relative empty:before:content-[attr(data-placeholder)] before:text-gray-400 before:absolute before:top-4 before:left-4"
          data-placeholder="Start typing..."
        />
      </div>
    );
  }
  