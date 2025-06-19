import { useParams } from 'react-router-dom'
import { useGetNoteQuery, useUpdateNoteMutation } from '@/features/note/noteApi'
import { TextField, Button, CircularProgress, Box, Snackbar, Alert, Typography, IconButton } from '@mui/material'
import { useEffect, useState } from 'react'
import EditIcon from '@mui/icons-material/Edit'
import { useAppSelector } from "@/app/hooks"

const NoteDetailPage = () => {
    const { id } = useParams<{ id: string }>()
    const { token } = useAppSelector(state => state.auth)

    // id может быть undefined, поэтому skip нужен
    const { data: note, isLoading, error, refetch } = useGetNoteQuery(id!, {
        skip: !id,
    })

    const [updateNote, { isSuccess }] = useUpdateNoteMutation()
    const [showSuccess, setShowSuccess] = useState(false)
    const [isEditing, setIsEditing] = useState(false)

    const [title, setTitle] = useState('')
    const [content, setContent] = useState('')

    useEffect(() => {
        if (note) {
            setTitle(note.title)
            setContent(note.content)
        }
    }, [note])

    useEffect(() => {
        if (token) {
            refetch()
        }
    }, [token])

    useEffect(() => {
        if (isSuccess) {
            setShowSuccess(true)
            setTimeout(() => setShowSuccess(false), 3000)
        }
    }, [isSuccess])

    if (isLoading) return <CircularProgress />
    if (error || !note) return <div>Ошибка загрузки заметки (не авторизован или заметка не найдена)</div>

    const handleSubmit = async () => {
        try {
            await updateNote({ id: note.id, title, content }).unwrap()
            setIsEditing(false)
        } catch (err) {
            console.error('Ошибка обновления заметки', err)
        }
    }

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
                    onChange={(e) => setTitle(e.target.value)}
                    fullWidth
                    autoFocus
                />
            )}

            {!isEditing ? (
                <Box>
                    <Typography variant="body1" whiteSpace="pre-line">{content}</Typography>
                </Box>
            ) : (
                <TextField
                    label="Текст"
                    multiline
                    rows={6}
                    value={content}
                    onChange={(e) => setContent(e.target.value)}
                    fullWidth
                />
            )}

            <p><strong>Уровень запоминания:</strong> {note.memoryLevel}</p>
            <p><strong>Следующее повторение:</strong> {note.next_review_at}</p>

            {isEditing && (
                <Box display="flex" gap={1}>
                    <Button variant="contained" onClick={handleSubmit}>Сохранить</Button>
                    <Button variant="outlined" onClick={() => setIsEditing(false)}>Отмена</Button>
                </Box>
            )}
        </Box>
    )
}

export default NoteDetailPage
