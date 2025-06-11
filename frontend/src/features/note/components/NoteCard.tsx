// NoteCard.tsx
import { Card, CardContent, Typography } from '@mui/material'
import { Note } from '../noteApi'
import dayjs from 'dayjs'

const NoteCard = ({ note }: { note: Note }) => (
    <Card>
        <CardContent>
            <Typography variant="h6">{note.title}</Typography>
            <Typography variant="body2" color="text.secondary">
                Создана: {dayjs(note.created_at).format('D MMMM YYYY')}
            </Typography>
        </CardContent>
    </Card>
)

export default NoteCard
