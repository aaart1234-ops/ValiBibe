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
import ReactQuill from 'react-quill'
import 'react-quill/dist/quill.snow.css'
import { useCreateNoteMutation } from '@/features/note/noteApi'

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
        <Dialog open={open} onClose={onClose} fullWidth maxWidth="md" onKeyDown={handleKeyDown} >
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
                    <ReactQuill
                        value={content}
                        onChange={setContent}
                        theme="snow"
                        style={{ minHeight: '200px', height: '200px' }}
                        placeholder="Содержание..."
                    />
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
