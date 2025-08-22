import { useCallback, useRef } from 'react'

interface Options {
    onLongPress: (e: React.MouseEvent | React.TouchEvent, target: EventTarget & HTMLElement) => void
    delay?: number
}

export const useLongPress = ({ onLongPress, delay = 500 }: Options) => {
    const timerRef = useRef<number>()

    const start = useCallback(
        (e: React.MouseEvent | React.TouchEvent) => {
            const target = e.currentTarget as HTMLElement
            e.persist?.()
            timerRef.current = window.setTimeout(() => {
                onLongPress(e, target) // теперь всегда даём target
            }, delay)
        },
        [onLongPress, delay]
    )

    const clear = useCallback(() => {
        if (timerRef.current) {
            clearTimeout(timerRef.current)
            timerRef.current = undefined
        }
    }, [])

    return {
        onMouseDown: start,
        onTouchStart: start,
        onMouseUp: clear,
        onMouseLeave: clear,
        onTouchEnd: clear,
    }
}
