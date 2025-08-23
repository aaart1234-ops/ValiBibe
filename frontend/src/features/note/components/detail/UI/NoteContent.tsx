import { Box } from '@mui/material'
import { RichTextEditor } from '@mantine/tiptap'
import { Editor, EditorContent } from '@tiptap/react'

interface NoteContentProps {
    editor: Editor | null
    isEditing: boolean
}

export const NoteContent = ({ editor, isEditing }: NoteContentProps) => {
    if (!editor) return null

    return !isEditing ? (
        <Box sx={{ border: '1px solid #ccc', borderRadius: 1, p: 1 }}>
            <div dangerouslySetInnerHTML={{ __html: editor.getHTML() }} />
        </Box>
    ) : (
        <RichTextEditor editor={editor} style={{ minHeight: 200, height: 'auto', padding: '10px' }}>
            <RichTextEditor.Toolbar sticky stickyOffset={60}>
                <RichTextEditor.Bold />
                <RichTextEditor.Italic />
                <RichTextEditor.Underline />
                <RichTextEditor.H1 />
                <RichTextEditor.H2 />
                <RichTextEditor.Link />
                <RichTextEditor.Highlight />
                <RichTextEditor.BulletList />
                <RichTextEditor.OrderedList />
                <RichTextEditor.Blockquote />
                <RichTextEditor.ClearFormatting />
            </RichTextEditor.Toolbar>
            <EditorContent editor={editor} />
        </RichTextEditor>
    )
}
