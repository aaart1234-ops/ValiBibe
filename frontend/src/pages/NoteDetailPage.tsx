import { useParams } from 'react-router-dom'
import { CircularProgress, Box } from '@mui/material'
import { useEditor } from '@tiptap/react'
import StarterKit from '@tiptap/starter-kit'
import Underline from '@tiptap/extension-underline'
import Link from '@tiptap/extension-link'
import Highlight from '@tiptap/extension-highlight'
import Image from '@tiptap/extension-image'
import { useEffect, useRef, useState } from 'react'

import { useGetNoteQuery } from '@/features/note/noteApi'
import { useNoteActions } from '@/features/note/hooks/useNoteActions'
import { NoteTitle } from '@/features/note/components/detail/UI/NoteTitle'
import { NoteContent } from '@/features/note/components/detail/UI/NoteContent'
import { NoteActions } from '@/features/note/components/detail/NoteActions'
import { ConfirmDialog } from '@/features/note/components/detail/ConfirmDialog'
import { SnackbarSuccess } from '@/features/note/components/detail/SnackbarSuccess'
import { NoteMoreMenu } from '@/features/note/components/detail/UI/NoteMoreMenu'
import { useDetailPage } from '@/context/DetailPageContext'

const NoteDetailPage = () => {
    const { id } = useParams<{ id: string }>()
    const { isEditing, setEditing } = useDetailPage()

    const [title, setTitle] = useState('')
    const [confirmDialog, setConfirmDialog] = useState<'delete' | 'archive' | null>(null)
    const [snackbar, setSnackbar] = useState<{ message: string; actionText?: string; onAction?: () => void } | null>(null)

    const titleRef = useRef<HTMLInputElement>(null)

    const editor = useEditor({
        extensions: [StarterKit, Underline, Link, Highlight, Image],
        content: '',
        editable: false,
    })

    const { data: note, isLoading, error } = useGetNoteQuery(id!, { skip: !id })
    const { update, remove, archive, unarchive, isSaving, isUnarchiving } = useNoteActions(id!)
    const [moreAnchorEl, setMoreAnchorEl] = useState<null | HTMLElement>(null)
    const openMore = Boolean(moreAnchorEl)

    const handleOpenMore = (e: React.MouseEvent<HTMLButtonElement>) => setMoreAnchorEl(e.currentTarget)
    const handleCloseMore = () => setMoreAnchorEl(null)

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

    // При уходе со страницы — точно выходим из режима редактирования
    useEffect(() => {
        return () => setEditing(false)
    }, [setEditing])

    const handleSubmit = async () => {
        if (!title.trim()) return alert('Введите заголовок')
        const content = editor?.getHTML() || ''
        if (!content.trim()) return alert('Введите текст заметки')
        try {
            await update(title, content)
            setSnackbar({ message: 'Заметка сохранена' })
            setEditing(false)
        } catch (err) {
            console.error('Ошибка обновления заметки', err)
        }
    }

    if (isLoading) return <CircularProgress />
    if (error || !note) return <div>Ошибка загрузки заметки</div>

    return (
        <Box display="flex" flexDirection="column" gap={2} p={2}>
            <SnackbarSuccess
                open={!!snackbar}
                message={snackbar?.message || ''}
                onClose={() => setSnackbar(null)}
                actionText={snackbar?.actionText}
                onAction={snackbar?.onAction}
            />

            <ConfirmDialog
                type={confirmDialog}
                open={!!confirmDialog}
                onClose={() => setConfirmDialog(null)}
                onConfirm={confirmDialog === 'delete' ? remove : archive}
            />

            <NoteTitle
                title={title}
                isEditing={isEditing}
                onChange={setTitle}
                onEditToggle={setEditing}
                titleRef={titleRef}
                onMoreClick={handleOpenMore}
            />

            <NoteMoreMenu
                anchorEl={moreAnchorEl}
                open={openMore}
                onClose={handleCloseMore}
                noteArchived={!!note.archived}
                onArchive={() => setConfirmDialog('archive')}
                onUnarchive={unarchive}
                onDelete={() => setConfirmDialog('delete')}
                onDuplicate={() => console.log('TODO: duplicate')}
                onMove={() => console.log('TODO: move')}
            />

            <NoteContent editor={editor} isEditing={isEditing} />

            <Box mt={0}>
                <p>
                    <strong>🧠 Уровень запоминания:</strong> {note.memoryLevel} из 100
                </p>
            </Box>

            <NoteActions
                isEditing={isEditing}
                noteArchived={!!note.archived}
                isSaving={isSaving}
                isUnarchiving={isUnarchiving}
                onSave={handleSubmit}
                onCancelEdit={() => setEditing(false)}
                onDelete={() => setConfirmDialog('delete')}
                onArchive={() => setConfirmDialog('archive')}
                onUnarchive={unarchive}
            />
        </Box>
    )
}

export default NoteDetailPage
