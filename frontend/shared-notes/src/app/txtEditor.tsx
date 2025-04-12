"use client"

import {useEditor, EditorContent} from "@tiptap/react"
import StarterKit from "@tiptap/starter-kit"
import React, { useCallback, useEffect, useState } from 'react'

import Collaboration from "@tiptap/extension-collaboration"
import * as Y from "yjs"

const doc = new Y.Doc()

const TipTap = () => {
    const editor = useEditor({
        extensions: [
            StarterKit, 
            //Collaboration.configure({document: doc})
        ],
        content: "<p>Hello</p>",
        editorProps: {
            attributes:{
                class:
                "prose prose-sm"
            }
        }
    });
    
    if(!editor){
        return null
    }

    return(
        <>
            <div className="flex dlex-col">
                <div className="w-full p-2">
                    <h2>Editor</h2>
                    <EditorContent editor={editor} className="border rounded-lg p-2" />
                </div>
            </div>    
        </>
    );
};

export default TipTap;