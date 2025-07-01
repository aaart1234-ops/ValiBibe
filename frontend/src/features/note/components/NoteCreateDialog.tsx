// components/CreateNoteModal.tsx
import {
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    Button,
    TextField,
    Box,
    CircularProgress,
    IconButton
} from '@mui/material'
import CloseIcon from '@mui/icons-material/Close'
import React, { useEffect, useRef, useState } from 'react'
import { useCreateNoteMutation } from '@/features/note/noteApi'
import { RichTextEditor } from '@mantine/tiptap'
import { useEditor } from '@tiptap/react'
import StarterKit from '@tiptap/starter-kit'
import Underline from '@tiptap/extension-underline'
import Link from '@tiptap/extension-link'
import Highlight from '@tiptap/extension-highlight'

interface Props {
    open: boolean
    onClose: () => void
    onCreated: () => void
}

const CreateNoteModal: React.FC<Props> = ({ open, onClose, onCreated }) => {
    const [title, setTitle] = useState('')
    const [content, setContent] = useState('')
    const [error, setError] = useState('')
    const [createNote, { isLoading }] = useCreateNoteMutation()
    const titleRef = useRef<HTMLInputElement>(null)

    useEffect(() => {
        if (open) {
            setTimeout(() => {
                titleRef.current?.focus()
            }, 100)
        }
    }, [open])

    const editor = useEditor({
        extensions: [
            StarterKit,
            Underline,
            Link,
            Highlight,
            ],
        content,
        onUpdate: ({ editor }) => setContent(editor.getHTML()),
    })

    const handleSubmit = async () => {
        if (!title.trim()) {
            setError('Заголовок обязателен')
            return
        }
        setError('')

        try {
            await createNote({ title, content }).unwrap()
            onCreated()
            onClose()
            setTitle('')
            setContent('')
            editor?.commands.setContent('') // Очистка редактора
        } catch (e) {
            setError('Ошибка при создании заметки')
        }
    }

    const handleKeyDown = (e: React.KeyboardEvent) => {
        if (e.key === 'Escape') {
            onClose()
        }
    }

    return (
        <Dialog open={open} onClose={onClose} fullWidth maxWidth="md" onKeyDown={handleKeyDown}>
            <DialogTitle>
                Новая заметка
                <IconButton
                    aria-label="close"
                    onClick={onClose}
                    sx={{
                        position: 'absolute',
                        right: 8,
                        top: 8,
                        color: (theme) => theme.palette.grey[500],
                    }}
                >
                    <CloseIcon />
                </IconButton>
            </DialogTitle>
            <DialogContent dividers>
                <Box display="flex" flexDirection="column" gap={2}>
                    <TextField
                        label="Заголовок"
                        inputRef={titleRef}
                        value={title}
                        onChange={(e) => setTitle(e.target.value)}
                        fullWidth
                        error={!!error}
                        helperText={error}
                    />
                    <RichTextEditor editor={editor}>
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
                        <RichTextEditor.Content
                            style={{ minHeight: '200px', borderRadius: 8 }}
                        />
                    </RichTextEditor>
                </Box>
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose} disabled={isLoading}>
                    Отмена
                </Button>
                <Button onClick={handleSubmit} variant="contained" disabled={isLoading}>
                    {isLoading ? <CircularProgress size={24} /> : 'Сохранить'}
                </Button>
            </DialogActions>
        </Dialog>
    )
}

export default CreateNoteModal
