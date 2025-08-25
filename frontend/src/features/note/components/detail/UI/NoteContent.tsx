import { Box, Button } from '@mui/material'
import { RichTextEditor } from '@mantine/tiptap'
import { Editor, EditorContent } from '@tiptap/react'
import { useState } from 'react'

interface NoteContentProps {
    editor: Editor | null
    isEditing: boolean
}

export const NoteContent = ({ editor, isEditing }: NoteContentProps) => {
    const [expanded, setExpanded] = useState(false)

    if (!editor) return null

    const fullText = editor.getHTML()
    const MAX_PREVIEW = 300
    const isLong = fullText.length > MAX_PREVIEW

    // HTML для отображения (с форматированием)
    const fullHtml = editor.getHTML()

    // превью: обрезаем текст, но оставляем HTML "как есть"
    const truncatedHtml = (() => {
        if (!isLong) return fullHtml

        // берём текст до лимита
        const textPreview = fullText.slice(0, MAX_PREVIEW) + '…'
        // упрощённо вставляем в <p>, можно усложнить парсером HTML
        return `<p>${textPreview}</p>`
    })()

    return !isEditing ? (
        <Box>
            <div
                dangerouslySetInnerHTML={{
                    __html: expanded ? fullHtml : truncatedHtml,
                }}
            />
            {isLong && (
                <Button
                    size="small"
                    onClick={() => setExpanded(!expanded)}
                    sx={{ mt: 1 }}
                >
                    {expanded ? 'Скрыть' : 'Показать всё'}
                </Button>
            )}
        </Box>
    ) : (
        <RichTextEditor editor={editor} style={{ height: 'auto', padding: '0' }}>
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
            <EditorContent editor={editor} style={{ padding: '10px' }} />
        </RichTextEditor>
    )
}
