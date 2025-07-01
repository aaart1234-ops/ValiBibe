import { useParams } from 'react-router-dom'
import { useGetNoteQuery, useUpdateNoteMutation } from '@/features/note/noteApi'
import {
    TextField,
    Button,
    CircularProgress,
    Box,
    Snackbar,
    Alert,
    Typography,
    IconButton,
} from '@mui/material'
import { useEffect, useRef, useState } from 'react'
import EditIcon from '@mui/icons-material/Edit'
import { useAppSelector } from '@/app/hooks'

import { useEditor, EditorContent } from '@tiptap/react'
import StarterKit from '@tiptap/starter-kit'
import Underline from '@tiptap/extension-underline'
import { RichTextEditor } from '@mantine/tiptap'
import Link from '@tiptap/extension-link'
import Highlight from '@tiptap/extension-highlight'
import Image from '@tiptap/extension-image'

const NoteDetailPage = () => {
    const { id } = useParams<{ id: string }>()
    const { token } = useAppSelector((state) => state.auth)

    const { data: note, isLoading, error, refetch } = useGetNoteQuery(id!, { skip: !id })
    const [updateNote, { isLoading: isSaving }] = useUpdateNoteMutation()

    const [title, setTitle] = useState('')
    const [isEditing, setIsEditing] = useState(false)
    const [showSuccess, setShowSuccess] = useState(false)
    const titleRef = useRef<HTMLInputElement>(null)

    // Инициализируем редактор без начального контента, заполним позже
    const editor = useEditor({
        extensions: [
            StarterKit,
            Underline,
            Link,
            Highlight,
            Image,
        ],
        content: '', // пустая строка, заполним позже
        editable: false, // старт не в режиме редактирования
    })

    // Когда заметка загружена - обновляем title и содержимое редактора
    useEffect(() => {
        if (note && editor) {
            setTitle(note.title)
            // setContent в редакторе - аккуратно обновляем, не сломать undo stack
            editor.commands.setContent(note.content || '', false)
        }
    }, [note, editor])

    // При переключении режима редактирования меняем доступность редактора
    useEffect(() => {
        if (editor) {
            editor.setEditable(isEditing)
            if (isEditing) {
                titleRef.current?.focus()
            }
        }
    }, [isEditing, editor])

    // При смене токена обновляем заметку с сервера
    useEffect(() => {
        if (token) refetch()
    }, [token, refetch])

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

    if (isLoading) return <CircularProgress />
    if (error || !note) return <div>Ошибка загрузки заметки (не авторизован или заметка не найдена)</div>

    return (
        <Box display="flex" flexDirection="column" gap={2} p={2}>
            <Snackbar open={showSuccess} autoHideDuration={3000} onClose={() => setShowSuccess(false)}>
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
                    {/* Безопасно проверяем редактор */}
                    <div dangerouslySetInnerHTML={{ __html: editor ? editor.getHTML() : '' }} />
                </Box>
            ) : (
                <RichTextEditor editor={editor} style={{ minHeight: 200, height: 200 }}>
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
                    <strong>Уровень запоминания:</strong> {note.memoryLevel}
                </p>
                <p>
                    <strong>Следующее повторение:</strong> {note.next_review_at}
                </p>
            </Box>

            {isEditing && (
                <Box display="flex" gap={1}>
                    <Button variant="contained" onClick={handleSubmit} disabled={isSaving}>
                        {isSaving ? <CircularProgress size={20} /> : 'Сохранить'}
                    </Button>
                    <Button variant="outlined" onClick={() => setIsEditing(false)}>
                        Отмена
                    </Button>
                </Box>
            )}
        </Box>
    )
}

export default NoteDetailPage
