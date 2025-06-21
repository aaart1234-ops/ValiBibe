import { useParams } from 'react-router-dom'
import { useGetNoteQuery, useUpdateNoteMutation } from '@/features/note/noteApi'
import { TextField, Button, CircularProgress, Box, Snackbar, Alert, Typography, IconButton } from '@mui/material'
import { useEffect, useRef, useState } from 'react'
import EditIcon from '@mui/icons-material/Edit'
import { useAppSelector } from "@/app/hooks"
import ReactQuill from 'react-quill'
import 'react-quill/dist/quill.snow.css'

const NoteDetailPage = () => {
    const { id } = useParams<{ id: string }>()
    const { token } = useAppSelector(state => state.auth)

    const { data: note, isLoading, error, refetch } = useGetNoteQuery(id!, { skip: !id })
    const [updateNote, { isLoading: isSaving }] = useUpdateNoteMutation()

    const [title, setTitle] = useState('')
    const [content, setContent] = useState('')
    const [isEditing, setIsEditing] = useState(false)
    const [showSuccess, setShowSuccess] = useState(false)
    const titleRef = useRef<HTMLInputElement>(null)

    useEffect(() => {
        if (note) {
            setTitle(note.title)
            setContent(note.content)
        }
    }, [note])

    useEffect(() => {
        if (token) refetch()
    }, [token])

    useEffect(() => {
        if (isEditing) titleRef.current?.focus()
    }, [isEditing])

    const handleSubmit = async () => {
        if (!title.trim()) return alert("Введите заголовок")
        if (!content.trim()) return alert("Введите текст заметки")

        try {
            await updateNote({ id: note!.id, title, content }).unwrap()
            setShowSuccess(true)
            setIsEditing(false)
        } catch (err) {
            console.error('Ошибка обновления заметки', err)
        }
    }

    if (isLoading) return <CircularProgress />
    if (error || !note) return <div>Ошибка загрузки заметки (не авторизован или заметка не найдена)</div>

    return (
        <Box display="flex" flexDirection="column" gap={2} p={2}>
            <Snackbar open={showSuccess} autoHideDuration={3000}>
                <Alert severity="success">Заметка сохранена</Alert>
            </Snackbar>

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

            {!isEditing ? (
                <Box sx={{ border: '1px solid #ccc', borderRadius: 2, p: 2 }}>
                    <div dangerouslySetInnerHTML={{ __html: content }} />
                </Box>
            ) : (
                <ReactQuill value={content} onChange={setContent}
                            style={{ minHeight: '300px', height: '300px' }}/>
            )}

            <Box mt={4}>
                <p><strong>Уровень запоминания:</strong> {note.memoryLevel}</p>
                <p><strong>Следующее повторение:</strong> {note.next_review_at}</p>
            </Box>

            {isEditing && (
                <Box display="flex" gap={1}>
                    <Button
                        variant="contained"
                        onClick={handleSubmit}
                        disabled={isSaving}
                    >
                        {isSaving ? <CircularProgress size={20} /> : 'Сохранить'}
                    </Button>
                    <Button variant="outlined" onClick={() => setIsEditing(false)}>Отмена</Button>
                </Box>
            )}
        </Box>
    )
}

export default NoteDetailPage
