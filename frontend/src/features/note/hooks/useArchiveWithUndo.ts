import { useCallback, useEffect, useRef, useState } from 'react'
import { useArchiveNoteMutation, Note } from '@/features/note/noteApi'

const ACTION_QUEUE_KEY = 'note_action_queue'
const SNACKBAR_TIMEOUT_MS = 4000

type PendingAction = {
    type: 'archive'
    noteId: string
    createdAt: number
}

type Options = {
    onCommitted?: () => void // вызываем после успешной архивации (для refetch)
}

export function useArchiveWithUndo(options?: Options) {
    const onCommitted = options?.onCommitted
    const [archiveNote] = useArchiveNoteMutation()

    const [optimisticArchivedIds, setOptimisticArchivedIds] = useState<Set<string>>(new Set())
    const [snackbarOpen, setSnackbarOpen] = useState(false)
    const [lastArchivedNote, setLastArchivedNote] = useState<Note | null>(null)
    const undoTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null)

    // ---- localStorage queue ----
    const getQueue = useCallback((): PendingAction[] => {
        try {
            return JSON.parse(localStorage.getItem(ACTION_QUEUE_KEY) || '[]') as PendingAction[]
        } catch {
            return []
        }
    }, [])

    const setQueue = useCallback((q: PendingAction[]) => {
        localStorage.setItem(ACTION_QUEUE_KEY, JSON.stringify(q))
    }, [])

    const pushToQueue = useCallback((action: PendingAction) => {
        const q = getQueue()
        q.push(action)
        setQueue(q)
    }, [getQueue, setQueue])

    const processQueue = useCallback(async () => {
        const q = getQueue()
        if (q.length === 0) return

        const remaining: PendingAction[] = []
        for (const action of q) {
            try {
                if (action.type === 'archive') {
                    await archiveNote(action.noteId).unwrap()
                }
                // успех — не добавляем в remaining
            } catch (err) {
                if (!navigator.onLine) {
                    remaining.push(action, ...q.slice(q.indexOf(action) + 1))
                    break
                } else {
                    remaining.push(action)
                }
            }
        }
        setQueue(remaining)
        if (remaining.length === 0) {
            onCommitted?.()
        }
    }, [archiveNote, getQueue, setQueue, onCommitted])

    useEffect(() => {
        const onOnline = () => {
            processQueue().catch((e) => console.error('processQueue error', e))
        }
        window.addEventListener('online', onOnline)
        return () => window.removeEventListener('online', onOnline)
    }, [processQueue])

    // ---- API для дочерних компонентов ----
    const handleRequestArchive = useCallback((note: Note) => {
        // оптимистично скрываем
        setOptimisticArchivedIds((prev) => {
            const next = new Set(prev)
            next.add(note.id)
            return next
        })
        setLastArchivedNote(note)
        setSnackbarOpen(true)

        if (undoTimerRef.current) {
            clearTimeout(undoTimerRef.current)
            undoTimerRef.current = null
        }

        // через таймаут пытаемся зафиксировать действие
        undoTimerRef.current = setTimeout(async () => {
            try {
                await archiveNote(note.id).unwrap()
                onCommitted?.()
            } catch {
                // офлайн/ошибка — в очередь
                pushToQueue({
                    type: 'archive',
                    noteId: note.id,
                    createdAt: Date.now(),
                })
            } finally {
                setOptimisticArchivedIds((prev) => {
                    const next = new Set(prev)
                    next.delete(note.id)
                    return next
                })
                setLastArchivedNote(null)
                setSnackbarOpen(false)
                undoTimerRef.current = null
            }
        }, SNACKBAR_TIMEOUT_MS)
    }, [archiveNote, onCommitted, pushToQueue])

    const handleUndo = useCallback(() => {
        if (undoTimerRef.current) {
            clearTimeout(undoTimerRef.current)
            undoTimerRef.current = null
        }
        setOptimisticArchivedIds((prev) => {
            const next = new Set(prev)
            if (lastArchivedNote) next.delete(lastArchivedNote.id)
            return next
        })
        setLastArchivedNote(null)
        setSnackbarOpen(false)
    }, [lastArchivedNote])

    const handleSnackbarClose = useCallback((_ev?: any, reason?: string) => {
        if (reason === 'clickaway') return
        setSnackbarOpen(false)
    }, [])

    return {
        optimisticArchivedIds,
        snackbarOpen,
        handleSnackbarClose,
        handleRequestArchive,
        handleUndo,
        autoHideDuration: SNACKBAR_TIMEOUT_MS,
    }
}