import { useNavigate } from 'react-router-dom'
import {
    useUpdateNoteMutation,
    useDeleteNoteMutation,
    useArchiveNoteMutation,
    useUnarchiveNoteMutation,
} from '@/features/note/noteApi'

export function useNoteActions(noteId: string) {
    const navigate = useNavigate()
    const [updateNote, { isLoading: isSaving }] = useUpdateNoteMutation()
    const [deleteNote] = useDeleteNoteMutation()
    const [archiveNote] = useArchiveNoteMutation()
    const [unarchiveNote, { isLoading: isUnarchiving }] = useUnarchiveNoteMutation()

    const update = async (title: string, content: string) => {
        return updateNote({ id: noteId, title, content }).unwrap()
    }

    const remove = async () => {
        await deleteNote(noteId).unwrap()
        navigate('/notes')
    }

    const archive = async () => {
        await archiveNote(noteId).unwrap()
        navigate('/notes')
    }

    const unarchive = async () => {
        await unarchiveNote(noteId).unwrap()
        navigate('/notes')
    }

    return { update, remove, archive, unarchive, isSaving, isUnarchiving }
}
