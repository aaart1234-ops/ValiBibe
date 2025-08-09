import { useParams, useNavigate } from 'react-router-dom'
import {
    useGetNoteQuery,
    useUpdateNoteMutation,
    useDeleteNoteMutation,
    useArchiveNoteMutation,
    useUnarchiveNoteMutation,
} from '@/features/note/noteApi'
import {
    TextField,
    Button,
    CircularProgress,
    Box,
    Fab,
    Snackbar,
    Alert,
    Typography,
    IconButton,
    Dialog,
    DialogActions,
    DialogContent,
    DialogContentText,
    DialogTitle,
} from '@mui/material'
import EditIcon from '@mui/icons-material/Edit'
import DeleteIcon from '@mui/icons-material/Delete'
import ArchiveIcon from '@mui/icons-material/Archive'
import ArrowBackIcon from '@mui/icons-material/ArrowBack'

import { useAppSelector } from '@/app/hooks'
import { useEditor, EditorContent } from '@tiptap/react'
import StarterKit from '@tiptap/starter-kit'
import Underline from '@tiptap/extension-underline'
import Link from '@tiptap/extension-link'
import Highlight from '@tiptap/extension-highlight'
import Image from '@tiptap/extension-image'
import { RichTextEditor } from '@mantine/tiptap'
import { useEffect, useRef, useState } from 'react'

