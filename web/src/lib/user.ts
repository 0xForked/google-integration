interface User {
    id: number,
    username: string
}

interface Availability {
    id: number
    label: string
    timezone: string
    days: AvailabilityDays[]
}

interface AvailabilityDays {
    id: number
    enable: number
    day: number
    start_time: number
    end_time: number
}

interface EventType {
    id: number
    enable: number
    title: string
    description: string
    duration: number
    availability: Availability
}