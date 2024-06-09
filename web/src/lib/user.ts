interface User {
  id: number
  username: string
  event_types?: EventType[]
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
  is_google_available: boolean
  is_microsoft_available: boolean
}