const NoteDetailPage = () => {
    const { id } = useParams<{ id: string }>()
    const navigate = useNavigate()
    const { token } = useAppSelector((state) => state.auth)

    const [wasDeleted, setWasDeleted] = useState(false)
    const [updateNote, { isLoading: isSaving }] = useUpdateNoteMutation()
    const [deleteNote] = useDeleteNoteMutation()
    const [archiveNote] = useArchiveNoteMutation()
    const [unarchiveNote, { isLoading: isUnarchiving }] = useUnarchiveNoteMutation()

    const [title, setTitle] = useState('')
    const [isEditing, setIsEditing] = useState(false)
    const [showSuccess, setShowSuccess] = useState(false)
    const [confirmDialog, setConfirmDialog] = useState<'delete' | 'archive' | null>(null)

    const titleRef = useRef<HTMLInputElement>(null)

    const editor = useEditor({
        extensions: [StarterKit, Underline, Link, Highlight, Image],
        content: '',
        editable: false,
    })

    const { data: note, isLoading, error } = useGetNoteQuery(id!, {
        skip: !id || wasDeleted,
    })

    const [showArchiveSnackbar, setShowArchiveSnackbar] = useState(false)
    const undoRef = useRef(false) // Для отслеживания Undo

    useEffect(() => {
        if (note && editor) {
            setTitle(note.title)
            editor.commands.setContent(note.content || '', false)
        }
    }, [note, editor])

    useEffect(() => {
        if (editor) {
            editor.setEditable(isEditing)
            if (isEditing) titleRef.current?.focus()
        }
    }, [isEditing, editor])

    useEffect(() => {
        if (wasDeleted) {
            navigate('/notes')
        }
    }, [wasDeleted, navigate])

    const handleSubmit = async () => {
        if (!title.trim()) return alert('Введите заголовок')
        const content = editor?.getHTML() || ''
        if (!content.trim()) return alert('Введите текст заметки')

        try {
            await updateNote({ id: note!.id, title, content }).unwrap()
            setShowSuccess(true)
            setIsEditing(false)
        } catch (err) {
            console.error('Ошибка обновления заметки', err)
        }
    }

    const handleDelete = async () => {
        try {
            setWasDeleted(true)
            await deleteNote(note!.id).unwrap()
        } catch (err) {
            console.error('Ошибка удаления заметки', err)
            setWasDeleted(false)
        }
    }

    const handleArchive = async () => {
        try {
            undoRef.current = false
            await archiveNote(note!.id).unwrap()
            setShowArchiveSnackbar(true)

            // Ждем 5 секунд — если не было Undo, уходим
            setTimeout(() => {
                if (!undoRef.current) {
                    navigate('/notes')
                }
            }, 5000)
        } catch (err) {
            console.error('Ошибка архивирования', err)
        }
    }

    const handleUndoArchive = async () => {
        try {
            undoRef.current = true
            await unarchiveNote(note!.id).unwrap()
            setShowArchiveSnackbar(false)
        } catch (err) {
            console.error('Ошибка при Undo архивирования', err)
        }
    }


    if (isLoading) return <CircularProgress />
    if (error || !note) return <div>Ошибка загрузки заметки</div>

    return (
        <Box display="flex" flexDirection="column" gap={2} p={2}>
            <Snackbar open={showSuccess} autoHideDuration={3000} onClose={() => setShowSuccess(false)}>
                <Alert severity="success">Заметка сохранена</Alert>
            </Snackbar>

            <Snackbar
                open={showArchiveSnackbar}
                autoHideDuration={5000}
                onClose={() => setShowArchiveSnackbar(false)}
                message="Заметка архивирована"
                action={
                    <Button color="secondary" size="small" onClick={handleUndoArchive}>
                        Отмена
                    </Button>
                }
            />

            <Dialog open={!!confirmDialog} onClose={() => setConfirmDialog(null)}>
                <DialogTitle>
                    {confirmDialog === 'delete' ? 'Удалить заметку?' : 'Архивировать заметку?'}
                </DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        {confirmDialog === 'delete'
                            ? 'Вы уверены, что хотите удалить эту заметку? Это действие необратимо.'
                            : 'После архивирования заметка будет скрыта из основного списка.'}
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setConfirmDialog(null)}>Отмена</Button>
                    <Button
                        onClick={confirmDialog === 'delete' ? handleDelete : handleArchive}
                        color={confirmDialog === 'delete' ? 'error' : 'primary'}
                    >
                        Подтвердить
                    </Button>
                </DialogActions>
            </Dialog>

            <Box sx={{ pl: 6 }}>
                {!isEditing ? (
                    <Box display="flex" alignItems="center" gap={1}>
                        <Typography variant="h5">{title}</Typography>
                        <IconButton onClick={() => setIsEditing(true)} size="small">
                            <EditIcon fontSize="small" />
                        </IconButton>
                    </Box>
                ) : (
                    <TextField
                        label="Заголовок"
                        value={title}
                        inputRef={titleRef}
                        onChange={(e) => setTitle(e.target.value)}
                        fullWidth
                    />
                )}
            </Box>

            {!isEditing ? (
                <Box sx={{ border: '1px solid #ccc', borderRadius: 2, p: 2 }}>
                    <div dangerouslySetInnerHTML={{ __html: editor ? editor.getHTML() : '' }} />
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
            )}

            <Box mt={4}>
                <p>
                    <strong>🧠 Уровень запоминания:</strong> {note.memoryLevel}
                </p>
            </Box>

            {isEditing ? (
                <Box display="flex" gap={1}>
                    <Button variant="contained" onClick={handleSubmit} disabled={isSaving}>
                        {isSaving ? <CircularProgress size={20} /> : 'Сохранить'}
                    </Button>
                    <Button variant="outlined" onClick={() => setIsEditing(false)}>
                        Отмена
                    </Button>
                </Box>
            ) : (
                <Box display="flex" gap={1}>
                    <Button
                        variant="outlined"
                        color="error"
                        startIcon={<DeleteIcon />}
                        onClick={() => setConfirmDialog('delete')}
                    >
                        Удалить
                    </Button>
                    {!note?.archived && (
                        <Button
                            variant="outlined"
                            startIcon={<ArchiveIcon />}
                            onClick={handleArchive}
                        >
                            Архивировать
                        </Button>
                    )}
                    {note?.archived && (
                        <Button
                            variant="contained"
                            color="secondary"
                            onClick={async () => {
                                try {
                                    await unarchiveNote(note.id).unwrap()
                                    navigate('/notes') // или refetch, если ты не хочешь редирект
                                } catch (e) {
                                    console.error('Ошибка при разархивировании:', e)
                                }
                            }}
                            disabled={isUnarchiving}
                        >
                            {isUnarchiving ? 'Восстановление...' : 'Из архива'}
                        </Button>
                    )}
                </Box>
            )}

            <Fab
                color="primary"
                onClick={() => navigate(-1)}
                sx={{
                    position: 'fixed',
                    top: 76,
                    left: 12,
                    zIndex: 1000,
                    width: 40,
                    height: 40,
                }}
            >
                <ArrowBackIcon />
            </Fab>

        </Box>
    )
}

export default NoteDetailPage
